# RFC-001: TODO CRUD API Specification

| Campo         | Valor                          |
|---------------|--------------------------------|
| **Título**    | TODO CRUD API                  |
| **Estado**    | Draft                          |
| **Módulo**    | `todo-api`                     |
| **Autor**     | Sebastian                      |
| **Fecha**     | 2026-01-27                     |

---

## 1. Resumen

Este documento especifica la implementación de una API REST para gestión de tareas (TODOs) con operaciones CRUD completas: Create, Read, Update y Delete.

---

## 2. Motivación

Crear una API simple y funcional para practicar conceptos de:
- Arquitectura en capas (Controller → Service → Domain)
- Manejo de requests/responses HTTP
- Validación de datos
- Patrones de diseño en Go

---

## 3. Arquitectura

### 3.1 Estructura de Carpetas

```
todo-api/
├── cmd/
│   └── main.go                 # Punto de entrada
├── pkg/
│   └── todo/
│       ├── domain/
│       │   └── todo.go         # Entidad y tipos
│       ├── controller/
│       │   └── controller.go   # Handlers HTTP
│       ├── service/
│       │   └── service.go      # Lógica de negocio
│       └── routes.go           # Registro de rutas
├── web/                        # Abstracciones HTTP (existente)
└── boot/                       # Inicialización (existente)
```

### 3.2 Flujo de Datos

```
┌─────────────────────────────────────────────────────────────────┐
│                         HTTP Request                            │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      CONTROLLER LAYER                           │
│  • Parsea request (body, params, query)                         │
│  • Valida input                                                 │
│  • Llama al service                                             │
│  • Formatea response                                            │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                       SERVICE LAYER                             │
│  • Lógica de negocio                                            │
│  • Almacenamiento de datos (in-memory)                          │
│  • Generación de IDs                                            │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                       DOMAIN LAYER                              │
│  • Entidades (Todo)                                             │
│  • Value Objects (Status, Priority)                             │
└─────────────────────────────────────────────────────────────────┘
```

---

## 4. Modelo de Dominio

### 4.1 Entidad: Todo

```go
type Todo struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Description string     `json:"description,omitempty"`
    Status      Status     `json:"status"`
    Priority    Priority   `json:"priority"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}
```

### 4.2 Enumeraciones

#### Status (Estado de la tarea)

| Valor         | Descripción                    |
|---------------|--------------------------------|
| `pending`     | Tarea pendiente (por defecto)  |
| `in_progress` | Tarea en progreso              |
| `completed`   | Tarea completada               |

```go
type Status string

const (
    StatusPending    Status = "pending"
    StatusInProgress Status = "in_progress"
    StatusCompleted  Status = "completed"
)
```

#### Priority (Prioridad)

| Valor    | Descripción          |
|----------|----------------------|
| `low`    | Prioridad baja       |
| `medium` | Prioridad media      |
| `high`   | Prioridad alta       |

```go
type Priority string

const (
    PriorityLow    Priority = "low"
    PriorityMedium Priority = "medium"
    PriorityHigh   Priority = "high"
)
```

### 4.3 Validaciones del Dominio

| Campo         | Reglas                                              |
|---------------|-----------------------------------------------------|
| `id`          | UUID v4, generado automáticamente                   |
| `title`       | Requerido, 1-100 caracteres                         |
| `description` | Opcional, máximo 500 caracteres                     |
| `status`      | Debe ser valor válido del enum                      |
| `priority`    | Debe ser valor válido del enum                      |
| `created_at`  | Generado automáticamente al crear                   |
| `updated_at`  | Actualizado automáticamente en cada modificación    |

---

## 5. API Endpoints

### 5.1 Base URL

```
http://localhost:8080/api/v1
```

### 5.2 Resumen de Endpoints

| Método   | Endpoint              | Descripción              |
|----------|-----------------------|--------------------------|
| `GET`    | `/todos`              | Listar todos los TODOs   |
| `GET`    | `/todos/:id`          | Obtener un TODO por ID   |
| `POST`   | `/todos`              | Crear un nuevo TODO      |
| `PATCH`  | `/todos/:id`          | Actualizar un TODO       |
| `DELETE` | `/todos/:id`          | Eliminar un TODO         |

---

## 6. Especificación de Endpoints

### 6.1 GET /todos - Listar TODOs

**Descripción**: Retorna la lista de todos los TODOs.

**Query Parameters** (opcionales):

| Parámetro  | Tipo     | Descripción                          |
|------------|----------|--------------------------------------|
| `status`   | string   | Filtrar por estado                   |
| `priority` | string   | Filtrar por prioridad                |

**Request**:
```http
GET /api/v1/todos?status=pending&priority=high
```

**Response 200 OK**:
```json
{
    "data": [
        {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "title": "Comprar leche",
            "description": "Ir al supermercado",
            "status": "pending",
            "priority": "high",
            "created_at": "2026-01-27T10:00:00Z",
            "updated_at": "2026-01-27T10:00:00Z"
        }
    ],
    "total": 1
}
```

---

### 6.2 GET /todos/:id - Obtener TODO por ID

**Descripción**: Retorna un TODO específico.

**Path Parameters**:

| Parámetro | Tipo   | Descripción    |
|-----------|--------|----------------|
| `id`      | string | UUID del TODO  |

**Request**:
```http
GET /api/v1/todos/550e8400-e29b-41d4-a716-446655440000
```

**Response 200 OK**:
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Comprar leche",
    "description": "Ir al supermercado",
    "status": "pending",
    "priority": "high",
    "created_at": "2026-01-27T10:00:00Z",
    "updated_at": "2026-01-27T10:00:00Z"
}
```

