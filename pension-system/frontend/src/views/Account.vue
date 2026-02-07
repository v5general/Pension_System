<template>
  <div class="account">
    <el-row :gutter="20">
      <!-- 用户信息卡片 -->
      <el-col :span="8">
        <el-card class="profile-card modern-card" shadow="never">
          <div class="profile-header">
            <div class="profile-avatar">
              <el-icon :size="48"><Avatar /></el-icon>
            </div>
            <div class="profile-info">
              <h2 class="profile-name">{{ userStore.user?.name }}</h2>
              <p class="profile-role">{{ getRoleName(userStore.user?.role) }}</p>
            </div>
          </div>
          <div class="profile-stats">
            <div class="stat-item">
              <div class="stat-value">{{ formatDate(userStore.user?.created_at) }}</div>
              <div class="stat-label">注册时间</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 账户操作卡片 -->
      <el-col :span="16">
        <el-card class="actions-card modern-card" shadow="never">
          <template #header>
            <div class="card-title-wrapper">
              <el-icon><Setting /></el-icon>
              <span>账户设置</span>
            </div>
          </template>

          <!-- 修改用户名 -->
          <div class="action-section">
            <h3 class="section-title">
              <el-icon><Edit /></el-icon>
              修改用户名
            </h3>
            <p class="section-desc">修改您在系统中显示的用户名</p>
            <el-form ref="nameFormRef" :model="nameForm" :rules="nameRules" label-width="100px" class="action-form">
              <el-form-item label="新用户名" prop="newName">
                <el-input v-model="nameForm.newName" placeholder="请输入新用户名" />
              </el-form-item>
              <el-form-item label="确认密码" prop="password">
                <el-input v-model="nameForm.password" type="password" placeholder="请输入当前密码" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="nameLoading" @click="handleUpdateName">保存修改</el-button>
              </el-form-item>
            </el-form>
          </div>

          <el-divider />

          <!-- 修改密码 -->
          <div class="action-section">
            <h3 class="section-title">
              <el-icon><Lock /></el-icon>
              修改密码
            </h3>
            <p class="section-desc">定期修改密码以保护账号安全</p>
            <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="100px" class="action-form">
              <el-form-item label="当前密码" prop="currentPassword">
                <el-input v-model="passwordForm.currentPassword" type="password" placeholder="请输入当前密码" />
              </el-form-item>
              <el-form-item label="新密码" prop="newPassword">
                <el-input v-model="passwordForm.newPassword" type="password" placeholder="请输入新密码（至少6位）" />
              </el-form-item>
              <el-form-item label="确认密码" prop="confirmPassword">
                <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="请再次输入新密码" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="passwordLoading" @click="handleUpdatePassword">保存修改</el-button>
              </el-form-item>
            </el-form>
          </div>

          <el-divider />

          <!-- 注销账号 -->
          <div class="action-section danger-section">
            <h3 class="section-title danger-title">
              <el-icon><WarningFilled /></el-icon>
              注销账号
            </h3>
            <p class="section-desc">注销后您的所有数据将被永久删除，无法恢复！</p>
            <el-alert type="error" :closable="false" show-icon class="delete-warning">
              <template #default>
                <div>注销账号将：</div>
                <ul style="margin: 8px 0 0 20px; padding: 0;">
                  <li>删除您的账号信息</li>
                  <li>清除您提交的问题记录</li>
                  <li>清除您的投票数据</li>
                  <li>您将无法再登录系统</li>
                </ul>
              </template>
            </el-alert>
            <el-form ref="deleteFormRef" :model="deleteForm" :rules="deleteRules" label-width="100px" class="action-form" style="margin-top: 20px;">
              <el-form-item label="确认密码" prop="password">
                <el-input v-model="deleteForm.password" type="password" placeholder="请输入当前密码" />
              </el-form-item>
              <el-form-item label="确认文本" prop="confirmText">
                <el-input v-model="deleteForm.confirmText" placeholder="请输入「确认注销」以确认" />
              </el-form-item>
              <el-form-item>
                <el-button type="danger" :loading="deleteLoading" @click="handleDeleteAccount">注销账号</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
import { Avatar, Setting, Edit, Lock, WarningFilled } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import api from '@/api'

const router = useRouter()
const userStore = useUserStore()

const nameFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()
const deleteFormRef = ref<FormInstance>()

const nameLoading = ref(false)
const passwordLoading = ref(false)
const deleteLoading = ref(false)

