<template>
  <div class="policy-form-view">
    <PageHeader
      :title="isEdit ? '编辑政策模板' : '新建政策模板'"
      :sub-title="isEdit ? `模板ID: ${route.params.id}` : '创建新的政策模板'"
    >
      <template #extra>
        <el-button @click="handleBack">返回</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ isEdit ? '保存修改' : '创建模板' }}
        </el-button>
      </template>
    </PageHeader>

    <el-card class="form-card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        :disabled="loading"
      >
        <!-- 基本信息 -->
        <div class="form-section">
          <div class="section-title">基本信息</div>
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="模板名称" prop="name">
                <el-input v-model="form.name" placeholder="请输入模板名称" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="所属通道" prop="channel_id">
                <ChannelSelect v-model="form.channel_id" style="width: 100%" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="状态" prop="status">
                <el-radio-group v-model="form.status">
                  <el-radio :label="1">启用</el-radio>
                  <el-radio :label="0">禁用</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="设为默认">
                <el-switch v-model="form.is_default" />
                <span class="form-tip">设为该通道的默认政策模板</span>
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <!-- 费率配置 -->
        <div class="form-section">
          <div class="section-title">费率配置</div>
          <el-row :gutter="24">
            <el-col :span="8">
              <el-form-item label="贷记卡费率" prop="credit_rate">
                <el-input-number
                  v-model="form.credit_rate"
                  :min="0"
                  :max="100"
                  :precision="4"
                  :step="0.01"
                  style="width: 160px"
                />
                <span class="form-unit">%</span>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="借记卡费率" prop="debit_rate">
                <el-input-number
                  v-model="form.debit_rate"
                  :min="0"
                  :max="100"
                  :precision="4"
                  :step="0.01"
                  style="width: 160px"
                />
                <span class="form-unit">%</span>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="借记卡封顶" prop="debit_cap">
                <el-input-number
                  v-model="form.debit_cap"
                  :min="0"
                  :max="1000"
                  :precision="2"
                  :step="1"
                  style="width: 160px"
                />
                <span class="form-unit">元</span>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="24">
            <el-col :span="8">
              <el-form-item label="云闪付费率" prop="qrcode_rate">
                <el-input-number
                  v-model="form.qrcode_rate"
                  :min="0"
                  :max="100"
                  :precision="4"
                  :step="0.01"
                  style="width: 160px"
                />
                <span class="form-unit">%</span>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="扫码费率" prop="scan_rate">
                <el-input-number
                  v-model="form.scan_rate"
                  :min="0"
                  :max="100"
                  :precision="4"
                  :step="0.01"
                  style="width: 160px"
                />
                <span class="form-unit">%</span>
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <!-- 分润配置 -->
        <div class="form-section">
          <div class="section-title">
            分润配置
            <el-button type="primary" link :icon="Plus" @click="addProfitRule">
              添加规则
            </el-button>
          </div>

          <el-table :data="form.profit_rules" border>
            <el-table-column prop="level" label="代理层级" width="120" align="center">
              <template #default="{ row, $index }">
                <el-select v-model="row.level" placeholder="选择层级" size="small">
                  <el-option label="一级代理" :value="1" />
                  <el-option label="二级代理" :value="2" />
                  <el-option label="三级代理" :value="3" />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column prop="profit_type" label="分润类型" width="140" align="center">
              <template #default="{ row }">
                <el-select v-model="row.profit_type" placeholder="选择类型" size="small">
                  <el-option label="固定金额" value="fixed" />
                  <el-option label="费率差" value="rate_diff" />
                  <el-option label="交易比例" value="percentage" />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column prop="profit_value" label="分润值" width="160" align="center">
              <template #default="{ row }">
                <el-input-number
                  v-model="row.profit_value"
                  :min="0"
                  :precision="4"
                  :step="0.01"
                  size="small"
                  style="width: 120px"
                />
              </template>
            </el-table-column>
            <el-table-column prop="min_amount" label="最小交易额" width="140" align="center">
              <template #default="{ row }">
                <el-input-number
                  v-model="row.min_amount"
                  :min="0"
                  :precision="2"
                  size="small"
                  style="width: 120px"
                />
              </template>
            </el-table-column>
            <el-table-column prop="max_amount" label="最大交易额" width="140" align="center">
              <template #default="{ row }">
                <el-input-number
                  v-model="row.max_amount"
                  :min="0"
                  :precision="2"
                  size="small"
                  style="width: 120px"
                />
              </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" min-width="150">
              <template #default="{ row }">
                <el-input v-model="row.remark" placeholder="备注" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="80" align="center" fixed="right">
              <template #default="{ $index }">
                <el-button type="danger" link @click="removeProfitRule($index)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 备注信息 -->
        <div class="form-section">
          <div class="section-title">备注信息</div>
          <el-form-item label="模板说明" prop="description">
            <el-input
              v-model="form.description"
              type="textarea"
              :rows="3"
              placeholder="请输入模板说明"
            />
          </el-form-item>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import ChannelSelect from '@/components/Common/ChannelSelect.vue'
