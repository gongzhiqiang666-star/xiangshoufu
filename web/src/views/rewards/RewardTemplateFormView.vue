<template>
  <div class="reward-template-form">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-button @click="goBack" :icon="ArrowLeft">返回</el-button>
          <span>{{ isEdit ? '编辑奖励模版' : '新建奖励模版' }}</span>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="120px"
        class="form-container"
      >
        <el-form-item label="模版名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入模版名称" maxlength="100" />
        </el-form-item>

        <el-form-item label="时间类型" prop="time_type">
          <el-radio-group v-model="formData.time_type" :disabled="isEdit">
            <el-radio value="days">按天数</el-radio>
            <el-radio value="months">按自然月</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="维度类型" prop="dimension_type">
          <el-radio-group v-model="formData.dimension_type" :disabled="isEdit">
            <el-radio value="amount">按金额</el-radio>
            <el-radio value="count">按笔数</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="交易类型" prop="trans_types">
          <el-checkbox-group v-model="transTypesArray">
            <el-checkbox value="scan">扫码</el-checkbox>
            <el-checkbox value="debit">借记卡</el-checkbox>
            <el-checkbox value="credit">贷记卡</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="金额限制">
          <el-row :gutter="10">
            <el-col :span="11">
              <el-input-number
                v-model="amountMinYuan"
                :min="0"
                :precision="2"
                placeholder="最小金额"
                style="width: 100%"
              />
            </el-col>
            <el-col :span="2" class="text-center">至</el-col>
            <el-col :span="11">
              <el-input-number
                v-model="amountMaxYuan"
                :min="0"
                :precision="2"
                placeholder="最大金额"
                style="width: 100%"
              />
            </el-col>
          </el-row>
          <div class="form-tip">只有符合金额条件的交易才参与奖励计算</div>
        </el-form-item>

        <el-form-item label="断档开关" prop="allow_gap">
          <el-switch v-model="formData.allow_gap" />
          <span class="switch-label">
            {{ formData.allow_gap ? '允许断档（阶段2不达标，阶段3仍可获得奖励）' : '不允许断档（阶段2不达标，阶段3也无法获得奖励）' }}
          </span>
        </el-form-item>

        <el-divider content-position="left">阶段配置</el-divider>

        <div class="stages-container">
          <div v-for="(stage, index) in formData.stages" :key="index" class="stage-item">
            <el-row :gutter="10" align="middle">
              <el-col :span="2">
                <span class="stage-order">阶段{{ index + 1 }}</span>
              </el-col>
              <el-col :span="4">
                <el-form-item :prop="`stages.${index}.start_value`" :rules="stageRules.start_value" label-width="0">
                  <el-input-number
                    v-model="stage.start_value"
                    :min="1"
                    placeholder="开始"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="1" class="text-center">-</el-col>
              <el-col :span="4">
                <el-form-item :prop="`stages.${index}.end_value`" :rules="stageRules.end_value" label-width="0">
                  <el-input-number
                    v-model="stage.end_value"
                    :min="stage.start_value"
                    placeholder="结束"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="1" class="text-center">
                {{ formData.time_type === 'days' ? '天' : '月' }}
              </el-col>
              <el-col :span="5">
                <el-form-item :prop="`stages.${index}.target_value`" :rules="stageRules.target_value" label-width="0">
                  <el-input-number
                    v-model="stage.target_value"
                    :min="1"
                    placeholder="达标值"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="1" class="text-center">
                {{ formData.dimension_type === 'amount' ? '分' : '笔' }}
              </el-col>
              <el-col :span="4">
                <el-form-item :prop="`stages.${index}.reward_amount`" :rules="stageRules.reward_amount" label-width="0">
                  <el-input-number
                    v-model="stage.reward_amount"
                    :min="1"
                    placeholder="奖励(分)"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="2">
                <el-button
                  type="danger"
                  :icon="Delete"
                  circle
                  size="small"
                  @click="removeStage(index)"
                  :disabled="formData.stages.length <= 1"
                />
              </el-col>
            </el-row>
          </div>

          <el-button type="primary" plain @click="addStage" :icon="Plus">
            添加阶段
          </el-button>
        </div>

        <el-form-item class="form-actions">
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ isEdit ? '保存修改' : '创建模版' }}
          </el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Plus, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import {
  getRewardTemplateDetail,
  createRewardTemplate,
  updateRewardTemplate,
  type CreateRewardTemplateRequest,
  type RewardStage,
} from '@/api/reward'

