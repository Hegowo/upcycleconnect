package models

import "math"

const CommissionRate = 0.10

func ComputeCommission(amountCents int64) (int64, int64) {
	commission := int64(math.Round(float64(amountCents) * CommissionRate))
	if commission > amountCents {
		commission = amountCents
	}
	return commission, amountCents - commission
}
