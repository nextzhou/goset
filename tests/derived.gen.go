package examples

import (
	"encoding/json"
	"fmt"
	"sort"
)

type IntSet struct {
	elements map[Int]struct{}
}

func NewIntSet(capacity int) *IntSet {
	set := new(IntSet)
	if capacity > 0 {
		set.elements = make(map[Int]struct{}, capacity)
	} else {
		set.elements = make(map[Int]struct{})
	}
	return set
}

func NewIntSetFromSlice(items []Int) *IntSet {
	set := NewIntSet(len(items))
	for _, item := range items {
		set.Put(item)
	}
	return set
}

func (set *IntSet) Extend(items ...Int) {
	for _, item := range items {
		set.Put(item)
	}
}

func (set *IntSet) Len() int {
	if set == nil {
		return 0
	}
	return len(set.elements)
}

func (set *IntSet) IsEmpty() bool {
	return set.Len() == 0
}

func (set *IntSet) ToSlice() []Int {
	if set == nil {
		return nil
	}
	s := make([]Int, 0, set.Len())
	set.ForEach(func(item Int) {
		s = append(s, item)
	})
	return s
}

func (set *IntSet) Put(key Int) {
	set.elements[key] = struct{}{}
}

func (set *IntSet) Clear() {
	set.elements = make(map[Int]struct{})
}

func (set *IntSet) Clone() *IntSet {
	cloned := NewIntSet(set.Len())
	for item := range set.elements {
		cloned.elements[item] = struct{}{}
	}
	return cloned
}

func (set *IntSet) Difference(another *IntSet) *IntSet {
	difference := NewIntSet(0)
	set.ForEach(func(item Int) {
		if !another.Contains(item) {
			difference.Put(item)
		}
	})
	return difference
}

func (set *IntSet) Equal(another *IntSet) bool {
	if set.Len() != another.Len() {
		return false
	}
	for item := range set.elements {
		if !another.Contains(item) {
			return false
		}
	}
	return true
}

func (set *IntSet) Intersect(another *IntSet) *IntSet {
	intersection := NewIntSet(0)
	if set.Len() < another.Len() {
		for item := range set.elements {
			if another.Contains(item) {
				intersection.Put(item)
			}
		}
	} else {
		for item := range another.elements {
			if set.Contains(item) {
				intersection.Put(item)
			}
		}
	}
	return intersection
}

func (set *IntSet) Union(another *IntSet) *IntSet {
	union := set.Clone()
	union.InPlaceUnion(another)
	return union
}

func (set *IntSet) InPlaceUnion(another *IntSet) {
	another.ForEach(func(item Int) {
		set.Put(item)
	})
}

func (set *IntSet) IsProperSubsetOf(another *IntSet) bool {
	return !set.Equal(another) && set.IsSubsetOf(another)
}

func (set *IntSet) IsProperSupersetOf(another *IntSet) bool {
	return !set.Equal(another) && set.IsSupersetOf(another)
}

func (set *IntSet) IsSubsetOf(another *IntSet) bool {
	if set.Len() > another.Len() {
		return false
	}
	for item := range set.elements {
		if !another.Contains(item) {
			return false
		}
	}
	return true
}

func (set *IntSet) IsSupersetOf(another *IntSet) bool {
	return another.IsSubsetOf(set)
}

func (set *IntSet) ForEach(f func(Int)) {
	if set.IsEmpty() {
		return
	}
	for item := range set.elements {
		f(item)
	}
}

func (set *IntSet) Filter(f func(Int) bool) *IntSet {
	result := NewIntSet(0)
	set.ForEach(func(item Int) {
		if f(item) {
			result.Put(item)
		}
	})
	return result
}

func (set *IntSet) Remove(key Int) {
	delete(set.elements, key)
}

func (set *IntSet) Contains(key Int) bool {
	_, ok := set.elements[key]
	return ok
}

func (set *IntSet) ContainsAny(keys ...Int) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set *IntSet) ContainsAll(keys ...Int) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}

func (set *IntSet) String() string {
	return fmt.Sprint(set.ToSlice())
}

