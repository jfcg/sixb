package sixb

import (
	"testing"
)

func Test1(t *testing.T) {
	for i := 0; i < 256; i++ {
		c := byte(i)
		if c != Sb2an(An2sb(c)) {
			t.Fatal("inverses do not work", c)
		}
	}

	l := "0:@Zaz"
	for i := 0; i < 6; i += 2 {
		for c := l[i]; c <= l[i+1]; c++ {
			if An2sb(c) > 63 {
				t.Fatal("domain error")
			}
		}
	}

	for c := byte(' '); c < '0'; c++ {
		if An2sb(c) > 127 {
			t.Fatal("domain error")
		}
	}
}

func Test2(t *testing.T) {
	x := []byte("qwert12345")
	y := Bs2is(x)
	if len(y) != 1 || cap(y) != 1 || y[0] != 3689065420975077233 {
		t.Fatal("slice conversion error")
	}
	z := Is2bs(y)
	if len(z) != 8 || cap(z) != 8 || &z[0] != &x[0] {
		t.Fatal("slice conversion error")
	}
}
