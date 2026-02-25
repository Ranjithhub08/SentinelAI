package repository

// Repository defines the interface for data access
type Repository interface {
}

type repositoryImpl struct {
}

// New returns a new Repository implementation
func New() Repository {
	return &repositoryImpl{}
}
