package models

import (
	"encoding/json"
	"testing"
)

func TestSimCashbacks_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
		wantErr  bool
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "empty array",
			input:    []byte("[]"),
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "valid json array",
			input:    []byte(`[{"tier_order":1,"cashback_amount":3000},{"tier_order":2,"cashback_amount":2500}]`),
			expected: 2,
			wantErr:  false,
		},
		{
			name:     "string input",
			input:    `[{"tier_order":1,"cashback_amount":3000}]`,
			expected: 1,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s SimCashbacks
			err := s.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(s) != tt.expected {
				t.Errorf("Scan() got %d items, want %d", len(s), tt.expected)
			}
		})
	}
}

func TestSimCashbacks_Value(t *testing.T) {
	tests := []struct {
		name    string
		input   SimCashbacks
		wantErr bool
	}{
		{
			name:    "nil value",
			input:   nil,
			wantErr: false,
		},
		{
			name:    "empty array",
			input:   SimCashbacks{},
			wantErr: false,
		},
		{
			name: "with items",
			input: SimCashbacks{
				{TierOrder: 1, CashbackAmount: 3000},
				{TierOrder: 2, CashbackAmount: 2500},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("Value() returned nil")
			}
		})
	}
}

func TestSimCashbacks_GetCashbackAmount(t *testing.T) {
	cashbacks := SimCashbacks{
		{TierOrder: 1, CashbackAmount: 3000},
		{TierOrder: 2, CashbackAmount: 2500},
		{TierOrder: 3, CashbackAmount: 2000},
		{TierOrder: 4, CashbackAmount: 1500},
	}

	tests := []struct {
		name      string
		tierOrder int
		expected  int64
	}{
		{"first tier", 1, 3000},
		{"second tier", 2, 2500},
		{"third tier", 3, 2000},
		{"fourth tier", 4, 1500},
		{"beyond configured tiers", 5, 1500}, // should return last tier
		{"beyond configured tiers", 10, 1500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cashbacks.GetCashbackAmount(tt.tierOrder)
			if result != tt.expected {
				t.Errorf("GetCashbackAmount(%d) = %d, want %d", tt.tierOrder, result, tt.expected)
			}
		})
	}
}

func TestSimCashbacks_SetCashbackAmount(t *testing.T) {
	t.Run("update existing tier", func(t *testing.T) {
		cashbacks := SimCashbacks{
			{TierOrder: 1, CashbackAmount: 3000},
		}
		cashbacks.SetCashbackAmount(1, 3500)
		if cashbacks[0].CashbackAmount != 3500 {
			t.Errorf("SetCashbackAmount() did not update, got %d, want 3500", cashbacks[0].CashbackAmount)
		}
	})

	t.Run("add new tier", func(t *testing.T) {
		cashbacks := SimCashbacks{
			{TierOrder: 1, CashbackAmount: 3000},
		}
		cashbacks.SetCashbackAmount(2, 2500)
		if len(cashbacks) != 2 {
			t.Errorf("SetCashbackAmount() did not add new tier, got %d items", len(cashbacks))
		}
		if cashbacks[1].CashbackAmount != 2500 {
			t.Errorf("SetCashbackAmount() new tier amount = %d, want 2500", cashbacks[1].CashbackAmount)
		}
	})
}

func TestValidateRateRange(t *testing.T) {
	tests := []struct {
		name    string
		rate    string
		minRate string
		maxRate string
		wantErr bool
	}{
		{"rate within range", "0.50", "0.38", "0.60", false},
		{"rate at min", "0.38", "0.38", "0.60", false},
		{"rate at max", "0.60", "0.38", "0.60", false},
		{"rate below min", "0.30", "0.38", "0.60", true},
		{"rate above max", "0.70", "0.38", "0.60", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRateRange(tt.rate, tt.minRate, tt.maxRate)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSettlementRate(t *testing.T) {
	tests := []struct {
		name           string
		rate           string
		upperRate      string
		channelMinRate string
		channelMaxRate string
		wantErr        bool
	}{
		{"rate >= upper and within channel range", "0.55", "0.50", "0.38", "0.60", false},
		{"rate == upper and within channel range", "0.50", "0.50", "0.38", "0.60", false},
		{"rate < upper", "0.45", "0.50", "0.38", "0.60", true},
		{"rate below channel min", "0.30", "0.30", "0.38", "0.60", true},
		{"rate above channel max", "0.70", "0.50", "0.38", "0.60", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSettlementRate(tt.rate, tt.upperRate, tt.channelMinRate, tt.channelMaxRate)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSettlementRate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCashbackAmount(t *testing.T) {
	tests := []struct {
		name        string
		cashback    int64
		maxCashback int64
		wantErr     bool
	}{
		{"within limit", 3000, 5000, false},
		{"at limit", 5000, 5000, false},
		{"zero", 0, 5000, false},
		{"above limit", 6000, 5000, true},
		{"negative", -100, 5000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCashbackAmount(tt.cashback, tt.maxCashback)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCashbackAmount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSettlementCashback(t *testing.T) {
	tests := []struct {
		name           string
		cashback       int64
		upperCashback  int64
		channelMax     int64
		wantErr        bool
	}{
		{"within all limits", 3000, 4000, 5000, false},
		{"at upper limit", 4000, 4000, 5000, false},
		{"above upper limit", 4500, 4000, 5000, true},
		{"above channel max", 6000, 6000, 5000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSettlementCashback(tt.cashback, tt.upperCashback, tt.channelMax)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSettlementCashback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseRateToFloat(t *testing.T) {
	tests := []struct {
		name     string
		rate     string
		expected float64
	}{
		{"normal rate", "0.55", 0.55},
		{"zero", "0", 0},
		{"integer", "1", 1},
		{"high precision", "0.5678", 0.5678},
		{"empty string", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseRateToFloat(tt.rate)
			if result != tt.expected {
				t.Errorf("ParseRateToFloat(%s) = %f, want %f", tt.rate, result, tt.expected)
			}
		})
	}
}

func TestSimCashbackItem_JSON(t *testing.T) {
	item := SimCashbackItem{
		TierOrder:      1,
		CashbackAmount: 3000,
	}

	// Test marshaling
	data, err := json.Marshal(item)
	if err != nil {
		t.Errorf("json.Marshal() error = %v", err)
	}

	expected := `{"tier_order":1,"cashback_amount":3000}`
	if string(data) != expected {
		t.Errorf("json.Marshal() = %s, want %s", string(data), expected)
	}

	// Test unmarshaling
	var result SimCashbackItem
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
	}
	if result.TierOrder != 1 || result.CashbackAmount != 3000 {
		t.Errorf("json.Unmarshal() got %+v, want %+v", result, item)
	}
}

func TestChannelRateConfig_TableName(t *testing.T) {
	config := ChannelRateConfig{}
	if config.TableName() != "channel_rate_configs" {
		t.Errorf("TableName() = %s, want channel_rate_configs", config.TableName())
	}
}

func TestChannelSimCashbackTier_TableName(t *testing.T) {
	tier := ChannelSimCashbackTier{}
	if tier.TableName() != "channel_sim_cashback_tiers" {
		t.Errorf("TableName() = %s, want channel_sim_cashback_tiers", tier.TableName())
	}
}
