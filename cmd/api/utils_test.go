package main

import "testing"

func Test_SameSlices(t *testing.T) {
	a := []string{"foo", "bar", "sufurki"}
	b := []string{"foo", "bar", "sufurki"}

	if !compareSlices(a, b) {
		t.Error("slices should be equal")
	}
}

func Test_DifferentLengthSlices(t *testing.T) {
	a := []string{"foo", "bar", "sufurki"}
	b := []string{"foo", "bar"}

	if compareSlices(a, b) {
		t.Error("slices should be different")
	}
}

func Test_DifferentSlicesSameLength(t *testing.T) {
	a := []string{"foo", "bar", "sufurki"}
	b := []string{"foo", "bar", "not-sufurki"}

	if compareSlices(a, b) {
		t.Error("slices should be different")
	}
}

// TODO: test wrapper json
func Test_WrapperJson(t *testing.T) {

}

// TODO: test error json
func Test_ErrorJson(t *testing.T) {

}
