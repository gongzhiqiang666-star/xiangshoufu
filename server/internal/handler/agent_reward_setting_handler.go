package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"
)

// AgentRewardSettingHandler 代理商奖励配置处理器
type AgentRewardSettingHandler struct {
	service      *service.AgentRewardSettingService
	changeLogSvc *service.PriceChangeLogService
}

// NewAgentRewardSettingHandler 创建代理商奖励配置处理器
func NewAgentRewardSettingHandler(
	service *service.AgentRewardSettingService,
	changeLogSvc *service.PriceChangeLogService,
) *AgentRewardSettingHandler {
	return &AgentRewardSettingHandler{
		service:      service,
		changeLogSvc: changeLogSvc,
	}
}

// List 获取奖励配置列表
// @Summary 获取奖励配置列表
// @Tags 奖励配置管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/reward-settings [get]
func (h *AgentRewardSettingHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	settings, total, err := h.service.List(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, settings, total, page, pageSize)
}

// GetByID 获取奖励配置详情
// @Summary 获取奖励配置详情
// @Tags 奖励配置管理
// @Accept json
// @Produce json
// @Param id path int64 true "奖励配置ID"
// @Success 200 {object} models.AgentRewardSetting
// @Router /api/v1/reward-settings/{id} [get]
func (h *AgentRewardSettingHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	setting, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(c, "奖励配置不存在")
		return
	}

	response.Success(c, setting)
}

// Create 创建奖励配置
// @Summary 创建奖励配置
// @Tags 奖励配置管理
// @Accept json
// @Produce json
// @Param body body models.AgentRewardSettingRequest true "创建奖励配置请求"
// @Success 200 {object} models.AgentRewardSetting
// @Router /api/v1/reward-settings [post]
func (h *AgentRewardSettingHandler) Create(c *gin.Context) {
	var req models.AgentRewardSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取操作者信息
	operatorID := getOperatorID(c)
	operatorName := getOperatorName(c)
	source := getSource(c)

	setting, err := h.service.CreateFromTemplate(
		req.AgentID,
		req.TemplateID,
		req.RewardAmount,
		nil,
		operatorID,
		operatorName,
		source,
	)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, setting)
}

// UpdateActivation 更新激活奖励
// @Summary 更新激活奖励
// @Tags 奖励配置管理
// @Accept json
// @Produce json
// @Param id path int64 true "奖励配置ID"
// @Param body body models.UpdateActivationRewardRequest true "更新激活奖励请求"
// @Success 200 {object} models.AgentRewardSetting
// @Router /api/v1/reward-settings/{id}/activation [put]
func (h *AgentRewardSettingHandler) UpdateActivation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req models.UpdateActivationRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取操作者信息
	operatorID := getOperatorID(c)
	operatorName := getOperatorName(c)
	source := getSource(c)
	ipAddress := c.ClientIP()

	setting, err := h.service.UpdateActivationReward(id, &req, operatorID, operatorName, source, ipAddress)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, setting)
}

// GetChangeLogs 获取奖励配置调价记录
// @Summary 获取奖励配置调价记录
// @Tags 奖励配置管理
// @Accept json
// @Produce json
// @Param id path int64 true "奖励配置ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} models.PriceChangeLogListResponse
// @Router /api/v1/reward-settings/{id}/change-logs [get]
func (h *AgentRewardSettingHandler) GetChangeLogs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	resp, err := h.changeLogSvc.ListByRewardSetting(id, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resp)
}
