package model

import "testing"

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		// Valid names
		{"simple lowercase", "myapp", true},
		{"simple uppercase", "MYAPP", true},
		{"mixed case", "MyApp", true},
		{"with numbers", "app123", true},
		{"with hyphen", "my-app", true},
		{"with underscore", "my_app", true},
		{"mixed special", "My-App_123", true},
		{"single char", "a", true},
		{"max length 128", string(make([]byte, 128)), false}, // need valid chars

		// Invalid names
		{"empty string", "", false},
		{"with space", "my app", false},
		{"with dot", "my.app", false},
		{"with slash", "my/app", false},
		{"with colon", "my:app", false},
		{"with at", "my@app", false},
		{"too long", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", false}, // 129 chars
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidName(tt.input)
			if got != tt.want {
				t.Errorf("IsValidName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidName_MaxLength(t *testing.T) {
	// Test exactly 128 characters (should be valid)
	validMaxLen := make([]byte, 128)
	for i := range validMaxLen {
		validMaxLen[i] = 'a'
	}
	if !IsValidName(string(validMaxLen)) {
		t.Errorf("IsValidName with 128 chars should be valid")
	}

	// Test 129 characters (should be invalid)
	invalidLen := make([]byte, 129)
	for i := range invalidLen {
		invalidLen[i] = 'a'
	}
	if IsValidName(string(invalidLen)) {
		t.Errorf("IsValidName with 129 chars should be invalid")
	}
}