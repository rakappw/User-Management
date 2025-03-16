package presenter

type LoginResponse struct {
	Token string `json:"token"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}
