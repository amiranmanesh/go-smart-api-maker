package encrypting

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash := GetHashedPassword("password")
	assert.NotEmpty(t, hash)
}

func TestVerifyPassword(t *testing.T) {
	hash := GetHashedPassword("password")
	assert.True(t, CheckPassword(hash, "password"))
	assert.False(t, CheckPassword(hash, "password"))
}
