//go:build wireinject

package server

import (
	"github.com/fardinabir/go-svc-boilerplate/internal/cases"
	"github.com/fardinabir/go-svc-boilerplate/internal/user"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// InitializeHandlers is the Wire injector: it declares WHAT each domain needs and
// lets Wire generate the wiring into wire_gen.go. The body is replaced entirely by
// generated code — do not edit. Run `make wire` after changing a provider set.
func InitializeHandlers(db *gorm.DB) (*Handlers, error) {
	wire.Build(
		user.ProviderSet,
		cases.ProviderSet,
		// cases.Service needs a UserReader; user.Service satisfies it.
		wire.Bind(new(cases.UserReader), new(user.Service)),
		wire.Struct(new(Handlers), "*"),
	)
	return nil, nil
}
