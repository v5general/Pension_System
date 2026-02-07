<template>
  <el-container class="layout-container">
    <el-aside width="260px" class="sidebar">
      <div class="logo-section">
        <div class="logo-icon">
          <el-icon :size="36"><UserFilled /></el-icon>
        </div>
        <div class="logo-text">
          <h2>智慧养老系统</h2>
          <p>Pension System</p>
        </div>
      </div>

      <el-menu
        :default-active="activeMenu"
        router
        class="sidebar-menu"
      >
        <el-menu-item index="/dashboard">
          <template #title>
            <div class="menu-item-content">
              <el-icon><DataBoard /></el-icon>
              <span>数据汇总</span>
            </div>
          </template>
        </el-menu-item>
        <el-menu-item index="/data-form">
          <template #title>
            <div class="menu-item-content">
              <el-icon><DocumentAdd /></el-icon>
              <span>数据录入</span>
            </div>
          </template>
        </el-menu-item>
        <el-menu-item index="/survey" v-if="userStore.isAdmin()">
          <template #title>
            <div class="menu-item-content">
              <el-icon><DataAnalysis /></el-icon>
              <span>民意调查</span>
            </div>
          </template>
        </el-menu-item>
        <el-menu-item index="/issue">
          <template #title>
            <div class="menu-item-content">
              <el-icon><ChatLineSquare /></el-icon>
              <span>问题收集</span>
            </div>
          </template>
        </el-menu-item>
        <el-menu-item index="/data-manage" v-if="userStore.isAdmin()">
          <template #title>
            <div class="menu-item-content">
              <el-icon><Management /></el-icon>
              <span>数据管理</span>
            </div>
          </template>
        </el-menu-item>
        <el-menu-item index="/user-manage" v-if="userStore.isAdmin()">
          <template #title>
            <div class="menu-item-content">
              <el-icon><User /></el-icon>
              <span>用户管理</span>
            </div>
          </template>
        </el-menu-item>
      </el-menu>

      <div class="sidebar-footer">
        <div class="user-mini-card" @click="goToAccount" style="cursor: pointer">
          <div class="user-avatar">
            <el-icon><Avatar /></el-icon>
          </div>
          <div class="user-info">
            <div class="user-name">{{ userStore.user?.name }}</div>
            <div class="user-role">{{ getRoleName(userStore.user?.role) }}</div>
          </div>
          <div class="user-arrow">
            <el-icon><ArrowRight /></el-icon>
          </div>
        </div>
      </div>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <h1 class="page-title">{{ currentPageTitle }}</h1>
        </div>
        <div class="header-right">
          <el-button text @click="handleCommand('logout')" class="logout-btn">
            <el-icon><SwitchButton /></el-icon>
            <span>退出登录</span>
          </el-button>
        </div>
      </el-header>

      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  DataBoard, DocumentAdd, DataAnalysis, ChatLineSquare,
  Management, UserFilled, Avatar, SwitchButton, ArrowRight, User
} from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const currentPageTitle = computed(() => route.meta.title as string || '智慧养老系统')

const getRoleName = (role?: string) => {
  const roleMap: Record<string, string> = {
    admin: '管理员',
    operator: '操作员',
    user: '普通用户'
  }
  return roleMap[role || ''] || '用户'
}

const handleCommand = (command: string) => {
  if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      userStore.logout()
      ElMessage.success('已退出登录')
      router.push('/login')
    }).catch(() => {})
  }
}

const goToAccount = () => {
  router.push('/account')
}
</script>

<style scoped>
.layout-container {
  width: 100%;
  height: 100vh;
}

.sidebar {
  background: linear-gradient(180deg, #0369a1 0%, #0891b2 20%, #0d9488 40%, #06b6d4 60%, #14b8a6 80%, #0d9488 100%);
  display: flex;
  flex-direction: column;
  border-right: 1px solid rgba(255, 255, 255, 0.12);
  position: relative;
  overflow: hidden;
}

.sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background:
    radial-gradient(circle at 30% 20%, rgba(255, 255, 255, 0.15) 0%, transparent 50%),
    radial-gradient(circle at 70% 80%, rgba(167, 243, 208, 0.2) 0%, transparent 50%);
  pointer-events: none;
}

