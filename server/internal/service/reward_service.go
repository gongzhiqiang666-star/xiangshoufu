package service

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"gorm.io/gorm"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// RewardService 奖励管理服务
// 处理奖励模版CRUD、代理商奖励比例、终端奖励进度跟踪、奖励计算与发放
type RewardService struct {
	db                   *gorm.DB
	templateRepo         *repository.GormRewardPolicyTemplateRepository
	stageRepo            *repository.GormRewardStageRepository
	agentRateRepo        *repository.GormAgentRewardRateRepository
	progressRepo         *repository.GormTerminalRewardProgressRepository
	stageRewardRepo      *repository.GormTerminalStageRewardRepository
	distributionRepo     *repository.GormRewardDistributionRepository
	overflowLogRepo      *repository.GormRewardOverflowLogRepository
	agentRepo            repository.AgentRepository
	transactionRepo      repository.TransactionRepository
	walletService        *WalletService
}

// NewRewardService 创建奖励服务
func NewRewardService(
	db *gorm.DB,
	templateRepo *repository.GormRewardPolicyTemplateRepository,
	stageRepo *repository.GormRewardStageRepository,
	agentRateRepo *repository.GormAgentRewardRateRepository,
	progressRepo *repository.GormTerminalRewardProgressRepository,
	stageRewardRepo *repository.GormTerminalStageRewardRepository,
	distributionRepo *repository.GormRewardDistributionRepository,
	overflowLogRepo *repository.GormRewardOverflowLogRepository,
	agentRepo repository.AgentRepository,
	transactionRepo repository.TransactionRepository,
	walletService *WalletService,
) *RewardService {
	return &RewardService{
		db:               db,
		templateRepo:     templateRepo,
		stageRepo:        stageRepo,
		agentRateRepo:    agentRateRepo,
		progressRepo:     progressRepo,
		stageRewardRepo:  stageRewardRepo,
		distributionRepo: distributionRepo,
		overflowLogRepo:  overflowLogRepo,
		agentRepo:        agentRepo,
		transactionRepo:  transactionRepo,
		walletService:    walletService,
	}
}

// ============================================================
// 奖励政策模版管理
// ============================================================

