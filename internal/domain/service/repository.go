package service

// Repository is the interface for the service repository
type Repository interface {
	// All returns all services
	All() ([]*Service, error)
	// One returns one service
	One(id Id) (*Service, error)

	// Oneshot persists one service
	Oneshot(service *Service) error
	// Longrun persists one service
	Longrun(service *Service) error
	// Bundle persists one service
	Bundle(service *Service) error

	// Remove removes a service
	Remove(service *Service) error
}
