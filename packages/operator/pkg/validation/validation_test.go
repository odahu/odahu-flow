package validation

import (
	"strings"
	"testing"
)

func TestValidateID(t *testing.T) {
	tests := []struct {
		name    string
		IDArgs  string
		wantErr bool
	}{
		{
			"Valid iD",
			"some-valid-id",
			false,
		},
		{
			"If ID contains 1 symbol than it is valid",
			"s",
			false,
		},
		{
			"If ID starts with '-' than it is not valid",
			"-some-valid-id",
			true,
		},
		{
			"If ID ends with '-' than it is not valid",
			"some-valid-id-",
			true,
		},
		{
			"Empty ID is not valid",
			"",
			true,
		},
		{
			"If ID contains more 64 symbols, than it is not valid",
			strings.Repeat("a", 64),
			true,
		},
		{
			"Check max length of id",
			strings.Repeat("a", 63),
			false,
		},
		{
			"Check 3 length of id",
			strings.Repeat("a", 3),
			false,
		},
		{
			"If ID contains upper symbols, than it is not valid",
			"SOME-id",
			true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateID(tt.IDArgs); (err != nil) != tt.wantErr {
				t.Errorf("ValidateID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
