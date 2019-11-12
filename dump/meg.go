package dump

import (
	"fmt"
	"math"
)

// Meg prints out a number in human readable form, e.g. 20, 20K, 20M, 20G.
// Returns "NaN" if the input is not a number (int, float).
func Meg(n interface{}) string {
	const (
		kilo float64 = 1024.0
		mega float64 = 1024.0 * kilo
		giga float64 = 1024.0 * mega
		tera float64 = 1024.0 * giga
	)

	megInt64 := func(x float64) string {
		xAbs := math.Abs(x)
		switch {
		case xAbs > tera:
			return fmt.Sprintf("%.2fT", x/tera)
		case xAbs > giga:
			return fmt.Sprintf("%.2fG", x/giga)
		case xAbs > mega:
			return fmt.Sprintf("%.2fM", x/mega)
		case xAbs > kilo:
			return fmt.Sprintf("%.2fK", x/kilo)
		default:
			return fmt.Sprintf("%d", int(x))
		}
	}

	// in the order of likelihood
	switch t := n.(type) {
	case uint64:
		return megInt64(float64(t))
	case int64:
		return megInt64(float64(t))

	case int:
		return megInt64(float64(t))

	case float64:
		return megInt64(float64(t))
	case float32:
		return megInt64(float64(t))

	case int32:
		return megInt64(float64(t))
	case uint32:
		return megInt64(float64(t))

	case int16:
		return megInt64(float64(t))

	case uint16:
		return megInt64(float64(t))

	case int8:
		return megInt64(float64(t))
	case uint8: // also byte
		return megInt64(float64(t))

	default:
		return "NaN"
	}
}
