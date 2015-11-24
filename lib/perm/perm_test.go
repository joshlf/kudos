package perm

import (
	"math/rand"
	"testing"
)

func TestParse(t *testing.T) {
	rand.Seed(30682)
	on := "rwxrwxrwx"
	off := "---------"
	for i := 0; i < 1000; i++ {
		var str string
		for i, c := range on {
			if rand.Int()%2 == 0 {
				c = rune(off[i])
			}
			str += string(c)
		}
		got := Parse(str).String()[1:]
		if got != str {
			t.Errorf("parsed %v, got %v", str, got)
		}
	}

	on = "rwx"
	off = "---"
	for i := 0; i < 1000; i++ {
		var str string
		for i, c := range on {
			if rand.Int()%2 == 0 {
				c = rune(off[i])
			}
			str += string(c)
		}
		got := ParseSingle(str).String()[7:]
		if got != str {
			t.Errorf("parsed %v, got %v", str, got)
		}
	}
}
