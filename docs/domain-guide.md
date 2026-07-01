# Domain Guide

Complete reference for adding domains, wiring dependencies, and reading across domain boundaries.

---

## Adding a New Domain

Every domain follows the same template. To add a domain named `billing`:

### 1. Create `internal/billing/` with these files

```go
// model.go
package billing

import "time"

type Invoice struct {
    ID        int       `gorm:"primaryKey" json:"id"`
    CaseID    int       `gorm:"not null"   json:"case_id"`
    Amount    float64   `gorm:"not null"   json:"amount"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
```

```go
// repository.go
package billing

import "gorm.io/gorm"

type Repository interface {
    Create(inv *Invoice) error
    FindByCaseID(caseID int) ([]Invoice, error)
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(inv *Invoice) error { return r.db.Create(inv).Error }
func (r *repository) FindByCaseID(id int) ([]Invoice, error) {
    var invs []Invoice
    return invs, r.db.Where("case_id = ?", id).Find(&invs).Error
}
```

```go
// service.go
package billing

type Service interface {
    CreateInvoice(inv *Invoice) error
    InvoicesForCase(caseID int) ([]Invoice, error)
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) CreateInvoice(inv *Invoice) error          { return s.repo.Create(inv) }
func (s *service) InvoicesForCase(id int) ([]Invoice, error) { return s.repo.FindByCaseID(id) }
```

```go
// handler.go
package billing

import (
    "net/http"

    apierr "github.com/fardinabir/go-svc-boilerplate/internal/errors"
    "github.com/fardinabir/go-svc-boilerplate/pkg/response"
    "github.com/fardinabir/go-svc-boilerplate/pkg/web"
    "github.com/labstack/echo/v4"
)

type Handler interface {
    CreateInvoice(c echo.Context) error
    ListByCaseID(c echo.Context) error
}

type handler struct {
    web.Base
    service Service
}

func NewHandler(s Service) Handler { return &handler{service: s} }

func (h *handler) CreateInvoice(c echo.Context) error {
    var inv Invoice
    if err := h.MustBind(c, &inv); err != nil {
        return response.Respond(c, apierr.ErrBadRequest, err.Error())
    }
    if err := h.service.CreateInvoice(&inv); err != nil {
        return response.Respond(c, apierr.ErrInternalServerError)
    }
    return c.JSON(http.StatusCreated, response.ResponseData{Data: inv})
}
```

```go
// routes.go
package billing

import "github.com/labstack/echo/v4"

func RegisterRoutes(g *echo.Group, h Handler) {
    invoices := g.Group("/invoices")
    invoices.POST("", h.CreateInvoice)
    invoices.GET("/case/:id", h.ListByCaseID)
}
```

```go
// wire_providers.go
package billing

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewRepository, NewService, NewHandler)
```

### 2. Register the model for auto-migration

In `internal/db/db.go`:

```go
db.AutoMigrate(&user.User{}, &cases.Case{}, &billing.Invoice{})
```

### 3. Add the domain to the Wire injector

In `internal/server/wire.go`:

```go
func InitializeHandlers(db *gorm.DB) (*Handlers, error) {
    wire.Build(
        user.ProviderSet,
        cases.ProviderSet,
        billing.ProviderSet,                                    // ← add
        wire.Bind(new(cases.UserReader), new(user.Service)),
        wire.Struct(new(Handlers), "*"),
    )
    return nil, nil
}
```

### 4. Add the handler to the `Handlers` struct

In `internal/server/api.go`:

```go
type Handlers struct {
    User    user.Handler
    Cases   cases.Handler
    Billing billing.Handler    // ← add
}
```

### 5. Mount the routes and register validators

In `setupRoutes` in `internal/server/api.go`:

```go
billing.RegisterRoutes(api, h.Billing)    // ← add
```

If the domain defines custom validation tags via a `RegisterValidations` function, also append it to the validator setup in the same file:

```go
e.Validator = web.NewCustomValidator(
    user.RegisterValidations,
    billing.RegisterValidations,    // ← add if domain has custom tags
)
```

### 6. Regenerate Wire and add a migration

```bash
make wire
```

Add `migrations/ddl/003_create_billing_schema.sql` if the domain needs tables beyond what GORM auto-migration produces.

---

## Cross-Domain Reads

When domain A needs data from domain B, A declares a **port interface** that it owns. Domain B's service satisfies it. The composition root binds them at build time. Domain A never imports domain B's model or repository — no import cycle, no tight coupling.

This also means the dependency is swappable: if `billing` later splits into its own service, `cases.CaseReader` can be backed by an HTTP client instead of a direct service call, with zero changes to `billing`.

### Step 1 — Declare the port in the consumer

Create `internal/billing/ports.go`:

```go
package billing

