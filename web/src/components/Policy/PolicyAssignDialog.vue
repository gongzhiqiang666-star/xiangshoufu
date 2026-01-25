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

      <!-- 动态费率配置 -->
      <div v-loading="rateTypesLoading">
        <template v-if="rateTypes.length > 0 && formData.template_id">
          <el-form-item
            v-for="rateType in rateTypes"
            :key="rateType.code"
            :label="rateType.name"
          >
            <el-input-number
              v-model="formData.rate_configs[rateType.code].rate"
              :min="parseFloat(rateType.min_rate)"
              :max="parseFloat(rateType.max_rate)"
              :precision="4"
              :step="0.01"
              style="width: 200px"
            />
            <span class="rate-unit">% ({{ rateType.min_rate }}~{{ rateType.max_rate }})</span>
          </el-form-item>
        </template>
      </div>

      <el-alert
        v-if="selectedTemplate && rateTypes.length > 0"
        type="info"
        :closable="false"
        style="margin-top: 10px"
      >
        <template #title>
          模板参考费率：
          <template v-for="(rateType, index) in rateTypes" :key="rateType.code">
            {{ rateType.name }} {{ selectedTemplate.rate_configs?.[rateType.code]?.rate || '0' }}%<template v-if="index < rateTypes.length - 1">，</template>
          </template>
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
import { getChannelRateTypes } from '@/api/channel'
import type { RateTypeDefinition, RateConfigs } from '@/types/policy'

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
  rate_configs: RateConfigs
}

const channels = ref<Channel[]>([])
const templates = ref<PolicyTemplate[]>([])
const selectedTemplate = ref<PolicyTemplate | null>(null)

// 动态费率类型
const rateTypes = ref<RateTypeDefinition[]>([])
const rateTypesLoading = ref(false)

const formData = reactive({
  agent_name: '',
  channel_id: null as number | null,
  template_id: null as number | null,
  rate_configs: {} as RateConfigs,
})

const rules: FormRules = {
  channel_id: [{ required: true, message: '请选择通道', trigger: 'change' }],
  template_id: [{ required: true, message: '请选择政策模板', trigger: 'change' }],
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
  formData.rate_configs = {}
  rateTypes.value = []

  if (channelId) {
    try {
      // 并行加载模板和费率类型
      rateTypesLoading.value = true
      const [templatesRes, rateTypesRes] = await Promise.all([
        getTemplatesByChannel(channelId),
        getChannelRateTypes(channelId)
      ])
      templates.value = templatesRes as PolicyTemplate[]
      rateTypes.value = rateTypesRes
      // 初始化费率配置
      initRateConfigs()
    } catch (error) {
      console.error('Fetch templates/rate types error:', error)
    } finally {
      rateTypesLoading.value = false
    }
  }
}

// 初始化费率配置
function initRateConfigs() {
  const configs: RateConfigs = {}
  for (const rt of rateTypes.value) {
    if (formData.rate_configs[rt.code]) {
      configs[rt.code] = formData.rate_configs[rt.code]
    } else {
      configs[rt.code] = { rate: rt.min_rate }
    }
  }
  formData.rate_configs = configs
}

function handleTemplateChange(templateId: number) {
  const template = templates.value.find(t => t.id === templateId)
  if (template) {
    selectedTemplate.value = template
    // 使用模板的费率配置作为默认值
    if (template.rate_configs) {
      formData.rate_configs = { ...template.rate_configs }
    }
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
        rate_configs: formData.rate_configs,
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
  formData.rate_configs = {}
  selectedTemplate.value = null
  templates.value = []
  rateTypes.value = []
}
</script>

<style lang="scss" scoped>
.rate-unit {
  margin-left: 10px;
  color: #909399;
  font-size: 12px;
}
</style>
