SELECT id, title, description, status, priority, created_at, updated_at
FROM todos
WHERE ($1::VARCHAR IS NULL OR status = $1)
  AND ($2::VARCHAR IS NULL OR priority = $2)
ORDER BY created_at DESC;
