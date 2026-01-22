package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/jobs"
)

// JobHandler 任务管理处理器
type JobHandler struct {
	configRepo   repository.JobConfigRepository
	logRepo      repository.JobExecutionLogRepository
	jobRegistry  map[string]*jobs.JobWrapper // 注册的任务
}

// NewJobHandler 创建任务管理处理器
func NewJobHandler(
	configRepo repository.JobConfigRepository,
	logRepo repository.JobExecutionLogRepository,
) *JobHandler {
	return &JobHandler{
		configRepo:  configRepo,
		logRepo:     logRepo,
		jobRegistry: make(map[string]*jobs.JobWrapper),
	}
}

// RegisterJob 注册任务
func (h *JobHandler) RegisterJob(wrapper *jobs.JobWrapper) {
	h.jobRegistry[wrapper.GetJobName()] = wrapper
}

// RegisterRoutes 注册路由
func (h *JobHandler) RegisterRoutes(rg *gin.RouterGroup) {
	jobGroup := rg.Group("/jobs")
	{
		jobGroup.GET("", h.ListJobs)
		jobGroup.GET("/:name", h.GetJob)
		jobGroup.PUT("/:name/config", h.UpdateJobConfig)
		jobGroup.POST("/:name/trigger", h.TriggerJob)
		jobGroup.PUT("/:name/enable", h.EnableJob)
	}

	logGroup := rg.Group("/job-logs")
	{
		logGroup.GET("", h.ListJobLogs)
		logGroup.GET("/:id", h.GetJobLog)
		logGroup.GET("/stats", h.GetJobStats)
	}
}

