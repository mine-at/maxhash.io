// SPDX-FileCopyrightText: 2025 maxhash.io <dev@maxhash.io>
//
// SPDX-License-Identifier: AGPL-3.0-only

package http_test

import (
	"testing"

	"github.com/mine-at/maxhash.io/http"
)

func TestIsValidBitcoinAddress(t *testing.T) {
	tests := []struct {
		address string
		want    bool
	}{
		// Valid mainnet addresses
		{"bc1pd5zyzu4cgdw0270ykue34dfpy8ezuc0laannduy33vlvz66ss2tqqcyzqx", true},

		// Invalid addresses
		{"", false}, // empty
		{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa12345678901234567890", false}, // too long
		{"4A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", false},                     // invalid prefix
		{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa!", false},                    // invalid character
		{"tb2qfm6xv4z8t7w8w7w8w7w8w7w8w7w8w7w8w7w8w7", false},             // invalid testnet prefix
	}

	for _, tt := range tests {
		got := http.IsValidBitcoinAddress(tt.address)
		if got != tt.want {
			t.Errorf("IsValidBitcoinAddress(%q) = %v, want %v", tt.address, got, tt.want)
		}
	}
}
