package repository

// Repository defines the interface for data access
type Repository interface {
	// Add data access methods here later
}

type repositoryImpl struct {
	// e.g., db connection
}

// New returns a new Repository implementation
func New() Repository {
	return &repositoryImpl{}
}