func (set *IntSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *IntSet) UnmarshalJSON(b []byte) error {
	s := make([]Int, 0)
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*set = *NewIntSetFromSlice(s)
	return nil
}

type intOrderSet struct {
	elements        map[Int2]uint32
	elementSequence []Int2
}

func NewIntOrderSet(capacity int) *intOrderSet {
	set := new(intOrderSet)
	if capacity > 0 {
		set.elements = make(map[Int2]uint32, capacity)
		set.elementSequence = make([]Int2, 0, capacity)
	} else {
		set.elements = make(map[Int2]uint32)
	}
	return set
}

func NewIntOrderSetFromSlice(items []Int2) *intOrderSet {
	set := NewIntOrderSet(len(items))
	for _, item := range items {
		set.Put(item)
	}
	return set
}

func (set *intOrderSet) Extend(items ...Int2) {
	for _, item := range items {
		set.Put(item)
	}
}

func (set *intOrderSet) Len() int {
	if set == nil {
		return 0
	}
	return len(set.elements)
}

func (set *intOrderSet) IsEmpty() bool {
	return set.Len() == 0
}

func (set *intOrderSet) ToSlice() []Int2 {
	if set == nil {
		return nil
	}
	s := make([]Int2, set.Len())
	for idx, item := range set.elementSequence {
		s[idx] = item
	}
	return s
}

// NOTICE: efficient but unsafe
func (set *intOrderSet) ToSliceRef() []Int2 {
	return set.elementSequence
}

func (set *intOrderSet) Put(key Int2) {
	if _, ok := set.elements[key]; !ok {
		set.elements[key] = uint32(len(set.elementSequence))
		set.elementSequence = append(set.elementSequence, key)
	}
}

func (set *intOrderSet) Clear() {
	set.elements = make(map[Int2]uint32)
	set.elementSequence = set.elementSequence[:0]
}

func (set *intOrderSet) Clone() *intOrderSet {
	cloned := NewIntOrderSet(set.Len())
	for idx, item := range set.elementSequence {
		cloned.elements[item] = uint32(idx)
		cloned.elementSequence = append(cloned.elementSequence, item)
	}
	return cloned
}

func (set *intOrderSet) Difference(another *intOrderSet) *intOrderSet {
	difference := NewIntOrderSet(0)
	set.ForEach(func(item Int2) {
		if !another.Contains(item) {
			difference.Put(item)
		}
	})
	return difference
}

func (set *intOrderSet) Equal(another *intOrderSet) bool {
	if set.Len() != another.Len() {
		return false
	}
	for item := range set.elements {
		if !another.Contains(item) {
			return false
		}
	}
	return true
}

// TODO keep order
func (set *intOrderSet) Intersect(another *intOrderSet) *intOrderSet {
	intersection := NewIntOrderSet(0)
	if set.Len() < another.Len() {
		for item := range set.elements {
			if another.Contains(item) {
				intersection.Put(item)
			}
		}
	} else {
		for item := range another.elements {
			if set.Contains(item) {
				intersection.Put(item)
			}
		}
	}
	return intersection
}

func (set *intOrderSet) Union(another *intOrderSet) *intOrderSet {
	union := set.Clone()
	union.InPlaceUnion(another)
	return union
}

func (set *intOrderSet) InPlaceUnion(another *intOrderSet) {
	another.ForEach(func(item Int2) {
		set.Put(item)
	})
}

func (set *intOrderSet) IsProperSubsetOf(another *intOrderSet) bool {
	return !set.Equal(another) && set.IsSubsetOf(another)
}

func (set *intOrderSet) IsProperSupersetOf(another *intOrderSet) bool {
	return !set.Equal(another) && set.IsSupersetOf(another)
}

func (set *intOrderSet) IsSubsetOf(another *intOrderSet) bool {
	if set.Len() > another.Len() {
		return false
	}
	for item := range set.elements {
		if !another.Contains(item) {
			return false
		}
	}
	return true
}

func (set *intOrderSet) IsSupersetOf(another *intOrderSet) bool {
	return another.IsSubsetOf(set)
}

