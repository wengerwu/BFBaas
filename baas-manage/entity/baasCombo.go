package entity

type BaasCombo struct {
	Id          string  `json:"id" xorm:"not null VARCHAR(50)"`
	DisplayName string  `json:"displayName" xorm:"not null VARCHAR(50) 'DisplayName'"`
	Sort        int     `json:"sort" xorm:"not null int(11)"`
	Price       float64 `json:"price" xorm:"int(18,2)"`
	Remark      string  `json:"remark" xorm:"not null VARCHAR(500)"`
}
