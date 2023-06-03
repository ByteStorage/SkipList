// Copyright 2023 Qi Shen. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.
// Package myskiplist implements a skip list data structure.
//
// Skip lists are probabilistic data structures that provide.
// the same asymptotic complexity as balanced trees and are.
// simpler and faster in practice.
//
// See https://en.wikipedia.org/wiki/Skip_list for more information.
// This package provides a generic implementation of skip lists.
// The keys in the skip list can be of any type that implements.
// the sort.Interface interface.
// The following example shows how to use the package:

// package main
//
// import (
// 	"fmt"
// 	"github.com/qishenonly/SkipList"
// )
//
// func main() {
// 	// Create a new skip list
// 	list := skipist.New(skiplist.Int)
//
// 	// Insert elements into the skip list
// 	list.Insert(3, "three")
// 	list.Insert(1, "one")
// 	list.Insert(2, "two")
//
// 	// Get element from the skip list
// 	value, ok := list.Get(2)
// 	if ok {
// 		fmt.Println(value)
// 	}
//
// 	// Remove element from the skip list
// 	list.Remove(2)
//
// 	// Get element from the skip list
// 	value, ok = list.Get(2)
// 	if ok {
// 		fmt.Println(value)
// 	}
// }


package SkipList

import (
	"errors"
	"math/rand"
	"reflect"
	"sort"
	"strings"
)

// Default maximum level for the skip list
var DefaultMaxLevel = 48

// Node represents a node in the skip list
type node struct {
	key     interface{} // Key of the node
	value   interface{} // Value of the node
	forward []*node     // Forward pointers of the node
}

// SkipList represents the skip list structure
type SkipList struct {
	head    *node         // Head node of the skip list
	level   int           // Current level of the skip list
	length  int           // Length of the skip list (number of nodes)
	keyType reflect.Type  // Type of the keys in the skip list
}

// SkipListIterator represents the iterator for the skip list
type SkipListIterator struct {
	list   *SkipList  // The skip list associated with the iterator
	node   *node      // Current node being iterated
	isHead bool       // Flag to indicate if the current node is the head node
}

// NewSkipList creates a new skip list with the specified key type
func NewSkipList(keyType reflect.Type) *SkipList {
	head := &node{
		forward: make([]*node, 1),
	}
	return &SkipList{
		head:    head,
		level:   1,
		length:  0,
		keyType: keyType,
	}
}

// randomLevel generates a random level for the new node in the skip list
func (s *SkipList) randomLevel() int {
	level := 1
	for rand.Float64() < 0.5 && level < 32 {
		level++
	}
	return level
}

// compareInt compares two integers and returns the comparison result
func compareInt(a, b interface{}) int {
	keyA, ok := a.(int)
	if !ok {
		return 0
	}

	keyB, ok := b.(int)
	if !ok {
		return 0
	}

	if keyA < keyB {
		return -1
	} else if keyA > keyB {
		return 1
	} else {
		return 0
	}
}

// compareString compares two strings and returns the comparison result
func compareString(a, b interface{}) int {
	keyA, ok := a.(string)
	if !ok {
		return 0
	}

	keyB, ok := b.(string)
	if !ok {
		return 0
	}

	return strings.Compare(keyA, keyB)
}

// compare compares two keys and returns the comparison result
func (s *SkipList) compare(a, b interface{}) int {
	switch a := a.(type) {
	case int:
		b, ok := b.(int)
		if !ok {
			return 0
		}
		if a < b {
			return -1
		} else if a > b {
			return 1
		} else {
			return 0
		}
	case string:
		b, ok := b.(string)
		if !ok {
			return 0
		}
		return strings.Compare(a, b)
	default:
		return 0
	}
}

// Insert inserts a new key-value pair into the skip list
func (s *SkipList) Insert(key, value interface{}) error {
	if key == nil {
		return errors.New("Key cannot be nil")
	}

	update := make([]*node, s.level)
	current := s.head

	for i := s.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && s.compare(current.forward[i].key, key) < 0 {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]

	if current != nil && s.compare(current.key, key) == 0 {
		current.value = value
	} else {
		level := s.randomLevel()

		if level > s.level {
			for i := s.level; i < level; i++ {
				update[i] = s.head
			}
			s.level = level
		}

		newNode := &node{
			key:     key,
			value:   value,
			forward: make([]*node, level),
		}

		for i := 0; i < level; i++ {
			newNode.forward[i] = update[i].forward[i]
			update[i].forward[i] = newNode
		}

		s.length++
	}

	return nil
}

// Search searches for a key in the skip list and returns the corresponding value
func (s *SkipList) Search(key interface{}) (interface{}, error) {
	if key == nil {
		return nil, errors.New("Key cannot be nil")
	}

	current := s.head

	for i := s.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && s.compare(current.forward[i].key, key) < 0 {
			current = current.forward[i]
		}
	}

	current = current.forward[0]

	if current != nil && s.compare(current.key, key) == 0 {
		return current.value, nil
	}

	return nil, errors.New("Key not found")
}

