package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"xiangshoufu/internal/models"
)

// TestValidateStages 测试阶段配置验证
func TestValidateStages(t *testing.T) {
	service := &RewardService{}

	tests := []struct {
		name      string
		timeType  models.TimeType
		stages    []models.CreateStageRequest
		wantError bool
		errorMsg  string
	}{
		{
			name:     "正常的阶段配置",
			timeType: models.TimeTypeDays,
			stages: []models.CreateStageRequest{
				{StageOrder: 1, StartValue: 1, EndValue: 10, TargetValue: 10000, RewardAmount: 5000},
				{StageOrder: 2, StartValue: 11, EndValue: 20, TargetValue: 20000, RewardAmount: 10000},
				{StageOrder: 3, StartValue: 21, EndValue: 30, TargetValue: 30000, RewardAmount: 15000},
			},
			wantError: false,
		},
		{
			name:      "空阶段配置",
			timeType:  models.TimeTypeDays,
			stages:    []models.CreateStageRequest{},
			wantError: true,
			errorMsg:  "至少需要一个阶段配置",
		},
		{
			name:     "阶段顺序不连续",
			timeType: models.TimeTypeDays,
			stages: []models.CreateStageRequest{
				{StageOrder: 1, StartValue: 1, EndValue: 10, TargetValue: 10000, RewardAmount: 5000},
				{StageOrder: 3, StartValue: 11, EndValue: 20, TargetValue: 20000, RewardAmount: 10000},
			},
			wantError: true,
			errorMsg:  "阶段顺序必须从1开始连续递增",
		},
		{
			name:     "阶段时间重叠",
			timeType: models.TimeTypeDays,
			stages: []models.CreateStageRequest{
				{StageOrder: 1, StartValue: 1, EndValue: 10, TargetValue: 10000, RewardAmount: 5000},
				{StageOrder: 2, StartValue: 8, EndValue: 20, TargetValue: 20000, RewardAmount: 10000},
			},
			wantError: true,
			errorMsg:  "重叠",
		},
		{
			name:     "结束值小于开始值",
			timeType: models.TimeTypeDays,
			stages: []models.CreateStageRequest{
				{StageOrder: 1, StartValue: 10, EndValue: 5, TargetValue: 10000, RewardAmount: 5000},
			},
			wantError: true,
			errorMsg:  "结束值",
		},
		{
			name:     "目标值为0",
			timeType: models.TimeTypeDays,
			stages: []models.CreateStageRequest{
				{StageOrder: 1, StartValue: 1, EndValue: 10, TargetValue: 0, RewardAmount: 5000},
			},
			wantError: true,
			errorMsg:  "目标值必须大于0",
		},
		{
			name:     "奖励金额为0",
			timeType: models.TimeTypeDays,
			stages: []models.CreateStageRequest{
				{StageOrder: 1, StartValue: 1, EndValue: 10, TargetValue: 10000, RewardAmount: 0},
			},
			wantError: true,
			errorMsg:  "奖励金额必须大于0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateStages(tt.timeType, tt.stages)
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCalculateStageTime 测试阶段时间计算
func TestCalculateStageTime(t *testing.T) {
	service := &RewardService{}

	// 固定绑定时间：2026-01-15 10:30:00
	bindTime := time.Date(2026, 1, 15, 10, 30, 0, 0, time.Local)

	tests := []struct {
		name           string
		timeType       models.TimeType
		startValue     int
		endValue       int
		expectedStart  time.Time
		expectedEndDay int // 只验证结束日期的天
	}{
		{
			name:           "按天数-第1-10天",
			timeType:       models.TimeTypeDays,
			startValue:     1,
			endValue:       10,
			expectedStart:  time.Date(2026, 1, 15, 0, 0, 0, 0, time.Local), // 绑定当天算第1天
			expectedEndDay: 24, // 1月15日 + 9天 = 1月24日
		},
		{
			name:           "按天数-第11-20天",
			timeType:       models.TimeTypeDays,
			startValue:     11,
			endValue:       20,
			expectedStart:  time.Date(2026, 1, 25, 0, 0, 0, 0, time.Local),
			expectedEndDay: 3, // 2月3日
		},
		{
			name:           "按自然月-第1月",
			timeType:       models.TimeTypeMonths,
			startValue:     1,
			endValue:       1,
			expectedStart:  time.Date(2026, 1, 15, 0, 0, 0, 0, time.Local),
			expectedEndDay: 14, // 2月14日（第1个月的最后一天）
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stageStart, stageEnd := service.calculateStageTime(bindTime, tt.timeType, tt.startValue, tt.endValue)

			assert.Equal(t, tt.expectedStart, stageStart, "阶段开始时间不匹配")
			assert.Equal(t, tt.expectedEndDay, stageEnd.Day(), "阶段结束日期不匹配")
			assert.Equal(t, 23, stageEnd.Hour(), "阶段结束时间应为23点")
			assert.Equal(t, 59, stageEnd.Minute(), "阶段结束分钟应为59")
			assert.Equal(t, 59, stageEnd.Second(), "阶段结束秒应为59")
		})
	}
}

// TestTemplateSnapshot 测试模版快照序列化
func TestTemplateSnapshot(t *testing.T) {
	snapshot := models.TemplateSnapshot{
		ID:            1,
		Name:          "测试模版",
		TimeType:      models.TimeTypeDays,
		DimensionType: models.DimensionTypeAmount,
		TransTypes:    "scan,debit,credit",
		AmountMin:     0,
		AmountMax:     nil,
		AllowGap:      false,
		Stages: []*models.RewardStage{
			{ID: 1, TemplateID: 1, StageOrder: 1, StartValue: 1, EndValue: 10, TargetValue: 10000, RewardAmount: 5000},
		},
	}

	// 测试Value()
	value, err := snapshot.Value()
	assert.NoError(t, err)
	assert.NotNil(t, value)

	// 测试Scan()
	var scanned models.TemplateSnapshot
	err = scanned.Scan(value)
	assert.NoError(t, err)
	assert.Equal(t, snapshot.ID, scanned.ID)
	assert.Equal(t, snapshot.Name, scanned.Name)
	assert.Equal(t, snapshot.TimeType, scanned.TimeType)
	assert.Equal(t, len(snapshot.Stages), len(scanned.Stages))
}

// TestAgentChain 测试代理商链序列化
func TestAgentChain(t *testing.T) {
	chain := models.AgentChain{
		{AgentID: 1, AgentName: "顶级代理商", Level: 1, RewardRate: 0.05},
		{AgentID: 2, AgentName: "中级代理商", Level: 2, RewardRate: 0.10},
		{AgentID: 3, AgentName: "终端代理商", Level: 3, RewardRate: 0.00},
	}

	// 测试Value()
	value, err := chain.Value()
	assert.NoError(t, err)
	assert.NotNil(t, value)

	// 测试Scan()
	var scanned models.AgentChain
	err = scanned.Scan(value)
	assert.NoError(t, err)
	assert.Equal(t, len(chain), len(scanned))
	assert.Equal(t, chain[0].AgentID, scanned[0].AgentID)
	assert.Equal(t, chain[1].RewardRate, scanned[1].RewardRate)
}

// TestRewardRateValidation 测试奖励比例验证
func TestRewardRateValidation(t *testing.T) {
	tests := []struct {
		name      string
		rate      float64
		wantValid bool
	}{
		{"0%", 0.0, true},
		{"10%", 0.10, true},
		{"50%", 0.50, true},
		{"100%", 1.0, true},
		{"负数", -0.1, false},
		{"超过100%", 1.1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.rate >= 0 && tt.rate <= 1
			assert.Equal(t, tt.wantValid, isValid)
		})
	}
}

// TestFixedPoolDistribution 测试固定池分配计算
func TestFixedPoolDistribution(t *testing.T) {
	// 模拟固定池分配
	// 总奖励：100元
	// A(顶级,5%) → B(中级,10%) → C(终端归属)
	totalReward := int64(10000) // 100元 = 10000分

	rates := map[string]float64{
		"A": 0.05, // 5%
		"B": 0.10, // 10%
		"C": 0.00, // 终端归属拿剩余
	}

	// 计算分配
	aAmount := int64(float64(totalReward) * rates["A"]) // 500分 = 5元
	bAmount := int64(float64(totalReward) * rates["B"]) // 1000分 = 10元
	cAmount := totalReward - aAmount - bAmount          // 8500分 = 85元

	assert.Equal(t, int64(500), aAmount, "A应得5元")
	assert.Equal(t, int64(1000), bAmount, "B应得10元")
	assert.Equal(t, int64(8500), cAmount, "C应得85元")
	assert.Equal(t, totalReward, aAmount+bAmount+cAmount, "总和应等于100元")
}

// TestOverflowDetection 测试溢出检测
func TestOverflowDetection(t *testing.T) {
	tests := []struct {
		name       string
		rates      []float64
		wantError  bool
	}{
		{
			name:      "正常比例之和15%",
			rates:     []float64{0.05, 0.10},
			wantError: false,
		},
		{
			name:      "比例之和刚好100%",
			rates:     []float64{0.30, 0.30, 0.40},
			wantError: false,
		},
		{
			name:      "比例之和超过100%",
			rates:     []float64{0.50, 0.60},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var totalRate float64
			for _, r := range tt.rates {
				totalRate += r
			}
			hasOverflow := totalRate > 1.0
			assert.Equal(t, tt.wantError, hasOverflow)
		})
	}
}
