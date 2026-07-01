package cases

// UserReader is the slice of the user domain that cases needs. It is declared here,
// in the consuming domain, and owns only the shape cases requires. user.Service
// satisfies it; the composition root binds them with wire.Bind. cases never imports
// the user model or repository, so there is no import cycle and the dependency can
// later be swapped for a remote call if user is split into its own service.
type UserReader interface {
	EmailByID(id int) (string, error)
}
