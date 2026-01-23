<template>
  <div class="deduction-create-view">
    <PageHeader title="发起代扣" sub-title="创建伙伴代扣计划">
      <template #extra>
        <el-button @click="handleBack">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          确认发起
        </el-button>
      </template>
    </PageHeader>

    <el-card class="form-card">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
        label-position="right"
      >
        <el-form-item label="被扣款代理商" prop="deductee_id">
          <AgentSelect
            v-model="formData.deductee_id"
            placeholder="请选择被扣款代理商"
            style="width: 400px"
          />
        </el-form-item>

        <el-form-item label="计划类型" prop="plan_type">
          <el-radio-group v-model="formData.plan_type">
            <el-radio :value="2">伙伴代扣</el-radio>
            <el-radio :value="3">押金代扣</el-radio>
          </el-radio-group>
          <div class="form-tip">
            货款代扣请在终端下发时设置，此处仅支持伙伴代扣和押金代扣
          </div>
        </el-form-item>

        <el-form-item label="代扣总金额" prop="total_amount">
          <el-input-number
            v-model="formData.total_amount"
            :min="1"
            :max="10000000"
            :precision="2"
            :step="100"
            placeholder="请输入代扣总金额（元）"
            style="width: 300px"
          />
          <span class="unit">元</span>
        </el-form-item>

        <el-form-item label="分期期数" prop="total_periods">
          <el-input-number
            v-model="formData.total_periods"
            :min="1"
            :max="120"
            :step="1"
            placeholder="请输入期数"
            style="width: 200px"
          />
          <span class="unit">期</span>
          <div class="form-tip">
            每期扣款金额: ¥{{ periodAmount }}
          </div>
        </el-form-item>

        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="formData.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息（可选）"
            maxlength="200"
            show-word-limit
            style="width: 500px"
          />
        </el-form-item>

        <el-divider />

        <!-- 协议确认 -->
        <el-form-item label="协议确认" prop="agreement">
          <div class="agreement-section">
            <el-checkbox v-model="formData.agreement">
              我已阅读并同意
              <el-link type="primary" @click="showAgreement">《代扣服务协议》</el-link>
            </el-checkbox>
          </div>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 预览信息 -->
    <el-card v-if="formData.deductee_id && formData.total_amount > 0" class="preview-card">
      <template #header>代扣计划预览</template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="计划类型">
          <el-tag :type="formData.plan_type === 2 ? 'success' : 'warning'" size="small">
            {{ formData.plan_type === 2 ? '伙伴代扣' : '押金代扣' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="被扣款方">
          待选择代理商
        </el-descriptions-item>
        <el-descriptions-item label="代扣总额">
          <span class="highlight">¥{{ formatAmount(formData.total_amount * 100) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="分期期数">
          {{ formData.total_periods }} 期
        </el-descriptions-item>
        <el-descriptions-item label="每期金额">
          ¥{{ periodAmount }}
        </el-descriptions-item>
        <el-descriptions-item label="扣款频率">
          每日扣款
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 协议弹窗 -->
    <el-dialog
      v-model="agreementVisible"
      title="代扣服务协议"
      width="600px"
      :close-on-click-modal="false"
    >
      <div class="agreement-content" ref="agreementRef" @scroll="handleAgreementScroll">
        <h3>代扣服务协议</h3>
        <p>甲方（扣款方）：{{ currentAgentName }}</p>
        <p>乙方（被扣款方）：被扣款代理商</p>

        <h4>第一条 服务内容</h4>
        <p>甲方根据与乙方的业务往来关系，按照本协议约定的方式和金额，从乙方的分润钱包和/或服务费钱包中代扣相应款项。</p>

        <h4>第二条 代扣金额及期限</h4>
        <p>1. 代扣总金额：以实际发起金额为准；</p>
        <p>2. 代扣期数：以实际发起期数为准；</p>
        <p>3. 每期扣款金额 = 代扣总金额 ÷ 代扣期数；</p>
        <p>4. 扣款频率：每日从乙方账户中自动扣款，直至扣完为止。</p>

        <h4>第三条 扣款规则</h4>
        <p>1. 扣款优先级：优先扣除分润钱包余额，分润余额不足时扣除服务费钱包余额；</p>
        <p>2. 部分扣款：当乙方钱包余额不足时，系统将扣除全部可用余额，剩余部分顺延至下期继续扣除；</p>
        <p>3. 扣款上限：每期扣款不超过应扣金额。</p>

        <h4>第四条 协议生效</h4>
        <p>本协议自乙方点击"同意"按钮接受后生效，双方均应遵守协议约定。</p>

        <h4>第五条 违约责任</h4>
        <p>任何一方违反本协议约定的，应承担相应的违约责任。</p>

        <h4>第六条 争议解决</h4>
        <p>本协议履行过程中发生的争议，双方应友好协商解决；协商不成的，可向有管辖权的人民法院提起诉讼。</p>

        <div class="agreement-footer">
          <p>签署日期：{{ new Date().toLocaleDateString() }}</p>
        </div>
      </div>
      <template #footer>
        <el-button @click="agreementVisible = false">关闭</el-button>
        <el-button type="primary" :disabled="!agreementRead" @click="confirmAgreement">
          {{ agreementRead ? '我已阅读并同意' : '请阅读完协议' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import AgentSelect from '@/components/Common/AgentSelect.vue'
import { createDeductionPlan } from '@/api/deduction'
import { formatAmount } from '@/utils/format'

const router = useRouter()

const formRef = ref<FormInstance>()
const submitting = ref(false)
const agreementVisible = ref(false)
const agreementRead = ref(false)
const agreementRef = ref<HTMLElement>()

// 当前登录代理商名称（实际应从store获取）
const currentAgentName = ref('当前代理商')

// 表单数据
const formData = reactive({
  deductee_id: undefined as number | undefined,
  plan_type: 2, // 默认伙伴代扣
  total_amount: 0, // 元
  total_periods: 12,
  remark: '',
  agreement: false,
})

// 每期金额
const periodAmount = computed(() => {
  if (!formData.total_amount || !formData.total_periods) return '0.00'
  return (formData.total_amount / formData.total_periods).toFixed(2)
})

// 表单验证规则
const formRules: FormRules = {
  deductee_id: [
    { required: true, message: '请选择被扣款代理商', trigger: 'change' },
  ],
  plan_type: [
    { required: true, message: '请选择计划类型', trigger: 'change' },
  ],
  total_amount: [
    { required: true, message: '请输入代扣总金额', trigger: 'blur' },
    { type: 'number', min: 0.01, message: '金额必须大于0', trigger: 'blur' },
  ],
  total_periods: [
    { required: true, message: '请输入分期期数', trigger: 'blur' },
    { type: 'number', min: 1, max: 120, message: '期数必须在1-120之间', trigger: 'blur' },
  ],
  agreement: [
    {
      validator: (rule, value, callback) => {
        if (!value) {
          callback(new Error('请阅读并同意代扣服务协议'))
        } else {
          callback()
        }
      },
      trigger: 'change',
    },
  ],
}

// 显示协议
function showAgreement() {
  agreementVisible.value = true
  agreementRead.value = false
}

// 协议滚动事件
function handleAgreementScroll() {
  if (!agreementRef.value) return
  const el = agreementRef.value
  // 滚动到底部附近时标记为已阅读
  if (el.scrollHeight - el.scrollTop - el.clientHeight < 50) {
    agreementRead.value = true
  }
}

// 确认协议
function confirmAgreement() {
  formData.agreement = true
  agreementVisible.value = false
}

// 返回
function handleBack() {
  router.push('/deductions/list')
}

// 提交
async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    submitting.value = true
    await createDeductionPlan({
      deductee_id: formData.deductee_id!,
      plan_type: formData.plan_type as 1 | 2 | 3,
      total_amount: Math.round(formData.total_amount * 100), // 转换为分
      total_periods: formData.total_periods,
      remark: formData.remark,
    })

    ElMessage.success('代扣计划创建成功')
    router.push('/deductions/list')
  } catch (error) {
    console.error('Create deduction plan error:', error)
  } finally {
    submitting.value = false
  }
}
</script>

<style lang="scss" scoped>
.deduction-create-view {
  padding: 0;
}

.form-card {
  margin-bottom: $spacing-md;

  .unit {
    margin-left: $spacing-sm;
    color: $text-secondary;
  }

  .form-tip {
    font-size: 12px;
    color: $text-secondary;
    margin-top: $spacing-xs;
  }

  .agreement-section {
    .el-link {
      vertical-align: baseline;
    }
  }
}

.preview-card {
  margin-bottom: $spacing-md;

  .highlight {
    color: $primary-color;
    font-weight: 600;
    font-size: 16px;
  }
}

.agreement-content {
  max-height: 400px;
  overflow-y: auto;
  padding: $spacing-md;
  background: $bg-color;
  border-radius: $border-radius-sm;

  h3 {
    text-align: center;
    margin-bottom: $spacing-md;
  }

  h4 {
    margin: $spacing-md 0 $spacing-sm;
    color: $text-primary;
  }

  p {
    margin-bottom: $spacing-sm;
    line-height: 1.8;
    color: $text-secondary;
  }

  .agreement-footer {
    margin-top: $spacing-lg;
    padding-top: $spacing-md;
    border-top: 1px dashed $border-color;
    text-align: right;
  }
}
</style>
