CREATE TABLE invoices(
    user_id VARCHAR(256) NOT NULL,
    file_id VARCHAR(256) NOT NULL,
    customer_id INT NOT NULL,
    period_start DATE NOT NULL,
    paid_plan VARCHAR(10) NOT NULL,
    paid_amount REAL NOT NULL,
    period_end DATE NOT NULL
);