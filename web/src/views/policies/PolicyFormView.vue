<template>
  <div class="policy-form-view">
    <PageHeader
      :title="isEdit ? '编辑政策模板' : '新建政策模板'"
      :sub-title="isEdit ? `模板ID: ${route.params.id}` : '创建新的政策模板（包含4块政策配置）'"
    >
      <template #extra>
        <el-button @click="handleBack">返回</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ isEdit ? '保存修改' : '创建模板' }}
        </el-button>
      </template>
    </PageHeader>

    <el-card class="form-card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        :disabled="loading"
      >
        <!-- 基本信息 -->
        <div class="form-section">
          <div class="section-title">基本信息</div>
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="模板名称" prop="template_name">
                <el-input v-model="form.template_name" placeholder="请输入模板名称" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="所属通道" prop="channel_id">
                <ChannelSelect v-model="form.channel_id" style="width: 100%" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="状态" prop="status">
                <el-radio-group v-model="form.status">
                  <el-radio :value="1">启用</el-radio>
                  <el-radio :value="0">禁用</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="设为默认">
                <el-switch v-model="form.is_default" />
                <span class="form-tip">设为该通道的默认政策模板</span>
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <!-- 1. 成本（费率配置） -->
        <div class="form-section">
          <div class="section-title">
            <span>1. 成本配置（费率）</span>
            <el-tag type="success" size="small">分润钱包</el-tag>
          </div>
          <el-row :gutter="24" v-loading="rateTypesLoading">
            <template v-if="rateTypes && rateTypes.length > 0 && Object.keys(form.rate_configs).length > 0">
              <el-col :span="8" v-for="rateType in rateTypes" :key="rateType.code">
                <el-form-item :label="rateType.name" v-if="form.rate_configs[rateType.code]">
                  <el-input-number
                    v-model="form.rate_configs[rateType.code].rate"
                    :min="parseFloat(rateType.min_rate)"
                    :max="parseFloat(rateType.max_rate)"
                    :precision="4"
                    :step="0.01"
                    style="width: 140px"
                  />
                  <span class="form-unit">%</span>
                  <span class="rate-range">({{ rateType.min_rate }}~{{ rateType.max_rate }})</span>
                </el-form-item>
              </el-col>
            </template>
            <el-col :span="24" v-else-if="!rateTypesLoading">
              <el-empty description="请先选择通道" :image-size="60" />
            </el-col>
          </el-row>
        </div>

        <!-- 2. 押金返现配置 -->
        <div class="form-section">
          <div class="section-title">
            <span>2. 押金返现配置</span>
            <el-tag type="info" size="small">服务费钱包</el-tag>
          </div>
          <DepositCashbackEditor v-model="form.deposit_cashbacks" />
        </div>

        <!-- 3. 流量卡返现配置 -->
        <div class="form-section">
          <div class="section-title">
            <span>3. 流量卡返现配置</span>
            <el-tag type="info" size="small">服务费钱包</el-tag>
          </div>
          <SimCashbackEditor v-model="form.sim_cashback" />
        </div>

        <!-- 4. 激活奖励配置 -->
        <div class="form-section">
          <div class="section-title">
            <span>4. 激活奖励配置</span>
            <el-tag type="warning" size="small">奖励钱包</el-tag>
          </div>
          <ActivationRewardEditor v-model="form.activation_rewards" />
        </div>

        <!-- 5. 费率阶梯配置（可选） -->
        <div class="form-section">
          <div class="section-title">
            <span>5. 费率阶梯配置（代理商调价）</span>
            <el-tag size="small">可选</el-tag>
          </div>
          <RateStageEditor v-model="form.rate_stages" />
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import ChannelSelect from '@/components/Common/ChannelSelect.vue'
import {
  DepositCashbackEditor,
  SimCashbackEditor,
  ActivationRewardEditor,
  RateStageEditor,
} from '@/components/Policy'
import { getPolicyTemplateDetail, createPolicyTemplate, updatePolicyTemplate } from '@/api/policy'
import { getChannelRateTypes } from '@/api/channel'
import type { RateTypeDefinition, RateConfigs } from '@/types/policy'

const route = useRoute()
const router = useRouter()

// 是否编辑模式
const isEdit = computed(() => route.name === 'PolicyEdit')

// 加载状态
const loading = ref(false)
const submitting = ref(false)

// 动态费率类型
const rateTypes = ref<RateTypeDefinition[]>([])
const rateTypesLoading = ref(false)

// 表单引用
const formRef = ref<FormInstance>()

// 押金返现类型
interface DepositCashbackItem {
  deposit_amount: number
  cashback_amount: number
}

// 流量卡返现类型
interface SimCashbackConfig {
  first_time_cashback: number
  second_time_cashback: number
  third_plus_cashback: number
  sim_fee_amount?: number
}

// 激活奖励类型
interface ActivationRewardItem {
  reward_name: string
  min_register_days: number
  max_register_days: number
  target_amount: number
  reward_amount: number
  priority: number
}

