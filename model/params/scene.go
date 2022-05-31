package params

import (
	"pluto/model/table"
)

type GetSceneList struct {
	PageInfo
	Name string `json:"name" form:"name"` // 名称
}

type CreateScene struct {
	Scene     table.Scene       `json:"scene" form:"scene"`
	SceneNode []table.SceneNode `json:"sceneNode" form:"sceneNode"`
}
