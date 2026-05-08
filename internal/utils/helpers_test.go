package utils

import "testing"

func TestToWSLPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`D:\Obras\Manga`, `/mnt/d/Obras/Manga`},
		{`C:\Users\John\manga`, `/mnt/c/Users/John/manga`},
		{`D:/Obras/Manga`, `/mnt/d/Obras/Manga`},
		{`/home/user/manga`, `/home/user/manga`},
		{``, ``},
		{`D:`, `D:`},
	}

	for _, tt := range tests {
		result := ToWSLPath(tt.input)
		if result != tt.expected {
			t.Errorf("ToWSLPath(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}
