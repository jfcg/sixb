package sixb

import "testing"

func Test1(t *testing.T) {
	for i := 255; i >= 0; i-- {
		c := byte(i)
		d := An2sb(c)
		if c != Sb2an(d) {
			t.Fatal("inverse does not work", i)
		}
	}

	l := "0:@Zaz"
	for i := 4; i >= 0; i -= 2 {
		for c := l[i]; c <= l[i+1]; c++ {
			if An2sb(c) > 63 {
				t.Fatal("domain error", c)
			}
		}
	}
}

func Test2(t *testing.T) {
	x := []byte("qwert12345")
	y := Bs2is(x)
	z := Is2bs(y)
	if len(y) != 1 || cap(y) != 1 || y[0] != 3689065420975077233 ||
		len(z) != 8 || cap(z) != 8 || &z[0] != &x[0] {
		t.Fatal("slice conversion error")
	}
}
