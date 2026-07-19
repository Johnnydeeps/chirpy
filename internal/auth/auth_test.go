package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

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

func TestJWT(t *testing.T) {
	// test token creation and validation functionality
	testID := uuid.New()
	signedToken, _ := MakeJWT(testID, "secret", time.Hour)
	signedTokenID, err := ValidateJWT(signedToken, "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if testID != signedTokenID {
		t.Errorf("expected %v, got %v", testID, signedTokenID)
	}
}

func TestJWTExpiredToken(t *testing.T) {
	// test an expired token with -time.Hour
	expiredToken, _ := MakeJWT(uuid.New(), "secret", -time.Hour)
	_, err := ValidateJWT(expiredToken, "secret")
	if err == nil {
		t.Errorf("expected error for expired token, got none")
	}
}

func TestJWTWrongSecretToken(t *testing.T) {
	// test if the sever secret token is intentionally invalid
	token, _ := MakeJWT(uuid.New(), "correct-secret", time.Hour)
	_, err := ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Errorf("expected error for wrong secret, got none")
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headerVal string
		wantToken string
		wantErr   bool
	}{
		{"valid", "Bearer abc123", "abc123", false},
		{"missing", "", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			headers := http.Header{}
			if test.headerVal != "" {
				headers.Set("Authorization", test.headerVal)
			}
			got, err := GetBearerToken(headers)
			if (err != nil) != test.wantErr {
				t.Errorf("wantErr %v, got err %v", test.wantErr, err)
			}
			if got != test.wantToken {
				t.Errorf("want %q, got %q", test.wantToken, got)
			}

		})
	}
}