// CreateRewardTemplate 创建奖励政策模版
func (s *RewardService) CreateRewardTemplate(req *models.CreateRewardTemplateRequest) (*models.RewardPolicyTemplate, error) {
	// 1. 验证阶段配置
	if err := s.validateStages(req.TimeType, req.Stages); err != nil {
		return nil, err
	}

	// 2. 创建模版
	template := &models.RewardPolicyTemplate{
		Name:          req.Name,
		TimeType:      req.TimeType,
		DimensionType: req.DimensionType,
		TransTypes:    req.TransTypes,
		AmountMin:     req.AmountMin,
		AmountMax:     req.AmountMax,
		AllowGap:      req.AllowGap,
		Description:   req.Description,
		Enabled:       true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 使用事务
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 创建模版
		if err := tx.Create(template).Error; err != nil {
			return fmt.Errorf("创建模版失败: %w", err)
		}

		// 创建阶段配置
		stages := make([]*models.RewardStage, 0, len(req.Stages))
		for _, stageReq := range req.Stages {
			stage := &models.RewardStage{
				TemplateID:   template.ID,
				StageOrder:   stageReq.StageOrder,
				StartValue:   stageReq.StartValue,
				EndValue:     stageReq.EndValue,
				TargetValue:  stageReq.TargetValue,
				RewardAmount: stageReq.RewardAmount,
				CreatedAt:    time.Now(),
			}
			stages = append(stages, stage)
		}

		if len(stages) > 0 {
			if err := tx.Create(&stages).Error; err != nil {
				return fmt.Errorf("创建阶段配置失败: %w", err)
			}
		}

		template.Stages = stages
		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Printf("[RewardService] 创建奖励模版成功: ID=%d, Name=%s", template.ID, template.Name)
	return template, nil
}

// UpdateRewardTemplate 更新奖励政策模版
func (s *RewardService) UpdateRewardTemplate(id int64, req *models.UpdateRewardTemplateRequest) (*models.RewardPolicyTemplate, error) {
	// 1. 查询现有模版
	template, err := s.templateRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("模版不存在: %w", err)
	}

	// 2. 如果更新阶段，验证阶段配置
	if len(req.Stages) > 0 {
		if err := s.validateStages(template.TimeType, req.Stages); err != nil {
			return nil, err
		}
	}

	// 3. 使用事务更新
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 更新模版字段
		if req.Name != "" {
			template.Name = req.Name
		}
		if req.TransTypes != "" {
			template.TransTypes = req.TransTypes
		}
		template.AmountMin = req.AmountMin
		template.AmountMax = req.AmountMax
		if req.AllowGap != nil {
			template.AllowGap = *req.AllowGap
		}
		if req.Enabled != nil {
			template.Enabled = *req.Enabled
		}
		if req.Description != "" {
			template.Description = req.Description
		}
		template.UpdatedAt = time.Now()

		if err := tx.Save(template).Error; err != nil {
			return fmt.Errorf("更新模版失败: %w", err)
		}

		// 如果有阶段配置，先删除旧的再创建新的
		if len(req.Stages) > 0 {
			if err := tx.Where("template_id = ?", id).Delete(&models.RewardStage{}).Error; err != nil {
				return fmt.Errorf("删除旧阶段配置失败: %w", err)
			}

			stages := make([]*models.RewardStage, 0, len(req.Stages))
			for _, stageReq := range req.Stages {
				stage := &models.RewardStage{
					TemplateID:   template.ID,
					StageOrder:   stageReq.StageOrder,
					StartValue:   stageReq.StartValue,
					EndValue:     stageReq.EndValue,
					TargetValue:  stageReq.TargetValue,
					RewardAmount: stageReq.RewardAmount,
					CreatedAt:    time.Now(),
				}
				stages = append(stages, stage)
			}

			if err := tx.Create(&stages).Error; err != nil {
				return fmt.Errorf("创建新阶段配置失败: %w", err)
			}

			template.Stages = stages
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Printf("[RewardService] 更新奖励模版成功: ID=%d", id)
	return template, nil
}

// GetRewardTemplateDetail 获取奖励模版详情
func (s *RewardService) GetRewardTemplateDetail(id int64) (*models.RewardTemplateDetail, error) {
	template, err := s.templateRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("模版不存在: %w", err)
	}

	stages, err := s.stageRepo.FindByTemplateID(id)
	if err != nil {
		return nil, fmt.Errorf("查询阶段配置失败: %w", err)
	}

	return &models.RewardTemplateDetail{
		RewardPolicyTemplate: *template,
		Stages:               stages,
	}, nil
}

// GetRewardTemplateList 获取奖励模版列表
func (s *RewardService) GetRewardTemplateList(enabled *bool, page, pageSize int) ([]*models.RewardTemplateListItem, int64, error) {
	offset := (page - 1) * pageSize
	templates, total, err := s.templateRepo.FindAll(enabled, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// 转换为列表项
	items := make([]*models.RewardTemplateListItem, 0, len(templates))
	for _, t := range templates {
		// 查询阶段数量
		stages, _ := s.stageRepo.FindByTemplateID(t.ID)
		item := &models.RewardTemplateListItem{
			ID:            t.ID,
			Name:          t.Name,
			TimeType:      t.TimeType,
			DimensionType: t.DimensionType,
			TransTypes:    t.TransTypes,
			AllowGap:      t.AllowGap,
			Enabled:       t.Enabled,
			StageCount:    len(stages),
			CreatedAt:     t.CreatedAt,
			UpdatedAt:     t.UpdatedAt,
		}
		items = append(items, item)
	}

	return items, total, nil
}

// UpdateRewardTemplateEnabled 更新模版启用状态
func (s *RewardService) UpdateRewardTemplateEnabled(id int64, enabled bool) error {
	return s.templateRepo.UpdateEnabled(id, enabled)
}

// DeleteRewardTemplate 删除奖励模版
func (s *RewardService) DeleteRewardTemplate(id int64) error {
	// 检查是否有关联的进度记录
	// TODO: 添加检查逻辑

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 先删除阶段配置
		if err := tx.Where("template_id = ?", id).Delete(&models.RewardStage{}).Error; err != nil {
			return err
		}
		// 再删除模版
		return tx.Delete(&models.RewardPolicyTemplate{}, id).Error
	})
}

