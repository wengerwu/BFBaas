package entity

type ReqToken struct {
	ReqToken  string `json:"req_token"`
	LoginName string `json:"account"`
	Password  string `json:"password"`
}

type LoginFrom struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	VerificationCode string `json:"verificationCode"`
	LoginNum         int    `json:"loginNum"`
	Timestamp        int64  `json:"timestamp"`
	Sign             string `json:"sign"`
}

type LoginPhoneFrom struct {
	Phone     string `json:"phone"`
	OldPhone  string `json:"oldPhone"`
	SmsPhone  string `json:"smsPhone"`
	SmsCode   string `json:"smsCode"`
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
}

type SendSmsFrom struct {
	Phone     string `json:"phone"`
	SmsPhone  string `json:"SmsPhone"`
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
}

type RegisterPhoneFrom struct {
	Phone     string `json:"phone"`
	OldPhone  string `json:"oldPhone"`
	SmsPhone  string `json:"smsPhone"`
	SmsCode   string `json:"smsCode"`
	PassWord  string `json:"passWord"`
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
}

type NotPasswordTokenFrom struct {
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
}

type UserTokenInfo struct {
	UserId  string
	Account string
	Token   string
}

type ResponseResult struct {
	Result              interface{} `json:"result"`
	TargetUrl           string      `json:"targetUrl"`
	Success             bool        `json:"success"`
	Error               Error       `json:"error"`
	UnAuthorizedRequest bool        `json:"unAuthorizedRequest"`
	Abp                 bool        `json:"__abp"`
}

type Error struct {
	Code             int         `json:"code"`
	Message          string      `json:"message"`
	Details          interface{} `json:"details"`
	ValidationErrors interface{} `json:"validationErrors"`
}

type Token struct {
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

type JWTToken struct {
	Nbf      int64    `json:"nbf"`
	Exp      int64    `json:"exp"`
	Iss      string   `json:"iss"`
	Aud      []string `json:"aud"`
	ClientId string   `json:"client_id"`
	Sub      string   `json:"sub"`
	AuthTime int64    `json:"auth_time"`
	Idp      string   `json:"idp"`
	Account  string   `json:"account"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	Sex      string   `json:"sex"`
	Photo    string   `json:"Photo"`
	Scope    []string `json:"scope"`
	Amr      []string `json:"amr"`
}
