package auth

import "testing"

func TestCheckPasswordHash(t *testing.T) {
	password1 := "Password1!"
	password2 := "Password2!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name string
		password string
		hash string
		err bool
		match bool
	}{
		{
			name: "Correct password",
			password: password1,
			hash: hash1,
			err: false,
			match: true,
		},
		{
			name: "Incorrect password",
			password: "admin",
			hash: hash1,
			err: false,
			match: false,
		},
		{
			name: "Wrong hash used",
			password: password1,
			hash: hash2,
			err: false,
			match: false,
		},
		{
			name: "Empty string",
			password: "",
			hash: hash1,
			err: false,
			match: false,
		},
		{
			name: "Invalid hash",
			password: password1,
			hash: "invalidHash",
			err: true,
			match: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			match, err := CheckPasswordHash(test.password, test.hash)
			if (err != nil) != test.err {
				t.Errorf("CheckPasswordHash error: %v, err: %v", err, test.err)
			}
			if !test.err && match != test.match {
				t.Errorf("CheckPasswordHash expects &v, got %v", test.match, match)
			}
		})
	}
}