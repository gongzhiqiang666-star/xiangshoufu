<template>
  <div class="job-list-view">
    <el-card class="page-header">
      <div class="header-content">
        <h2>定时任务管理</h2>
        <p class="description">管理系统定时任务，查看运行状态，配置重试和告警策略</p>
      </div>
    </el-card>

    <el-card class="main-content">
      <el-table
        v-loading="loading"
        :data="jobList"
        stripe
        border
        style="width: 100%"
      >
        <el-table-column prop="job_name" label="任务名称" width="220">
          <template #default="{ row }">
            <el-link type="primary" @click="showJobDetail(row)">
              {{ row.job_name }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="job_desc" label="任务描述" min-width="200" />
        <el-table-column prop="interval_seconds" label="执行间隔" width="120">
          <template #default="{ row }">
            {{ formatInterval(row.interval_seconds) }}
          </template>
        </el-table-column>
        <el-table-column prop="is_enabled" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_enabled"
              @change="handleEnableChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="is_running" label="运行状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_running" type="warning" effect="dark">
              <el-icon class="is-loading"><Loading /></el-icon>
              运行中
            </el-tag>
            <el-tag v-else type="info">空闲</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_status" label="上次结果" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.last_status === 1" type="success">成功</el-tag>
            <el-tag v-else-if="row.last_status === 2" type="danger">失败</el-tag>
            <el-tag v-else-if="row.last_status === 3" type="warning">运行中</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="last_run_at" label="上次执行" width="180">
          <template #default="{ row }">
            {{ row.last_run_at ? formatDateTime(row.last_run_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              :disabled="row.is_running"
              @click="handleTrigger(row)"
            >
              立即执行
            </el-button>
            <el-button
              type="default"
              size="small"
              @click="showConfigDialog(row)"
            >
              配置
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 任务详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="currentJob?.job_name + ' - 任务详情'"
      width="800px"
    >
      <div v-if="jobDetail" class="job-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="任务名称">{{ jobDetail.config.job_name }}</el-descriptions-item>
          <el-descriptions-item label="任务描述">{{ jobDetail.config.job_desc }}</el-descriptions-item>
          <el-descriptions-item label="执行间隔">{{ formatInterval(jobDetail.config.interval_seconds) }}</el-descriptions-item>
          <el-descriptions-item label="是否启用">
            <el-tag :type="jobDetail.config.is_enabled ? 'success' : 'info'">
              {{ jobDetail.config.is_enabled ? '启用' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="最大重试次数">{{ jobDetail.config.max_retries }} 次</el-descriptions-item>
          <el-descriptions-item label="告警阈值">连续失败 {{ jobDetail.config.alert_threshold }} 次</el-descriptions-item>
          <el-descriptions-item label="超时时间">{{ jobDetail.config.timeout_seconds }} 秒</el-descriptions-item>
          <el-descriptions-item label="运行状态">
            <el-tag :type="jobDetail.is_running ? 'warning' : 'info'">
              {{ jobDetail.is_running ? '运行中' : '空闲' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <h4 style="margin: 20px 0 10px">最近执行记录</h4>
        <el-table :data="jobDetail.latest_logs" size="small" border>
          <el-table-column prop="started_at" label="开始时间" width="180">
            <template #default="{ row }">
              {{ formatDateTime(row.started_at) }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="80" align="center">
            <template #default="{ row }">
              <el-tag v-if="row.status === 1" type="success" size="small">成功</el-tag>
              <el-tag v-else-if="row.status === 2" type="danger" size="small">失败</el-tag>
              <el-tag v-else type="warning" size="small">运行中</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="duration_ms" label="耗时" width="100">
            <template #default="{ row }">
              {{ row.duration_ms ? (row.duration_ms / 1000).toFixed(2) + 's' : '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="processed_count" label="处理/成功/失败" width="130">
            <template #default="{ row }">
              {{ row.processed_count }}/{{ row.success_count }}/{{ row.fail_count }}
            </template>
          </el-table-column>
          <el-table-column prop="error_message" label="错误信息" show-overflow-tooltip />
        </el-table>
      </div>
    </el-dialog>

    <!-- 配置对话框 -->
    <el-dialog
      v-model="configDialogVisible"
      title="任务配置"
      width="500px"
    >
      <el-form
        ref="configFormRef"
        :model="configForm"
        :rules="configRules"
        label-width="120px"
      >
        <el-form-item label="任务描述" prop="job_desc">
          <el-input v-model="configForm.job_desc" placeholder="请输入任务描述" />
        </el-form-item>
        <el-form-item label="执行间隔(秒)" prop="interval_seconds">
          <el-input-number
            v-model="configForm.interval_seconds"
            :min="60"
            :max="86400"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="最大重试次数" prop="max_retries">
          <el-input-number
            v-model="configForm.max_retries"
            :min="0"
            :max="10"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="重试间隔(秒)" prop="retry_interval">
          <el-input-number
            v-model="configForm.retry_interval"
            :min="10"
            :max="3600"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="告警阈值" prop="alert_threshold">
          <el-input-number
            v-model="configForm.alert_threshold"
            :min="1"
            :max="10"
            style="width: 100%"
          />
          <div class="form-tip">连续失败达到此次数时触发告警</div>
        </el-form-item>
        <el-form-item label="超时时间(秒)" prop="timeout_seconds">
          <el-input-number
            v-model="configForm.timeout_seconds"
            :min="60"
            :max="7200"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSaveConfig">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { getJobs, getJob, updateJobConfig, triggerJob, enableJob } from '@/api/job'
import type { JobListItem, JobDetail, UpdateJobConfigRequest } from '@/types/job'

const loading = ref(false)
const saving = ref(false)
const jobList = ref<JobListItem[]>([])
const currentJob = ref<JobListItem | null>(null)
const jobDetail = ref<JobDetail | null>(null)
const detailDialogVisible = ref(false)
const configDialogVisible = ref(false)

const configForm = ref<UpdateJobConfigRequest>({
  job_desc: '',
  interval_seconds: 300,
  max_retries: 3,
  retry_interval: 60,
  alert_threshold: 3,
  timeout_seconds: 3600
})

const configRules = {
  interval_seconds: [{ required: true, message: '请输入执行间隔', trigger: 'blur' }],
  max_retries: [{ required: true, message: '请输入最大重试次数', trigger: 'blur' }]
}

// 加载任务列表
const loadJobs = async () => {
  loading.value = true
  try {
    jobList.value = await getJobs()
  } catch (error) {
    console.error('加载任务列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 格式化执行间隔
const formatInterval = (seconds: number) => {
  if (seconds < 60) return `${seconds}秒`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}分钟`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}小时`
  return `${Math.floor(seconds / 86400)}天`
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

// 显示任务详情
const showJobDetail = async (job: JobListItem) => {
  currentJob.value = job
  detailDialogVisible.value = true
  try {
    jobDetail.value = await getJob(job.job_name)
  } catch (error) {
    console.error('加载任务详情失败:', error)
  }
}

// 显示配置对话框
const showConfigDialog = async (job: JobListItem) => {
  currentJob.value = job
  try {
    const detail = await getJob(job.job_name)
    configForm.value = {
      job_desc: detail.config.job_desc,
      interval_seconds: detail.config.interval_seconds,
      max_retries: detail.config.max_retries,
      retry_interval: detail.config.retry_interval,
      alert_threshold: detail.config.alert_threshold,
      timeout_seconds: detail.config.timeout_seconds
    }
    configDialogVisible.value = true
  } catch (error) {
    console.error('加载任务配置失败:', error)
  }
}

// 保存配置
const handleSaveConfig = async () => {
  if (!currentJob.value) return
  saving.value = true
  try {
    await updateJobConfig(currentJob.value.job_name, configForm.value)
    ElMessage.success('配置保存成功')
    configDialogVisible.value = false
    loadJobs()
  } catch (error) {
    console.error('保存配置失败:', error)
  } finally {
    saving.value = false
  }
}

// 启用/禁用任务
const handleEnableChange = async (job: JobListItem) => {
  try {
    await enableJob(job.job_name, job.is_enabled)
    ElMessage.success(job.is_enabled ? '任务已启用' : '任务已禁用')
  } catch (error) {
    job.is_enabled = !job.is_enabled
    console.error('更新任务状态失败:', error)
  }
}

// 手动触发任务
const handleTrigger = async (job: JobListItem) => {
  try {
    await ElMessageBox.confirm(
      `确定要立即执行任务 "${job.job_name}" 吗？`,
      '确认执行',
      { type: 'warning' }
    )
    await triggerJob(job.job_name)
    ElMessage.success('任务已触发，请稍后查看执行结果')
    setTimeout(loadJobs, 2000)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('触发任务失败:', error)
    }
  }
}

onMounted(() => {
  loadJobs()
})
</script>

<style scoped>
.job-list-view {
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

.main-content {
  min-height: 400px;
}

.job-detail {
  max-height: 500px;
  overflow-y: auto;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>