const route = useRoute()
const router = useRouter()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const isEdit = computed(() => !!route.params.id)
const templateId = computed(() => Number(route.params.id) || 0)

const formData = reactive<CreateRewardTemplateRequest>({
  name: '',
  time_type: 'days',
  dimension_type: 'amount',
  trans_types: '',
  amount_min: undefined,
  amount_max: undefined,
  allow_gap: false,
  stages: [{ stage_order: 1, start_value: 1, end_value: 10, target_value: 10000, reward_amount: 5000 }],
})

const transTypesArray = computed({
  get: () => formData.trans_types ? formData.trans_types.split(',') : [],
  set: (val: string[]) => { formData.trans_types = val.join(',') }
})

const amountMinYuan = computed({
  get: () => formData.amount_min !== undefined ? formData.amount_min / 100 : undefined,
  set: (val: number | undefined) => { formData.amount_min = val !== undefined ? Math.round(val * 100) : undefined }
})

const amountMaxYuan = computed({
  get: () => formData.amount_max !== undefined ? formData.amount_max / 100 : undefined,
  set: (val: number | undefined) => { formData.amount_max = val !== undefined ? Math.round(val * 100) : undefined }
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入模版名称', trigger: 'blur' }],
  time_type: [{ required: true, message: '请选择时间类型', trigger: 'change' }],
  dimension_type: [{ required: true, message: '请选择维度类型', trigger: 'change' }],
  trans_types: [{ required: true, message: '请选择交易类型', trigger: 'change' }],
}

const stageRules = {
  start_value: [{ required: true, message: '必填', trigger: 'blur' }],
  end_value: [{ required: true, message: '必填', trigger: 'blur' }],
  target_value: [{ required: true, message: '必填', trigger: 'blur' }],
  reward_amount: [{ required: true, message: '必填', trigger: 'blur' }],
}

const addStage = () => {
  const lastStage = formData.stages[formData.stages.length - 1]
  formData.stages.push({
    stage_order: formData.stages.length + 1,
    start_value: lastStage.end_value + 1,
    end_value: lastStage.end_value + 10,
    target_value: lastStage.target_value,
    reward_amount: lastStage.reward_amount,
  })
}

const removeStage = (index: number) => {
  formData.stages.splice(index, 1)
  formData.stages.forEach((s, i) => { s.stage_order = i + 1 })
}

const goBack = () => {
  router.push('/rewards/templates')
}

const loadTemplate = async () => {
  if (!templateId.value) return
  try {
    const detail = await getRewardTemplateDetail(templateId.value)
    formData.name = detail.name
    formData.time_type = detail.time_type
    formData.dimension_type = detail.dimension_type
    formData.trans_types = detail.trans_types
    formData.amount_min = detail.amount_min ?? undefined
    formData.amount_max = detail.amount_max ?? undefined
    formData.allow_gap = detail.allow_gap
    formData.stages = detail.stages?.map((s: RewardStage) => ({
      stage_order: s.stage_order,
      start_value: s.start_value,
      end_value: s.end_value,
      target_value: s.target_value,
      reward_amount: s.reward_amount,
    })) || []
  } catch (err) {
    ElMessage.error('加载模版失败')
    goBack()
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()

  if (formData.stages.length === 0) {
    ElMessage.error('请至少配置一个阶段')
    return
  }

  submitting.value = true
  try {
    if (isEdit.value) {
      await updateRewardTemplate(templateId.value, formData)
      ElMessage.success('修改成功')
    } else {
      await createRewardTemplate(formData)
      ElMessage.success('创建成功')
    }
    goBack()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  if (isEdit.value) {
    loadTemplate()
  }
})
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  gap: 16px;
}

.form-container {
  max-width: 900px;
}

.form-tip {
  color: #909399;
  font-size: 12px;
  margin-top: 4px;
}

.switch-label {
  margin-left: 10px;
  color: #606266;
  font-size: 13px;
}

.stages-container {
  padding-left: 120px;
}

.stage-item {
  margin-bottom: 10px;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;
}

.stage-order {
  font-weight: bold;
  color: #409eff;
}

.text-center {
  text-align: center;
  line-height: 32px;
}

.form-actions {
  margin-top: 30px;
}
</style>
