package entity

import (
	"fmt"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

type JsonTime time.Time

type BaasOrder struct {
	Id              string   `json:"id" xorm:"not null VARCHAR(50) 'Id'"`
	OrderNo         string   `json:"orderNo" xorm:"not null VARCHAR(50) 'OrderNo'"`
	PayOrderNo      string   `json:"payOrderNo" xorm:"VARCHAR(50) 'PayOrderNo'"`
	BaasComboId     string   `json:"baasComboId" xorm:"not null VARCHAR(50) 'BaasComboId'"`
	BaasComboName   string   `json:"baasComboName" xorm:"not null VARCHAR(50) 'BaasComboName'"`
	Body            string   `json:"body" xorm:"not null VARCHAR(200) 'Body'"`
	IsDelete        int      `json:"isDelete" xorm:"not null TINYINT(4) 'IsDelete'"`
	STATUS          int8     `json:"status" xorm:"not null tinyint 'STATUS'"`
	Price           float64  `json:"price" xorm:"not null decimal(18,2) 'Price'"`
	TagName         string   `json:"tagName" xorm:"not null VARCHAR(50) 'TagName'"`
	Question        string   `json:"question" xorm:"VARCHAR(2000) 'Question'"`
	Phone           string   `json:"phone" xorm:"VARCHAR(50) 'Phone'"`
	CreationTime    JsonTime `json:"creationTime" xorm:"not null DateTime 'CreationTime'"`
	CreatorUserId   string   `json:"creatorUserId" xorm:"not null varchar(50) 'CreatorUserId'"`
	CreatorUserName string   `json:"creatorUserName" xorm:"not null VARCHAR(50) 'CreatorUserName'"`
	DisposeTime     JsonTime `json:"disposeTime" xorm:"DateTime 'DisposeTime'"`
	DisposeUserId   string   `json:"disposeUserId" xorm:"VARCHAR(50) 'DisposeUserId'"`
	DisposeUserName string   `json:"disposeUserName" xorm:"VARCHAR(50) 'DisposeUserName'"`
}

func (jsonTime JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(jsonTime).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (jsonTime *JsonTime) UnmarshalJSON(data []byte) (err error) {
	newTime, err := time.ParseInLocation("\""+timeFormat+"\"", string(data), time.Local)
	*jsonTime = JsonTime(newTime)
	return
}
