<template>
  <div class="log-list-view">
    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <template #extra>
        <el-button :icon="Download" @click="handleExport">导出日志</el-button>
      </template>
      <el-form-item label="操作用户">
        <el-input v-model="searchForm.username" placeholder="用户名" clearable />
      </el-form-item>
      <el-form-item label="操作模块">
        <el-select v-model="searchForm.module" placeholder="请选择模块" clearable style="width: 150px">
          <el-option
            v-for="(config, key) in LOG_MODULE_CONFIG"
            :key="key"
            :label="config.label"
            :value="key"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="日期范围">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          @change="handleDateChange"
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
      @refresh="fetchData"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="操作用户" width="120">
        <template #default="{ row }">
          <div>{{ row.nickname }}</div>
          <div class="sub-text">{{ row.username }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="module" label="模块" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getModuleTagType(row.module)" size="small">
            {{ getModuleLabel(row.module) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="action" label="操作" width="120" />
      <el-table-column prop="method" label="请求方式" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="getMethodType(row.method)" size="small" effect="plain">
            {{ row.method }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="path" label="请求路径" min-width="200" show-overflow-tooltip />
      <el-table-column prop="ip" label="IP地址" width="140" />
      <el-table-column prop="response_code" label="状态码" width="80" align="center">
        <template #default="{ row }">
          <el-tag
            :type="row.response_code >= 200 && row.response_code < 300 ? 'success' : 'danger'"
            size="small"
          >
            {{ row.response_code }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="response_time" label="耗时" width="80" align="right">
        <template #default="{ row }">
          <span :class="{ 'slow-request': row.response_time > 1000 }">
            {{ row.response_time }}ms
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="操作时间" width="170" />

      <template #action="{ row }">
        <el-button type="primary" link @click="handleView(row)">详情</el-button>
      </template>
    </ProTable>

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailDialogVisible" title="日志详情" width="700px">
      <el-descriptions v-if="currentLog" :column="2" border>
        <el-descriptions-item label="日志ID">{{ currentLog.id }}</el-descriptions-item>
        <el-descriptions-item label="操作时间">{{ currentLog.created_at }}</el-descriptions-item>
        <el-descriptions-item label="操作用户">
          {{ currentLog.nickname }} ({{ currentLog.username }})
        </el-descriptions-item>
        <el-descriptions-item label="操作模块">
          <el-tag :type="getModuleTagType(currentLog.module)" size="small">
            {{ getModuleLabel(currentLog.module) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="操作描述" :span="2">{{ currentLog.action }}</el-descriptions-item>
        <el-descriptions-item label="请求方式">
          <el-tag :type="getMethodType(currentLog.method)" size="small" effect="plain">
            {{ currentLog.method }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态码">
          <el-tag
            :type="currentLog.response_code >= 200 && currentLog.response_code < 300 ? 'success' : 'danger'"
            size="small"
          >
            {{ currentLog.response_code }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="请求路径" :span="2">{{ currentLog.path }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ currentLog.ip }}</el-descriptions-item>
        <el-descriptions-item label="响应耗时">
          <span :class="{ 'slow-request': currentLog.response_time > 1000 }">
            {{ currentLog.response_time }}ms
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="User-Agent" :span="2">
          <span class="user-agent">{{ currentLog.user_agent }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="请求参数" :span="2">
          <div class="request-body">
            <pre>{{ formatJson(currentLog.request_body) }}</pre>
          </div>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import { getLogs, exportLogs } from '@/api/system'
import type { OperationLog, LogModule } from '@/types/system'
import { LOG_MODULE_CONFIG } from '@/types/system'

// 搜索表单
const searchForm = reactive({
  username: '',
  module: undefined as LogModule | undefined,
  start_date: '',
  end_date: '',
})

const dateRange = ref<[string, string] | null>(null)

// 表格数据
const tableData = ref<OperationLog[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 详情弹窗
const detailDialogVisible = ref(false)
const currentLog = ref<OperationLog | null>(null)

// 获取模块标签类型
type TagType = 'primary' | 'success' | 'warning' | 'info' | 'danger'
function getModuleTagType(module: LogModule): TagType {
  const typeMap: Record<string, TagType> = {
    '#409eff': 'primary',
    '#67c23a': 'success',
    '#e6a23c': 'warning',
    '#f56c6c': 'danger',
    '#909399': 'info',
  }
  return typeMap[LOG_MODULE_CONFIG[module]?.color] || 'info'
}

// 获取模块名称
function getModuleLabel(module: LogModule) {
  return LOG_MODULE_CONFIG[module]?.label || module
}

// 获取请求方式类型
function getMethodType(method: string): TagType {
  const methodTypes: Record<string, TagType> = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger',
    PATCH: 'info',
  }
  return methodTypes[method] || 'info'
}

// 格式化JSON
function formatJson(str: string) {
  if (!str) return '-'
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

// 处理日期变化
function handleDateChange(val: [string, string] | null) {
  if (val) {
    searchForm.start_date = val[0]
    searchForm.end_date = val[1]
  } else {
    searchForm.start_date = ''
    searchForm.end_date = ''
  }
}

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getLogs({
      ...searchForm,
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list || []
    total.value = res.total
  } catch (error) {
    console.error('Fetch logs error:', error)
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
  dateRange.value = null
  page.value = 1
  fetchData()
}

// 查看详情
function handleView(row: OperationLog) {
  currentLog.value = row
  detailDialogVisible.value = true
}

// 导出
async function handleExport() {
  try {
    const blob = await exportLogs({
      ...searchForm,
    })
    const url = window.URL.createObjectURL(blob as Blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `操作日志_${new Date().toISOString().slice(0, 10)}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    ElMessage.error('导出失败')
    console.error('Export logs error:', error)
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.log-list-view {
  padding: 0;
}

.sub-text {
  font-size: 12px;
  color: $text-placeholder;
}

.slow-request {
  color: $danger-color;
  font-weight: 600;
}

.user-agent {
  font-size: 12px;
  color: $text-secondary;
  word-break: break-all;
}

.request-body {
  max-height: 200px;
  overflow: auto;
  background: $bg-color;
  border-radius: $border-radius-sm;
  padding: $spacing-sm;

  pre {
    margin: 0;
    font-size: 12px;
    white-space: pre-wrap;
    word-break: break-all;
  }
}
</style>
