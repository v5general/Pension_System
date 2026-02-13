<template>
  <div class="my-data">
    <el-card class="form-card">
      <template #header>
        <div class="card-header">
          <span>{{ editId !== null ? '编辑我的信息' : '录入我的信息' }}</span>
          <el-button v-if="editId !== null" @click="cancelEdit" size="small">取消编辑</el-button>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="120px"
        label-position="right"
      >
        <el-divider content-position="left">基本信息</el-divider>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="formData.name" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="性别" prop="gender">
              <el-select v-model="formData.gender" placeholder="请选择性别" style="width: 100%">
                <el-option label="男" value="男" />
                <el-option label="女" value="女" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="年龄" prop="age">
              <el-input-number v-model="formData.age" :min="0" :max="150" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="身份证号" prop="id_card">
              <el-input v-model="formData.id_card" placeholder="请输入身份证号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话" prop="phone">
              <el-input v-model="formData.phone" placeholder="请输入联系电话" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">健康与护理</el-divider>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="健康状况" prop="health_condition">
              <el-select v-model="formData.health_condition" placeholder="请选择健康状况" style="width: 100%">
                <el-option label="健康" value="健康" />
                <el-option label="良好" value="良好" />
                <el-option label="一般" value="一般" />
                <el-option label="较差" value="较差" />
                <el-option label="需要护理" value="需要护理" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="护理等级" prop="care_level">
              <el-select v-model="formData.care_level" placeholder="请选择护理等级" style="width: 100%">
                <el-option label="无需护理" value="无需护理" />
                <el-option label="轻度护理" value="轻度护理" />
                <el-option label="中度护理" value="中度护理" />
                <el-option label="重度护理" value="重度护理" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">经济与居住</el-divider>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="经济来源" prop="economic_source">
              <el-select v-model="formData.economic_source" placeholder="请选择经济来源" style="width: 100%">
                <el-option label="退休金" value="退休金" />
                <el-option label="子女赡养" value="子女赡养" />
                <el-option label="社会救助" value="社会救助" />
                <el-option label="其他" value="其他" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="居住情况" prop="living_condition">
              <el-input v-model="formData.living_condition" placeholder="如：独居、与子女同住等" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="医疗保险" prop="medical_insurance">
              <el-input v-model="formData.medical_insurance" placeholder="请输入医疗保险类型" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="居住地址" prop="address">
              <el-input v-model="formData.address" placeholder="请输入详细地址" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">其他信息</el-divider>

        <el-form-item label="备注" prop="remarks">
          <el-input
            v-model="formData.remarks"
            type="textarea"
            :rows="4"
            placeholder="请输入其他需要说明的信息"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitLoading" @click="handleSubmit" size="large">
            {{ editId !== null ? '更新保存' : '提交保存' }}
          </el-button>
          <el-button @click="handleReset" size="large">{{ editId !== null ? '取消编辑' : '重置' }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Data List Card -->
    <el-card class="list-card">
      <template #header>
        <div class="card-header">
          <span>我的数据列表</span>
          <div class="header-actions">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索姓名、身份证号、手机号"
              clearable
              style="width: 250px; margin-right: 10px"
              @clear="loadData"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" :icon="Search" @click="loadData">搜索</el-button>
          </div>
        </div>
      </template>

      <el-table :data="tableData" border style="width: 100%" v-loading="loading" :row-class-name="tableRowClassName">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="姓名" width="100" />
        <el-table-column prop="gender" label="性别" width="60" />
        <el-table-column prop="age" label="年龄" width="60" />
        <el-table-column prop="id_card" label="身份证号" width="180" />
        <el-table-column prop="phone" label="联系电话" width="130" />
        <el-table-column prop="health_condition" label="健康状况" width="100" />
        <el-table-column prop="care_level" label="护理等级" width="100" />
        <el-table-column prop="economic_source" label="经济来源" width="100" />
        <el-table-column prop="address" label="地址" min-width="150" show-overflow-tooltip />
        <el-table-column prop="created_at" label="录入时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)" :icon="Edit">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)" :icon="Delete">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        class="pagination"
        @size-change="loadData"
        @current-change="loadData"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
import { Search, Edit, Delete } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import api from '@/api'

const userStore = useUserStore()
const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const loading = ref(false)
const tableData = ref<any[]>([])
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const editId = ref<number | null>(null)

