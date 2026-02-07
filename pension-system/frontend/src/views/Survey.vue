<template>
  <div class="survey">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>民意调查管理</span>
          <el-button type="primary" @click="showCreateDialog = true">新建调查</el-button>
        </div>
      </template>

      <el-table :data="surveys" border style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="title" label="调查标题" min-width="200" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column label="时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.start_date) }} - {{ formatDate(row.end_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_votes" label="投票数" width="80" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewResults(row)">查看结果</el-button>
            <el-button
              v-if="row.status === 'draft'"
              type="success"
              size="small"
              @click="updateStatus(row, 'active')"
            >
              发布
            </el-button>
            <el-button
              v-if="row.status === 'active'"
              type="warning"
              size="small"
              @click="updateStatus(row, 'closed')"
            >
              结束
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        style="margin-top: 20px"
        @current-change="loadSurveys"
      />
    </el-card>

    <!-- Create Survey Dialog -->
    <el-dialog v-model="showCreateDialog" title="新建调查" width="600px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="调查标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入调查标题" />
        </el-form-item>
        <el-form-item label="调查描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入调查描述（可选）" />
        </el-form-item>
        <el-form-item label="调查选项" prop="options" required>
          <div style="margin-bottom: 8px; color: #909399; font-size: 12px;">
            请填写调查的可选答案，至少需要2个选项
          </div>
          <div v-for="(_option, index) in options" :key="index" style="margin-bottom: 10px">
            <el-input v-model="options[index]" :placeholder="`选项 ${index + 1}`">
              <template #prepend>
                <span style="width: 60px; text-align: center">选项{{ index + 1 }}</span>
              </template>
              <template #append>
                <el-button :icon="Delete" @click="removeOption(index)" />
              </template>
            </el-input>
          </div>
          <el-button type="dashed" @click="addOption" style="width: 100%">
            <el-icon><Plus /></el-icon> 添加选项
          </el-button>
        </el-form-item>
        <el-form-item label="开始日期" prop="startDate">
          <el-date-picker v-model="form.startDate" type="date" placeholder="选择日期" style="width: 100%" />
        </el-form-item>
        <el-form-item label="结束日期" prop="endDate">
          <el-date-picker v-model="form.endDate" type="date" placeholder="选择日期" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleCreate">创建</el-button>
      </template>
    </el-dialog>

    <!-- Results Dialog -->
    <el-dialog v-model="showResultsDialog" title="调查结果" width="600px">
      <div v-if="currentResults">
        <h3 style="margin-bottom: 8px">{{ currentSurvey?.title }}</h3>
        <p style="color: #606060; margin-bottom: 20px">
          总投票数: <strong>{{ currentResults.total_votes || 0 }}</strong>
        </p>
        <div v-if="currentResults.total_votes > 0" ref="chartRef" style="width: 100%; height: 400px; margin-top: 20px"></div>
        <el-empty v-else description="暂无投票数据" :image-size="200" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import api from '@/api'

const loading = ref(false)
const submitLoading = ref(false)
const surveys = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const showCreateDialog = ref(false)
const showResultsDialog = ref(false)
const currentSurvey = ref<any>(null)
const currentResults = ref<any>(null)
const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

const formRef = ref<FormInstance>()
const form = reactive({
  title: '',
  description: '',
  startDate: '',
  endDate: ''
})
const options = ref<string[]>(['', ''])

const rules: FormRules = {
  title: [{ required: true, message: '请输入调查标题', trigger: 'blur' }],
  startDate: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
  endDate: [{ required: true, message: '请选择结束日期', trigger: 'change' }]
}

const loadSurveys = async () => {
  loading.value = true
  try {
    const response = await api.getSurveyList(currentPage.value, pageSize.value)
    if (response.success) {
      surveys.value = response.data || []
      total.value = response.total || 0
    }
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('zh-CN')
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    draft: 'info',
    active: 'success',
    closed: 'warning'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    draft: '草稿',
    active: '进行中',
    closed: '已结束'
  }
  return map[status] || status
}

