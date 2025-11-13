-- name: GetLoan :one
select *
from loans
where id = $1
;

-- name: CountLoansByUser :one
select count(*)
from loans
where user_id = $1
;

-- name: ListLoansByUser :many
select *
from loans
where user_id = $1
order by id desc
limit $2
offset $3
;