.logo-section {
  padding: 28px 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  position: relative;
  z-index: 1;
}

.logo-icon {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  background: linear-gradient(135deg, #0891b2 0%, #06b6d4 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
  box-shadow: 0 6px 16px rgba(8, 145, 178, 0.4);
}

.logo-text h2 {
  color: white;
  font-size: 18px;
  font-weight: 700;
  margin: 0 0 4px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.logo-text p {
  color: rgba(255, 255, 255, 0.95);
  font-size: 11px;
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 1px;
}

:deep(.sidebar-menu) {
  background: transparent;
  border: none;
  flex: 1;
  padding: 20px 12px;
  position: relative;
  z-index: 1;
}

:deep(.sidebar-menu .el-menu-item) {
  background: transparent;
  color: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  margin-bottom: 8px;
  padding: 0 16px !important;
  height: 48px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

:deep(.sidebar-menu .el-menu-item::before) {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 4px;
  height: 100%;
  background: linear-gradient(180deg, #06b6d4, #0891b2);
  border-radius: 4px 0 0 4px;
  transform: scaleY(0);
  transition: transform 0.3s ease;
}

:deep(.sidebar-menu .el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  transform: translateX(4px);
}

:deep(.sidebar-menu .el-menu-item:hover::before) {
  transform: scaleY(1);
}

:deep(.sidebar-menu .el-menu-item.is-active) {
  background: rgba(255, 255, 255, 0.25);
  color: white;
  box-shadow: 0 6px 20px rgba(8, 145, 178, 0.25);
  transform: translateX(4px);
}

:deep(.sidebar-menu .el-menu-item.is-active::before) {
  transform: scaleY(1);
}

:deep(.sidebar-menu .el-menu-item .el-icon) {
  margin-right: 10px;
  font-size: 18px;
}

.menu-item-content {
  display: flex;
  align-items: center;
  width: 100%;
}

.sidebar-footer {
  padding: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
  position: relative;
  z-index: 1;
}

.user-mini-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 14px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.user-mini-card:hover {
  background: rgba(255, 255, 255, 0.25);
  transform: translateX(4px);
  box-shadow: 0 6px 20px rgba(3, 105, 161, 0.3);
}

.user-arrow {
  color: rgba(255, 255, 255, 0.7);
  font-size: 16px;
  transition: all 0.3s ease;
}

.user-mini-card:hover .user-arrow {
  color: white;
  transform: translateX(4px);
}

.user-avatar {
  width: 42px;
  height: 42px;
  border-radius: 12px;
  background: linear-gradient(135deg, #0284c7 0%, #0369a1 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(3, 105, 161, 0.35);
}

.user-avatar .el-icon {
  font-size: 22px;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  color: white;
  font-size: 14px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.user-role {
  color: rgba(255, 255, 255, 0.85);
  font-size: 12px;
  margin-top: 2px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(90deg, #e0e7ff 0%, #dbeafe 50%, #e2e8f0 100%);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  padding: 0 32px;
  height: 64px;
}

.header-left {
  display: flex;
  align-items: center;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, #0369a1 0%, #0891b2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.logout-btn {
  color: #606060;
  font-size: 14px;
  padding: 10px 18px;
  border-radius: 10px;
  transition: all 0.3s ease;
  font-weight: 500;
}

.logout-btn:hover {
  background: linear-gradient(135deg, #dc2626, #ef4444);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.3);
}

.logout-btn .el-icon {
  margin-right: 6px;
}

.main-content {
  background: linear-gradient(135deg, #e0e7ff 0%, #dbeafe 25%, #e2e8f0 50%, #cbd5e1 75%, #e0e7ff 100%);
  background-size: 100% 100%;
  overflow-y: auto;
  padding: 24px;
}
</style>
