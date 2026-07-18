package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("didn't expect an error, but got: %v", err)
	}
	if hash == "" {
		t.Errorf("expected a non-empty result")
	}
	// checking correct behaviour checkpasswordhash should return true based on the above var hash.
	check, err := CheckPasswordHash("password", hash)
	if err != nil {
		t.Errorf("didn't expect an error, but got: %v", err)
	}
	if check == false {
		t.Errorf("expected check to pass as true")
	}
	// checking if password is incorrect, should return false

	check2, _ := CheckPasswordHash("123456", hash)
	if check2 == true {
		t.Errorf("expected check to pass as false, diliberate password mismatch to hash")
	}

}
