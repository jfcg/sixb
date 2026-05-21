/*	Copyright (c) 2019-present, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sixb

// None represents empty struct with zero size.
type None struct{}

// Set of elements of type K. Create empty set with
//
//	s := Set[K]{}
type Set[K comparable] map[K]None

// NewSet creates a set from given elements.
func NewSet[K comparable](els ...K) Set[K] {
	s := Set[K]{}
	s.Add(els...)
	return s
}

// Size returns number of elements in set.
func (s Set[K]) Size() int {
	return len(s)
}

// HasAny returns true only if set contains any of given elements.
// If els is empty, it returns false by definition.
func (s Set[K]) HasAny(els ...K) bool {
	for _, e := range els {
		if _, ok := s[e]; ok {
			return true
		}
	}
	return false
}

// HasAll returns true only if set contains all of given elements.
// If els is empty, it returns true by definition.
func (s Set[K]) HasAll(els ...K) bool {
	for _, e := range els {
		if _, ok := s[e]; !ok {
			return false
		}
	}
	return true
}

// Add elements to set.
func (s Set[K]) Add(els ...K) {
	for _, e := range els {
		s[e] = None{}
	}
}

// Remove elements from set.
func (s Set[K]) Remove(els ...K) {
	for _, e := range els {
		delete(s, e)
	}
}

// CircleQ is a circular, fifo queue of items of type T with a maximum capacity. Create a new one like:
//
//	// It has four distinct states:
//	var q CircleQ[int]  // null   : Push & Pop wont work
//	q.Reset(2)          // empty  : Only Push will work
//	q.Push(45)          // partial: Push & Pop will work
//	q.Push(46)          // full   : Only Pop will work
//	item, ok := q.Pop() // returns 45, true
//	r := NewCircleQ[int](2) // same as Reset
type CircleQ[T any] struct {
	data []T
	// end = begin: empty queue
	// end = ^0   : full queue
	begin, end uint32
}

// NewCircleQ creates a circular, fifo queue of given maximum capacity.
func NewCircleQ[T any](capacity uint32) (r CircleQ[T]) {
	r.Reset(capacity)
	return
}

// Reset the queue to empty with new maximum capacity.
func (q *CircleQ[T]) Reset(capacity uint32) {
	q.data = make([]T, capacity)
	q.begin = 0
	q.end = 0
}

// Capacity is the maximum number of items the queue can store.
func (q *CircleQ[T]) Capacity() int {
	return len(q.data)
}

// Size is the number of current items in the queue.
func (q *CircleQ[T]) Size() int {
	cap := q.Capacity()
	end := q.end
	if ^end == 0 { // full ?
		return cap
	}
	size := int(end) - int(q.begin)
	if size < 0 {
		return cap + size
	}
	return size
}

// Copy returns a deep copy of the queue.
func (q *CircleQ[T]) Copy() CircleQ[T] {
	r := CircleQ[T]{
		data:  make([]T, q.Capacity()),
		begin: q.begin,
		end:   q.end,
	}
	copy(r.data, q.data)
	return r
}

// Push item to the queue if space is available.
func (q *CircleQ[T]) Push(item T) (ok bool) {
	data := q.data
	end := q.end
	if len(data) == 0 || ^end == 0 { // null, full ?
		return false
	}
	data[end] = item
	end++
	if int(end) >= len(data) {
		end = 0
	}
	if end == q.begin {
		end = ^uint32(0) // full
	}
	q.end = end
	return true
}

// Pop item from the queue if available.
func (q *CircleQ[T]) Pop() (item T, ok bool) {
	begin := q.begin
	end := q.end
	if end == begin { // empty ?
		return
	}
	data := q.data
	item = data[begin]
	// clean the slot: T may contain pointers, strings, slices, etc.
	var zero T
	data[begin] = zero
	if ^end == 0 { // full ?
		q.end = begin // not anymore
	}
	begin++
	if int(begin) >= len(data) {
		begin = 0
	}
	q.begin = begin
	return item, true
}
