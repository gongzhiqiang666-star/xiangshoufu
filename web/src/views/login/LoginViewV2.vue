<template>
  <div class="login-v2">
    <!-- 全屏背景 -->
    <div class="fullscreen-bg">
      <div class="bg-overlay"></div>
      <div class="grid-pattern"></div>
      <div class="glow-sphere sphere-1"></div>
      <div class="glow-sphere sphere-2"></div>
      <div class="glow-sphere sphere-3"></div>
    </div>

    <!-- 主内容区 -->
    <div class="login-content">
      <!-- 左侧：品牌展示 -->
      <div class="brand-panel">
        <div class="brand-wrapper">
          <!-- Logo -->
          <div class="logo-container">
            <svg viewBox="0 0 80 80" class="logo-svg">
              <defs>
                <linearGradient id="goldGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#F59E0B" />
                  <stop offset="50%" stop-color="#FBBF24" />
                  <stop offset="100%" stop-color="#F59E0B" />
                </linearGradient>
                <filter id="glow">
                  <feGaussianBlur stdDeviation="3" result="coloredBlur"/>
                  <feMerge>
                    <feMergeNode in="coloredBlur"/>
                    <feMergeNode in="SourceGraphic"/>
                  </feMerge>
                </filter>
              </defs>
              <circle cx="40" cy="40" r="36" fill="none" stroke="url(#goldGrad)" stroke-width="2" filter="url(#glow)"/>
              <text x="40" y="52" text-anchor="middle" fill="url(#goldGrad)" font-size="36" font-weight="bold" filter="url(#glow)">¥</text>
            </svg>
          </div>

          <h1 class="brand-name">享收付</h1>
          <p class="brand-tagline">智能分润 · 财富增长</p>
          <p class="brand-slogan">XiangShouFu Payment Platform</p>

          <!-- 特性列表 -->
          <div class="features">
            <div class="feature-item">
              <div class="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                </svg>
              </div>
              <div class="feature-text">
                <span class="feature-title">银行级安全</span>
                <span class="feature-desc">256位SSL加密传输</span>
              </div>
            </div>

            <div class="feature-item">
              <div class="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <polyline points="12,6 12,12 16,14"/>
                </svg>
              </div>
              <div class="feature-text">
                <span class="feature-title">实时结算</span>
                <span class="feature-desc">毫秒级分润计算</span>
              </div>
            </div>

            <div class="feature-item">
              <div class="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                  <circle cx="9" cy="7" r="4"/>
                  <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                  <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                </svg>
              </div>
              <div class="feature-text">
                <span class="feature-title">多级代理</span>
                <span class="feature-desc">智能分润体系</span>
              </div>
            </div>
          </div>

        </div>
      </div>

      <!-- 右侧：登录表单 -->
      <div class="login-panel">
        <div class="login-card">
          <!-- 欢迎语 -->
          <div class="welcome-section">
            <h2 class="welcome-title">欢迎回来</h2>
            <p class="welcome-desc">登录开启您的财富之旅</p>
          </div>

          <!-- 登录表单 -->
          <el-form
            ref="formRef"
            :model="loginForm"
            :rules="loginRules"
            class="login-form"
            @keyup.enter="handleLogin"
          >
            <div class="form-group">
              <label class="form-label">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                  <circle cx="12" cy="7" r="4"/>
                </svg>
                账户名称
              </label>
              <el-form-item prop="username">
                <el-input
                  v-model="loginForm.username"
                  placeholder="请输入用户名"
                  size="large"
                  clearable
                  class="gold-input"
                >
                  <template #prefix>
                    <span class="input-prefix-line"></span>
                  </template>
                </el-input>
              </el-form-item>
            </div>

            <div class="form-group">
              <label class="form-label">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
                登录密码
              </label>
              <el-form-item prop="password">
                <el-input
                  v-model="loginForm.password"
                  type="password"
                  placeholder="请输入密码"
                  size="large"
                  show-password
                  class="gold-input"
                >
                  <template #prefix>
                    <span class="input-prefix-line"></span>
                  </template>
                </el-input>
              </el-form-item>
            </div>

            <div class="form-options">
              <el-checkbox v-model="rememberMe" class="gold-checkbox">
                记住登录状态
              </el-checkbox>
            </div>

            <el-button
              type="primary"
              class="login-button"
              :loading="loading"
              @click="handleLogin"
            >
              <template v-if="!loading">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="btn-icon">
                  <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/>
                  <polyline points="10 17 15 12 10 7"/>
                  <line x1="15" y1="12" x2="3" y2="12"/>
                </svg>
                立即登录
              </template>
              <template v-else>
                正在验证...
              </template>
            </el-button>
          </el-form>

          <!-- 安全徽章 -->
          <div class="security-badges">
            <div class="badge">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                <path d="M9 12l2 2 4-4"/>
              </svg>
              <span>SSL加密</span>
            </div>
            <div class="badge">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              <span>数据安全</span>
            </div>
            <div class="badge">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <polyline points="12 6 12 12 16 14"/>
              </svg>
              <span>7×24小时</span>
            </div>
          </div>
        </div>

        <!-- 底部版权 -->
        <div class="footer">
          <p>© 2024 享收付 XiangShouFu. All rights reserved.</p>
        </div>
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
    await formRef.value.validate()
    loading.value = true

    await userStore.login(loginForm.username, loginForm.password)
    ElMessage.success('登录成功，欢迎回来！')

    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect)
  } catch (error: any) {
    console.error('Login error:', error)
  } finally {
    loading.value = false
  }
}