func (set *intOrderSet) ForEach(f func(Int2)) {
	if set.IsEmpty() {
		return
	}
	for _, item := range set.elementSequence {
		f(item)
	}
}

func (set *intOrderSet) Filter(f func(Int2) bool) *intOrderSet {
	result := NewIntOrderSet(0)
	set.ForEach(func(item Int2) {
		if f(item) {
			result.Put(item)
		}
	})
	return result
}

func (set *intOrderSet) Remove(key Int2) {
	if idx, ok := set.elements[key]; ok {
		l := set.Len()
		delete(set.elements, key)
		for ; idx < uint32(l-1); idx++ {
			item := set.elementSequence[idx+1]
			set.elementSequence[idx] = item
			set.elements[item] = idx
		}
		set.elementSequence = set.elementSequence[:l-1]
	}
}

func (set *intOrderSet) Contains(key Int2) bool {
	_, ok := set.elements[key]
	return ok
}

func (set *intOrderSet) ContainsAny(keys ...Int2) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set *intOrderSet) ContainsAll(keys ...Int2) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}

func (set *intOrderSet) DoUntil(f func(Int2) bool) int {
	for idx, item := range set.elementSequence {
		if f(item) {
			return idx
		}
	}
	return -1
}

func (set *intOrderSet) DoWhile(f func(Int2) bool) int {
	for idx, item := range set.elementSequence {
		if !f(item) {
			return idx
		}
	}
	return -1
}

func (set *intOrderSet) String() string {
	return fmt.Sprint(set.elementSequence)
}

func (set *intOrderSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *intOrderSet) UnmarshalJSON(b []byte) error {
	s := make([]Int2, 0)
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*set = *NewIntOrderSetFromSlice(s)
	return nil
}

type Int3Set struct {
	cmp             func(i, j Int3) bool
	elements        map[Int3]uint32
	elementSequence []Int3
}

func NewInt3Set(capacity int, cmp func(i, j Int3) bool) *Int3Set {
	set := new(Int3Set)
	if capacity > 0 {
		set.elements = make(map[Int3]uint32, capacity)
		set.elementSequence = make([]Int3, 0, capacity)
	} else {
		set.elements = make(map[Int3]uint32)
	}
	set.cmp = cmp
	return set
}

func NewInt3SetFromSlice(items []Int3, cmp func(i, j Int3) bool) *Int3Set {
	set := NewInt3Set(len(items), cmp)
	for _, item := range items {
		set.Put(item)
	}
	return set
}

func NewAscendingInt3Set(capacity int) *Int3Set {
	return NewInt3Set(capacity, func(i, j Int3) bool { return i < j })
}

func NewDescendingInt3Set(capacity int) *Int3Set {
	return NewInt3Set(capacity, func(i, j Int3) bool { return i > j })
}

func NewAscendingInt3SetFromSlice(items []Int3) *Int3Set {
	return NewInt3SetFromSlice(items, func(i, j Int3) bool { return i < j })
}

func NewDescendingInt3SetFromSlice(items []Int3) *Int3Set {
	return NewInt3SetFromSlice(items, func(i, j Int3) bool { return i > j })
}

func (set *Int3Set) Extend(items ...Int3) {
	for _, item := range items {
		set.Put(item)
	}
}

func (set *Int3Set) Len() int {
	if set == nil {
		return 0
	}
	return len(set.elements)
}

func (set *Int3Set) IsEmpty() bool {
	return set.Len() == 0
}

func (set *Int3Set) ToSlice() []Int3 {
	if set == nil {
		return nil
	}
	s := make([]Int3, 0, set.Len())
	set.ForEach(func(item Int3) {
		s = append(s, item)
	})
	return s
}

func (set *Int3Set) Put(key Int3) {
	if _, ok := set.elements[key]; !ok {
		idx := sort.Search(len(set.elementSequence), func(i int) bool {
			return set.cmp(key, set.elementSequence[i])
		})
		l := len(set.elementSequence)
		set.elementSequence = append(set.elementSequence, key)
		for i := l; i > idx; i-- {
			set.elements[set.elementSequence[i]] = uint32(i + 1)
			set.elementSequence[i] = set.elementSequence[i-1]
		}
		set.elements[set.elementSequence[idx]] = uint32(idx + 1)
		set.elementSequence[idx] = key
		set.elements[key] = uint32(idx)
	}
}

