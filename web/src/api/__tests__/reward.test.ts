import { describe, it, expect, vi, beforeEach } from 'vitest'
import * as request from '../request'

// Mock request模块
vi.mock('../request', () => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  del: vi.fn(),
}))

import {
  getRewardTemplates,
  getRewardTemplateDetail,
  createRewardTemplate,
  updateRewardTemplate,
  deleteRewardTemplate,
  updateRewardTemplateStatus,
  getAgentRewardAmount,
  setAgentRewardAmount,
  getTerminalRewardProgress,
  getOverflowLogs,
  resolveOverflowLog,
} from '../reward'

describe('reward API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  // ============================================================
  // 奖励模版管理测试
  // ============================================================

  describe('getRewardTemplates', () => {
    it('应该正确调用获取模版列表API', async () => {
      const mockResponse = {
        list: [{ id: 1, name: '测试模版' }],
        total: 1,
        page: 1,
        page_size: 20,
      }
      vi.mocked(request.get).mockResolvedValue(mockResponse)

      const result = await getRewardTemplates({ page: 1, page_size: 20 })

      expect(request.get).toHaveBeenCalledWith('/v1/rewards/templates', { page: 1, page_size: 20 })
      expect(result.list).toHaveLength(1)
      expect(result.total).toBe(1)
    })

    it('应该支持按启用状态筛选', async () => {
      vi.mocked(request.get).mockResolvedValue({ list: [], total: 0 })

      await getRewardTemplates({ enabled: true, page: 1 })

      expect(request.get).toHaveBeenCalledWith('/v1/rewards/templates', { enabled: true, page: 1 })
    })
  })

  describe('getRewardTemplateDetail', () => {
    it('应该正确获取模版详情', async () => {
      const mockTemplate = {
        id: 1,
        name: '测试模版',
        time_type: 'days',
        dimension_type: 'amount',
        stages: [{ stage_order: 1, start_value: 1, end_value: 10 }],
      }
      vi.mocked(request.get).mockResolvedValue(mockTemplate)

      const result = await getRewardTemplateDetail(1)

      expect(request.get).toHaveBeenCalledWith('/v1/rewards/templates/1')
      expect(result.name).toBe('测试模版')
      expect(result.stages).toHaveLength(1)
    })
  })

  describe('createRewardTemplate', () => {
    it('应该正确创建模版', async () => {
      const newTemplate = {
        name: '新模版',
        time_type: 'days' as const,
        dimension_type: 'amount' as const,
        trans_types: 'scan,debit',
        allow_gap: false,
        stages: [{ stage_order: 1, start_value: 1, end_value: 10, target_value: 10000, reward_amount: 5000 }],
      }
      vi.mocked(request.post).mockResolvedValue({ id: 2 })

      const result = await createRewardTemplate(newTemplate)

      expect(request.post).toHaveBeenCalledWith('/v1/rewards/templates', newTemplate)
      expect(result.id).toBe(2)
    })
  })

  describe('updateRewardTemplate', () => {
    it('应该正确更新模版', async () => {
      const updateData = {
        name: '更新后的模版',
        time_type: 'days' as const,
        dimension_type: 'amount' as const,
        trans_types: 'scan',
        allow_gap: true,
        stages: [],
      }
      vi.mocked(request.put).mockResolvedValue(undefined)

      await updateRewardTemplate(1, updateData)

      expect(request.put).toHaveBeenCalledWith('/v1/rewards/templates/1', updateData)
    })
  })

  describe('deleteRewardTemplate', () => {
    it('应该正确删除模版', async () => {
      vi.mocked(request.del).mockResolvedValue(undefined)

      await deleteRewardTemplate(1)

      expect(request.del).toHaveBeenCalledWith('/v1/rewards/templates/1')
    })
  })

  describe('updateRewardTemplateStatus', () => {
    it('应该正确更新模版状态', async () => {
      vi.mocked(request.put).mockResolvedValue(undefined)

      await updateRewardTemplateStatus(1, true)

      expect(request.put).toHaveBeenCalledWith('/v1/rewards/templates/1/status', { enabled: true })
    })

    it('应该支持禁用模版', async () => {
      vi.mocked(request.put).mockResolvedValue(undefined)

      await updateRewardTemplateStatus(1, false)

      expect(request.put).toHaveBeenCalledWith('/v1/rewards/templates/1/status', { enabled: false })
    })
  })

  // ============================================================
  // 代理商奖励金额配置测试（差额分配模式）
  // ============================================================

  describe('getAgentRewardAmount', () => {
    it('应该正确获取代理商奖励金额配置', async () => {
      const mockRate = { agent_id: 1, template_id: 1, reward_amount: 10000 }
      vi.mocked(request.get).mockResolvedValue(mockRate)

      const result = await getAgentRewardAmount(1, 1)

      expect(request.get).toHaveBeenCalledWith('/v1/rewards/agents/1/amount', { template_id: 1 })
      expect(result.reward_amount).toBe(10000)
    })
  })

  describe('setAgentRewardAmount', () => {
    it('应该正确设置代理商奖励金额', async () => {
      vi.mocked(request.put).mockResolvedValue(undefined)

      await setAgentRewardAmount(1, 1, 15000)

      expect(request.put).toHaveBeenCalledWith('/v1/rewards/agents/1/amount', { template_id: 1, reward_amount: 15000 })
    })

    it('应该支持设置0金额', async () => {
      vi.mocked(request.put).mockResolvedValue(undefined)

      await setAgentRewardAmount(1, 1, 0)

      expect(request.put).toHaveBeenCalledWith('/v1/rewards/agents/1/amount', { template_id: 1, reward_amount: 0 })
    })
  })

  // ============================================================
  // 终端奖励进度测试
  // ============================================================

  describe('getTerminalRewardProgress', () => {
    it('应该正确获取终端奖励进度', async () => {
      const mockProgress = {
        terminal_sn: 'SN001',
        current_stage: 2,
        status: 'active',
      }
      vi.mocked(request.get).mockResolvedValue(mockProgress)

      const result = await getTerminalRewardProgress('SN001')

      expect(request.get).toHaveBeenCalledWith('/v1/rewards/terminals/SN001/progress')
      expect(result.current_stage).toBe(2)
    })
  })

  // ============================================================
  // 溢出日志测试
  // ============================================================

  describe('getOverflowLogs', () => {
    it('应该正确获取溢出日志列表', async () => {
      const mockLogs = {
        list: [{ id: 1, terminal_sn: 'SN001', total_rate: 1.1 }],
        total: 1,
      }
      vi.mocked(request.get).mockResolvedValue(mockLogs)

      const result = await getOverflowLogs({ page: 1, page_size: 20 })

      expect(request.get).toHaveBeenCalledWith('/v1/rewards/overflow-logs', { page: 1, page_size: 20 })
      expect(result.list).toHaveLength(1)
    })
  })

  describe('resolveOverflowLog', () => {
    it('应该正确解决溢出日志', async () => {
      vi.mocked(request.post).mockResolvedValue(undefined)

      await resolveOverflowLog(1)

      expect(request.post).toHaveBeenCalledWith('/v1/rewards/overflow-logs/1/resolve')
    })
  })
})
