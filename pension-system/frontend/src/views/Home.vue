<template>
  <el-container class="layout-container">
    <el-aside width="200px" class="sidebar">
      <div class="logo">
        <h2>智慧养老</h2>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <span>数据汇总</span>
        </el-menu-item>
        <el-menu-item index="/data-form">
          <el-icon><Edit /></el-icon>
          <span>数据录入</span>
        </el-menu-item>
        <el-menu-item index="/survey" v-if="userStore.isAdmin()">
          <el-icon><Poll /></el-icon>
          <span>民意调查</span>
        </el-menu-item>
        <el-menu-item index="/issue">
          <el-icon><ChatDotRound /></el-icon>
          <span>问题收集</span>
        </el-menu-item>
        <el-menu-item index="/data-manage" v-if="userStore.isAdmin()">
          <el-icon><Management /></el-icon>
          <span>数据管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <span class="page-title">{{ currentPageTitle }}</span>
        </div>
        <div class="header-right">
          <span class="user-info">{{ userStore.user?.name }}</span>
          <el-dropdown @command="handleCommand">
            <el-icon class="user-icon"><Avatar /></el-icon>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
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
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const currentPageTitle = computed(() => route.meta.title as string || '智慧养老系统')

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
</script>

<style scoped>
.layout-container {
  width: 100%;
  height: 100vh;
}

.sidebar {
  background-color: #304156;
  overflow-x: hidden;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #2b3a4a;
}

.logo h2 {
  color: #fff;
  font-size: 20px;
  margin: 0;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
}

.page-title {
  font-size: 18px;
  font-weight: 500;
  color: #333;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-info {
  font-size: 14px;
  color: #666;
}

.user-icon {
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.main-content {
  background-color: #f0f2f5;
  overflow-y: auto;
  padding: 20px;
}
</style>
