package category

// Category is an interface for forming users into groups.
// Category is the nation, which users are citizens of.
// Membership in a category is maintained by a sorted list.
type Category interface {
	// Commands
	AddUser(uid interface{}) error
	RemoveUser(uid interface{}) error

	Contains(uid interface{}) (bool, error)
}

type Group struct {
}
