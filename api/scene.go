package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pluto/global"
	"pluto/model/params"
	"pluto/model/reply"
	"pluto/model/table"
	"pluto/utils"
)

type SceneCtl struct{}

// @Tags Scene
// @Summary 新增话术
// @Produce  application/json
// @Param body body params.CreateScene true  "新增话术"
// @Success 200
// @Router /scene/createScene [post]
func (s *SceneCtl) CreateScene(c *gin.Context) {
	var r params.CreateScene
	_ = c.BindJSON(&r)

	err := SceneService.CreateScene(r.Scene, r.SceneNode)
	if err != nil {
		global.GVA_LOG.Error("新增话术失败!", zap.Any("err", err))
		reply.FailWithMessage("新增话术失败", c)
		return
	}
	reply.Ok(c)
}

// @Tags Scene
// @Summary 修改话术节点
// @Produce  application/json
// @Param body body table.SceneNode true  "新增话术"
// @Success 200
// @Router /scene/setSceneNode [put]
func (s *SceneCtl) SetSceneNode(c *gin.Context) {
	var r table.SceneNode
	_ = c.BindJSON(&r)

	if err := utils.Verify(r, utils.IdVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}

	err := SceneService.SetSceneNode(r)
	if err != nil {
		global.GVA_LOG.Error("修改话术节点失败!", zap.Any("err", err))
		reply.FailWithMessage("修改话术节点失败", c)
		return
	}
	reply.Ok(c)
}

// @Tags Scene
// @Summary 修改话术
// @Produce  application/json
// @Param body body table.Scene true  "修改话术"
// @Success 200
// @Router /scene/setScene [put]
func (s *SceneCtl) SetScene(c *gin.Context) {
	var r table.Scene
	_ = c.BindJSON(&r)

	if err := utils.Verify(r, utils.IdVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}

	err := SceneService.SetScene(r)
	if err != nil {
		global.GVA_LOG.Error("修改话术失败!", zap.Any("err", err))
		reply.FailWithMessage("修改话术失败", c)
		return
	}
	reply.Ok(c)
}

// @Tags Scene
// @Summary 查询话术列表
// @Produce  application/json
// @Param body body params.GetSceneList true  "查询话术列表"
// @Success 200
// @Router /scene/getSceneList [get]
func (s *SceneCtl) GetSceneList(c *gin.Context) {
	var pageInfo params.GetSceneList
	_ = c.Bind(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := SceneService.GetSceneList(pageInfo); err != nil {
		global.GVA_LOG.Error("查询话术列表失败!", zap.Any("err", err))
		reply.FailWithMessage("查询话术列表失败", c)
	} else {
		reply.OkWithDetailed(params.PageResult{
			List:   list,
			Total:  total,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
		}, "获取成功", c)
	}
}

// @Tags Scene
// @Summary 查询话术节点详情详情
// @Produce  application/json
// @Param id query string true "id"
// @Success 200
// @Router /scene/getSceneInfoByID [get]
func (s *SceneCtl) GetSceneInfoByID(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		reply.FailWithMessage("ID不能为空", c)
		return
	}
	if scene, err := SceneService.GetSceneInfoByID(id); err != nil {
		global.GVA_LOG.Error("查询话术节点详情失败!", zap.Any("err", err))
		reply.FailWithMessage("查询话术节点详情失败", c)
	} else {
		reply.OkWithData(scene, c)
	}
}

// @Tags Scene
// @Summary 删除话术
// @Produce  application/json
// @Param id query string true "id"
// @Success 200
// @Router /scene/seleteScene [delete]
func (s *SceneCtl) DeleteScene(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		reply.FailWithMessage("ID不能为空", c)
		return
	}
	if err := SceneService.DeleteScene(id); err != nil {
		global.GVA_LOG.Error("删除话术失败!", zap.Any("err", err))
		reply.FailWithMessage("删除话术失败", c)
	} else {
		reply.Ok(c)
	}
}
