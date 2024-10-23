package model

type ValidateTokenResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Data struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int    `json:"expiry"`
	TokenType    string `json:"token_type"`
}
