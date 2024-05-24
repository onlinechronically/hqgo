package hqtrivia

// POST /verify-existing-phone
type VerificationStart struct {
	Method string `json:"method"`
	Phone  string `json:"phone"`
}

// POST /apiUrl/verifications/<verificationID>
type VerificationFinish struct {
	Code string `json:"code"`
}

// POST /usernames/available
type CheckUsername struct {
	Username string `json:"username"`
}

// POST /referral-code/valid
type CheckReferral struct {
	ReferralCode string `json:"referralCode"`
}

// POST /tokens
type RefreshTokens struct {
	Token string `json:"token"`
}

// POST /users
type RegisterAccount struct {
	Country           string `json:"country"`
	TimeZone          string `json:"timeZone"`
	Locale            string `json:"locale"`
	ReferringUsername string `json:"referringUsername,omitempty"`
	Username          string `json:"username"`
	VerificationID    string `json:"verificationId"`
}

// PATCH /users/me
type ChangeUsername struct {
	Username string `json:"username"`
}
