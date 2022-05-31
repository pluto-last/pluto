package service

import (
	"pluto/global"
	"pluto/model/table"
)

type OperationRecordService struct{}

func (operationRecordService *OperationRecordService) CreateSysOperationRecord(sysOperationRecord table.OperationRecord) (err error) {
	err = global.GVA_DB.Create(&sysOperationRecord).Error
	return err
}
