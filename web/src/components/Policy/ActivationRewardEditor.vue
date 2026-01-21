<template>
  <div class="activation-reward-editor">
    <el-alert type="success" :closable="false" class="mb-4">
      <template #title>
        <el-icon><InfoFilled /></el-icon>
        配置激活奖励：商户入网后达成交易量目标，给代理商发放奖励
      </template>
    </el-alert>

    <div v-if="!rewardItems.length" class="empty-state">
      <el-empty description="暂无激活奖励配置">
        <el-button type="primary" @click="addReward">添加奖励</el-button>
      </el-empty>
    </div>

    <div v-else>
      <el-card v-for="(item, index) in rewardItems" :key="index" class="reward-card">
        <template #header>
          <div class="card-header">
            <span>奖励 {{ index + 1 }}</span>
            <el-button type="danger" text @click="removeReward(index)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </template>
        
        <el-form label-position="top">
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="奖励名称">
                <el-input v-model="item.rewardName" placeholder="如：首刷奖励" @change="emitChange" />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="最少入网天数">
                <el-input-number v-model="item.minDays" :min="0" :max="365" @change="emitChange" />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="最多入网天数">
                <el-input-number v-model="item.maxDays" :min="item.minDays" :max="365" @change="emitChange" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="目标交易量（万元）">
                <el-input-number 
                  v-model="item.targetAmountWan" 
                  :min="0" 
                  :precision="2"
                  :step="0.1"
                  @change="emitChange" 
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="奖励金额（元）">
                <el-input-number 
                  v-model="item.rewardAmountYuan" 
                  :min="0" 
                  :max="getMaxReward(index)"
                  :precision="2"
                  :step="1"
                  @change="emitChange" 
                />
                <div class="limit-hint" v-if="limits?.[index]">
                  最高 ¥{{ limits[index] }}
                </div>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-card>

      <el-button type="primary" plain @click="addReward" class="add-btn">
        <el-icon><Plus /></el-icon>
        添加奖励
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { InfoFilled, Delete, Plus } from '@element-plus/icons-vue'

interface ActivationRewardItem {
  reward_name: string
  min_register_days: number
  max_register_days: number
  target_amount: number
  reward_amount: number
}

interface RewardEditItem {
  rewardName: string
  minDays: number
  maxDays: number
  targetAmountWan: number
  rewardAmountYuan: number
}

const props = defineProps<{
  modelValue: ActivationRewardItem[]
  limits?: number[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: ActivationRewardItem[]): void
}>()

const rewardItems = ref<RewardEditItem[]>([])

onMounted(() => {
  if (props.modelValue?.length) {
    rewardItems.value = props.modelValue.map(toEditItem)
  }
})

watch(() => props.modelValue, (newVal) => {
  if (newVal?.length) {
    rewardItems.value = newVal.map(toEditItem)
  }
}, { deep: true })

function toEditItem(item: ActivationRewardItem): RewardEditItem {
  return {
    rewardName: item.reward_name,
    minDays: item.min_register_days,
    maxDays: item.max_register_days,
    targetAmountWan: item.target_amount / 1000000,
    rewardAmountYuan: item.reward_amount / 100,
  }
}

function toApiItem(item: RewardEditItem): ActivationRewardItem {
  return {
    reward_name: item.rewardName,
    min_register_days: item.minDays,
    max_register_days: item.maxDays,
    target_amount: item.targetAmountWan * 1000000,
    reward_amount: item.rewardAmountYuan * 100,
  }
}

function getMaxReward(index: number): number {
  return props.limits?.[index] ?? 9999
}

function addReward() {
  rewardItems.value.push({
    rewardName: '',
    minDays: 0,
    maxDays: 30,
    targetAmountWan: 1,
    rewardAmountYuan: 0,
  })
  emitChange()
}

function removeReward(index: number) {
  rewardItems.value.splice(index, 1)
  emitChange()
}

function emitChange() {
  emit('update:modelValue', rewardItems.value.map(toApiItem))
}
</script>

<style scoped>
.activation-reward-editor {
  padding: 16px;
  background: #f0f9eb;
  border-radius: 8px;
}
.reward-card {
  margin-bottom: 16px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.limit-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
.empty-state {
  padding: 32px;
  text-align: center;
}
.add-btn {
  width: 100%;
}
.mb-4 {
  margin-bottom: 16px;
}
</style>
