-- name: CreateApplication :one
INSERT INTO loan_applications(
  user_id,
  type,
  vehicle_vin,
  vehicle_name,
  currency_code,
  price,
  down_payment,
  net_price,
  margin_rate,
  term_months,
  monthly_payment,
  status
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12  
) RETURNING *;

-- name: GetApplication :one
select *
from loan_applications
where id = $1
;

-- name: CountApplicationsByUser :one
select count(*)
from loan_applications
where user_id = $1
;

-- name: ListApplicationsByUser :many
select *
from loan_applications
where user_id = $1
order by id desc
limit $2
offset $3
;
