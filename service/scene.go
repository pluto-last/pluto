package service

import (
	"errors"
	"pluto/global"
	"pluto/middleware/db"
	"pluto/model/params"
	"pluto/model/table"
)

type SceneService struct{}

// CreateScene 新增话术
func (s *SceneService) CreateScene(scene table.Scene, sceneNode []table.SceneNode) (err error) {
	tx := db.Begin(global.GVA_DB)
	defer tx.RollbackIfFailed()

	// 保存话术表
	err = tx.Create(&scene).Error
	if err != nil {
		return
	}

	// 保存话术节点表
	for i, item := range sceneNode {
		item.SceneID = scene.ID
		err = tx.Create(&sceneNode[i]).Error
		if err != nil {
			return
		}
	}
	tx.Commit()
	return nil
}

// SetSceneNode 修改话术节点
func (s *SceneService) SetSceneNode(info table.SceneNode) (err error) {
	if info.ID == "" {
		err = errors.New("id is null")
		return
	}
	err = global.GVA_DB.Table("m_scene_node").Where("id = ?", info.ID).Updates(&info).Error
	return err
}

// SetScene 修改话术
func (s *SceneService) SetScene(info table.Scene) (err error) {
	if info.ID == "" {
		err = errors.New("id is null")
		return
	}
	err = global.GVA_DB.Table("m_scene").Where("id = ?", info.ID).Updates(&info).Error
	return err
}

// GetSceneList 查询话术列表
func (s *SceneService) GetSceneList(info params.GetSceneList) (list []table.Scene, total int64, err error) {
	db := global.GVA_DB.Model(&table.Scene{})
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

// GetSceneInfoByID 查询话术节点详情详情
func (s *SceneService) GetSceneInfoByID(id string) (scene table.Scene, err error) {
	err = global.GVA_DB.Where("id = ?", id).Preload("SceneNode").First(&scene).Error
	return
}

// DeleteScene 删除话术
func (s *SceneService) DeleteScene(id string) (err error) {

	tx := db.Begin(global.GVA_DB)
	defer tx.RollbackIfFailed()
	err = global.GVA_DB.Where("id = ?", id).Delete(&table.Scene{}).Error
	if err != nil {
		return
	}

	err = global.GVA_DB.Where("scene_id = ?", id).Delete(&table.SceneNode{}).Error
	if err != nil {
		return
	}

	tx.Commit()
	return
}
