package domain

type Referral struct {
	ReferrerID     *string
	Balance        int
	CommissionRate int
	Count          int
}