// validateStages 验证阶段配置
func (s *RewardService) validateStages(timeType models.TimeType, stages []models.CreateStageRequest) error {
	if len(stages) == 0 {
		return errors.New("至少需要一个阶段配置")
	}

	// 按阶段顺序排序
	sort.Slice(stages, func(i, j int) bool {
		return stages[i].StageOrder < stages[j].StageOrder
	})

	// 检查阶段顺序连续性和时间范围不重叠
	for i, stage := range stages {
		if stage.StageOrder != i+1 {
			return fmt.Errorf("阶段顺序必须从1开始连续递增，当前第%d个阶段顺序为%d", i+1, stage.StageOrder)
		}

		if stage.EndValue < stage.StartValue {
			return fmt.Errorf("阶段%d的结束值(%d)不能小于开始值(%d)", stage.StageOrder, stage.EndValue, stage.StartValue)
		}

		if stage.TargetValue <= 0 {
			return fmt.Errorf("阶段%d的目标值必须大于0", stage.StageOrder)
		}

		if stage.RewardAmount <= 0 {
			return fmt.Errorf("阶段%d的奖励金额必须大于0", stage.StageOrder)
		}

		// 检查与前一阶段是否重叠
		if i > 0 {
			prevStage := stages[i-1]
			if stage.StartValue <= prevStage.EndValue {
				return fmt.Errorf("阶段%d的开始值(%d)与阶段%d的结束值(%d)重叠",
					stage.StageOrder, stage.StartValue, prevStage.StageOrder, prevStage.EndValue)
			}
		}
	}

	return nil
}

// ============================================================
// 代理商奖励金额配置（差额分配模式）
// ============================================================

// SetAgentRewardAmount 设置代理商奖励金额（上级给下级配置的奖励金额）
// 差额分配规则：A给B配置100元，B给C配置30元，C达标时：C得30，B得70（100-30）
func (s *RewardService) SetAgentRewardAmount(req *models.AgentRewardRateRequest) error {
	// 1. 验证代理商存在
	agent, err := s.agentRepo.FindByID(req.AgentID)
	if err != nil {
		return fmt.Errorf("代理商不存在: %w", err)
	}

	// 2. 验证模版存在
	_, err = s.templateRepo.FindByID(req.TemplateID)
	if err != nil {
		return fmt.Errorf("奖励模版不存在: %w", err)
	}

	// 3. 验证金额范围（不能为负数）
	if req.RewardAmount < 0 {
		return errors.New("奖励金额不能为负数")
	}

	// 4. 验证不超过上级配置的金额（配置时验证）
	if err := s.validateAgentRewardAmount(agent, req.TemplateID, req.RewardAmount); err != nil {
		return err
	}

	// 5. 保存或更新
	rate := &models.AgentRewardRate{
		AgentID:      req.AgentID,
		TemplateID:   req.TemplateID,
		RewardAmount: req.RewardAmount,
		UpdatedAt:    time.Now(),
	}

	if err := s.agentRateRepo.Upsert(rate); err != nil {
		return fmt.Errorf("保存奖励金额配置失败: %w", err)
	}

	log.Printf("[RewardService] 设置代理商奖励金额: AgentID=%d, TemplateID=%d, Amount=%d分",
		req.AgentID, req.TemplateID, req.RewardAmount)
	return nil
}

// GetAgentRewardAmount 获取代理商奖励金额配置
func (s *RewardService) GetAgentRewardAmount(agentID, templateID int64) (*models.AgentRewardRate, error) {
	return s.agentRateRepo.FindByAgentAndTemplate(agentID, templateID)
}

