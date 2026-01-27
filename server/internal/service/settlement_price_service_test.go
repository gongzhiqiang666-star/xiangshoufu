package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"xiangshoufu/internal/models"
)

// MockSettlementPriceRepository 模拟结算价仓库
type MockSettlementPriceRepository struct {
	mock.Mock
}

func (m *MockSettlementPriceRepository) Create(price *models.SettlementPrice) error {
	args := m.Called(price)
	return args.Error(0)
}

func (m *MockSettlementPriceRepository) Update(price *models.SettlementPrice) error {
	args := m.Called(price)
	return args.Error(0)
}

func (m *MockSettlementPriceRepository) GetByID(id int64) (*models.SettlementPrice, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SettlementPrice), args.Error(1)
}

func (m *MockSettlementPriceRepository) GetByAgentAndChannel(agentID, channelID int64, brandCode string) (*models.SettlementPrice, error) {
	args := m.Called(agentID, channelID, brandCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SettlementPrice), args.Error(1)
}

func (m *MockSettlementPriceRepository) List(req *models.SettlementPriceListRequest) ([]models.SettlementPrice, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]models.SettlementPrice), args.Get(1).(int64), args.Error(2)
}

func (m *MockSettlementPriceRepository) ListByAgent(agentID int64) ([]models.SettlementPrice, error) {
	args := m.Called(agentID)
	return args.Get(0).([]models.SettlementPrice), args.Error(1)
}

func (m *MockSettlementPriceRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockPriceChangeLogRepository 模拟调价记录仓库
type MockPriceChangeLogRepository struct {
	mock.Mock
}

func (m *MockPriceChangeLogRepository) Create(log *models.PriceChangeLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockPriceChangeLogRepository) GetByID(id int64) (*models.PriceChangeLog, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceChangeLog), args.Error(1)
}

func (m *MockPriceChangeLogRepository) List(req *models.PriceChangeLogListRequest) ([]models.PriceChangeLog, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]models.PriceChangeLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockPriceChangeLogRepository) ListByAgent(agentID int64, page, pageSize int) ([]models.PriceChangeLog, int64, error) {
	args := m.Called(agentID, page, pageSize)
	return args.Get(0).([]models.PriceChangeLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockPriceChangeLogRepository) ListBySettlementPrice(settlementPriceID int64, page, pageSize int) ([]models.PriceChangeLog, int64, error) {
	args := m.Called(settlementPriceID, page, pageSize)
	return args.Get(0).([]models.PriceChangeLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockPriceChangeLogRepository) ListByRewardSetting(rewardSettingID int64, page, pageSize int) ([]models.PriceChangeLog, int64, error) {
	args := m.Called(rewardSettingID, page, pageSize)
	return args.Get(0).([]models.PriceChangeLog), args.Get(1).(int64), args.Error(2)
}

// ============================================================
// SettlementPriceService 测试用例
// ============================================================

func TestSettlementPriceService_CreateFromTemplate(t *testing.T) {
	t.Run("should create settlement price from template", func(t *testing.T) {
		mockRepo := new(MockSettlementPriceRepository)
		mockLogRepo := new(MockPriceChangeLogRepository)

		service := NewSettlementPriceService(mockRepo, mockLogRepo, nil)

		// 模拟不存在的结算价
		mockRepo.On("GetByAgentAndChannel", int64(1), int64(1), "").Return(nil, assert.AnError)
		mockRepo.On("Create", mock.AnythingOfType("*models.SettlementPrice")).Return(nil)
		mockLogRepo.On("Create", mock.AnythingOfType("*models.PriceChangeLog")).Return(nil)

		price, err := service.CreateFromTemplate(1, 1, nil, "", nil, 1, "admin", "PC")

		assert.NoError(t, err)
		assert.NotNil(t, price)
		assert.Equal(t, int64(1), price.AgentID)
		assert.Equal(t, int64(1), price.ChannelID)
		assert.Equal(t, 1, price.Version)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when settlement price already exists", func(t *testing.T) {
		mockRepo := new(MockSettlementPriceRepository)
		mockLogRepo := new(MockPriceChangeLogRepository)

		service := NewSettlementPriceService(mockRepo, mockLogRepo, nil)

		existingPrice := &models.SettlementPrice{ID: 1, AgentID: 1, ChannelID: 1}
		mockRepo.On("GetByAgentAndChannel", int64(1), int64(1), "").Return(existingPrice, nil)

		_, err := service.CreateFromTemplate(1, 1, nil, "", nil, 1, "admin", "PC")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "结算价已存在")
	})
}

func TestSettlementPriceService_UpdateRate(t *testing.T) {
	t.Run("should update rate and create change log", func(t *testing.T) {
		mockRepo := new(MockSettlementPriceRepository)
		mockLogRepo := new(MockPriceChangeLogRepository)

		service := NewSettlementPriceService(mockRepo, mockLogRepo, nil)

		existingPrice := &models.SettlementPrice{
			ID:        1,
			AgentID:   1,
			ChannelID: 1,
			Version:   1,
			RateConfigs: models.RateConfigs{
				"credit": {Rate: "0.55"},
			},
		}

		mockRepo.On("GetByID", int64(1)).Return(existingPrice, nil)
		mockRepo.On("Update", mock.AnythingOfType("*models.SettlementPrice")).Return(nil)
		mockLogRepo.On("Create", mock.AnythingOfType("*models.PriceChangeLog")).Return(nil)

		newRate := "0.60"
		req := &models.UpdateRateRequest{
			CreditRate: &newRate,
		}

		price, err := service.UpdateRate(1, req, 1, "admin", "PC", "127.0.0.1")

		assert.NoError(t, err)
		assert.NotNil(t, price)
		assert.Equal(t, 2, price.Version) // 版本号应该+1
		assert.Equal(t, "0.60", *price.CreditRate)
		mockRepo.AssertExpectations(t)
		mockLogRepo.AssertExpectations(t)
	})

	t.Run("should return error when settlement price not found", func(t *testing.T) {
		mockRepo := new(MockSettlementPriceRepository)
		mockLogRepo := new(MockPriceChangeLogRepository)

		service := NewSettlementPriceService(mockRepo, mockLogRepo, nil)

		mockRepo.On("GetByID", int64(999)).Return(nil, assert.AnError)

		req := &models.UpdateRateRequest{}
		_, err := service.UpdateRate(999, req, 1, "admin", "PC", "127.0.0.1")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "结算价不存在")
	})
}

