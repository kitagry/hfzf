package main

import (
	"reflect"
	"testing"
)

func TestPointPlace(t *testing.T) {
	s1 := "aiueo"
	s2 := "auo"
	points := PointPlace(s1, s2)
	if !reflect.DeepEqual(points, []int{0, 2, 4}) {
		t.Fatal("PointPlace is something wrong")
	}
}
