package models

type TwAccountLog struct {
	Id           int64  `json:"id"`
	UserId       string `json:"user_id"`
	Name         string `json:"name"`
	Account      string `json:"account"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiredAt    int64  `json:"expired_at"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

func (tLog *TwAccountLog) TableName() string {
	return "tw_account_log"
}
