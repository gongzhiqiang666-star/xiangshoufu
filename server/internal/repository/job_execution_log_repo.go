package repository

import (
	"time"

	"gorm.io/gorm"
	"xiangshoufu/internal/models"
)

// JobExecutionLogRepository 任务执行日志仓库接口
type JobExecutionLogRepository interface {
	Create(log *models.JobExecutionLog) error
	Update(log *models.JobExecutionLog) error
	FindByID(id int64) (*models.JobExecutionLog, error)
	FindByJobName(jobName string, limit, offset int) ([]*models.JobExecutionLog, error)
	FindByJobNameAndDateRange(jobName string, startDate, endDate time.Time, limit, offset int) ([]*models.JobExecutionLog, error)
	CountByJobName(jobName string) (int64, error)
	CountByJobNameAndDateRange(jobName string, startDate, endDate time.Time) (int64, error)
	FindLatestByJobName(jobName string) (*models.JobExecutionLog, error)
	FindRunningJobs() ([]*models.JobExecutionLog, error)
	GetStats(startDate, endDate time.Time) ([]*JobExecutionStats, error)
	DeleteOlderThan(days int) (int64, error)
}

// JobExecutionStats 任务执行统计
type JobExecutionStats struct {
	JobName      string  `json:"job_name"`
	TotalCount   int64   `json:"total_count"`
	SuccessCount int64   `json:"success_count"`
	FailCount    int64   `json:"fail_count"`
	AvgDuration  float64 `json:"avg_duration_ms"`
	MaxDuration  int64   `json:"max_duration_ms"`
	MinDuration  int64   `json:"min_duration_ms"`
}

// GormJobExecutionLogRepository 任务执行日志仓库GORM实现
type GormJobExecutionLogRepository struct {
	db *gorm.DB
}

// NewGormJobExecutionLogRepository 创建任务执行日志仓库
func NewGormJobExecutionLogRepository(db *gorm.DB) *GormJobExecutionLogRepository {
	return &GormJobExecutionLogRepository{db: db}
}

// Create 创建执行日志
func (r *GormJobExecutionLogRepository) Create(log *models.JobExecutionLog) error {
	return r.db.Create(log).Error
}

// Update 更新执行日志
func (r *GormJobExecutionLogRepository) Update(log *models.JobExecutionLog) error {
	return r.db.Save(log).Error
}

// FindByID 根据ID查询日志
func (r *GormJobExecutionLogRepository) FindByID(id int64) (*models.JobExecutionLog, error) {
	var log models.JobExecutionLog
	err := r.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// FindByJobName 根据任务名称查询日志（分页）
func (r *GormJobExecutionLogRepository) FindByJobName(jobName string, limit, offset int) ([]*models.JobExecutionLog, error) {
	var logs []*models.JobExecutionLog
	query := r.db.Model(&models.JobExecutionLog{})
	if jobName != "" {
		query = query.Where("job_name = ?", jobName)
	}
	err := query.Order("started_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, err
}

// FindByJobNameAndDateRange 根据任务名称和日期范围查询日志
func (r *GormJobExecutionLogRepository) FindByJobNameAndDateRange(jobName string, startDate, endDate time.Time, limit, offset int) ([]*models.JobExecutionLog, error) {
	var logs []*models.JobExecutionLog
	query := r.db.Model(&models.JobExecutionLog{}).
		Where("started_at >= ? AND started_at < ?", startDate, endDate)
	if jobName != "" {
		query = query.Where("job_name = ?", jobName)
	}
	err := query.Order("started_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, err
}

// CountByJobName 统计任务日志数量
func (r *GormJobExecutionLogRepository) CountByJobName(jobName string) (int64, error) {
	var count int64
	query := r.db.Model(&models.JobExecutionLog{})
	if jobName != "" {
		query = query.Where("job_name = ?", jobName)
	}
	err := query.Count(&count).Error
	return count, err
}

// CountByJobNameAndDateRange 统计日期范围内任务日志数量
func (r *GormJobExecutionLogRepository) CountByJobNameAndDateRange(jobName string, startDate, endDate time.Time) (int64, error) {
	var count int64
	query := r.db.Model(&models.JobExecutionLog{}).
		Where("started_at >= ? AND started_at < ?", startDate, endDate)
	if jobName != "" {
		query = query.Where("job_name = ?", jobName)
	}
	err := query.Count(&count).Error
	return count, err
}

// FindLatestByJobName 查询任务最新的执行日志
func (r *GormJobExecutionLogRepository) FindLatestByJobName(jobName string) (*models.JobExecutionLog, error) {
	var log models.JobExecutionLog
	err := r.db.Where("job_name = ?", jobName).
		Order("started_at DESC").
		First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// FindRunningJobs 查询正在运行的任务
func (r *GormJobExecutionLogRepository) FindRunningJobs() ([]*models.JobExecutionLog, error) {
	var logs []*models.JobExecutionLog
	err := r.db.Where("status = ?", models.JobStatusRunning).
		Order("started_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetStats 获取任务执行统计
func (r *GormJobExecutionLogRepository) GetStats(startDate, endDate time.Time) ([]*JobExecutionStats, error) {
	var stats []*JobExecutionStats
	err := r.db.Model(&models.JobExecutionLog{}).
		Select(`
			job_name,
			COUNT(*) as total_count,
			SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) as success_count,
			SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) as fail_count,
			AVG(duration_ms) as avg_duration,
			MAX(duration_ms) as max_duration,
			MIN(duration_ms) as min_duration
		`).
		Where("started_at >= ? AND started_at < ?", startDate, endDate).
		Group("job_name").
		Order("job_name").
		Scan(&stats).Error
	return stats, err
}

// DeleteOlderThan 删除指定天数之前的日志
func (r *GormJobExecutionLogRepository) DeleteOlderThan(days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := r.db.Where("created_at < ?", cutoff).Delete(&models.JobExecutionLog{})
	return result.RowsAffected, result.Error
}

// Ensure interface compliance
var _ JobExecutionLogRepository = (*GormJobExecutionLogRepository)(nil)
