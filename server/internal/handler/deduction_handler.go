package handler

import (
	"net/http"
	"strconv"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// DeductionHandler 代扣管理Handler
type DeductionHandler struct {
	deductionService *service.DeductionService
}

// NewDeductionHandler 创建代扣Handler
func NewDeductionHandler(deductionService *service.DeductionService) *DeductionHandler {
	return &DeductionHandler{
		deductionService: deductionService,
	}
}

// CreateDeductionPlanRequest API请求
type CreateDeductionPlanRequest struct {
	DeductorID   int64  `json:"deductor_id" binding:"required"`   // 扣款方代理商ID
	DeducteeID   int64  `json:"deductee_id" binding:"required"`   // 被扣款方代理商ID
	PlanType     int16  `json:"plan_type" binding:"required"`     // 计划类型：1货款 2伙伴 3押金
	TotalAmount  int64  `json:"total_amount" binding:"required"`  // 总金额（分）
	TotalPeriods int    `json:"total_periods" binding:"required"` // 总期数
	RelatedType  string `json:"related_type"`                     // 关联类型
	RelatedID    int64  `json:"related_id"`                       // 关联ID
	Remark       string `json:"remark"`                           // 备注
}

// CreateDeductionPlan 创建代扣计划
// @Summary 创建代扣计划
// @Description 创建代扣计划，支持伙伴代扣（任意代理商之间，不限层级关系）
// @Tags 代扣管理
// @Accept json
// @Produce json
// @Param request body CreateDeductionPlanRequest true "创建代扣计划请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/deduction/plans [post]
func (h *DeductionHandler) CreateDeductionPlan(c *gin.Context) {
	var req CreateDeductionPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前用户ID（从JWT或session中获取）
	createdBy := getCurrentUserID(c)

	serviceReq := &service.CreateDeductionPlanRequest{
		DeductorID:   req.DeductorID,
		DeducteeID:   req.DeducteeID,
		PlanType:     req.PlanType,
		TotalAmount:  req.TotalAmount,
		TotalPeriods: req.TotalPeriods,
		RelatedType:  req.RelatedType,
		RelatedID:    req.RelatedID,
		Remark:       req.Remark,
		CreatedBy:    createdBy,
	}

	plan, err := h.deductionService.CreateDeductionPlan(serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    plan,
	})
}

// GetDeductionPlan 获取代扣计划详情
// @Summary 获取代扣计划详情
// @Tags 代扣管理
// @Produce json
// @Param id path int true "计划ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/deduction/plans/{id} [get]
func (h *DeductionHandler) GetDeductionPlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	detail, err := h.deductionService.GetPlanByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	// 转换代扣记录为前端友好格式
	records := make([]gin.H, 0, len(detail.Records))
	for _, r := range detail.Records {
		records = append(records, gin.H{
			"id":           r.ID,
			"period_num":   r.PeriodNum,
			"amount":       r.Amount,
			"amount_yuan":  float64(r.Amount) / 100,
			"status":       r.Status,
			"status_name":  DeductionRecordStatusDesc(r.Status),
			"scheduled_at": r.ScheduledAt,
			"executed_at":  r.DeductedAt,
			"fail_reason":  r.FailReason,
		})
	}

	plan := detail.Plan
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"plan": gin.H{
				"id":               plan.ID,
				"plan_no":          plan.PlanNo,
				"deductor_id":      plan.DeductorID,
				"deductee_id":      plan.DeducteeID,
				"plan_type":        plan.PlanType,
				"plan_type_name":   DeductionPlanTypeDesc(plan.PlanType),
				"total_amount":     plan.TotalAmount,
				"total_amount_yuan": float64(plan.TotalAmount) / 100,
				"deducted_amount":  plan.DeductedAmount,
				"deducted_yuan":    float64(plan.DeductedAmount) / 100,
				"remaining_amount": plan.RemainingAmount,
				"remaining_yuan":   float64(plan.RemainingAmount) / 100,
				"total_periods":    plan.TotalPeriods,
				"current_period":   plan.CurrentPeriod,
				"period_amount":    plan.PeriodAmount,
				"period_yuan":      float64(plan.PeriodAmount) / 100,
				"status":           plan.Status,
				"status_name":      DeductionPlanStatusDesc(plan.Status),
				"related_type":     plan.RelatedType,
				"related_id":       plan.RelatedID,
				"remark":           plan.Remark,
				"created_at":       plan.CreatedAt,
				"updated_at":       plan.UpdatedAt,
			},
			"records": records,
		},
	})
}

// PauseDeductionPlan 暂停代扣计划
// @Summary 暂停代扣计划
// @Tags 代扣管理
// @Produce json
// @Param id path int true "计划ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/deduction/plans/{id}/pause [post]
func (h *DeductionHandler) PauseDeductionPlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	if err := h.deductionService.PauseDeductionPlan(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "暂停成功",
	})
}

