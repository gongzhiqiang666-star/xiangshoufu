<template>
  <div class="send-message-view">
    <PageHeader title="消息管理" sub-title="发送消息">
      <template #extra>
        <el-button @click="handleBack">返回列表</el-button>
      </template>
    </PageHeader>

    <el-card>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        style="max-width: 800px"
      >
        <el-form-item label="消息标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入消息标题" maxlength="64" show-word-limit />
        </el-form-item>

        <el-form-item label="消息内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            placeholder="请输入消息内容"
            :rows="6"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="消息类型" prop="message_type">
          <el-select v-model="form.message_type" placeholder="请选择消息类型" style="width: 200px">
            <el-option label="系统公告" :value="6" />
          </el-select>
          <span class="form-tip">管理员只能发送系统公告类型的消息</span>
        </el-form-item>

        <el-form-item label="有效期" prop="expire_days">
          <el-slider
            v-model="form.expire_days"
            :min="1"
            :max="30"
            :marks="expireMarks"
            show-input
            style="width: 400px"
          />
          <span class="form-tip">消息有效期（天），过期后自动清理</span>
        </el-form-item>

        <el-form-item label="发送范围" prop="send_scope">
          <el-radio-group v-model="form.send_scope">
            <el-radio value="all">全部代理商</el-radio>
            <el-radio value="agents">指定代理商</el-radio>
            <el-radio value="level">指定层级</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="form.send_scope === 'agents'" label="选择代理商" prop="agent_ids">
          <el-select
            v-model="form.agent_ids"
            multiple
            filterable
            remote
            reserve-keyword
            placeholder="请搜索并选择代理商"
            :remote-method="searchAgents"
            :loading="agentLoading"
            style="width: 100%"
          >
            <el-option
              v-for="agent in agentOptions"
              :key="agent.id"
              :label="`${agent.agent_name} (${agent.agent_no})`"
              :value="agent.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item v-if="form.send_scope === 'level'" label="选择层级" prop="level">
          <el-select v-model="form.level" placeholder="请选择层级" style="width: 200px">
            <el-option label="一级代理" :value="1" />
            <el-option label="二级代理" :value="2" />
            <el-option label="三级代理" :value="3" />
            <el-option label="四级代理" :value="4" />
            <el-option label="五级代理" :value="5" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">
            发送消息
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import { sendMessage } from '@/api/message'
import { searchAgentList } from '@/api/agent'
import type { SendMessageRequest, MessageTypeValue } from '@/types'

const router = useRouter()

// 表单引用
const formRef = ref<FormInstance>()

// 表单数据
const form = reactive<SendMessageRequest & { agent_ids: number[] }>({
  title: '',
  content: '',
  message_type: 6 as MessageTypeValue, // 默认系统公告
  expire_days: 3,
  send_scope: 'all',
  agent_ids: [],
  level: undefined,
})

// 有效期标记
const expireMarks = {
  1: '1天',
  3: '3天',
  7: '7天',
  14: '14天',
  30: '30天',
}

// 验证规则
const rules = reactive<FormRules>({
  title: [
    { required: true, message: '请输入消息标题', trigger: 'blur' },
    { min: 2, max: 64, message: '标题长度为2-64个字符', trigger: 'blur' },
  ],
  content: [
    { required: true, message: '请输入消息内容', trigger: 'blur' },
    { min: 1, max: 500, message: '内容长度不能超过500个字符', trigger: 'blur' },
  ],
  message_type: [{ required: true, message: '请选择消息类型', trigger: 'change' }],
  send_scope: [{ required: true, message: '请选择发送范围', trigger: 'change' }],
  agent_ids: [
    {
      validator: (_: any, value: number[], callback: any) => {
        if (form.send_scope === 'agents' && (!value || value.length === 0)) {
          callback(new Error('请选择代理商'))
        } else {
          callback()
        }
      },
      trigger: 'change',
    },
  ],
  level: [
    {
      validator: (_: any, value: number, callback: any) => {
        if (form.send_scope === 'level' && !value) {
          callback(new Error('请选择层级'))
        } else {
          callback()
        }
      },
      trigger: 'change',
    },
  ],
})

// 代理商选项
const agentOptions = ref<Array<{ id: number; agent_name: string; agent_no: string }>>([])
const agentLoading = ref(false)

// 提交状态
const submitting = ref(false)

// 搜索代理商
async function searchAgents(query: string) {
  if (query.length < 1) {
    agentOptions.value = []
    return
  }

  agentLoading.value = true
  try {
    const res = await searchAgentList({ keyword: query, page: 1, page_size: 20 })
    agentOptions.value = res.list || []
  } catch (error) {
    console.error('Search agents error:', error)
    agentOptions.value = []
  } finally {
    agentLoading.value = false
  }
}

// 提交
async function handleSubmit() {
  if (!formRef.value) return

  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const data: SendMessageRequest = {
      title: form.title,
      content: form.content,
      message_type: form.message_type,
      expire_days: form.expire_days,
      send_scope: form.send_scope,
    }

    if (form.send_scope === 'agents') {
      data.agent_ids = form.agent_ids
    } else if (form.send_scope === 'level') {
      data.level = form.level
    }

    const res = await sendMessage(data)
    ElMessage.success(`消息发送成功，已发送给 ${res.sent_count} 个代理商`)
    router.push('/system/messages')
  } catch (error: any) {
    ElMessage.error(error.message || '发送失败')
  } finally {
    submitting.value = false
  }
}

// 重置
function handleReset() {
  formRef.value?.resetFields()
  form.send_scope = 'all'
  form.agent_ids = []
  form.level = undefined
}

// 返回
function handleBack() {
  router.push('/system/messages')
}
</script>

<style lang="scss" scoped>
.send-message-view {
  padding: 0;
}

.form-tip {
  margin-left: 12px;
  color: #909399;
  font-size: 12px;
}
</style>
