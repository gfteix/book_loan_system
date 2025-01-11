package loans

import (
	"database/sql"

	"github.com/gfteix/book_loan_system/types"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateLoan(loan types.Loan) error {
	id := uuid.NewString()
	_, err := r.db.Exec("INSERT INTO loans (id, user_id, book_item_id, status, expiring_date, return_date, loan_date, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		id, loan.UserId, loan.BookItemId, loan.Status, loan.ExpiringDate, loan.ReturnDate, loan.LoanDate, loan.CreatedAt)
	if err != nil {
		return err
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
