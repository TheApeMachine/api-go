package decimal

import "testing"

func BenchmarkDecimalAddEqualScale(b *testing.B) {
	left, _ := NewFromString("12345.6789")
	right, _ := NewFromString("0.1234")

	b.ReportAllocs()

	for b.Loop() {
		_ = left.Add(right)
	}
}

func BenchmarkDecimalAddMixedScale(b *testing.B) {
	left, _ := NewFromString("12345.6789")
	right, _ := NewFromString("0.12")

	b.ReportAllocs()

	for b.Loop() {
		_ = left.Add(right)
	}
}

func BenchmarkDecimalAddPreScaledOperands(b *testing.B) {
	left, _ := NewFromString("12345.6789")
	right, _ := NewFromString("0.12")

	b.ReportAllocs()

	for b.Loop() {
		_ = left.SetScale(4).Add(right.SetScale(4))
	}
}

func BenchmarkDecimalSubEqualScale(b *testing.B) {
	left, _ := NewFromString("12345.6789")
	right, _ := NewFromString("0.1234")

	b.ReportAllocs()

	for b.Loop() {
		_ = left.Sub(right)
	}
}

func BenchmarkDecimalMulEqualScale(b *testing.B) {
	left, _ := NewFromString("12345.6789")
	right, _ := NewFromString("0.1234")

	b.ReportAllocs()

	for b.Loop() {
		_ = left.Mul(right)
	}
}

func BenchmarkDecimalDivEqualScale(b *testing.B) {
	left, _ := NewFromString("12345.6789")
	right, _ := NewFromString("0.1234")

	b.ReportAllocs()

	for b.Loop() {
		_ = left.Div(right)
	}
}

func BenchmarkDecimalFloat64(b *testing.B) {
	value, _ := NewFromString("61234.12345678")

	b.ReportAllocs()

	for b.Loop() {
		_ = value.Float64()
	}
}
