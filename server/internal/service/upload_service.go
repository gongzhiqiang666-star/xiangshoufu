package service

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

// UploadService 上传服务接口
type UploadService interface {
	UploadImage(file *multipart.FileHeader, module string, uploadedBy int64) (*UploadResult, error)
	GetImageConfig(module string) models.ImageConfig
}

// UploadResult 上传结果
type UploadResult struct {
	ID           int64  `json:"id"`
	OriginalName string `json:"original_name"`
	FileURL      string `json:"file_url"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	FileSize     int64  `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	MimeType     string `json:"mime_type"`
}

// uploadService 上传服务实现
type uploadService struct {
	repo       repository.UploadedFileRepository
	uploadDir  string // 上传目录
	baseURL    string // 访问基础URL
}

// NewUploadService 创建上传服务实例
func NewUploadService(repo repository.UploadedFileRepository, uploadDir, baseURL string) UploadService {
	return &uploadService{
		repo:      repo,
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

// UploadImage 上传图片
func (s *uploadService) UploadImage(file *multipart.FileHeader, module string, uploadedBy int64) (*UploadResult, error) {
	// 获取配置
	config := s.GetImageConfig(module)

	// 验证文件大小
	if file.Size > config.MaxSize {
		return nil, fmt.Errorf("文件大小超过限制，最大允许 %dMB", config.MaxSize/1024/1024)
	}

	// 验证文件类型
	mimeType := file.Header.Get("Content-Type")
	if !s.isAllowedType(mimeType, config.AllowedTypes) {
		return nil, fmt.Errorf("不支持的文件类型: %s", mimeType)
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 解码图片获取尺寸
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("解码图片失败: %w", err)
	}

	// 重置文件读取位置
	src.Seek(0, 0)

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 生成文件名和路径
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		ext = s.getExtByMimeType(mimeType)
	}
	storedName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dateDir := time.Now().Format("2006/01/02")
	relativePath := filepath.Join(module, dateDir, storedName)
	fullPath := filepath.Join(s.uploadDir, relativePath)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	// 处理图片（压缩/调整尺寸）
	processedImg := img
	needResize := width > config.MaxWidth || height > config.MaxHeight
	if needResize {
		processedImg = imaging.Fit(img, config.MaxWidth, config.MaxHeight, imaging.Lanczos)
		bounds = processedImg.Bounds()
		width = bounds.Dx()
		height = bounds.Dy()
	}

	// 保存处理后的图片
	if err := s.saveImage(processedImg, fullPath, config.Quality); err != nil {
		return nil, fmt.Errorf("保存图片失败: %w", err)
	}

	// 获取保存后的文件大小
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}
	fileSize := fileInfo.Size()

	// 生成访问URL
	fileURL := fmt.Sprintf("%s/%s", s.baseURL, relativePath)

	// 生成缩略图（如果需要）
	var thumbnailURL string
	if config.GenerateThumb && config.ThumbWidth > 0 {
		thumbName := fmt.Sprintf("thumb_%s", storedName)
		thumbRelativePath := filepath.Join(module, dateDir, thumbName)
		thumbFullPath := filepath.Join(s.uploadDir, thumbRelativePath)

		thumbImg := imaging.Resize(img, config.ThumbWidth, 0, imaging.Lanczos)
		if err := s.saveImage(thumbImg, thumbFullPath, config.Quality); err != nil {
			// 缩略图生成失败不影响主流程，记录日志即可
			fmt.Printf("生成缩略图失败: %v\n", err)
		} else {
			thumbnailURL = fmt.Sprintf("%s/%s", s.baseURL, thumbRelativePath)
		}
	}

	// 保存上传记录
	uploadedFile := &models.UploadedFile{
		OriginalName: file.Filename,
		StoredName:   storedName,
		FilePath:     relativePath,
		FileURL:      fileURL,
		FileSize:     fileSize,
		MimeType:     mimeType,
		Width:        width,
		Height:       height,
		Module:       module,
		UploadedBy:   uploadedBy,
	}

	if err := s.repo.Create(uploadedFile); err != nil {
		// 删除已上传的文件
		os.Remove(fullPath)
		return nil, fmt.Errorf("保存上传记录失败: %w", err)
	}

	return &UploadResult{
		ID:           uploadedFile.ID,
		OriginalName: file.Filename,
		FileURL:      fileURL,
		ThumbnailURL: thumbnailURL,
		FileSize:     fileSize,
		Width:        width,
		Height:       height,
		MimeType:     mimeType,
	}, nil
}

// GetImageConfig 获取图片配置
func (s *uploadService) GetImageConfig(module string) models.ImageConfig {
	switch module {
	case "banner":
		return models.BannerImageConfig
	case "poster":
		return models.PosterImageConfig
	default:
		// 默认配置
		return models.ImageConfig{
			MaxSize:       5 * 1024 * 1024,
			MaxWidth:      2000,
			MaxHeight:     2000,
			Quality:       85,
			AllowedTypes:  []string{"image/jpeg", "image/png", "image/webp", "image/gif"},
			GenerateThumb: false,
		}
	}
}

// isAllowedType 检查是否允许的类型
func (s *uploadService) isAllowedType(mimeType string, allowedTypes []string) bool {
	for _, t := range allowedTypes {
		if t == mimeType {
			return true
		}
	}
	return false
}

// getExtByMimeType 根据MIME类型获取扩展名
func (s *uploadService) getExtByMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ".jpg"
	}
}

// saveImage 保存图片
func (s *uploadService) saveImage(img image.Image, path string, quality int) error {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".jpg", ".jpeg":
		return imaging.Save(img, path, imaging.JPEGQuality(quality))
	case ".png":
		return imaging.Save(img, path, imaging.PNGCompressionLevel(0))
	case ".webp":
		// imaging 不直接支持 webp，转为 jpg
		path = strings.TrimSuffix(path, ext) + ".jpg"
		return imaging.Save(img, path, imaging.JPEGQuality(quality))
	default:
		return imaging.Save(img, path, imaging.JPEGQuality(quality))
	}
}

// copyFile 复制文件
func copyFile(src io.Reader, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
