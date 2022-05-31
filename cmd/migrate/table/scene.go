package table

import "pluto/global"

// Scene 话术表
type Scene struct {
	global.UUID
	Name      string      `json:"name" form:"name"`         // 话术名称
	Describe  string      `json:"describe" form:"describe"` // 描述
	SceneNode []SceneNode `gorm:"ForeignKey:ID;AssociationForeignKey:SceneID" json:"sceneNode" `
}

func (Scene) TableName() string {
	return "m_scene"
}

// SceneNode 话术流程节点表
type SceneNode struct {
	global.UUID
	SceneID string `gorm:"index" json:"sceneID" form:"sceneID"` // 话术ID
	Name    string `json:"name" form:"name"`                    // 节点名称
	Text    string `json:"text" form:"text"`                    // 文字
}

func (SceneNode) TableName() string {
	return "m_scene_node"
}

// NodeBranch 话术节点分支表
type NodeBranch struct {
	global.UUID
	Type       string `json:"type"`                         // 分支类型 肯定分支，拒绝分支，默认分支
	SceneID    string `gorm:"index" json:"sceneID"`         // 话术ID
	NodeID     string `json:"nodeID"`                       // 节点
	NextNodeID string `gorm:"nextNodeID" json:"nextNodeID"` // 分支所连接子节点ID
}

func (NodeBranch) TableName() string {
	return "m_node_branch"
}
