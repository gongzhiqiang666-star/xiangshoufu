<template>
  <div class="alert-config-view">
    <el-card class="page-header">
      <div class="header-content">
        <div>
          <h2>告警配置</h2>
          <p class="description">配置任务告警通道，支持钉钉、企业微信、邮件通知</p>
        </div>
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          新增配置
        </el-button>
      </div>
    </el-card>

    <el-card class="main-content">
      <el-table
        v-loading="loading"
        :data="configList"
        stripe
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="配置名称" width="180" />
        <el-table-column prop="channel_type" label="通道类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.channel_type === 1" type="primary">钉钉</el-tag>
            <el-tag v-else-if="row.channel_type === 2" type="success">企业微信</el-tag>
            <el-tag v-else-if="row.channel_type === 3" type="warning">邮件</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="webhook_url" label="Webhook地址" min-width="200">
          <template #default="{ row }">
            <span v-if="row.channel_type !== 3">{{ row.webhook_url || '-' }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="email_addresses" label="邮箱地址" min-width="200">
          <template #default="{ row }">
            <span v-if="row.channel_type === 3">{{ row.email_addresses || '-' }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="is_enabled" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_enabled"
              @change="handleEnableChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleTest(row)">
              测试
            </el-button>
            <el-button type="primary" link size="small" @click="showEditDialog(row)">
              编辑
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑告警配置' : '新增告警配置'"
      width="600px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="120px"
      >
        <el-form-item label="配置名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入配置名称" />
        </el-form-item>
        <el-form-item label="通道类型" prop="channel_type">
          <el-radio-group v-model="form.channel_type" :disabled="isEdit">
            <el-radio :label="1">钉钉</el-radio>
            <el-radio :label="2">企业微信</el-radio>
            <el-radio :label="3">邮件</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 钉钉/企微配置 -->
        <template v-if="form.channel_type === 1 || form.channel_type === 2">
          <el-form-item label="Webhook地址" prop="webhook_url">
            <el-input
              v-model="form.webhook_url"
              placeholder="请输入Webhook地址"
              type="textarea"
              :rows="2"
            />
          </el-form-item>
          <el-form-item v-if="form.channel_type === 1" label="签名密钥">
            <el-input
              v-model="form.webhook_secret"
              placeholder="可选，钉钉加签密钥"
            />
            <div class="form-tip">如果机器人开启了加签验证，请填写密钥</div>
          </el-form-item>
        </template>

        <!-- 邮件配置 -->
        <template v-if="form.channel_type === 3">
          <el-form-item label="收件人邮箱" prop="email_addresses">
            <el-input
              v-model="form.email_addresses"
              placeholder="多个邮箱用逗号分隔"
              type="textarea"
              :rows="2"
            />
          </el-form-item>
          <el-form-item label="SMTP服务器" prop="email_smtp_host">
            <el-input v-model="form.email_smtp_host" placeholder="如: smtp.qq.com" />
          </el-form-item>
          <el-form-item label="SMTP端口" prop="email_smtp_port">
            <el-input-number v-model="form.email_smtp_port" :min="1" :max="65535" />
          </el-form-item>
          <el-form-item label="发件人账号" prop="email_username">
            <el-input v-model="form.email_username" placeholder="SMTP登录账号" />
          </el-form-item>
          <el-form-item label="发件人密码" prop="email_password">
            <el-input
              v-model="form.email_password"
              type="password"
              placeholder="SMTP登录密码或授权码"
              show-password
            />
            <div class="form-tip">编辑时留空表示不修改密码</div>
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSubmit">
          {{ isEdit ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getAlertConfigs,
  createAlertConfig,
  updateAlertConfig,
  deleteAlertConfig,
  enableAlertConfig,
  testAlertConfig
} from '@/api/job'
import type { AlertConfig, CreateAlertConfigRequest } from '@/types/job'

const loading = ref(false)
const saving = ref(false)
const configList = ref<AlertConfig[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentConfig = ref<AlertConfig | null>(null)
const formRef = ref<FormInstance>()

const form = reactive<CreateAlertConfigRequest>({
  name: '',
  channel_type: 1,
  webhook_url: '',
  webhook_secret: '',
  email_addresses: '',
  email_smtp_host: '',
  email_smtp_port: 465,
  email_username: '',
  email_password: ''
})

const formRules = {
  name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
  channel_type: [{ required: true, message: '请选择通道类型', trigger: 'change' }],
  webhook_url: [
    {
      validator: (_: any, value: string, callback: Function) => {
        if ((form.channel_type === 1 || form.channel_type === 2) && !value) {
          callback(new Error('请输入Webhook地址'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  email_addresses: [
    {
      validator: (_: any, value: string, callback: Function) => {
        if (form.channel_type === 3 && !value) {
          callback(new Error('请输入收件人邮箱'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  email_smtp_host: [
    {
      validator: (_: any, value: string, callback: Function) => {
        if (form.channel_type === 3 && !value) {
          callback(new Error('请输入SMTP服务器'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 加载配置列表
const loadConfigs = async () => {
  loading.value = true
  try {
    configList.value = await getAlertConfigs()
  } catch (error) {
    console.error('加载告警配置失败:', error)
  } finally {
    loading.value = false
  }
}

// 格式化日期时间
const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 重置表单
const resetForm = () => {
  form.name = ''
  form.channel_type = 1
  form.webhook_url = ''
  form.webhook_secret = ''
  form.email_addresses = ''
  form.email_smtp_host = ''
  form.email_smtp_port = 465
  form.email_username = ''
  form.email_password = ''
}

// 显示新增对话框
const showCreateDialog = () => {
  isEdit.value = false
  currentConfig.value = null
  resetForm()
  dialogVisible.value = true
}

// 显示编辑对话框
const showEditDialog = (config: AlertConfig) => {
  isEdit.value = true
  currentConfig.value = config
  form.name = config.name
  form.channel_type = config.channel_type
  form.webhook_url = config.webhook_url || ''
  form.webhook_secret = ''
  form.email_addresses = config.email_addresses || ''
  form.email_smtp_host = config.email_smtp_host || ''
  form.email_smtp_port = config.email_smtp_port || 465
  form.email_username = config.email_username || ''
  form.email_password = ''
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  saving.value = true
  try {
    if (isEdit.value && currentConfig.value) {
      await updateAlertConfig(currentConfig.value.id, form)
      ElMessage.success('更新成功')
    } else {
      await createAlertConfig(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadConfigs()
  } catch (error) {
    console.error('保存失败:', error)
  } finally {
    saving.value = false
  }
}

// 启用/禁用
const handleEnableChange = async (config: AlertConfig) => {
  try {
    await enableAlertConfig(config.id, config.is_enabled)
    ElMessage.success(config.is_enabled ? '配置已启用' : '配置已禁用')
  } catch (error) {
    config.is_enabled = !config.is_enabled
    console.error('更新状态失败:', error)
  }
}

// 测试告警
const handleTest = async (config: AlertConfig) => {
  try {
    await ElMessageBox.confirm(
      `确定要发送测试告警到 "${config.name}" 吗？`,
      '测试告警',
      { type: 'info' }
    )
    await testAlertConfig(config.id)
    ElMessage.success('测试消息已发送，请检查接收端')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('测试失败:', error)
    }
  }
}

// 删除配置
const handleDelete = async (config: AlertConfig) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除告警配置 "${config.name}" 吗？`,
      '确认删除',
      { type: 'warning' }
    )
    await deleteAlertConfig(config.id)
    ElMessage.success('删除成功')
    loadConfigs()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

onMounted(() => {
  loadConfigs()
})
</script>

<style scoped>
.alert-config-view {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content h2 {
  margin: 0 0 8px 0;
  font-size: 20px;
}

.header-content .description {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.main-content {
  min-height: 400px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>
