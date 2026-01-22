package handler

import (
	"net/http"
	"strconv"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

// PosterHandler 海报处理器
type PosterHandler struct {
	posterService service.PosterService
}

// NewPosterHandler 创建海报处理器实例
func NewPosterHandler(posterService service.PosterService) *PosterHandler {
	return &PosterHandler{
		posterService: posterService,
	}
}

// ========== 分类管理接口 ==========

// ListCategories 获取分类列表
// @Summary 获取海报分类列表
// @Tags 海报分类管理
// @Accept json
// @Produce json
// @Param status query int false "状态"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/poster-categories [get]
func (h *PosterHandler) ListCategories(c *gin.Context) {
	var req models.PosterCategoryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	categories, err := h.posterService.GetCategories(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	list := make([]*models.PosterCategoryResponse, 0, len(categories))
	for _, c := range categories {
		list = append(list, c.ToResponse())
	}

	response.Success(c, list)
}

// CreateCategory 创建分类
// @Summary 创建海报分类
// @Tags 海报分类管理
// @Accept json
// @Produce json
// @Param body body models.PosterCategoryCreateRequest true "创建请求"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/poster-categories [post]
func (h *PosterHandler) CreateCategory(c *gin.Context) {
	var req models.PosterCategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	category, err := h.posterService.CreateCategory(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}

	response.Success(c, category.ToResponse())
}

// UpdateCategory 更新分类
// @Summary 更新海报分类
// @Tags 海报分类管理
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param body body models.PosterCategoryUpdateRequest true "更新请求"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/poster-categories/{id} [put]
func (h *PosterHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	var req models.PosterCategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	category, err := h.posterService.UpdateCategory(id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	response.Success(c, category.ToResponse())
}

// DeleteCategory 删除分类
// @Summary 删除海报分类
// @Tags 海报分类管理
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/poster-categories/{id} [delete]
func (h *PosterHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := h.posterService.DeleteCategory(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// ========== 海报管理接口 ==========

// List 获取海报列表
// @Summary 获取海报列表
// @Tags 海报管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param category_id query int false "分类ID"
// @Param status query int false "状态"
// @Param keyword query string false "关键词"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/posters [get]
func (h *PosterHandler) List(c *gin.Context) {
	var req models.PosterListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	posters, total, err := h.posterService.GetList(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	list := make([]*models.PosterResponse, 0, len(posters))
	for _, p := range posters {
		list = append(list, p.ToResponse())
	}

	response.SuccessList(c, list, total)
}

// Get 获取海报详情
// @Summary 获取海报详情
// @Tags 海报管理
// @Accept json
// @Produce json
// @Param id path int true "海报ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/posters/{id} [get]
func (h *PosterHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	poster, err := h.posterService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "海报不存在")
		return
	}

	response.Success(c, poster.ToResponse())
}

// Create 创建海报
// @Summary 创建海报
// @Tags 海报管理
// @Accept json
// @Produce json
// @Param body body models.PosterCreateRequest true "创建请求"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/posters [post]
func (h *PosterHandler) Create(c *gin.Context) {
	var req models.PosterCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := getCurrentUserID(c)

	poster, err := h.posterService.Create(&req, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}

	response.Success(c, poster.ToResponse())
}

// Update 更新海报
// @Summary 更新海报
// @Tags 海报管理
// @Accept json
// @Produce json
// @Param id path int true "海报ID"
// @Param body body models.PosterUpdateRequest true "更新请求"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/posters/{id} [put]
func (h *PosterHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	var req models.PosterUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	poster, err := h.posterService.Update(id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	response.Success(c, poster.ToResponse())
}

// Delete 删除海报
// @Summary 删除海报
// @Tags 海报管理
// @Accept json
// @Produce json
// @Param id path int true "海报ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/posters/{id} [delete]
func (h *PosterHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := h.posterService.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// BatchImport 批量导入海报
// @Summary 批量导入海报
// @Tags 海报管理
// @Accept json
// @Produce json
// @Param body body models.PosterBatchImportRequest true "导入请求"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/posters/batch-import [post]
func (h *PosterHandler) BatchImport(c *gin.Context) {
	var req models.PosterBatchImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := getCurrentUserID(c)

	count, err := h.posterService.BatchImport(&req, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "导入失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"imported_count": count})
}

// ========== APP端接口 ==========

// GetActiveCategories 获取有效分类列表（APP端）
// @Summary 获取有效分类列表（含数量）
// @Tags 海报展示
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/posters/categories [get]
func (h *PosterHandler) GetActiveCategories(c *gin.Context) {
	categories, err := h.posterService.GetActiveCategories()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取列表失败: "+err.Error())
		return
	}

	response.Success(c, categories)
}

// GetActivePosters 获取有效海报列表（APP端）
// @Summary 获取有效海报列表
// @Tags 海报展示
// @Accept json
// @Produce json
// @Param category_id query int false "分类ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/posters [get]
func (h *PosterHandler) GetActivePosters(c *gin.Context) {
	var categoryID *int64
	if cidStr := c.Query("category_id"); cidStr != "" {
		cid, err := strconv.ParseInt(cidStr, 10, 64)
		if err == nil {
			categoryID = &cid
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	posters, total, err := h.posterService.GetActivePosters(categoryID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取列表失败: "+err.Error())
		return
	}

	// 转换为APP端响应格式
	list := make([]*models.AppPosterResponse, 0, len(posters))
	for _, p := range posters {
		list = append(list, p.ToAppResponse())
	}

	response.SuccessList(c, list, total)
}

// GetActivePosterDetail 获取海报详情（APP端）
// @Summary 获取海报详情
// @Tags 海报展示
// @Accept json
// @Produce json
// @Param id path int true "海报ID"
// @Success 200 {object} response.Response
// @Router /api/v1/posters/{id} [get]
func (h *PosterHandler) GetActivePosterDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	poster, err := h.posterService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "海报不存在")
		return
	}

	response.Success(c, poster.ToAppResponse())
}

// RecordDownload 记录下载
// @Summary 记录海报下载
// @Tags 海报展示
// @Accept json
// @Produce json
// @Param id path int true "海报ID"
// @Success 200 {object} response.Response
// @Router /api/v1/posters/{id}/download [post]
func (h *PosterHandler) RecordDownload(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	// 异步记录下载，不阻塞响应
	go h.posterService.RecordDownload(id)

	response.Success(c, nil)
}

// RecordShare 记录分享
// @Summary 记录海报分享
// @Tags 海报展示
// @Accept json
// @Produce json
// @Param id path int true "海报ID"
// @Success 200 {object} response.Response
// @Router /api/v1/posters/{id}/share [post]
func (h *PosterHandler) RecordShare(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	// 异步记录分享，不阻塞响应
	go h.posterService.RecordShare(id)

	response.Success(c, nil)
}