const addOption = () => {
  options.value.push('')
}

const removeOption = (index: number) => {
  if (options.value.length > 2) {
    options.value.splice(index, 1)
  } else {
    ElMessage.warning('至少需要两个选项')
  }
}

const handleCreate = async () => {
  if (!formRef.value) return

  // First validate form fields
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    // Then validate options separately
    const validOptions = options.value.filter(o => o.trim() !== '')
    if (validOptions.length < 2) {
      ElMessage.warning('请至少填写两个有效选项')
      return
    }

    submitLoading.value = true
    try {
      const startDate = typeof form.startDate === 'string'
        ? form.startDate
        : (form.startDate as Date).toISOString().split('T')[0]
      const endDate = typeof form.endDate === 'string'
        ? form.endDate
        : (form.endDate as Date).toISOString().split('T')[0]

      const response = await api.createSurvey(
        form.title,
        form.description,
        JSON.stringify(validOptions),
        startDate,
        endDate,
        1
      )

      if (response.success) {
        ElMessage.success('创建成功')
        showCreateDialog.value = false
        formRef.value?.resetFields()
        options.value = ['', '']
        form.title = ''
        form.description = ''
        form.startDate = ''
        form.endDate = ''
        loadSurveys()
      } else {
        ElMessage.error(response.message)
      }
    } catch (error) {
      ElMessage.error('创建失败')
    } finally {
      submitLoading.value = false
    }
  })
}

const updateStatus = async (survey: any, status: string) => {
  try {
    const response = await api.updateSurveyStatus(survey.id, status)
    if (response.success) {
      ElMessage.success('操作成功')
      loadSurveys()
    }
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleDelete = (survey: any) => {
  ElMessageBox.confirm('确定要删除这个调查吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const response = await api.deleteSurvey(survey.id)
      if (response.success) {
        ElMessage.success('删除成功')
        loadSurveys()
      }
    } catch (error) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

const viewResults = async (survey: any) => {
  currentSurvey.value = survey
  try {
    const response = await api.getSurveyResults(survey.id)
    if (response.success) {
      currentResults.value = response.data
      showResultsDialog.value = true
      await nextTick()
      initChart()
    } else {
      ElMessage.error(response.message || '加载结果失败')
    }
  } catch (error) {
    ElMessage.error('加载结果失败')
  }
}

const initChart = () => {
  if (!chartRef.value || !currentResults.value || !currentSurvey.value) return
  if (!currentResults.value.total_votes || currentResults.value.total_votes === 0) return

  if (!chart) {
    chart = echarts.init(chartRef.value)
  }

  // Parse options from survey - handle both string and already parsed data
  let optionList: string[] = []
  try {
    if (typeof currentSurvey.value.options === 'string') {
      optionList = JSON.parse(currentSurvey.value.options || '[]')
    } else if (Array.isArray(currentSurvey.value.options)) {
      optionList = currentSurvey.value.options
    }
  } catch (e) {
    console.error('Failed to parse options:', e)
    optionList = []
  }

  // Build data for chart
  const data = optionList.map((option: string) => ({
    name: option,
    value: currentResults.value.option_counts?.[option] || 0
  }))

  // Dark cyan-blue color palette
  const colors = ['#0369a1', '#0891b2', '#0d9488', '#06b6d4', '#1d4ed8', '#0284c7', '#14b8a6', '#4f46e5']

  chart.setOption({
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center',
      textStyle: { color: '#606060' }
    },
    color: colors,
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      data,
      itemStyle: {
        borderRadius: 8,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: true,
        formatter: '{b}\n{c}票 ({d}%)',
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
      }
    }]
  })
}

onMounted(() => {
  loadSurveys()
})

onUnmounted(() => {
  chart?.dispose()
})
</script>

<style scoped>
.survey {
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