// 费率阶梯类型
interface RateStageItem {
  stage_name: string
  apply_to: number
  min_days: number
  max_days: number
  credit_rate_delta: number
  debit_rate_delta: number
  unionpay_rate_delta: number
  wechat_rate_delta: number
  alipay_rate_delta: number
  priority: number
}

// 表单数据
const form = reactive({
  template_name: '',
  channel_id: undefined as number | undefined,
  status: 1,
  is_default: false,
  // 1. 成本（动态费率配置）
  rate_configs: {} as RateConfigs,
  // 2. 押金返现
  deposit_cashbacks: [] as DepositCashbackItem[],
  // 3. 流量卡返现
  sim_cashback: null as SimCashbackConfig | null,
  // 4. 激活奖励
  activation_rewards: [] as ActivationRewardItem[],
  // 5. 费率阶梯
  rate_stages: [] as RateStageItem[],
})

// 表单验证规则
const rules: FormRules = {
  template_name: [
    { required: true, message: '请输入模板名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度为2-50个字符', trigger: 'blur' },
  ],
  channel_id: [
    { required: true, message: '请选择所属通道', trigger: 'change' },
  ],
}

// 监听通道变化，动态加载费率类型
watch(() => form.channel_id, async (channelId) => {
  if (channelId) {
    rateTypesLoading.value = true
    try {
      const result = await getChannelRateTypes(channelId)
      rateTypes.value = result || []
      // 每次加载完费率类型后都初始化配置
      initRateConfigs()
    } catch (error) {
      console.error('Failed to load rate types:', error)
      rateTypes.value = []
      form.rate_configs = {}
    } finally {
      rateTypesLoading.value = false
    }
  } else {
    rateTypes.value = []
    form.rate_configs = {}
  }
})

// 初始化费率配置
function initRateConfigs() {
  if (!rateTypes.value || rateTypes.value.length === 0) {
    form.rate_configs = {}
    return
  }
  const configs: RateConfigs = {}
  for (const rt of rateTypes.value) {
    if (form.rate_configs[rt.code]) {
      // 保留已有配置
      configs[rt.code] = form.rate_configs[rt.code]
    } else {
      // 使用最小费率作为默认值
      configs[rt.code] = { rate: rt.min_rate }
    }
  }
  form.rate_configs = configs
}

// 获取模板详情
async function fetchDetail() {
  if (!isEdit.value) return

  loading.value = true
  try {
    const data = await getPolicyTemplateDetail(Number(route.params.id))
    Object.assign(form, {
      template_name: data.template_name,
      channel_id: data.channel_id,
      status: data.status,
      is_default: data.is_default,
      // 动态费率配置
      rate_configs: data.rate_configs || {},
      deposit_cashbacks: data.deposit_cashbacks || [],
      sim_cashback: data.sim_cashback || null,
      activation_rewards: data.activation_rewards || [],
      rate_stages: data.rate_stages || [],
    })
    // 加载费率类型后初始化配置
    if (data.channel_id) {
      rateTypesLoading.value = true
      try {
        rateTypes.value = await getChannelRateTypes(data.channel_id)
        initRateConfigs()
      } finally {
        rateTypesLoading.value = false
      }
    }
  } catch (error) {
    console.error('Fetch policy template error:', error)
    ElMessage.error('获取模板详情失败')
  } finally {
    loading.value = false
  }
}

// 返回列表
function handleBack() {
  router.push('/policies/list')
}

// 提交表单
async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  // 验证费率配置
  if (rateTypes.value.length > 0 && Object.keys(form.rate_configs).length === 0) {
    ElMessage.warning('请配置费率')
    return
  }

  submitting.value = true
  try {
    const submitData = {
      template_name: form.template_name,
      channel_id: form.channel_id,
      is_default: form.is_default,
      // 动态费率配置
      rate_configs: form.rate_configs,
      deposit_cashbacks: form.deposit_cashbacks,
      sim_cashback: form.sim_cashback,
      activation_rewards: form.activation_rewards,
      rate_stages: form.rate_stages,
    }

    if (isEdit.value) {
      await updatePolicyTemplate(Number(route.params.id), submitData)
      ElMessage.success('更新成功')
    } else {
      await createPolicyTemplate(submitData)
      ElMessage.success('创建成功')
    }
    router.push('/policies/list')
  } catch (error) {
    console.error('Submit policy template error:', error)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchDetail()
})
</script>

<style lang="scss" scoped>
.policy-form-view {
  padding: 0;
}

.form-card {
  margin-top: 16px;
}

.form-section {
  margin-bottom: 32px;

  &:last-child {
    margin-bottom: 0;
  }

  .section-title {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid #ebeef5;
    display: flex;
    align-items: center;
    gap: 12px;
  }
}

.form-unit {
  margin-left: 8px;
  color: #909399;
}

.form-tip {
  margin-left: 16px;
  font-size: 12px;
  color: #c0c4cc;
}

.rate-range {
  margin-left: 8px;
  font-size: 12px;
  color: #909399;
}
</style>
