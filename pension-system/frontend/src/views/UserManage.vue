<template>
  <div class="user-manage">
    <el-card class="modern-card">
      <template #header>
        <div class="card-header">
          <div class="card-title-wrapper">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </div>
        </div>
      </template>

      <el-table :data="users" border style="width: 100%" v-loading="loading">
        <el-table-column label="序号" width="80" align="center">
          <template #default="{ $index }">
            {{ (currentPage - 1) * pageSize + $index + 1 }}
          </template>
        </el-table-column>
        <el-table-column prop="name" label="用户名" min-width="120" />
        <el-table-column prop="username" label="账号" min-width="120" />
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="row.role === 'admin' ? 'danger' : 'primary'">
              {{ getRoleName(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="canModifyUser(row)"
              type="primary"
              size="small"
              @click="handleResetPassword(row)"
            >
              重置密码
            </el-button>
            <el-button
              v-if="canDeleteUser(row)"
              type="danger"
              size="small"
              @click="handleDeleteUser(row)"
            >
              注销
            </el-button>
            <span v-if="!canModifyUser(row) && !canDeleteUser(row)" style="color: #999; font-size: 12px;">-</span>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        style="margin-top: 20px"
        @current-change="loadUsers"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import api from '@/api'

const userStore = useUserStore()
const loading = ref(false)
const users = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const getRoleName = (role: string) => {
  const roleMap: Record<string, string> = {
    admin: '管理员',
    operator: '操作员',
    user: '普通用户'
  }
  return roleMap[role] || '用户'
}

const formatDateTime = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

const canDeleteUser = (row: any) => {
  // Cannot delete self
  if (row.id === userStore.user?.id) return false
  return true
}

const canModifyUser = (row: any) => {
  // Cannot modify self (use account settings page)
  if (row.id === userStore.user?.id) return false
  return true
}

const loadUsers = async () => {
  loading.value = true
  try {
    const response = await api.getUserList(currentPage.value, pageSize.value)
    if (response.success) {
      users.value = response.data || []
      total.value = response.total || 0
    }
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const handleDeleteUser = (user: any) => {
  ElMessageBox.confirm(
    `确定要注销用户「${user.name}」的账号吗？注销后该用户将无法登录系统，此操作不可恢复！`,
    '注销用户账号',
    {
      confirmButtonText: '确定注销',
      cancelButtonText: '取消',
      type: 'error',
      confirmButtonClass: 'el-button--danger'
    }
  ).then(async () => {
    try {
      const response = await api.deleteUserAccount(user.id, userStore.user?.id || 1)
      if (response.success) {
        ElMessage.success('用户账号已注销')
        loadUsers()
      } else {
        ElMessage.error(response.message)
      }
    } catch (error) {
      ElMessage.error('注销失败')
    }
  }).catch(() => {})
}

const handleResetPassword = (user: any) => {
  ElMessageBox.confirm(
    `确定要将用户「${user.name}」的密码重置为 123456 吗？`,
    '重置密码',
    {
      confirmButtonText: '确定重置',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const response = await api.resetUserPassword(user.id)
      if (response.success) {
        ElMessage.success(response.message || '密码已重置为：123456')
      } else {
        ElMessage.error(response.message)
      }
    } catch (error) {
      ElMessage.error('重置失败，请检查网络连接')
    }
  }).catch(() => {})
}

// Initial load
loadUsers()
</script>

<style scoped>
.user-manage {
  width: 100%;
}

.modern-card {
  border-radius: 16px;
  border: 1px solid rgba(0, 0, 0, 0.05);
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

:deep(.modern-card .el-card__header) {
  padding: 18px 22px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  background: linear-gradient(to bottom, #fafafa, #ffffff);
}

:deep(.modern-card .el-card__body) {
  padding: 22px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #0369a1 0%, #0d9488 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

:deep(.el-button--danger) {
  background: linear-gradient(135deg, #dc2626 0%, #ef4444 100%);
  border: none;
  border-radius: 8px;
}

:deep(.el-button--danger:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.3);
}

:deep(.el-button--primary) {
  background: linear-gradient(135deg, #0369a1 0%, #0891b2 100%);
  border: none;
  border-radius: 8px;
}

:deep(.el-button--primary:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(8, 145, 178, 0.3);
}
</style>