func TestSettlementPriceService_GetAgentRate(t *testing.T) {
	tests := []struct {
		name      string
		agentID   int64
		channelID int64
		rateType  string
		mockPrice *models.SettlementPrice
		mockErr   error
		expected  string
		wantErr   bool
	}{
		{
			name:      "credit card rate from RateConfigs",
			agentID:   1,
			channelID: 1,
			rateType:  "credit",
			mockPrice: &models.SettlementPrice{
				RateConfigs: models.RateConfigs{
					"credit": {Rate: "0.60"},
				},
			},
			expected: "0.60",
			wantErr:  false,
		},
		{
			name:      "debit card rate from legacy field",
			agentID:   1,
			channelID: 1,
			rateType:  "debit",
			mockPrice: &models.SettlementPrice{
				DebitRate: strPtr("0.50"),
			},
			expected: "0.50",
			wantErr:  false,
		},
		{
			name:      "non-existent agent",
			agentID:   999,
			channelID: 1,
			rateType:  "credit",
			mockErr:   assert.AnError,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSettlementPriceRepository)
			mockLogRepo := new(MockPriceChangeLogRepository)

			service := NewSettlementPriceService(mockRepo, mockLogRepo, nil)

			mockRepo.On("GetByAgentAndChannel", tt.agentID, tt.channelID, "").Return(tt.mockPrice, tt.mockErr)

			rate, err := service.GetAgentRate(tt.agentID, tt.channelID, "", tt.rateType)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, rate)
			}
		})
	}
}

func TestSettlementPriceService_UpdateDepositCashback(t *testing.T) {
	t.Run("should update deposit cashback config", func(t *testing.T) {
		mockRepo := new(MockSettlementPriceRepository)
		mockLogRepo := new(MockPriceChangeLogRepository)

		service := NewSettlementPriceService(mockRepo, mockLogRepo, nil)

		existingPrice := &models.SettlementPrice{
			ID:        1,
			AgentID:   1,
			ChannelID: 1,
			Version:   1,
			DepositCashbacks: models.DepositCashbacks{
				{DepositAmount: 9900, CashbackAmount: 5000},
			},
		}

		mockRepo.On("GetByID", int64(1)).Return(existingPrice, nil)
		mockRepo.On("Update", mock.AnythingOfType("*models.SettlementPrice")).Return(nil)
		mockLogRepo.On("Create", mock.AnythingOfType("*models.PriceChangeLog")).Return(nil)

		req := &models.UpdateDepositCashbackRequest{
			DepositCashbacks: models.DepositCashbacks{
				{DepositAmount: 9900, CashbackAmount: 6000},
				{DepositAmount: 19900, CashbackAmount: 12000},
			},
		}

		price, err := service.UpdateDepositCashback(1, req, 1, "admin", "PC", "127.0.0.1")

		assert.NoError(t, err)
		assert.NotNil(t, price)
		assert.Len(t, price.DepositCashbacks, 2)
		assert.Equal(t, int64(6000), price.DepositCashbacks[0].CashbackAmount)
	})
}

