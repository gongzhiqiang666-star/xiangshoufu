<template>
  <div class="job-log-view">
    <el-card class="page-header">
      <div class="header-content">
        <h2>执行日志</h2>
        <p class="description">查看定时任务执行记录，分析任务运行情况</p>
      </div>
    </el-card>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="任务名称">
          <el-select
            v-model="filterForm.job_name"
            placeholder="全部任务"
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="job in jobList"
              :key="job.job_name"
              :label="job.job_name"
              :value="job.job_name"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="filterForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 280px"
          />
        </el-form-item>
        <el-form-item label="执行状态">
          <el-select
            v-model="filterForm.status"
            placeholder="全部状态"
            clearable
            style="width: 120px"
          >
            <el-option label="成功" :value="1" />
            <el-option label="失败" :value="2" />
            <el-option label="运行中" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="main-content">
      <el-table
        v-loading="loading"
        :data="logList"
        stripe
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="job_name" label="任务名称" width="220" />
        <el-table-column prop="started_at" label="开始时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.started_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="ended_at" label="结束时间" width="180">
          <template #default="{ row }">
            {{ row.ended_at ? formatDateTime(row.ended_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="duration_ms" label="耗时" width="100">
          <template #default="{ row }">
            {{ row.duration_ms ? (row.duration_ms / 1000).toFixed(2) + 's' : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.status === 1" type="success">成功</el-tag>
            <el-tag v-else-if="row.status === 2" type="danger">失败</el-tag>
            <el-tag v-else type="warning">
              <el-icon class="is-loading"><Loading /></el-icon>
              运行中
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="trigger_type" label="触发方式" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.trigger_type === 1" type="info">自动</el-tag>
            <el-tag v-else type="warning">手动</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="processed_count" label="处理条数" width="100" align="center" />
        <el-table-column prop="success_count" label="成功" width="80" align="center">
          <template #default="{ row }">
            <span class="text-success">{{ row.success_count }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="fail_count" label="失败" width="80" align="center">
          <template #default="{ row }">
            <span :class="row.fail_count > 0 ? 'text-danger' : ''">{{ row.fail_count }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="retry_count" label="重试次数" width="90" align="center" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="showLogDetail(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadLogs"
          @current-change="loadLogs"
        />
      </div>
    </el-card>

    <!-- 日志详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="执行日志详情"
      width="700px"
    >
      <div v-if="currentLog" class="log-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="日志ID">{{ currentLog.id }}</el-descriptions-item>
          <el-descriptions-item label="任务名称">{{ currentLog.job_name }}</el-descriptions-item>
          <el-descriptions-item label="开始时间">{{ formatDateTime(currentLog.started_at) }}</el-descriptions-item>
          <el-descriptions-item label="结束时间">
            {{ currentLog.ended_at ? formatDateTime(currentLog.ended_at) : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="执行耗时">
            {{ currentLog.duration_ms ? (currentLog.duration_ms / 1000).toFixed(2) + '秒' : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="执行状态">
            <el-tag v-if="currentLog.status === 1" type="success">成功</el-tag>
            <el-tag v-else-if="currentLog.status === 2" type="danger">失败</el-tag>
            <el-tag v-else type="warning">运行中</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="触发方式">
            {{ currentLog.trigger_type === 1 ? '自动触发' : '手动触发' }}
          </el-descriptions-item>
          <el-descriptions-item label="重试次数">{{ currentLog.retry_count }} 次</el-descriptions-item>
          <el-descriptions-item label="处理条数">{{ currentLog.processed_count }}</el-descriptions-item>
          <el-descriptions-item label="成功/失败">
            <span class="text-success">{{ currentLog.success_count }}</span> /
            <span class="text-danger">{{ currentLog.fail_count }}</span>
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="currentLog.error_message" class="error-section">
          <h4>错误信息</h4>
          <el-input
            v-model="currentLog.error_message"
            type="textarea"
            :rows="3"
            readonly
          />
        </div>

        <div v-if="currentLog.error_stack" class="error-section">
          <h4>错误堆栈</h4>
          <el-input
            v-model="currentLog.error_stack"
            type="textarea"
            :rows="8"
            readonly
          />
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { Loading } from '@element-plus/icons-vue'
import { getJobLogs, getJobs } from '@/api/job'
import type { JobListItem, JobExecutionLog } from '@/types/job'

const loading = ref(false)
const logList = ref<JobExecutionLog[]>([])
const jobList = ref<JobListItem[]>([])
const currentLog = ref<JobExecutionLog | null>(null)
const detailDialogVisible = ref(false)

const filterForm = reactive({
  job_name: '',
  dateRange: [] as string[],
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 加载任务列表（用于筛选）
const loadJobList = async () => {
  try {
    jobList.value = await getJobs()
  } catch (error) {
    console.error('加载任务列表失败:', error)
  }
}

// 加载执行日志
const loadLogs = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (filterForm.job_name) {
      params.job_name = filterForm.job_name
    }
    if (filterForm.dateRange && filterForm.dateRange.length === 2) {
      params.start_date = filterForm.dateRange[0]
      params.end_date = filterForm.dateRange[1]
    }
    if (filterForm.status !== undefined) {
      params.status = filterForm.status
    }

    const result = await getJobLogs(params)
    logList.value = result.list || []
    pagination.total = result.total || 0
  } catch (error) {
    console.error('加载执行日志失败:', error)
  } finally {
    loading.value = false
  }
}

// 格式化日期时间
const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 查询
const handleSearch = () => {
  pagination.page = 1
  loadLogs()
}

// 重置
const handleReset = () => {
  filterForm.job_name = ''
  filterForm.dateRange = []
  filterForm.status = undefined
  pagination.page = 1
  loadLogs()
}

// 显示日志详情
const showLogDetail = (log: JobExecutionLog) => {
  currentLog.value = log
  detailDialogVisible.value = true
}

onMounted(() => {
  loadJobList()
  loadLogs()
})
</script>

<style scoped>
.job-log-view {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.header-content h2 {
  margin: 0 0 8px 0;
  font-size: 20px;
}

.header-content .description {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.main-content {
  min-height: 400px;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.log-detail {
  max-height: 500px;
  overflow-y: auto;
}

.error-section {
  margin-top: 20px;
}

.error-section h4 {
  margin: 0 0 10px 0;
  font-size: 14px;
  color: #303133;
}

.text-success {
  color: #67c23a;
}

.text-danger {
  color: #f56c6c;
}
</style>
