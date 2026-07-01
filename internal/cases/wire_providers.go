package cases

import "github.com/google/wire"

// ProviderSet groups the cases domain's constructors for Wire. NewService also needs
// a UserReader; the composition root supplies it via wire.Bind to user.Service.
var ProviderSet = wire.NewSet(
	NewRepository,
	NewService,
	NewHandler,
)
