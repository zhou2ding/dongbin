package tools

import "testing"

func TestGetDigitOfStr(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"case1", "A23", "23"},
		{"case2", "AB", ""},
		{"case3", "gw-0046-mc", "0046"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := GetDigitOfStr(c.input)
			if got != c.want {
				t.Errorf("want: %v, but got: %v", c.want, got)
			}
		})
	}
}