import { getPolicyTemplate, createPolicyTemplate, updatePolicyTemplate } from '@/api/policy'

const route = useRoute()
const router = useRouter()

// 是否编辑模式
const isEdit = computed(() => route.name === 'PolicyEdit')

// 加载状态
const loading = ref(false)
const submitting = ref(false)

// 表单引用
const formRef = ref<FormInstance>()

// 分润规则类型
interface ProfitRule {
  level: number
  profit_type: 'fixed' | 'rate_diff' | 'percentage'
  profit_value: number
  min_amount: number
  max_amount: number
  remark: string
}

// 表单数据
const form = reactive({
  name: '',
  channel_id: undefined as number | undefined,
  status: 1,
  is_default: false,
  credit_rate: 0.6,
  debit_rate: 0.6,
  debit_cap: 20,
  qrcode_rate: 0.38,
  scan_rate: 0.38,
  profit_rules: [] as ProfitRule[],
  description: '',
})

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: '请输入模板名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度为2-50个字符', trigger: 'blur' },
  ],
  channel_id: [
    { required: true, message: '请选择所属通道', trigger: 'change' },
  ],
  credit_rate: [
    { required: true, message: '请输入贷记卡费率', trigger: 'blur' },
  ],
  debit_rate: [
    { required: true, message: '请输入借记卡费率', trigger: 'blur' },
  ],
  debit_cap: [
    { required: true, message: '请输入借记卡封顶', trigger: 'blur' },
  ],
}

// 添加分润规则
function addProfitRule() {
  form.profit_rules.push({
    level: 1,
    profit_type: 'rate_diff',
    profit_value: 0,
    min_amount: 0,
    max_amount: 0,
    remark: '',
  })
}

// 删除分润规则
function removeProfitRule(index: number) {
  form.profit_rules.splice(index, 1)
}

// 获取模板详情
async function fetchDetail() {
  if (!isEdit.value) return

  loading.value = true
  try {
    const data = await getPolicyTemplate(Number(route.params.id))
    Object.assign(form, {
      name: data.name,
      channel_id: data.channel_id,
      status: data.status,
      is_default: data.is_default,
      credit_rate: data.credit_rate * 100,
      debit_rate: data.debit_rate * 100,
      debit_cap: data.debit_cap,
      qrcode_rate: (data.qrcode_rate || 0) * 100,
      scan_rate: (data.scan_rate || 0) * 100,
      profit_rules: data.profit_rules || [],
      description: data.description || '',
    })
  } catch (error) {
    console.error('Fetch policy template error:', error)
    ElMessage.error('获取模板详情失败')
  } finally {
    loading.value = false
  }
}

// 返回列表
function handleBack() {
  router.push('/policies/list')
}

// 提交表单
async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const submitData = {
      ...form,
      credit_rate: form.credit_rate / 100,
      debit_rate: form.debit_rate / 100,
      qrcode_rate: form.qrcode_rate / 100,
      scan_rate: form.scan_rate / 100,
    }

    if (isEdit.value) {
      await updatePolicyTemplate(Number(route.params.id), submitData)
      ElMessage.success('更新成功')
    } else {
      await createPolicyTemplate(submitData)
      ElMessage.success('创建成功')
    }
    router.push('/policies/list')
  } catch (error) {
    console.error('Submit policy template error:', error)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchDetail()
})
</script>

<style lang="scss" scoped>
.policy-form-view {
  padding: 0;
}

.form-card {
  margin-top: $spacing-md;
}

.form-section {
  margin-bottom: $spacing-xl;

  &:last-child {
    margin-bottom: 0;
  }

  .section-title {
    font-size: 16px;
    font-weight: 600;
    color: $text-primary;
    margin-bottom: $spacing-md;
    padding-bottom: $spacing-sm;
    border-bottom: 1px solid $border-color;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
}

.form-unit {
  margin-left: $spacing-sm;
  color: $text-secondary;
}

.form-tip {
  margin-left: $spacing-md;
  font-size: 12px;
  color: $text-placeholder;
}
</style>
