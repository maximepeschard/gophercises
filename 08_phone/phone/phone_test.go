package phone

import "testing"

func TestParseNumber(t *testing.T) {
	cases := []struct{ in, want string }{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"123-456-7890", "1234567890"},
		{"1234567892", "1234567892"},
		{"(123)456-7892", "1234567892"},
	}

	for _, c := range cases {
		got, err := Normalize(c.in)
		if err != nil {
			t.Errorf("ParseNumber(%s) returns non-nil error", c.in)
		} else if got != c.want {
			t.Errorf("ParseNumber(%s) == %s, want %s", c.in, got, c.want)
		}
	}
}

func TestParseNumberInvalid(t *testing.T) {
	cases := []string{"123", "12345678900", "helloworld"}

	for _, c := range cases {
		_, err := Normalize(c)
		if err == nil {
			t.Errorf("ParseNumber(%s) returns nil error", c)
		}
	}
}
