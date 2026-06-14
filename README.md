# Smart Pension System

[![License](https://img.shields.io/badge/license-GPL--3.0-green.svg)](LICENSE)

**English** | [简体中文](./README.zh-CN.md)

A desktop application built with Wails v2 + Vue 3 + Element Plus for elderly information management, opinion polling, and issue collection.

## Tech Stack

### Backend
- Go 1.21+
- Wails v2
- GORM (SQLite)
- bcrypt password hashing

### Frontend
- Vue 3.4+
- TypeScript
- Element Plus UI
- ECharts
- Pinia state management
- Vue Router

## Project Structure

```
pension-system/
├── main.go                 # Application entry
├── app.go                  # Application context
├── go.mod                  # Go module definition
├── wails.json              # Wails configuration
├── database/               # Database configuration
│   └── db.go
├── models/                 # Data models
│   ├── user.go            # User model
│   ├── elderly.go         # Elderly model
│   ├── survey.go          # Survey model
│   └── issue.go           # Issue model
├── controllers/            # Controllers
│   ├── auth.go            # Authentication
│   ├── data.go            # Data management
│   ├── survey.go          # Opinion surveys
│   └── issue.go           # Issue collection
├── utils/                  # Utility functions
│   └── crypto.go          # Password hashing
└── frontend/              # Frontend code
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    ├── tsconfig.json
    └── src/
        ├── main.ts
        ├── App.vue
        ├── api/           # API wrappers
        │   └── index.ts
        ├── router/        # Routing configuration
        │   └── index.ts
        ├── store/         # State management
        │   └── user.ts
        └── views/         # Page components
            ├── Login.vue
            ├── Home.vue
            ├── DataSummary.vue
            ├── DataManage.vue
            ├── DataForm.vue
            ├── Survey.vue
            └── Issue.vue
```

## Modules

### 1. Authentication
- User login/registration
- Role-based access control (Administrator, Operator, Regular User)
- Default admin account: admin / admin123

### 2. Data Summary
- Total elderly count statistics
- Gender distribution chart
- Age distribution chart
- Health status distribution
- Care level distribution

### 3. Data Management (Administrator)
- Elderly information CRUD
- Data search and pagination
- Data import/export

### 4. Data Entry
- Elderly information form
- Form validation
- Data saving

### 5. Opinion Surveys (Administrator)
- Create surveys
- Publish/close surveys
- View voting results
- Chart visualization

### 6. Issue Collection
- Submit issues
- Issue categorization and priority
- Issue replies
- Issue closure

## Installation and Usage

### Prerequisites
1. Go 1.21 or later
2. Node.js 18 or later
3. Wails CLI

### Install Wails CLI
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Install Dependencies

#### Backend dependencies
```bash
cd D:\Workspace\PensionSystem\pension-system
go mod download
```

#### Frontend dependencies
```bash
cd frontend
npm install
```

### Development Mode
```bash
wails dev
```

### Build Application
```bash
wails build
```

## Database

The application uses SQLite. The database file `pension_system.db` is created automatically on first run.

### Tables
- `users` - User table
- `elderly` - Elderly information table
- `surveys` - Survey table
- `survey_votes` - Survey vote table
- `issues` - Issue collection table

## Default Account

- **Username**: admin
- **Password**: admin123
- **Role**: Administrator

## Development Guide

### Adding a New API Method

1. Add the method to the appropriate controller under `controllers/`
2. Register the controller in the `Bind` call in `main.go`
3. Add the API call method in `frontend/src/api/index.ts`

### Adding a New Page

1. Create a Vue component under `frontend/src/views/`
2. Add a route in `frontend/src/router/index.ts`
3. Add a menu item in `frontend/src/views/Home.vue`

## License

This project is licensed under the [GPL-3.0 License](./LICENSE).
