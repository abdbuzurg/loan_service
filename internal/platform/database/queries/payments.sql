-- name: CreatePayment :one
INSERT INTO payments(
  loan_id,
  currency_code,
  payment_date,
  amount,
  method,
  status,
  transaction_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $6
) RETURNING *;

-- name: CountPayments :one
select count(*)
from payments
where loan_id = $1
;

-- name: ListPaymentsByLoan :many
select *
from payments
where loan_id = $1
order by id desc
limit $2
offset $3
;
