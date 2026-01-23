<template>
  <el-dialog
    v-model="visible"
    title="修改费率"
    width="500px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item label="商户编号">
        <span class="merchant-info">{{ merchant?.merchant_no }}</span>
      </el-form-item>
      <el-form-item label="商户名称">
        <span class="merchant-info">{{ merchant?.merchant_name }}</span>
      </el-form-item>
      <el-form-item label="当前贷记卡费率">
        <span class="current-rate">{{ formatRate(merchant?.credit_rate) }}</span>
      </el-form-item>
      <el-form-item label="当前借记卡费率">
        <span class="current-rate">{{ formatRate(merchant?.debit_rate) }}</span>
      </el-form-item>
      <el-divider />
      <el-form-item label="新贷记卡费率" prop="credit_rate">
        <el-input-number
          v-model="form.credit_rate"
          :min="0"
          :max="1"
          :step="0.0001"
          :precision="4"
          :controls="false"
          style="width: 200px"
        />
        <span class="rate-unit">（如：0.0060 表示 0.60%）</span>
      </el-form-item>
      <el-form-item label="新借记卡费率" prop="debit_rate">
        <el-input-number
          v-model="form.debit_rate"
          :min="0"
          :max="1"
          :step="0.0001"
          :precision="4"
          :controls="false"
          style="width: 200px"
        />
        <span class="rate-unit">（如：0.0055 表示 0.55%）</span>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">
        确认修改
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { updateMerchantRate } from '@/api/merchant'

interface MerchantInfo {
  id: number
  merchant_no: string
  merchant_name: string
  credit_rate: string
  debit_rate: string
}

const props = defineProps<{
  modelValue: boolean
  merchant: MerchantInfo | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'success': []
}>()

const visible = ref(false)
const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = reactive({
  credit_rate: 0,
  debit_rate: 0,
})

const rules: FormRules = {
  credit_rate: [
    { required: true, message: '请输入贷记卡费率', trigger: 'blur' },
    { type: 'number', min: 0, max: 0.1, message: '费率范围: 0 ~ 0.1', trigger: 'blur' },
  ],
  debit_rate: [
    { required: true, message: '请输入借记卡费率', trigger: 'blur' },
    { type: 'number', min: 0, max: 0.1, message: '费率范围: 0 ~ 0.1', trigger: 'blur' },
  ],
}

// 格式化费率显示
function formatRate(rate: string | undefined): string {
  if (!rate) return '-'
  const num = parseFloat(rate)
  if (isNaN(num)) return rate
  return `${(num * 100).toFixed(2)}%`
}

// 同步 visible
watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.merchant) {
    // 初始化表单
    form.credit_rate = parseFloat(props.merchant.credit_rate) || 0
    form.debit_rate = parseFloat(props.merchant.debit_rate) || 0
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

// 关闭弹窗
function handleClose() {
  visible.value = false
  formRef.value?.resetFields()
}

// 提交
async function handleSubmit() {
  if (!props.merchant) return

  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  submitting.value = true
  try {
    const result = await updateMerchantRate(props.merchant.id, {
      credit_rate: form.credit_rate,
      debit_rate: form.debit_rate,
    })

    // 显示同步结果
    if (result.sync_success) {
      ElMessage.success('费率修改成功，已同步到支付通道')
    } else {
      ElMessage.warning(`费率修改成功，但通道同步失败：${result.sync_message}，请人工处理`)
    }

    emit('success')
    handleClose()
  } catch (error: any) {
    console.error('Update rate error:', error)
    ElMessage.error(error.message || '费率修改失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style lang="scss" scoped>
.merchant-info {
  color: #606266;
  font-weight: 500;
}

.current-rate {
  color: #409eff;
  font-weight: 600;
}

.rate-unit {
  margin-left: 8px;
  color: #909399;
  font-size: 12px;
}

.el-divider {
  margin: 16px 0;
}
</style>
