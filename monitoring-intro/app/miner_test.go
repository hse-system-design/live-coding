package main

import "testing"

func TestHasLeadingZeroBits(t *testing.T) {
	var array [32]byte
	//slice := array[:]

	if !hasLeadingZeroBits(array, 32*8) {
		t.FailNow()
	}
	array[31] = 1
	if hasLeadingZeroBits(array, 32*8) {
		t.FailNow()
	}
	if !hasLeadingZeroBits(array, 32*8-1) {
		t.FailNow()
	}

	array[0] = 2
	if hasLeadingZeroBits(array, 7) {
		t.FailNow()
	}
	if !hasLeadingZeroBits(array, 6) {
		t.FailNow()
	}
}

func BenchmarkGeneratePair_4096_2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mineKey(4096, 2)
	}
}

func BenchmarkGeneratePair_4096_8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mineKey(4096, 8)
	}
}

func BenchmarkGeneratePair_4096_10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mineKey(4096, 10)
	}
}

func BenchmarkGeneratePair_4096_12(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mineKey(4096, 12)
	}
}

func BenchmarkGeneratePair_4096_13(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mineKey(4096, 13)
	}
}

func BenchmarkGeneratePair_4096_16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mineKey(4096, 16)
	}
}
