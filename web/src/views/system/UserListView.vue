<template>
  <div class="user-list-view">
    <PageHeader title="系统管理" sub-title="用户管理">
      <template #extra>
        <el-button type="primary" :icon="Plus" @click="handleCreate">新建用户</el-button>
      </template>
    </PageHeader>

    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="关键词">
        <el-input v-model="searchForm.keyword" placeholder="用户名/昵称/手机号" clearable />
      </el-form-item>
      <el-form-item label="角色">
        <el-select v-model="searchForm.role" placeholder="请选择角色" clearable>
          <el-option
            v-for="(config, key) in USER_ROLE_CONFIG"
            :key="key"
            :label="config.label"
            :value="key"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
          <el-option label="正常" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
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
      <el-table-column prop="username" label="用户名" width="120" />
      <el-table-column prop="nickname" label="昵称" width="120" />
      <el-table-column prop="phone" label="手机号" width="130" />
      <el-table-column prop="email" label="邮箱" min-width="180" />
      <el-table-column prop="role" label="角色" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getRoleTagType(row.role)" size="small">
            {{ getRoleLabel(row.role) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-switch
            v-model="row.status"
            :active-value="1"
            :inactive-value="0"
            @change="handleStatusChange(row)"
          />
        </template>
      </el-table-column>
      <el-table-column prop="last_login_at" label="最后登录" width="170">
        <template #default="{ row }">
          <div v-if="row.last_login_at">
            <div>{{ row.last_login_at }}</div>
            <div class="login-ip">{{ row.last_login_ip }}</div>
          </div>
          <span v-else class="text-placeholder">未登录</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="170" />

      <template #action="{ row }">
        <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
        <el-button type="warning" link @click="handleResetPassword(row)">重置密码</el-button>
        <el-button
          type="danger"
          link
          :disabled="row.role === 'admin'"
          @click="handleDelete(row)"
        >
          删除
        </el-button>
      </template>
    </ProTable>

    <!-- 新建/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑用户' : '新建用户'"
      width="500px"
      @closed="handleDialogClosed"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="form.username"
            placeholder="请输入用户名"
            :disabled="isEdit"
          />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" placeholder="请选择角色" style="width: 100%">
            <el-option
              v-for="(config, key) in USER_ROLE_CONFIG"
              :key="key"
              :label="config.label"
              :value="key"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>

    <!-- 重置密码弹窗 -->
    <el-dialog v-model="resetPasswordDialogVisible" title="重置密码" width="400px">
      <el-form ref="resetFormRef" :model="resetForm" :rules="resetRules" label-width="80px">
        <el-form-item label="新密码" prop="password">
          <el-input
            v-model="resetForm.password"
            type="password"
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="resetForm.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetPasswordDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="resetting" @click="handleSubmitResetPassword">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import { getUsers, createUser, updateUser, deleteUser, resetPassword, toggleUserStatus } from '@/api/system'
import type { SystemUser, UserRole } from '@/types/system'
import { USER_ROLE_CONFIG } from '@/types/system'

// 搜索表单
const searchForm = reactive({
  keyword: '',
  role: undefined as UserRole | undefined,
  status: undefined as number | undefined,
})

// 表格数据
const tableData = ref<SystemUser[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 弹窗
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const currentUser = ref<SystemUser | null>(null)

// 表单
const form = reactive({
  username: '',
  password: '',
  nickname: '',
  phone: '',
  email: '',
  role: 'readonly' as UserRole,
})

// 表单验证规则
const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度为3-20个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度为6-20个字符', trigger: 'blur' },
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' },
  ],
}

// 重置密码弹窗
const resetPasswordDialogVisible = ref(false)
const resetting = ref(false)
const resetFormRef = ref<FormInstance>()
const resetForm = reactive({
  password: '',
  confirmPassword: '',
})

// 重置密码验证规则
const resetRules: FormRules = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度为6-20个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== resetForm.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
}

// 获取角色标签类型
function getRoleTagType(role: UserRole) {
  const typeMap: Record<string, string> = {
    '#409eff': 'primary',
    '#67c23a': 'success',
    '#e6a23c': 'warning',
    '#909399': 'info',
  }
  return typeMap[USER_ROLE_CONFIG[role]?.color] || ''
}

// 获取角色名称
function getRoleLabel(role: UserRole) {
  return USER_ROLE_CONFIG[role]?.label || role
}

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getUsers({
      ...searchForm,
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('Fetch users error:', error)
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
  isEdit.value = false
  currentUser.value = null
  form.username = ''
  form.password = ''
  form.nickname = ''
  form.phone = ''
  form.email = ''
  form.role = 'readonly'
  dialogVisible.value = true
}

// 编辑
function handleEdit(row: SystemUser) {
  isEdit.value = true
  currentUser.value = row
  form.username = row.username
  form.password = ''
  form.nickname = row.nickname
  form.phone = row.phone
  form.email = row.email
  form.role = row.role
  dialogVisible.value = true
}

// 弹窗关闭
function handleDialogClosed() {
  formRef.value?.resetFields()
}

// 提交表单
async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value && currentUser.value) {
      await updateUser(currentUser.value.id, {
        nickname: form.nickname,
        phone: form.phone,
        email: form.email,
        role: form.role,
      })
      ElMessage.success('更新成功')
    } else {
      await createUser({
        username: form.username,
        password: form.password,
        nickname: form.nickname,
        phone: form.phone,
        email: form.email,
        role: form.role,
      })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('Submit error:', error)
  } finally {
    submitting.value = false
  }
}

// 状态变更
async function handleStatusChange(row: SystemUser) {
  try {
    await toggleUserStatus(row.id, row.status)
    ElMessage.success(row.status === 1 ? '已启用' : '已禁用')
  } catch (error) {
    // 恢复状态
    row.status = row.status === 1 ? 0 : 1
    console.error('Toggle status error:', error)
  }
}

// 重置密码
function handleResetPassword(row: SystemUser) {
  currentUser.value = row
  resetForm.password = ''
  resetForm.confirmPassword = ''
  resetPasswordDialogVisible.value = true
}

// 提交重置密码
async function handleSubmitResetPassword() {
  const valid = await resetFormRef.value?.validate().catch(() => false)
  if (!valid) return

  if (!currentUser.value) return

  resetting.value = true
  try {
    await resetPassword(currentUser.value.id, resetForm.password)
    ElMessage.success('密码重置成功')
    resetPasswordDialogVisible.value = false
  } catch (error) {
    console.error('Reset password error:', error)
  } finally {
    resetting.value = false
  }
}

// 删除用户
async function handleDelete(row: SystemUser) {
  try {
    await ElMessageBox.confirm(`确定要删除用户 "${row.nickname}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await deleteUser(row.id)
    ElMessage.success('删除成功')
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
.user-list-view {
  padding: 0;
}

.login-ip {
  font-size: 12px;
  color: $text-placeholder;
}

.text-placeholder {
  color: $text-placeholder;
}
</style>
