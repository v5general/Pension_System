# 智慧养老系统

[![License](https://img.shields.io/badge/license-GPL--3.0-green.svg)](LICENSE)

[English](./README.md) | **简体中文**

基于 Wails v2 + Vue 3 + Element Plus 开发的桌面应用程序，用于老年人信息管理、民意调查和问题收集。

## 技术栈

### 后端
- Go 1.21+
- Wails v2
- GORM (SQLite)
- bcrypt 密码加密

### 前端
- Vue 3.4+
- TypeScript
- Element Plus UI
- ECharts 图表
- Pinia 状态管理
- Vue Router

## 项目结构

```
pension-system/
├── main.go                 # 应用入口
├── app.go                  # 应用上下文
├── go.mod                  # Go 模块配置
├── wails.json              # Wails 配置
├── database/               # 数据库配置
│   └── db.go
├── models/                 # 数据模型
│   ├── user.go            # 用户模型
│   ├── elderly.go         # 老年人模型
│   ├── survey.go          # 民意调查模型
│   └── issue.go           # 问题收集模型
├── controllers/            # 控制器
│   ├── auth.go            # 用户认证
│   ├── data.go            # 数据管理
│   ├── survey.go          # 民意调查
│   └── issue.go           # 问题收集
├── utils/                  # 工具函数
│   └── crypto.go          # 密码加密
└── frontend/              # 前端代码
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    ├── tsconfig.json
    └── src/
        ├── main.ts
        ├── App.vue
        ├── api/           # API 封装
        │   └── index.ts
        ├── router/        # 路由配置
        │   └── index.ts
        ├── store/         # 状态管理
        │   └── user.ts
        └── views/         # 页面组件
            ├── Login.vue
            ├── Home.vue
            ├── DataSummary.vue
            ├── DataManage.vue
            ├── DataForm.vue
            ├── Survey.vue
            └── Issue.vue
```

## 功能模块

### 1. 用户认证
- 用户登录/注册
- 角色权限管理（管理员、操作员、普通用户）
- 默认管理员账号: admin / admin123

### 2. 数据汇总
- 老年人总数统计
- 性别分布图表
- 年龄分布图表
- 健康状况分布
- 护理等级分布

### 3. 数据管理（管理员）
- 老年人信息 CRUD
- 数据搜索和分页
- 数据导入导出

### 4. 数据录入
- 老年人信息表单
- 表单验证
- 数据保存

### 5. 民意调查（管理员）
- 创建调查
- 发布/结束调查
- 查看投票结果
- 图表展示

### 6. 问题收集
- 提交问题
- 问题分类和优先级
- 问题回复
- 问题关闭

## 安装和运行

### 前置要求
1. Go 1.21 或更高版本
2. Node.js 18 或更高版本
3. Wails CLI

### 安装 Wails CLI
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 安装依赖

#### 后端依赖
```bash
cd D:\Workspace\PensionSystem\pension-system
go mod download
```

#### 前端依赖
```bash
cd frontend
npm install
```

### 开发模式运行
```bash
wails dev
```

### 构建应用
```bash
wails build
```

## 数据库

应用使用 SQLite 数据库，数据库文件 `pension_system.db` 会在首次运行时自动创建。

### 数据表
- `users` - 用户表
- `elderly` - 老年人信息表
- `surveys` - 民意调查表
- `survey_votes` - 调查投票表
- `issues` - 问题收集表

## 默认账号

- **用户名**: admin
- **密码**: admin123
- **角色**: 管理员

## 开发说明

### 添加新的 API 方法

1. 在 `controllers/` 目录下的相应控制器中添加方法
2. 在 `main.go` 的 `Bind` 中注册控制器
3. 在 `frontend/src/api/index.ts` 中添加 API 调用方法

### 添加新的页面

1. 在 `frontend/src/views/` 创建 Vue 组件
2. 在 `frontend/src/router/index.ts` 中添加路由
3. 在 `frontend/src/views/Home.vue` 中添加菜单项

## 许可证

本项目基于 [GPL-3.0](./LICENSE) 协议开源。
