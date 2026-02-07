<template>
  <div class="data-summary">
    <!-- Stats Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <div class="stat-card modern-card primary">
          <div class="stat-icon">
            <el-icon :size="32"><User /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-label">老年人总数</div>
            <div class="stat-number">{{ summaryData.total_elderly }}</div>
            <div class="stat-trend">
              <el-icon><TrendCharts /></el-icon>
              <span>实时统计</span>
            </div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card modern-card success">
          <div class="stat-icon">
            <el-icon :size="32"><Male /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-label">男性人数</div>
            <div class="stat-number">{{ summaryData.gender_distribution?.['男'] || 0 }}</div>
            <div class="stat-trend">
              <el-icon><PieChart /></el-icon>
              <span>{{ getPercentage(summaryData.gender_distribution?.['男'] || 0) }}%</span>
            </div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card modern-card warning">
          <div class="stat-icon">
            <el-icon :size="32"><Female /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-label">女性人数</div>
            <div class="stat-number">{{ summaryData.gender_distribution?.['女'] || 0 }}</div>
            <div class="stat-trend">
              <el-icon><PieChart /></el-icon>
              <span>{{ getPercentage(summaryData.gender_distribution?.['女'] || 0) }}%</span>
            </div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card modern-card info">
          <div class="stat-icon">
            <el-icon :size="32"><DataAnalysis /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-label">数据更新</div>
            <div class="stat-number">今天</div>
            <div class="stat-trend">
              <el-icon><Clock /></el-icon>
              <span>实时数据</span>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Charts Row -->
    <el-row :gutter="16" class="charts-row">
      <el-col :span="12">
        <el-card class="chart-card modern-card" shadow="never">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <el-icon><PieChart /></el-icon>
                <span>性别分布</span>
              </div>
              <el-button text @click="loadSummary" class="refresh-btn" :loading="loading">
                <el-icon><Refresh /></el-icon>
              </el-button>
            </div>
          </template>
          <div ref="genderChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card modern-card" shadow="never">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <el-icon><Histogram /></el-icon>
                <span>年龄分布</span>
              </div>
              <el-button text @click="loadSummary" class="refresh-btn" :loading="loading">
                <el-icon><Refresh /></el-icon>
              </el-button>
            </div>
          </template>
          <div ref="ageChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="charts-row">
      <el-col :span="12">
        <el-card class="chart-card modern-card" shadow="never">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <el-icon><Suitcase /></el-icon>
                <span>护理等级分布</span>
              </div>
              <el-button text @click="loadSummary" class="refresh-btn" :loading="loading">
                <el-icon><Refresh /></el-icon>
              </el-button>
            </div>
          </template>
          <div ref="careChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card modern-card" shadow="never">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <el-icon><FirstAidKit /></el-icon>
                <span>健康状况分布</span>
              </div>
              <el-button text @click="loadSummary" class="refresh-btn" :loading="loading">
                <el-icon><Refresh /></el-icon>
              </el-button>
            </div>
          </template>
          <div ref="healthChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, onActivated } from 'vue'
import * as echarts from 'echarts'
import {
  User, Male, Female, DataAnalysis, TrendCharts,
  PieChart, Clock, Histogram, Suitcase, FirstAidKit, Refresh
} from '@element-plus/icons-vue'
import api from '@/api'

const genderChartRef = ref<HTMLElement>()
const ageChartRef = ref<HTMLElement>()
const careChartRef = ref<HTMLElement>()
const healthChartRef = ref<HTMLElement>()

let genderChart: echarts.ECharts | null = null
let ageChart: echarts.ECharts | null = null
let careChart: echarts.ECharts | null = null
let healthChart: echarts.ECharts | null = null
let refreshTimer: number | null = null

const loading = ref(false)

const summaryData = ref<any>({
  total_elderly: 0,
  gender_distribution: {},
  age_distribution: {},
  health_conditions: {},
  care_levels: {}
})

