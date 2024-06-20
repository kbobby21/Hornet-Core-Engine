package sessmanager

type Repository interface {
	// stores hash:email and returns the hash
	StoreEmailHash(string) (string, error)

	// gets back email from hash
	GetEmailFromHash(string) (string, error)
}
