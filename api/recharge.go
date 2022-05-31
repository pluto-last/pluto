package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pluto/global"
	"pluto/model/params"
	"pluto/model/reply"
	"pluto/utils"
)

type RechargeCtl struct{}

// @Tags Recharge
// @Summary 查询充值记录列表
// @Produce  application/json
// @Param body body params.GetRechargeList true  "查询充值记录列表"
// @Success 200
// @Router /recharge/getRechargeList [get]
func (r *RechargeCtl) GetRechargeList(c *gin.Context) {
	var pageInfo params.GetRechargeList
	_ = c.Bind(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := RechargeService.GetRechargeList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取充值记录列表失败!", zap.Any("err", err))
		reply.FailWithMessage("获取充值记录列表失败", c)
	} else {
		reply.OkWithDetailed(params.PageResult{
			List:   list,
			Total:  total,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
		}, "获取成功", c)
	}
}

// @Tags Recharge
// @Summary 管理员充值
// @Produce  application/json
// @Param body body params.CreateRecharge true  "管理员充值"
// @Success 200
// @Router /recharge/createRecharge [post]
func (r *RechargeCtl) CreateRecharge(c *gin.Context) {
	var req params.CreateRecharge
	_ = c.Bind(&req)
	if err := utils.Verify(req, utils.Recharge); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if err := RechargeService.Recharge(req); err != nil {
		global.GVA_LOG.Error("管理员充值失败!", zap.Any("err", err))
		reply.FailWithMessage("管理员充值失败", c)
	}
	reply.Ok(c)
}