const getPercentage = (value: number) => {
  if (summaryData.value.total_elderly === 0) return 0
  return ((value / summaryData.value.total_elderly) * 100).toFixed(1)
}

const loadSummary = async () => {
  loading.value = true
  try {
    const response = await api.getSummary()
    if (response.success) {
      summaryData.value = response.data
      initCharts()
    }
  } catch (error) {
    console.error('Failed to load summary:', error)
  } finally {
    loading.value = false
  }
}

const initCharts = () => {
  // Darker cyan-blue cool tone color palette
  const colors = [
    '#0369a1',   // Dark Sky Blue
    '#0891b2',   // Dark Cyan
    '#0d9488',   // Dark Teal
    '#06b6d4',   // Cyan
    '#1d4ed8',   // Deep Blue
    '#0284c7'    // Sky
  ]

  // Gender Chart
  if (genderChartRef.value) {
    if (!genderChart) {
      genderChart = echarts.init(genderChartRef.value)
    }
    const genderData = Object.entries(summaryData.value.gender_distribution || {})
      .map(([key, value]) => ({ name: key, value: value as number }))
      .filter(item => item.value > 0)
    genderChart.setOption({
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      legend: {
        bottom: 0,
        left: 'center',
        textStyle: { color: '#606060' }
      },
      color: colors,
      series: [{
        type: 'pie',
        radius: ['35%', '65%'],
        center: ['50%', '45%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: true,
          formatter: '{b}\n{c}人',
          fontSize: 13,
          color: '#606060'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold'
          },
          itemStyle: {
            shadowBlur: 15,
            shadowOffsetX: 0,
            shadowColor: 'rgba(3, 105, 161, 0.4)'
          }
        },
        data: genderData
      }]
    }, { notMerge: true })
  }

  // Age Chart
  if (ageChartRef.value) {
    if (!ageChart) {
      ageChart = echarts.init(ageChartRef.value)
    }
    const ageData = Object.entries(summaryData.value.age_distribution || {})
      .filter(([, value]) => (value as number) > 0)
    ageChart.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: { type: 'shadow' }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '15%',
        top: '10%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: ageData.map(([key]) => key),
        axisLabel: { interval: 0, rotate: 0, fontSize: 12, color: '#606060' },
        axisLine: { lineStyle: { color: '#e0e0e0' } },
        axisTick: { show: false }
      },
      yAxis: {
        type: 'value',
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { lineStyle: { color: '#f0f0f0' } },
        axisLabel: { color: '#606060' }
      },
      series: [{
        data: ageData.map(([, value]) => value),
        type: 'bar',
        barWidth: '50%',
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#06b6d4' },
            { offset: 0.5, color: '#0d9488' },
            { offset: 1, color: '#0369a1' }
          ]),
          borderRadius: [8, 8, 0, 0]
        },
        emphasis: {
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#1d4ed8' },
              { offset: 0.5, color: '#0891b2' },
              { offset: 1, color: '#4f46e5' }
            ])
          }
        }
      }]
    }, { notMerge: true })
  }

  // Care Level Chart
  if (careChartRef.value) {
    if (!careChart) {
      careChart = echarts.init(careChartRef.value)
    }
    const careData = Object.entries(summaryData.value.care_levels || {})
      .filter(([, value]) => (value as number) > 0)
    careChart.setOption({
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      legend: {
        bottom: 0,
        left: 'center',
        textStyle: { color: '#606060' }
      },
      color: colors.slice(2, 5),
      series: [{
        type: 'pie',
        radius: ['35%', '65%'],
        center: ['50%', '45%'],
        data: careData.map(([key, value]) => ({ name: key || '未分级', value })),
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          formatter: '{b}\n{c}人',
          fontSize: 12,
          color: '#606060'
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 15,
            shadowOffsetX: 0,
            shadowColor: 'rgba(13, 148, 136, 0.4)'
          }
        }
      }]
    }, { notMerge: true })
  }

  // Health Condition Chart
  if (healthChartRef.value) {
    if (!healthChart) {
      healthChart = echarts.init(healthChartRef.value)
    }
    const healthData = Object.entries(summaryData.value.health_conditions || {})
      .filter(([, value]) => (value as number) > 0)
    healthChart.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: { type: 'shadow' }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '15%',
        top: '10%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: healthData.map(([key]) => key || '未填写'),
        axisLabel: { interval: 0, rotate: 20, fontSize: 11, color: '#606060' },
        axisLine: { lineStyle: { color: '#e0e0e0' } },
        axisTick: { show: false }
      },
      yAxis: {
        type: 'value',
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { lineStyle: { color: '#f0f0f0' } },
        axisLabel: { color: '#606060' }
      },
      series: [{
        data: healthData.map(([, value]) => value),
        type: 'bar',
        barWidth: '45%',
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#0891b2' },
            { offset: 0.5, color: '#0d9488' },
            { offset: 1, color: '#1d4ed8' }
          ]),
          borderRadius: [8, 8, 0, 0]
        },
        emphasis: {
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#06b6d4' },
              { offset: 0.5, color: '#0284c7' },
              { offset: 1, color: '#4f46e5' }
            ])
          }
        }
      }]
    }, { notMerge: true })
  }
}

