package packer

import (
	"reflect"
	"testing"
)

func TestSolve_Examples(t *testing.T) {
	packSizes := []int{250, 500, 1000, 2000, 5000}

	tests := []struct {
		name      string
		amount    int
		wantTotal int
		wantPacks map[string]int
	}{
		{
			name:      "Order 1",
			amount:    1,
			wantTotal: 250,
			wantPacks: map[string]int{"250": 1},
		},
		{
			name:      "Order 250",
			amount:    250,
			wantTotal: 250,
			wantPacks: map[string]int{"250": 1},
		},
		{
			name:      "Order 251",
			amount:    251,
			wantTotal: 500,
			wantPacks: map[string]int{"500": 1},
		},
		{
			name:      "Order 501",
			amount:    501,
			wantTotal: 750,
			wantPacks: map[string]int{"500": 1, "250": 1},
		},
		{
			name:      "Order 12001",
			amount:    12001,
			wantTotal: 12250,
			wantPacks: map[string]int{"5000": 2, "2000": 1, "250": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Solve(tt.amount, packSizes)
			if err != nil {
				t.Fatalf("Solve returned error: %v", err)
			}
			if res.ShippedTotal != tt.wantTotal {
				t.Errorf("ShippedTotal = %d, want %d", res.ShippedTotal, tt.wantTotal)
			}
			if !reflect.DeepEqual(res.Packs, tt.wantPacks) {
				t.Errorf("Packs = %v, want %v", res.Packs, tt.wantPacks)
			}
		})
	}
}

func TestSolve_InvalidInputs(t *testing.T) {
	_, err := Solve(0, []int{100, 200})
	if err == nil {
		t.Error("expected error for amount=0")
	}

	_, err = Solve(100, []int{})
	if err == nil {
		t.Error("expected error for empty sizes")
	}

	_, err = Solve(100, []int{-10, 0})
	if err == nil {
		t.Error("expected error for invalid sizes")
	}
}

func TestSolve_EdgeCase(t *testing.T) {
	// provided edge case
	sizes := []int{23,31,53}
	amount := 500000
	res, err := Solve(amount, sizes)
	if err != nil {
		t.Fatal(err)
	}
	// expected from the task
	expected := map[string]int{"23":2, "31":7, "53":9429}
	if res.ShippedTotal != 500000 {
		t.Fatalf("expected shipped total 500000, got %d", res.ShippedTotal)
	}
	if !reflect.DeepEqual(res.Packs, expected) {
		t.Fatalf("packs mismatch: expected %v got %v", expected, res.Packs)
	}
}