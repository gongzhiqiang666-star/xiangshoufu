<template>
  <div class="login-container">
    <!-- 左侧品牌展示区 -->
    <div class="brand-section">
      <!-- 背景装饰 -->
      <div class="bg-decoration">
        <div class="grid-pattern"></div>
        <div class="glow-orb glow-1"></div>
        <div class="glow-orb glow-2"></div>
      </div>

      <div class="brand-content">
        <!-- Logo -->
        <div class="brand-logo">
          <svg viewBox="0 0 56 56" class="logo-icon">
            <defs>
              <linearGradient id="logoGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stop-color="#F59E0B" />
                <stop offset="100%" stop-color="#FBBF24" />
              </linearGradient>
            </defs>
            <rect x="4" y="4" width="48" height="48" rx="12" fill="url(#logoGrad)" />
            <path d="M18 28 L26 36 L38 20" stroke="#0F172A" stroke-width="4" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>

        <h1 class="brand-title">享收付</h1>
        <p class="brand-subtitle">XiangShouFu Payment Platform</p>

        <!-- 特性列表 -->
        <div class="features">
          <div class="feature-item">
            <svg class="feature-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
            <div class="feature-text">
              <span class="feature-title">银行级安全</span>
              <span class="feature-desc">256位SSL加密传输</span>
            </div>
          </div>

          <div class="feature-item">
            <svg class="feature-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <polyline points="12,6 12,12 16,14"/>
            </svg>
            <div class="feature-text">
              <span class="feature-title">实时结算</span>
              <span class="feature-desc">毫秒级分润计算</span>
            </div>
          </div>

          <div class="feature-item">
            <svg class="feature-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
            </svg>
            <div class="feature-text">
              <span class="feature-title">多级代理</span>
              <span class="feature-desc">智能分润体系</span>
            </div>
          </div>
        </div>

      </div>
    </div>

    <!-- 右侧登录区域 -->
    <div class="login-section">
      <div class="login-card">
        <!-- 欢迎文本 -->
        <div class="welcome-header">
          <h2 class="welcome-title">欢迎登录</h2>
          <p class="welcome-subtitle">请输入您的账户信息</p>
        </div>

        <!-- 登录表单 -->
        <el-form
          ref="formRef"
          :model="loginForm"
          :rules="loginRules"
          class="login-form"
          @keyup.enter="handleLogin"
        >
          <div class="form-field">
            <label class="field-label" for="username">
              <svg class="label-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
              账户名称
            </label>
            <el-form-item prop="username">
              <el-input
                id="username"
                v-model="loginForm.username"
                placeholder="请输入用户名"
                size="large"
                clearable
                class="custom-input"
              />
            </el-form-item>
          </div>

          <div class="form-field">
            <label class="field-label" for="password">
              <svg class="label-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              登录密码
            </label>
            <el-form-item prop="password">
              <el-input
                id="password"
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                size="large"
                show-password
                class="custom-input"
              />
            </el-form-item>
          </div>

          <div class="form-options">
            <el-checkbox v-model="rememberMe" class="remember-checkbox">
              记住登录状态
            </el-checkbox>
          </div>

          <el-button
            type="primary"
            class="login-btn"
            :loading="loading"
            @click="handleLogin"
          >
            <span v-if="!loading">安全登录</span>
            <span v-else>正在验证...</span>
          </el-button>
        </el-form>

        <!-- 安全提示 -->
        <div class="security-footer">
          <div class="security-badge">
            <svg class="security-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
              <path d="M9 12l2 2 4-4"/>
            </svg>
            <span>SSL安全加密</span>
          </div>
          <div class="security-badge">
            <svg class="security-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
            </svg>
            <span>隐私保护</span>
          </div>
        </div>
      </div>

      <!-- 底部版权 -->
      <div class="copyright">
        <p>© 2024 享收付 XiangShouFu. All rights reserved.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
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
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
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
// ============================================
// Design System: UI/UX Pro Max - Fintech Theme
// Colors: Gold (#F59E0B) + Deep Blue (#0F172A)
// Typography: Poppins / Open Sans
// ============================================

// 设计令牌
$color-primary: #F59E0B;
$color-primary-light: #FBBF24;
$color-primary-dark: #D97706;
$color-background: #0F172A;
$color-background-light: #1E293B;
$color-surface: #334155;
$color-text-primary: #F8FAFC;
$color-text-secondary: #94A3B8;
$color-text-muted: #64748B;
$color-border: #475569;
$color-success: #10B981;

// 布局
.login-container {
  min-height: 100vh;
  display: flex;
  font-family: 'Open Sans', -apple-system, BlinkMacSystemFont, sans-serif;
}

// ==========================================
// 左侧品牌区域
// ==========================================
.brand-section {
  flex: 1.2;
  background: $color-background;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  padding: 60px;
}

// 背景装饰
.bg-decoration {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.grid-pattern {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(248, 250, 252, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(248, 250, 252, 0.03) 1px, transparent 1px);
  background-size: 60px 60px;
}

.glow-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
}

.glow-1 {
  width: 400px;
  height: 400px;
  background: $color-primary;
  top: -100px;
  right: -100px;
  animation: float 8s ease-in-out infinite;
}

.glow-2 {
  width: 300px;
  height: 300px;
  background: $color-primary-dark;
  bottom: -50px;
  left: -50px;
  animation: float 10s ease-in-out infinite reverse;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(20px, 20px); }
}

// 品牌内容
.brand-content {
  position: relative;
  z-index: 2;
  max-width: 480px;
}

.brand-logo {
  margin-bottom: 32px;

  .logo-icon {
    width: 72px;
    height: 72px;
    filter: drop-shadow(0 8px 24px rgba(245, 158, 11, 0.3));
  }
}

