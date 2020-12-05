package model

import "testing"

// TestUser ...
func TestUser(t *testing.T) *UModel {
	t.Helper()

	return &UModel{
		Email:    "testuser@example.org",
		Password: "testpassword",
		Role:     "user",
	}
}
