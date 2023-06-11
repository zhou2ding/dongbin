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

func TestSplitCharAndNum(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want1 string
		want2 int
	}{
		{"case1", "gw-0043-mc-1", "gw--mc-1", 43},
		{"case2", "gw-0002-db", "gw--db", 2},
		{"case3", "123", "", 123},
		{"case4", "gw-5", "gw-", 5},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got1, got2 := SplitCharAndNum(c.input)
			if got1 != c.want1 || got2 != c.want2 {
				t.Errorf("want1: %v, wang2: %v, but go1t: %v, got2:%v\n", c.want1, c.want2, got1, got2)
			}
		})
	}
}
