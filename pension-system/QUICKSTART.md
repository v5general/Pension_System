# 智慧养老系统 - 快速开始指南

## 项目概述

本项目是一个基于 Wails v2 + Vue 3 开发的智慧养老管理系统，提供老年人信息管理、数据汇总统计、民意调查和问题收集等功能。

## 系统要求

- **Go**: 1.21 或更高版本
- **Node.js**: 18 或更高版本
- **Wails CLI**: 最新版本
- **操作系统**: Windows 10/11, macOS, 或 Linux

## 快速开始

### 1. 安装 Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 2. 安装项目依赖

#### Windows 用户
双击运行 `install.bat` 文件

#### macOS/Linux 用户
```bash
# 安装 Go 依赖
go mod download

# 安装前端依赖
cd frontend
npm install
cd ..
```

### 3. 启动开发服务器

#### Windows 用户
双击运行 `build.bat` 文件

#### macOS/Linux 用户
```bash
wails dev
```

应用将自动编译并在默认浏览器中打开。

### 4. 构建生产版本

```bash
wails build
```

构建完成后，可执行文件将位于 `build/bin/` 目录中。

## 默认登录信息

```
用户名: admin
密码: admin123
```

首次登录后，建议立即修改默认密码。

## 主要功能

### 1. 数据汇总（Dashboard）
- 老年人总数统计
- 性别分布图表
- 年龄分布图表
- 健康状况分布
- 护理等级分布

### 2. 数据录入
- 老年人信息录入表单
- 支持基本信息、健康状况、经济来源等多个字段
- 表单验证和自动保存

### 3. 数据管理（管理员权限）
- 老年人信息的增删改查
- 数据搜索和筛选
- 分页展示

### 4. 民意调查（管理员权限）
- 创建和管理调查
- 设置调查选项和时间
- 查看投票结果和统计图表

### 5. 问题收集
- 提交问题和反馈
- 问题分类和优先级设置
- 问题回复和跟踪

## 项目结构

```
pension-system/
├── main.go              # 应用入口
├── app.go              # 应用上下文
├── controllers/        # 业务逻辑控制器
├── models/            # 数据模型
├── database/          # 数据库配置
├── utils/             # 工具函数
└── frontend/          # 前端代码
    └── src/
        ├── api/       # API 封装
        ├── views/     # 页面组件
        ├── router/    # 路由配置
        └── store/     # 状态管理
```

## 开发说明

### 添加新的 API 方法

1. 在 `controllers/` 中添加控制器方法
2. 在 `main.go` 的 `Bind` 中注册
3. 在 `frontend/src/api/index.ts` 中添加前端调用

### 添加新的页面

1. 在 `frontend/src/views/` 创建 Vue 组件
2. 在 `frontend/src/router/index.ts` 添加路由
3. 在 `frontend/src/views/Home.vue` 添加菜单

## 常见问题

### Q: 编译时出现 Go 模块错误
A: 运行 `go mod tidy` 清理依赖

### Q: 前端资源加载失败
A: 删除 `frontend/node_modules` 后重新运行 `npm install`

### Q: 数据库连接失败
A: 确保应用有写入权限，数据库文件会在首次运行时自动创建

### Q: Wails dev 启动失败
A: 确保 Wails CLI 已正确安装，运行 `wails doctor` 检查环境

## 技术支持

如遇到问题，请检查：
1. Go 版本是否符合要求
2. Node.js 版本是否符合要求
3. Wails CLI 是否正确安装
4. 防火墙是否阻止了应用运行

## 许可证

MIT License
