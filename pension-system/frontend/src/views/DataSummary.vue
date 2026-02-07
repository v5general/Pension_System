<template>
  <div class="data-summary">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon total">
              <el-icon size="40"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ summaryData.total_elderly }}</div>
              <div class="stat-label">老年人总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="18">
        <el-card class="chart-card">
          <template #header>
            <span>性别分布</span>
          </template>
          <div ref="genderChartRef" style="width: 100%; height: 300px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>年龄分布</span>
          </template>
          <div ref="ageChartRef" style="width: 100%; height: 350px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>护理等级分布</span>
          </template>
          <div ref="careChartRef" style="width: 100%; height: 350px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row style="margin-top: 20px">
      <el-col :span="24">
        <el-card class="chart-card">
          <template #header>
            <span>健康状况分布</span>
          </template>
          <div ref="healthChartRef" style="width: 100%; height: 350px"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import api from '@/api'

const genderChartRef = ref<HTMLElement>()
const ageChartRef = ref<HTMLElement>()
const careChartRef = ref<HTMLElement>()
const healthChartRef = ref<HTMLElement>()

let genderChart: echarts.ECharts | null = null
let ageChart: echarts.ECharts | null = null
let careChart: echarts.ECharts | null = null
let healthChart: echarts.ECharts | null = null

const summaryData = ref<any>({
  total_elderly: 0,
  gender_distribution: {},
  age_distribution: {},
  health_conditions: {},
  care_levels: {}
})

const loadSummary = async () => {
  try {
    const response = await api.getSummary()
    if (response.success) {
      summaryData.value = response.data
      initCharts()
    }
  } catch (error) {
    console.error('Failed to load summary:', error)
  }
}

const initCharts = () => {
  // Gender Chart
  if (genderChartRef.value) {
    if (!genderChart) {
      genderChart = echarts.init(genderChartRef.value)
    }
    const genderData = Object.entries(summaryData.value.gender_distribution).map(([key, value]) => ({
      name: key,
      value: value
    }))
    genderChart.setOption({
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        right: 10,
        top: 'center'
      },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: true,
          formatter: '{b}: {c}'
        },
        data: genderData
      }]
    })
  }

  // Age Chart
  if (ageChartRef.value) {
    if (!ageChart) {
      ageChart = echarts.init(ageChartRef.value)
    }
    const ageData = Object.entries(summaryData.value.age_distribution)
    ageChart.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: { type: 'shadow' }
      },
      xAxis: {
        type: 'category',
        data: ageData.map(([key]) => key),
        axisLabel: { interval: 0, rotate: 0 }
      },
      yAxis: {
        type: 'value'
      },
      series: [{
        data: ageData.map(([, value]) => value),
        type: 'bar',
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#83bff6' },
            { offset: 1, color: '#188df0' }
          ])
        }
      }]
    })
  }

  // Care Level Chart
  if (careChartRef.value) {
    if (!careChart) {
      careChart = echarts.init(careChartRef.value)
    }
    const careData = Object.entries(summaryData.value.care_levels)
    careChart.setOption({
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      series: [{
        type: 'pie',
        radius: '70%',
        data: careData.map(([key, value]) => ({ name: key, value })),
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }]
    })
  }

  // Health Condition Chart
  if (healthChartRef.value) {
    if (!healthChart) {
      healthChart = echarts.init(healthChartRef.value)
    }
    const healthData = Object.entries(summaryData.value.health_conditions)
    healthChart.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: { type: 'shadow' }
      },
      xAxis: {
        type: 'category',
        data: healthData.map(([key]) => key || '未填写'),
        axisLabel: { interval: 0, rotate: 30 }
      },
      yAxis: {
        type: 'value'
      },
      series: [{
        data: healthData.map(([, value]) => value),
        type: 'bar',
        itemStyle: {
          color: '#67c23a'
        }
      }]
    })
  }
}

const handleResize = () => {
  genderChart?.resize()
  ageChart?.resize()
  careChart?.resize()
  healthChart?.resize()
}

onMounted(() => {
  loadSummary()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  genderChart?.dispose()
  ageChart?.dispose()
  careChart?.dispose()
  healthChart?.dispose()
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.data-summary {
  width: 100%;
}

.stat-card {
  height: 380px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 80px;
  height: 80px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 36px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #999;
  margin-top: 5px;
}

.chart-card {
  height: 400px;
}
</style>
