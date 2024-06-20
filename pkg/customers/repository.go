package customers

type Repository interface {
	CustomersAdd(email string) error
}
