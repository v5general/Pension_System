<template>
  <div class="data-manage">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>数据管理</span>
          <el-button type="primary" @click="handleAdd">新增</el-button>
        </div>
      </template>

      <el-row :gutter="20" class="search-bar">
        <el-col :span="8">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索姓名、身份证号、手机号"
            clearable
            @clear="loadData"
          >
            <template #append>
              <el-button :icon="Search" @click="loadData" />
            </template>
          </el-input>
        </el-col>
      </el-row>

      <el-table :data="tableData" border style="width: 100%; margin-top: 20px" v-loading="loading">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="姓名" width="100" />
        <el-table-column prop="gender" label="性别" width="60" />
        <el-table-column prop="age" label="年龄" width="60" />
        <el-table-column prop="id_card" label="身份证号" width="180" />
        <el-table-column prop="phone" label="联系电话" width="130" />
        <el-table-column prop="health_condition" label="健康状况" width="120" />
        <el-table-column prop="care_level" label="护理等级" width="100" />
        <el-table-column prop="address" label="地址" min-width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px; justify-content: flex-end"
        @size-change="loadData"
        @current-change="loadData"
      />
    </el-card>

    <!-- Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="800px"
      @close="resetForm"
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="formData.name" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="性别" prop="gender">
              <el-select v-model="formData.gender" placeholder="请选择" style="width: 100%">
                <el-option label="男" value="男" />
                <el-option label="女" value="女" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="年龄" prop="age">
              <el-input-number v-model="formData.age" :min="0" :max="150" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="身份证号" prop="id_card">
              <el-input v-model="formData.id_card" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="联系电话" prop="phone">
              <el-input v-model="formData.phone" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="健康状况" prop="health_condition">
              <el-select v-model="formData.health_condition" placeholder="请选择" style="width: 100%">
                <el-option label="健康" value="健康" />
                <el-option label="良好" value="良好" />
                <el-option label="一般" value="一般" />
                <el-option label="较差" value="较差" />
                <el-option label="需要护理" value="需要护理" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="经济来源" prop="economic_source">
              <el-select v-model="formData.economic_source" placeholder="请选择" style="width: 100%">
                <el-option label="退休金" value="退休金" />
                <el-option label="子女赡养" value="子女赡养" />
                <el-option label="社会救助" value="社会救助" />
                <el-option label="其他" value="其他" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="护理等级" prop="care_level">
              <el-select v-model="formData.care_level" placeholder="请选择" style="width: 100%">
                <el-option label="无需护理" value="无需护理" />
                <el-option label="轻度护理" value="轻度护理" />
                <el-option label="中度护理" value="中度护理" />
                <el-option label="重度护理" value="重度护理" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="居住情况" prop="living_condition">
          <el-input v-model="formData.living_condition" />
        </el-form-item>

        <el-form-item label="医疗保险" prop="medical_insurance">
          <el-input v-model="formData.medical_insurance" />
        </el-form-item>

        <el-form-item label="地址" prop="address">
          <el-input v-model="formData.address" type="textarea" :rows="2" />
        </el-form-item>

        <el-form-item label="备注" prop="remarks">
          <el-input v-model="formData.remarks" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import api from '@/api'

const loading = ref(false)
const submitLoading = ref(false)
const tableData = ref<any[]>([])
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const dialogVisible = ref(false)
const dialogTitle = ref('新增')
const editId = ref<number | null>(null)

const formRef = ref<FormInstance>()
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
    const response = await api.getElderlyList(currentPage.value, pageSize.value, searchKeyword.value)
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

const handleAdd = () => {
  dialogTitle.value = '新增'
  editId.value = null
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑'
  editId.value = row.id
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm('确定要删除这条记录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const response = await api.deleteElderly(row.id)
      if (response.success) {
        ElMessage.success('删除成功')
        loadData()
      }
    } catch (error) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const data = { ...formData, recorder_id: 1 }
      let response

      if (editId.value) {
        response = await api.updateElderly(editId.value, data)
      } else {
        response = await api.createElderly(data)
      }

      if (response.success) {
        ElMessage.success(editId.value ? '更新成功' : '创建成功')
        dialogVisible.value = false
        loadData()
      } else {
        ElMessage.error(response.message)
      }
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

const resetForm = () => {
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

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.data-manage {
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-bar {
  margin-bottom: 20px;
}
</style>
