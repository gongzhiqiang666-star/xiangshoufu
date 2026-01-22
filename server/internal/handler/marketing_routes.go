package handler

import (
	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// RegisterBannerRoutes 注册Banner路由
func RegisterBannerRoutes(router *gin.RouterGroup, bannerHandler *BannerHandler, authService *service.AuthService) {
	// APP端接口（公开或需要登录）
	banners := router.Group("/banners")
	{
		banners.GET("", bannerHandler.GetActive)              // 获取有效Banner列表
		banners.POST("/:id/click", bannerHandler.RecordClick) // 记录点击
	}

	// 管理端接口（需要认证）
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(authService))
	{
		adminBanners := admin.Group("/banners")
		{
			adminBanners.GET("", bannerHandler.List)                  // 获取列表
			adminBanners.GET("/:id", bannerHandler.Get)               // 获取详情
			adminBanners.POST("", bannerHandler.Create)               // 创建
			adminBanners.PUT("/:id", bannerHandler.Update)            // 更新
			adminBanners.DELETE("/:id", bannerHandler.Delete)         // 删除
			adminBanners.PUT("/:id/status", bannerHandler.UpdateStatus) // 状态切换
			adminBanners.PUT("/sort", bannerHandler.UpdateSort)       // 批量排序
		}
	}
}

// RegisterPosterRoutes 注册海报路由
func RegisterPosterRoutes(router *gin.RouterGroup, posterHandler *PosterHandler, authService *service.AuthService) {
	// APP端接口
	posters := router.Group("/posters")
	{
		posters.GET("/categories", posterHandler.GetActiveCategories)       // 获取分类列表
		posters.GET("", posterHandler.GetActivePosters)                     // 获取海报列表
		posters.GET("/:id", posterHandler.GetActivePosterDetail)            // 获取海报详情
		posters.POST("/:id/download", posterHandler.RecordDownload)         // 记录下载
		posters.POST("/:id/share", posterHandler.RecordShare)               // 记录分享
	}

	// 管理端接口（需要认证）
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(authService))
	{
		// 分类管理
		categories := admin.Group("/poster-categories")
		{
			categories.GET("", posterHandler.ListCategories)          // 获取列表
			categories.POST("", posterHandler.CreateCategory)         // 创建
			categories.PUT("/:id", posterHandler.UpdateCategory)      // 更新
			categories.DELETE("/:id", posterHandler.DeleteCategory)   // 删除
		}

		// 海报管理
		adminPosters := admin.Group("/posters")
		{
			adminPosters.GET("", posterHandler.List)                         // 获取列表
			adminPosters.GET("/:id", posterHandler.Get)                      // 获取详情
			adminPosters.POST("", posterHandler.Create)                      // 创建
			adminPosters.PUT("/:id", posterHandler.Update)                   // 更新
			adminPosters.DELETE("/:id", posterHandler.Delete)                // 删除
			adminPosters.POST("/batch-import", posterHandler.BatchImport)    // 批量导入
		}
	}
}

// RegisterUploadRoutes 注册上传路由
func RegisterUploadRoutes(router *gin.RouterGroup, uploadHandler *UploadHandler, authService *service.AuthService) {
	// 上传接口（需要认证）
	upload := router.Group("/upload")
	upload.Use(middleware.AuthMiddleware(authService))
	{
		upload.POST("/image", uploadHandler.UploadImage) // 上传图片
	}
}
