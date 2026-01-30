<template>
  <div class="reward-template-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>奖励政策模版</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建模版
          </el-button>
        </div>
      </template>

      <!-- 筛选栏 -->
      <el-form :inline="true" class="filter-form">
        <el-form-item label="状态">
          <el-select v-model="filters.enabled" placeholder="全部" clearable @change="loadData" style="width: 100px">
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 表格 -->
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="模版名称" min-width="150" />
        <el-table-column label="时间类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.time_type === 'days' ? 'primary' : 'success'" size="small">
              {{ row.time_type === 'days' ? '按天数' : '按自然月' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="维度类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.dimension_type === 'amount' ? 'warning' : 'info'" size="small">
              {{ row.dimension_type === 'amount' ? '按金额' : '按笔数' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="交易类型" width="150">
          <template #default="{ row }">
            {{ formatTransTypes(row.trans_types) }}
          </template>
        </el-table-column>
        <el-table-column label="断档开关" width="100">
          <template #default="{ row }">
            <el-tag :type="row.allow_gap ? 'success' : 'danger'" size="small">
              {{ row.allow_gap ? '允许断档' : '不允许' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-switch
              v-model="row.enabled"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm
              title="确定要删除该模版吗？"
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <!-- 查看详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="模版详情" width="700px">
      <template v-if="currentTemplate">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="模版名称">{{ currentTemplate.name }}</el-descriptions-item>
          <el-descriptions-item label="时间类型">
            {{ currentTemplate.time_type === 'days' ? '按天数' : '按自然月' }}
          </el-descriptions-item>
          <el-descriptions-item label="维度类型">
            {{ currentTemplate.dimension_type === 'amount' ? '按金额' : '按笔数' }}
          </el-descriptions-item>
          <el-descriptions-item label="交易类型">
            {{ formatTransTypes(currentTemplate.trans_types) }}
          </el-descriptions-item>
          <el-descriptions-item label="金额限制">
            {{ formatAmountRange(currentTemplate.amount_min, currentTemplate.amount_max) }}
          </el-descriptions-item>
          <el-descriptions-item label="断档开关">
            {{ currentTemplate.allow_gap ? '允许断档' : '不允许断档' }}
          </el-descriptions-item>
        </el-descriptions>

        <div class="stages-section">
          <h4>阶段配置</h4>
          <el-table :data="currentTemplate.stages" border size="small">
            <el-table-column prop="stage_order" label="阶段" width="80" />
            <el-table-column label="时间范围">
              <template #default="{ row }">
                {{ row.start_value }} - {{ row.end_value }}
                {{ currentTemplate.time_type === 'days' ? '天' : '月' }}
              </template>
            </el-table-column>
            <el-table-column label="达标值">
              <template #default="{ row }">
                {{ currentTemplate.dimension_type === 'amount'
                  ? `${(row.target_value / 100).toFixed(2)}元`
                  : `${row.target_value}笔` }}
              </template>
            </el-table-column>
            <el-table-column label="奖励金额">
              <template #default="{ row }">
                {{ (row.reward_amount / 100).toFixed(2) }}元
              </template>
            </el-table-column>
          </el-table>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getRewardTemplates,
  getRewardTemplateDetail,
  deleteRewardTemplate,
  updateRewardTemplateStatus,
  type RewardTemplate,
} from '@/api/reward'

const router = useRouter()

const loading = ref(false)
const tableData = ref<RewardTemplate[]>([])
const detailDialogVisible = ref(false)
const currentTemplate = ref<RewardTemplate | null>(null)

const filters = reactive({
  enabled: undefined as boolean | undefined,
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
})

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const formatTransTypes = (types: string) => {
  if (!types) return '-'
  const typeMap: Record<string, string> = {
    scan: '扫码',
    debit: '借记卡',
    credit: '贷记卡',
  }
  return types.split(',').map(t => typeMap[t] || t).join('、')
}

const formatAmountRange = (min: number | null, max: number | null) => {
  if (min === null && max === null) return '不限'
  if (min !== null && max !== null) {
    return `${(min / 100).toFixed(2)}元 - ${(max / 100).toFixed(2)}元`
  }
  if (min !== null) return `≥${(min / 100).toFixed(2)}元`
  if (max !== null) return `<${(max / 100).toFixed(2)}元`
  return '-'
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getRewardTemplates({
      enabled: filters.enabled,
      page: pagination.page,
      page_size: pagination.pageSize,
    })
    tableData.value = res.list || []
    pagination.total = res.total || 0
  } catch (err) {
    console.error('加载数据失败:', err)
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  filters.enabled = undefined
  pagination.page = 1
  loadData()
}

const handleCreate = () => {
  router.push('/rewards/templates/create')
}

const handleView = async (row: RewardTemplate) => {
  try {
    const detail = await getRewardTemplateDetail(row.id)
    currentTemplate.value = detail
    detailDialogVisible.value = true
  } catch (err) {
    ElMessage.error('获取详情失败')
  }
}

const handleEdit = (row: RewardTemplate) => {
  router.push(`/rewards/templates/${row.id}/edit`)
}

const handleDelete = async (row: RewardTemplate) => {
  try {
    await deleteRewardTemplate(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (err) {
    ElMessage.error('删除失败')
  }
}

const handleStatusChange = async (row: RewardTemplate) => {
  try {
    await updateRewardTemplateStatus(row.id, row.enabled)
    ElMessage.success(row.enabled ? '已启用' : '已禁用')
  } catch (err) {
    row.enabled = !row.enabled
    ElMessage.error('状态更新失败')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-form {
  margin-bottom: 16px;
}

.pagination-wrapper {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.stages-section {
  margin-top: 20px;
}

.stages-section h4 {
  margin-bottom: 10px;
  color: #303133;
}
</style>
