package alfred

import "testing"

func TestSubString(t *testing.T) {
	word := "bingowashisnameo"
	if word[0:5] != "bino" {
		t.Fatalf(word[0:5])
	}
}
