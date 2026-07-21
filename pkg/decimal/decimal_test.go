package decimal

import (
	"math"
	"math/big"
	"testing"
)

func TestNewFromString(t *testing.T) {
	if d, err := NewFromString("1.015"); err != nil {
		t.Error(err)
	} else if d.GetIncrement() != 1 {
		t.Errorf("d.GetIncrement() != 1, got %d", d.GetIncrement())
	} else if d.RawBigInt().Cmp(new(big.Int).SetInt64(1015)) != 0 {
		t.Errorf("d.RawBigInt() != 1015, got %s", d.integer)
	} else if d.GetScale() != 3 {
		t.Errorf("d.GetScale() != 3, got %d", d.GetScale())
	} else if d.String() != "1.015" {
		t.Errorf("d.String() != 1.015, got %s", d.String())
	}
}

func TestMath(t *testing.T) {
	if d, err := NewFromString("1.015"); err != nil {
		t.Error(err)
	} else if d = d.Add(NewFromInt64(1)); d.String() != "2.015" {
		t.Errorf("Add(1) != 2.015, got %s", d)
	} else if d = d.Sub(NewFromInt64(1)); d.String() != "1.015" {
		t.Errorf("Sub(1) != 1.015, got %s", d)
	} else if d = d.Mul(NewFromInt64(2)); d.String() != "2.030" {
		t.Errorf("Mul(2) != 2.030, got %s", d)
	} else if d = d.Div(NewFromInt64(2)); d.String() != "1.015" {
		t.Errorf("Div(2) != 1.015, got %s", d)
	} else if d = d.Pow(NewFromInt64(2)); d.String() != "1.030" {
		t.Errorf("Pow(2) != 1.030, got %s", d)
	} else if d = d.Add(NewFromInt64(-1)); d.String() != "0.030" {
		t.Errorf("Add(-1) != 0.030, got %s", d)
	}
}

func TestRounding(t *testing.T) {
	if d, err := NewFromString("1.002"); err != nil {
		t.Error(err)
	} else if d = d.SetIncrement(5); d.String() != "1.000" {
		t.Errorf("SetIncrement(5) != 1.000, got %s", d)
	} else if d = d.OffsetTicks(NewFromInt64(1)); d.String() != "1.005" {
		t.Errorf("OffsetTicks(1) != 1.005, got %s", d)
	} else if d = d.SetScale(2); d.String() != "1.00" {
		t.Errorf("SetScale(2) != 1.00, got %s", d)
	}
}

func TestEqualScaleMathPreservesOperands(t *testing.T) {
	left, err := NewFromString("12.3400")

	if err != nil {
		t.Fatal(err)
	}

	right, err := NewFromString("0.1200")

	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		calculate func() *Decimal
		expected  string
	}{
		{
			name:      "add",
			calculate: func() *Decimal { return left.Add(right) },
			expected:  "12.4600",
		},
		{
			name:      "subtract",
			calculate: func() *Decimal { return left.Sub(right) },
			expected:  "12.2200",
		},
		{
			name:      "multiply",
			calculate: func() *Decimal { return left.Mul(right) },
			expected:  "1.4808",
		},
		{
			name:      "divide",
			calculate: func() *Decimal { return left.Div(right) },
			expected:  "102.8333",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.calculate()

			if result.String() != test.expected {
				t.Fatalf("result = %s, want %s", result, test.expected)
			}

			if left.String() != "12.3400" {
				t.Fatalf("left operand mutated to %s", left)
			}

			if right.String() != "0.1200" {
				t.Fatalf("right operand mutated to %s", right)
			}
		})
	}
}

func TestMixedScaleMathUsesReceiverScale(t *testing.T) {
	receiver, err := NewFromString("1.234")

	if err != nil {
		t.Fatal(err)
	}

	operand, err := NewFromString("0.1")

	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		calculate func() *Decimal
		expected  string
	}{
		{
			name:      "add",
			calculate: func() *Decimal { return receiver.Add(operand) },
			expected:  "1.334",
		},
		{
			name:      "subtract",
			calculate: func() *Decimal { return receiver.Sub(operand) },
			expected:  "1.134",
		},
		{
			name:      "multiply",
			calculate: func() *Decimal { return receiver.Mul(operand) },
			expected:  "0.123",
		},
		{
			name:      "divide",
			calculate: func() *Decimal { return receiver.Div(operand) },
			expected:  "12.340",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := test.calculate().String(); result != test.expected {
				t.Fatalf("result = %s, want %s", result, test.expected)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	t.Run("exact integer range", func(t *testing.T) {
		value, err := NewFromString("61234.12345678")

		if err != nil {
			t.Fatal(err)
		}

		if difference := math.Abs(value.Float64() - 61234.12345678); difference > 1e-12 {
			t.Fatalf("Float64 difference = %g", difference)
		}
	})

	t.Run("large rational range", func(t *testing.T) {
		value, err := NewFromString("123456789012345678901.123")

		if err != nil {
			t.Fatal(err)
		}

		expected, _ := value.Rat().Float64()

		if value.Float64() != expected {
			t.Fatalf("Float64 = %g, want %g", value.Float64(), expected)
		}
	})
}
