<template>
  <div class="terminal-type-list">
    <!-- 搜索和操作区 -->
    <el-card class="search-card" shadow="never">
      <el-form :inline="true" :model="queryParams" class="search-form">
        <el-form-item label="所属通道">
          <el-select
            v-model="queryParams.channel_id"
            placeholder="全部通道"
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="channel in channels"
              :key="channel.id"
              :label="channel.channel_name"
              :value="channel.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input
            v-model="queryParams.keyword"
            placeholder="品牌/型号"
            clearable
            style="width: 180px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
      <div class="action-buttons">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          新增终端类型
        </el-button>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card class="table-card" shadow="never">
      <el-table
        v-loading="loading"
        :data="tableData"
        stripe
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="channel_name" label="所属通道" width="120" />
        <el-table-column prop="brand_name" label="品牌名称" width="120" />
        <el-table-column prop="brand_code" label="品牌编码" width="120" />
        <el-table-column prop="model_name" label="型号名称" width="120">
          <template #default="{ row }">
            {{ row.model_name || row.model_code }}
          </template>
        </el-table-column>
        <el-table-column prop="model_code" label="型号编码" width="120" />
        <el-table-column prop="full_name" label="完整名称" min-width="150" />
        <el-table-column prop="description" label="描述" min-width="150" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button
              :type="row.status === 1 ? 'warning' : 'success'"
              link
              size="small"
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-popconfirm
              title="确定删除该终端类型吗？"
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button type="danger" link size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'create' ? '新增终端类型' : '编辑终端类型'"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="所属通道" prop="channel_id">
          <el-select
            v-model="formData.channel_id"
            placeholder="请选择通道"
            :disabled="dialogType === 'edit'"
            style="width: 100%"
          >
            <el-option
              v-for="channel in channels"
              :key="channel.id"
              :label="channel.channel_name"
              :value="channel.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="品牌编码" prop="brand_code">
          <el-input v-model="formData.brand_code" placeholder="如：NEWLAND" />
        </el-form-item>
        <el-form-item label="品牌名称" prop="brand_name">
          <el-input v-model="formData.brand_name" placeholder="如：新大陆" />
        </el-form-item>
        <el-form-item label="型号编码" prop="model_code">
          <el-input v-model="formData.model_code" placeholder="如：ME31" />
        </el-form-item>
        <el-form-item label="型号名称" prop="model_name">
          <el-input v-model="formData.model_name" placeholder="可选，如：智能POS" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="终端类型描述（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { TerminalType, TerminalTypeQueryParams, CreateTerminalTypeParams } from '@/types/terminalType'
import {
  getTerminalTypes,
  createTerminalType,
  updateTerminalType,
  updateTerminalTypeStatus,
  deleteTerminalType,
} from '@/api/terminalType'
import { getChannels } from '@/api/agent-channel'

// 通道列表
interface Channel {
  id: number
  channel_name: string
  channel_code: string
}
const channels = ref<Channel[]>([])

// 查询参数
const queryParams = reactive<TerminalTypeQueryParams>({
  channel_id: undefined,
  status: undefined,
  keyword: '',
  page: 1,
  page_size: 20,
})

// 表格数据
const loading = ref(false)
const tableData = ref<TerminalType[]>([])
const total = ref(0)

// 对话框
const dialogVisible = ref(false)
const dialogType = ref<'create' | 'edit'>('create')
const submitLoading = ref(false)
const formRef = ref<FormInstance>()
const currentEditId = ref<number>(0)

// 表单数据
const formData = reactive<CreateTerminalTypeParams>({
  channel_id: 0,
  brand_code: '',
  brand_name: '',
  model_code: '',
  model_name: '',
  description: '',
})

// 表单验证规则
const formRules: FormRules = {
  channel_id: [{ required: true, message: '请选择所属通道', trigger: 'change' }],
  brand_code: [{ required: true, message: '请输入品牌编码', trigger: 'blur' }],
  brand_name: [{ required: true, message: '请输入品牌名称', trigger: 'blur' }],
  model_code: [{ required: true, message: '请输入型号编码', trigger: 'blur' }],
}

// 获取通道列表
const fetchChannels = async () => {
  try {
    const res = await getChannels()
    channels.value = res.list || []
  } catch (error) {
    console.error('获取通道列表失败:', error)
  }
}

// 获取终端类型列表
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getTerminalTypes(queryParams)
    tableData.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取终端类型列表失败:', error)
    ElMessage.error('获取终端类型列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  queryParams.page = 1
  fetchData()
}

// 重置
const handleReset = () => {
  queryParams.channel_id = undefined
  queryParams.status = undefined
  queryParams.keyword = ''
  queryParams.page = 1
  fetchData()
}

// 分页
const handleSizeChange = (size: number) => {
  queryParams.page_size = size
  fetchData()
}

const handleCurrentChange = (page: number) => {
  queryParams.page = page
  fetchData()
}

// 重置表单
const resetForm = () => {
  formData.channel_id = 0
  formData.brand_code = ''
  formData.brand_name = ''
  formData.model_code = ''
  formData.model_name = ''
  formData.description = ''
}

// 新增
const handleCreate = () => {
  dialogType.value = 'create'
  resetForm()
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: TerminalType) => {
  dialogType.value = 'edit'
  currentEditId.value = row.id
  formData.channel_id = row.channel_id
  formData.brand_code = row.brand_code
  formData.brand_name = row.brand_name
  formData.model_code = row.model_code
  formData.model_name = row.model_name
  formData.description = row.description
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      if (dialogType.value === 'create') {
        await createTerminalType(formData)
        ElMessage.success('创建成功')
      } else {
        await updateTerminalType(currentEditId.value, formData)
        ElMessage.success('更新成功')
      }
      dialogVisible.value = false
      fetchData()
    } catch (error: any) {
      ElMessage.error(error.message || '操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

// 切换状态
const handleToggleStatus = async (row: TerminalType) => {
  const newStatus = row.status === 1 ? 0 : 1
  try {
    await updateTerminalTypeStatus(row.id, newStatus)
    ElMessage.success(newStatus === 1 ? '启用成功' : '禁用成功')
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

// 删除
const handleDelete = async (row: TerminalType) => {
  try {
    await deleteTerminalType(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
  }
}

// 初始化
onMounted(() => {
  fetchChannels()
  fetchData()
})
</script>

<style scoped lang="scss">
.terminal-type-list {
  padding: 20px;
}

.search-card {
  margin-bottom: 20px;

  .search-form {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }

  .action-buttons {
    margin-top: 16px;
    display: flex;
    gap: 10px;
  }
}

.table-card {
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
