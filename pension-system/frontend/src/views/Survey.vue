<template>
  <div class="survey-page">
    <!-- 管理员/操作员界面 -->
    <template v-if="userStore.isAdmin()">
      <el-card class="main-card">
        <template #header>
          <div class="card-header">
            <span class="card-title">民意调查管理</span>
            <el-button type="primary" @click="showCreateDialog = true">新建调查</el-button>
          </div>
        </template>

        <el-table :data="surveys" border stripe style="width: 100%" v-loading="loading">
          <el-table-column prop="id" label="ID" width="60" align="center" />
          <el-table-column prop="title" label="调查标题" min-width="200" />
          <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
          <el-table-column label="时间" width="180" align="center">
            <template #default="{ row }">
              {{ formatDate(row.start_date) }} - {{ formatDate(row.end_date) }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="total_votes" label="投票数" width="80" align="center" />
          <el-table-column label="操作" width="280" fixed="right" align="center">
            <template #default="{ row }">
              <el-button size="small" @click="viewResults(row)">查看结果</el-button>
              <el-button v-if="row.status === 'draft'" type="success" size="small" @click="updateStatus(row, 'active')">发布</el-button>
              <el-button v-if="row.status === 'active'" type="warning" size="small" @click="updateStatus(row, 'closed')">结束</el-button>
              <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          style="margin-top: 20px; justify-content: center"
          @current-change="loadSurveys"
        />
      </el-card>
    </template>

    <!-- 普通用户界面 -->
    <template v-else>
      <!-- 无调查时显示 -->
      <el-card v-if="activeSurveys.length === 0 && !loading" class="empty-card">
        <el-empty description="暂无进行中的调查" />
      </el-card>

      <!-- 调查列表 -->
      <div v-for="survey in activeSurveys" :key="survey.id" class="survey-item">
        <div class="survey-header">
          <h3 class="survey-title">{{ survey.title }}</h3>
          <el-tag type="success" size="small">进行中</el-tag>
        </div>

        <p v-if="survey.description" class="survey-description">{{ survey.description }}</p>

        <div class="survey-info">
          <span>开始: {{ formatDate(survey.start_date) }}</span>
          <span>结束: {{ formatDate(survey.end_date) }}</span>
          <span>参与: {{ survey.total_votes || 0 }} 人</span>
        </div>

        <!-- 已投票提示 -->
        <div v-if="hasVoted(survey.id)" class="voted-message">
          <el-alert type="success" :closable="false" show-icon>
            您已参与此调查，感谢您的参与！
          </el-alert>
        </div>

        <!-- 投票选项 -->
        <div v-else class="vote-section">
          <div class="vote-label">请选择您的答案:</div>
          <el-radio-group v-model="userVotes[survey.id]" class="vote-list">
            <div
              v-for="(option, index) in getSurveyOptions(survey)"
              :key="index"
              class="vote-item"
              :class="{ 'is-selected': userVotes[survey.id] === option }"
              @click="selectOption(survey.id, option)"
            >
              <div class="vote-item-content">
                <span class="vote-index">{{ String(index + 1).padStart(2, '0') }}</span>
                <el-radio :label="option" class="vote-radio">{{ option }}</el-radio>
              </div>
            </div>
          </el-radio-group>
          <div class="vote-actions">
            <el-button
              type="primary"
              size="large"
              :disabled="!userVotes[survey.id]"
              :loading="voteLoading[survey.id]"
              @click="submitVote(survey)"
            >
              提交投票
            </el-button>
          </div>
        </div>
      </div>

      <!-- 加载中 -->
      <div v-if="loading" v-loading="true" class="loading-container"></div>
    </template>

    <!-- 创建调查对话框 -->
    <el-dialog v-model="showCreateDialog" title="新建调查" width="600px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="调查标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入调查标题" />
        </el-form-item>
        <el-form-item label="调查描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入调查描述（可选）" />
        </el-form-item>
        <el-form-item label="调查选项" required>
          <div class="option-hint">请填写调查的可选答案，至少需要2个选项</div>
          <div v-for="(_option, index) in options" :key="index" class="option-input-row">
            <el-input v-model="options[index]" :placeholder="`选项 ${index + 1}`">
              <template #prepend>选项{{ index + 1 }}</template>
              <template #append>
                <el-button :icon="Delete" @click="removeOption(index)" :disabled="options.length <= 2" />
              </template>
            </el-input>
          </div>
          <el-button type="dashed" @click="addOption" class="add-option-btn">
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

    <!-- 结果对话框 -->
    <el-dialog v-model="showResultsDialog" title="调查结果" width="600px">
      <div v-if="currentResults">
        <h3>{{ currentSurvey?.title }}</h3>
        <p class="result-summary">总投票数: <strong>{{ currentResults.total_votes || 0 }}</strong></p>
        <div v-if="currentResults.total_votes > 0" ref="chartRef" class="chart-container"></div>
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
import { useUserStore } from '@/store/user'

const userStore = useUserStore()

const loading = ref(false)
const submitLoading = ref(false)
const surveys = ref<any[]>([])
const activeSurveys = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const showCreateDialog = ref(false)
const showResultsDialog = ref(false)
const currentSurvey = ref<any>(null)
const currentResults = ref<any>(null)
const chartRef = ref<HTMLElement>()
const userVotes = ref<Record<number, string>>({})
const voteLoading = ref<Record<number, boolean>>({})
const votedSurveyIds = ref<Set<number>>(new Set())

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

const loadActiveSurveys = async () => {
  loading.value = true
  try {
    const response = await api.getActiveSurveys()
    if (response.success) {
      activeSurveys.value = response.data || []
      for (const survey of activeSurveys.value) {
        userVotes.value[survey.id] = ''
        voteLoading.value[survey.id] = false
      }
    }
  } catch (error) {
    ElMessage.error('加载调查失败')
  } finally {
    loading.value = false
  }
}

const getSurveyOptions = (survey: any) => {
  try {
    if (typeof survey.options === 'string') {
      return JSON.parse(survey.options || '[]')
    } else if (Array.isArray(survey.options)) {
      return survey.options
    }
  } catch (e) {
    console.error('Failed to parse options:', e)
  }
  return []
}

const hasVoted = (surveyId: number) => {
  return votedSurveyIds.value.has(surveyId)
}

const selectOption = (surveyId: number, option: string) => {
  userVotes.value[surveyId] = option
}

const submitVote = async (survey: any) => {
  const selectedOption = userVotes.value[survey.id]
  if (!selectedOption) {
    ElMessage.warning('请先选择一个选项')
    return
  }

  // 获取当前用户
  const user = userStore.user || userStore.getUser()
  console.log('[DEBUG] userStore.user:', userStore.user)
  console.log('[DEBUG] userStore.getUser():', userStore.getUser())
  console.log('[DEBUG] Final user:', user)
  console.log('[DEBUG] User ID:', user?.id)
  console.log('[DEBUG] Survey ID:', survey.id)
  console.log('[DEBUG] Selected option:', selectedOption)

  if (!user) {
    ElMessage.error('用户未登录，请先登录')
    return
  }

  if (!user.id || user.id === 0) {
    ElMessage.error('用户信息异常，请重新登录')
    return
  }

  voteLoading.value[survey.id] = true
  try {
    // Convert selectedOption to string for backend
    const optionString = String(selectedOption)
    console.log('[DEBUG] Calling API vote:', survey.id, optionString, user.id)
    const response = await api.vote(survey.id, optionString, user.id)
    console.log('[DEBUG] API response:', response)

    if (response.success) {
      ElMessage.success('投票成功！')
      votedSurveyIds.value.add(survey.id)
      survey.total_votes = (survey.total_votes || 0) + 1
    } else {
      ElMessage.error(response.message || '投票失败')
    }
  } catch (error) {
    console.error('[DEBUG] Vote error:', error)
    ElMessage.error('投票失败，请稍后重试')
  } finally {
    voteLoading.value[survey.id] = false
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

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    const validOptions = options.value.filter(o => o.trim() !== '')
    if (validOptions.length < 2) {
      ElMessage.warning('请至少填写两个有效选项')
      return
    }

    submitLoading.value = true
    try {
      const formatDate = (date: Date | string): string => {
        if (!date) return ''
        if (typeof date === 'string') return date
        const d = new Date(date)
        const year = d.getFullYear()
        const month = String(d.getMonth() + 1).padStart(2, '0')
        const day = String(d.getDate()).padStart(2, '0')
        return `${year}-${month}-${day}`
      }
      const startDate = formatDate(form.startDate)
      const endDate = formatDate(form.endDate)

      const response = await api.createSurvey(
        form.title,
        form.description,
        JSON.stringify(validOptions),
        startDate,
        endDate,
        (userStore.user || userStore.getUser())?.id || 1
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
  console.log('[DEBUG] View results for survey:', survey.id)
  try {
    const response = await api.getSurveyResults(survey.id)
    console.log('[DEBUG] GetSurveyResults response:', response)

    if (response.success) {
      currentResults.value = response.data
      showResultsDialog.value = true
      await nextTick()
      initChart()
    } else {
      ElMessage.error(response.message || '加载结果失败')
    }
  } catch (error) {
    console.error('[DEBUG] GetSurveyResults error:', error)
    ElMessage.error('加载结果失败')
  }
}

const initChart = () => {
  if (!chartRef.value || !currentResults.value || !currentSurvey.value) return
  if (!currentResults.value.total_votes || currentResults.value.total_votes === 0) return

  if (!chart) {
    chart = echarts.init(chartRef.value)
  }

  let optionList: string[] = []
  try {
    if (typeof currentSurvey.value.options === 'string') {
      optionList = JSON.parse(currentSurvey.value.options || '[]')
    } else if (Array.isArray(currentSurvey.value.options)) {
      optionList = currentSurvey.value.options
    }
  } catch (e) {
    optionList = []
  }

  const data = optionList.map((option: string) => ({
    name: option,
    value: currentResults.value.option_counts?.[option] || 0
  }))

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
  if (userStore.isAdmin()) {
    loadSurveys()
  } else {
    loadActiveSurveys()
  }
})

onUnmounted(() => {
  chart?.dispose()
})
</script>

<style scoped>
.survey-page {
  width: 100%;
}

.main-card {
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
}

.empty-card {
  border-radius: 8px;
}

.loading-container {
  height: 200px;
}

/* 调查项样式 */
.survey-item {
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 16px;
}

.survey-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.survey-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.survey-description {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #606266;
}

.survey-info {
  display: flex;
  gap: 20px;
  margin-bottom: 16px;
  padding: 10px 12px;
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 13px;
  color: #606266;
}

.voted-message {
  margin-top: 12px;
}

/* 投票区域 */
.vote-section {
  margin-top: 16px;
}

.vote-label {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
}

.vote-list {
  display: block;
  width: 100%;
}

.vote-item {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  margin-bottom: 10px;
  cursor: pointer;
  transition: all 0.2s;
  background: white;
}

.vote-item:hover {
  border-color: #409eff;
  background: #ecf5ff;
}

.vote-item.is-selected {
  border-color: #409eff;
  background: #ecf5ff;
}

.vote-item-content {
  display: flex;
  align-items: center;
  padding: 12px 16px;
}

.vote-index {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  height: 28px;
  background: #e4e7ed;
  border-radius: 4px;
  color: #606266;
  font-size: 13px;
  font-weight: 600;
  margin-right: 12px;
  flex-shrink: 0;
}

.vote-item.is-selected .vote-index {
  background: #409eff;
  color: white;
}

.vote-radio {
  margin: 0;
  flex: 1;
}

:deep(.vote-radio .el-radio__label) {
  font-size: 14px;
  color: #303133;
}

:deep(.vote-radio.is-checked .el-radio__label) {
  color: #409eff;
  font-weight: 600;
}

.vote-actions {
  margin-top: 16px;
  text-align: center;
}

.vote-actions .el-button {
  min-width: 150px;
}

/* 创建调查对话框样式 */
.option-hint {
  font-size: 12px;
  color: #909399;
  margin-bottom: 8px;
}

.option-input-row {
  margin-bottom: 10px;
}

.add-option-btn {
  width: 100%;
}

.chart-container {
  width: 100%;
  height: 400px;
  margin-top: 20px;
}

.result-summary {
  color: #606060;
  margin-bottom: 20px;
}
</style>
