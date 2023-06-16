package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
			assert.Equal(t, got, c.want)
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
	ass := assert.New(t)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got1, got2 := SplitCharAndNum(c.input)
			if got1 != c.want1 || got2 != c.want2 {
				ass.Equal(got1, c.want1)
				ass.Equal(got2, c.want2)
			}
		})
	}
}
