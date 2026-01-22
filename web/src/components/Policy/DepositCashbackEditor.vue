<template>
  <div class="deposit-cashback-editor">
    <el-alert type="info" :closable="false" class="mb-4">
      <template #title>
        <el-icon><InfoFilled /></el-icon>
        配置押金返现：商户缴纳押金后，按此配置返现给代理商
      </template>
    </el-alert>

    <el-form label-position="top">
      <el-row :gutter="16">
        <el-col :span="6" v-for="item in depositItems" :key="item.deposit">
          <el-form-item :label="`押金 ¥${item.deposit} 返现`">
            <el-input-number
              v-model="item.cashback"
              :min="0"
              :max="limits?.[item.deposit] ?? 999"
              :precision="2"
              :step="1"
              controls-position="right"
              @change="emitChange"
            />
            <div class="limit-hint" v-if="limits?.[item.deposit]">
              最高 ¥{{ limits[item.deposit] }}
            </div>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { InfoFilled } from '@element-plus/icons-vue'

interface DepositCashbackItem {
  deposit_amount: number
  cashback_amount: number
}

const props = defineProps<{
  modelValue: DepositCashbackItem[]
  limits?: Record<number, number>
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: DepositCashbackItem[]): void
}>()

// 默认的押金档位（分转换为元显示）
const depositItems = ref([
  { deposit: 0, cashback: 0 },
  { deposit: 99, cashback: 0 },
  { deposit: 199, cashback: 0 },
  { deposit: 299, cashback: 0 },
])

onMounted(() => {
  // 从props初始化
  if (props.modelValue?.length) {
    for (const item of props.modelValue) {
      const deposit = item.deposit_amount / 100
      const target = depositItems.value.find(d => d.deposit === deposit)
      if (target) {
        target.cashback = item.cashback_amount / 100
      }
    }
  }
})

watch(() => props.modelValue, (newVal) => {
  if (newVal?.length) {
    for (const item of newVal) {
      const deposit = item.deposit_amount / 100
      const target = depositItems.value.find(d => d.deposit === deposit)
      if (target) {
        target.cashback = item.cashback_amount / 100
      }
    }
  }
}, { deep: true })

function emitChange() {
  const result: DepositCashbackItem[] = depositItems.value.map(item => ({
    deposit_amount: item.deposit * 100,
    cashback_amount: item.cashback * 100,
  }))
  emit('update:modelValue', result)
}
</script>

<style scoped>
.deposit-cashback-editor {
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
}
.limit-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
.mb-4 {
  margin-bottom: 16px;
}
</style>