// 页面加载
onMounted(() => {
  if (userStore.isLoggedIn) {
    router.push('/dashboard')
  }
})
</script>

<style lang="scss" scoped>
// ============================================
// 设计系统：Dark Mode + Gold Theme
// 主色：Gold #F59E0B
// 背景：全屏深色渐变
// ============================================

$gold-primary: #F59E0B;
$gold-light: #FBBF24;
$gold-dark: #D97706;
$purple-accent: #8B5CF6;
$green-success: #10B981;
$bg-dark: #0F172A;
$bg-darker: #020617;
$bg-card: #1E293B;
$bg-surface: #334155;
$text-primary: #F8FAFC;
$text-secondary: #94A3B8;
$text-muted: #64748B;

.login-v2 {
  min-height: 100vh;
  width: 100%;
  position: relative;
  overflow: hidden;
}

// ==========================================
// 全屏背景
// ==========================================
.fullscreen-bg {
  position: fixed;
  inset: 0;
  background: linear-gradient(135deg, $bg-darker 0%, $bg-dark 50%, #1a1f35 100%);
  z-index: 0;
}

.bg-overlay {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(ellipse at 30% 20%, rgba($gold-primary, 0.08) 0%, transparent 50%),
    radial-gradient(ellipse at 70% 80%, rgba($purple-accent, 0.06) 0%, transparent 50%);
}

.grid-pattern {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(248, 250, 252, 0.015) 1px, transparent 1px),
    linear-gradient(90deg, rgba(248, 250, 252, 0.015) 1px, transparent 1px);
  background-size: 100px 100px;
}

// 光晕球
.glow-sphere {
  position: absolute;
  border-radius: 50%;
  filter: blur(120px);
  pointer-events: none;

  &.sphere-1 {
    width: 600px;
    height: 600px;
    background: rgba($gold-primary, 0.12);
    top: -200px;
    right: -150px;
    animation: float-slow 20s ease-in-out infinite;
  }

  &.sphere-2 {
    width: 500px;
    height: 500px;
    background: rgba($purple-accent, 0.08);
    bottom: -200px;
    left: -150px;
    animation: float-slow 25s ease-in-out infinite reverse;
  }

  &.sphere-3 {
    width: 400px;
    height: 400px;
    background: rgba($gold-primary, 0.06);
    top: 40%;
    left: 30%;
    animation: float-slow 30s ease-in-out infinite;
  }
}

@keyframes float-slow {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(30px, 30px) scale(1.1); }
}

// ==========================================
// 主内容布局
// ==========================================
.login-content {
  position: relative;
  z-index: 1;
  display: flex;
  min-height: 100vh;
  width: 100%;
}

// ==========================================
// 左侧品牌面板
// ==========================================
.brand-panel {
  flex: 1.2;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px;
}

.brand-wrapper {
  max-width: 500px;
}

.logo-container {
  margin-bottom: 32px;

  .logo-svg {
    width: 100px;
    height: 100px;
  }
}

