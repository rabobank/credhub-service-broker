package model

type UaaTokenIntrospection struct {
	UserId    string   `json:"user_id"`
	UserName  string   `json:"user_name"`
	Email     string   `json:"email"`
	ClientId  string   `json:"client_id"`
	Exp       int      `json:"exp"`
	Scope     []string `json:"scope"`
	Jti       string   `json:"jti"`
	Aud       []string `json:"aud"`
	Sub       string   `json:"sub"`
	Iss       string   `json:"iss"`
	Iat       int      `json:"iat"`
	Cid       string   `json:"cid"`
	GrantType string   `json:"grant_type"`
	Azp       string   `json:"azp"`
	AuthTime  int      `json:"auth_time"`
	Zid       string   `json:"zid"`
	RevSig    string   `json:"rev_sig"`
	Origin    string   `json:"origin"`
	Revocable bool     `json:"revocable"`
	Active    bool     `json:"active"`
}
