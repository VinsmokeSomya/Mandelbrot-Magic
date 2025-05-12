package main

import (
	"time"
)

// xorshift random

// Initialize randState with the current Unix time in nanoseconds
var randState = uint64(time.Now().UnixNano())

// Generates a random uint64 using the xorshift algorithm
func RandUint64() uint64 {
	randState = ((randState ^ (randState << 13)) ^ (randState >> 7)) ^ (randState << 17)
	return randState
}

// Generates a random float64 between 0 and 1 using the xorshift algorithm
func RandFloat64() float64 {
	return float64(RandUint64()/2) / (1 << 63)
}
