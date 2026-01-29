<template>
  <div class="poster-category">
    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="fetchData" @reset="fetchData">
      <template #extra>
        <el-button type="primary" :icon="Plus" @click="handleCreate">新增分类</el-button>
      </template>
    </SearchForm>

    <!-- 数据表格 -->
    <ProTable
      :data="tableData"
      :loading="loading"
      :total="tableData.length"
      :show-pagination="false"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="分类名称" min-width="150" />
      <el-table-column prop="sort_order" label="排序" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="poster_count" label="海报数量" width="100" />
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDateTime(row.created_at) }}
        </template>
      </el-table-column>

      <template #action="{ row }">
        <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
        <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
      </template>
    </ProTable>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingId ? '编辑分类' : '新增分类'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入分类名称" maxlength="50" show-word-limit />
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="9999" />
          <span class="form-tip">数值越大越靠前</span>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import {
  getPosterCategories,
  createPosterCategory,
  updatePosterCategory,
  deletePosterCategory
} from '@/api/poster'
import type { PosterCategory, PosterCategoryCreateRequest } from '@/types/poster'
import { formatDateTime } from '@/utils/format'

const searchForm = reactive({})
const loading = ref(false)
const tableData = ref<PosterCategory[]>([])
const dialogVisible = ref(false)
const editingId = ref<number | null>(null)
const submitLoading = ref(false)

const formRef = ref<FormInstance>()
const formData = reactive<PosterCategoryCreateRequest>({
  name: '',
  sort_order: 0,
  status: 1
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { max: 50, message: '名称最多50个字符', trigger: 'blur' }
  ]
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getPosterCategories()
    tableData.value = res || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 新增
const handleCreate = () => {
  editingId.value = null
  formData.name = ''
  formData.sort_order = 0
  formData.status = 1
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: PosterCategory) => {
  editingId.value = row.id
  formData.name = row.name
  formData.sort_order = row.sort_order
  formData.status = row.status
  dialogVisible.value = true
}

// 删除
const handleDelete = async (row: PosterCategory) => {
  try {
    await ElMessageBox.confirm(`确定要删除分类"${row.name}"吗？删除后该分类下的海报将无法查看。`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deletePosterCategory(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 提交
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitLoading.value = true

    if (editingId.value) {
      await updatePosterCategory(editingId.value, formData)
      ElMessage.success('保存成功')
    } else {
      await createPosterCategory(formData)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (error: any) {
    if (error !== false) {
      ElMessage.error(editingId.value ? '保存失败' : '创建失败')
    }
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.poster-category {
  padding: 0;
}

.form-tip {
  color: #909399;
  font-size: 12px;
  margin-left: 10px;
}
</style>