func (set *Int3Set) Clear() {
	set.elements = make(map[Int3]uint32)
	set.elementSequence = set.elementSequence[:0]
}

func (set *Int3Set) Clone() *Int3Set {
	cloned := NewInt3Set(set.Len(), set.cmp)
	for idx, item := range set.elementSequence {
		cloned.elements[item] = uint32(idx)
		cloned.elementSequence = append(cloned.elementSequence, item)
	}
	return cloned
}

func (set *Int3Set) Difference(another *Int3Set) *Int3Set {
	difference := NewInt3Set(0, set.cmp)
	set.ForEach(func(item Int3) {
		if !another.Contains(item) {
			difference.Put(item)
		}
	})
	return difference
}

func (set *Int3Set) Equal(another *Int3Set) bool {
	if set.Len() != another.Len() {
		return false
	}
	for item := range set.elements {
		if !another.Contains(item) {
			return false
		}
	}
	return true
}

func (set *Int3Set) Intersect(another *Int3Set) *Int3Set {
	intersection := NewInt3Set(0, set.cmp)
	if set.Len() < another.Len() {
		for item := range set.elements {
			if another.Contains(item) {
				intersection.Put(item)
			}
		}
	} else {
		for item := range another.elements {
			if set.Contains(item) {
				intersection.Put(item)
			}
		}
	}
	return intersection
}

func (set *Int3Set) Union(another *Int3Set) *Int3Set {
	union := set.Clone()
	union.InPlaceUnion(another)
	return union
}

func (set *Int3Set) InPlaceUnion(another *Int3Set) {
	another.ForEach(func(item Int3) {
		set.Put(item)
	})
}

func (set *Int3Set) IsProperSubsetOf(another *Int3Set) bool {
	return !set.Equal(another) && set.IsSubsetOf(another)
}

func (set *Int3Set) IsProperSupersetOf(another *Int3Set) bool {
	return !set.Equal(another) && set.IsSupersetOf(another)
}

func (set *Int3Set) IsSubsetOf(another *Int3Set) bool {
	if set.Len() > another.Len() {
		return false
	}
	for item := range set.elements {
		if !another.Contains(item) {
			return false
		}
	}
	return true
}

func (set *Int3Set) IsSupersetOf(another *Int3Set) bool {
	return another.IsSubsetOf(set)
}

func (set *Int3Set) ForEach(f func(Int3)) {
	if set.IsEmpty() {
		return
	}
	for _, item := range set.elementSequence {
		f(item)
	}
}

func (set *Int3Set) Filter(f func(Int3) bool) *Int3Set {
	result := NewInt3Set(0, set.cmp)
	set.ForEach(func(item Int3) {
		if f(item) {
			result.Put(item)
		}
	})
	return result
}

func (set *Int3Set) Remove(key Int3) {
	if idx, ok := set.elements[key]; ok {
		l := set.Len()
		delete(set.elements, key)
		for ; idx < uint32(l-1); idx++ {
			item := set.elementSequence[idx+1]
			set.elementSequence[idx] = item
			set.elements[item] = idx
		}
		set.elementSequence = set.elementSequence[:l-1]
	}
}

func (set *Int3Set) Contains(key Int3) bool {
	_, ok := set.elements[key]
	return ok
}

func (set *Int3Set) ContainsAny(keys ...Int3) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set *Int3Set) ContainsAll(keys ...Int3) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}

func (set *Int3Set) DoUntil(f func(Int3) bool) int {
	for idx, item := range set.elementSequence {
		if f(item) {
			return idx
		}
	}
	return -1
}

func (set *Int3Set) DoWhile(f func(Int3) bool) int {
	for idx, item := range set.elementSequence {
		if !f(item) {
			return idx
		}
	}
	return -1
}

func (set *Int3Set) String() string {
	return fmt.Sprint(set.elementSequence)
}

func (set *Int3Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *Int3Set) UnmarshalJSON(b []byte) error {
	return fmt.Errorf("unsupported")
}
