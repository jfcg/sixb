package sixb

import (
	"github.com/jfcg/sorty"
	"testing"
)

// Txt2int collision test
func Test0(t *testing.T) {
	const N = 3 << 27
	hl := make([]uint64, N)  // hash list: 3 GiB ram
	bf := []byte{0, 0, 0, 0} // input buffer

	// fill hl with hashes of short utf8 text
	for i, l := N-1, 0; i >= 0; i-- {
		hl[i] = Txt2int(string(bf[:l]))

		// next utf8-ish input
		for k := 0; ; k++ {
			if bf[k] == 0 {
				l++ // increase input length
			}
			bf[k]++

			if bf[k] != 0 {
				break
			}
			bf[k]++ // continue with carry
		}
	}

	sorty.ArU8 = hl // sort hl
	sorty.SortU8()

	k := 0 // count collisions
	for i := N - 1; i > 0; i-- {
		if hl[i] == hl[i-1] {
			k++
		}
	}
	if k > 0 {
		t.Fatal("Txt2int has at least", k, "collisions for short utf8 inputs")
	}
}

// An2sb & Sb2an bijection & domain
func Test1(t *testing.T) {
	for i := 255; i >= 0; i-- {
		c := byte(i)
		d := An2sb[c]
		if c == d {
			t.Fatal("fixed point", i)
		}
		if c != Sb2an[d] {
			t.Fatal("inverse does not work", i)
		}
	}

	l := "0:@Zaz"
	for i := 4; i >= 0; i -= 2 {
		for c := l[i]; c <= l[i+1]; c++ {
			if An2sb[c] > 63 {
				t.Fatal("domain error", c)
			}
		}
	}
}

// Slice conversions
func Test2(t *testing.T) {
	v := "qwert123"
	w := v + "45"
	x := []byte(w)

	y := BtI8(x)
	z := I8tB(y)
	p := BtI4(x)
	q := I4tB(p)

	a := StI8(w)
	b := I8tS(a)
	r := StI4(w)
	s := I4tS(r)

	if len(y) != 1 || cap(y) != 1 || y[0] != 3689065420975077233 ||
		len(a) != 1 || cap(a) != 1 || a[0] != 3689065420975077233 ||
		len(p) != 2 || cap(p) != 2 || p[0] != 1919252337 || p[1] != 858927476 ||
		len(r) != 2 || cap(r) != 2 || r[0] != 1919252337 || r[1] != 858927476 ||
		len(z) != 8 || cap(z) != 8 || &z[0] != &x[0] ||
		len(q) != 8 || cap(q) != 8 || &q[0] != &x[0] ||
		b != v || s != v {
		t.Fatal("slice/string conversion error")
	}
}
