import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  getSettlementPrices,
  getSettlementPrice,
  createSettlementPrice,
  updateSettlementPriceRate,
  updateSettlementPriceDeposit,
  updateSettlementPriceSim,
  getPriceChangeLogs,
  getPriceChangeLog,
  type SettlementPrice,
  type SettlementPriceListResponse,
  type PriceChangeLog,
  type PriceChangeLogListResponse,
} from '../settlementPrice'
import * as request from '../request'

// Mock request module
vi.mock('../request', () => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
}))

describe('settlementPrice API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  // ============================================================
  // 结算价API测试
  // ============================================================

  describe('getSettlementPrices', () => {
    it('should fetch settlement price list successfully', async () => {
      const mockResponse: SettlementPriceListResponse = {
        list: [
          {
            id: 1,
            agent_id: 100,
            agent_name: '测试代理商',
            channel_id: 1,
            channel_name: '恒信通',
            brand_code: '',
            rate_configs: { credit: { rate: '0.60' } },
            deposit_cashbacks: [{ deposit_amount: 9900, cashback_amount: 5000 }],
            sim_first_cashback: 5000,
            sim_second_cashback: 3000,
            sim_third_plus_cashback: 2000,
            version: 1,
            status: 1,
            effective_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
          },
        ],
        total: 1,
        page: 1,
        size: 20,
      }

      vi.mocked(request.get).mockResolvedValue(mockResponse)

      const result = await getSettlementPrices({ page: 1, page_size: 20 })

      expect(request.get).toHaveBeenCalledWith('/v1/settlement-prices', { page: 1, page_size: 20 })
      expect(result.list).toHaveLength(1)
      expect(result.total).toBe(1)
    })

    it('should handle empty list', async () => {
      const mockResponse: SettlementPriceListResponse = {
        list: [],
        total: 0,
        page: 1,
        size: 20,
      }

      vi.mocked(request.get).mockResolvedValue(mockResponse)

      const result = await getSettlementPrices({})

      expect(result.list).toHaveLength(0)
      expect(result.total).toBe(0)
    })

    it('should pass filter parameters correctly', async () => {
      vi.mocked(request.get).mockResolvedValue({ list: [], total: 0, page: 1, size: 20 })

      await getSettlementPrices({
        agent_id: 100,
        channel_id: 1,
        status: 1,
        page: 2,
        page_size: 50,
      })

      expect(request.get).toHaveBeenCalledWith('/v1/settlement-prices', {
        agent_id: 100,
        channel_id: 1,
        status: 1,
        page: 2,
        page_size: 50,
      })
    })
  })

  describe('getSettlementPrice', () => {
    it('should fetch settlement price detail successfully', async () => {
      const mockPrice: SettlementPrice = {
        id: 1,
        agent_id: 100,
        agent_name: '测试代理商',
        channel_id: 1,
        channel_name: '恒信通',
        brand_code: '',
        rate_configs: { credit: { rate: '0.60' }, debit: { rate: '0.50' } },
        credit_rate: '0.60',
        debit_rate: '0.50',
        deposit_cashbacks: [],
        sim_first_cashback: 5000,
        sim_second_cashback: 3000,
        sim_third_plus_cashback: 2000,
        version: 1,
        status: 1,
        effective_at: '2024-01-01T00:00:00Z',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(request.get).mockResolvedValue(mockPrice)

      const result = await getSettlementPrice(1)

      expect(request.get).toHaveBeenCalledWith('/v1/settlement-prices/1')
      expect(result.id).toBe(1)
      expect(result.credit_rate).toBe('0.60')
    })

    it('should handle API error', async () => {
      vi.mocked(request.get).mockRejectedValue(new Error('Not found'))

      await expect(getSettlementPrice(999)).rejects.toThrow('Not found')
    })
  })

  describe('createSettlementPrice', () => {
    it('should create settlement price successfully', async () => {
      const mockPrice: SettlementPrice = {
        id: 1,
        agent_id: 100,
        channel_id: 1,
        agent_name: '',
        channel_name: '',
        brand_code: '',
        rate_configs: {},
        deposit_cashbacks: [],
        sim_first_cashback: 0,
        sim_second_cashback: 0,
        sim_third_plus_cashback: 0,
        version: 1,
        status: 1,
        effective_at: '',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(request.post).mockResolvedValue(mockPrice)

      const result = await createSettlementPrice({
        agent_id: 100,
        channel_id: 1,
      })

      expect(request.post).toHaveBeenCalledWith('/v1/settlement-prices', {
        agent_id: 100,
        channel_id: 1,
      })
      expect(result.id).toBe(1)
    })
  })

  describe('updateSettlementPriceRate', () => {
    it('should update rate successfully', async () => {
      const mockPrice: SettlementPrice = {
        id: 1,
        agent_id: 100,
        channel_id: 1,
        agent_name: '',
        channel_name: '',
        brand_code: '',
        rate_configs: { credit: { rate: '0.55' } },
        credit_rate: '0.55',
        deposit_cashbacks: [],
        sim_first_cashback: 0,
        sim_second_cashback: 0,
        sim_third_plus_cashback: 0,
        version: 2,
        status: 1,
        effective_at: '',
        created_at: '',
        updated_at: '2024-01-02T00:00:00Z',
      }

      vi.mocked(request.put).mockResolvedValue(mockPrice)

      const result = await updateSettlementPriceRate(1, { credit_rate: '0.55' })

      expect(request.put).toHaveBeenCalledWith('/v1/settlement-prices/1/rate', { credit_rate: '0.55' })
      expect(result.version).toBe(2)
      expect(result.credit_rate).toBe('0.55')
    })
  })

  describe('updateSettlementPriceDeposit', () => {
    it('should update deposit cashback successfully', async () => {
      const mockPrice: SettlementPrice = {
        id: 1,
        agent_id: 100,
        channel_id: 1,
        agent_name: '',
        channel_name: '',
        brand_code: '',
        rate_configs: {},
        deposit_cashbacks: [
          { deposit_amount: 9900, cashback_amount: 6000 },
          { deposit_amount: 19900, cashback_amount: 12000 },
        ],
        sim_first_cashback: 0,
        sim_second_cashback: 0,
        sim_third_plus_cashback: 0,
        version: 2,
        status: 1,
        effective_at: '',
        created_at: '',
        updated_at: '2024-01-02T00:00:00Z',
      }

      vi.mocked(request.put).mockResolvedValue(mockPrice)

      const result = await updateSettlementPriceDeposit(1, {
        deposit_cashbacks: [
          { deposit_amount: 9900, cashback_amount: 6000 },
          { deposit_amount: 19900, cashback_amount: 12000 },
        ],
      })

      expect(request.put).toHaveBeenCalledWith('/v1/settlement-prices/1/deposit', {
        deposit_cashbacks: [
          { deposit_amount: 9900, cashback_amount: 6000 },
          { deposit_amount: 19900, cashback_amount: 12000 },
        ],
      })
      expect(result.deposit_cashbacks).toHaveLength(2)
    })
  })

  describe('updateSettlementPriceSim', () => {
    it('should update sim cashback successfully', async () => {
      const mockPrice: SettlementPrice = {
        id: 1,
        agent_id: 100,
        channel_id: 1,
        agent_name: '',
        channel_name: '',
        brand_code: '',
        rate_configs: {},
        deposit_cashbacks: [],
        sim_first_cashback: 6000,
        sim_second_cashback: 4000,
        sim_third_plus_cashback: 3000,
        version: 2,
        status: 1,
        effective_at: '',
        created_at: '',
        updated_at: '2024-01-02T00:00:00Z',
      }

      vi.mocked(request.put).mockResolvedValue(mockPrice)

      const result = await updateSettlementPriceSim(1, {
        sim_first_cashback: 6000,
        sim_second_cashback: 4000,
        sim_third_plus_cashback: 3000,
      })

      expect(request.put).toHaveBeenCalledWith('/v1/settlement-prices/1/sim', {
        sim_first_cashback: 6000,
        sim_second_cashback: 4000,
        sim_third_plus_cashback: 3000,
      })
      expect(result.sim_first_cashback).toBe(6000)
    })
  })

  // ============================================================
  // 调价记录API测试
  // ============================================================

  describe('getPriceChangeLogs', () => {
    it('should fetch price change logs successfully', async () => {
      const mockResponse: PriceChangeLogListResponse = {
        list: [
          {
            id: 1,
            agent_id: 100,
            agent_name: '测试代理商',
            channel_id: 1,
            channel_name: '恒信通',
            change_type: 2,
            change_type_name: '费率调整',
            config_type: 1,
            config_type_name: '结算价',
            field_name: 'credit_rate',
            old_value: '0.60',
            new_value: '0.55',
            change_summary: '贷记卡费率: 0.60% → 0.55%',
            operator_name: 'admin',
            source: 'PC',
            created_at: '2024-01-02T10:00:00Z',
          },
        ],
        total: 1,
        page: 1,
        size: 20,
      }

      vi.mocked(request.get).mockResolvedValue(mockResponse)

      const result = await getPriceChangeLogs({ page: 1, page_size: 20 })

      expect(request.get).toHaveBeenCalledWith('/v1/price-change-logs', { page: 1, page_size: 20 })
      expect(result.list).toHaveLength(1)
      expect(result.list[0].change_type_name).toBe('费率调整')
    })

    it('should filter by change type', async () => {
      vi.mocked(request.get).mockResolvedValue({ list: [], total: 0, page: 1, size: 20 })

      await getPriceChangeLogs({ change_type: 2 })

      expect(request.get).toHaveBeenCalledWith('/v1/price-change-logs', { change_type: 2 })
    })

    it('should filter by date range', async () => {
      vi.mocked(request.get).mockResolvedValue({ list: [], total: 0, page: 1, size: 20 })

      await getPriceChangeLogs({
        start_date: '2024-01-01',
        end_date: '2024-01-31',
      })

      expect(request.get).toHaveBeenCalledWith('/v1/price-change-logs', {
        start_date: '2024-01-01',
        end_date: '2024-01-31',
      })
    })
  })

  describe('getPriceChangeLog', () => {
    it('should fetch price change log detail successfully', async () => {
      const mockLog: PriceChangeLog = {
        id: 1,
        agent_id: 100,
        agent_name: '测试代理商',
        channel_id: 1,
        channel_name: '恒信通',
        change_type: 2,
        change_type_name: '费率调整',
        config_type: 1,
        config_type_name: '结算价',
        field_name: 'credit_rate',
        old_value: '0.60',
        new_value: '0.55',
        change_summary: '贷记卡费率: 0.60% → 0.55%',
        operator_name: 'admin',
        source: 'PC',
        created_at: '2024-01-02T10:00:00Z',
      }

      vi.mocked(request.get).mockResolvedValue(mockLog)

      const result = await getPriceChangeLog(1)

      expect(request.get).toHaveBeenCalledWith('/v1/price-change-logs/1')
      expect(result.id).toBe(1)
      expect(result.change_type).toBe(2)
    })
  })
})
