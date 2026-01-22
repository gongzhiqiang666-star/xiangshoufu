<template>
  <div class="sim-cashback-editor">
    <el-alert type="warning" :closable="false" class="mb-4">
      <template #title>
        <el-icon><InfoFilled /></el-icon>
        配置流量卡返现：商户缴纳流量费（99元/年）后，按次数返现给代理商
      </template>
    </el-alert>

    <el-form label-position="top">
      <el-row :gutter="16">
        <el-col :span="8">
          <el-form-item label="首次返现">
            <el-input-number
              v-model="formData.firstCashback"
              :min="0"
              :max="limits?.firstCashback ?? 99"
              :precision="2"
              :step="1"
              controls-position="right"
              @change="emitChange"
            />
            <div class="limit-hint" v-if="limits?.firstCashback">
              最高 ¥{{ limits.firstCashback }}
            </div>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="二次返现">
            <el-input-number
              v-model="formData.secondCashback"
              :min="0"
              :max="limits?.secondCashback ?? 99"
              :precision="2"
              :step="1"
              controls-position="right"
              @change="emitChange"
            />
            <div class="limit-hint" v-if="limits?.secondCashback">
              最高 ¥{{ limits.secondCashback }}
            </div>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="后续返现（第三次及以后）">
            <el-input-number
              v-model="formData.thirdPlusCashback"
              :min="0"
              :max="limits?.thirdPlusCashback ?? 99"
              :precision="2"
              :step="1"
              controls-position="right"
              @change="emitChange"
            />
            <div class="limit-hint" v-if="limits?.thirdPlusCashback">
              最高 ¥{{ limits.thirdPlusCashback }}
            </div>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { reactive, watch, onMounted } from 'vue'
import { InfoFilled } from '@element-plus/icons-vue'

interface SimCashbackConfig {
  first_time_cashback: number
  second_time_cashback: number
  third_plus_cashback: number
}

interface Limits {
  firstCashback?: number
  secondCashback?: number
  thirdPlusCashback?: number
}

const props = defineProps<{
  modelValue?: SimCashbackConfig
  limits?: Limits
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: SimCashbackConfig): void
}>()

const formData = reactive({
  firstCashback: 0,
  secondCashback: 0,
  thirdPlusCashback: 0,
})

onMounted(() => {
  if (props.modelValue) {
    formData.firstCashback = props.modelValue.first_time_cashback / 100
    formData.secondCashback = props.modelValue.second_time_cashback / 100
    formData.thirdPlusCashback = props.modelValue.third_plus_cashback / 100
  }
})

watch(() => props.modelValue, (newVal) => {
  if (newVal) {
    formData.firstCashback = newVal.first_time_cashback / 100
    formData.secondCashback = newVal.second_time_cashback / 100
    formData.thirdPlusCashback = newVal.third_plus_cashback / 100
  }
}, { deep: true })

function emitChange() {
  emit('update:modelValue', {
    first_time_cashback: formData.firstCashback * 100,
    second_time_cashback: formData.secondCashback * 100,
    third_plus_cashback: formData.thirdPlusCashback * 100,
  })
}
</script>

<style scoped>
.sim-cashback-editor {
  padding: 16px;
  background: #fffbf0;
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
