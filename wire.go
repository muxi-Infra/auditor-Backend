//go:generate wire
//go:build wireinject

package main

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/client"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/config"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/controller"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/ioc"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/middleware"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/jwt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/viperx"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/cache"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/dao"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/router"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/service"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// wire.go

// 提供 dao.UserDAO 的 provider
func ProvideUserDAO(db *gorm.DB) dao.UserDAOInterface {
	return &dao.UserDAO{DB: db}
}
func ProvideRedisCache(c *ioc.RedisCache) cache.ProjectCacheInterface {
	return c
}

func InitWebServer(confPath string) *App {
	wire.Build(
		viperx.NewVipperSetting,
		config.NewAppConf,
		config.NewJWTConf,
		config.NewOAuthConf,
		config.NewDBConf,
		config.NewLogConf,
		config.NewCacheConf,
		config.NewPrometheusConf,
		config.NewMiddleWareConf,
		config.NewQiniuConf,
		// 初始化基础依赖
		ioc.InitDB,
		ioc.InitLogger,
		ioc.InitRedis,
		ioc.NewRedisCache,
		ioc.InitPrometheus,
		// 初始化具体模块
		dao.NewUserDAO,
		cache.NewProjectCache,
		jwt.NewRedisJWTHandler,
		service.NewAuthService,
		service.NewUserService,
		ProvideUserDAO,
		ProvideRedisCache,
		service.NewProjectService,
		service.NewItemService,
		service.NewTubeService,
		service.NewRemoveService,
		controller.NewOAuthController,
		controller.NewUserController,
		controller.NewProjectController,
		controller.NewItemController,
		controller.NewTuberController,
		controller.NewRemoveController,
		client.NewOAuthClient,
		router.NewRouter,

		// 中间件
		middleware.NewAuthMiddleware,
		middleware.NewLoggerMiddleware,
		middleware.NewCorsMiddleware,
		// 应用入口
		NewApp,
	)
	return &App{}
}
