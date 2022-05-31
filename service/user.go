package service

import (
	"errors"
	"pluto/global"
	"pluto/middleware/db"
	"pluto/model/constant"
	"pluto/model/params"
	"pluto/model/table"
	"pluto/utils"

	"github.com/jinzhu/gorm"
)

type UserService struct{}

// Login 用户登录
func (userService *UserService) Login(u *table.User) (err error, userInter *table.User) {
	var user table.User
	u.Password = utils.MD5V([]byte(u.Password))
	err = global.GVA_DB.Where("mobile = ? AND password = ?", u.Mobile, u.Password).First(&user).Error
	return err, &user
}

// Register 用户注册
func (userService *UserService) Register(u table.User) (userInter table.User, err error) {
	var user table.User
	if !errors.Is(global.GVA_DB.Where("mobile = ?", u.Mobile).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("手机号码已经已注册")
	}

	// 开启事务
	tx := db.Begin(global.GVA_DB)
	defer tx.RollbackIfFailed()
	err = tx.Create(&u).Error
	if err != nil {
		return u, err
	}

	wallet := table.Wallet{
		UserID:  u.ID,
		Balance: 0,
	}
	err = tx.Create(&wallet).Error
	tx.Commit()
	return u, err
}

// BatchCreateUser 批量生成用户
func (userService *UserService) BatchCreateUser(count int) (err error) {
	// 开启事务
	tx := db.Begin(global.GVA_DB)
	defer tx.RollbackIfFailed()
	for i := 0; i < count; i++ {
		mobile := GetActiveMobile()
		user := table.User{
			UserName: mobile,
			Mobile:   mobile,
			Password: utils.RandStr(8),
			Status:   constant.NormalStatus,
		}

		err = global.GVA_DB.Create(&user).Error
		if err != nil {
			return
		}

		wallet := table.Wallet{
			UserID:  user.ID,
			Balance: 0,
		}
		err = global.GVA_DB.Create(&wallet).Error
		if err != nil {
			return
		}
	}
	tx.Commit()
	return err
}

// GetActiveMobile 获取可用的手机号码
func GetActiveMobile() string {
	mobile := utils.GetRandomTel()
	var user table.User
	if global.GVA_DB.Where("mobile = ?", mobile).First(&user).RecordNotFound() {
		return mobile
	}
	return GetActiveMobile()
}

// ChangePassword 用户修改密码
func (userService *UserService) ChangePassword(userID, newPassword string) (err error) {
	return global.GVA_DB.Table("sys_user").Where("id = ?", userID).Update("password", newPassword).Error
}

// GetUserInfoList 查询用户列表
func (userService *UserService) GetUserInfoList(info params.GetUserList) (list []table.User, total int64, err error) {
	db := global.GVA_DB.Model(&table.User{})
	if info.Mobile != "" {
		db = db.Where("mobile like ?", "%"+info.Mobile+"%")
	}
	if info.UserName != "" {
		db = db.Where("user_name like ?", "%"+info.UserName+"%")
	}
	var userList []table.User
	err = db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = db.Limit(info.Limit).Offset(info.Offset).Order("created_at desc").Find(&userList).Error
	return userList, total, err
}

// GetUserInfoByID 根据id获取用户信息
func (userService *UserService) GetUserInfoByID(id string) (user table.User, err error) {
	err = global.GVA_DB.First(&user, "id = ?", id).Preload("WalletInfo").Error
	return user, err
}

// GetUserInfoByMobile 根据mobile获取用户信息
func (userService *UserService) GetUserInfoByMobile(mobile string) (user table.User, err error) {
	err = global.GVA_DB.First(&user, "mobile = ?", mobile).Error
	return user, err
}

// SetUserInfo 设置用户信息
func (userService *UserService) SetUserInfo(info params.SetUserInfo) (err error) {
	if info.ID == "" {
		err = errors.New("id is null")
		return
	}
	err = global.GVA_DB.Table("sys_user").Where("id = ?", info.ID).Updates(&info).Error
	return err
}

// DeleteUser 删除用户
func (userService *UserService) DeleteUser(id string) (err error) {
	return global.GVA_DB.Where("id = ?", id).Delete(&table.User{}).Error
}
