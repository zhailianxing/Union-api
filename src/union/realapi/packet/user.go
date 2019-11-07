package packet

type Session struct {
	UserId      uint32 `json:"user_id"`
	UserName    string `json:"user_name"`
	UserPhone   string `json:"user_phone"`
	TokenId     []byte `json:"token_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	AccType     int    `json:"acc_type"`
}

type User struct {
	UserId      uint32
	UserName    string
	UserPhone   string
	DataSignup  int64
	DisplayName string
	Email       string
	AccType     int
}
