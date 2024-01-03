package create

import "github.com/dazz/s6-cli/internal/domain/service"

func Create(iterator service.StepIterator) error {
	return iterator.Iterate()
}
