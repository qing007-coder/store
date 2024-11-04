package request

type UpdatePersonalInfoReq struct {
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Gender       string `json:"gender"`
	Sign         string `json:"sign"`
}

type ModifyPasswordReq struct {
	NewPassword      string `json:"new_password"`
	VerificationCode string `json:"verification_code"`
}

type AddReceiverAddressReq struct {
	DetailedAddress string `json:"detailed_address"`
	Receiver        string `json:"receiver"`
	Phone           string `json:"phone"`
	Label           string `json:"label"`
}

type DeleteReceiverAddressReq struct {
	ID string `json:"id"`
}

type UpdateReceiverAddressReq struct {
	ID              string `json:"id"`
	DetailedAddress string `json:"detailed_address"`
	Receiver        string `json:"receiver"`
	Phone           string `json:"phone"`
	Label           string `json:"label"`
}

type AddFavouritesReq struct {
	TargetID string `json:"target_id"`
	Category string `json:"category"`
}

type DeleteFavouritesReq struct {
	TargetID string `json:"target_id"`
}

type AddFootprintReq struct {
	TargetID string `json:"target_id"`
	Category string `json:"category"`
}

type DeleteFootprintReq struct {
	TargetID string `json:"target_id"`
}

type FollowMerchantReq struct {
	MerchantID string `json:"merchant_id"`
}

type CancelFollowReq struct {
	MerchantID string `json:"merchant_id"`
}
