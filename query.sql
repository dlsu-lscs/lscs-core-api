-- name: GetMember :one
SELECT * FROM members WHERE email = ?;

-- name: ListMembers :many
SELECT * FROM members ORDER BY email;

-- name: CheckEmailIfMember :one
SELECT email FROM members WHERE email = ?;
