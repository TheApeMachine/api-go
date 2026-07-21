package decimal

import "math/big"

// SetScale returns m with adjusted decimal places.
func (d *Decimal) SetScale(scale int64) *Decimal {
	if scale == d.scale {
		return d.Copy()
	}

	return &Decimal{
		integer:   d.integerAtScale(scale),
		scale:     scale,
		increment: d.increment,
		rounding:  d.rounding,
	}
}

// integerAtScale returns d's immutable integer directly when scales match and
// otherwise builds only the adjusted integer required by the operation.
func (d *Decimal) integerAtScale(scale int64) *big.Int {
	if d.scale == scale {
		return d.integer
	}

	result := new(big.Int).Set(d.integer)
	difference := scale - d.scale

	if result.Sign() == 0 {
		return result
	}

	absoluteDifference := difference

	if absoluteDifference < 0 {
		absoluteDifference = -absoluteDifference
	}

	factor := new(big.Int).Exp(
		big.NewInt(10),
		big.NewInt(absoluteDifference),
		nil,
	)

	if difference > 0 {
		result.Mul(result, factor)
		return d.roundIntegerToGranularity(result)
	}

	result = d.rounding(result, factor)

	return d.roundIntegerToGranularity(result)
}

// ScalingFactor returns 10 ^ decimals in [big.Int].
func (d *Decimal) ScalingFactor() *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(d.scale), nil)
}

// roundToGranularity returns the rounding of m to the granularity constraint.
func (d *Decimal) roundToGranularity() {
	d.integer = d.roundIntegerToGranularity(d.integer)
}

// roundIntegerToGranularity applies d's increment without requiring a
// temporary Decimal while arithmetic aligns an operand's integer.
func (d *Decimal) roundIntegerToGranularity(integer *big.Int) *big.Int {
	if d.increment <= 1 {
		return integer
	}

	tick := big.NewInt(d.increment)
	remainder := new(big.Int).Mod(integer, tick)
	half := new(big.Int).Div(tick, big.NewInt(2))
	remainderCmpHalf := remainder.Cmp(half)

	if remainderCmpHalf > 0 {
		return integer.Sub(integer, remainder).Add(integer, tick)
	}

	if remainderCmpHalf < 0 {
		return integer.Sub(integer, remainder)
	}

	roundedDown := new(big.Int).Sub(integer, remainder)
	quotient := new(big.Int).Div(roundedDown, tick)

	if quotient.Bit(0) == 0 {
		return roundedDown
	}

	return roundedDown.Add(roundedDown, tick)
}
