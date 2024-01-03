package service

// Repository is the interface for the service repository
type Repository interface {
	// All returns all services
	All() ([]*Service, error)
	// One returns one service
	One(id Id) (*Service, error)

	RemoveDirectory(path string) error

	RemoveFile(path string) error

	ServiceScriptFilePath(id Id) (string, error)

	ServicePath(id Id) string

	ServiceDependenciesPath(s *Service) (string, error)

	FileExists(file string) bool

	AllFileExists(files []string) bool

	AnyFileExists(files []string) bool

	CreateDependencies(service *Service) error

	FindDependenciesById(id Id) ([]string, error)

	CreateTypeFile(*Service) error

	CreateFile(path string, content []byte, perm int) error

	CreateDirectory(path string, perm int) error
}