// ResumeDeductionPlan 恢复代扣计划
// @Summary 恢复代扣计划
// @Tags 代扣管理
// @Produce json
// @Param id path int true "计划ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/deduction/plans/{id}/resume [post]
func (h *DeductionHandler) ResumeDeductionPlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	if err := h.deductionService.ResumeDeductionPlan(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "恢复成功",
	})
}

// CancelDeductionPlan 取消代扣计划
// @Summary 取消代扣计划
// @Tags 代扣管理
// @Produce json
// @Param id path int true "计划ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/deduction/plans/{id}/cancel [post]
func (h *DeductionHandler) CancelDeductionPlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	if err := h.deductionService.CancelDeductionPlan(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "取消成功",
	})
}

// CreateDeductionChainRequest 创建代扣链请求
type CreateDeductionChainRequest struct {
	DistributeID   int64   `json:"distribute_id" binding:"required"` // 终端下发记录ID
	TerminalSN     string  `json:"terminal_sn" binding:"required"`   // 终端SN
	AgentPath      []int64 `json:"agent_path" binding:"required"`    // 代理商路径
	TotalAmount    int64   `json:"total_amount" binding:"required"`  // 总金额
	TotalPeriods   int     `json:"total_periods" binding:"required"` // 总期数
}

// CreateDeductionChain 创建代扣链（跨级下发）
// @Summary 创建代扣链
// @Description 跨级下发时系统自动按层级生成A→B→C的货款代扣链
// @Tags 代扣管理
// @Accept json
// @Produce json
// @Param request body CreateDeductionChainRequest true "创建代扣链请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/deduction/chains [post]
func (h *DeductionHandler) CreateDeductionChain(c *gin.Context) {
	var req CreateDeductionChainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	createdBy := getCurrentUserID(c)

	serviceReq := &service.CreateDeductionChainRequest{
		DistributeID: req.DistributeID,
		TerminalSN:   req.TerminalSN,
		AgentPath:    req.AgentPath,
		TotalAmount:  req.TotalAmount,
		TotalPeriods: req.TotalPeriods,
		CreatedBy:    createdBy,
	}

	chain, err := h.deductionService.CreateDeductionChain(serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    chain,
	})
}

// ExecuteDailyDeduction 手动触发每日代扣（管理员接口）
// @Summary 手动触发每日代扣
// @Description 管理员手动触发每日代扣任务
// @Tags 代扣管理
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/deduction/execute [post]
func (h *DeductionHandler) ExecuteDailyDeduction(c *gin.Context) {
	if err := h.deductionService.ExecuteDailyDeduction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "执行成功",
	})
}

// RegisterDeductionRoutes 注册代扣管理路由
func RegisterDeductionRoutes(r *gin.RouterGroup, h *DeductionHandler) {
	deduction := r.Group("/deduction")
	{
		// 代扣计划
		deduction.POST("/plans", h.CreateDeductionPlan)
		deduction.GET("/plans/:id", h.GetDeductionPlan)
		deduction.POST("/plans/:id/pause", h.PauseDeductionPlan)
		deduction.POST("/plans/:id/resume", h.ResumeDeductionPlan)
		deduction.POST("/plans/:id/cancel", h.CancelDeductionPlan)

		// 代扣链
		deduction.POST("/chains", h.CreateDeductionChain)

		// 管理员接口
		deduction.POST("/execute", h.ExecuteDailyDeduction)
	}
}

// DeductionPlanTypeDesc 代扣计划类型描述
func DeductionPlanTypeDesc(planType int16) string {
	switch planType {
	case models.DeductionPlanTypeGoods:
		return "货款代扣"
	case models.DeductionPlanTypePartner:
		return "伙伴代扣"
	case models.DeductionPlanTypeDeposit:
		return "押金代扣"
	default:
		return "未知类型"
	}
}

// DeductionPlanStatusDesc 代扣计划状态描述
func DeductionPlanStatusDesc(status int16) string {
	switch status {
	case models.DeductionPlanStatusActive:
		return "进行中"
	case models.DeductionPlanStatusCompleted:
		return "已完成"
	case models.DeductionPlanStatusPaused:
		return "已暂停"
	case models.DeductionPlanStatusCancelled:
		return "已取消"
	default:
		return "未知状态"
	}
}

// DeductionRecordStatusDesc 代扣记录状态描述
func DeductionRecordStatusDesc(status int16) string {
	switch status {
	case models.DeductionRecordStatusPending:
		return "待扣款"
	case models.DeductionRecordStatusSuccess:
		return "扣款成功"
	case models.DeductionRecordStatusFailed:
		return "扣款失败"
	case models.DeductionRecordStatusPartialSuccess:
		return "部分成功"
	default:
		return "未知状态"
	}
}
