# Hexagonal Architecture Implementation Guide for go-invoice

## Overview

This guide provides a step-by-step curriculum for transforming the current simple layered go-invoice application into a proper Hexagonal Architecture following the Ports and Adapters pattern.

## Learning Resources

### Essential Reading (Complete before implementation)

1. **[Hexagonal Architecture - Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture)** 
   - The foundational article explaining Ports and Adapters
   - Key concepts: Inside-Outside asymmetry, Ports, Adapters
   - **Time**: 2-3 hours

2. **[How to implement Clean Architecture in Go - Three Dots Labs](https://threedots.tech/post/introducing-clean-architecture/)**
   - Practical Go implementation patterns
   - Folder structure and dependency injection
   - **Time**: 1-2 hours

3. **[Building RESTful API with Hexagonal Architecture in Go](https://dev.to/bagashiz/building-restful-api-with-hexagonal-architecture-in-go-1mij)**
   - Step-by-step Go implementation guide
   - Real-world examples
   - **Time**: 1-2 hours

### Key Concepts to Master

- **Ports**: Interfaces that define how your application communicates
- **Adapters**: Concrete implementations of ports
- **Dependency Inversion**: Inner layers depend on interfaces, outer layers implement them
- **Testability**: Core logic testable without external dependencies

## Target Architecture

```
Current: HTTP -> Service -> DB
Target:  Driving Adapters -> Core (Domain + Ports) -> Driven Adapters
```

### Target Folder Structure

```
/home/william/code/go-invoice/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point & DI
├── internal/
│   ├── domain/                        # Pure business logic
│   │   ├── invoice.go                 # Core domain entities
│   │   ├── line_item.go               # Value objects
│   │   └── tax.go                     # Business rules
│   ├── ports/
│   │   ├── inbound/
│   │   │   └── invoice_service.go     # Service interfaces (primary ports)
│   │   └── outbound/
│   │       ├── invoice_repo.go        # Repository interfaces (secondary ports)
│   │       ├── pdf_generator.go       # PDF generation interface
│   │       └── event_publisher.go     # Event publishing interface
│   └── adapters/
│       ├── driving/                   # Inbound adapters
│       │   └── http/
│       │       ├── handler.go         # HTTP handlers
│       │       └── dto.go             # Request/Response DTOs
│       └── driven/                    # Outbound adapters
│           ├── postgres/
│           │   └── invoice_repo.go    # Database implementation
│           ├── memory/
│           │   └── invoice_repo.go    # In-memory implementation for tests
│           ├── pdf/
│           │   └── generator.go       # PDF generation
│           └── eventbus/
│               └── publisher.go       # Event publishing
├── migrations/                        # Database migrations
├── docs/                             # API documentation
├── go.mod
└── go.sum
```

## Implementation Phases

### Phase 0: Baseline (0.5 day)

**Objective**: Establish migration foundation

**Tasks**:
1. Create a snapshot/branch of current working code
2. Set up testing framework and basic test structure
3. Document current API endpoints and expected behavior
4. Freeze current handlers as "legacy adapter" (keep them working)

**Deliverables**:
- Git branch: `hexagonal-migration-baseline`
- Basic test setup
- API documentation

### Phase 1: Domain Carve-out (1-2 days)

**Objective**: Extract pure business logic from current code

**Critical**: Start with failing tests to drive the extraction

**Tasks**:
1. Create `internal/domain/` package
2. Define core entities:
   ```go
   // internal/domain/invoice.go
   type Invoice struct {
       ID          InvoiceID
       CustomerID  string
       LineItems   []LineItem
       Status      InvoiceStatus
       CreatedAt   time.Time
       FinalizedAt *time.Time
   }
   
   func (i *Invoice) AddLineItem(item LineItem) error { ... }
   func (i *Invoice) CalculateTotal() Money { ... }
   func (i *Invoice) Finalize() error { ... }
   ```

3. Create value objects:
   ```go
   // internal/domain/value_objects.go
   type InvoiceID string
   type Money struct { Amount int64; Currency string }
   type InvoiceStatus string
   ```

4. Move business rules out of HTTP handlers
5. **No external imports** except standard library
6. Write comprehensive unit tests for domain logic

**Deliverables**:
- Pure domain package with business logic
- Comprehensive unit tests
- Zero external dependencies in domain

### Phase 2: Define Ports (0.5 day)

**Objective**: Create interface contracts for communication

**Tasks**:
1. Define inbound ports:
   ```go
   // internal/ports/inbound/invoice_service.go
   type InvoiceService interface {
       CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*Invoice, error)
       UpdateInvoice(ctx context.Context, id string, req UpdateInvoiceRequest) error
       FinalizeInvoice(ctx context.Context, id string) error
       GetInvoice(ctx context.Context, id string) (*Invoice, error)
       GeneratePDF(ctx context.Context, id string) ([]byte, error)
   }
   ```

2. Define outbound ports:
   ```go
   // internal/ports/outbound/invoice_repo.go
   type InvoiceRepository interface {
       Save(ctx context.Context, invoice *Invoice) error
       FindByID(ctx context.Context, id string) (*Invoice, error)
       List(ctx context.Context, filter ListFilter) ([]*Invoice, error)
   }
   
   // internal/ports/outbound/pdf_generator.go
   type PDFGenerator interface {
       Generate(ctx context.Context, invoice *Invoice) ([]byte, error)
   }
   
   // internal/ports/outbound/event_publisher.go
   type EventPublisher interface {
       PublishInvoiceCreated(ctx context.Context, invoice *Invoice) error
       PublishInvoiceFinalized(ctx context.Context, invoice *Invoice) error
   }
   ```

**Deliverables**:
- Complete port interfaces
- Clear separation of inbound vs outbound ports

### Phase 3: Domain-Level Tests (1 day)

**Objective**: Establish comprehensive test coverage for core logic

**Tasks**:
1. Write unit tests for all domain entities and business rules
2. Create mock implementations of outbound ports
3. Test invoice calculations, validation, state transitions
4. Test error conditions and edge cases
5. Achieve >90% test coverage for domain package

**Example**:
```go
func TestInvoice_CalculateTotal(t *testing.T) {
    invoice := &Invoice{
        LineItems: []LineItem{
            {Description: "Item 1", Quantity: 2, UnitPrice: Money{1000, "USD"}},
            {Description: "Item 2", Quantity: 1, UnitPrice: Money{500, "USD"}},
        },
    }
    
    total := invoice.CalculateTotal()
    expected := Money{2500, "USD"}
    assert.Equal(t, expected, total)
}
```

**Deliverables**:
- Comprehensive test suite for domain logic
- Mock implementations for testing
- Documentation of business rules

### Phase 4: Driven Adapters (2-3 days)

**Objective**: Implement outbound adapters (secondary adapters)

**Tasks**:
1. **PostgreSQL Repository**:
   ```go
   // internal/adapters/driven/postgres/invoice_repo.go
   type InvoiceRepository struct {
       db *sql.DB
   }
   
   func (r *InvoiceRepository) Save(ctx context.Context, invoice *domain.Invoice) error {
       // SQL implementation
   }
   ```

2. **Memory Repository** (for testing):
   ```go
   // internal/adapters/driven/memory/invoice_repo.go
   type InvoiceRepository struct {
       invoices map[string]*domain.Invoice
       mu       sync.RWMutex
   }
   ```

3. **PDF Generator Adapter**:
   ```go
   // internal/adapters/driven/pdf/generator.go
   type Generator struct {
       templatePath string
   }
   
   func (g *Generator) Generate(ctx context.Context, invoice *domain.Invoice) ([]byte, error) {
       // PDF generation logic
   }
   ```

4. **Event Publisher Adapter**:
   ```go
   // internal/adapters/driven/eventbus/publisher.go
   type Publisher struct {
       bus *eventbus.Bus
   }
   ```

5. Write integration tests for each adapter
6. Ensure all adapters satisfy their port interfaces

**Deliverables**:
- Complete outbound adapter implementations
- Integration tests for each adapter
- Database migration scripts

### Phase 5: Driving Adapters (1-1.5 days)

**Objective**: Implement inbound adapters (primary adapters)

**Tasks**:
1. **Create DTO layer**:
   ```go
   // internal/adapters/driving/http/dto.go
   type CreateInvoiceRequest struct {
       CustomerID string     `json:"customer_id"`
       LineItems  []LineItem `json:"line_items"`
   }
   
   type InvoiceResponse struct {
       ID         string    `json:"id"`
       CustomerID string    `json:"customer_id"`
       Total      string    `json:"total"`
       Status     string    `json:"status"`
       CreatedAt  time.Time `json:"created_at"`
   }
   ```

2. **HTTP Handlers**:
   ```go
   // internal/adapters/driving/http/handler.go
   type Handler struct {
       invoiceService ports.InvoiceService
   }
   
   func (h *Handler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
       var req CreateInvoiceRequest
       // Parse request, call service, return response
   }
   ```

3. **DTO ↔ Domain mapping functions**
4. Unit tests for handlers using httptest
5. Error handling and HTTP status code mapping

**Deliverables**:
- HTTP adapter with proper DTO mapping
- Unit tests for all handlers
- Error handling strategy

### Phase 6: Wire & Bootstrap (0.5 day)

**Objective**: Connect all components with dependency injection

**Tasks**:
1. **Create composition root**:
   ```go
   // cmd/server/main.go
   func main() {
       // Initialize driven adapters
       db := setupDatabase()
       invoiceRepo := postgres.NewInvoiceRepository(db)
       pdfGen := pdf.NewGenerator("./templates")
       eventPub := eventbus.NewPublisher(setupEventBus())
       
       // Initialize service with injected dependencies
       invoiceService := services.NewInvoiceService(invoiceRepo, pdfGen, eventPub)
       
       // Initialize driving adapters
       httpHandler := http.NewHandler(invoiceService)
       
       // Start server
       startHTTPServer(httpHandler)
   }
   ```

2. Consider using `github.com/google/wire` for compile-time DI
3. Configuration management (environment variables, config files)
4. Graceful shutdown handling

**Deliverables**:
- Complete application wiring
- Configuration management
- Startup/shutdown procedures

### Phase 7: Acceptance Tests (1 day)

**Objective**: End-to-end validation of complete system

**Tasks**:
1. **Docker Compose setup**:
   ```yaml
   # docker-compose.test.yml
   services:
     postgres:
       image: postgres:15
       environment:
         POSTGRES_DB: invoice_test
     app:
       build: .
       depends_on:
         - postgres
   ```

2. **End-to-end tests**:
   ```go
   func TestCreateInvoiceE2E(t *testing.T) {
       // Start test environment
       // Make HTTP request
       // Verify database state
       // Check generated events
   }
   ```

3. API contract tests
4. Performance baseline tests
5. Documentation updates

**Deliverables**:
- Complete end-to-end test suite
- Docker-based test environment
- Updated API documentation

## Migration Strategy

### Gradual Migration Approach

1. **Keep legacy handlers working** during migration
2. **Route new endpoints** to hexagonal implementation
3. **Migrate endpoints incrementally**:
   - Start with `POST /invoice` (create)
   - Then `GET /invoice/{id}` (read)
   - Finally complex operations (finalize, PDF)

4. **Dual-write pattern** for database changes if needed
5. **Feature flags** to switch between old and new implementations

### Risk Mitigation

- **Comprehensive testing** at each phase
- **Rollback plan** for each phase
- **Performance monitoring** during migration
- **Gradual traffic shifting** in production

## Success Criteria

### Technical Metrics

- ✅ **100% test coverage** for domain logic
- ✅ **<100ms P95 latency** maintained
- ✅ **Zero breaking changes** to existing API
- ✅ **All integration tests passing**

### Architectural Goals

- ✅ **Technology independence**: Core logic has no framework dependencies
- ✅ **Testability**: Can unit test business logic without external dependencies
- ✅ **Flexibility**: Can swap database/PDF library without changing core logic
- ✅ **Maintainability**: Clear separation of concerns and single responsibility

## Time Estimates

| Phase | Duration | Effort | Critical Path |
|-------|----------|--------|---------------|
| Phase 0: Baseline | 0.5 day | Setup | Testing framework |
| Phase 1: Domain | 1-2 days | High | Business logic extraction |
| Phase 2: Ports | 0.5 day | Medium | Interface definition |
| Phase 3: Tests | 1 day | High | Domain test coverage |
| Phase 4: Driven Adapters | 2-3 days | High | Database integration |
| Phase 5: Driving Adapters | 1-1.5 days | Medium | HTTP layer |
| Phase 6: Wiring | 0.5 day | Low | Dependency injection |
| Phase 7: E2E Tests | 1 day | Medium | Integration validation |
| **Buffer & Review** | 1 day | - | Documentation |

**Total: 8-10 person-days**

## Tools and Dependencies

### Development Tools

- **Testing**: Go's built-in testing + testify
- **Mocking**: Go interfaces (no external mocking library needed)
- **Database**: PostgreSQL with `database/sql` or `sqlx`
- **HTTP**: Standard library `net/http` or Gin
- **Dependency Injection**: Manual or `github.com/google/wire`

### Quality Assurance

- **Linting**: `golangci-lint`
- **Formatting**: `gofmt`
- **Security**: `govulncheck`
- **Coverage**: `go test -cover`

### Documentation

- **API Docs**: Swagger/OpenAPI with `swaggo`
- **Architecture**: Mermaid diagrams
- **Code**: Godoc comments

## Next Steps

1. **Read all learning resources** (Phase 0)
2. **Set up development environment** with testing framework
3. **Create migration branch** and begin Phase 1
4. **Follow phases sequentially** - do not skip ahead
5. **Test thoroughly** at each phase before proceeding
6. **Document decisions** and architectural choices

## Troubleshooting

### Common Pitfalls

- **Starting with HTTP layer**: Always extract domain first
- **Tight coupling**: Ensure domain has no external dependencies
- **Missing DTO mapping**: Keep JSON schema changes from leaking to domain
- **Inadequate testing**: Each phase must have comprehensive tests
- **Big bang migration**: Migrate incrementally, keep old code working

### When to Seek Help

- **Complex business rules**: If invoice calculations become intricate
- **Performance issues**: If tests run slowly or database queries are inefficient
- **Dependency cycles**: If packages import each other
- **Test failures**: If existing functionality breaks during migration

This guide provides a comprehensive roadmap for successfully implementing Hexagonal Architecture in the go-invoice application while maintaining system stability and ensuring high code quality.
