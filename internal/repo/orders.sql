-- name: CreateOrder :exec
INSERT INTO orders (
        guid,
        user_guid,
        number,
        amount,
        status,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, now(), now());