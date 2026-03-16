# Aura ERP - Go Backend

A RESTful API backend for the Aura ERP system, built with Go and Gin framework.

## Features

- **RESTful API** with Gin web framework
- **PostgreSQL** database integration
- **JWT Authentication** with custom token implementation
- **CORS** support for frontend integration
- **Complete CRUD operations** for all resources
- **Nested routes** for items (proposal items, order items)
- **Audit logging** for entity changes
- **Connection pooling** for efficient database access

## Project Structure

```
backend-go/
├── config/
│   └── database.go           # Database initialization and connection pooling
├── controllers/
│   ├── auth_controller.go    # Authentication endpoints
│   ├── user_controller.go    # User management endpoints
│   ├── client_controller.go  # Client management endpoints
│   ├── product_controller.go # Product management endpoints
│   ├── section_controller.go # Section management endpoints
│   ├── proposal_controller.go # Proposal & proposal items endpoints
│   ├── order_controller.go   # Order & order items endpoints
│   └── audit_controller.go   # Audit log endpoints
├── models/
│   ├── user.go               # User data models
│   ├── client.go             # Client data models
│   ├── product.go            # Product data models
│   ├── section.go            # Section data models
│   ├── proposal.go           # Proposal & ProposalItem data models
│   ├── order.go              # Order & OrderItem data models
│   └── audit.go              # AuditLog data models
├── services/
│   ├── auth_service.go       # Authentication logic & token management
│   ├── user_service.go       # User business logic
│   ├── client_service.go     # Client business logic
│   ├── product_service.go    # Product business logic
│   ├── section_service.go    # Section business logic
│   ├── proposal_service.go   # Proposal business logic
│   ├── order_service.go      # Order business logic
│   └── audit_service.go      # Audit logging logic
├── routes/
│   └── routes.go             # API route definitions
├── main.go                   # Application entry point
├── go.mod                    # Go module definition
├── go.sum                    # Go module dependencies
├── Dockerfile                # Docker image configuration
├── .env.example              # Example environment variables
└── README.md                 # This file
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Docker & Docker Compose (optional)

## Environment Variables

Create a `.env` file based on `.env.example`:

```bash
DATABASE_URL=postgres://myuser:mypassword@localhost:5433/crud_db?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
PORT=5000
```

## Installation & Setup

### Option 1: Local Development

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Ensure PostgreSQL is running:**
   The database should be accessible at the `DATABASE_URL` specified in your `.env` file.

3. **Initialize the database:**
   Run the SQL script from the project root:
   ```bash
   psql -U myuser -d crud_db -f init.sql
   ```

4. **Run the server:**
   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:5000`

### Option 2: Docker Compose

From the project root directory:

```bash
docker-compose up
```

This will:
- Start PostgreSQL
- Initialize the database with `init.sql`
- Build and run the Go backend
- Build and run the frontend
- Start pgAdmin for database management

Access the services:
- **Backend API:** http://localhost:5000
- **Frontend:** http://localhost:5173
- **pgAdmin:** http://localhost:5050

## Build for Production

```bash
go build -o aura-erp-backend main.go
```

Or with Docker:

```bash
docker build -t aura-erp-backend:latest -f Dockerfile .
```

## API Endpoints

### Health Check
- `GET /api/health` - Server health status

### Authentication
- `POST /api/auth/login` - User login
- `GET /api/auth/verify` - Verify JWT token
- `POST /api/auth/logout` - User logout

### Users
- `GET /api/users` - List all users
- `GET /api/users/:id` - Get user by ID
- `POST /api/users` - Create new user
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

### Clients
- `GET /api/clients` - List all clients
- `GET /api/clients/:id` - Get client by ID
- `POST /api/clients` - Create new client
- `PUT /api/clients/:id` - Update client
- `DELETE /api/clients/:id` - Delete client

### Products
- `GET /api/products` - List all products
- `GET /api/products/:id` - Get product by ID
- `POST /api/products` - Create new product
- `PUT /api/products/:id` - Update product
- `DELETE /api/products/:id` - Delete product

### Sections
- `GET /api/sections` - List all sections
- `GET /api/sections/:id` - Get section by ID
- `POST /api/sections` - Create new section
- `PUT /api/sections/:id` - Update section
- `DELETE /api/sections/:id` - Delete section

### Proposals
- `GET /api/proposals` - List all proposals
- `GET /api/proposals/:id` - Get proposal with items
- `POST /api/proposals` - Create new proposal
- `PUT /api/proposals/:id` - Update proposal
- `DELETE /api/proposals/:id` - Delete proposal

