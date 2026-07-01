package user

import "github.com/google/wire"

// ProviderSet groups the user domain's constructors for Wire. The composition
// root adds this to its injector to build the user handler's repo→service→handler chain.
var ProviderSet = wire.NewSet(
	NewRepository,
	NewService,
	NewHandler,
)
