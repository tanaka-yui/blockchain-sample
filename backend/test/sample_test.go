package test

import (
	"fmt"
	"log"
	"testing"
)

// normal test
func Test_Exec(t *testing.T) {
	log.Println("Test_Exec Start")
}

// go benchmark sample
func Benchmark_Performance(b *testing.B) {
	log.Println("Benchmark_Performance Start")
	log.Println(b.N)
	b.ResetTimer()
	var base []string
	for i := 0; i < b.N; i++ {
		base = append(base, fmt.Sprintf("no%d", i))
	}
	log.Println("Benchmark_Performance End")
}
