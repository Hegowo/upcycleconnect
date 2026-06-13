package models

import "math"

// CommissionRate is the platform's cut on each paid prestation transaction.
// The spec (Annexe 1) allows 5-10%; we apply 10%.
const CommissionRate = 0.10

// ComputeCommission returns (commissionCents, netCents) for a gross amount.
// Commission is rounded to the nearest cent; net is the remainder owed to the provider.
func ComputeCommission(amountCents int64) (int64, int64) {
	commission := int64(math.Round(float64(amountCents) * CommissionRate))
	if commission > amountCents {
		commission = amountCents
	}
	return commission, amountCents - commission
}