const handleResize = () => {
  genderChart?.resize()
  ageChart?.resize()
  careChart?.resize()
  healthChart?.resize()
}

// Auto refresh every 10 seconds
const startAutoRefresh = () => {
  refreshTimer = window.setInterval(() => {
    loadSummary()
  }, 10000)
}

const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

onMounted(() => {
  loadSummary()
  startAutoRefresh()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  stopAutoRefresh()
  genderChart?.dispose()
  ageChart?.dispose()
  careChart?.dispose()
  healthChart?.dispose()
  window.removeEventListener('resize', handleResize)
})

// Refresh when page is activated
onActivated(() => {
  loadSummary()
})
</script>

<style scoped>
.data-summary {
  width: 100%;
}

.stats-row {
  margin-bottom: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  border-radius: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 28px rgba(3, 105, 161, 0.2);
}

.stat-card.modern-card {
  background: white;
  border: 1px solid rgba(0, 0, 0, 0.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.stat-card.primary .stat-icon {
  background: linear-gradient(135deg, #0369a1 0%, #0891b2 100%);
  color: white;
}

.stat-card.success .stat-icon {
  background: linear-gradient(135deg, #06b6d4 0%, #0d9488 100%);
  color: white;
}

.stat-card.warning .stat-icon {
  background: linear-gradient(135deg, #0284c7 0%, #1d4ed8 100%);
  color: white;
}

.stat-card.info .stat-icon {
  background: linear-gradient(135deg, #14b8a6 0%, #0d9488 100%);
  color: white;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
}

.stat-content {
  flex: 1;
}

.stat-label {
  font-size: 13px;
  color: #888;
  margin-bottom: 4px;
  font-weight: 500;
}

.stat-number {
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, #0369a1 0%, #0d9488 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1.2;
  margin-bottom: 6px;
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #0369a1;
  font-weight: 500;
}

.charts-row {
  margin-bottom: 16px;
}

.chart-card {
  border-radius: 16px;
  border: 1px solid rgba(0, 0, 0, 0.05);
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

:deep(.chart-card .el-card__header) {
  padding: 18px 22px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  background: linear-gradient(to bottom, #fafafa, #ffffff);
}

:deep(.chart-card .el-card__body) {
  padding: 22px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  background: linear-gradient(135deg, #0369a1 0%, #0d9488 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.refresh-btn {
  color: #0369a1;
  transition: all 0.3s ease;
}

.refresh-btn:hover {
  color: #0d9488;
  transform: rotate(180deg);
}

.chart-container {
  width: 100%;
  height: 280px;
}
</style>
