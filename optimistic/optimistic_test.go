package optimistic

import (
	"testing"
)

func BenchmarkUpdate(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Update()
		}

	})
}
