package checker

import "testing"

func TestIsAbove(t *testing.T) {
	above, err := isAbove(100000.0, "BTC")
	if err != nil {
		t.Log(err.Error())
	}
	if above != false {
		t.Error("expect: flase, get: true")
	}
}
func TestIsBelow(t *testing.T) {
	below, err := isBelow(70000.0, "BTC")
	if err != nil {
		t.Log(err.Error())
	}
	if below != false {
		t.Error("expect: flase, get: true")
	}
}
