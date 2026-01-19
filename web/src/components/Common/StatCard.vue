<template>
  <div class="stat-card" :style="{ borderLeftColor: color }">
    <div class="card-content">
      <div class="card-header">
        <span class="card-title">{{ title }}</span>
        <el-icon class="card-icon" :style="{ color }">
          <component :is="icon" />
        </el-icon>
      </div>

      <div class="card-value">
        <span v-if="prefix" class="prefix">{{ prefix }}</span>
        <span class="value">{{ value }}</span>
        <span v-if="suffix" class="suffix">{{ suffix }}</span>
      </div>

      <div v-if="trend !== undefined" class="card-trend">
        <span :class="['trend-value', trendClass]">
          <el-icon>
            <CaretTop v-if="trend >= 0" />
            <CaretBottom v-else />
          </el-icon>
          {{ Math.abs(trend).toFixed(2) }}%
        </span>
        <span class="trend-label">{{ trendLabel }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { CaretTop, CaretBottom } from '@element-plus/icons-vue'

interface Props {
  title: string
  value: string | number
  prefix?: string
  suffix?: string
  trend?: number
  trendLabel?: string
  icon?: string
  color?: string
}

const props = withDefaults(defineProps<Props>(), {
  color: '#409eff',
  trendLabel: '较昨日',
})

// 趋势样式类
const trendClass = computed(() => {
  if (props.trend === undefined) return ''
  return props.trend >= 0 ? 'trend-up' : 'trend-down'
})
</script>

<style lang="scss" scoped>
.stat-card {
  background: $bg-white;
  border-radius: $border-radius-md;
  padding: $spacing-lg;
  box-shadow: $shadow-sm;
  border-left: 4px solid;
  transition: all $transition-normal;

  &:hover {
    box-shadow: $shadow-md;
    transform: translateY(-2px);
  }
}

.card-content {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 14px;
  color: $text-secondary;
}

.card-icon {
  font-size: 24px;
  opacity: 0.8;
}

.card-value {
  display: flex;
  align-items: baseline;
  gap: 4px;

  .prefix {
    font-size: 16px;
    color: $text-primary;
  }

  .value {
    font-size: 28px;
    font-weight: 600;
    color: $text-primary;
    line-height: 1.2;
  }

  .suffix {
    font-size: 14px;
    color: $text-secondary;
  }
}

.card-trend {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  font-size: 12px;
}

.trend-value {
  display: flex;
  align-items: center;

  &.trend-up {
    color: $success-color;
  }

  &.trend-down {
    color: $danger-color;
  }

  .el-icon {
    font-size: 12px;
  }
}

.trend-label {
  color: $text-secondary;
}
</style>
