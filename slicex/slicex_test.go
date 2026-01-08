package slicex

import "testing"

func TestContains(t *testing.T) {
	if !Contains([]int{1, 2, 3}, 2) {
		t.Fatal("expected true")
	}
	if Contains([]string{"a", "b"}, "c") {
		t.Fatal("expected false")
	}
}

func TestIndexOf(t *testing.T) {
	if idx := IndexOf([]int{5, 6, 7}, 6); idx != 1 {
		t.Fatalf("expected 1, got %d", idx)
	}
	if idx := IndexOf([]int{5, 6, 7}, 8); idx != -1 {
		t.Fatalf("expected -1, got %d", idx)
	}
}

func TestFilter(t *testing.T) {
	res := Filter([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 })
	if len(res) != 2 || res[0] != 2 || res[1] != 4 {
		t.Fatalf("unexpected result: %v", res)
	}
}

func TestMap(t *testing.T) {
	res := Map([]int{1, 2, 3}, func(v int) string { return string(rune('a' + v - 1)) })
	expected := []string{"a", "b", "c"}
	for i := range expected {
		if res[i] != expected[i] {
			t.Fatalf("unexpected result: %v", res)
		}
	}
}
