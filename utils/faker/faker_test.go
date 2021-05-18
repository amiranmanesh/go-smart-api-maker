package faker

import (
	"testing"
)

func TestFaker(tt *testing.T) {
	tt.Run("", func(t *testing.T) {
		result, _ := Call("ipv6")
		t.Log(result)
	})
}
