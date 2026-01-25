<template>
  <div class="overflow-logs">
    <el-card>
      <template #header>
        <span>奖励池溢出日志</span>
      </template>

      <el-alert
        type="warning"
        :closable="false"
        style="margin-bottom: 16px"
      >
        当代理商链上所有层级的奖励比例之和超过100%时，会产生溢出记录。请及时处理以确保奖励正常发放。
      </el-alert>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="terminal_sn" label="终端SN" width="150" />
        <el-table-column prop="stage_reward_id" label="阶段奖励ID" width="120" />
        <el-table-column label="总比例" width="100">
          <template #default="{ row }">
            <el-tag type="danger">{{ (row.total_rate * 100).toFixed(1) }}%</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.resolved ? 'success' : 'danger'">
              {{ row.resolved ? '已解决' : '待处理' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resolved_by" label="处理人" width="100" />
        <el-table-column label="处理时间" width="170">
          <template #default="{ row }">
            {{ row.resolved_at ? formatDateTime(row.resolved_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-popconfirm
              v-if="!row.resolved"
              title="确定标记为已解决吗？"
              @confirm="handleResolve(row)"
            >
              <template #reference>
                <el-button type="primary" link>解决</el-button>
              </template>
            </el-popconfirm>
            <span v-else class="text-muted">已处理</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getOverflowLogs, resolveOverflowLog, type RewardOverflowLog } from '@/api/reward'

const loading = ref(false)
const tableData = ref<RewardOverflowLog[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
})

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getOverflowLogs({
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

const handleResolve = async (row: RewardOverflowLog) => {
  try {
    await resolveOverflowLog(row.id)
    ElMessage.success('已标记为解决')
    loadData()
  } catch (err) {
    ElMessage.error('操作失败')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.pagination-wrapper {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.text-muted {
  color: #909399;
}
</style>
