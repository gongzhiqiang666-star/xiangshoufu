<template>
  <div class="deposit-tier-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>押金档位管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增档位
          </el-button>
        </div>
      </template>

      <!-- 搜索区域 -->
      <el-form :inline="true" :model="searchParams" class="search-form">
        <el-form-item label="通道">
          <el-select
            v-model="searchParams.channel_id"
            placeholder="请选择通道"
            clearable
            style="width: 200px"
            @change="handleSearch"
          >
            <el-option
              v-for="ch in channels"
              :key="ch.id"
              :label="ch.channel_name"
              :value="ch.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="品牌">
          <el-input
            v-model="searchParams.brand_code"
            placeholder="品牌编码"
            clearable
            style="width: 150px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 档位列表 -->
      <el-table :data="tierList" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="tier_code" label="档位编码" width="120" />
        <el-table-column prop="tier_name" label="档位名称" width="150" />
        <el-table-column label="押金金额" width="120">
          <template #default="{ row }">
            <span class="amount">{{ formatAmount(row.deposit_amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="brand_code" label="品牌编码" width="100">
          <template #default="{ row }">
            {{ row.brand_code || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm
              title="确定要删除该押金档位吗？"
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button link type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="通道" prop="channel_id">
          <el-select
            v-model="form.channel_id"
            placeholder="请选择通道"
            :disabled="isEdit"
            style="width: 100%"
          >
            <el-option
              v-for="ch in channels"
              :key="ch.id"
              :label="ch.channel_name"
              :value="ch.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="品牌编码" prop="brand_code">
          <el-input v-model="form.brand_code" placeholder="如：DEFAULT" />
        </el-form-item>
        <el-form-item label="档位编码" prop="tier_code">
          <el-input v-model="form.tier_code" placeholder="如：TIER_99" />
        </el-form-item>
        <el-form-item label="档位名称" prop="tier_name">
          <el-input v-model="form.tier_name" placeholder="如：99元档" />
        </el-form-item>
        <el-form-item label="押金金额" prop="deposit_amount">
          <el-input-number
            v-model="form.deposit_amount"
            :min="0"
            :step="100"
            style="width: 100%"
          />
          <div class="form-tip">单位：分（{{ formatAmount(form.deposit_amount) }}）</div>
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="form.sort_order" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item v-if="isEdit" label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getDepositTiers,
  createDepositTier,
  updateDepositTier,
  deleteDepositTier,
  type DepositTier,
  type CreateDepositTierRequest,
  type UpdateDepositTierRequest
} from '@/api/depositTier'
import { getChannelList } from '@/api/channel'
import type { Channel } from '@/types'

// 通道列表（从API加载）
const channels = ref<Channel[]>([])

// 搜索参数
const searchParams = reactive({
  channel_id: undefined as number | undefined,
  brand_code: ''
})

// 列表数据
const tierList = ref<DepositTier[]>([])
const loading = ref(false)

// 弹窗控制
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

// 表单数据
const form = reactive<CreateDepositTierRequest & { status?: number }>({
  channel_id: 0,
  brand_code: '',
  tier_code: '',
  tier_name: '',
  deposit_amount: 0,
  sort_order: 0,
  status: 1
})

// 弹窗标题
const dialogTitle = computed(() => isEdit.value ? '编辑押金档位' : '新增押金档位')

// 表单验证规则
const rules: FormRules = {
  channel_id: [{ required: true, message: '请选择通道', trigger: 'change' }],
  tier_code: [{ required: true, message: '请输入档位编码', trigger: 'blur' }],
  tier_name: [{ required: true, message: '请输入档位名称', trigger: 'blur' }],
  deposit_amount: [{ required: true, message: '请输入押金金额', trigger: 'blur' }]
}

// 格式化金额（分转元）
const formatAmount = (amount: number): string => {
  if (!amount) return '0.00元'
  return (amount / 100).toFixed(2) + '元'
}

// 格式化日期时间
const formatDateTime = (dateStr: string): string => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

// 查询列表
const fetchList = async () => {
  if (!searchParams.channel_id) {
    tierList.value = []
    return
  }

  loading.value = true
  try {
    const data = await getDepositTiers(searchParams.channel_id, searchParams.brand_code || undefined)
    tierList.value = data || []
  } catch (error: any) {
    ElMessage.error(error.message || '获取列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  fetchList()
}

// 重置
const handleReset = () => {
  searchParams.channel_id = undefined
  searchParams.brand_code = ''
  tierList.value = []
}

// 新增
const handleAdd = () => {
  isEdit.value = false
  editingId.value = null
  Object.assign(form, {
    channel_id: searchParams.channel_id || 0,
    brand_code: '',
    tier_code: '',
    tier_name: '',
    deposit_amount: 0,
    sort_order: 0,
    status: 1
  })
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: DepositTier) => {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(form, {
    channel_id: row.channel_id,
    brand_code: row.brand_code,
    tier_code: row.tier_code,
    tier_name: row.tier_name,
    deposit_amount: row.deposit_amount,
    sort_order: row.sort_order,
    status: row.status
  })
  dialogVisible.value = true
}

// 删除
const handleDelete = async (row: DepositTier) => {
  try {
    await deleteDepositTier(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value && editingId.value) {
      const updateData: UpdateDepositTierRequest = {
        tier_code: form.tier_code,
        tier_name: form.tier_name,
        deposit_amount: form.deposit_amount,
        sort_order: form.sort_order,
        status: form.status
      }
      await updateDepositTier(editingId.value, updateData)
      ElMessage.success('更新成功')
    } else {
      await createDepositTier(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchList()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  // 加载通道列表
  loadChannels()
})

// 加载通道列表
async function loadChannels() {
  try {
    channels.value = await getChannelList()
  } catch (e) {
    console.error('加载通道列表失败', e)
  }
}
</script>

<style scoped>
.deposit-tier-management {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.amount {
  color: #e6a23c;
  font-weight: 500;
}

.form-tip {
  color: #909399;
  font-size: 12px;
  margin-top: 4px;
}
</style>
