<template>
  <div class="issue">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>问题收集</span>
          <el-button type="primary" @click="showCreateDialog = true">提交问题</el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="全部" name="all"></el-tab-pane>
        <el-tab-pane label="待处理" name="pending"></el-tab-pane>
        <el-tab-pane label="处理中" name="processing"></el-tab-pane>
        <el-tab-pane label="已解决" name="resolved"></el-tab-pane>
        <el-tab-pane label="已关闭" name="closed"></el-tab-pane>
      </el-tabs>

      <el-table :data="issues" border style="width: 100%; margin-top: 20px" v-loading="loading">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="title" label="问题标题" min-width="200" />
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag>{{ row.category }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" width="100">
          <template #default="{ row }">
            <el-tag :type="getPriorityType(row.priority)">{{ getPriorityText(row.priority) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="submitter.name" label="提交人" width="100" />
        <el-table-column prop="created_at" label="提交时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewDetail(row)">详情</el-button>
            <el-button
              v-if="row.status !== 'closed'"
              type="success"
              size="small"
              @click="respondToIssue(row)"
            >
              回复
            </el-button>
            <el-button
              v-if="row.status !== 'closed' && isAdmin"
              type="danger"
              size="small"
              @click="closeIssue(row)"
            >
              关闭
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        style="margin-top: 20px"
        @current-change="loadIssues"
      />
    </el-card>

    <!-- Create Issue Dialog -->
    <el-dialog v-model="showCreateDialog" title="提交问题" width="600px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="问题标题" prop="title">
          <el-input v-model="form.title" placeholder="请简要描述问题" />
        </el-form-item>
        <el-form-item label="问题分类" prop="category">
          <el-select v-model="form.category" placeholder="请选择分类" style="width: 100%">
            <el-option label="生活服务" value="生活服务" />
            <el-option label="医疗健康" value="医疗健康" />
            <el-option label="文化娱乐" value="文化娱乐" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-radio-group v-model="form.priority">
            <el-radio label="low">低</el-radio>
            <el-radio label="normal">普通</el-radio>
            <el-radio label="high">高</el-radio>
            <el-radio label="urgent">紧急</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="详细描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="5" placeholder="请详细描述您遇到的问题" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleCreate">提交</el-button>
      </template>
    </el-dialog>

    <!-- Detail Dialog -->
    <el-dialog v-model="showDetailDialog" title="问题详情" width="700px">
      <div v-if="currentIssue">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="标题">{{ currentIssue.title }}</el-descriptions-item>
          <el-descriptions-item label="分类">{{ currentIssue.category }}</el-descriptions-item>
          <el-descriptions-item label="优先级">
            <el-tag :type="getPriorityType(currentIssue.priority)">{{ getPriorityText(currentIssue.priority) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(currentIssue.status)">{{ getStatusText(currentIssue.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="提交人">{{ currentIssue.submitter?.name }}</el-descriptions-item>
          <el-descriptions-item label="提交时间">{{ formatDateTime(currentIssue.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="描述">{{ currentIssue.description }}</el-descriptions-item>
          <el-descriptions-item v-if="currentIssue.response" label="回复">
            {{ currentIssue.response }}
          </el-descriptions-item>
          <el-descriptions-item v-if="currentIssue.handler" label="处理人">
            {{ currentIssue.handler.name }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <!-- Response Dialog -->
    <el-dialog v-model="showResponseDialog" title="回复问题" width="600px">
      <el-form ref="responseFormRef" :model="responseForm" :rules="responseRules" label-width="100px">
        <el-form-item label="回复内容" prop="response">
          <el-input v-model="responseForm.response" type="textarea" :rows="5" placeholder="请输入回复内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showResponseDialog = false">取消</el-button>
        <el-button type="primary" :loading="responseLoading" @click="handleRespond">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
import { useUserStore } from '@/store/user'
import api from '@/api'

const userStore = useUserStore()
const isAdmin = computed(() => userStore.isAdmin())

const loading = ref(false)
const submitLoading = ref(false)
const responseLoading = ref(false)
const issues = ref<any[]>([])
const activeTab = ref('all')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const showResponseDialog = ref(false)
const currentIssue = ref<any>(null)

const formRef = ref<FormInstance>()
const responseFormRef = ref<FormInstance>()

const form = reactive({
  title: '',
  category: '',
  priority: 'normal',
  description: ''
})

const responseForm = reactive({
  response: ''
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入问题标题', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  description: [{ required: true, message: '请详细描述问题', trigger: 'blur' }]
}

const responseRules: FormRules = {
  response: [{ required: true, message: '请输入回复内容', trigger: 'blur' }]
}

const loadIssues = async () => {
  loading.value = true
  try {
    const response = await api.getIssueList(currentPage.value, pageSize.value, activeTab.value)
    if (response.success) {
      issues.value = response.data || []
      total.value = response.total || 0
    }
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const handleTabChange = () => {
  currentPage.value = 1
  loadIssues()
}

const formatDateTime = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

const getPriorityType = (priority: string) => {
  const map: Record<string, string> = {
    low: 'info',
    normal: '',
    high: 'warning',
    urgent: 'danger'
  }
  return map[priority] || ''
}

const getPriorityText = (priority: string) => {
  const map: Record<string, string> = {
    low: '低',
    normal: '普通',
    high: '高',
    urgent: '紧急'
  }
  return map[priority] || priority
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'info',
    processing: 'warning',
    resolved: 'success',
    closed: 'danger'
  }
  return map[status] || ''
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理',
    processing: '处理中',
    resolved: '已解决',
    closed: '已关闭'
  }
  return map[status] || status
}

const handleCreate = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const userId = userStore.user?.id || 1
      const response = await api.createIssue(
        form.title,
        form.description,
        form.category,
        form.priority,
        userId
      )

      if (response.success) {
        ElMessage.success('提交成功')
        showCreateDialog.value = false
        formRef.value?.resetFields()
        loadIssues()
      } else {
        ElMessage.error(response.message)
      }
    } catch (error) {
      ElMessage.error('提交失败')
    } finally {
      submitLoading.value = false
    }
  })
}

const viewDetail = (issue: any) => {
  currentIssue.value = issue
  showDetailDialog.value = true
}

const respondToIssue = (issue: any) => {
  currentIssue.value = issue
  responseForm.response = ''
  showResponseDialog.value = true
}

const handleRespond = async () => {
  if (!responseFormRef.value) return

  await responseFormRef.value.validate(async (valid) => {
    if (!valid) return

    responseLoading.value = true
    try {
      const handlerId = userStore.user?.id || 1
      const response = await api.respondToIssue(
        currentIssue.value.id,
        responseForm.response,
        handlerId
      )

      if (response.success) {
        ElMessage.success('回复成功')
        showResponseDialog.value = false
        loadIssues()
      } else {
        ElMessage.error(response.message)
      }
    } catch (error) {
      ElMessage.error('回复失败')
    } finally {
      responseLoading.value = false
    }
  })
}

const closeIssue = (issue: any) => {
  ElMessageBox.confirm('确定要关闭这个问题吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const response = await api.closeIssue(issue.id)
      if (response.success) {
        ElMessage.success('已关闭')
        loadIssues()
      }
    } catch (error) {
      ElMessage.error('操作失败')
    }
  }).catch(() => {})
}

// Initial load
loadIssues()
</script>

<style scoped>
.issue {
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