func TestSettlementPriceService_List(t *testing.T) {
	t.Run("should return settlement price list with pagination", func(t *testing.T) {
		mockRepo := new(MockSettlementPriceRepository)
		mockLogRepo := new(MockPriceChangeLogRepository)

		service := NewSettlementPriceService(mockRepo, mockLogRepo, nil)

		prices := []models.SettlementPrice{
			{ID: 1, AgentID: 1, ChannelID: 1, Version: 1},
			{ID: 2, AgentID: 2, ChannelID: 1, Version: 1},
		}

		req := &models.SettlementPriceListRequest{Page: 1, PageSize: 20}
		mockRepo.On("List", req).Return(prices, int64(2), nil)

		resp, err := service.List(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.List, 2)
		assert.Equal(t, int64(2), resp.Total)
	})

	t.Run("should return empty list when no data", func(t *testing.T) {
		mockRepo := new(MockSettlementPriceRepository)
		mockLogRepo := new(MockPriceChangeLogRepository)

		service := NewSettlementPriceService(mockRepo, mockLogRepo, nil)

		req := &models.SettlementPriceListRequest{Page: 1, PageSize: 20}
		mockRepo.On("List", req).Return([]models.SettlementPrice{}, int64(0), nil)

		resp, err := service.List(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.List, 0)
		assert.Equal(t, int64(0), resp.Total)
	})
}

func TestSettlementPriceService_GetAgentDepositCashback(t *testing.T) {
	t.Run("should return cashback amount for matching deposit", func(t *testing.T) {
		mockRepo := new(MockSettlementPriceRepository)
		mockLogRepo := new(MockPriceChangeLogRepository)

		service := NewSettlementPriceService(mockRepo, mockLogRepo, nil)

		price := &models.SettlementPrice{
			DepositCashbacks: models.DepositCashbacks{
				{DepositAmount: 9900, CashbackAmount: 5000},
				{DepositAmount: 19900, CashbackAmount: 12000},
			},
		}

		mockRepo.On("GetByAgentAndChannel", int64(1), int64(1), "").Return(price, nil)

		cashback, err := service.GetAgentDepositCashback(1, 1, "", 9900)

		assert.NoError(t, err)
		assert.Equal(t, int64(5000), cashback)
	})

	t.Run("should return 0 for non-matching deposit", func(t *testing.T) {
		mockRepo := new(MockSettlementPriceRepository)
		mockLogRepo := new(MockPriceChangeLogRepository)

		service := NewSettlementPriceService(mockRepo, mockLogRepo, nil)

		price := &models.SettlementPrice{
			DepositCashbacks: models.DepositCashbacks{
				{DepositAmount: 9900, CashbackAmount: 5000},
			},
		}

		mockRepo.On("GetByAgentAndChannel", int64(1), int64(1), "").Return(price, nil)

		cashback, err := service.GetAgentDepositCashback(1, 1, "", 29900)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), cashback)
	})
}

func TestSettlementPriceService_GetAgentSimCashback(t *testing.T) {
	t.Run("should return sim cashback config", func(t *testing.T) {
		mockRepo := new(MockSettlementPriceRepository)
		mockLogRepo := new(MockPriceChangeLogRepository)

		service := NewSettlementPriceService(mockRepo, mockLogRepo, nil)

		price := &models.SettlementPrice{
			SimFirstCashback:     5000,
			SimSecondCashback:    3000,
			SimThirdPlusCashback: 2000,
		}

		mockRepo.On("GetByAgentAndChannel", int64(1), int64(1), "").Return(price, nil)

		first, second, third, err := service.GetAgentSimCashback(1, 1, "")

		assert.NoError(t, err)
		assert.Equal(t, int64(5000), first)
		assert.Equal(t, int64(3000), second)
		assert.Equal(t, int64(2000), third)
	})
}

// 辅助函数
func strPtr(s string) *string {
	return &s
}

// ============================================================
// PriceChangeLog 测试用例
// ============================================================

func TestPriceChangeLogService_List(t *testing.T) {
	t.Run("should return change logs with pagination", func(t *testing.T) {
		mockRepo := new(MockPriceChangeLogRepository)
		service := NewPriceChangeLogService(mockRepo)

		logs := []models.PriceChangeLog{
			{
				ID:            1,
				AgentID:       1,
				ChangeType:    models.ChangeTypeRate,
				ConfigType:    models.ConfigTypeSettlement,
				ChangeSummary: "费率调整",
				OperatorName:  "admin",
				CreatedAt:     time.Now(),
			},
		}

		req := &models.PriceChangeLogListRequest{Page: 1, PageSize: 20}
		mockRepo.On("List", req).Return(logs, int64(1), nil)

		resp, err := service.List(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.List, 1)
		assert.Equal(t, "费率调整", resp.List[0].ChangeTypeName)
	})
}
