package handler

import (
	"net/http"

	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

// UploadHandler 上传处理器
type UploadHandler struct {
	uploadService service.UploadService
}

// NewUploadHandler 创建上传处理器实例
func NewUploadHandler(uploadService service.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// UploadImage 上传图片
// @Summary 上传图片
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "图片文件"
// @Param module formData string false "模块 banner/poster"
// @Success 200 {object} response.Response
// @Router /api/v1/upload/image [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择要上传的文件")
		return
	}

	// 获取模块参数
	module := c.DefaultPostForm("module", "")

	// 获取当前用户ID
	userID := getCurrentUserID(c)

	// 调用上传服务
	result, err := h.uploadService.UploadImage(file, module, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, result)
}
