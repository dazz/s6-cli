package service

type Repository interface {
	All() ([]*Service, error)
	//Create(service Service) error
	//Remove(id Id) error
}
