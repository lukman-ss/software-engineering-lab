package currency

import (
	"fmt"
	"strconv"
	"strings"
)

// NaiveCurrency handles monetary values as float64 dollars.
// Fails floating point precision during parse/format roundtrips.
type NaiveCurrency struct{}

func (n NaiveCurrency) Format(dollars float64) string {
	return fmt.Sprintf("$%.2f", dollars)
}

func (n NaiveCurrency) Parse(s string) (float64, error) {
	s = strings.TrimPrefix(s, "$")
	return strconv.ParseFloat(s, 64)
}

// RobustAmount stores currency as integer cents to preserve exact values.
type RobustAmount struct {
	Cents int64
}

func NewRobustAmount(cents int64) RobustAmount {
	return RobustAmount{Cents: cents}
}

func (r RobustAmount) Format() string {
	absCents := r.Cents
	sign := ""
	if absCents < 0 {
		sign = "-"
		absCents = -absCents
	}
	dollars := absCents / 100
	cents := absCents % 100
	return fmt.Sprintf("%s$%d.%02d", sign, dollars, cents)
}

func ParseRobust(s string) (RobustAmount, error) {
	s = strings.TrimSpace(s)
	sign := int64(1)
	if strings.HasPrefix(s, "-") {
		sign = -1
		s = strings.TrimPrefix(s, "-")
	}
	if !strings.HasPrefix(s, "$") {
		return RobustAmount{}, fmt.Errorf("invalid format: missing $ prefix")
	}
	s = strings.TrimPrefix(s, "$")
	parts := strings.Split(s, ".")
	if len(parts) != 2 || len(parts[1]) != 2 {
		return RobustAmount{}, fmt.Errorf("invalid decimal format")
	}

	dollars, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return RobustAmount{}, err
	}
	cents, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return RobustAmount{}, err
	}

	total := (dollars*100 + cents) * sign
	return RobustAmount{Cents: total}, nil
}
