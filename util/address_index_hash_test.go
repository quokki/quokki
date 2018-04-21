package util

import (
	"testing"

	"encoding/hex"
)

func BenchmarkIndexHash(b *testing.B) {
	address, _ := hex.DecodeString("52B27304F46FB9E53D4EF35A2FA9101CD63683E4")
	for n := 0; n < b.N; n++ {
		GetAddressIndexHash(address, int64(n), "test")
	}
}
