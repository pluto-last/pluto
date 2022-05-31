package api

import (
	"pluto/global"
	"pluto/middleware/auth"
	"pluto/middleware/sms"
	"pluto/model/constant"
	"pluto/utils"
	"strconv"
	"time"

	JWT "pluto/middleware/jwt"
	"pluto/model/params"
	"pluto/model/reply"
	"pluto/model/table"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserCtl struct{}

// @Tags User
// @Summary 用户登录
// @Produce  application/json
// @Param body body params.Login true "用户名, 密码"
// @Success 200
// @Router /user/login [post]
func (u *UserCtl) Login(c *gin.Context) {
	var l params.Login
	// 参数绑定
	_ = c.Bind(&l)

	// 参数校验
	if err := utils.Verify(l, utils.LoginVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}

	rkey := "login_limit:" + l.Mobile
	loginLimit, _ := global.GVA_REDIS.GetInt(rkey)
	if loginLimit >= 5 {
		global.GVA_LOG.Info("用户登陆失败次数过多", zap.String("mobile", l.Mobile))
		reply.FailWithMessage("达到最大重试次数", c)
		return
	}

	user, err := userService.GetUserInfoByMobile(l.Mobile)
	if err != nil {
		reply.FailWithMessage("用户未找到", c)
		return
	}
	if auth.VerifyPasswd(user.Password, l.Password) {
		//if res, _ := auth.VerifyPasswd(user.Password, l.Password); !res {
		// 失败次数限制
		if global.GVA_REDIS.IsExist(rkey) {
			global.GVA_REDIS.GetSeqence(rkey)
		} else {
			global.GVA_REDIS.SetStringExpire(rkey, "1", time.Minute*10)
		}

		reply.FailWithMessage("密码错误！", c)
		return
	}

	if user.Status != constant.NormalStatus {
		reply.FailWithMessage("帐号被禁用！", c)
		return
	}

	if l.Source == "admin" && !roleAuthService.HasRole(user.ID, "Admin") {
		reply.FailWithMessage("您还不是管理员，没有权限访问！", c)
		return
	}

	u.tokenNext(c, user)

}

// 登录以后签发jwt
func (u *UserCtl) tokenNext(c *gin.Context, user table.User) {
	j := &JWT.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SignKey)} // 唯一签名
	claims := params.CustomClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,                                                    // 签名生效时间
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(utils.GetLoginTimeOut())).Unix(), // 每天0点过期
			Issuer:    "pluto",                                                                     // 签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GVA_LOG.Error("获取token失败!", zap.Any("err", err))
		reply.FailWithMessage("获取token失败", c)
		return
	}

	permissions, err := roleAuthService.GetUserAllPermissions(user.ID)
	if err != nil {
		global.GVA_LOG.Error("获取用户的所有权限失败!", zap.Any("err", err))
		reply.FailWithMessage("获取用户的所有权限失败！", c)
		return
	}
	reply.OkWithDetailed(params.LoginResponse{
		User:        user,
		Token:       token,
		Permissions: permissions,
	}, "登录成功", c)
	return
}

// @Tags User
// @Summary 用户注册
// @Produce  application/json
// @Param body body params.Register true  "手机号码，密码，短信验证码"
// @Success 200
// @Router /user/register [post]
func (u *UserCtl) Register(c *gin.Context) {
	var r params.Register
	_ = c.Bind(&r)
	if err := utils.Verify(r, utils.RegisterVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}

	// 校验验证码
	if err := sms.VerifyCaptcha(r.Mobile, r.Sms); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}

	_, err := userService.GetUserInfoByMobile(r.Mobile)
	if err == nil {
		reply.FailWithMessage("该手机号已经注册过，请直接登录", c)
		return
	}

	if err := auth.CheckGoodPassword(r.Password); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	user := table.User{
		Mobile:     r.Mobile,
		UserName:   r.Mobile,
		Password:   auth.GetPasswd(r.Password),
		RegisterIP: c.ClientIP(),
		Status:     constant.NormalStatus,
	}

	userReturn, err := userService.Register(user)
	if err != nil {
		global.GVA_LOG.Error("注册失败!", zap.Any("err", err))
		reply.FailWithMessage(err.Error(), c)
		return
	}
	reply.OkWithDetailed(userReturn, "注册成功", c)
}

// @Tags User
// @Summary 管理员批量生成用户
// @Produce  application/json
// @Param count query int true "数量"
// @Success 200
// @Router /user/batchCreateUser [post]
func (u *UserCtl) BatchCreateUser(c *gin.Context) {

	value := c.PostForm("count")

	count, _ := strconv.Atoi(value)
	if count <= 0 {
		reply.FailWithMessage("count不能小于等于0", c)
		return
	}

	err := userService.BatchCreateUser(count)
	if err != nil {
		global.GVA_LOG.Error("批量生成用户失败!", zap.Any("err", err))
		reply.FailWithMessage("批量生成用户失败", c)
		return
	}
	reply.Ok(c)
}

