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
	merchantRepo := repository.NewGormMerchantRepository(db)
	rateStagePolicyRepo := repository.NewGormRateStagePolicyRepository(db)
	terminalRepo := repository.NewGormTerminalRepository(db) // 提前初始化，供callbackProcessor使用

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

	// 6.1 初始化费率阶梯服务并注入到分润服务
	rateStagingService := service.NewRateStagingService(
		rateStagePolicyRepo,
		merchantRepo,
		agentRepo,
	)
	profitService.SetRateStagingService(rateStagingService)

	// 7. 初始化回调处理服务
	callbackProcessor := service.NewCallbackProcessor(
		factory,
		callbackRepo,
		transactionRepo,
		deviceFeeRepo,
		rateChangeRepo,
		merchantRepo,
		terminalRepo,
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

	// 8.1 初始化货款代扣相关Repository和Service
	goodsDeductionRepo := repository.NewGormGoodsDeductionRepository(db)
	goodsDeductionDetailRepo := repository.NewGormGoodsDeductionDetailRepository(db)
	goodsDeductionTerminalRepo := repository.NewGormGoodsDeductionTerminalRepository(db)
	goodsDeductionNotificationRepo := repository.NewGormGoodsDeductionNotificationRepository(db)

	goodsDeductionService := service.NewGoodsDeductionService(
		goodsDeductionRepo,
		goodsDeductionDetailRepo,
		goodsDeductionTerminalRepo,
		goodsDeductionNotificationRepo,
		walletRepo,
		walletLogRepo,
		agentRepo,
	)

	// 8.2 将货款代扣服务注入到分润服务（延迟注入，避免循环依赖）
	profitService.SetGoodsDeductionService(goodsDeductionService)

	// 9. 初始化终端下发相关Repository和Service
	terminalDistributeRepo := repository.NewGormTerminalDistributeRepository(db)
	terminalRecallRepo := repository.NewGormTerminalRecallRepository(db)
	terminalImportRecordRepo := repository.NewGormTerminalImportRecordRepository(db)

	terminalDistributeService := service.NewTerminalDistributeService(
		terminalRepo,
		terminalDistributeRepo,
		agentRepo,
		deductionService,
	)

	// 9.1 将货款代扣服务注入到终端划拨服务
	terminalDistributeService.SetGoodsDeductionService(goodsDeductionService)

	terminalService := service.NewTerminalService(
		terminalRepo,
		terminalRecallRepo,
		terminalImportRecordRepo,
		agentRepo,
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
	goodsDeductionHandler := handler.NewGoodsDeductionHandler(goodsDeductionService)

	// 12. 初始化PC端新增Repository
	userRepo := repository.NewGormUserRepository(db)
	refreshTokenRepo := repository.NewGormRefreshTokenRepository(db)
	loginLogRepo := repository.NewGormLoginLogRepository(db)

	// 13. 初始化认证服务
	authConfig := service.DefaultAuthConfig()
	authService := service.NewAuthService(authConfig, userRepo, refreshTokenRepo, loginLogRepo, agentRepo)

	// 14. 初始化PC端新增Service
	agentService := service.NewAgentService(agentRepo, agentPolicyRepo, walletRepo, transactionRepo, profitRepo)
	walletService := service.NewWalletService(walletRepo, walletLogRepo, agentRepo)

	// 15. 初始化PC端新增Handler
	authHandler := handler.NewAuthHandler(authService)
	authHandler.SetAgentService(agentService) // 注入代理商服务用于公开注册接口
	agentHandler := handler.NewAgentHandler(agentService)
	walletHandler := handler.NewWalletHandler(walletService)
	transactionHandler := handler.NewTransactionHandler(transactionRepo)
	profitHandler := handler.NewProfitHandler(profitRepo)
	dashboardHandler := handler.NewDashboardHandler(transactionRepo, profitRepo, agentRepo, walletRepo)
	messageHandler := handler.NewMessageHandler(messageRepo)

	// 15.1 初始化管理端消息Handler
	adminMessageHandler := handler.NewAdminMessageHandler(messageService, agentRepo)

	// 16. 初始化商户、终端Handler
	merchantService := service.NewMerchantService(merchantRepo, agentRepo, transactionRepo, terminalRepo)
	merchantHandler := handler.NewMerchantHandler(merchantRepo, transactionRepo, merchantService)
	terminalHandler := handler.NewTerminalHandler(terminalRepo, transactionRepo, terminalService)

	// 17. 初始化政策相关Repository
	policyTemplateRepo := repository.NewGormPolicyTemplateRepository(db)
	depositPolicyRepo := repository.NewGormDepositCashbackPolicyRepository(db)
	depositRecordRepo := repository.NewGormDepositCashbackRecordRepository(db)
	rewardPolicyRepo := repository.NewGormActivationRewardPolicyRepository(db)
	rewardRecordRepo := repository.NewGormActivationRewardRecordRepository(db)
	agentDepositPolicyRepo := repository.NewGormAgentDepositCashbackPolicyRepository(db)
	agentSimPolicyRepo := repository.NewGormAgentSimCashbackPolicyRepository(db)
	agentRewardPolicyRepo := repository.NewGormAgentActivationRewardPolicyRepository(db)
	simPolicyRepo := repository.NewGormSimCashbackPolicyRepository(db)
	channelRepo := repository.NewGormChannelRepository(db)

	// 17.1 延迟注入仓储到商户Handler（避免循环依赖）
	merchantHandler.SetAgentRepo(agentRepo)
	merchantHandler.SetChannelRepo(channelRepo)
	merchantHandler.SetTerminalRepo(terminalRepo)

	// 18. 初始化政策管理服务
	policyService := service.NewPolicyService(
		policyTemplateRepo,
		depositPolicyRepo,
		simPolicyRepo,
		rewardPolicyRepo,
		rateStagePolicyRepo,
		agentPolicyRepo,
		agentDepositPolicyRepo,
		agentSimPolicyRepo,
		agentRewardPolicyRepo,
		agentRepo,
	)

	// 19. 初始化押金返现服务
	depositCashbackService := service.NewDepositCashbackService(
		terminalRepo,
		merchantRepo,
		depositPolicyRepo,
		agentDepositPolicyRepo,
		depositRecordRepo,
		walletRepo,
		walletLogRepo,
		agentRepo,
		agentPolicyRepo,
		messageService,
		memQueue,
	)

	// 20. 初始化激活奖励服务
	activationRewardService := service.NewActivationRewardService(
		terminalRepo,
		merchantRepo,
		transactionRepo,
		rewardPolicyRepo,
		agentRewardPolicyRepo,
		rewardRecordRepo,
		walletRepo,
		walletLogRepo,
		agentRepo,
		agentPolicyRepo,
		messageService,
		memQueue,
	)

	// 20.1 初始化代理商通道服务
	agentChannelRepo := repository.NewGormAgentChannelRepository(db)
	agentChannelService := service.NewAgentChannelService(
		agentChannelRepo,
		channelRepo,
		agentRepo,
	)
	agentChannelHandler := handler.NewAgentChannelHandler(agentChannelService)

	// 20.2 初始化充值钱包服务
	chargingWalletRepo := repository.NewGormChargingWalletRepository(db)
	chargingWalletService := service.NewChargingWalletService(
		chargingWalletRepo,
		walletRepo,
		walletLogRepo,
		agentRepo,
	)
	chargingWalletHandler := handler.NewChargingWalletHandler(chargingWalletService)

	// 20.3 初始化沉淀钱包服务
	settlementWalletRepo := repository.NewGormSettlementWalletRepository(db)
	settlementWalletService := service.NewSettlementWalletService(
		settlementWalletRepo,
		chargingWalletRepo,
		walletRepo,
		walletLogRepo,
		agentRepo,
	)
	settlementWalletHandler := handler.NewSettlementWalletHandler(settlementWalletService)

	// 20.4 初始化税筹通道服务
	taxChannelRepo := repository.NewGormTaxChannelRepository(db)
	taxChannelService := service.NewTaxChannelService(taxChannelRepo)
	taxChannelHandler := handler.NewTaxChannelHandler(taxChannelService)

	// 20.5 初始化营销模块（Banner、海报）
	bannerRepo := repository.NewGormBannerRepository(db)
	posterRepo := repository.NewGormPosterRepository(db)
	posterCategoryRepo := repository.NewGormPosterCategoryRepository(db)
	uploadedFileRepo := repository.NewGormUploadedFileRepository(db)

	// 上传服务配置
	uploadDir := "./uploads"     // 上传目录
	uploadBaseURL := "/uploads"  // 访问基础URL
	uploadService := service.NewUploadService(uploadedFileRepo, uploadDir, uploadBaseURL)
	bannerService := service.NewBannerService(bannerRepo)
	posterService := service.NewPosterService(posterRepo, posterCategoryRepo)

	uploadHandler := handler.NewUploadHandler(uploadService)
	bannerHandler := handler.NewBannerHandler(bannerService)
	posterHandler := handler.NewPosterHandler(posterService)

	// 21. 初始化政策Handler
	policyHandler := handler.NewPolicyHandler(policyTemplateRepo, agentPolicyRepo, policyService)

	// 添加depositCashbackService到定时任务
	_ = depositCashbackService // 将在后续定时任务中使用

	// 22. 初始化监控服务
	metricsService := service.NewMetricsService(messageService, memQueue)

	// 23. 订阅队列消息
	setupQueueSubscribers(memQueue, callbackProcessor, profitService, messageService)

	// 24. 初始化定时任务
	scheduler := setupScheduler(
		metricsService, messageService, transactionRepo, profitService,
		callbackRepo, callbackProcessor, deductionService, simCashbackService,
		// 新增参数：激活奖励相关
		terminalRepo, channelRepo, activationRewardService,
		depositRecordRepo, rewardRecordRepo, walletRepo, simCashbackRecordRepo,
		// 新增参数：商户类型计算
		merchantRepo, merchantService,
		// 新增参数：流量费返现
		deviceFeeRepo,
	)
	scheduler.Start()

	// 15. 创建HTTP服务器
	router := setupRouter(
		callbackHandler, deductionHandler, terminalDistributeHandler, simCashbackHandler,
		goodsDeductionHandler,
		authHandler, agentHandler, walletHandler, transactionHandler, profitHandler, dashboardHandler, messageHandler,
		adminMessageHandler,
		merchantHandler, terminalHandler, policyHandler, agentChannelHandler,
		chargingWalletHandler, settlementWalletHandler, taxChannelHandler,
		uploadHandler, bannerHandler, posterHandler,
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
	// 新增参数：激活奖励相关
	terminalRepo *repository.GormTerminalRepository,
	channelRepo *repository.GormChannelRepository,
	activationRewardService *service.ActivationRewardService,
	depositRecordRepo *repository.GormDepositCashbackRecordRepository,
	rewardRecordRepo *repository.GormActivationRewardRecordRepository,
	walletRepo *repository.GormWalletRepository,
	simCashbackRecordRepo *repository.GormSimCashbackRecordRepository,
	// 新增参数：商户类型计算
	merchantRepo *repository.GormMerchantRepository,
	merchantService *service.MerchantService,
	// 新增参数：流量费返现
	deviceFeeRepo *repository.GormDeviceFeeRepository,
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
	partitionJob := jobs.NewPartitionManagerJob(callbackRepo.GetDB())
	scheduler.AddJob("partition_manager", 24*time.Hour, partitionJob.Run)

	// 每日代扣任务（每24小时执行一次）
	jobs.SetupDeductionJobs(scheduler, deductionService, simCashbackService, deviceFeeRepo)

	// ============================================================
	// 新增：政策相关定时任务
	// ============================================================

	// 激活奖励检查任务（每天凌晨2点执行，这里用24小时间隔）
	rewardCheckJob := jobs.NewRewardCheckJob(terminalRepo, channelRepo, activationRewardService)
	scheduler.AddJob("reward_check", 24*time.Hour, rewardCheckJob.Run)

	// 押金返现入账任务（每10分钟）
	depositCashbackJob := jobs.NewDepositCashbackJob(depositRecordRepo, walletRepo)
	scheduler.AddJob("deposit_cashback_settle", 10*time.Minute, depositCashbackJob.Run)

	// 激活奖励入账任务（每10分钟）
	rewardSettleJob := jobs.NewActivationRewardSettleJob(rewardRecordRepo, walletRepo)
	scheduler.AddJob("activation_reward_settle", 10*time.Minute, rewardSettleJob.Run)

	// 流量费返现入账任务（每10分钟）
	simSettleJob := jobs.NewSimCashbackSettleJob(simCashbackRecordRepo, walletRepo)
	scheduler.AddJob("sim_cashback_settle", 10*time.Minute, simSettleJob.Run)

	// 商户类型计算任务（每天凌晨3点执行，这里用24小时间隔）
	merchantTypeJob := jobs.NewMerchantTypeCalculatorJob(merchantRepo, merchantService)
	scheduler.AddJob("merchant_type_calculator", 24*time.Hour, merchantTypeJob.Run)

	return scheduler
}

// setupRouter 设置路由
func setupRouter(
	callbackHandler *handler.CallbackHandler,
	deductionHandler *handler.DeductionHandler,
	terminalDistributeHandler *handler.TerminalDistributeHandler,
	simCashbackHandler *handler.SimCashbackHandler,
	goodsDeductionHandler *handler.GoodsDeductionHandler,
	authHandler *handler.AuthHandler,
	agentHandler *handler.AgentHandler,
	walletHandler *handler.WalletHandler,
	transactionHandler *handler.TransactionHandler,
	profitHandler *handler.ProfitHandler,
	dashboardHandler *handler.DashboardHandler,
	messageHandler *handler.MessageHandler,
	adminMessageHandler *handler.AdminMessageHandler,
	merchantHandler *handler.MerchantHandler,
	terminalHandler *handler.TerminalHandler,
	policyHandler *handler.PolicyHandler,
	agentChannelHandler *handler.AgentChannelHandler,
	chargingWalletHandler *handler.ChargingWalletHandler,
	settlementWalletHandler *handler.SettlementWalletHandler,
	taxChannelHandler *handler.TaxChannelHandler,
	uploadHandler *handler.UploadHandler,
	bannerHandler *handler.BannerHandler,
	posterHandler *handler.PosterHandler,
	authService *service.AuthService,
	metricsService *service.MetricsService,
	swaggerEnabled bool,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 全局中间件
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())

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

		// 注册货款代扣路由
		handler.RegisterGoodsDeductionRoutes(apiV1, goodsDeductionHandler)

		// 注册PC端新增路由
		handler.RegisterAgentRoutes(apiV1, agentHandler, authService)
		handler.RegisterWalletRoutes(apiV1, walletHandler, authService)
		handler.RegisterTransactionRoutes(apiV1, transactionHandler, authService)
		handler.RegisterProfitRoutes(apiV1, profitHandler, authService)
		handler.RegisterDashboardRoutes(apiV1, dashboardHandler, authService)
		handler.RegisterMessageRoutes(apiV1, messageHandler, authService)

		// 注册管理端消息路由
		handler.RegisterAdminMessageRoutes(apiV1, adminMessageHandler, authService)

		// 注册商户、终端、政策路由
		handler.RegisterMerchantRoutes(apiV1, merchantHandler, authService)
		handler.RegisterTerminalRoutes(apiV1, terminalHandler, authService)
		handler.RegisterPolicyRoutes(apiV1, policyHandler, authService)

		// 注册代理商通道路由
		handler.RegisterAgentChannelRoutes(apiV1, agentChannelHandler, authService)

		// 注册钱包相关路由
		handler.RegisterChargingWalletRoutes(apiV1, chargingWalletHandler, authService)
		handler.RegisterSettlementWalletRoutes(apiV1, settlementWalletHandler, authService)
		handler.RegisterTaxChannelRoutes(apiV1, taxChannelHandler, authService)

		// 注册营销模块路由（Banner、海报、上传）
		handler.RegisterUploadRoutes(apiV1, uploadHandler, authService)
		handler.RegisterBannerRoutes(apiV1, bannerHandler, authService)
		handler.RegisterPosterRoutes(apiV1, posterHandler, authService)
	}

	// Swagger UI (可通过环境变量关闭)
	if swaggerEnabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		log.Println("Swagger UI enabled at /swagger/index.html")
	}

	return router
}