.brand-title {
  font-family: 'Poppins', sans-serif;
  font-size: 48px;
  font-weight: 600;
  color: $color-text-primary;
  letter-spacing: 4px;
  margin-bottom: 8px;
}

.brand-subtitle {
  font-size: 16px;
  color: $color-text-secondary;
  letter-spacing: 2px;
  margin-bottom: 48px;
}

// 特性列表
.features {
  display: flex;
  flex-direction: column;
  gap: 24px;
  margin-bottom: 48px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  transition: all 0.3s ease;
  cursor: default;

  &:hover {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(245, 158, 11, 0.3);
    transform: translateX(8px);
  }
}

.feature-icon {
  width: 40px;
  height: 40px;
  color: $color-primary;
  flex-shrink: 0;
}

.feature-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.feature-title {
  font-size: 15px;
  font-weight: 600;
  color: $color-text-primary;
}

.feature-desc {
  font-size: 13px;
  color: $color-text-secondary;
}

// ==========================================
// 右侧登录区域
// ==========================================
.login-section {
  flex: 0.8;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 60px;
  background: linear-gradient(180deg, #F8FAFC 0%, #F1F5F9 100%);
  min-width: 480px;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: #FFFFFF;
  border-radius: 20px;
  padding: 48px 40px;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.05),
    0 10px 15px -3px rgba(0, 0, 0, 0.08),
    0 20px 25px -5px rgba(0, 0, 0, 0.05);
}

// 欢迎文本
.welcome-header {
  text-align: center;
  margin-bottom: 36px;
}

.welcome-title {
  font-family: 'Poppins', sans-serif;
  font-size: 28px;
  font-weight: 600;
  color: $color-background;
  margin-bottom: 8px;
}

.welcome-subtitle {
  font-size: 15px;
  color: $color-text-muted;
}

// 表单字段
.form-field {
  margin-bottom: 24px;
}

.field-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: $color-background;
  margin-bottom: 10px;
}

.label-icon {
  width: 18px;
  height: 18px;
  color: $color-text-muted;
}

.custom-input {
  :deep(.el-input__wrapper) {
    background: #F8FAFC;
    border: 2px solid #E2E8F0;
    border-radius: 12px;
    padding: 4px 16px;
    box-shadow: none;
    transition: all 0.2s ease;

    &:hover {
      border-color: #CBD5E1;
    }

    &.is-focus {
      border-color: $color-primary;
      background: #FFFFFF;
      box-shadow: 0 0 0 4px rgba(245, 158, 11, 0.1);
    }
  }

  :deep(.el-input__inner) {
    font-size: 15px;
    color: $color-background;
    height: 48px;

    &::placeholder {
      color: #94A3B8;
    }
  }
}

// 表单选项
.form-options {
  margin-bottom: 28px;

  .remember-checkbox {
    :deep(.el-checkbox__label) {
      color: $color-text-muted;
      font-size: 14px;
    }

    :deep(.el-checkbox__inner) {
      border-radius: 4px;
      border-color: #CBD5E1;
      width: 18px;
      height: 18px;
    }

    :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
      background-color: $color-primary;
      border-color: $color-primary;
    }
  }
}

// 登录按钮
.login-btn {
  width: 100%;
  height: 52px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 2px;
  border-radius: 12px;
  background: linear-gradient(135deg, $color-primary 0%, $color-primary-dark 100%);
  border: none;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(245, 158, 11, 0.4);
  }

  &:active {
    transform: translateY(0);
  }

  &.is-loading {
    background: $color-primary-dark;
  }
}

// 安全提示
.security-footer {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #E2E8F0;
}

.security-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: $color-text-muted;
}

.security-icon {
  width: 16px;
  height: 16px;
  color: $color-success;
}

// 版权信息
.copyright {
  margin-top: 32px;

  p {
    font-size: 12px;
    color: #94A3B8;
    letter-spacing: 0.5px;
  }
}

// ==========================================
// 响应式设计
// ==========================================
@media (max-width: 1200px) {
  .brand-section {
    flex: 1;
    padding: 40px;
  }

  .login-section {
    min-width: 420px;
    padding: 40px;
  }
}

@media (max-width: 1024px) {
  .login-container {
    flex-direction: column;
  }

  .brand-section {
    flex: none;
    min-height: 320px;
    padding: 40px;

    .brand-title {
      font-size: 36px;
    }

    .features {
      display: none;
    }

    .stats {
      justify-content: center;
    }
  }

  .login-section {
    flex: 1;
    min-width: auto;
    padding: 40px 24px;
  }
}

@media (max-width: 480px) {
  .brand-section {
    min-height: 240px;
    padding: 32px 24px;

    .brand-logo .logo-icon {
      width: 56px;
      height: 56px;
    }

    .brand-title {
      font-size: 28px;
      letter-spacing: 2px;
    }

    .brand-subtitle {
      font-size: 14px;
      margin-bottom: 24px;
    }

    .stats {
      gap: 16px;
    }

    .stat-value {
      font-size: 22px;
    }
  }

  .login-card {
    padding: 32px 24px;
    border-radius: 16px;
  }

  .welcome-title {
    font-size: 24px;
  }

  .login-btn {
    height: 48px;
    letter-spacing: 1px;
  }

  .security-footer {
    flex-direction: column;
    align-items: center;
    gap: 12px;
  }
}

// 减少动画（无障碍）
@media (prefers-reduced-motion: reduce) {
  .glow-orb {
    animation: none;
  }

  .feature-item,
  .login-btn {
    transition: none;
  }
}
</style>
