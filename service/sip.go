package service

import (
	"errors"
	"pluto/global"
	"pluto/model/params"
	"pluto/model/table"
)

type SipService struct{}

// CreateSip 新增线路
func (s *SipService) CreateSip(sip table.Sip) (mas table.Sip, err error) {
	err = global.GVA_DB.Create(&sip).Error
	if err != nil {
		return sip, err
	}
	return sip, err
}

// GetSipList 查询线路列表
func (s *SipService) GetSipList(info params.GetSipList) (list []table.Sip, total int64, err error) {
	db := global.GVA_DB.Model(&table.Sip{})
	if info.Mobile != "" {
		db = db.Where("mobile like ?", "%"+info.Mobile+"%")
	}
	if info.Name != "" {
		db = db.Where("name like ?", "%"+info.Name+"%")
	}
	err = db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = db.Limit(info.Limit).Offset(info.Offset).Order("created_at desc").Find(&list).Error
	return list, total, err
}

// GetSipInfoByID 根据id获取线路信息
func (s *SipService) GetSipInfoByID(id string) (sip table.Sip, err error) {
	err = global.GVA_DB.First(&sip, "id = ?", id).Error
	return sip, err
}

// SetSipInfo 设置线路信息
func (s *SipService) SetSipInfo(info params.SetSipInfo) (err error) {
	if info.ID == "" {
		err = errors.New("id is null")
		return
	}
	err = global.GVA_DB.Table("m_sip").Where("id = ?", info.ID).Updates(&info).Error
	return err
}

// DeleteSip 删除线路
func (s *SipService) DeleteSip(id string) (err error) {
	return global.GVA_DB.Where("id = ?", id).Delete(&table.Sip{}).Error
}

// GetUserSips 查询用户下已分配的线路
func (s *SipService) GetUserSips(info params.GetUserSipList) (list []table.UserSip, total int64, err error) {
	db := global.GVA_DB.Model(&table.UserSip{})
	if info.UserID != "" {
		db = db.Where("user_id = ?", info.UserID)
	}
	if info.Name != "" {
		db = db.Joins(`join m_sip on m_sip.id = m_user_sip.sip_id and m_sip.delete_at is null and m_sip like ?`, "%"+info.Name+"%")
	}
	err = db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = db.Limit(info.Limit).Offset(info.Offset).Order("created_at desc").
		Preload("SipInfo").
		Find(&list).Error
	return list, total, err
}

// UserAddSip 给用户分配线路
func (s *SipService) UserAddSip(userSip table.UserSip) (mas table.UserSip, err error) {
	err = global.GVA_DB.Create(&userSip).Error
	if err != nil {
		return userSip, err
	}
	return userSip, err
}

// SetUserSip 设置用户线路信息
func (s *SipService) SetUserSip(info params.SetUserSip) (err error) {
	if info.ID == "" {
		err = errors.New("id is null")
		return
	}
	err = global.GVA_DB.Table("m_user_sip").Where("id = ?", info.ID).Updates(&info).Error
	return err
}

// UserDelSip 给用户删除线路
func (s *SipService) UserDelSip(id string) (err error) {
	return global.GVA_DB.Where("id = ?", id).Delete(&table.UserSip{}).Error
}
