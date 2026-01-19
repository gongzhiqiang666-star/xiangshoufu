<template>
  <div ref="chartRef" class="line-chart"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'

interface Props {
  dates: string[]
  data: {
    name: string
    values: number[]
    color?: string
  }[]
  title?: string
  height?: string
  showLegend?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  height: '350px',
  showLegend: true,
})

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

// 初始化图表
function initChart() {
  if (!chartRef.value) return

  chart = echarts.init(chartRef.value)
  updateChart()

  // 监听窗口大小变化
  window.addEventListener('resize', handleResize)
}

// 更新图表
function updateChart() {
  if (!chart) return

  const series = props.data.map((item) => ({
    name: item.name,
    type: 'line',
    smooth: true,
    symbol: 'circle',
    symbolSize: 6,
    itemStyle: {
      color: item.color,
    },
    lineStyle: {
      width: 2,
    },
    areaStyle: {
      opacity: 0.1,
    },
    data: item.values,
  }))

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
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      borderColor: '#e4e7ed',
      borderWidth: 1,
      textStyle: {
        color: '#303133',
      },
      axisPointer: {
        type: 'cross',
        crossStyle: {
          color: '#999',
        },
      },
    },
    legend: props.showLegend
      ? {
          data: props.data.map((item) => item.name),
          bottom: 0,
          itemGap: 20,
        }
      : undefined,
    grid: {
      left: '3%',
      right: '4%',
      top: props.title ? '15%' : '10%',
      bottom: props.showLegend ? '15%' : '10%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: props.dates,
      axisLine: {
        lineStyle: {
          color: '#dcdfe6',
        },
      },
      axisLabel: {
        color: '#606266',
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
        formatter: (value: number) => {
          if (value >= 10000) {
            return (value / 10000).toFixed(1) + 'w'
          }
          return value.toString()
        },
      },
    },
    series,
  }

  chart.setOption(option)
}

// 处理窗口大小变化
function handleResize() {
  chart?.resize()
}

// 监听数据变化
watch(
  () => [props.dates, props.data],
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
.line-chart {
  width: 100%;
  height: v-bind(height);
}
</style>
