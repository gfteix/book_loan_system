package loans

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gfteix/book_loan_system/types"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetBookItemById(ctx context.Context, tx *sql.Tx, id string) (*types.BookItem, error) {
	rows, err := tx.QueryContext(ctx, "SELECT id, status FROM book_items WHERE Id = $1 FOR UPDATE", id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		bookItem := new(types.BookItem)

		err := rows.Scan(
			&bookItem.Id,
			&bookItem.Status,
		)

		if err != nil {
			return nil, err
		}

		return bookItem, nil
	}

	return nil, nil
}

func (r *Repository) CreateLoan(ctx context.Context, loan types.Loan) error {
	fail := func(tx *sql.Tx, err error) error {
		fmt.Printf("transaction failure %v", err)

		er := tx.Rollback()

		if er != nil {
			fmt.Printf("rollback fail %v", er)
		}

		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		log.Printf("error while starting transaction %v", err)
		return err
	}

	bookItem, err := r.GetBookItemById(ctx, tx, loan.BookItemId)

	if err != nil {
		log.Printf("error while getting book item %v", err)
		return err
	}

	if bookItem == nil {
		return fmt.Errorf("no related book item found for %v", loan.BookItemId)
	}

	_, err = tx.ExecContext(ctx, "UPDATE book_items SET status = lent WHERE id = $1", loan.BookItemId)

	if err != nil {
		log.Printf("error while updating book item %v", err)
		return fail(tx, err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO loans (id, user_id, book_item_id, status, expiring_date, return_date, loan_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		uuid.NewString(), loan.UserId, loan.BookItemId, loan.Status, loan.ExpiringDate, loan.ReturnDate, loan.LoanDate)

	if err != nil {
		log.Printf("error while creating loan %v", err)
		return fail(tx, err)
	}

	return nil
}

func (r *Repository) GetLoan(id string) (*types.Loan, error) {
	rows, err := r.db.Query("SELECT id, user_id, book_item_id, status, expiring_date, return_date, loan_date, created_at FROM loans WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loan *types.Loan

	for rows.Next() {
		loan, err = scanRowIntoLoan(rows)
		if err != nil {
			return nil, err
		}
	}

	return loan, nil
}

func (r *Repository) GetLoans() ([]types.Loan, error) {
	rows, err := r.db.Query("SELECT id, user_id, book_item_id, status, expiring_date, return_date, loan_date, created_at FROM loans")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []types.Loan
	for rows.Next() {
		loan, err := scanRowIntoLoan(rows)
		if err != nil {
			return nil, err
		}
		loans = append(loans, *loan)
	}

	return loans, nil
}

func scanRowIntoLoan(rows *sql.Rows) (*types.Loan, error) {
	loan := new(types.Loan)
	err := rows.Scan(
		&loan.Id,
		&loan.UserId,
		&loan.BookItemId,
		&loan.Status,
		&loan.ExpiringDate,
		&loan.ReturnDate,
		&loan.LoanDate,
		&loan.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return loan, nil
}
