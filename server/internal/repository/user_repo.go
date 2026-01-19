package repository

import (
	"time"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	FindByID(id int64) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	UpdateLastLogin(id int64, ip string) error
}

// GormUserRepository GORM实现
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository 创建用户仓库
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *GormUserRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *GormUserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *GormUserRepository) UpdateLastLogin(id int64, ip string) error {
	now := time.Now()
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": ip,
	}).Error
}

// RefreshTokenRepository 刷新令牌仓库接口
type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	FindByToken(token string) (*models.RefreshToken, error)
	DeleteByToken(token string) error
	DeleteByUserID(userID int64) error
	DeleteExpired() (int64, error)
}

// GormRefreshTokenRepository GORM实现
type GormRefreshTokenRepository struct {
	db *gorm.DB
}

// NewGormRefreshTokenRepository 创建刷新令牌仓库
func NewGormRefreshTokenRepository(db *gorm.DB) *GormRefreshTokenRepository {
	return &GormRefreshTokenRepository{db: db}
}

func (r *GormRefreshTokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *GormRefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&rt).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &rt, err
}

func (r *GormRefreshTokenRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

func (r *GormRefreshTokenRepository) DeleteByUserID(userID int64) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

func (r *GormRefreshTokenRepository) DeleteExpired() (int64, error) {
	result := r.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{})
	return result.RowsAffected, result.Error
}

// LoginLogRepository 登录日志仓库接口
type LoginLogRepository interface {
	Create(log *models.LoginLog) error
	FindByUserID(userID int64, limit, offset int) ([]*models.LoginLog, error)
}

// GormLoginLogRepository GORM实现
type GormLoginLogRepository struct {
	db *gorm.DB
}

// NewGormLoginLogRepository 创建登录日志仓库
func NewGormLoginLogRepository(db *gorm.DB) *GormLoginLogRepository {
	return &GormLoginLogRepository{db: db}
}

func (r *GormLoginLogRepository) Create(log *models.LoginLog) error {
	return r.db.Create(log).Error
}

func (r *GormLoginLogRepository) FindByUserID(userID int64, limit, offset int) ([]*models.LoginLog, error) {
	var logs []*models.LoginLog
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&logs).Error
	return logs, err
}
