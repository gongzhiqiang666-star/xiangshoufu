<template>
  <div class="rate-stage-editor">
    <el-alert type="info" :closable="false" class="mb-4">
      <template #title>
        <el-icon><InfoFilled /></el-icon>
        配置费率阶梯：按商户/代理商入网天数自动调整费率
      </template>
    </el-alert>

    <div v-if="!stageItems.length" class="empty-state">
      <el-empty description="暂无费率阶梯配置">
        <el-button type="primary" @click="addStage">添加阶梯</el-button>
      </el-empty>
    </div>

    <div v-else>
      <el-table :data="stageItems" border>
        <el-table-column label="类型" width="150">
          <template #default="{ row }">
            <el-select v-model="row.stageType" @change="emitChange">
              <el-option :value="1" label="商户入网天数" />
              <el-option :value="2" label="代理商入网天数" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column label="开始天数" width="120">
          <template #default="{ row }">
            <el-input-number v-model="row.startDay" :min="0" :max="365" size="small" @change="emitChange" />
          </template>
        </el-table-column>
        <el-table-column label="结束天数" width="120">
          <template #default="{ row }">
            <el-input-number v-model="row.endDay" :min="row.startDay" :max="365" size="small" @change="emitChange" />
          </template>
        </el-table-column>
        <el-table-column label="贷记卡费率变化" width="140">
          <template #default="{ row }">
            <el-input-number 
              v-model="row.creditRateDelta" 
              :min="-50" 
              :max="50" 
              :step="1"
              size="small" 
              @change="emitChange" 
            />
            <span class="unit">万分点</span>
          </template>
        </el-table-column>
        <el-table-column label="借记卡费率变化" width="140">
          <template #default="{ row }">
            <el-input-number 
              v-model="row.debitRateDelta" 
              :min="-50" 
              :max="50" 
              :step="1"
              size="small" 
              @change="emitChange" 
            />
            <span class="unit">万分点</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80">
          <template #default="{ $index }">
            <el-button type="danger" text @click="removeStage($index)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-button type="primary" plain @click="addStage" class="add-btn">
        <el-icon><Plus /></el-icon>
        添加阶梯
      </el-button>
    </div>

    <div class="tips">
      <p><strong>说明：</strong></p>
      <p>• 正数表示费率上调，负数表示费率下调</p>
      <p>• 例如：贷记卡变化 +5 表示费率上调 0.05%</p>
      <p>• 阶梯按天数范围生效，不重叠</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { InfoFilled, Delete, Plus } from '@element-plus/icons-vue'

interface RateStageItem {
  stage_type: number
  start_day: number
  end_day: number
  credit_rate_delta: number
  debit_rate_delta: number
}

interface StageEditItem {
  stageType: number
  startDay: number
  endDay: number
  creditRateDelta: number
  debitRateDelta: number
}

const props = defineProps<{
  modelValue: RateStageItem[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: RateStageItem[]): void
}>()

const stageItems = ref<StageEditItem[]>([])

onMounted(() => {
  if (props.modelValue?.length) {
    stageItems.value = props.modelValue.map(toEditItem)
  }
})

watch(() => props.modelValue, (newVal) => {
  if (newVal?.length) {
    stageItems.value = newVal.map(toEditItem)
  }
}, { deep: true })

function toEditItem(item: RateStageItem): StageEditItem {
  return {
    stageType: item.stage_type,
    startDay: item.start_day,
    endDay: item.end_day,
    creditRateDelta: item.credit_rate_delta,
    debitRateDelta: item.debit_rate_delta,
  }
}

function toApiItem(item: StageEditItem): RateStageItem {
  return {
    stage_type: item.stageType,
    start_day: item.startDay,
    end_day: item.endDay,
    credit_rate_delta: item.creditRateDelta,
    debit_rate_delta: item.debitRateDelta,
  }
}

function addStage() {
  const lastEnd = stageItems.value.length > 0 
    ? stageItems.value[stageItems.value.length - 1].endDay 
    : 0
  stageItems.value.push({
    stageType: 1,
    startDay: lastEnd,
    endDay: lastEnd + 30,
    creditRateDelta: 0,
    debitRateDelta: 0,
  })
  emitChange()
}

function removeStage(index: number) {
  stageItems.value.splice(index, 1)
  emitChange()
}

function emitChange() {
  emit('update:modelValue', stageItems.value.map(toApiItem))
}
</script>

<style scoped>
.rate-stage-editor {
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
}
.unit {
  font-size: 12px;
  color: #909399;
  margin-left: 4px;
}
.empty-state {
  padding: 32px;
  text-align: center;
}
.add-btn {
  width: 100%;
  margin-top: 16px;
}
.tips {
  margin-top: 16px;
  padding: 12px;
  background: #fff;
  border-radius: 4px;
  font-size: 12px;
  color: #606266;
}
.tips p {
  margin: 4px 0;
}
.mb-4 {
  margin-bottom: 16px;
}
</style>