**Response 404 Not Found**:
```json
{
    "status": 404,
    "error": "Not Found",
    "message": "todo not found",
    "causes": []
}
```

---

### 6.3 POST /todos - Crear TODO

**Descripción**: Crea un nuevo TODO.

**Request Body**:

| Campo         | Tipo   | Requerido | Descripción                        |
|---------------|--------|-----------|-------------------------------------|
| `title`       | string | Sí        | Título de la tarea (1-100 chars)   |
| `description` | string | No        | Descripción (max 500 chars)        |
| `priority`    | string | No        | `low`, `medium`, `high` (default: `medium`) |

**Request**:
```http
POST /api/v1/todos
Content-Type: application/json

{
    "title": "Comprar leche",
    "description": "Ir al supermercado",
    "priority": "high"
}
```

**Response 201 Created**:
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Comprar leche",
    "description": "Ir al supermercado",
    "status": "pending",
    "priority": "high",
    "created_at": "2026-01-27T10:00:00Z",
    "updated_at": "2026-01-27T10:00:00Z"
}
```

**Response 400 Bad Request**:
```json
{
    "status": 400,
    "error": "Bad Request",
    "message": "validation failed",
    "causes": [
        "title is required",
        "title must be between 1 and 100 characters"
    ]
}
```

---

### 6.4 PATCH /todos/:id - Actualizar TODO

**Descripción**: Actualiza parcialmente un TODO existente.

**Path Parameters**:

| Parámetro | Tipo   | Descripción    |
|-----------|--------|----------------|
| `id`      | string | UUID del TODO  |

**Request Body** (todos los campos son opcionales):

| Campo         | Tipo   | Descripción                        |
|---------------|--------|------------------------------------|
| `title`       | string | Nuevo título (1-100 chars)         |
| `description` | string | Nueva descripción (max 500 chars)  |
| `status`      | string | `pending`, `in_progress`, `completed` |
| `priority`    | string | `low`, `medium`, `high`            |

**Request**:
```http
PATCH /api/v1/todos/550e8400-e29b-41d4-a716-446655440000
Content-Type: application/json

{
    "status": "completed"
}
```

**Response 200 OK**:
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Comprar leche",
    "description": "Ir al supermercado",
    "status": "completed",
    "priority": "high",
    "created_at": "2026-01-27T10:00:00Z",
    "updated_at": "2026-01-27T12:30:00Z"
}
```

**Response 404 Not Found**:
```json
{
    "status": 404,
    "error": "Not Found",
    "message": "todo not found",
    "causes": []
}
```

**Response 400 Bad Request**:
```json
{
    "status": 400,
    "error": "Bad Request",
    "message": "validation failed",
    "causes": [
        "invalid status value: must be pending, in_progress, or completed"
    ]
}
```

---

### 6.5 DELETE /todos/:id - Eliminar TODO

**Descripción**: Elimina un TODO existente.

**Path Parameters**:

| Parámetro | Tipo   | Descripción    |
|-----------|--------|----------------|
| `id`      | string | UUID del TODO  |

**Request**:
```http
DELETE /api/v1/todos/550e8400-e29b-41d4-a716-446655440000
```

**Response 204 No Content**:
```
(sin body)
```

**Response 404 Not Found**:
```json
{
    "status": 404,
    "error": "Not Found",
    "message": "todo not found",
    "causes": []
}
```

---

## 7. Códigos de Error

| Código | Nombre                | Descripción                           |
|--------|-----------------------|---------------------------------------|
| 200    | OK                    | Operación exitosa                     |
| 201    | Created               | Recurso creado exitosamente           |
| 204    | No Content            | Eliminación exitosa                   |
| 400    | Bad Request           | Error de validación en el request     |
| 404    | Not Found             | Recurso no encontrado                 |
| 500    | Internal Server Error | Error interno del servidor            |

---

## 8. Estructura de Errores

Todos los errores siguen el formato:

```json
{
    "status": <http_status_code>,
    "error": "<http_status_text>",
    "message": "<descripción_del_error>",
    "causes": ["<causa_1>", "<causa_2>"]
}
```

---

## 9. DTOs (Data Transfer Objects)

