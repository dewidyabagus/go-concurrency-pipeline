package main

import "testing"

func BenchmarkProceed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Proceed()
	}
}

func BenchmarkProceedConcurrency(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ProceedConcurrency()
	}
}
