<template>
  <div class="agent-form-view">
    <PageHeader :title="isEdit ? '编辑代理商' : '新增代理商'" show-back />

    <el-card v-loading="loading" class="form-card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        style="max-width: 800px"
      >
        <el-divider content-position="left">基本信息</el-divider>

        <el-form-item label="代理商名称" prop="agent_name">
          <el-input
            v-model="form.agent_name"
            placeholder="请输入代理商名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="联系人" prop="contact_name">
          <el-input
            v-model="form.contact_name"
            placeholder="请输入联系人姓名"
            maxlength="20"
          />
        </el-form-item>

        <el-form-item label="联系电话" prop="contact_phone">
          <el-input
            v-model="form.contact_phone"
            placeholder="请输入联系电话"
            maxlength="11"
          />
        </el-form-item>

        <el-form-item label="身份证号" prop="id_card_no">
          <el-input
            v-model="form.id_card_no"
            placeholder="请输入身份证号"
            maxlength="18"
          />
        </el-form-item>

        <el-form-item v-if="!isEdit" label="上级代理商" prop="parent_id">
          <el-select
            v-model="form.parent_id"
            filterable
            remote
            :remote-method="searchParentAgents"
            placeholder="选择上级代理商（不选则为自己的下级）"
            style="width: 100%"
            clearable
            :loading="searchLoading"
          >
            <el-option
              v-for="agent in parentAgentOptions"
              :key="agent.id"
              :label="`${agent.agent_name} (${agent.agent_no})`"
              :value="agent.id"
            />
          </el-select>
          <div class="form-tip">不选择上级代理商时，将默认创建为您的直属下级</div>
        </el-form-item>

        <el-divider content-position="left">结算信息</el-divider>

        <el-form-item label="开户银行" prop="bank_name">
          <el-input
            v-model="form.bank_name"
            placeholder="请输入开户银行名称"
            maxlength="50"
          />
        </el-form-item>

        <el-form-item label="开户名" prop="bank_account">
          <el-input
            v-model="form.bank_account"
            placeholder="请输入开户名"
            maxlength="30"
          />
        </el-form-item>

        <el-form-item label="银行卡号" prop="bank_card_no">
          <el-input
            v-model="form.bank_card_no"
            placeholder="请输入银行卡号"
            maxlength="25"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">
            {{ isEdit ? '保存修改' : '立即创建' }}
          </el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import { createAgent, getAgentDetail, updateAgentProfile, searchAgents } from '@/api/agent'
import type { Agent, AgentDetail } from '@/types'

const route = useRoute()
const router = useRouter()

const formRef = ref<FormInstance>()
const loading = ref(false)
const submitting = ref(false)
const searchLoading = ref(false)
const parentAgentOptions = ref<Agent[]>([])

// 判断是编辑还是新增
const agentId = computed(() => {
  const id = route.params.id as string
  return id && id !== 'new' ? Number(id) : null
})
const isEdit = computed(() => !!agentId.value)

// 表单数据
const form = reactive({
  agent_name: '',
  contact_name: '',
  contact_phone: '',
  id_card_no: '',
  bank_name: '',
  bank_account: '',
  bank_card_no: '',
  parent_id: undefined as number | undefined,
})

// 表单验证规则
const rules: FormRules = {
  agent_name: [
    { required: true, message: '请输入代理商名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度在2-50个字符之间', trigger: 'blur' },
  ],
  contact_name: [
    { required: true, message: '请输入联系人姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '姓名长度在2-20个字符之间', trigger: 'blur' },
  ],
  contact_phone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' },
  ],
  id_card_no: [
    { pattern: /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/, message: '请输入正确的身份证号', trigger: 'blur' },
  ],
  bank_card_no: [
    { pattern: /^\d{16,19}$/, message: '请输入正确的银行卡号', trigger: 'blur' },
  ],
}

// 加载编辑数据
async function loadAgentData() {
  if (!agentId.value) return

  loading.value = true
  try {
    const data = await getAgentDetail(agentId.value) as AgentDetail
    form.agent_name = data.agent_name
    form.contact_name = data.contact_name
    form.contact_phone = data.contact_phone
    form.id_card_no = data.id_card_no || ''
    form.bank_name = data.bank_name || ''
    form.bank_account = data.bank_account || ''
    form.bank_card_no = data.bank_card_no || ''
  } catch (error) {
    console.error('Load agent data error:', error)
    ElMessage.error('加载代理商数据失败')
  } finally {
    loading.value = false
  }
}

// 搜索上级代理商
async function searchParentAgents(keyword: string) {
  if (!keyword || keyword.length < 2) {
    parentAgentOptions.value = []
    return
  }

  searchLoading.value = true
  try {
    parentAgentOptions.value = await searchAgents(keyword)
  } catch (error) {
    console.error('Search agents error:', error)
  } finally {
    searchLoading.value = false
  }
}

// 提交表单
async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value) {
        // 编辑模式
        await updateAgentProfile({
          agent_name: form.agent_name,
          contact_name: form.contact_name,
          contact_phone: form.contact_phone,
          bank_name: form.bank_name,
          bank_account: form.bank_account,
          bank_card_no: form.bank_card_no,
        })
        ElMessage.success('修改成功')
        router.push(`/agents/${agentId.value}`)
      } else {
        // 新增模式
        const result = await createAgent({
          agent_name: form.agent_name,
          contact_name: form.contact_name,
          contact_phone: form.contact_phone,
          id_card_no: form.id_card_no,
          bank_name: form.bank_name,
          bank_account: form.bank_account,
          bank_card_no: form.bank_card_no,
          parent_id: form.parent_id,
        })
        ElMessage.success('创建成功')
        router.push(`/agents/${result.id}`)
      }
    } catch (error: unknown) {
      console.error('Submit form error:', error)
      const message = (error as { message?: string })?.message || (isEdit.value ? '修改失败' : '创建失败')
      ElMessage.error(message)
    } finally {
      submitting.value = false
    }
  })
}

// 取消
function handleCancel() {
  router.back()
}

onMounted(() => {
  if (isEdit.value) {
    loadAgentData()
  }
})
</script>

<style lang="scss" scoped>
.agent-form-view {
  padding: 0;
}

.form-card {
  :deep(.el-divider__text) {
    font-weight: 600;
    color: $text-primary;
  }
}

.form-tip {
  font-size: 12px;
  color: $text-secondary;
  margin-top: 4px;
}
</style>
