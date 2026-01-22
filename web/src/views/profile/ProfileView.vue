<template>
  <div class="profile-view" v-loading="loading">
    <!-- 基本信息卡片 -->
    <el-row :gutter="20">
      <el-col :xs="24" :lg="8">
        <el-card class="profile-card">
          <template #header>
            <div class="card-header">
              <span class="title">基本信息</span>
            </div>
          </template>
          <div class="user-info-section">
            <div class="avatar-wrapper">
              <el-avatar :size="80" class="user-avatar">
                {{ userInitial }}
              </el-avatar>
            </div>
            <div class="info-list">
              <div class="info-item">
                <span class="label">姓名</span>
                <span class="value">{{ profile?.contact_name || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">手机号</span>
                <span class="value">{{ maskPhone(profile?.contact_phone) }}</span>
              </div>
              <div class="info-item">
                <span class="label">服务商编号</span>
                <span class="value">{{ profile?.agent_no || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">服务商名称</span>
                <span class="value">{{ profile?.agent_name || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="label">代理层级</span>
                <el-tag size="small" :type="profile?.level === 1 ? 'warning' : 'info'">
                  {{ profile?.level }}级代理
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">入网时间</span>
                <span class="value">{{ formatDate(profile?.register_time) }}</span>
              </div>
              <div class="info-item">
                <span class="label">状态</span>
                <el-tag size="small" :type="profile?.status === 1 ? 'success' : 'danger'">
                  {{ profile?.status_name || '未知' }}
                </el-tag>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 结算卡信息 -->
      <el-col :xs="24" :lg="8">
        <el-card class="profile-card">
          <template #header>
            <div class="card-header">
              <span class="title">结算卡信息</span>
              <el-button type="primary" link @click="showBankCardDialog">
                <el-icon><Edit /></el-icon>
                修改
              </el-button>
            </div>
          </template>
          <div class="info-list">
            <div class="info-item">
              <span class="label">开户银行</span>
              <span class="value">{{ profile?.bank_name || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">开户名</span>
              <span class="value">{{ profile?.bank_account || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">银行卡号</span>
              <span class="value">{{ maskBankCard(profile?.bank_card_no) }}</span>
            </div>
            <div class="info-item">
              <span class="label">身份证号</span>
              <span class="value">{{ maskIdCard(profile?.id_card_no) }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 邀请码信息 -->
      <el-col :xs="24" :lg="8">
        <el-card class="profile-card">
          <template #header>
            <div class="card-header">
              <span class="title">邀请码</span>
            </div>
          </template>
          <div class="invite-section">
            <div class="invite-code-display">
              <span class="code">{{ inviteInfo?.invite_code || '-' }}</span>
              <el-button type="primary" link @click="copyInviteCode">
                <el-icon><CopyDocument /></el-icon>
                复制
              </el-button>
            </div>
            <div class="qr-code-section" v-if="inviteInfo?.qr_code_url">
              <img :src="inviteInfo.qr_code_url" alt="邀请二维码" class="qr-code-image" />
              <p class="qr-tip">扫码邀请下级代理</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 团队统计 -->
    <el-row :gutter="20" class="stats-section">
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-value">{{ profile?.direct_agent_count || 0 }}</div>
          <div class="stat-label">直属代理</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-value">{{ profile?.team_agent_count || 0 }}</div>
          <div class="stat-label">团队代理</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-value">{{ profile?.direct_merchant_count || 0 }}</div>
          <div class="stat-label">直属商户</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-value">{{ profile?.team_merchant_count || 0 }}</div>
          <div class="stat-label">团队商户</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷操作 -->
    <el-card class="quick-actions-card">
      <template #header>
        <div class="card-header">
          <span class="title">快捷操作</span>
        </div>
      </template>
      <div class="action-buttons">
        <el-button @click="$router.push('/change-password')">
          <el-icon><Lock /></el-icon>
          修改密码
        </el-button>
        <el-button @click="$router.push('/policies/list')">
          <el-icon><Document /></el-icon>
          查看政策
        </el-button>
        <el-button @click="$router.push('/wallets/list')">
          <el-icon><Wallet /></el-icon>
          我的钱包
        </el-button>
        <el-button @click="$router.push('/agents/list')">
          <el-icon><User /></el-icon>
          我的团队
        </el-button>
      </div>
    </el-card>

    <!-- 修改结算卡弹窗 -->
    <el-dialog v-model="bankCardDialogVisible" title="修改结算卡" width="500px">
      <el-form :model="bankCardForm" :rules="bankCardRules" ref="bankCardFormRef" label-width="100px">
        <el-form-item label="开户银行" prop="bank_name">
          <el-input v-model="bankCardForm.bank_name" placeholder="请输入开户银行" />
        </el-form-item>
        <el-form-item label="开户名" prop="bank_account">
          <el-input v-model="bankCardForm.bank_account" placeholder="请输入开户名" />
        </el-form-item>
        <el-form-item label="银行卡号" prop="bank_card_no">
          <el-input v-model="bankCardForm.bank_card_no" placeholder="请输入银行卡号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="bankCardDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleUpdateBankCard" :loading="updating">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Edit, CopyDocument, Lock, Document, Wallet, User } from '@element-plus/icons-vue'
import { getAgentDetail, getInviteCode, updateAgentProfile } from '@/api/agent'
import type { AgentDetail } from '@/types'

// 加载状态
const loading = ref(false)
const updating = ref(false)

// 个人资料
const profile = ref<AgentDetail | null>(null)

// 邀请码信息
const inviteInfo = ref<{ invite_code: string; qr_code_url: string } | null>(null)

// 用户名首字母
const userInitial = computed(() => {
  const name = profile.value?.contact_name || profile.value?.agent_name || 'U'
  return name.charAt(0).toUpperCase()
})

// 结算卡弹窗
const bankCardDialogVisible = ref(false)
const bankCardFormRef = ref<FormInstance>()
const bankCardForm = reactive({
  bank_name: '',
  bank_account: '',
  bank_card_no: '',
})

const bankCardRules: FormRules = {
  bank_name: [{ required: true, message: '请输入开户银行', trigger: 'blur' }],
  bank_account: [{ required: true, message: '请输入开户名', trigger: 'blur' }],
  bank_card_no: [
    { required: true, message: '请输入银行卡号', trigger: 'blur' },
    { pattern: /^\d{16,19}$/, message: '请输入正确的银行卡号', trigger: 'blur' },
  ],
}

// 脱敏函数
function maskPhone(phone?: string): string {
  if (!phone || phone.length !== 11) return phone || '-'
  return `${phone.substring(0, 3)}****${phone.substring(7)}`
}

function maskIdCard(idCard?: string): string {
  if (!idCard || idCard.length !== 18) return idCard || '-'
  return `${idCard.substring(0, 3)}***********${idCard.substring(14)}`
}

function maskBankCard(cardNo?: string): string {
  if (!cardNo || cardNo.length < 4) return cardNo || '-'
  return `**** **** **** ${cardNo.substring(cardNo.length - 4)}`
}

// 格式化日期
function formatDate(dateStr?: string): string {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

// 加载个人资料
async function loadProfile() {
  loading.value = true
  try {
    const data = await getAgentDetail(0) // 0 表示获取当前用户
    profile.value = data
  } catch (error) {
    console.error('Failed to load profile:', error)
    ElMessage.error('加载个人资料失败')
  } finally {
    loading.value = false
  }
}

// 加载邀请码
async function loadInviteCode() {
  try {
    const data = await getInviteCode()
    inviteInfo.value = data
  } catch (error) {
    console.error('Failed to load invite code:', error)
  }
}

// 显示修改结算卡弹窗
function showBankCardDialog() {
  if (profile.value) {
    bankCardForm.bank_name = profile.value.bank_name || ''
    bankCardForm.bank_account = profile.value.bank_account || ''
    bankCardForm.bank_card_no = profile.value.bank_card_no || ''
  }
  bankCardDialogVisible.value = true
}

// 更新结算卡
async function handleUpdateBankCard() {
  const valid = await bankCardFormRef.value?.validate()
  if (!valid) return

  updating.value = true
  try {
    await updateAgentProfile({
      bank_name: bankCardForm.bank_name,
      bank_account: bankCardForm.bank_account,
      bank_card_no: bankCardForm.bank_card_no,
    })
    ElMessage.success('结算卡信息更新成功')
    bankCardDialogVisible.value = false
    loadProfile()
  } catch (error) {
    console.error('Failed to update bank card:', error)
    ElMessage.error('更新失败，请重试')
  } finally {
    updating.value = false
  }
}

// 复制邀请码
function copyInviteCode() {
  if (!inviteInfo.value?.invite_code) return
  navigator.clipboard.writeText(inviteInfo.value.invite_code)
  ElMessage.success('邀请码已复制')
}

// 页面加载时获取数据
onMounted(() => {
  loadProfile()
  loadInviteCode()
})
</script>

<style lang="scss" scoped>
.profile-view {
  min-height: 100%;
}

.profile-card {
  margin-bottom: $spacing-lg;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .title {
      font-size: 16px;
      font-weight: 500;
      color: $text-primary;
    }
  }
}

.user-info-section {
  .avatar-wrapper {
    text-align: center;
    margin-bottom: $spacing-lg;
  }

  .user-avatar {
    background-color: $primary-color;
    color: #ffffff;
    font-size: 32px;
    font-weight: 600;
  }
}

.info-list {
  .info-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: $spacing-sm 0;
    border-bottom: 1px solid $border-color-lighter;

    &:last-child {
      border-bottom: none;
    }

    .label {
      color: $text-secondary;
      font-size: 14px;
    }

    .value {
      color: $text-primary;
      font-size: 14px;
      font-weight: 500;
    }
  }
}

.invite-section {
  text-align: center;

  .invite-code-display {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: $spacing-md;
    margin-bottom: $spacing-lg;

    .code {
      font-size: 24px;
      font-weight: 600;
      color: $primary-color;
      letter-spacing: 2px;
    }
  }

  .qr-code-section {
    .qr-code-image {
      width: 150px;
      height: 150px;
      border: 1px solid $border-color-lighter;
      border-radius: $border-radius-sm;
    }

    .qr-tip {
      margin-top: $spacing-sm;
      font-size: 12px;
      color: $text-secondary;
    }
  }
}

.stats-section {
  margin-bottom: $spacing-lg;

  .el-col {
    margin-bottom: $spacing-md;
  }
}

.stat-card {
  text-align: center;
  padding: $spacing-lg;

  .stat-value {
    font-size: 32px;
    font-weight: 600;
    color: $primary-color;
  }

  .stat-label {
    font-size: 14px;
    color: $text-secondary;
    margin-top: $spacing-xs;
  }
}

.quick-actions-card {
  .card-header {
    .title {
      font-size: 16px;
      font-weight: 500;
      color: $text-primary;
    }
  }

  .action-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: $spacing-md;

    .el-button {
      display: flex;
      align-items: center;
      gap: $spacing-xs;
    }
  }
}
</style>
