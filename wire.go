//go:generate wire
//go:build wireinject

package main

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/muxi-Infra/auditor-Backend/client"
	"github.com/muxi-Infra/auditor-Backend/config"
	"github.com/muxi-Infra/auditor-Backend/controller"
	confluentinc_llm "github.com/muxi-Infra/auditor-Backend/events/confluentinc-llm"
	"github.com/muxi-Infra/auditor-Backend/ioc"
	lcl "github.com/muxi-Infra/auditor-Backend/langchain/client"
	lc "github.com/muxi-Infra/auditor-Backend/langchain/config"
	"github.com/muxi-Infra/auditor-Backend/middleware"
	"github.com/muxi-Infra/auditor-Backend/pkg/jwt"
	"github.com/muxi-Infra/auditor-Backend/pkg/viperx"
	"github.com/muxi-Infra/auditor-Backend/repository/cache"
	"github.com/muxi-Infra/auditor-Backend/repository/dao"
	"github.com/muxi-Infra/auditor-Backend/server"
	"github.com/muxi-Infra/auditor-Backend/service"
)

// wire.go

// 提供 dao.UserDAO 的 provider
func ProvideUserDAO(db *gorm.DB) dao.UserDAOInterface {
	return &dao.UserDAO{DB: db}
}
func ProvideRedisCache(c *ioc.RedisCache) cache.CacheInterface {
	return c
}

func InitWebServer(confPath string) *App {
	wire.Build(
		viperx.NewVipperSettingFromNacos,
		config.NewAppConf,
		config.NewJWTConf,
		config.NewOAuthConf,
		config.NewDBConf,
		config.NewLogConf,
		config.NewCacheConf,
		config.NewPrometheusConf,
		config.NewMiddleWareConf,
		config.NewQiniuConf,
		config.NewKafkaConf,
		lc.NewMuxiAIConf,
		// 初始化基础依赖
		ioc.InitDB,
		ioc.InitLogger,
		ioc.InitRedis,
		ioc.NewRedisCache,
		ioc.InitPrometheus,
		ioc.InitProducer,
		// 初始化具体模块
		dao.NewUserDAO,
		dao.NewProjectDAO,
		dao.NewItemDao,
		dao.NewCommentDao,
		cache.NewProjectCache,
		jwt.NewRedisJWTHandler,
		confluentinc_llm.NewLlmProducer,
		service.NewAuthService,
		service.NewUserService,
		ProvideUserDAO,
		ProvideRedisCache,
		lcl.AuditAIConnect,
		service.NewProjectService,
		service.NewItemService,
		service.NewTubeService,
		service.NewRemoveService,
		service.NewLLMService,
		controller.NewOAuthController,
		controller.NewUserController,
		controller.NewProjectController,
		controller.NewItemController,
		controller.NewTuberController,
		controller.NewRemoveController,
		controller.NewLLMController,
		client.NewOAuthClient,
		server.NewServer,

		// 中间件
		middleware.NewAuthMiddleware,
		middleware.NewLoggerMiddleware,
		middleware.NewCorsMiddleware,
		// 应用入口
		NewApp,
	)
	return &App{}
}
