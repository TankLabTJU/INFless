package types

import "math"

const ACC = 0.000001
func Equal(a, b float64) bool {
	return math.Abs(a-b) < ACC
}

func Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > ACC
}

func Less(a, b float64) bool {
	return math.Max(a, b) == b && math.Abs(a-b) > ACC
}

func GreaterEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < ACC
}

func  LessEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < ACC
}