// @Tags User
// @Summary 管理员重置密码
// @Produce  application/json
// @Param id query string true "userID"
// @Success 200
// @Router /user/resetPasswordByID [post]
func (u *UserCtl) ResetPasswordByID(c *gin.Context) {
	id := c.PostForm("id")
	if id == "" {
		reply.FailWithMessage("userid不能为空", c)
		return
	}
	if err := userService.ChangePassword(id, utils.RandStr(8)); err != nil {
		global.GVA_LOG.Error("修改密码失败!", zap.Any("err", err))
		reply.FailWithMessage("修改密码失败", c)
		return
	}
	reply.OkWithMessage("修改成功", c)
}

// @Tags User
// @Summary 修改密码
// @Produce  application/json
// @Param body body params.Register true  "手机号码，密码，短信验证码"
// @Success 200
// @Router /user/resetPassword [post]
func (u *UserCtl) ResetPassword(c *gin.Context) {

	var r params.Register
	_ = c.Bind(&r)
	if err := utils.Verify(r, utils.RegisterVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}

	// 校验验证码
	//if err := sms.VerifyCaptcha(r.Mobile, r.Sms); err != nil {
	//	reply.FailWithMessage(err.Error(), c)
	//	return
	//}
	if err := auth.CheckGoodPassword(r.Password); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}

	user, err := userService.GetUserInfoByMobile(r.Mobile)
	if err != nil {
		reply.FailWithMessage("用户未找到", c)
		return
	}

	if err := userService.ChangePassword(user.ID, auth.GetPasswd(r.Password)); err != nil {
		global.GVA_LOG.Error("修改密码失败!", zap.Any("err", err))
		reply.FailWithMessage("修改密码失败", c)
		return
	}
	reply.OkWithMessage("修改成功", c)
}

// @Tags User
// @Summary 分页获取用户列表
// @Produce  application/json
// @Param body query params.GetUserList true "分页获取用户列表"
// @Success 200
// @Router /user/getUserList [get]
func (u *UserCtl) GetUserList(c *gin.Context) {
	var pageInfo params.GetUserList
	_ = c.Bind(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := userService.GetUserInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取用户列表失败!", zap.Any("err", err))
		reply.FailWithMessage("获取用户列表失败", c)
	} else {
		reply.OkWithDetailed(params.PageResult{
			List:   list,
			Total:  total,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
		}, "获取成功", c)
	}
}

// @Tags User
// @Summary 获取用户信息详情
// @Produce  application/json
// @Param id query string true "用户id"
// @Success 200
// @Router /user/getUserInfo [get]
func (u *UserCtl) GetUserInfo(c *gin.Context) {
	userID := c.Query("id")
	if userID == "" {
		reply.FailWithMessage("user_id不能为空", c)
		return
	}
	if user, err := userService.GetUserInfoByID(userID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		reply.FailWithMessage("获取失败", c)
	} else {
		reply.OkWithDetailed(gin.H{"userInfo": user}, "获取成功", c)
	}
}

// @Tags User
// @Summary 修改用户信息
// @Produce  application/json
// @Param body body params.SetUserInfo false  "修改用户信息"
// @Success 200
// @Router /user/setUserInfo [put]
func (u *UserCtl) SetUserInfo(c *gin.Context) {
	var r params.SetUserInfo
	_ = c.Bind(&r)
	if err := utils.Verify(r, utils.IdVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}

	//user := table.User{
	//	UserName:  r.UserName,
	//	HeaderImg: r.HeaderImg,
	//	Note:      r.Note,
	//	Status:    r.Status,
	//}
	//user.ID = r.ID

	if err := userService.SetUserInfo(r); err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Any("err", err))
		reply.FailWithMessage("设置失败", c)
		return
	}
	user, err := userService.GetUserInfoByID(r.ID)
	if err != nil {
		global.GVA_LOG.Error("获取用户信息失败!", zap.Any("err", err))
		reply.FailWithMessage("获取用户信息失败", c)
		return
	}
	reply.OkWithData(user, c)
}

// @Tags User
// @Summary 删除用户
// @Produce  application/json
// @Param id query string true "用户id"
// @Success 200
// @Router /user/deleteUser [delete]
func (u *UserCtl) DeleteUser(c *gin.Context) {
	userID := c.Query("id")
	if userID == "" {
		reply.FailWithMessage("user_id不能为空", c)
		return
	}
	if err := userService.DeleteUser(userID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		reply.FailWithMessage("删除失败", c)
	} else {
		reply.OkWithMessage("删除成功", c)
	}
}
