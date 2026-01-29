<template>
  <div class="agent-list-view">
    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="通道">
        <ChannelSelect v-model="searchForm.channel_id" style="width: 150px" />
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="searchForm.status" placeholder="请选择状态" clearable style="width: 100px">
          <el-option label="正常" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
      </el-form-item>
      <el-form-item label="关键词">
        <el-input
          v-model="searchForm.keyword"
          placeholder="代理商名称/手机号/编号"
          clearable
          style="width: 180px"
        />
      </el-form-item>
      <template #extra>
        <el-button type="primary" :icon="Plus" @click="handleCreate">新增代理商</el-button>
      </template>
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
      <el-table-column prop="agent_no" label="代理商编号" width="140" />
      <el-table-column prop="agent_name" label="代理商名称" min-width="120" />
      <el-table-column prop="contact_phone" label="手机号" width="130" />
      <el-table-column prop="level" label="层级" width="80" align="center">
        <template #default="{ row }">
          <el-tag size="small">{{ row.level }}级</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="direct_agent_count" label="直属代理" width="90" align="center" />
      <el-table-column prop="direct_merchant_count" label="直属商户" width="90" align="center" />
      <el-table-column prop="status" label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status_name || (row.status === 1 ? '正常' : '禁用') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="register_time" label="注册时间" width="170" />

      <template #action="{ row }">
        <el-button type="primary" link @click="handleView(row)">详情</el-button>
        <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
        <el-button
          :type="row.status === 1 ? 'danger' : 'success'"
          link
          @click="handleToggleStatus(row)"
        >
          {{ row.status === 1 ? '禁用' : '启用' }}
        </el-button>
      </template>
    </ProTable>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import ChannelSelect from '@/components/Common/ChannelSelect.vue'
import { getAgents, updateAgentStatus } from '@/api/agent'
import type { Agent } from '@/types'

const router = useRouter()

// 搜索表单
const searchForm = reactive({
  channel_id: undefined as number | undefined,
  status: undefined as number | undefined,
  keyword: '',
})

// 表格数据
const tableData = ref<Agent[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getAgents({
      ...searchForm,
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('Fetch agents error:', error)
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
function handleView(row: Agent) {
  router.push(`/agents/${row.id}`)
}

// 新增代理商
function handleCreate() {
  router.push('/agents/create')
}

// 编辑
function handleEdit(row: Agent) {
  router.push(`/agents/${row.id}/edit`)
}

// 切换状态
async function handleToggleStatus(row: Agent) {
  const newStatus = row.status === 1 ? 0 : 1
  const action = newStatus === 1 ? '启用' : '禁用'

  try {
    await ElMessageBox.confirm(`确定要${action}代理商 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await updateAgentStatus(row.id, newStatus)
    ElMessage.success(`${action}成功`)
    fetchData()
  } catch (error) {
    // 用户取消或请求失败
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.agent-list-view {
  padding: 0;
}
</style>
