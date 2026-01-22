package handler

import (
	"net/http"
	"strconv"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

// BannerHandler Banner处理器
type BannerHandler struct {
	bannerService service.BannerService
}

// NewBannerHandler 创建Banner处理器实例
func NewBannerHandler(bannerService service.BannerService) *BannerHandler {
	return &BannerHandler{
		bannerService: bannerService,
	}
}

// ========== 管理端接口 ==========

// List 获取Banner列表
// @Summary 获取Banner列表
// @Tags Banner管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param status query int false "状态"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/banners [get]
func (h *BannerHandler) List(c *gin.Context) {
	var req models.BannerListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	banners, total, err := h.bannerService.GetList(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	list := make([]*models.BannerResponse, 0, len(banners))
	for _, b := range banners {
		list = append(list, b.ToResponse())
	}

	response.SuccessList(c, list, total)
}

// Get 获取Banner详情
// @Summary 获取Banner详情
// @Tags Banner管理
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/banners/{id} [get]
func (h *BannerHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	banner, err := h.bannerService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Banner不存在")
		return
	}

	response.Success(c, banner.ToResponse())
}

// Create 创建Banner
// @Summary 创建Banner
// @Tags Banner管理
// @Accept json
// @Produce json
// @Param body body models.BannerCreateRequest true "创建请求"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/banners [post]
func (h *BannerHandler) Create(c *gin.Context) {
	var req models.BannerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 从上下文获取当前用户ID
	userID := getCurrentUserID(c)

	banner, err := h.bannerService.Create(&req, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}

	response.Success(c, banner.ToResponse())
}

// Update 更新Banner
// @Summary 更新Banner
// @Tags Banner管理
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Param body body models.BannerUpdateRequest true "更新请求"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/banners/{id} [put]
func (h *BannerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	var req models.BannerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	banner, err := h.bannerService.Update(id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	response.Success(c, banner.ToResponse())
}

// Delete 删除Banner
// @Summary 删除Banner
// @Tags Banner管理
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/banners/{id} [delete]
func (h *BannerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := h.bannerService.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateStatus 更新状态
// @Summary 更新Banner状态
// @Tags Banner管理
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Param body body models.BannerStatusRequest true "状态请求"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/banners/{id}/status [put]
func (h *BannerHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	var req models.BannerStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := h.bannerService.UpdateStatus(id, req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新状态失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateSort 批量更新排序
// @Summary 批量更新Banner排序
// @Tags Banner管理
// @Accept json
// @Produce json
// @Param body body models.BannerSortRequest true "排序请求"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/banners/sort [put]
func (h *BannerHandler) UpdateSort(c *gin.Context) {
	var req models.BannerSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := h.bannerService.UpdateSortOrder(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新排序失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// ========== APP端接口 ==========

// GetActive 获取有效Banner列表（APP端）
// @Summary 获取有效Banner列表
// @Tags Banner展示
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/banners [get]
func (h *BannerHandler) GetActive(c *gin.Context) {
	banners, err := h.bannerService.GetActiveBanners()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取列表失败: "+err.Error())
		return
	}

	// 转换为APP端响应格式
	list := make([]*models.AppBannerResponse, 0, len(banners))
	for _, b := range banners {
		list = append(list, b.ToAppResponse())
	}

	response.Success(c, list)
}

// RecordClick 记录点击
// @Summary 记录Banner点击
// @Tags Banner展示
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.Response
// @Router /api/v1/banners/{id}/click [post]
func (h *BannerHandler) RecordClick(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	// 异步记录点击，不阻塞响应
	go h.bannerService.RecordClick(id)

	response.Success(c, nil)
}

// getCurrentUserID 从上下文获取当前用户ID
func getCurrentUserID(c *gin.Context) int64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	if id, ok := userID.(int64); ok {
		return id
	}
	return 0
}
