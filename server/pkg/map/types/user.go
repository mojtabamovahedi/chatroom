package types

type UserRole uint

// all roles that user can get
const (
	ADMIN UserRole = iota + 1
	USER
)


// User represents a user in chatroom with ID, name and role
type User struct {
	Id   string
	Name string
	Role UserRole
}
