INSERT INTO todos (title, description, status, priority)
VALUES ($1, $2, $3, $4)
RETURNING id, title, description, status, priority, created_at, updated_at;