### 9.1 Request DTOs

#### CreateTodoRequest

```go
type CreateTodoRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Priority    string `json:"priority"`
}
```

#### UpdateTodoRequest

```go
type UpdateTodoRequest struct {
    Title       *string `json:"title,omitempty"`
    Description *string `json:"description,omitempty"`
    Status      *string `json:"status,omitempty"`
    Priority    *string `json:"priority,omitempty"`
}
```

### 9.2 Response DTOs

#### TodoResponse

```go
type TodoResponse struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description,omitempty"`
    Status      string `json:"status"`
    Priority    string `json:"priority"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}
```

#### ListTodosResponse

```go
type ListTodosResponse struct {
    Data  []TodoResponse `json:"data"`
    Total int            `json:"total"`
}
```

---

## 10. Interfaces del Service

```go
type TodoService interface {
    // List retorna todos los TODOs, opcionalmente filtrados
    List(ctx context.Context, filters TodoFilters) ([]Todo, error)

    // GetByID retorna un TODO por su ID
    GetByID(ctx context.Context, id string) (Todo, error)

    // Create crea un nuevo TODO
    Create(ctx context.Context, req CreateTodoRequest) (Todo, error)

    // Update actualiza un TODO existente
    Update(ctx context.Context, id string, req UpdateTodoRequest) (Todo, error)

    // Delete elimina un TODO por su ID
    Delete(ctx context.Context, id string) error
}

type TodoFilters struct {
    Status   *Status
    Priority *Priority
}
```

---

## 11. Errores del Dominio

```go
var (
    // ErrTodoNotFound se retorna cuando no se encuentra el TODO
    ErrTodoNotFound = errors.New("todo not found")

    // ErrInvalidTitle se retorna cuando el título es inválido
    ErrInvalidTitle = errors.New("invalid title")

    // ErrInvalidStatus se retorna cuando el status es inválido
    ErrInvalidStatus = errors.New("invalid status")

    // ErrInvalidPriority se retorna cuando la prioridad es inválida
    ErrInvalidPriority = errors.New("invalid priority")
)
```

---

## 12. Plan de Implementación

### Fase 1: Domain Layer
1. Crear `pkg/todo/domain/todo.go` con entidades y enums
2. Crear `pkg/todo/domain/errors.go` con errores del dominio

### Fase 2: Service Layer
3. Crear `pkg/todo/service/service.go` con interfaz e implementación in-memory

### Fase 3: Controller Layer
4. Crear `pkg/todo/controller/controller.go` con handlers HTTP
5. Crear `pkg/todo/controller/dto.go` con request/response DTOs

### Fase 4: Routes
6. Crear `pkg/todo/routes.go` para registrar las rutas
7. Actualizar `cmd/main.go` para montar las rutas

### Fase 5: Testing
8. Probar endpoints con curl/Postman

---

## 13. Ejemplos de Uso con cURL

```bash
# Crear TODO
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Mi primera tarea", "priority": "high"}'

# Listar TODOs
curl http://localhost:8080/api/v1/todos

# Obtener TODO por ID
curl http://localhost:8080/api/v1/todos/550e8400-e29b-41d4-a716-446655440000

# Actualizar TODO
curl -X PATCH http://localhost:8080/api/v1/todos/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'

# Eliminar TODO
curl -X DELETE http://localhost:8080/api/v1/todos/550e8400-e29b-41d4-a716-446655440000
```

---

## 14. Consideraciones Futuras

- [ ] Persistencia en PostgreSQL
- [ ] Paginación en listado
- [ ] Ordenamiento por fecha/prioridad
- [ ] Búsqueda por texto en título/descripción
- [ ] Autenticación/Autorización
- [ ] Rate limiting
- [ ] Logging estructurado
- [ ] Métricas y observabilidad

---

## Apéndice A: Diagrama de Secuencia - Crear TODO

```
┌──────┐          ┌────────────┐          ┌─────────┐          ┌────────┐
│Client│          │ Controller │          │ Service │          │ Domain │
└──┬───┘          └─────┬──────┘          └────┬────┘          └───┬────┘
   │                    │                      │                   │
   │ POST /todos        │                      │                   │
   │ {title, priority}  │                      │                   │
   │───────────────────>│                      │                   │
   │                    │                      │                   │
   │                    │ Parse & Validate     │                   │
   │                    │─────────────────────>│                   │
   │                    │                      │                   │
   │                    │                      │ Create Todo       │
   │                    │                      │──────────────────>│
   │                    │                      │                   │
   │                    │                      │    Todo Entity    │
   │                    │                      │<──────────────────│
   │                    │                      │                   │
   │                    │     Todo Created     │                   │
   │                    │<─────────────────────│                   │
   │                    │                      │                   │
   │  201 Created       │                      │                   │
   │  {todo object}     │                      │                   │
   │<───────────────────│                      │                   │
   │                    │                      │                   │
```

---

**Fin del documento**
