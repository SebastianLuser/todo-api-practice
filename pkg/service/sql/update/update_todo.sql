UPDATE todos
SET
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    status = COALESCE($4, status),
    priority = COALESCE($5, priority),
    updated_at = NOW()
WHERE id = $1
RETURNING id, title, description, status, priority, created_at, updated_at;
