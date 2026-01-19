<template>
  <div class="policy-list-view">
    <PageHeader title="政策管理" sub-title="政策模板">
      <template #extra>
        <el-button type="primary" :icon="Plus" @click="handleCreate">新建模板</el-button>
      </template>
    </PageHeader>

    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="通道">
        <ChannelSelect v-model="searchForm.channel_id" />
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
          <el-option label="启用" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
      </el-form-item>
      <el-form-item label="关键词">
        <el-input v-model="searchForm.keyword" placeholder="模板名称" clearable />
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
      <el-table-column prop="name" label="模板名称" min-width="150" />
      <el-table-column prop="channel_name" label="所属通道" width="120" />
      <el-table-column prop="credit_rate" label="贷记卡费率" width="120" align="center">
        <template #default="{ row }">
          {{ (row.credit_rate * 100).toFixed(2) }}%
        </template>
      </el-table-column>
      <el-table-column prop="debit_rate" label="借记卡费率" width="120" align="center">
        <template #default="{ row }">
          {{ (row.debit_rate * 100).toFixed(2) }}%
        </template>
      </el-table-column>
      <el-table-column prop="debit_cap" label="借记卡封顶" width="100" align="center">
        <template #default="{ row }">
          ¥{{ row.debit_cap }}
        </template>
      </el-table-column>
      <el-table-column prop="is_default" label="默认" width="80" align="center">
        <template #default="{ row }">
          <el-icon v-if="row.is_default" class="default-icon"><StarFilled /></el-icon>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>

      <template #action="{ row }">
        <el-button type="primary" link @click="handleView(row)">详情</el-button>
        <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
        <el-button type="primary" link @click="handleCopy(row)">复制</el-button>
        <el-button
          v-if="!row.is_default"
          type="success"
          link
          @click="handleSetDefault(row)"
        >
          设为默认
        </el-button>
      </template>
    </ProTable>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, StarFilled } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import ChannelSelect from '@/components/Common/ChannelSelect.vue'
import { getPolicyTemplates, copyPolicyTemplate, setDefaultTemplate } from '@/api/policy'
import type { PolicyTemplate } from '@/types'

const router = useRouter()

// 搜索表单
const searchForm = reactive({
  channel_id: undefined as number | undefined,
  status: undefined as number | undefined,
  keyword: '',
})

// 表格数据
const tableData = ref<PolicyTemplate[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getPolicyTemplates({
      ...searchForm,
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('Fetch policy templates error:', error)
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

// 新建
function handleCreate() {
  router.push('/policies/create')
}

// 查看详情
function handleView(row: PolicyTemplate) {
  router.push(`/policies/${row.id}`)
}

// 编辑
function handleEdit(row: PolicyTemplate) {
  router.push(`/policies/${row.id}/edit`)
}

// 复制
async function handleCopy(row: PolicyTemplate) {
  try {
    const { value } = await ElMessageBox.prompt('请输入新模板名称', '复制模板', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputValue: `${row.name}_副本`,
      inputPattern: /\S+/,
      inputErrorMessage: '请输入模板名称',
    })

    await copyPolicyTemplate(row.id, value)
    ElMessage.success('复制成功')
    fetchData()
  } catch (error) {
    // 用户取消
  }
}

// 设为默认
async function handleSetDefault(row: PolicyTemplate) {
  try {
    await ElMessageBox.confirm(`确定要将 "${row.name}" 设为默认模板吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await setDefaultTemplate(row.id)
    ElMessage.success('设置成功')
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
.policy-list-view {
  padding: 0;
}

.default-icon {
  color: $warning-color;
  font-size: 16px;
}
</style>
