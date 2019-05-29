package main

import (
	"reflect"
	"testing"
)

func TestPointPlace(t *testing.T) {
	s1 := "aiueo"
	s2 := "auo"
	points := pointPlace(s1, s2)
	if !reflect.DeepEqual(points, []int{0, 2, 4}) {
		t.Fatal("pointPlace is something wrong")
	}
}
