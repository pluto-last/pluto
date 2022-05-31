package global

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"pluto/config"
	"pluto/middleware/cache"
)

var (
	GVA_DB     *gorm.DB
	GVA_REDIS  *cache.RCache
	GVA_CONFIG *config.Specification
	GVA_LOG    *zap.Logger
)
