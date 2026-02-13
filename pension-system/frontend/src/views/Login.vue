<template>
  <div class="login-container">
    <!-- Left side - Branding -->
    <div class="brand-section">
      <div class="brand-content">
        <div class="brand-icon">
          <el-icon :size="64"><UserFilled /></el-icon>
        </div>
        <h1 class="brand-title">智慧养老系统</h1>
        <p class="brand-subtitle">Smart Pension System</p>
        <p class="brand-description">专业的老年人信息管理平台</p>
        <div class="brand-features">
          <div class="feature-item feature-1">
            <el-icon><DataAnalysis /></el-icon>
            <span>数据分析</span>
          </div>
          <div class="feature-item feature-2">
            <el-icon><Document /></el-icon>
            <span>信息管理</span>
          </div>
          <div class="feature-item feature-3">
            <el-icon><ChatDotRound /></el-icon>
            <span>民意调查</span>
          </div>
          <div class="feature-item feature-4">
            <el-icon><QuestionFilled /></el-icon>
            <span>问题反馈</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Right side - Login form -->
    <div class="login-section">
      <div class="login-content">
        <div class="login-header">
          <h2>欢迎登录</h2>
          <p>请输入您的账号信息</p>
        </div>

        <el-form ref="loginFormRef" :model="loginForm" :rules="rules" class="login-form">
          <el-form-item prop="name">
            <el-input
              v-model="loginForm.name"
              placeholder="账户名"
              size="large"
              class="modern-input"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="密码"
              size="large"
              class="modern-input"
              @keyup.enter="handleLogin"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              class="login-button"
              @click="handleLogin"
            >
              <span v-if="!loading">登录</span>
              <span v-else>登录中...</span>
            </el-button>
          </el-form-item>

          <el-form-item>
            <div class="login-footer">
              <span>还没有账号？</span>
              <el-link type="primary" @click="showRegister = true">立即注册</el-link>
            </div>
          </el-form-item>
        </el-form>

        <div class="demo-info">
          <el-alert type="info" :closable="false">
            <template #default>
              <div class="demo-info-content">
                <el-icon><InfoFilled /></el-icon>
                <span>默认管理员账号: <strong>admin</strong> / <strong>admin123</strong></span>
              </div>
            </template>
          </el-alert>
        </div>
      </div>
    </div>

    <!-- Register Dialog -->
    <el-dialog v-model="showRegister" title="注册新账号" width="450px" class="modern-dialog">
      <el-form ref="registerFormRef" :model="registerForm" :rules="registerRules" label-position="top">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="registerForm.username" placeholder="请输入用户名">
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="账户名" prop="name">
          <el-input v-model="registerForm.name" placeholder="请输入账户名">
            <template #prefix>
              <el-icon><Avatar /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="registerForm.password" type="password" placeholder="请输入密码">
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="registerForm.confirmPassword" type="password" placeholder="请再次输入密码">
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showRegister = false" size="large">取消</el-button>
          <el-button type="primary" :loading="registerLoading" @click="handleRegister" size="large">
            注册
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import {
  User, Lock, UserFilled, DataAnalysis, Document,
  ChatDotRound, QuestionFilled, InfoFilled, Avatar
} from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import api from '@/api'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref<FormInstance>()
const registerFormRef = ref<FormInstance>()
const loading = ref(false)
const registerLoading = ref(false)
const showRegister = ref(false)

const loginForm = reactive({
  name: '',
  password: ''
})

const registerForm = reactive({
  username: '',
  name: '',
  password: '',
  confirmPassword: ''
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入账户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const registerRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度为3-20个字符', trigger: 'blur' }
  ],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (value !== registerForm.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return

  await loginFormRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const response = await api.login(loginForm.name, loginForm.password)
      if (response.success) {
        userStore.setUser(response.user)
        ElMessage.success('登录成功')
        router.push('/dashboard')
      } else {
        ElMessage.error(response.message)
      }
    } catch (error) {
      ElMessage.error('登录失败，请检查网络连接')
    } finally {
      loading.value = false
    }
  })
}

const handleRegister = async () => {
  if (!registerFormRef.value) return

  await registerFormRef.value.validate(async (valid) => {
    if (!valid) return

    registerLoading.value = true
    try {
      const response = await api.register(
        registerForm.username,
        registerForm.password,
        registerForm.name
      )
      if (response.success) {
        ElMessage.success('注册成功，请登录')
        showRegister.value = false
        loginForm.name = registerForm.name
        loginForm.password = ''
        registerFormRef.value?.resetFields()
      } else {
        ElMessage.error(response.message)
      }
    } catch (error) {
      ElMessage.error('注册失败，请检查网络连接')
    } finally {
      registerLoading.value = false
    }
  })
}
</script>

