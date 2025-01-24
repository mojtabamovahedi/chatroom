package nanoId

import (
	"github.com/aidarkhanov/nanoid"
)

const (
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	size     = 5
)


// generate a nanoID with numbers and english alphabet with 5 length
func GenerateId() (string, error) {
	return nanoid.Generate(alphabet, size)
}
