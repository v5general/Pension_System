<template>
  <div class="data-form">
    <el-card>
      <template #header>
        <span>老年人信息录入</span>
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
            提交保存
          </el-button>
          <el-button @click="handleReset" size="large">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import api from '@/api'

const formRef = ref<FormInstance>()
const submitLoading = ref(false)

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

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      // Convert empty strings to null for nullable fields
      const data = {
        ...formData,
        recorder_id: 1,
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
      const response = await api.createElderly(data)

      if (response.success) {
        ElMessage.success('录入成功')
        handleReset()
      } else {
        ElMessage.error(response.message)
      }
    } catch (error) {
      ElMessage.error('录入失败，请检查网络连接')
    } finally {
      submitLoading.value = false
    }
  })
}

const handleReset = () => {
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
</script>

<style scoped>
.data-form {
  max-width: 1200px;
  margin: 0 auto;
}

:deep(.el-divider) {
  margin: 20px 0;
}

:deep(.el-divider__text) {
  font-weight: bold;
  color: #409eff;
}
</style>
