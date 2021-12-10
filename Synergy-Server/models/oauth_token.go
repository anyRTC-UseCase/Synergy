package models

type OauthToken struct {
	Token    string
	UserId   string
	Secret   string
	ExpireAt int64
	CreateAt int64
	Revoked  bool
}

type Token struct {
	Token string `json:"access_token"`
}

func GetOauthTokenByToken(token string) (ot *OauthToken) {
	ot = new(OauthToken)
	return
}
