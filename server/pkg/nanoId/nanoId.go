package nanoId

import (
	"github.com/aidarkhanov/nanoid"
)

const (
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	size     = 5
)

func GenerateId() (string, error) {
	return nanoid.Generate(alphabet, size)
}
