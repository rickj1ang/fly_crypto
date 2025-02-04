package checker

import (
	"sync"
	"testing"
)

func TestIsAbove(t *testing.T) {
	var m sync.Map
	m.Store("BTC", 95000.0)
	above := isAbove(100000.0, "BTC", &m)

	if above != false {
		t.Error("expect: flase, get: true")
	}
}
func TestIsBelow(t *testing.T) {
	var m sync.Map
	m.Store("BTC", 95000.0)
	below := isBelow(70000.0, "BTC", &m)

	if below != false {
		t.Error("expect: flase, get: true")
	}
}