// CaseReader is the slice of the cases domain that billing needs.
// Declared here (in the consumer), not in cases — billing owns this contract.
type CaseReader interface {
    StatusByID(caseID int) (string, error)
}
```

Keep ports narrow. Only declare the methods this domain actually needs.

### Step 2 — Inject the port into the service

In `internal/billing/service.go`:

```go
type service struct {
    repo   Repository
    cases  CaseReader   // ← cross-domain dependency via port
}

func NewService(repo Repository, cases CaseReader) Service {
    return &service{repo: repo, cases: cases}
}
```

### Step 3 — Implement the method in the provider

In `internal/cases/service.go`, add the method if not already present:

```go
func (s *service) StatusByID(caseID int) (string, error) {
    c, err := s.repo.FindByID(caseID)
    if err != nil {
        return "", err
    }
    return c.Status, nil
}
```

`cases.Service` now satisfies `billing.CaseReader` implicitly — no changes to its interface declaration.

### Step 4 — Bind at the composition root

In `internal/server/wire.go`:

```go
wire.Bind(new(billing.CaseReader), new(cases.Service)),
```

Wire generates the injection automatically. `billing` never imports `cases`.

### The rule

| ✅ Allowed | ❌ Not allowed |
|---|---|
| `billing` imports `billing.CaseReader` (its own port) | `billing` imports `cases.Service` |
| `cases.Service` satisfies `billing.CaseReader` implicitly | `billing` imports `cases.Repository` |
| Composition root imports both domains | Any domain imports another domain's `model.go` |

---

## Wire DI Reference

[Google Wire](https://github.com/google/wire) generates compile-time dependency wiring. No reflection — a missing provider or type mismatch is a build error, not a nil panic at runtime.

### How it works

Wire reads `internal/server/wire.go` (guarded by `//go:build wireinject`) and generates `internal/server/wire_gen.go`. The generated file is committed and used in normal builds. The injector file is excluded from normal compilation by the build tag.

```
wire.go  (injector declaration, not compiled normally)
    ↓  wire generates
wire_gen.go  (concrete wiring, compiled normally)
```

### Provider sets

Each domain declares a provider set in `wire_providers.go`:

```go
var ProviderSet = wire.NewSet(NewRepository, NewService, NewHandler)
```

Wire reads these and resolves the full dependency graph. A provider is any function `func NewX(deps...) X` — Wire matches output types to input types.

### When to run `make wire`

Run `make wire` after any of these changes:

| Change | Reason |
|---|---|
| Add a new domain's `ProviderSet` to `wire.Build(...)` | Wire needs to generate wiring for the new providers |
| Add or remove a constructor parameter | The resolved type graph changes |
| Add a `wire.Bind(...)` for a new interface | Wire needs to know which concrete type satisfies the interface |
| Rename a `New*` function | Old provider reference breaks the graph |

You do **not** need to run `make wire` for:
- Changes inside a function body (not signature)
- Adding methods to a type (not constructors)
- Changes to business logic

### Troubleshooting

**`wire: no provider found for <type>`** — a constructor's dependency has no registered provider. Add the missing `New*` function to the relevant `ProviderSet`, or add a `wire.Bind` if it's an interface.

**`wire: cycle detected`** — two providers depend on each other. This usually means a cross-domain dependency was imported directly instead of through a port interface.

**Generated file out of date** — if `wire_gen.go` doesn't compile, run `make wire` and commit the updated file.
