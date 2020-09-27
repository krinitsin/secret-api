package secret_api

// Response - базовая структура ответат от API secret_api
type Response struct {
	Vendor       string `json:"Vendor"`
	APIVersion   string `json:"APIVersion"`
	Code         int    `json:"ResponseCode"`
	ReplyMessage string `json:"ReplyMessage"`

	APMAC  string `json:"AP-MAC"`
	UserIP string `json:"UE-IP"`
	MAC    string `json:"UE-MAC"`
	SSID   string `json:"SSID"`
}

type AuthResp struct {
	Response
	Data string `json:"Data"`
}

type StatusResp struct {
	Response
	GuestUser       string `json:"GuestUser"`
	SmartClientInfo string `json:"SmartClientInfo"`
}
