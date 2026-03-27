package stack

import "testing"

func benchmarkSizes() []struct {
	name string
	size int
} {
	return []struct {
		name string
		size int
	}{
		{"n=100", 100},
		{"n=1000", 1000},
		{"n=10000", 10000},
	}
}

// Prevent compiler optimization of benchmark results.
var benchResult any

func BenchmarkPush(b *testing.B) {
	for _, tc := range benchmarkSizes() {
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				var s Stack[int]
				for i := range tc.size {
					s.Push(i)
				}
			}
		})
	}
}

func BenchmarkPop(b *testing.B) {
	for _, tc := range benchmarkSizes() {
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				b.StopTimer()

				var s Stack[int]
				for i := range tc.size {
					s.Push(i)
				}

				b.StartTimer()

				for range tc.size {
					s.Pop()
				}
			}
		})
	}
}

func BenchmarkPushLeft(b *testing.B) {
	for _, tc := range benchmarkSizes() {
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				var s Stack[int]
				for i := range tc.size {
					s.PushLeft(i)
				}
			}
		})
	}
}

func BenchmarkPeek(b *testing.B) {
	var s Stack[int]
	for i := range 1000 {
		s.Push(i)
	}

	b.ResetTimer()

	for range b.N {
		benchResult = s.Peek()
	}
}

func BenchmarkReverse(b *testing.B) {
	for _, tc := range benchmarkSizes() {
		b.Run(tc.name, func(b *testing.B) {
			var s Stack[int]
			for i := range tc.size {
				s.Push(i)
			}

			b.ResetTimer()

			for range b.N {
				s.Reverse()
			}
		})
	}
}
