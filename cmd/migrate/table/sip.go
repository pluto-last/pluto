package table

import (
	"pluto/global"
	"time"
)

// Sip 线路表
type Sip struct {
	global.UUID
	Name        string `json:"name"`        // 线路名称
	Mobile      string `json:"mobile"`      // 线路显示号码
	IntervalSec int    `json:"intervalSec"` // 呼叫间隔
	SipIP       string `json:"sipIP"`       // 线路IP
	SipPort     string `json:"sipPort"`     // 线路端口
	Note        string `json:"note"`        //备注
}

func (Sip) TableName() string {
	return "m_sip"
}

// UserSip 用户线路表
type UserSip struct {
	global.UUID
	UserID     string    `json:"userID"`
	SipID      string    `json:"sipID"`
	ExpireAt   time.Time `json:"expireAt"`   // 过期时间
	Concurrent int       `json:"concurrent"` // 并发数
	Price      int       `json:"price"`      // 单价 每分钟
	SipInfo    Sip       `gorm:"ForeignKey:SipID" json:"sipInfo"`
}

func (UserSip) TableName() string {
	return "m_user_sip"
}
