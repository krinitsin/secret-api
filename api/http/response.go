package http

type Response struct {
	Userip string `json:"user_ip,omitempty"`
	Mac    string `json:"mac,omitempty"`
}
