<template>
  <div class="message-list-view">
    <PageHeader title="消息管理" sub-title="消息列表">
      <template #extra>
        <el-button type="primary" :icon="Plus" @click="handleSendMessage">发送消息</el-button>
      </template>
    </PageHeader>

    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="消息类型">
        <el-select v-model="searchForm.message_type" placeholder="请选择" clearable>
          <el-option
            v-for="(config, key) in MESSAGE_TYPE_CONFIG"
            :key="key"
            :label="config.label"
            :value="Number(key)"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="阅读状态">
        <el-select v-model="searchForm.is_read" placeholder="请选择" clearable>
          <el-option label="未读" :value="false" />
          <el-option label="已读" :value="true" />
        </el-select>
      </el-form-item>
    </SearchForm>

    <!-- 统计汇总 -->
    <el-card class="summary-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="summary-item">
            <span class="label">总消息数:</span>
            <span class="value">{{ stats.total }}</span>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="summary-item">
            <span class="label">未读消息:</span>
            <span class="value warning">{{ stats.unread_total }}</span>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="summary-item">
            <span class="label">分润类:</span>
            <span class="value">{{ stats.profit_count }}</span>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="summary-item">
            <span class="label">系统类:</span>
            <span class="value">{{ stats.system_count }}</span>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 表格 -->
    <ProTable
      :data="tableData"
      :loading="loading"
      :total="total"
      v-model:page="page"
      v-model:page-size="pageSize"
      @refresh="fetchData"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="agent_id" label="代理商ID" width="100" />
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column prop="message_type" label="消息类型" width="120" align="center">
        <template #default="{ row }">
          <el-tag :type="getTypeTag(row.message_type)" size="small">
            {{ row.type_name || getTypeLabel(row.message_type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="is_read" label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_read ? 'info' : 'danger'" size="small">
            {{ row.is_read ? '已读' : '未读' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="expire_at" label="过期时间" width="170">
        <template #default="{ row }">
          {{ formatDate(row.expire_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="170">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>

      <template #action="{ row }">
        <el-button type="primary" link @click="handleView(row)">详情</el-button>
        <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
      </template>
    </ProTable>

    <!-- 消息详情弹窗 -->
    <el-dialog v-model="detailDialogVisible" title="消息详情" width="600px">
      <el-descriptions :column="2" border v-if="currentMessage">
        <el-descriptions-item label="ID">{{ currentMessage.id }}</el-descriptions-item>
        <el-descriptions-item label="代理商ID">{{ currentMessage.agent_id }}</el-descriptions-item>
        <el-descriptions-item label="标题" :span="2">{{ currentMessage.title }}</el-descriptions-item>
        <el-descriptions-item label="消息类型">
          <el-tag :type="getTypeTag(currentMessage.message_type)" size="small">
            {{ currentMessage.type_name }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="阅读状态">
          <el-tag :type="currentMessage.is_read ? 'info' : 'danger'" size="small">
            {{ currentMessage.is_read ? '已读' : '未读' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="内容" :span="2">{{ currentMessage.content }}</el-descriptions-item>
        <el-descriptions-item label="过期时间">{{ formatDate(currentMessage.expire_at) }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentMessage.created_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import { getAdminMessages, getAdminMessageDetail, deleteMessage } from '@/api/message'
import { formatDate } from '@/utils/format'
import type { Message, MessageTypeValue, MessageStats } from '@/types'
import { MESSAGE_TYPE_CONFIG } from '@/types/message'

const router = useRouter()

// 统计数据
const stats = ref<MessageStats>({
  total: 0,
  unread_total: 0,
  profit_count: 0,
  register_count: 0,
  consumption_count: 0,
  system_count: 0,
})

// 搜索表单
const searchForm = reactive({
  message_type: undefined as MessageTypeValue | undefined,
  is_read: undefined as boolean | undefined,
})

// 表格数据
const tableData = ref<Message[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 详情弹窗
const detailDialogVisible = ref(false)
const currentMessage = ref<Message | null>(null)

// 获取类型标签颜色
function getTypeTag(type: MessageTypeValue): string {
  const config = MESSAGE_TYPE_CONFIG[type]
  if (!config) return ''
  const colorMap: Record<string, string> = {
    '#67c23a': 'success',
    '#e6a23c': 'warning',
    '#409eff': 'primary',
    '#909399': 'info',
    '#f56c6c': 'danger',
  }
  return colorMap[config.color] || ''
}

function getTypeLabel(type: MessageTypeValue): string {
  return MESSAGE_TYPE_CONFIG[type]?.label || String(type)
}

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getAdminMessages({
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list
    total.value = res.total

    // 更新统计（简化：使用列表数据统计）
    stats.value.total = res.total
    stats.value.unread_total = res.list.filter((m: Message) => !m.is_read).length
  } catch (error) {
    console.error('Fetch messages error:', error)
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
async function handleView(row: Message) {
  try {
    currentMessage.value = await getAdminMessageDetail(row.id)
    detailDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取消息详情失败')
  }
}

// 删除
async function handleDelete(row: Message) {
  try {
    await ElMessageBox.confirm('确定要删除该消息吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await deleteMessage(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 发送消息
function handleSendMessage() {
  router.push('/system/messages/send')
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.message-list-view {
  padding: 0;
}

.summary-card {
  margin-bottom: 16px;

  .summary-item {
    display: flex;
    align-items: center;
    gap: 8px;

    .label {
      color: #909399;
    }

    .value {
      font-size: 18px;
      font-weight: 600;
      color: #409eff;

      &.warning {
        color: #e6a23c;
      }
    }
  }
}
</style>
