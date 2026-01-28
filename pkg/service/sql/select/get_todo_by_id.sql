SELECT id, title, description, status, priority, created_at, updated_at
FROM todos
WHERE id = $1;
