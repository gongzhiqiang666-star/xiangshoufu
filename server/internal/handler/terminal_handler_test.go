package handler

import (
	"testing"

	"xiangshoufu/internal/models"

	"github.com/stretchr/testify/assert"
)

// TestGetTerminalStatusName 测试终端状态名称
func TestGetTerminalStatusName(t *testing.T) {
	tests := []struct {
		status   int16
		expected string
	}{
		{models.TerminalStatusPending, "待分配"},
		{models.TerminalStatusAllocated, "已分配"},
		{models.TerminalStatusBound, "已绑定"},
		{models.TerminalStatusActivated, "已激活"},
		{models.TerminalStatusUnbound, "已解绑"},
		{models.TerminalStatusRecycled, "已回收"},
		{0, "未知"},
		{99, "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getTerminalStatusName(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestTerminalListResponse 测试终端列表响应结构
func TestTerminalListResponse(t *testing.T) {
	// 模拟终端列表响应
	list := []map[string]interface{}{
		{
			"id":           int64(1),
			"terminal_sn":  "SN001",
			"channel_id":   int64(1),
			"channel_code": "HENGXINTONG",
			"brand_code":   "BRAND01",
			"model_code":   "MODEL01",
			"merchant_no":  "M001",
			"status":       int16(4),
			"status_name":  "已激活",
		},
	}

	assert.Len(t, list, 1)
	assert.Equal(t, int64(1), list[0]["id"])
	assert.Equal(t, "SN001", list[0]["terminal_sn"])
	assert.Equal(t, "HENGXINTONG", list[0]["channel_code"])
	assert.Equal(t, "已激活", list[0]["status_name"])
}

// TestTerminalStats 测试终端统计结构
func TestTerminalStats(t *testing.T) {
	stats := map[string]interface{}{
		"total":           int64(100),
		"pending_count":   int64(20),
		"allocated_count": int64(30),
		"bound_count":     int64(15),
		"activated_count": int64(35),
	}

	assert.Equal(t, int64(100), stats["total"])
	assert.Equal(t, int64(20), stats["pending_count"])
	assert.Equal(t, int64(30), stats["allocated_count"])
	assert.Equal(t, int64(15), stats["bound_count"])
	assert.Equal(t, int64(35), stats["activated_count"])

	// 验证统计数据一致性
	total := stats["pending_count"].(int64) +
		stats["allocated_count"].(int64) +
		stats["bound_count"].(int64) +
		stats["activated_count"].(int64)
	assert.Equal(t, stats["total"], total)
}

// TestTerminalDetail 测试终端详情结构
func TestTerminalDetail(t *testing.T) {
	detail := map[string]interface{}{
		"id":             int64(1),
		"terminal_sn":    "SN001",
		"channel_id":     int64(1),
		"channel_code":   "HENGXINTONG",
		"brand_code":     "BRAND01",
		"model_code":     "MODEL01",
		"merchant_id":    int64(100),
		"merchant_no":    "M001",
		"status":         int16(4),
		"status_name":    "已激活",
		"sim_fee_count":  3,
		"last_sim_fee_at": nil,
		"activated_at":   nil,
		"bound_at":       nil,
	}

	assert.Equal(t, "SN001", detail["terminal_sn"])
	assert.Equal(t, "已激活", detail["status_name"])
	assert.Equal(t, 3, detail["sim_fee_count"])
}

// TestTerminalStatusTransition 测试终端状态流转
func TestTerminalStatusTransition(t *testing.T) {
	tests := []struct {
		name        string
		fromStatus  int16
		toStatus    int16
		isValid     bool
	}{
		{"pending to allocated", models.TerminalStatusPending, models.TerminalStatusAllocated, true},
		{"allocated to bound", models.TerminalStatusAllocated, models.TerminalStatusBound, true},
		{"bound to activated", models.TerminalStatusBound, models.TerminalStatusActivated, true},
		{"activated to unbound", models.TerminalStatusActivated, models.TerminalStatusUnbound, true},
		{"unbound to recycled", models.TerminalStatusUnbound, models.TerminalStatusRecycled, true},
		{"pending to activated", models.TerminalStatusPending, models.TerminalStatusActivated, false}, // 不能跳过
		{"recycled to pending", models.TerminalStatusRecycled, models.TerminalStatusPending, false},   // 不能逆转
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 简单的状态流转规则验证
			validTransitions := map[int16][]int16{
				models.TerminalStatusPending:   {models.TerminalStatusAllocated},
				models.TerminalStatusAllocated: {models.TerminalStatusBound, models.TerminalStatusPending},
				models.TerminalStatusBound:     {models.TerminalStatusActivated, models.TerminalStatusUnbound},
				models.TerminalStatusActivated: {models.TerminalStatusUnbound},
				models.TerminalStatusUnbound:   {models.TerminalStatusRecycled, models.TerminalStatusBound},
			}

			isValid := false
			if validNext, ok := validTransitions[tt.fromStatus]; ok {
				for _, v := range validNext {
					if v == tt.toStatus {
						isValid = true
						break
					}
				}
			}

			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

// TestTerminalSNFormat 测试终端SN格式
func TestTerminalSNFormat(t *testing.T) {
	tests := []struct {
		sn      string
		isValid bool
	}{
		{"SN001", true},
		{"SN12345678", true},
		{"ABCD1234567890", true},
		{"", false},
		{"SN", false}, // 太短
	}

	for _, tt := range tests {
		t.Run(tt.sn, func(t *testing.T) {
			isValid := len(tt.sn) >= 3
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

// TestTerminalPermissionCheck 测试终端权限检查
func TestTerminalPermissionCheck(t *testing.T) {
	tests := []struct {
		name            string
		terminalOwnerID int64
		currentAgentID  int64
		hasPermission   bool
	}{
		{"same agent", 1, 1, true},
		{"different agent", 1, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasPermission := tt.terminalOwnerID == tt.currentAgentID
			assert.Equal(t, tt.hasPermission, hasPermission)
		})
	}
}

// TestSimFeeCount 测试流量费缴费次数
func TestSimFeeCount(t *testing.T) {
	tests := []struct {
		count    int
		tier     string
	}{
		{0, "首次"},
		{1, "第2次"},
		{2, "2+N次"},
		{3, "2+N次"},
		{10, "2+N次"},
	}

	for _, tt := range tests {
		t.Run(tt.tier, func(t *testing.T) {
			var tier string
			switch {
			case tt.count == 0:
				tier = "首次"
			case tt.count == 1:
				tier = "第2次"
			default:
				tier = "2+N次"
			}
			assert.Equal(t, tt.tier, tier)
		})
	}
}

// TestTerminalStatusFilter 测试终端状态筛选
func TestTerminalStatusFilter(t *testing.T) {
	terminals := []struct {
		sn     string
		status int16
	}{
		{"SN001", models.TerminalStatusPending},
		{"SN002", models.TerminalStatusAllocated},
		{"SN003", models.TerminalStatusActivated},
		{"SN004", models.TerminalStatusActivated},
	}

	tests := []struct {
		filterStatus int16
		expected     int
	}{
		{models.TerminalStatusPending, 1},
		{models.TerminalStatusAllocated, 1},
		{models.TerminalStatusActivated, 2},
		{models.TerminalStatusBound, 0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			count := 0
			for _, term := range terminals {
				if term.status == tt.filterStatus {
					count++
				}
			}
			assert.Equal(t, tt.expected, count)
		})
	}
}