const nameForm = reactive({
  newName: '',
  password: ''
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const deleteForm = reactive({
  password: '',
  confirmText: ''
})

const nameRules: FormRules = {
  newName: [
    { required: true, message: '请输入新用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度为2-20个字符', trigger: 'blur' }
  ],
  password: [{ required: true, message: '请输入当前密码', trigger: 'blur' }]
}

const passwordRules: FormRules = {
  currentPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const deleteRules: FormRules = {
  password: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  confirmText: [
    { required: true, message: '请输入确认文本', trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (value !== '确认注销') {
          callback(new Error('请输入「确认注销」以确认操作'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const getRoleName = (role?: string) => {
  const roleMap: Record<string, string> = {
    admin: '管理员',
    operator: '操作员',
    user: '普通用户'
  }
  return roleMap[role || ''] || '用户'
}

const formatDate = (date?: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const handleUpdateName = async () => {
  if (!nameFormRef.value) return

  await nameFormRef.value.validate(async (valid) => {
    if (!valid) return

    nameLoading.value = true
    try {
      const response = await api.updateName(nameForm.newName, nameForm.password, userStore.user?.id || 1)
      if (response.success) {
        userStore.setUser(response.user)
        ElMessage.success('用户名修改成功')
        nameFormRef.value?.resetFields()
      } else {
        ElMessage.error(response.message)
      }
    } catch (error) {
      ElMessage.error('修改失败，请检查网络连接')
    } finally {
      nameLoading.value = false
    }
  })
}

const handleUpdatePassword = async () => {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate(async (valid) => {
    if (!valid) return

    passwordLoading.value = true
    try {
      const response = await api.updatePassword(passwordForm.currentPassword, passwordForm.newPassword, userStore.user?.id || 1)
      if (response.success) {
        ElMessage.success('密码修改成功，请重新登录')
        passwordFormRef.value?.resetFields()
        setTimeout(() => {
          userStore.logout()
          router.push('/login')
        }, 1500)
      } else {
        ElMessage.error(response.message)
      }
    } catch (error) {
      ElMessage.error('修改失败，请检查网络连接')
    } finally {
      passwordLoading.value = false
    }
  })
}

const handleDeleteAccount = async () => {
  if (!deleteFormRef.value) return

  await deleteFormRef.value.validate(async (valid) => {
    if (!valid) return

    await ElMessageBox.confirm(
      '此操作将永久注销您的账号，所有数据将被删除且无法恢复！确定要继续吗？',
      '危险操作确认',
      {
        confirmButtonText: '确定注销',
        cancelButtonText: '取消',
        type: 'error',
        confirmButtonClass: 'el-button--danger'
      }
    ).catch(() => {
      ElMessage.info('已取消注销')
      return
    })

    deleteLoading.value = true
    try {
      const response = await api.deleteAccount(deleteForm.password, userStore.user?.id || 1)
      if (response.success) {
        ElMessage.success('账号已注销')
        userStore.logout()
        router.push('/login')
      } else {
        ElMessage.error(response.message)
      }
    } catch (error) {
      ElMessage.error('注销失败，请检查网络连接')
    } finally {
      deleteLoading.value = false
    }
  })
}
</script>

<style scoped>
.account {
  width: 100%;
  padding: 10px;
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

/* 用户信息卡片 */
.profile-card {
  background: linear-gradient(135deg, #0369a1 0%, #0891b2 100%);
  border: none;
}

:deep(.profile-card .el-card__body) {
  padding: 30px;
}

.profile-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 24px;
}

.profile-avatar {
  width: 80px;
  height: 80px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-bottom: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.profile-info {
  text-align: center;
}

.profile-name {
  color: white;
  font-size: 20px;
  font-weight: 700;
  margin: 0 0 8px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.profile-role {
  color: rgba(255, 255, 255, 0.85);
  font-size: 14px;
  margin: 0;
}

.profile-stats {
  display: flex;
  justify-content: center;
}

.stat-item {
  text-align: center;
}

.stat-value {
  color: white;
  font-size: 14px;
  font-weight: 600;
}

.stat-label {
  color: rgba(255, 255, 255, 0.75);
  font-size: 12px;
  margin-top: 4px;
}

/* 操作卡片 */
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

.action-section {
  padding: 10px 0;
}

.action-section:first-child {
  padding-top: 0;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #0369a1 0%, #0891b2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 8px 0;
}

.section-desc {
  color: #666;
  font-size: 13px;
  margin: 0 0 16px 0;
}

.action-form {
  max-width: 500px;
}

:deep(.action-form .el-form-item__label) {
  font-weight: 500;
  color: #333;
}

:deep(.action-form .el-input__wrapper) {
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.action-form .el-input__wrapper:hover) {
  box-shadow: 0 2px 8px rgba(8, 145, 178, 0.15);
}

:deep(.action-form .el-input__wrapper.is-focus) {
  box-shadow: 0 2px 12px rgba(3, 105, 161, 0.2);
}

:deep(.action-form .el-button--primary) {
  background: linear-gradient(135deg, #0369a1 0%, #0891b2 100%);
  border: none;
  border-radius: 8px;
  padding: 10px 24px;
}

:deep(.action-form .el-button--primary:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(8, 145, 178, 0.3);
}

/* 危险区域 */
.danger-section {
  margin-top: 30px;
  padding-top: 20px;
}

.danger-title {
  background: linear-gradient(135deg, #dc2626 0%, #ef4444 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.delete-warning {
  margin: 16px 0;
}

:deep(.delete-warning .el-alert__content) {
  font-size: 13px;
}

:deep(.action-form .el-button--danger) {
  background: linear-gradient(135deg, #dc2626 0%, #ef4444 100%);
  border: none;
  border-radius: 8px;
  padding: 10px 24px;
}

:deep(.action-form .el-button--danger:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.3);
}

:deep(.el-divider) {
  margin: 24px 0;
}
</style>
