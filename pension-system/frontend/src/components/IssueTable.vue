<template>
  <div>
    <el-table :data="issues" border style="width: 100%" v-loading="loading" :row-class-name="tableRowClassName">
      <el-table-column label="序号" width="80" align="center">
        <template #default="{ $index }">
          {{ (currentPage - 1) * pageSize + $index + 1 }}
        </template>
      </el-table-column>
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
      <el-table-column label="操作" width="280" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="$emit('viewDetail', row)">详情</el-button>
          <el-button
            v-if="row.status !== 'closed'"
            type="success"
            size="small"
            @click="$emit('respondIssue', row)"
          >
            回复
          </el-button>
          <el-button
            v-if="row.status !== 'closed' && isAdmin"
            type="warning"
            size="small"
            @click="$emit('closeIssue', row)"
          >
            关闭
          </el-button>
          <el-button
            v-if="isAdmin"
            type="danger"
            size="small"
            @click="$emit('deleteIssue', row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      :current-page="currentPage"
      :page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      style="margin-top: 20px"
      @current-change="$emit('pageChange')"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useUserStore } from '@/store/user'

interface Props {
  issues: any[]
  loading: boolean
  currentPage: number
  pageSize: number
  total: number
}

defineProps<Props>()
defineEmits<{
  viewDetail: [issue: any]
  respondIssue: [issue: any]
  closeIssue: [issue: any]
  deleteIssue: [issue: any]
  pageChange: []
}>()

const userStore = useUserStore()
const isAdmin = computed(() => userStore.isAdmin())

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

const tableRowClassName = ({ rowIndex }: { rowIndex: number }) => {
  return rowIndex % 2 === 0 ? 'even-row' : 'odd-row'
}
</script>

<style scoped>
:deep(.el-table) {
  border-radius: 4px;
  overflow: hidden;
}

:deep(.el-table tr.even-row) {
  background: #fafafa;
}

:deep(.el-table tr.odd-row) {
  background: #ffffff;
}

:deep(.el-table tbody tr:hover) {
  background: #ecf5ff !important;
}

:deep(.el-tag) {
  border-radius: 4px;
}
</style>
