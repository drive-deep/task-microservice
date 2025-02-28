package repositories

type Repository[T any] interface {
	Create(entity *T) error
	GetByID(id string) (*T, error)
	GetAll(filter map[string]interface{}, sort string, page, pageSize int) ([]T, error)
	Update(entity *T) error
	Delete(id string) error
}
