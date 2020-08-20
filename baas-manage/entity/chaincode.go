package entity

type Chaincode struct {
	Id            int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	ChaincodeName string `json:"chaincodeName" xorm:"not null VARCHAR(64)"`
	ChainId       int    `json:"chain_id" xorm:"not null INT(11)"`
	ChannelId     int    `json:"channelId" xorm:"not null INT(11)"`
	OpenId        string `json:"open_id" xorm:"not null" VARCHAR(50)`
	Created       int64  `json:"created" xorm:"not null BIGINT(20)"`
	Version       string `json:"version" xorm:"VARCHAR(10)"`
	Status        int    `json:"status" xorm:"default 0 INT(11)"`
	GithubPath    string `json:"githubPath" xorm:"VARCHAR(256)"`
	Args          string `json:"args" xorm:"not null VARCHAR(500)"`
	Policy        string `json:"policy" xorm:"not null VARCHAR(200)"`
	Secret        string `json:"secret" xorm:"not null VARCHAR(50)"`
}
