package table

import (
	"pluto/global"
	"time"
)

// MonthlyBill 月账单
type MonthlyBill struct {
	global.UUID
	UserID      string    `gorm:"index" json:"userID"`
	Date        time.Time `json:"date"`
	DurationMin int       `json:"durationMin"`
	Cost        int       `json:"cost"`
	Count       int       `json:"count"`
}

func (MonthlyBill) TableName() string {
	return "m_month_bill"
}
