<template>
  <div class="login-container">
    <div class="login-card">
      <!-- Logo和标题 -->
      <div class="login-header">
        <img src="/vite.svg" alt="Logo" class="logo" />
        <h1 class="title">享收付管理系统</h1>
        <p class="subtitle">Agent Profit Sharing Management System</p>
      </div>

      <!-- 登录表单 -->
      <el-form
        ref="formRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        size="large"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            :prefix-icon="User"
            clearable
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <div class="form-options">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <el-link type="primary" underline="never">忘记密码?</el-link>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            class="login-btn"
            :loading="loading"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登 录' }}
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 底部信息 -->
      <div class="login-footer">
        <p>© 2024 享收付. All rights reserved.</p>
      </div>
    </div>

    <!-- 背景装饰 -->
    <div class="login-bg">
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
      <div class="bg-shape shape-3"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 表单引用
const formRef = ref<FormInstance>()

// 登录表单数据
const loginForm = reactive({
  username: '',
  password: '',
})

// 记住我
const rememberMe = ref(false)

// 加载状态
const loading = ref(false)

// 表单验证规则
const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在3-20个字符之间', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在6-20个字符之间', trigger: 'blur' },
  ],
}

// 处理登录
async function handleLogin() {
  if (!formRef.value) return

  try {
    // 表单验证
    await formRef.value.validate()

    loading.value = true

    // 调用登录接口
    await userStore.login(loginForm.username, loginForm.password)

    ElMessage.success('登录成功')

    // 跳转到目标页面
    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect)
  } catch (error: any) {
    console.error('Login error:', error)
    // API层已经显示了错误消息
  } finally {
    loading.value = false
  }
}

// 页面加载时检查是否已登录
onMounted(() => {
  if (userStore.isLoggedIn) {
    router.push('/dashboard')
  }
})
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

.login-card {
  width: 420px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  z-index: 10;
  backdrop-filter: blur(10px);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;

  .logo {
    width: 64px;
    height: 64px;
    margin-bottom: 16px;
  }

  .title {
    font-size: 24px;
    font-weight: 600;
    color: $text-primary;
    margin-bottom: 8px;
  }

  .subtitle {
    font-size: 14px;
    color: $text-secondary;
  }
}

.login-form {
  .form-options {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .login-btn {
    width: 100%;
    height: 44px;
    font-size: 16px;
    letter-spacing: 2px;
  }
}

.login-footer {
  text-align: center;
  margin-top: 24px;

  p {
    font-size: 12px;
    color: $text-secondary;
  }
}

// 背景装饰
.login-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  background: #ffffff;
}

.shape-1 {
  width: 400px;
  height: 400px;
  top: -100px;
  left: -100px;
  animation: float 8s ease-in-out infinite;
}

.shape-2 {
  width: 300px;
  height: 300px;
  bottom: -50px;
  right: -50px;
  animation: float 6s ease-in-out infinite reverse;
}

.shape-3 {
  width: 200px;
  height: 200px;
  top: 50%;
  right: 20%;
  animation: float 7s ease-in-out infinite;
}

@keyframes float {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-20px);
  }
}

// 响应式
@media (max-width: 480px) {
  .login-card {
    width: 90%;
    padding: 30px 20px;
  }

  .login-header {
    .title {
      font-size: 20px;
    }
  }
}
</style>
