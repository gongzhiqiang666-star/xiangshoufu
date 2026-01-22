<template>
  <div class="change-password-view">
    <el-card class="password-card">
      <template #header>
        <div class="card-header">
          <span class="title">修改密码</span>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        class="password-form"
      >
        <el-form-item label="当前密码" prop="old_password">
          <el-input
            v-model="form.old_password"
            type="password"
            placeholder="请输入当前密码"
            show-password
          />
        </el-form-item>

        <el-form-item label="新密码" prop="new_password">
          <el-input
            v-model="form.new_password"
            type="password"
            placeholder="请输入新密码（至少6位）"
            show-password
          />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirm_password">
          <el-input
            v-model="form.confirm_password"
            type="password"
            placeholder="请再次输入新密码"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            确认修改
          </el-button>
          <el-button @click="handleReset">重置</el-button>
          <el-button @click="$router.back()">返回</el-button>
        </el-form-item>
      </el-form>

      <div class="password-tips">
        <h4>密码设置建议：</h4>
        <ul>
          <li>密码长度至少6位</li>
          <li>建议包含数字、字母和特殊字符</li>
          <li>避免使用连续数字或简单组合</li>
          <li>定期更换密码以保障账户安全</li>
        </ul>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { changePassword } from '@/api/auth'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

// 验证确认密码
const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== form.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  old_password: [
    { required: true, message: '请输入当前密码', trigger: 'blur' },
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' },
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

async function handleSubmit() {
  const valid = await formRef.value?.validate()
  if (!valid) return

  loading.value = true
  try {
    await changePassword({
      old_password: form.old_password,
      new_password: form.new_password,
    })
    ElMessage.success('密码修改成功，请重新登录')
    // 修改密码成功后退出登录
    await userStore.logout()
    router.push('/login')
  } catch (error: any) {
    console.error('Failed to change password:', error)
    ElMessage.error(error.message || '密码修改失败，请检查当前密码是否正确')
  } finally {
    loading.value = false
  }
}

function handleReset() {
  formRef.value?.resetFields()
}
</script>

<style lang="scss" scoped>
.change-password-view {
  max-width: 600px;
  margin: 0 auto;
}

.password-card {
  .card-header {
    .title {
      font-size: 18px;
      font-weight: 500;
      color: $text-primary;
    }
  }
}

.password-form {
  max-width: 400px;
  margin-bottom: $spacing-xl;
}

.password-tips {
  padding: $spacing-lg;
  background: $bg-color;
  border-radius: $border-radius-sm;

  h4 {
    margin: 0 0 $spacing-md 0;
    font-size: 14px;
    color: $text-primary;
  }

  ul {
    margin: 0;
    padding-left: $spacing-lg;

    li {
      font-size: 13px;
      color: $text-secondary;
      line-height: 1.8;
    }
  }
}
</style>
