package service

import (
	"errors"
	"pluto/global"
	"pluto/middleware/db"
	"pluto/model/constant"
	"pluto/model/params"
	"pluto/model/table"
)

type RechargeService struct{}

// Recharge 充值
func (r *RechargeService) Recharge(info params.CreateRecharge) (err error) {
	var wallet table.Wallet
	if global.GVA_DB.Where("user_id = ?", info.UserID).First(&wallet).RecordNotFound() {
		return errors.New("钱包未找到")
	}

	tx := db.Begin(global.GVA_DB)
	defer tx.RollbackIfFailed()

	// 给钱包充值
	wallet.Balance += info.Amount
	err = tx.Save(&wallet).Error
	if err != nil {
		return err
	}

	// 保存充值记录
	recharge := table.Recharge{
		UserID:  info.UserID,
		Type:    info.Type,
		Amount:  info.Amount,
		Balance: wallet.Balance,
		Status:  constant.NormalStatus,
	}
	err = tx.Create(&recharge).Error
	if err != nil {
		return err
	}

	tx.Commit()
	return err
}

// GetRechargeList 获取充值记录列表
func (r *RechargeService) GetRechargeList(info params.GetRechargeList) (list []table.Recharge, total int64, err error) {
	db := global.GVA_DB.Model(&table.Recharge{})
	if info.Type != "" {
		db = db.Where("type = ?", info.Type)
	}
	if info.Mobile != "" {
		db = db.Joins(` join sys_user on sys_user.id = m_recharge.user_id and sys_user.mobile like ?`, "%"+info.Mobile+"%")
	}
	err = db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = db.Limit(info.Limit).Offset(info.Offset).Order("created_at desc").
		Preload("UserInfo").
		Find(&list).Error
	return list, total, err
}