// Delete deletes a key from the skip list
func (s *SkipList) Delete(key interface{}) error {
	if key == nil {
		return errors.New("Key cannot be nil")
	}

	update := make([]*node, s.level)
	current := s.head

	for i := s.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && s.compare(current.forward[i].key, key) < 0 {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]

	if current != nil && s.compare(current.key, key) == 0 {
		for i := 0; i < s.level; i++ {
			if update[i].forward[i] != current {
				break
			}
			update[i].forward[i] = current.forward[i]
		}

		for s.level > 1 && s.head.forward[s.level-1] == nil {
			s.level--
		}

		s.length--

		return nil
	}

	return errors.New("Key not found")
}

// Length returns the length of the skip list
func (s *SkipList) Length() int {
	return s.length
}

// Iterator returns a new iterator for the skip list
func (s *SkipList) Iterator() *SkipListIterator {
	return &SkipListIterator{
		list:   s,
		node:   s.head,
		isHead: true,
	}
}

// Next moves the iterator to the next node in the skip list and returns true if successful
func (it *SkipListIterator) Next() bool {
	if it.node.forward[0] != nil {
		it.node = it.node.forward[0]
		it.isHead = false
		return true
	}
	return false
}

// Key returns the key of the current node being iterated
func (it *SkipListIterator) Key() interface{} {
	if it.isHead {
		return nil
	}
	return it.node.key
}

// Value returns the value of the current node being iterated
func (it *SkipListIterator) Value() interface{} {
	if it.isHead {
		return nil
	}
	return it.node.value
}

// Clear Reset resets the iterator to the beginning of the skip list
func (s *SkipList) Clear() {
	s.head.forward = make([]*node, DefaultMaxLevel)
	s.level = 1
	s.length = 0
}


// MinString returns the minimum string key in the skip list,
// along with a boolean indicating if a key was found.
func (s *SkipList) MinString() (string, bool) {
	if s.length == 0 {
		return "", false
	}

	current := s.head.forward[0]
	minKey := ""
	for current != nil {
		if key, ok := current.key.(string); ok {
			if minKey == "" || key < minKey {
				minKey = key
			}
		}
		current = current.forward[0]
	}

	if minKey != "" {
		return minKey, true
	}

	return "", false
}

// MaxString returns the maximum string key in the skip list,
// along with a boolean indicating if a key was found.
func (s *SkipList) MaxString() (string, bool) {
	if s.length == 0 {
		return "", false
	}

	current := s.head.forward[0]
	maxKey := ""
	for current != nil {
		if key, ok := current.key.(string); ok {
			if maxKey == "" || key > maxKey {
				maxKey = key
			}
		}
		current = current.forward[0]
	}

	if maxKey != "" {
		return maxKey, true
	}

	return "", false
}

// MaxInt returns the maximum int key in the skip list,
// along with a boolean indicating if a key was found.
func (s *SkipList) MaxInt() (int, bool) {
	if s.length == 0 {
		return 0, false
	}

	current := s.head
	for i := s.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i] != s.head {
			current = current.forward[i]
		}
	}

	if key, ok := current.key.(int); ok {
		return key, true
	}

	return 0, false
}

// MinInt returns the minimum int key in the skip list,
// along with a boolean indicating if a key was found.
func (s *SkipList) MinInt() (int, bool) {
	if s.length == 0 {
		return 0, false
	}

	current := s.head.forward[0]
	for current != nil && current != s.head {
		if key, ok := current.key.(int); ok {
			return key, true
		}
		current = current.forward[0]
	}

	return 0, false
}

// SortByValue returns a slice of values in the skip list sorted by their values.
// If reverse is true, the values are sorted in descending order; otherwise,
// they are sorted in ascending order.
func (s *SkipList) SortByValue(reverse bool) []interface{} {
	if s.length == 0 {
		return nil
	}

	current := s.head.forward[0]
	result := make([]interface{}, 0, s.length)

	for current != nil {
		result = append(result, current.value)
		current = current.forward[0]
	}

	if reverse {
		sort.Slice(result, func(i, j int) bool {
			return compareKeysOrValues(result[j], result[i])
		})
	} else {
		sort.Slice(result, func(i, j int) bool {
			return compareKeysOrValues(result[i], result[j])
		})
	}

	return result
}

// SortByKey returns a slice of keys in the skip list sorted by their keys.
// If reverse is true, the keys are sorted in descending order; otherwise,
// they are sorted in ascending order.
func (s *SkipList) SortByKey(reverse bool) []interface{} {
	if s.length == 0 {
		return nil
	}

	current := s.head.forward[0]
	result := make([]interface{}, 0, s.length)

	for current != nil {
		result = append(result, current.key)
		current = current.forward[0]
	}

	if reverse {
		sort.Slice(result, func(i, j int) bool {
			return compareKeysOrValues(result[j], result[i])
		})
	} else {
		sort.Slice(result, func(i, j int) bool {
			return compareKeysOrValues(result[i], result[j])
		})
	}

	return result
}

// compareKeysOrValues compares two keys or values for sorting purposes.
// It supports comparison of int and string types.
func compareKeysOrValues(a, b interface{}) bool {
	intA, okA := a.(int)
	intB, okB := b.(int)
	if okA && okB {
		return intA < intB
	}

	strA, okA := a.(string)
	strB, okB := b.(string)
	if okA && okB {
		return strA < strB
	}

	return false
}