// validateAgentRewardAmount 验证代理商奖励金额不超过上级配置的金额
// 差额分配规则：下级配置的金额不能超过上级给自己配置的金额
func (s *RewardService) validateAgentRewardAmount(agent *repository.Agent, templateID int64, newAmount int64) error {
	// 如果没有上级（顶级代理商），则不需要验证
	if agent.ParentID == 0 {
		return nil
	}

	// 获取上级给当前代理商配置的金额
	parentRate, err := s.agentRateRepo.FindByAgentAndTemplate(agent.ID, templateID)
	if err != nil {
		// 如果上级还没有给当前代理商配置金额，则当前代理商也不能配置
		return errors.New("上级尚未配置奖励金额，无法给下级配置")
	}

	// 下级配置的金额不能超过上级给自己配置的金额
	if newAmount > parentRate.RewardAmount {
		return fmt.Errorf("配置的奖励金额(%d分)不能超过上级配置的金额(%d分)",
			newAmount, parentRate.RewardAmount)
	}

	return nil
}

// ============================================================
// 终端奖励进度管理
// ============================================================

// InitTerminalRewardProgress 初始化终端奖励进度（终端绑定时调用）
func (s *RewardService) InitTerminalRewardProgress(terminalSN string, terminalID *int64, agentID int64, templateID int64) (*models.TerminalRewardProgress, error) {
	// 1. 检查是否已有进行中的进度
	existing, err := s.progressRepo.FindActiveByTerminalSN(terminalSN)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("终端%s已有进行中的奖励进度", terminalSN)
	}

	// 2. 获取模版详情
	template, err := s.templateRepo.FindByID(templateID)
	if err != nil {
		return nil, fmt.Errorf("奖励模版不存在: %w", err)
	}

	if !template.Enabled {
		return nil, errors.New("奖励模版已禁用")
	}

	stages, err := s.stageRepo.FindByTemplateID(templateID)
	if err != nil {
		return nil, fmt.Errorf("获取阶段配置失败: %w", err)
	}

	// 3. 创建模版快照
	snapshot := models.TemplateSnapshot{
		ID:            template.ID,
		Name:          template.Name,
		TimeType:      template.TimeType,
		DimensionType: template.DimensionType,
		TransTypes:    template.TransTypes,
		AmountMin:     template.AmountMin,
		AmountMax:     template.AmountMax,
		AllowGap:      template.AllowGap,
		Stages:        stages,
	}

	bindTime := time.Now()

	// 4. 创建进度记录
	progress := &models.TerminalRewardProgress{
		TerminalSN:        terminalSN,
		TerminalID:        terminalID,
		TemplateID:        templateID,
		TemplateSnapshot:  snapshot,
		BindAgentID:       agentID,
		BindTime:          bindTime,
		CurrentStage:      1,
		LastAchievedStage: 0,
		Status:            models.RewardProgressStatusActive,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// 5. 使用事务创建进度和阶段奖励记录
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(progress).Error; err != nil {
			return fmt.Errorf("创建进度记录失败: %w", err)
		}

		// 创建各阶段的奖励记录
		stageRewards := make([]*models.TerminalStageReward, 0, len(stages))
		for _, stage := range stages {
			stageStart, stageEnd := s.calculateStageTime(bindTime, template.TimeType, stage.StartValue, stage.EndValue)

			reward := &models.TerminalStageReward{
				ProgressID:   progress.ID,
				TerminalSN:   terminalSN,
				StageOrder:   stage.StageOrder,
				StageStart:   stageStart,
				StageEnd:     stageEnd,
				TargetValue:  stage.TargetValue,
				ActualValue:  0,
				IsAchieved:   false,
				RewardAmount: &stage.RewardAmount,
				Status:       models.StageRewardStatusPending,
				GapBlocked:   false,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			stageRewards = append(stageRewards, reward)
		}

		if err := tx.Create(&stageRewards).Error; err != nil {
			return fmt.Errorf("创建阶段奖励记录失败: %w", err)
		}

		progress.StageRewards = stageRewards
		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Printf("[RewardService] 初始化终端奖励进度: TerminalSN=%s, TemplateID=%d, AgentID=%d", terminalSN, templateID, agentID)
	return progress, nil
}

// calculateStageTime 计算阶段时间范围
func (s *RewardService) calculateStageTime(bindTime time.Time, timeType models.TimeType, startValue, endValue int) (time.Time, time.Time) {
	// 绑定当天算第1天，从绑定日0:00:00开始
	bindDate := time.Date(bindTime.Year(), bindTime.Month(), bindTime.Day(), 0, 0, 0, 0, bindTime.Location())

	var stageStart, stageEnd time.Time

	if timeType == models.TimeTypeDays {
		// 按天数：第N天 = 绑定日 + (N-1)天
		stageStart = bindDate.AddDate(0, 0, startValue-1)
		stageEnd = bindDate.AddDate(0, 0, endValue-1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		// 按自然月：第N月 = 绑定日 + (N-1)月
		stageStart = bindDate.AddDate(0, startValue-1, 0)
		// 结束时间为该月最后一天23:59:59
		endMonth := bindDate.AddDate(0, endValue, 0)
		stageEnd = endMonth.AddDate(0, 0, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	return stageStart, stageEnd
}

// TerminateTerminalRewardProgress 终止终端奖励进度（终端解绑时调用）
func (s *RewardService) TerminateTerminalRewardProgress(terminalSN string) error {
	progress, err := s.progressRepo.FindActiveByTerminalSN(terminalSN)
	if err != nil {
		// 没有进行中的进度，直接返回
		return nil
	}

	// 中途解绑不给奖励
	if err := s.progressRepo.Terminate(progress.ID); err != nil {
		return fmt.Errorf("终止进度失败: %w", err)
	}

	log.Printf("[RewardService] 终止终端奖励进度: TerminalSN=%s, ProgressID=%d", terminalSN, progress.ID)
	return nil
}

// GetTerminalRewardProgress 获取终端奖励进度详情
func (s *RewardService) GetTerminalRewardProgress(terminalSN string) (*models.TerminalRewardProgressDetail, error) {
	progress, err := s.progressRepo.FindActiveByTerminalSN(terminalSN)
	if err != nil {
		return nil, fmt.Errorf("未找到进行中的奖励进度: %w", err)
	}

	stageRewards, err := s.stageRewardRepo.FindByProgressID(progress.ID)
	if err != nil {
		return nil, fmt.Errorf("获取阶段奖励记录失败: %w", err)
	}

	// 获取代理商名称
	agent, _ := s.agentRepo.FindByID(progress.BindAgentID)
	agentName := ""
	if agent != nil {
		agentName = agent.AgentName
	}

	return &models.TerminalRewardProgressDetail{
		TerminalRewardProgress: *progress,
		TemplateName:           progress.TemplateSnapshot.Name,
		AgentName:              agentName,
		StageRewards:           stageRewards,
	}, nil
}

// ============================================================
// 奖励计算与发放（定时任务调用）
// ============================================================

// ProcessPendingStageRewards 处理待检查的阶段奖励（定时任务入口）
func (s *RewardService) ProcessPendingStageRewards(batchSize int) (int, int, error) {
	now := time.Now()
	successCount := 0
	failCount := 0

	// 分批处理
	offset := 0
	for {
		rewards, err := s.stageRewardRepo.FindPendingByStageEnd(now, batchSize)
		if err != nil {
			return successCount, failCount, fmt.Errorf("查询待处理奖励失败: %w", err)
		}

		if len(rewards) == 0 {
			break
		}

		for _, reward := range rewards {
			if err := s.processStageReward(reward); err != nil {
				log.Printf("[RewardService] 处理阶段奖励失败: ID=%d, Error=%v", reward.ID, err)
				failCount++
				// 重试3次后标记失败
				continue
			}
			successCount++
		}

		offset += batchSize
		if len(rewards) < batchSize {
			break
		}
	}

	log.Printf("[RewardService] 处理阶段奖励完成: 成功=%d, 失败=%d", successCount, failCount)
	return successCount, failCount, nil
}

// processStageReward 处理单个阶段奖励
func (s *RewardService) processStageReward(reward *models.TerminalStageReward) error {
	// 1. 获取进度记录
	progress, err := s.progressRepo.FindByID(reward.ProgressID)
	if err != nil {
		return fmt.Errorf("获取进度记录失败: %w", err)
	}

	// 检查进度状态
	if progress.Status != models.RewardProgressStatusActive {
		// 进度已终止，标记阶段失败
		return s.stageRewardRepo.UpdateStatus(reward.ID, models.StageRewardStatusFailed)
	}

	snapshot := progress.TemplateSnapshot

	// 2. 检查断档阻断
	if !snapshot.AllowGap && reward.StageOrder > 1 {
		// 不允许断档，检查上一阶段是否达标
		prevRewards, _ := s.stageRewardRepo.FindByProgressID(progress.ID)
		for _, pr := range prevRewards {
			if pr.StageOrder == reward.StageOrder-1 && !pr.IsAchieved {
				// 上一阶段未达标，阻断当前阶段
				return s.db.Transaction(func(tx *gorm.DB) error {
					if err := tx.Model(&models.TerminalStageReward{}).
						Where("id = ?", reward.ID).
						Updates(map[string]interface{}{
							"status":      models.StageRewardStatusGapBlocked,
							"gap_blocked": true,
							"updated_at":  time.Now(),
						}).Error; err != nil {
						return err
					}
					log.Printf("[RewardService] 阶段奖励被断档阻断: ID=%d, Stage=%d", reward.ID, reward.StageOrder)
					return nil
				})
			}
		}
	}

	// 3. 计算实际值（查询交易统计）
	actualValue, err := s.calculateActualValue(reward.TerminalSN, snapshot, reward.StageStart, reward.StageEnd)
	if err != nil {
		return fmt.Errorf("计算实际值失败: %w", err)
	}

	// 4. 判断是否达标
	isAchieved := actualValue >= reward.TargetValue

	// 5. 更新阶段奖励记录
	return s.db.Transaction(func(tx *gorm.DB) error {
		var newStatus models.StageRewardStatus
		if isAchieved {
			newStatus = models.StageRewardStatusAchieved
		} else {
			newStatus = models.StageRewardStatusFailed
		}

		if err := tx.Model(&models.TerminalStageReward{}).
			Where("id = ?", reward.ID).
			Updates(map[string]interface{}{
				"actual_value": actualValue,
				"is_achieved":  isAchieved,
				"status":       newStatus,
				"updated_at":   time.Now(),
			}).Error; err != nil {
			return err
		}

		// 如果达标，进行级差分配
		if isAchieved && reward.RewardAmount != nil && *reward.RewardAmount > 0 {
			if err := s.distributeReward(tx, reward, progress); err != nil {
				return fmt.Errorf("奖励分配失败: %w", err)
			}
		}

		// 更新进度记录
		if isAchieved {
			if err := tx.Model(&models.TerminalRewardProgress{}).
				Where("id = ?", progress.ID).
				Updates(map[string]interface{}{
					"last_achieved_stage": reward.StageOrder,
					"updated_at":          time.Now(),
				}).Error; err != nil {
				return err
			}
		}

		log.Printf("[RewardService] 阶段奖励处理完成: ID=%d, Stage=%d, Achieved=%v, Actual=%d, Target=%d",
			reward.ID, reward.StageOrder, isAchieved, actualValue, reward.TargetValue)
		return nil
	})
}

// calculateActualValue 计算实际值（查询交易统计）
func (s *RewardService) calculateActualValue(terminalSN string, snapshot models.TemplateSnapshot, start, end time.Time) (int64, error) {
	// TODO: 根据snapshot的配置查询交易记录并统计
	// - 按dimension_type决定是统计金额还是笔数
	// - 按trans_types筛选交易类型
	// - 按amount_min/amount_max筛选交易金额

	// 临时实现：直接查询终端总交易金额
	if snapshot.DimensionType == models.DimensionTypeAmount {
		return s.transactionRepo.GetTerminalTotalTradeAmount(terminalSN)
	}

	// 按笔数统计 - TODO: 实现
	return 0, nil
}

// distributeReward 差额分配奖励
// 规则：A给B配置100元，B给C配置30元，C达标时：C得30，B得70（100-30），总金额100
func (s *RewardService) distributeReward(tx *gorm.DB, reward *models.TerminalStageReward, progress *models.TerminalRewardProgress) error {
	if reward.RewardAmount == nil || *reward.RewardAmount <= 0 {
		return nil
	}

	templateID := progress.TemplateID

	// 1. 获取代理商链（从终端归属代理商往上）
	ancestors, err := s.agentRepo.FindAncestors(progress.BindAgentID)
	if err != nil {
		return fmt.Errorf("获取代理商链失败: %w", err)
	}

	// 2. 构建完整的代理商链（包括终端归属代理商）
	// 顺序：终端归属代理商 -> 上级 -> 上上级 -> ... -> 顶级
	allAgentIDs := make([]int64, 0, len(ancestors)+1)
	allAgentIDs = append(allAgentIDs, progress.BindAgentID)
	for _, a := range ancestors {
		allAgentIDs = append(allAgentIDs, a.ID)
	}

	// 3. 获取各级代理商针对该模版的奖励金额配置
	rates, err := s.agentRateRepo.FindByAgentIDsAndTemplate(allAgentIDs, templateID)
	if err != nil {
		return fmt.Errorf("获取奖励金额配置失败: %w", err)
	}

	// 构建金额映射：agent_id -> 上级给该代理商配置的奖励金额
	amountMap := make(map[int64]int64)
	for _, r := range rates {
		amountMap[r.AgentID] = r.RewardAmount
	}

	// 4. 差额分配计算
	// 规则：每个代理商得到的金额 = 上级给自己配置的金额 - 自己给下级配置的金额
	distributions := make([]*models.RewardDistribution, 0)
	agentChain := make(models.AgentChain, 0)

	// 从终端归属代理商开始，往上遍历
	var totalDistributed int64
	prevAmount := int64(0) // 下级获得的金额（用于计算差额）

	for level, agentID := range allAgentIDs {
		// 获取上级给当前代理商配置的金额
		configuredAmount, hasConfig := amountMap[agentID]
		if !hasConfig {
			// 如果没有配置，跳过（顶级代理商可能没有上级给他配置）
			continue
		}

		// 当前代理商获得的金额 = 上级给自己配置的金额 - 自己给下级配置的金额
		earnedAmount := configuredAmount - prevAmount

		if earnedAmount < 0 {
			// 配置错误：下级金额超过上级配置，记录异常
			agentChain = append(agentChain, models.AgentChainInfo{
				AgentID:      agentID,
				Level:        level + 1,
				RewardRate:   float64(configuredAmount) / float64(*reward.RewardAmount),
			})

			overflowLog := &models.RewardOverflowLog{
				TerminalSN:    reward.TerminalSN,
				StageRewardID: &reward.ID,
				AgentChain:    agentChain,
				TotalRate:     0,
				RewardAmount:  *reward.RewardAmount,
				ErrorMessage:  fmt.Sprintf("差额分配异常：代理商%d配置金额(%d)小于下级获得金额(%d)", agentID, configuredAmount, prevAmount),
				CreatedAt:     time.Now(),
			}
			if err := tx.Create(overflowLog).Error; err != nil {
				log.Printf("[RewardService] 记录溢出日志失败: %v", err)
			}
			return fmt.Errorf("差额分配异常：配置金额不足以覆盖下级奖励")
		}

		if earnedAmount > 0 {
			dist := &models.RewardDistribution{
				StageRewardID: reward.ID,
				TerminalSN:    reward.TerminalSN,
				AgentID:       agentID,
				AgentLevel:    level + 1,
				RewardRate:    float64(earnedAmount) / float64(*reward.RewardAmount),
				RewardAmount:  earnedAmount,
				WalletStatus:  0,
				CreatedAt:     time.Now(),
			}
			distributions = append(distributions, dist)
			totalDistributed += earnedAmount
		}

		// 更新prevAmount为当前代理商配置的金额
		prevAmount = configuredAmount
	}

	// 5. 验证总分配金额不超过奖励金额
	if totalDistributed > *reward.RewardAmount {
		overflowLog := &models.RewardOverflowLog{
			TerminalSN:    reward.TerminalSN,
			StageRewardID: &reward.ID,
			AgentChain:    agentChain,
			TotalRate:     float64(totalDistributed) / float64(*reward.RewardAmount),
			RewardAmount:  *reward.RewardAmount,
			ErrorMessage:  fmt.Sprintf("分配总额(%d)超过奖励金额(%d)", totalDistributed, *reward.RewardAmount),
			CreatedAt:     time.Now(),
		}
		if err := tx.Create(overflowLog).Error; err != nil {
			log.Printf("[RewardService] 记录溢出日志失败: %v", err)
		}
		return fmt.Errorf("奖励池溢出: 分配总额(%d分)超过奖励金额(%d分)", totalDistributed, *reward.RewardAmount)
	}

	// 6. 批量创建分配记录
	if len(distributions) > 0 {
		if err := tx.Create(&distributions).Error; err != nil {
			return fmt.Errorf("创建分配记录失败: %w", err)
		}
	}

	log.Printf("[RewardService] 差额分配完成: StageRewardID=%d, Total=%d, Distributed=%d, Records=%d",
		reward.ID, *reward.RewardAmount, totalDistributed, len(distributions))
	return nil
}

// ============================================================
// 奖励入账（将分配记录入账到钱包）
// ============================================================

// SettleRewardDistributions 结算奖励分配（入账到奖励钱包）
func (s *RewardService) SettleRewardDistributions(batchSize int) (int, int, error) {
	successCount := 0
	failCount := 0

	distributions, err := s.distributionRepo.FindPendingWallet(batchSize)
	if err != nil {
		return 0, 0, fmt.Errorf("查询待入账分配失败: %w", err)
	}

	for _, dist := range distributions {
		// 入账到奖励钱包（wallet_type=3）
		if s.walletService != nil {
			walletRecordID, err := s.walletService.AddRewardWalletBalance(dist.AgentID, dist.RewardAmount, fmt.Sprintf("阶段奖励-%s", dist.TerminalSN))
			if err != nil {
				log.Printf("[RewardService] 奖励入账失败: DistID=%d, Error=%v", dist.ID, err)
				failCount++
				continue
			}

			if err := s.distributionRepo.UpdateWalletStatus(dist.ID, walletRecordID, 1); err != nil {
				log.Printf("[RewardService] 更新入账状态失败: DistID=%d, Error=%v", dist.ID, err)
				failCount++
				continue
			}
		}

		successCount++
	}

	if successCount > 0 || failCount > 0 {
		log.Printf("[RewardService] 奖励入账完成: 成功=%d, 失败=%d", successCount, failCount)
	}
	return successCount, failCount, nil
}

// ============================================================
// 溢出日志管理
// ============================================================

// GetUnresolvedOverflowLogs 获取未解决的溢出日志
func (s *RewardService) GetUnresolvedOverflowLogs(page, pageSize int) ([]*models.RewardOverflowLog, int64, error) {
	offset := (page - 1) * pageSize
	return s.overflowLogRepo.FindUnresolved(pageSize, offset)
}

// ResolveOverflowLog 解决溢出日志
func (s *RewardService) ResolveOverflowLog(id int64, resolvedBy string) error {
	return s.overflowLogRepo.Resolve(id, resolvedBy)
}
