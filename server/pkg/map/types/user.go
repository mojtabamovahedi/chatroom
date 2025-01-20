package types

type UserRole uint

const (
	ADMIN UserRole = iota + 1
	USER
)

type User struct {
	Id   string
	Name string
	Role UserRole
}
