<template>
  <div class="price-change-log-list">
    <div class="page-header">
      <h2>调价记录</h2>
    </div>

    <!-- 搜索栏 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="代理商ID">
          <el-input v-model.number="searchForm.agent_id" placeholder="请输入代理商ID" clearable />
        </el-form-item>
        <el-form-item label="通道ID">
          <el-input v-model.number="searchForm.channel_id" placeholder="请输入通道ID" clearable />
        </el-form-item>
        <el-form-item label="变更类型">
          <el-select v-model="searchForm.change_type" placeholder="请选择" clearable>
            <el-option label="初始化" :value="1" />
            <el-option label="费率调整" :value="2" />
            <el-option label="押金返现调整" :value="3" />
            <el-option label="流量费返现调整" :value="4" />
            <el-option label="激活奖励调整" :value="5" />
          </el-select>
        </el-form-item>
        <el-form-item label="配置类型">
          <el-select v-model="searchForm.config_type" placeholder="请选择" clearable>
            <el-option label="结算价" :value="1" />
            <el-option label="奖励配置" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="开始日期">
          <el-date-picker v-model="searchForm.start_date" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item label="结束日期">
          <el-date-picker v-model="searchForm.end_date" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card class="table-card">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="agent_id" label="代理商ID" width="100" />
        <el-table-column prop="agent_name" label="代理商名称" min-width="120" />
        <el-table-column prop="channel_id" label="通道ID" width="80" />
        <el-table-column prop="channel_name" label="通道名称" min-width="100" />
        <el-table-column prop="change_type_name" label="变更类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getChangeTypeTagType(row.change_type)">
              {{ row.change_type_name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="config_type_name" label="配置类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.config_type === 1 ? 'primary' : 'warning'">
              {{ row.config_type_name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="change_summary" label="变更摘要" min-width="150" />
        <el-table-column prop="operator_name" label="操作人" width="100" />
        <el-table-column prop="source" label="来源" width="80">
          <template #default="{ row }">
            <el-tag :type="row.source === 'PC' ? 'info' : 'success'" size="small">
              {{ row.source }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="操作时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleViewDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="调价记录详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ currentLog?.id }}</el-descriptions-item>
        <el-descriptions-item label="代理商">{{ currentLog?.agent_name }} (ID: {{ currentLog?.agent_id }})</el-descriptions-item>
        <el-descriptions-item label="通道">{{ currentLog?.channel_name }}</el-descriptions-item>
        <el-descriptions-item label="变更类型">{{ currentLog?.change_type_name }}</el-descriptions-item>
        <el-descriptions-item label="配置类型">{{ currentLog?.config_type_name }}</el-descriptions-item>
        <el-descriptions-item label="变更字段">{{ currentLog?.field_name }}</el-descriptions-item>
        <el-descriptions-item label="变更摘要" :span="2">{{ currentLog?.change_summary }}</el-descriptions-item>
        <el-descriptions-item label="操作人">{{ currentLog?.operator_name }}</el-descriptions-item>
        <el-descriptions-item label="操作来源">{{ currentLog?.source }}</el-descriptions-item>
        <el-descriptions-item label="操作时间" :span="2">{{ formatDateTime(currentLog?.created_at || '') }}</el-descriptions-item>
      </el-descriptions>
      <template v-if="currentLog?.old_value || currentLog?.new_value">
        <el-divider content-position="left">变更内容</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <div class="value-label">变更前</div>
            <pre class="value-content">{{ formatJson(currentLog?.old_value) }}</pre>
          </el-col>
          <el-col :span="12">
            <div class="value-label">变更后</div>
            <pre class="value-content">{{ formatJson(currentLog?.new_value) }}</pre>
          </el-col>
        </el-row>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getPriceChangeLogs, type PriceChangeLog } from '@/api/settlementPrice'

// 搜索表单
const searchForm = reactive({
  agent_id: undefined as number | undefined,
  channel_id: undefined as number | undefined,
  change_type: undefined as number | undefined,
  config_type: undefined as number | undefined,
  start_date: '',
  end_date: '',
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
})

// 数据
const tableData = ref<PriceChangeLog[]>([])
const loading = ref(false)

// 详情对话框
const detailDialogVisible = ref(false)
const currentLog = ref<PriceChangeLog | null>(null)

// 格式化日期时间
const formatDateTime = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

// 格式化JSON
const formatJson = (str: string | undefined) => {
  if (!str) return ''
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

// 获取变更类型标签颜色
const getChangeTypeTagType = (type: number): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const typeMap: Record<number, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    1: 'info',
    2: 'primary',
    3: 'success',
    4: 'warning',
    5: 'danger',
  }
  return typeMap[type] || 'info'
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const resp = await getPriceChangeLogs({
      agent_id: searchForm.agent_id,
      channel_id: searchForm.channel_id,
      change_type: searchForm.change_type,
      config_type: searchForm.config_type,
      start_date: searchForm.start_date,
      end_date: searchForm.end_date,
      page: pagination.page,
      page_size: pagination.pageSize,
    })
    tableData.value = resp.list || []
    pagination.total = resp.total
  } catch (e) {
    console.error('加载调价记录失败', e)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadData()
}

// 重置
const handleReset = () => {
  searchForm.agent_id = undefined
  searchForm.channel_id = undefined
  searchForm.change_type = undefined
  searchForm.config_type = undefined
  searchForm.start_date = ''
  searchForm.end_date = ''
  handleSearch()
}

// 分页
const handleSizeChange = () => {
  pagination.page = 1
  loadData()
}
const handleCurrentChange = () => {
  loadData()
}

// 查看详情
const handleViewDetail = (row: PriceChangeLog) => {
  currentLog.value = row
  detailDialogVisible.value = true
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.price-change-log-list {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
}

.search-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.value-label {
  font-weight: bold;
  margin-bottom: 10px;
  color: #606266;
}

.value-content {
  background: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  font-size: 12px;
  max-height: 300px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
