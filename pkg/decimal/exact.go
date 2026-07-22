package decimal

import "math/big"

// ExactMul multiplies at the sum of the operand scales so every finite product
// digit is retained. Scale-zero banker rounding on integer products is avoided
// by keeping at least one fractional place.
func ExactMul(left, right *Decimal) *Decimal {
	if left == nil || right == nil {
		return nil
	}

	scale := max(int64(1), left.GetScale()+right.GetScale())

	return left.SetScale(scale).Mul(right)
}

// ExactDiv divides at the greater of the operands' precision and DefaultScale
// so receiver-scale truncation cannot erase the finer operand.
func ExactDiv(left, right *Decimal) *Decimal {
	if left == nil || right == nil {
		return nil
	}

	scale := max(
		int64(DefaultScale),
		left.GetScale(),
		right.GetScale(),
	)

	return left.SetScale(scale).Div(right)
}

// ExactDivFloor divides and floors the quotient at the caller's executable
// scale so sell quantities never round upward past a funded boundary.
func ExactDivFloor(left, right *Decimal, scale int64) *Decimal {
	if left == nil || right == nil {
		return nil
	}

	workingScale := max(scale, left.GetScale(), right.GetScale())
	dividend := left.SetRounding(
		func(integer *big.Int, factor *big.Int) *big.Int {
			return new(big.Int).Quo(integer, factor)
		},
	).SetScale(workingScale)
	quotient := dividend.Div(right)

	return quotient.Copy().SetRounding(func(integer *big.Int, factor *big.Int) *big.Int {
		return new(big.Int).Quo(integer, factor)
	}).SetScale(scale)
}