<style scoped>
.login-container {
  width: 100%;
  height: 100vh;
  display: flex;
  background: #f5f7fa;
}

.brand-section {
  flex: 1;
  background: linear-gradient(135deg, #0369a1 0%, #0891b2 25%, #0d9488 50%, #1d4ed8 75%, #4f46e5 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.brand-section::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
  animation: float 25s ease-in-out infinite;
}

.brand-section::after {
  content: '';
  position: absolute;
  bottom: -50%;
  right: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(167, 243, 208, 0.15) 0%, transparent 70%);
  animation: float 30s ease-in-out infinite reverse;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) rotate(0deg); }
  50% { transform: translate(20px, 20px) rotate(120deg); }
}

.brand-content {
  text-align: center;
  color: white;
  z-index: 1;
  padding: 60px;
}

.brand-icon {
  margin-bottom: 24px;
  opacity: 0.95;
}

.brand-title {
  font-size: 42px;
  font-weight: 600;
  margin: 0 0 12px 0;
  letter-spacing: 1px;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.15);
}

.brand-subtitle {
  font-size: 16px;
  margin: 0 0 8px 0;
  opacity: 0.95;
  letter-spacing: 1px;
}

.brand-description {
  font-size: 14px;
  margin: 0 0 48px 0;
  opacity: 0.85;
}

.brand-features {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  max-width: 320px;
  margin: 0 auto;
}

.feature-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 16px;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.feature-item:hover {
  background: rgba(255, 255, 255, 0.25);
  transform: translateY(-4px) scale(1.02);
  box-shadow: 0 12px 24px rgba(0, 96, 100, 0.15);
}

.feature-1 { background: rgba(13, 148, 136, 0.25); }
.feature-2 { background: rgba(8, 145, 178, 0.25); }
.feature-3 { background: rgba(6, 182, 212, 0.25); }
.feature-4 { background: rgba(29, 78, 216, 0.25); }

.feature-item .el-icon {
  font-size: 28px;
}

.feature-item span {
  font-size: 13px;
  font-weight: 500;
}

.login-section {
  width: 480px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #dbeafe 0%, #cbd5e1 50%, #e0e7ff 100%);
  padding: 40px;
}

.login-content {
  width: 100%;
  max-width: 400px;
}

.login-header {
  margin-bottom: 40px;
}

.login-header h2 {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #0369a1 0%, #0891b2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 8px 0;
}

.login-header p {
  font-size: 14px;
  color: #606060;
  margin: 0;
}

.login-form {
  margin-bottom: 24px;
}

:deep(.modern-input .el-input__wrapper) {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 2px solid transparent;
  transition: all 0.3s ease;
  background: white;
}

:deep(.modern-input .el-input__wrapper:hover) {
  box-shadow: 0 4px 12px rgba(8, 145, 178, 0.15);
  border-color: rgba(8, 145, 178, 0.2);
}

:deep(.modern-input .el-input__wrapper.is-focus) {
  box-shadow: 0 4px 16px rgba(3, 105, 161, 0.25);
  border-color: #0891b2;
}

.login-button {
  width: 100%;
  height: 48px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #0369a1 0%, #0891b2 100%);
  border: none;
  transition: all 0.3s ease;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(8, 145, 178, 0.4);
}

.login-footer {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: #606060;
}

.login-footer span {
  margin-right: 8px;
}

.demo-info {
  margin-top: 24px;
}

:deep(.demo-info .el-alert) {
  border-radius: 12px;
  border: none;
  background: linear-gradient(135deg, #e0f7fa, #fff);
}

.demo-info-content {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.demo-info-content strong {
  background: linear-gradient(135deg, #0369a1, #0891b2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: 700;
}

:deep(.modern-dialog) {
  border-radius: 16px;
}

:deep(.modern-dialog .el-dialog__header) {
  padding: 24px 24px 16px;
}

:deep(.modern-dialog .el-dialog__body) {
  padding: 16px 24px;
}

:deep(.modern-dialog .el-dialog__footer) {
  padding: 16px 24px 24px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.el-form-item__label) {
  font-weight: 600;
  color: #333;
}
</style>
