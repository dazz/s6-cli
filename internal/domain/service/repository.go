package service

// Repository is the interface for the service repository
type Repository interface {

	// All returns all services
	All() ([]*Service, error)
	//Create(service Service) error
	//Remove(id Id) error
}