// JobListResponse 任务列表响应
type JobListResponse struct {
	ID              int64     `json:"id"`
	JobName         string    `json:"job_name"`
	JobDesc         string    `json:"job_desc"`
	IntervalSeconds int       `json:"interval_seconds"`
	IsEnabled       bool      `json:"is_enabled"`
	MaxRetries      int       `json:"max_retries"`
	AlertThreshold  int       `json:"alert_threshold"`
	IsRunning       bool      `json:"is_running"`
	LastRunAt       *time.Time `json:"last_run_at"`
	LastStatus      *int16    `json:"last_status"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ListJobs 任务列表
// @Summary 获取任务列表
// @Tags 任务管理
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/jobs [get]
func (h *JobHandler) ListJobs(c *gin.Context) {
	configs, err := h.configRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询任务配置失败"})
		return
	}

	var list []JobListResponse
	for _, config := range configs {
		item := JobListResponse{
			ID:              config.ID,
			JobName:         config.JobName,
			JobDesc:         config.JobDesc,
			IntervalSeconds: config.IntervalSeconds,
			IsEnabled:       config.IsEnabled,
			MaxRetries:      config.MaxRetries,
			AlertThreshold:  config.AlertThreshold,
			UpdatedAt:       config.UpdatedAt,
		}

		// 检查是否正在运行
		if wrapper, ok := h.jobRegistry[config.JobName]; ok {
			item.IsRunning = wrapper.IsRunning()
		}

		// 获取最新执行日志
		if latestLog, err := h.logRepo.FindLatestByJobName(config.JobName); err == nil && latestLog != nil {
			item.LastRunAt = &latestLog.StartedAt
			item.LastStatus = &latestLog.Status
		}

		list = append(list, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": list,
	})
}

// JobDetailResponse 任务详情响应
type JobDetailResponse struct {
	Config     *models.JobConfig      `json:"config"`
	IsRunning  bool                   `json:"is_running"`
	LatestLogs []*models.JobExecutionLog `json:"latest_logs"`
}

// GetJob 任务详情
// @Summary 获取任务详情
// @Tags 任务管理
// @Param name path string true "任务名称"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/jobs/{name} [get]
func (h *JobHandler) GetJob(c *gin.Context) {
	jobName := c.Param("name")

	config, err := h.configRepo.FindByName(jobName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	resp := JobDetailResponse{
		Config: config,
	}

	// 检查是否正在运行
	if wrapper, ok := h.jobRegistry[jobName]; ok {
		resp.IsRunning = wrapper.IsRunning()
	}

	// 获取最近10条执行日志
	logs, _ := h.logRepo.FindByJobName(jobName, 10, 0)
	resp.LatestLogs = logs

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": resp,
	})
}

// UpdateJobConfigRequest 更新任务配置请求
type UpdateJobConfigRequest struct {
	JobDesc         string `json:"job_desc"`
	IntervalSeconds int    `json:"interval_seconds"`
	MaxRetries      int    `json:"max_retries"`
	RetryInterval   int    `json:"retry_interval"`
	AlertThreshold  int    `json:"alert_threshold"`
	TimeoutSeconds  int    `json:"timeout_seconds"`
}

// UpdateJobConfig 更新任务配置
// @Summary 更新任务配置
// @Tags 任务管理
// @Param name path string true "任务名称"
// @Param body body UpdateJobConfigRequest true "配置信息"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/jobs/{name}/config [put]
func (h *JobHandler) UpdateJobConfig(c *gin.Context) {
	jobName := c.Param("name")

	config, err := h.configRepo.FindByName(jobName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	var req UpdateJobConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 更新配置
	if req.JobDesc != "" {
		config.JobDesc = req.JobDesc
	}
	if req.IntervalSeconds > 0 {
		config.IntervalSeconds = req.IntervalSeconds
	}
	if req.MaxRetries > 0 {
		config.MaxRetries = req.MaxRetries
	}
	if req.RetryInterval > 0 {
		config.RetryInterval = req.RetryInterval
	}
	if req.AlertThreshold > 0 {
		config.AlertThreshold = req.AlertThreshold
	}
	if req.TimeoutSeconds > 0 {
		config.TimeoutSeconds = req.TimeoutSeconds
	}

	if err := h.configRepo.Update(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// TriggerJob 手动触发任务
// @Summary 手动触发任务
// @Tags 任务管理
// @Param name path string true "任务名称"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/jobs/{name}/trigger [post]
func (h *JobHandler) TriggerJob(c *gin.Context) {
	jobName := c.Param("name")

	wrapper, ok := h.jobRegistry[jobName]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在或未注册"})
		return
	}

	if wrapper.IsRunning() {
		c.JSON(http.StatusConflict, gin.H{"error": "任务正在运行中"})
		return
	}

	// 异步执行
	go wrapper.RunManual()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "任务已触发",
	})
}

// EnableJobRequest 启用/禁用任务请求
type EnableJobRequest struct {
	IsEnabled bool `json:"is_enabled"`
}

// EnableJob 启用/禁用任务
// @Summary 启用/禁用任务
// @Tags 任务管理
// @Param name path string true "任务名称"
// @Param body body EnableJobRequest true "启用状态"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/jobs/{name}/enable [put]
func (h *JobHandler) EnableJob(c *gin.Context) {
	jobName := c.Param("name")

	var req EnableJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.configRepo.UpdateEnabled(jobName, req.IsEnabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	status := "启用"
	if !req.IsEnabled {
		status = "禁用"
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "任务已" + status,
	})
}

// ListJobLogsRequest 执行日志列表请求
type ListJobLogsRequest struct {
	JobName   string `form:"job_name"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Status    int16  `form:"status"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
}

// ListJobLogs 执行日志列表
// @Summary 获取执行日志列表
// @Tags 任务管理
// @Param job_name query string false "任务名称"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/job-logs [get]
func (h *JobHandler) ListJobLogs(c *gin.Context) {
	var req ListJobLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	var logs []*models.JobExecutionLog
	var total int64
	var err error

	if req.StartDate != "" && req.EndDate != "" {
		startDate, _ := time.Parse("2006-01-02", req.StartDate)
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		endDate = endDate.Add(24 * time.Hour) // 包含结束日期

		logs, err = h.logRepo.FindByJobNameAndDateRange(req.JobName, startDate, endDate, req.PageSize, offset)
		if err == nil {
			total, _ = h.logRepo.CountByJobNameAndDateRange(req.JobName, startDate, endDate)
		}
	} else {
		logs, err = h.logRepo.FindByJobName(req.JobName, req.PageSize, offset)
		if err == nil {
			total, _ = h.logRepo.CountByJobName(req.JobName)
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// GetJobLog 执行日志详情
// @Summary 获取执行日志详情
// @Tags 任务管理
// @Param id path int true "日志ID"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/job-logs/{id} [get]
func (h *JobHandler) GetJobLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	log, err := h.logRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "日志不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": log,
	})
}

// GetJobStatsRequest 任务统计请求
type GetJobStatsRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

// GetJobStats 任务执行统计
// @Summary 获取任务执行统计
// @Tags 任务管理
// @Param start_date query string true "开始日期"
// @Param end_date query string true "结束日期"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/job-logs/stats [get]
func (h *JobHandler) GetJobStats(c *gin.Context) {
	var req GetJobStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供开始日期和结束日期"})
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)
	endDate = endDate.Add(24 * time.Hour)

	stats, err := h.logRepo.GetStats(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
	})
}