**Proposal Items** (nested):
- `GET /api/proposals/:proposalId/items` - Get proposal items
- `POST /api/proposals/:proposalId/items` - Create proposal item
- `PUT /api/proposals/items/:id` - Update proposal item
- `DELETE /api/proposals/items/:id` - Delete proposal item

### Orders
- `GET /api/orders` - List all orders
- `GET /api/orders/:id` - Get order with items
- `POST /api/orders` - Create new order
- `PUT /api/orders/:id` - Update order
- `DELETE /api/orders/:id` - Delete order

**Order Items** (nested):
- `GET /api/orders/:orderId/items` - Get order items
- `POST /api/orders/:orderId/items` - Create order item
- `PUT /api/orders/items/:id` - Update order item
- `DELETE /api/orders/items/:id` - Delete order item

### Audit Logs
- `GET /api/audit-log` - List audit logs (with pagination)
  - Query params: `limit` (default: 50), `offset` (default: 0)
- `GET /api/audit-log/:entityType/:entityId` - Get logs for specific entity
- `POST /api/audit-log` - Create audit log entry

## Authentication

The API uses JWT-based authentication. After login:

```bash
curl -X POST http://localhost:5000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```

Response:
```json
{
  "token": "eyJ...",
  "user": {
    "id": 1,
    "name": "User Name",
    "email": "user@example.com",
    "role": "admin"
  }
}
```

Use the token in subsequent requests:

```bash
curl -X GET http://localhost:5000/api/users \
  -H "Authorization: Bearer eyJ..."
```

## Database Schema

The PostgreSQL database includes the following tables:
- `users` - System users
- `clients` - Client information
- `sections` - Work sections
- `products` - Product catalog
- `proposals` - Sales proposals
- `proposal_items` - Items in proposals
- `orders` - Customer orders
- `order_items` - Items in orders
- `audit_log` - Change history for all entities

See `init.sql` in the project root for the complete schema.

## Error Handling

All endpoints return standard HTTP status codes:
- `200 OK` - Successful GET/PUT
- `201 Created` - Successful POST
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Missing/invalid token
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

Error responses include a JSON body:
```json
{
  "error": "Descriptive error message"
}
```

## CORS Configuration

CORS is enabled for all origins. Modify the `corsMiddleware()` function in `main.go` to restrict origins in production:

```go
c.Writer.Header().Set("Access-Control-Allow-Origin", "https://yourdomain.com")
```

## Database Connection Pooling

The connection pool is configured with:
- Max open connections: 25
- Max idle connections: 5

Adjust in `config/database.go` if needed for your load.

## Development Notes

### Adding a New Resource

1. Create model file in `models/` (e.g., `models/newresource.go`)
2. Create service file in `services/` (e.g., `services/newresource_service.go`)
3. Create controller file in `controllers/` (e.g., `controllers/newresource_controller.go`)
4. Add routes in `routes/routes.go`
5. Ensure database tables exist in PostgreSQL

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
goimports -w .
```

### Linting

```bash
golangci-lint run ./...
```

## Performance Considerations

- Connection pooling is enabled for database efficiency
- SQL queries use parameterized statements to prevent SQL injection
- Indexes should be added to frequently queried fields (see `init.sql`)
- Consider caching for read-heavy operations

## Security Notes

1. **JWT Secret:** Change `JWT_SECRET` environment variable in production
2. **CORS:** Restrict origins to your frontend domain
3. **Password Hashing:** Currently uses SHA256; consider using bcrypt for stronger security
4. **Database Credentials:** Never commit `.env` files; use environment variables
5. **HTTPS:** Ensure the API is served over HTTPS in production

## Troubleshooting

### Database Connection Error
- Ensure PostgreSQL is running
- Verify `DATABASE_URL` is correct
- Check database credentials

### Port Already in Use
- Change the `PORT` environment variable
- Or kill the process using the port: `lsof -i :5000`

### CORS Errors
- Verify frontend URL is allowed in CORS configuration
- Check `Access-Control-Allow-Origin` header

## Deployment

### Docker

```bash
docker build -t aura-erp-backend .
docker run -e DATABASE_URL="..." -e JWT_SECRET="..." -p 5000:5000 aura-erp-backend
```

### Kubernetes

Create appropriate deployment manifests with:
- ConfigMap for `JWT_SECRET`
- Secret for `DATABASE_URL`
- Service for port exposure
- Probes for health checking

## License

ISC License - See LICENSE file

## Support

For issues and questions, please refer to the project documentation or contact the development team.
