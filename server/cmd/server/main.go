package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "xiangshoufu/swagger" // swagger docs
	"xiangshoufu/internal/async"
	"xiangshoufu/internal/cache"
	"xiangshoufu/internal/channel"
	"xiangshoufu/internal/channel/hengxintong"
	"xiangshoufu/internal/handler"
	"xiangshoufu/internal/jobs"
	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"
)

// Config 应用配置
type Config struct {
	DatabaseURL     string
	ServerPort      string
	HxtPublicKey    string
	AlertWebhookURL string
	SwaggerEnabled  bool // 是否启用Swagger UI
}

// @title           8通道回调服务 API
// @version         1.0
// @description     代理商分润管理系统 - 支付通道回调处理服务
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@xiangshoufu.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	log.Println("Starting 8-Channel Callback Server...")

	// 加载配置
	config := loadConfig()

	// 1. 初始化数据库连接
	db, err := initDatabase(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully")

	// 2. 初始化基础组件
	localCache := cache.NewLocalCache(nil)
	memQueue := async.NewMemoryQueue(nil)

	// 3. 初始化适配器工厂并注册适配器
	factory := channel.GetFactory()
	registerAdapters(factory, config)

	// 4. 初始化Repository
	callbackRepo := repository.NewGormRawCallbackRepository(db)
	messageRepo := repository.NewGormMessageRepository(db)
	deviceFeeRepo := repository.NewGormDeviceFeeRepository(db)
	rateChangeRepo := repository.NewGormRateChangeRepository(db)
	transactionRepo := repository.NewGormTransactionRepository(db)
	profitRepo := repository.NewGormProfitRecordRepository(db)
	walletRepo := repository.NewGormWalletRepository(db)
	walletLogRepo := repository.NewGormWalletLogRepository(db)
	agentRepo := repository.NewGormAgentRepository(db)
	agentPolicyRepo := repository.NewGormAgentPolicyRepository(db)

	// 5. 初始化消息服务
	pushConfig := &service.PushConfig{
		Enabled:    false, // 开发环境关闭推送
		WebhookURL: config.AlertWebhookURL,
	}
	messageService := service.NewMessageService(messageRepo, pushConfig)

	// 6. 初始化分润服务
	profitService := service.NewProfitService(
		transactionRepo,
		profitRepo,
		walletRepo,
		walletLogRepo,
		agentRepo,
		agentPolicyRepo,
		messageService,
		memQueue,
	)

	// 7. 初始化回调处理服务
	callbackProcessor := service.NewCallbackProcessor(
		factory,
		callbackRepo,
		transactionRepo,
		deviceFeeRepo,
		rateChangeRepo,
		profitService,
		memQueue,
	)

	// 8. 初始化代扣相关Repository和Service
	deductionPlanRepo := repository.NewGormDeductionPlanRepository(db)
	deductionRecordRepo := repository.NewGormDeductionRecordRepository(db)
	deductionChainRepo := repository.NewGormDeductionChainRepository(db)
	deductionChainItemRepo := repository.NewGormDeductionChainItemRepository(db)

	deductionService := service.NewDeductionService(
		deductionPlanRepo,
		deductionRecordRepo,
		deductionChainRepo,
		deductionChainItemRepo,
		walletRepo,
		walletLogRepo,
		agentRepo,
	)

	// 9. 初始化终端下发相关Repository和Service
	terminalRepo := repository.NewGormTerminalRepository(db)
	terminalDistributeRepo := repository.NewGormTerminalDistributeRepository(db)

	terminalDistributeService := service.NewTerminalDistributeService(
		terminalRepo,
		terminalDistributeRepo,
		agentRepo,
		deductionService,
	)

	// 10. 初始化流量费返现相关Repository和Service
	simCashbackPolicyRepo := repository.NewGormSimCashbackPolicyRepository(db)
	simCashbackRecordRepo := repository.NewGormSimCashbackRecordRepository(db)

	simCashbackService := service.NewSimCashbackService(
		terminalRepo,
		simCashbackPolicyRepo,
		simCashbackRecordRepo,
		deviceFeeRepo,
		walletRepo,
		walletLogRepo,
		agentRepo,
		agentPolicyRepo,
		messageService,
		memQueue,
	)

	// 11. 初始化Handler
	callbackHandler := handler.NewCallbackHandler(factory, localCache, memQueue, callbackRepo)
	deductionHandler := handler.NewDeductionHandler(deductionService)
	terminalDistributeHandler := handler.NewTerminalDistributeHandler(terminalDistributeService)
	simCashbackHandler := handler.NewSimCashbackHandler(simCashbackService)

	// 12. 初始化PC端新增Repository
	userRepo := repository.NewGormUserRepository(db)
	refreshTokenRepo := repository.NewGormRefreshTokenRepository(db)
	loginLogRepo := repository.NewGormLoginLogRepository(db)

	// 13. 初始化认证服务
	authConfig := service.DefaultAuthConfig()
	authService := service.NewAuthService(authConfig, userRepo, refreshTokenRepo, loginLogRepo, agentRepo)

	// 14. 初始化PC端新增Service
	agentService := service.NewAgentService(agentRepo, agentPolicyRepo, walletRepo, transactionRepo)
	walletService := service.NewWalletService(walletRepo, walletLogRepo, agentRepo)

	// 15. 初始化PC端新增Handler
	authHandler := handler.NewAuthHandler(authService)
	agentHandler := handler.NewAgentHandler(agentService)
	walletHandler := handler.NewWalletHandler(walletService)
	transactionHandler := handler.NewTransactionHandler(transactionRepo)
	profitHandler := handler.NewProfitHandler(profitRepo)
	dashboardHandler := handler.NewDashboardHandler(transactionRepo, profitRepo, agentRepo, walletRepo)
	messageHandler := handler.NewMessageHandler(messageRepo)

	// 16. 初始化商户、终端、政策Handler
	merchantRepo := repository.NewGormMerchantRepository(db)
	merchantHandler := handler.NewMerchantHandler(merchantRepo, transactionRepo)
	terminalHandler := handler.NewTerminalHandler(terminalRepo, transactionRepo)
	policyTemplateRepo := repository.NewGormPolicyTemplateRepository(db)
	policyHandler := handler.NewPolicyHandler(policyTemplateRepo, agentPolicyRepo)

	// 17. 初始化监控服务
	metricsService := service.NewMetricsService(messageService, memQueue)

	// 13. 订阅队列消息
	setupQueueSubscribers(memQueue, callbackProcessor, profitService, messageService)

	// 14. 初始化定时任务
	scheduler := setupScheduler(metricsService, messageService, transactionRepo, profitService, callbackRepo, callbackProcessor, deductionService, simCashbackService)
	scheduler.Start()

	// 15. 创建HTTP服务器
	router := setupRouter(
		callbackHandler, deductionHandler, terminalDistributeHandler, simCashbackHandler,
		authHandler, agentHandler, walletHandler, transactionHandler, profitHandler, dashboardHandler, messageHandler,
		merchantHandler, terminalHandler, policyHandler,
		authService, metricsService, config.SwaggerEnabled,
	)

	srv := &http.Server{
		Addr:         ":" + config.ServerPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 13. 启动服务器
	go func() {
		log.Printf("Server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 14. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 停止定时任务
	scheduler.Stop()

	// 关闭队列
	memQueue.Close()

	// 关闭缓存
	localCache.Close()

	// 关闭HTTP服务器
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

// loadConfig 加载配置
func loadConfig() *Config {
	config := &Config{
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		ServerPort:      os.Getenv("SERVER_PORT"),
		HxtPublicKey:    os.Getenv("HENGXINTONG_PUBLIC_KEY"),
		AlertWebhookURL: os.Getenv("ALERT_WEBHOOK_URL"),
		SwaggerEnabled:  os.Getenv("SWAGGER_ENABLED") != "false", // 默认启用，生产环境设为false关闭
	}

	// 默认值
	if config.DatabaseURL == "" {
		// macOS Homebrew PostgreSQL 默认使用系统用户名，无需密码
		config.DatabaseURL = "postgres://apple@localhost:5432/xiangshoufu?sslmode=disable"
	}
	if config.ServerPort == "" {
		config.ServerPort = "8080"
	}

	return config
}

// initDatabase 初始化数据库连接
func initDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取底层sql.DB并配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// registerAdapters 注册所有通道适配器
func registerAdapters(factory *channel.AdapterFactory, config *Config) {
	// 恒信通适配器
	hxtConfig := &channel.ChannelConfig{
		ChannelCode: channel.ChannelCodeHengxintong,
		ChannelName: "恒信通",
		PublicKey:   config.HxtPublicKey,
		Enabled:     true,
	}
	hxtAdapter, err := hengxintong.NewAdapter(hxtConfig)
	if err != nil {
		log.Printf("Failed to create Hengxintong adapter: %v", err)
	} else {
		factory.Register(hxtAdapter)
		log.Printf("Registered adapter: %s", hxtAdapter.GetChannelName())
	}

	// TODO: 注册其他7个通道适配器
	// factory.Register(lakala.NewAdapter(...))
	// factory.Register(yeahka.NewAdapter(...))

	log.Printf("Total registered adapters: %d", len(factory.GetSupportedChannels()))
}

// setupQueueSubscribers 设置队列订阅者
func setupQueueSubscribers(
	queue *async.MemoryQueue,
	callbackProcessor *service.CallbackProcessor,
	profitService *service.ProfitService,
	msgService *service.MessageService,
) {
	// 原始回调处理队列
	queue.Subscribe(async.TopicRawCallback, func(msg []byte) error {
		return callbackProcessor.ProcessMessage(msg)
	})

	// 分润计算队列
	queue.Subscribe(async.TopicProfitCalc, func(msg []byte) error {
		return profitService.ProcessMessage(msg)
	})

	// 消息通知队列
	queue.Subscribe(async.TopicNotification, func(msg []byte) error {
		return msgService.ProcessMessage(msg)
	})
}

// setupScheduler 设置定时任务调度器
func setupScheduler(
	metricsService *service.MetricsService,
	messageService *service.MessageService,
	transactionRepo *repository.GormTransactionRepository,
	profitService *service.ProfitService,
	callbackRepo *repository.GormRawCallbackRepository,
	callbackProcessor *service.CallbackProcessor,
	deductionService *service.DeductionService,
	simCashbackService *service.SimCashbackService,
) *jobs.Scheduler {
	scheduler := jobs.NewScheduler()

	// 告警检查（每分钟）
	alertJob := service.NewAlertCheckerJob(metricsService)
	scheduler.AddJob("alert_checker", 1*time.Minute, alertJob.Run)

	// 分润计算兜底（每5分钟）
	profitCalcJob := jobs.NewProfitCalculatorJob(transactionRepo, profitService)
	scheduler.AddJob("profit_calculator", 5*time.Minute, profitCalcJob.Run)

	// 回调重试（每5分钟）
	retryJob := jobs.NewCallbackRetryJob(callbackRepo, callbackProcessor)
	scheduler.AddJob("callback_retry", 5*time.Minute, retryJob.Run)

	// 消息清理（每6小时）
	cleanupJob := jobs.NewMessageCleanupJob(messageService)
	scheduler.AddJob("message_cleanup", 6*time.Hour, cleanupJob.Run)

	// 分区管理（每天检查一次）
	partitionJob := jobs.NewPartitionManagerJob()
	scheduler.AddJob("partition_manager", 24*time.Hour, partitionJob.Run)

	// 每日代扣任务（每24小时执行一次）
	jobs.SetupDeductionJobs(scheduler, deductionService, simCashbackService)

	return scheduler
}

// setupRouter 设置路由
func setupRouter(
	callbackHandler *handler.CallbackHandler,
	deductionHandler *handler.DeductionHandler,
	terminalDistributeHandler *handler.TerminalDistributeHandler,
	simCashbackHandler *handler.SimCashbackHandler,
	authHandler *handler.AuthHandler,
	agentHandler *handler.AgentHandler,
	walletHandler *handler.WalletHandler,
	transactionHandler *handler.TransactionHandler,
	profitHandler *handler.ProfitHandler,
	dashboardHandler *handler.DashboardHandler,
	messageHandler *handler.MessageHandler,
	merchantHandler *handler.MerchantHandler,
	terminalHandler *handler.TerminalHandler,
	policyHandler *handler.PolicyHandler,
	authService *service.AuthService,
	metricsService *service.MetricsService,
	swaggerEnabled bool,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 全局中间件
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggingMiddleware())

	// 限流中间件（每秒1000个请求，桶容量2000）
	globalLimiter := middleware.NewRateLimiter(1000, 2000)
	router.Use(middleware.RateLimitMiddleware(globalLimiter))

	// 健康检查
	// @Summary 健康检查
	// @Description 检查服务是否正常运行
	// @Tags 系统
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /health [get]
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 监控指标
	// @Summary 获取监控指标
	// @Description 获取系统运行指标和各通道统计
	// @Tags 系统
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /metrics [get]
	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, metricsService.GetAllMetrics())
	})

	// 通道回调入口
	// POST /callback/:channel_code
	router.POST("/callback/:channel_code", callbackHandler.HandleCallback)

	// API v1 路由组
	apiV1 := router.Group("/api/v1")
	{
		// 注册认证路由（公开）
		handler.RegisterAuthRoutes(apiV1, authHandler, authService)

		// 注册代扣管理路由
		handler.RegisterDeductionRoutes(apiV1, deductionHandler)

		// 注册终端下发路由
		handler.RegisterTerminalDistributeRoutes(apiV1, terminalDistributeHandler)

		// 注册流量费返现路由
		handler.RegisterSimCashbackRoutes(apiV1, simCashbackHandler)

		// 注册PC端新增路由
		handler.RegisterAgentRoutes(apiV1, agentHandler, authService)
		handler.RegisterWalletRoutes(apiV1, walletHandler, authService)
		handler.RegisterTransactionRoutes(apiV1, transactionHandler, authService)
		handler.RegisterProfitRoutes(apiV1, profitHandler, authService)
		handler.RegisterDashboardRoutes(apiV1, dashboardHandler, authService)
		handler.RegisterMessageRoutes(apiV1, messageHandler, authService)

		// 注册商户、终端、政策路由
		handler.RegisterMerchantRoutes(apiV1, merchantHandler, authService)
		handler.RegisterTerminalRoutes(apiV1, terminalHandler, authService)
		handler.RegisterPolicyRoutes(apiV1, policyHandler, authService)
	}

	// Swagger UI (可通过环境变量关闭)
	if swaggerEnabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		log.Println("Swagger UI enabled at /swagger/index.html")
	}

	return router
}
