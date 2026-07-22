package decimal

import "testing"

func TestExactMulPreservesProductDigits(t *testing.T) {
	price, err := NewFromString("0.00012345")

	if err != nil {
		t.Fatal(err)
	}

	quantity, err := NewFromString("0.00000067")

	if err != nil {
		t.Fatal(err)
	}

	if got := ExactMul(price, quantity).String(); got != "0.0000000000827115" {
		t.Fatalf("ExactMul = %s", got)
	}
}

func TestExactDivPreservesFinerOperand(t *testing.T) {
	left, err := NewFromString("1.2")

	if err != nil {
		t.Fatal(err)
	}

	right, err := NewFromString("1.234")

	if err != nil {
		t.Fatal(err)
	}

	// Receiver-scale Div would truncate right to scale 1 before dividing.
	if got := ExactDiv(left, right).String(); got != "0.972447325770" {
		t.Fatalf("ExactDiv = %s", got)
	}
}

func TestExactDivFloorDoesNotExceedBudget(t *testing.T) {
	capacity, _ := NewFromString("1.000")
	mark, _ := NewFromString("0.333")
	quantity := ExactDivFloor(capacity, mark, 3)

	if ExactMul(mark, quantity).Cmp(capacity) > 0 {
		t.Fatalf("floored quantity %s exceeds capacity at mark", quantity)
	}
}
