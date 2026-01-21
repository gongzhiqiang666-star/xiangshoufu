<template>
  <div class="merchant-list-view">
    <PageHeader title="商户管理" sub-title="商户列表">
      <template #extra>
        <el-button :icon="Download" :loading="exporting" @click="handleExport">导出Excel</el-button>
      </template>
    </PageHeader>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <el-row :gutter="16">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon total">
                <el-icon size="24"><Shop /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.total_count }}</div>
                <div class="stat-label">商户总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon direct">
                <el-icon size="24"><User /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.direct_count }}</div>
                <div class="stat-label">直营商户</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon team">
                <el-icon size="24"><UserFilled /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.team_count }}</div>
                <div class="stat-label">团队商户</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon today">
                <el-icon size="24"><Plus /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.today_new_count }}</div>
                <div class="stat-label">今日新增</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="通道">
        <ChannelSelect v-model="searchForm.channel_id" />
      </el-form-item>
      <el-form-item label="商户类型">
        <el-select v-model="searchForm.merchant_type" placeholder="请选择类型" clearable>
          <el-option label="忠诚商户" value="loyal" />
          <el-option label="优质商户" value="quality" />
          <el-option label="潜力商户" value="potential" />
          <el-option label="一般商户" value="normal" />
          <el-option label="低活跃" value="low_active" />
          <el-option label="30天无交易" value="inactive" />
        </el-select>
      </el-form-item>
      <el-form-item label="归属类型">
        <el-radio-group v-model="searchForm.owner_type">
          <el-radio-button label="all">全部</el-radio-button>
          <el-radio-button label="direct">直营</el-radio-button>
          <el-radio-button label="team">团队</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="关键词">
        <el-input
          v-model="searchForm.keyword"
          placeholder="商户名称/编号/机具号"
          clearable
        />
      </el-form-item>
    </SearchForm>

    <!-- 表格 -->
    <ProTable
      :data="tableData"
      :loading="loading"
      :total="total"
      v-model:page="page"
      v-model:page-size="pageSize"
      :show-export="true"
      @refresh="fetchData"
      @export="handleExport"
    >
      <el-table-column prop="merchant_code" label="商户编号" width="120" />
      <el-table-column prop="name" label="商户姓名" width="100" />
      <el-table-column prop="phone_masked" label="手机号" width="130" />
      <el-table-column prop="is_direct" label="归属类型" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_direct ? 'primary' : ''" size="small">
            {{ row.is_direct ? '直营' : '团队' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="agent_name" label="所属代理" width="100" />
      <el-table-column prop="merchant_type" label="商户类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getMerchantTypeTag(row.merchant_type)" size="small">
            {{ getMerchantTypeLabel(row.merchant_type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="terminal_sn" label="机具号" width="130" />
      <el-table-column prop="channel_name" label="通道" width="100" />
      <el-table-column prop="activated_at" label="激活时间" width="170" />

      <template #action="{ row }">
        <el-button type="primary" link @click="handleView(row)">详情</el-button>
        <el-button type="primary" link @click="handleRegister(row)">
          {{ row.registered_phone ? '编辑登记' : '登记' }}
        </el-button>
      </template>
    </ProTable>

    <!-- 商户登记弹窗 -->
    <el-dialog
      v-model="registerDialogVisible"
      :title="currentMerchant?.registered_phone ? '编辑商户登记' : '商户登记'"
      width="500px"
    >
      <el-form :model="registerForm" label-width="100px">
        <el-form-item label="商户编号">
          {{ currentMerchant?.merchant_code }}
        </el-form-item>
        <el-form-item label="商户姓名">
          {{ currentMerchant?.name }}
        </el-form-item>
        <el-form-item label="完整手机号" required>
          <el-input v-model="registerForm.phone" placeholder="请输入完整手机号" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="registerForm.remark"
            type="textarea"
            placeholder="请输入备注信息"
            :rows="3"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="registerDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitRegister">确认登记</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Download, Shop, User, UserFilled, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import ChannelSelect from '@/components/Common/ChannelSelect.vue'
import { getMerchants, registerMerchant, updateMerchantRegister, getExtendedMerchantStats, exportMerchants } from '@/api/merchant'
import type { ExtendedMerchantStats } from '@/api/merchant'
import type { Merchant, MerchantType } from '@/types'
import { MERCHANT_TYPE_CONFIG } from '@/types/merchant'

const router = useRouter()

// 搜索表单
const searchForm = reactive({
  channel_id: undefined as number | undefined,
  merchant_type: undefined as MerchantType | undefined,
  owner_type: 'all' as 'all' | 'direct' | 'team',
  keyword: '',
})

// 统计数据
const stats = ref<ExtendedMerchantStats>({
  total_count: 0,
  active_count: 0,
  pending_count: 0,
  disabled_count: 0,
  direct_count: 0,
  team_count: 0,
  today_new_count: 0,
  loyal_count: 0,
  quality_count: 0,
  potential_count: 0,
  normal_count: 0,
  low_active_count: 0,
  inactive_count: 0,
})

// 表格数据
const tableData = ref<Merchant[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const exporting = ref(false)

// 登记弹窗
const registerDialogVisible = ref(false)
const currentMerchant = ref<Merchant | null>(null)
const registerForm = reactive({
  phone: '',
  remark: '',
})

// 获取商户类型标签颜色
function getMerchantTypeTag(type: MerchantType) {
  const config = MERCHANT_TYPE_CONFIG[type]
  if (!config) return ''

  const colorMap: Record<string, string> = {
    '#67c23a': 'success',
    '#409eff': 'primary',
    '#e6a23c': 'warning',
    '#f56c6c': 'danger',
    '#909399': 'info',
    '#c0c4cc': 'info',
  }
  return colorMap[config.color] || ''
}

// 获取商户类型标签文本
function getMerchantTypeLabel(type: MerchantType) {
  return MERCHANT_TYPE_CONFIG[type]?.label || type
}

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getMerchants({
      ...searchForm,
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('Fetch merchants error:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
function handleSearch() {
  page.value = 1
  fetchData()
}

// 重置
function handleReset() {
  page.value = 1
  fetchData()
}

// 查看详情
function handleView(row: Merchant) {
  router.push(`/merchants/${row.id}`)
}

// 登记
function handleRegister(row: Merchant) {
  currentMerchant.value = row
  registerForm.phone = ''
  registerForm.remark = ''
  registerDialogVisible.value = true
}

// 提交登记
async function handleSubmitRegister() {
  if (!currentMerchant.value) return

  if (!registerForm.phone) {
    ElMessage.warning('请输入完整手机号')
    return
  }

  try {
    if ((currentMerchant.value as any).registered_phone) {
      await updateMerchantRegister(currentMerchant.value.id, registerForm)
    } else {
      await registerMerchant(currentMerchant.value.id, registerForm)
    }
    ElMessage.success('登记成功')
    registerDialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('Register merchant error:', error)
  }
}

// 导出
async function handleExport() {
  exporting.value = true
  try {
    const params: any = { ...searchForm }
    if (params.owner_type === 'direct') {
      params.is_direct = true
    } else if (params.owner_type === 'team') {
      params.is_direct = false
    }
    delete params.owner_type

    const blob = await exportMerchants(params)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `商户列表_${new Date().toISOString().slice(0, 10)}.csv`
    a.click()
    URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('Export merchants error:', error)
    ElMessage.error('导出失败')
  } finally {
    exporting.value = false
  }
}

// 获取统计数据
async function fetchStats() {
  try {
    stats.value = await getExtendedMerchantStats()
  } catch (error) {
    console.error('Fetch stats error:', error)
  }
}

onMounted(() => {
  fetchData()
  fetchStats()
})
</script>

<style lang="scss" scoped>
.merchant-list-view {
  padding: 0;
}

.stats-cards {
  margin-bottom: 16px;

  .stat-card {
    :deep(.el-card__body) {
      padding: 16px;
    }
  }

  .stat-content {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .stat-icon {
    width: 48px;
    height: 48px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;

    &.total {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    }

    &.direct {
      background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
    }

    &.team {
      background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
    }

    &.today {
      background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
    }
  }

  .stat-info {
    flex: 1;
  }

  .stat-value {
    font-size: 24px;
    font-weight: 600;
    color: #303133;
    line-height: 1.2;
  }

  .stat-label {
    font-size: 13px;
    color: #909399;
    margin-top: 4px;
  }
}
</style>
