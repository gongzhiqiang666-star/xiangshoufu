<template>
  <el-dialog
    v-model="visible"
    title="分配政策模板"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
      <el-form-item label="代理商" prop="agent_name">
        <el-input v-model="formData.agent_name" disabled />
      </el-form-item>

      <el-form-item label="通道" prop="channel_id">
        <el-select
          v-model="formData.channel_id"
          placeholder="请选择通道"
          style="width: 100%"
          @change="handleChannelChange"
        >
          <el-option
            v-for="channel in channels"
            :key="channel.id"
            :label="channel.channel_name"
            :value="channel.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="政策模板" prop="template_id">
        <el-select
          v-model="formData.template_id"
          placeholder="请选择政策模板"
          style="width: 100%"
          :disabled="!formData.channel_id"
          @change="handleTemplateChange"
        >
          <el-option
            v-for="template in templates"
            :key="template.id"
            :label="template.name"
            :value="template.id"
          >
            <span>{{ template.name }}</span>
            <el-tag v-if="template.is_default" size="small" type="success" style="margin-left: 8px">
              默认
            </el-tag>
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="贷记卡费率" prop="credit_rate">
        <el-input-number
          v-model="formData.credit_rate"
          :min="0"
          :max="10"
          :precision="4"
          :step="0.001"
          style="width: 200px"
        />
        <span class="rate-unit">% (范围: 0 - 10%)</span>
      </el-form-item>

      <el-form-item label="借记卡费率" prop="debit_rate">
        <el-input-number
          v-model="formData.debit_rate"
          :min="0"
          :max="10"
          :precision="4"
          :step="0.001"
          style="width: 200px"
        />
        <span class="rate-unit">% (范围: 0 - 10%)</span>
      </el-form-item>

      <el-alert
        v-if="selectedTemplate"
        type="info"
        :closable="false"
        style="margin-top: 10px"
      >
        <template #title>
          模板参考费率：贷记卡 {{ (selectedTemplate.credit_rate * 100).toFixed(4) }}%，
          借记卡 {{ (selectedTemplate.debit_rate * 100).toFixed(4) }}%
        </template>
      </el-alert>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="saving" @click="handleSave">
        确定
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getChannels } from '@/api/agent-channel'
import { getTemplatesByChannel, assignAgentPolicy } from '@/api/policy'

interface Props {
  modelValue: boolean
  agentId: number
  agentName: string
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const visible = ref(false)
const saving = ref(false)
const formRef = ref<FormInstance>()

interface Channel {
  id: number
  channel_name: string
  channel_code: string
}

interface PolicyTemplate {
  id: number
  name: string
  is_default: boolean
  credit_rate: number
  debit_rate: number
}

const channels = ref<Channel[]>([])
const templates = ref<PolicyTemplate[]>([])
const selectedTemplate = ref<PolicyTemplate | null>(null)

const formData = reactive({
  agent_name: '',
  channel_id: null as number | null,
  template_id: null as number | null,
  credit_rate: 0.6,
  debit_rate: 0.6,
})

const rules: FormRules = {
  channel_id: [{ required: true, message: '请选择通道', trigger: 'change' }],
  template_id: [{ required: true, message: '请选择政策模板', trigger: 'change' }],
  credit_rate: [{ required: true, message: '请输入贷记卡费率', trigger: 'blur' }],
  debit_rate: [{ required: true, message: '请输入借记卡费率', trigger: 'blur' }],
}

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) {
    formData.agent_name = props.agentName
    fetchChannels()
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

async function fetchChannels() {
  try {
    const res = await getChannels()
    channels.value = res as Channel[]
  } catch (error) {
    console.error('Fetch channels error:', error)
  }
}

async function handleChannelChange(channelId: number) {
  formData.template_id = null
  selectedTemplate.value = null
  templates.value = []

  if (channelId) {
    try {
      const res = await getTemplatesByChannel(channelId)
      templates.value = res as PolicyTemplate[]
    } catch (error) {
      console.error('Fetch templates error:', error)
    }
  }
}

function handleTemplateChange(templateId: number) {
  const template = templates.value.find(t => t.id === templateId)
  if (template) {
    selectedTemplate.value = template
    // 使用模板的费率作为默认值（转换为百分比显示）
    formData.credit_rate = template.credit_rate * 100
    formData.debit_rate = template.debit_rate * 100
  }
}

async function handleSave() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    saving.value = true
    try {
      await assignAgentPolicy(props.agentId, {
        channel_id: formData.channel_id!,
        template_id: formData.template_id!,
        credit_rate: formData.credit_rate / 100, // 转换回小数
        debit_rate: formData.debit_rate / 100,
      })
      ElMessage.success('政策分配成功')
      emit('success')
      handleClose()
    } catch (error) {
      ElMessage.error('政策分配失败')
    } finally {
      saving.value = false
    }
  })
}

function handleClose() {
  visible.value = false
  formRef.value?.resetFields()
  formData.channel_id = null
  formData.template_id = null
  formData.credit_rate = 0.6
  formData.debit_rate = 0.6
  selectedTemplate.value = null
  templates.value = []
}
</script>

<style lang="scss" scoped>
.rate-unit {
  margin-left: 10px;
  color: #909399;
  font-size: 12px;
}
</style>
