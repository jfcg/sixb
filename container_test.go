/*	Copyright (c) 2019-present, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sixb

import (
	"testing"
	"unsafe"
)

func TestSet(t *testing.T) {
	if unsafe.Sizeof(None{}) != 0 {
		t.Fatal("None type should have zero size")
	}

	s := NewSet[int]()
	if s.Size() != 0 {
		t.Fatal("empty set should have no elements")
	}

	// add 1,3,5 to set
	for i, n := 1, 1; i <= 5; i, n = i+2, n+1 {
		s.Add(i)
		// increasing size?
		if s := s.Size(); s != n {
			t.Fatal("expected size:", n, "got:", s)
		}
	}

	// already present
	for i := 5; i >= 1; i -= 2 {
		s.Add(i)
		if s := s.Size(); s != 3 {
			t.Fatal("expected size: 3 got:", s)
		}
	}

	// should have 1,3,5
	for i := 5; i >= 1; i -= 2 {
		if !s.HasAny(i) || !s.HasAny(i, 2) || !s.HasAny(i, 2, 4) {
			t.Fatal("set should have:", i)
		}
	}
	if !s.HasAll(3) || !s.HasAll(5, 1) || !s.HasAll(3, 5, 1) || s.HasAll(3, 4, 5, 1) {
		t.Fatal("set should have: 1,3,5")
	}

	// should not have 2,4,6,7,8
	for i := 2; i <= 6; i += 2 {
		if s.HasAny(i) || s.HasAny(i, 7) || s.HasAny(i, 7, 8) ||
			s.HasAll(i) || s.HasAll(i, 1) || s.HasAll(i, 3, 5) {
			t.Fatal("set should not have:", i)
		}
	}

	s.Add(4, 2, 6)
	if s := s.Size(); s != 6 {
		t.Fatal("expected size: 6 got:", s)
	}
	// should have 1,2,3,4,5,6
	for i := 6; i >= 1; i-- {
		if !s.HasAny(i) || !s.HasAny(i+1, i) ||
			!s.HasAny(i, 8) || !s.HasAny(9, i, 8) {
			t.Fatal("set should have:", i)
		}
	}
	if !s.HasAll(2, 3, 1) || !s.HasAll(5, 4, 2, 3) ||
		!s.HasAll(3, 5, 1, 2, 4, 6) {
		t.Fatal("set should have: 1,2,3,4,5,6")
	}

	s.Remove(1)
	if s.Size() != 5 || s.HasAny(1) || s.HasAny(7, 1) ||
		!s.HasAll(6, 4, 2, 5, 3) {
		t.Fatal("set should not have: 1")
	}
	s.Remove(3, 2)
	if s.Size() != 3 || s.HasAny(2, 1) || s.HasAny(3) ||
		s.HasAny(3, 1, 2) || !s.HasAll(6, 4, 5) {
		t.Fatal("set should not have: 1,2,3")
	}
}

func TestCircleQ(t *testing.T) {
	// CircleQ should have four distinct states:

	// null: Push & Pop wont work
	var q CircleQ[int]
	if q.Capacity() != 0 || q.Size() != 0 {
		t.Fatal("null state: wrong capacity or size")
	}
	if q.Push(6) {
		t.Fatal("null state: should not Push")
	}
	if _, ok := q.Pop(); ok {
		t.Fatal("null state: should not Pop")
	}

	q.Reset(2) // capacity should be 2

	// empty: Only Push will work
	if q.Capacity() != 2 || q.Size() != 0 {
		t.Fatal("empty state: wrong capacity or size")
	}
	if _, ok := q.Pop(); ok {
		t.Fatal("empty state: should not Pop")
	}
	if !q.Push(45) {
		t.Fatal("empty state: should Push")
	}

	// partial: Push & Pop will work
	if q.Capacity() != 2 || q.Size() != 1 {
		t.Fatal("partial state: wrong capacity or size")
	}
	// deep copy for independent Push & Pop calls
	c := q.Copy()
	if &c.data[0] == &q.data[0] || c.Capacity() != q.Capacity() ||
		c.begin != q.begin || c.end != q.end {
		t.Fatal("invalid deep copy")
	}
	if item, ok := c.Pop(); !ok || item != 45 {
		t.Fatal("partial state: should Pop top item")
	}
	// c should be empty here
	if c.Capacity() != 2 || c.Size() != 0 {
		t.Fatal("empty state: wrong capacity or size")
	}

	if !q.Push(46) {
		t.Fatal("partial state: should Push")
	}

	// full: Only Pop will work
	if q.Capacity() != 2 || q.Size() != 2 {
		t.Fatal("full state: wrong capacity or size")
	}
	if q.Push(47) {
		t.Fatal("full state: should not Push")
	}
	if item, ok := q.Pop(); !ok || item != 45 {
		t.Fatal("full state: should Pop top item")
	}

	// partial: Push & Pop will work
	if q.Capacity() != 2 || q.Size() != 1 {
		t.Fatal("partial state: wrong capacity or size")
	}
	// deep copy for independent Push & Pop calls
	c = q.Copy()
	if !c.Push(47) {
		t.Fatal("partial state: should Push")
	}
	// c should be full here
	if item, ok := c.Pop(); !ok || item != 46 {
		t.Fatal("full state: should Pop top item")
	}
	if item, ok := c.Pop(); !ok || item != 47 {
		t.Fatal("partial state: should Pop top item")
	}
	// c is now empty
	if c.Capacity() != 2 || c.Size() != 0 {
		t.Fatal("empty state: wrong capacity or size")
	}

	// q should be partial here
	if item, ok := q.Pop(); !ok || item != 46 {
		t.Fatal("partial state: should Pop top item")
	}
	// q is now empty
	if q.Capacity() != 2 || q.Size() != 0 {
		t.Fatal("empty state: wrong capacity or size")
	}

	for _, v := range q.data {
		if v != 0 {
			t.Fatal("reference to past items")
		}
	}
	for _, v := range c.data {
		if v != 0 {
			t.Fatal("reference to past items")
		}
	}

	c = NewCircleQ[int](3) // same as Reset
	if c.Capacity() != 3 || c.Size() != 0 {
		t.Fatal("wrong capacity or size with NewCircleQ")
	}
}
