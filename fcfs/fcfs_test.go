package fcfs

import "testing"

func TestFirstComeFirstServe(t *testing.T) {
	msg := FirstComeFirstServe()
	if msg == "" {
		t.Fatalf("FirstComeFirstServe() = %s, want an non empty string", msg)
	}
}