// Get current user ID
const currentUserId = ref<number>(1)

const formData = reactive({
  name: '',
  gender: '',
  age: 0,
  id_card: '',
  phone: '',
  address: '',
  health_condition: '',
  economic_source: '',
  living_condition: '',
  care_level: '',
  medical_insurance: '',
  remarks: ''
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
  age: [{ required: true, message: '请输入年龄', trigger: 'blur' }]
}

const loadData = async () => {
  loading.value = true
  try {
    const response = await api.getMyElderlyList(currentPage.value, pageSize.value, searchKeyword.value, currentUserId.value)
    if (response.success) {
      tableData.value = response.data || []
      total.value = response.total || 0
    }
  } catch (error) {
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const data = {
        ...formData,
        recorder_id: currentUserId.value,
        id_card: formData.id_card || null,
        phone: formData.phone || null,
        address: formData.address || null,
        health_condition: formData.health_condition || null,
        economic_source: formData.economic_source || null,
        living_condition: formData.living_condition || null,
        care_level: formData.care_level || null,
        medical_insurance: formData.medical_insurance || null,
        remarks: formData.remarks || null
      }

      let response
      if (editId.value) {
        response = await api.updateElderly(editId.value, data)
      } else {
        response = await api.createElderly(data)
      }

      if (response.success) {
        ElMessage.success(editId.value ? '更新成功' : '录入成功')
        handleReset()
        loadData()
      } else {
        ElMessage.error(response.message)
      }
    } catch (error) {
      ElMessage.error('操作失败，请检查网络连接')
    } finally {
      submitLoading.value = false
    }
  })
}

const handleReset = () => {
  editId.value = null
  formRef.value?.resetFields()
  Object.assign(formData, {
    name: '',
    gender: '',
    age: 0,
    id_card: '',
    phone: '',
    address: '',
    health_condition: '',
    economic_source: '',
    living_condition: '',
    care_level: '',
    medical_insurance: '',
    remarks: ''
  })
}

const cancelEdit = () => {
  handleReset()
}

const handleEdit = (row: any) => {
  // Check if the record belongs to current user
  if (row.recorder_id !== currentUserId.value) {
    ElMessage.warning('您只能编辑自己录入的数据')
    return
  }

  editId.value = row.id
  Object.assign(formData, {
    name: row.name || '',
    gender: row.gender || '',
    age: row.age || 0,
    id_card: row.id_card || '',
    phone: row.phone || '',
    address: row.address || '',
    health_condition: row.health_condition || '',
    economic_source: row.economic_source || '',
    living_condition: row.living_condition || '',
    care_level: row.care_level || '',
    medical_insurance: row.medical_insurance || '',
    remarks: row.remarks || ''
  })
  // Scroll to form
  window.scrollTo({ top: 0, behavior: 'smooth' })
  ElMessage.info('已载入数据，修改后点击"更新保存"')
}

const handleDelete = (row: any) => {
  // Check if the record belongs to current user
  if (row.recorder_id !== currentUserId.value) {
    ElMessage.warning('您只能删除自己录入的数据')
    return
  }

  ElMessageBox.confirm('确定要删除这条记录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const response = await api.deleteElderly(row.id)
      if (response.success) {
        ElMessage.success('删除成功')
        if (editId.value === row.id) {
          handleReset()
        }
        loadData()
      }
    } catch (error) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

const formatDateTime = (date: string) => {
  if (!date) return ''
  return new Date(date).toLocaleString('zh-CN')
}

const tableRowClassName = ({ rowIndex }: { rowIndex: number }) => {
  return rowIndex % 2 === 0 ? 'even-row' : 'odd-row'
}

onMounted(() => {
  // Get current user ID from store
  const user = userStore.getUser()
  if (user) {
    currentUserId.value = user.id
  }
  loadData()
})
</script>

<style scoped>
.my-data {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-card,
.list-card {
  border-radius: 8px;
}

:deep(.el-card__header) {
  padding: 16px 20px;
}

:deep(.el-card__body) {
  padding: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header span {
  font-size: 16px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  align-items: center;
}

.header-actions .el-input {
  width: 250px;
}

:deep(.el-divider) {
  margin: 20px 0;
}

:deep(.el-divider__text) {
  font-weight: bold;
  padding: 0 12px;
}

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

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