.brand-name {
  font-size: 56px;
  font-weight: 700;
  letter-spacing: 8px;
  margin-bottom: 16px;
  background: linear-gradient(135deg, $text-primary 0%, $gold-light 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.brand-tagline {
  font-size: 20px;
  color: $gold-primary;
  letter-spacing: 6px;
  margin-bottom: 8px;
}

.brand-slogan {
  font-size: 14px;
  color: $text-muted;
  letter-spacing: 2px;
  margin-bottom: 60px;
}

// 特性列表
.features {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-bottom: 48px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 18px 24px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  transition: all 0.3s ease;
  cursor: default;

  &:hover {
    background: rgba(255, 255, 255, 0.06);
    border-color: rgba($gold-primary, 0.3);
    transform: translateX(8px);
  }
}

.feature-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: rgba($gold-primary, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  color: $gold-primary;
  flex-shrink: 0;

  svg {
    width: 24px;
    height: 24px;
  }
}

.feature-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.feature-title {
  font-size: 16px;
  font-weight: 600;
  color: $text-primary;
}

.feature-desc {
  font-size: 13px;
  color: $text-secondary;
}

// 合作伙伴
.partners {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.partners-label {
  font-size: 13px;
  color: $text-muted;
}

.partner-list {
  display: flex;
  gap: 10px;
}

.partner {
  font-size: 12px;
  color: $text-secondary;
  padding: 6px 14px;
  background: rgba(255, 255, 255, 0.04);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.06);
}

// ==========================================
// 右侧登录面板
// ==========================================
.login-panel {
  flex: 0.8;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 60px;
  background: rgba($bg-card, 0.4);
  backdrop-filter: blur(40px);
  border-left: 1px solid rgba(255, 255, 255, 0.05);
  min-width: 480px;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: rgba($bg-card, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 24px;
  padding: 48px 40px;
  box-shadow:
    0 25px 50px rgba(0, 0, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.05);
}

// 欢迎区
.welcome-section {
  text-align: center;
  margin-bottom: 40px;
}

.welcome-title {
  font-size: 28px;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: 8px;
}

.welcome-desc {
  font-size: 15px;
  color: $text-secondary;
}

// 表单
.login-form {
  :deep(.el-form-item) {
    margin-bottom: 0;
  }

  :deep(.el-form-item__error) {
    padding-top: 6px;
  }
}

.form-group {
  margin-bottom: 24px;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
  color: $text-primary;
  margin-bottom: 10px;

  svg {
    width: 18px;
    height: 18px;
    color: $gold-primary;
  }
}

.gold-input {
  :deep(.el-input__wrapper) {
    background: $bg-surface;
    border: 2px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    padding: 4px 16px;
    box-shadow: none;
    transition: all 0.3s ease;

    &:hover {
      border-color: rgba($gold-primary, 0.3);
    }

    &.is-focus {
      border-color: $gold-primary;
      background: rgba($gold-primary, 0.05);
      box-shadow: 0 0 0 4px rgba($gold-primary, 0.1);
    }
  }

  :deep(.el-input__inner) {
    font-size: 15px;
    color: $text-primary;
    height: 48px;

    &::placeholder {
      color: $text-muted;
    }
  }

  :deep(.el-input__prefix) {
    .input-prefix-line {
      width: 3px;
      height: 20px;
      background: linear-gradient(180deg, $gold-primary, $gold-dark);
      border-radius: 2px;
      margin-right: 8px;
    }
  }
}

.form-options {
  margin-bottom: 28px;
}

.gold-checkbox {
  :deep(.el-checkbox__label) {
    color: $text-secondary;
    font-size: 14px;
  }

  :deep(.el-checkbox__inner) {
    background: transparent;
    border-color: $text-muted;
    border-radius: 4px;
    width: 18px;
    height: 18px;
  }

  :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
    background: $gold-primary;
    border-color: $gold-primary;
  }
}

.login-button {
  width: 100%;
  height: 52px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 2px;
  border-radius: 12px;
  background: linear-gradient(135deg, $gold-primary 0%, $gold-dark 100%);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;

  .btn-icon {
    width: 20px;
    height: 20px;
  }

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 12px 28px rgba($gold-primary, 0.4);
  }

  &:active {
    transform: translateY(0);
  }

  &.is-loading {
    background: $gold-dark;
    cursor: wait;
  }
}

// 安全徽章
.security-badges {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.badge {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: $text-muted;

  svg {
    width: 16px;
    height: 16px;
    color: $green-success;
  }
}

// 底部
.footer {
  margin-top: 32px;

  p {
    font-size: 12px;
    color: $text-muted;
    letter-spacing: 0.5px;
  }
}

// ==========================================
// 响应式
// ==========================================
@media (max-width: 1200px) {
  .brand-panel {
    padding: 40px;
  }

  .login-panel {
    min-width: 420px;
    padding: 40px;
  }

  .brand-name {
    font-size: 48px;
  }
}

@media (max-width: 1024px) {
  .login-content {
    flex-direction: column;
  }

  .brand-panel {
    padding: 48px 24px;

    .features {
      display: none;
    }

    .partners {
      justify-content: center;
    }
  }

  .brand-wrapper {
    text-align: center;

    .logo-container {
      display: flex;
      justify-content: center;
    }
  }

  .login-panel {
    min-width: auto;
    border-left: none;
    border-top: 1px solid rgba(255, 255, 255, 0.05);
  }
}

@media (max-width: 640px) {
  .brand-panel {
    padding: 32px 20px;
  }

  .brand-name {
    font-size: 36px;
    letter-spacing: 4px;
  }

  .brand-tagline {
    font-size: 16px;
    letter-spacing: 4px;
  }

  .brand-slogan {
    margin-bottom: 32px;
  }

  .login-panel {
    padding: 32px 20px;
  }

  .login-card {
    padding: 32px 24px;
    border-radius: 20px;
  }

  .welcome-title {
    font-size: 24px;
  }

  .security-badges {
    flex-wrap: wrap;
    gap: 12px;
  }
}

// 减少动画
@media (prefers-reduced-motion: reduce) {
  .glow-sphere {
    animation: none;
  }

  .feature-item,
  .login-button {
    transition: none;
  }
}
</style>
