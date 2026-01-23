<template>
  <div ref="chartRef" class="bar-chart"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'

interface Props {
  data: {
    categories: string[]
    values: number[]
  }
  title?: string
  height?: string
  color?: string
}

const props = withDefaults(defineProps<Props>(), {
  height: '300px',
  color: '#409eff',
})

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

const colors = [
  '#67c23a',
  '#409eff',
  '#e6a23c',
  '#f56c6c',
  '#909399',
  '#00d4ff',
]

// 初始化图表
function initChart() {
  if (!chartRef.value) return

  chart = echarts.init(chartRef.value)
  updateChart()

  window.addEventListener('resize', handleResize)
}

// 更新图表
function updateChart() {
  if (!chart) return

  const option: echarts.EChartsOption = {
    title: props.title
      ? {
          text: props.title,
          left: 'center',
          textStyle: {
            fontSize: 16,
            fontWeight: 500,
            color: '#303133',
          },
        }
      : undefined,
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
      formatter: (params: any) => {
        const data = params[0]
        return `${data.name}<br/>数量: ${data.value}`
      },
    },
    grid: {
      left: '3%',
      right: '4%',
      top: props.title ? '15%' : '10%',
      bottom: '10%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: props.data.categories,
      axisLine: {
        lineStyle: {
          color: '#dcdfe6',
        },
      },
      axisLabel: {
        color: '#606266',
        fontSize: 12,
        interval: 0,
        rotate: props.data.categories.length > 6 ? 30 : 0,
      },
    },
    yAxis: {
      type: 'value',
      axisLine: {
        show: false,
      },
      axisTick: {
        show: false,
      },
      splitLine: {
        lineStyle: {
          color: '#e4e7ed',
          type: 'dashed',
        },
      },
      axisLabel: {
        color: '#606266',
      },
    },
    series: [
      {
        type: 'bar',
        data: props.data.values.map((value, index) => ({
          value,
          itemStyle: {
            color: colors[index % colors.length],
            borderRadius: [4, 4, 0, 0],
          },
        })),
        barWidth: '40%',
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.2)',
          },
        },
      },
    ],
  }

  chart.setOption(option)
}

function handleResize() {
  chart?.resize()
}

watch(
  () => props.data,
  () => {
    updateChart()
  },
  { deep: true }
)

onMounted(() => {
  initChart()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chart?.dispose()
})
</script>

<style lang="scss" scoped>
.bar-chart {
  width: 100%;
  height: v-bind(height);
}
</style>
