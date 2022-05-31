package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pluto/global"
	"pluto/model/params"
	"pluto/model/reply"
	"pluto/model/table"
	"pluto/utils"
	"time"
)

type SipCtl struct{}

// @Tags Sip
// @Summary 创建线路
// @Produce  application/json
// @Param body body params.SetSipInfo true  "创建线路"
// @Success 200
// @Router /sip/createSip [post]
func (s *SipCtl) CreateSip(c *gin.Context) {
	var r params.SetSipInfo
	_ = c.Bind(&r)
	if err := utils.Verify(r, utils.CreateSipVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}

	sip := table.Sip{
		Mobile:      r.Mobile,
		Name:        r.Name,
		IntervalSec: r.IntervalSec,
		SipIP:       r.SipIP,
		SipPort:     r.SipPort,
		Note:        r.Note,
	}

	sipInfo, err := SipService.CreateSip(sip)
	if err != nil {
		global.GVA_LOG.Error("创建线路失败!", zap.Any("err", err))
		reply.FailWithMessage("创建线路失败", c)
		return
	}
	reply.OkWithData(sipInfo, c)
}

// @Tags Sip
// @Summary 查询线路列表
// @Produce  application/json
// @Param body body params.GetSipList true  "查询线路列表"
// @Success 200
// @Router /sip/getSipList [get]
func (s *SipCtl) GetSipList(c *gin.Context) {
	var pageInfo params.GetSipList
	_ = c.Bind(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := SipService.GetSipList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取线路列表失败!", zap.Any("err", err))
		reply.FailWithMessage("获取线路列表失败", c)
	} else {
		reply.OkWithDetailed(params.PageResult{
			List:   list,
			Total:  total,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
		}, "获取成功", c)
	}
}

// @Tags Sip
// @Summary 获取线路信息详情
// @Produce  application/json
// @Param id query string true "线路id"
// @Success 200
// @Router /sip/getUserInfo [get]
func (u *UserCtl) GetSipInfo(c *gin.Context) {
	sipID := c.Query("id")
	if sipID == "" {
		reply.FailWithMessage("sipID不能为空", c)
		return
	}
	if sip, err := SipService.GetSipInfoByID(sipID); err != nil {
		global.GVA_LOG.Error("获取线路信息失败!", zap.Any("err", err))
		reply.FailWithMessage("获取线路信息失败", c)
	} else {
		reply.OkWithData(sip, c)
	}
}

// @Tags Sip
// @Summary 设置线路信息
// @Produce  application/json
// @Param body body params.SetSipInfo false  "设置线路信息"
// @Success 200
// @Router /sip/setSipInfo [put]
func (s *SipCtl) SetSipInfo(c *gin.Context) {
	var r params.SetSipInfo
	_ = c.Bind(&r)
	if err := utils.Verify(r, utils.IdVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}

	if err := SipService.SetSipInfo(r); err != nil {
		global.GVA_LOG.Error("设置线路信息失败!", zap.Any("err", err))
		reply.FailWithMessage("设置线路信息失败!", c)
		return
	}
	sip, err := SipService.GetSipInfoByID(r.ID)
	if err != nil {
		global.GVA_LOG.Error("获取线路信息失败!", zap.Any("err", err))
		reply.FailWithMessage("获取线路信息失败", c)
		return
	}
	reply.OkWithData(sip, c)
}

// @Tags Sip
// @Summary 删除线路
// @Produce  application/json
// @Param id query string true "线路id"
// @Success 200
// @Router /sip/deleteSip [delete]
func (s *SipCtl) DeleteSip(c *gin.Context) {
	sipID := c.Query("id")
	if sipID == "" {
		reply.FailWithMessage("sipID不能为空", c)
		return
	}
	if err := SipService.DeleteSip(sipID); err != nil {
		global.GVA_LOG.Error("删除线路失败!", zap.Any("err", err))
		reply.FailWithMessage("删除线路失败", c)
	} else {
		reply.OkWithMessage("删除线路成功", c)
	}
}

// @Tags Sip
// @Summary 查询用户下已分配的线路
// @Produce  application/json
// @Param body body params.GetUserSipList true  "查询线路列表"
// @Success 200
// @Router /sip/getUserSips [get]
func (s *SipCtl) GetUserSips(c *gin.Context) {
	var pageInfo params.GetUserSipList
	_ = c.Bind(&pageInfo)
	if err := utils.Verify(pageInfo, utils.GetUserSipList); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := SipService.GetUserSips(pageInfo); err != nil {
		global.GVA_LOG.Error("获取用户下的线路列表失败!", zap.Any("err", err))
		reply.FailWithMessage("获取用户下的线路列表失败", c)
	} else {
		reply.OkWithDetailed(params.PageResult{
			List:   list,
			Total:  total,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
		}, "获取成功", c)
	}
}

// @Tags Sip
// @Summary 给用户分配线路
// @Produce  application/json
// @Param body body params.SetUserSip false  "给用户分配线路"
// @Success 200
// @Router /sip/userAddSip [post]
func (s *SipCtl) UserAddSip(c *gin.Context) {
	var r params.SetUserSip
	_ = c.Bind(&r)
	if err := utils.Verify(r, utils.UserAddSipVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	date, _ := time.Parse(time.RFC3339, r.ExpireAt)

	info := table.UserSip{
		UserID:     r.UserID,
		SipID:      r.SipID,
		ExpireAt:   date,
		Concurrent: r.Concurrent,
		Price:      r.Price,
	}

	if userSip, err := SipService.UserAddSip(info); err != nil {
		global.GVA_LOG.Error("给用户分配线路失败!", zap.Any("err", err))
		reply.FailWithMessage("给用户分配线路失败!", c)
		return
	} else {
		reply.OkWithData(userSip, c)
	}
}

// @Tags Sip
// @Summary 设置分配给用户线路信息
// @Produce  application/json
// @Param body body params.SetUserSip false  "设置分配给用户线路信息"
// @Success 200
// @Router /sip/setUserSip [put]
func (s *SipCtl) SetUserSip(c *gin.Context) {
	var r params.SetUserSip
	_ = c.Bind(&r)
	if err := utils.Verify(r, utils.IdVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if err := SipService.SetUserSip(r); err != nil {
		global.GVA_LOG.Error("设置用户线路信息失败!", zap.Any("err", err))
		reply.FailWithMessage("设置用户线路信息失败!", c)
		return
	}
	reply.Ok(c)
}

// @Tags Sip
// @Summary 删除用户下的线路
// @Produce  application/json
// @Param id query string true "id"
// @Success 200
// @Router /sip/userDelSip [delete]
func (s *SipCtl) UserDelSip(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		reply.FailWithMessage("id不能为空", c)
		return
	}
	if err := SipService.UserDelSip(id); err != nil {
		global.GVA_LOG.Error("删除线路失败!", zap.Any("err", err))
		reply.FailWithMessage("删除线路失败", c)
	} else {
		reply.OkWithMessage("删除线路成功", c)
	}
}
