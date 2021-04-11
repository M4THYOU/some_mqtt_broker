package utils

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func checkIsIntInSlice(t *testing.T, i int, arr []int, expected bool) {
	res := IsIntInSlice(i, arr)
	if res != expected {
		t.Fatalf("TestIsIntInSlice expected %v, got %v for %d in %v", expected, res, i, arr)
	}
}
func TestIsIntInSlice(t *testing.T) {
	arr := []int{2, 3, 99}
	checkIsIntInSlice(t, 2, arr, true)
	checkIsIntInSlice(t, 3, arr, true)
	checkIsIntInSlice(t, 99, arr, true)
	checkIsIntInSlice(t, 4, arr, false)
	checkIsIntInSlice(t, 1, arr, false)
}

func checkReadBytesToSlice(t *testing.T, count int, buf, expected []byte, shouldPass bool) {
	rdr := bufio.NewReader(bytes.NewReader(buf))
	res, err := ReadBytesToSlice(count, rdr)
	if err != nil && shouldPass {
		t.Fatalf("failed to read bytes to slice: %v", err.Error())
	} else if err == nil && !shouldPass {
		t.Fatalf("should have failed to read bytes to slice: %v", res)
	} else if !cmp.Equal(res, expected) && shouldPass {
		t.Fatalf("Got:\n%v\nExpected:\n%v", res, expected)
	}
}
func TestReadBytesToSlice(t *testing.T) {
	buf := []byte{2, 6, 'M', 'm'}
	expected := buf
	checkReadBytesToSlice(t, 4, buf, expected, true)
	checkReadBytesToSlice(t, 5, buf, expected, false)
	expected = buf[:3]
	checkReadBytesToSlice(t, 3, buf, expected, true)
}
