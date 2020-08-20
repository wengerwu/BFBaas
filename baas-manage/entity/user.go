package entity

type User struct {
	Id       int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Userid   string `json:"userid" xorm:"not null" VARCHAR(100)`
	OpenId   string `json:"open_id" xorm:"not null" VARCHAR(50)`
	Account  string `json:"account" xorm:"not null unique VARCHAR(30)"`
	Password string `json:"password" xorm:"not null VARCHAR(100)"`
	Avatar   string `json:"avatar" xorm:"VARCHAR(200)"`
	Name     string `json:"name" xorm:"not null VARCHAR(20)"`
	Sex      int    `json:"sex" xorm:"not null INT(11)"`
	Phone    string `json:"phone" xorm:"not null VARCHAR(30)"`
	UserType int    `json:"user_type" xorm:"not null INT(11)"`
	Created  int64  `json:"created" xorm:"not null BIGINT(20)"`
}

type UserInfo struct {
	Userid       string   `json:"userid"`
	Roles        []string `json:"roles"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
	Name         string   `json:"name"`
	Account      string   `json:"account"`
}

type UserRole struct {
	UserId  string `json:"userId" xorm:"not null pk VARCHAR(100)"`
	RoleKey string `json:"roleKey" xorm:"not null VARCHAR(20)"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type UserDetail struct {
	User
	Roles []string `json:"roles"`
}

type UserIdAccount struct {
	OpenId  string `json:"open_id"`
	Account string `json:"account"`
}

type ChangePwdFrom struct {
	OldPassword     string `json:"oldPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}
