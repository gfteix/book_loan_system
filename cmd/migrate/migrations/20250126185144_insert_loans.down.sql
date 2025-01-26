DELETE FROM loans WHERE expiring_date = CURRENT_DATE + INTERVAL '1 day';
