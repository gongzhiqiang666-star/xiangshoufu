<template>
  <div class="terminal-list-view">
    <PageHeader title="终端管理" sub-title="终端列表">
      <template #extra>
        <el-button :icon="Upload" @click="handleImport">入库</el-button>
        <el-button :icon="Bottom" @click="handleBatchDispatch">批量下发</el-button>
        <el-button :icon="Top" @click="handleBatchRecall">批量回拨</el-button>
      </template>
    </PageHeader>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card">
          <div class="stat-value">{{ stats.total }}</div>
          <div class="stat-label">终端总数</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card activated">
          <div class="stat-value">{{ stats.activated }}</div>
          <div class="stat-label">已激活</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card">
          <div class="stat-value">{{ stats.distributed }}</div>
          <div class="stat-label">未激活</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card stock">
          <div class="stat-value">{{ stats.stock }}</div>
          <div class="stat-label">库存</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card">
          <div class="stat-value">{{ stats.yesterday_activated }}</div>
          <div class="stat-label">昨日激活</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card today">
          <div class="stat-value">{{ stats.today_activated }}</div>
          <div class="stat-label">今日激活</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card">
          <div class="stat-value">{{ stats.month_activated }}</div>
          <div class="stat-label">本月激活</div>
        </div>
      </el-col>
    </el-row>

    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="通道">
        <ChannelSelect v-model="searchForm.channel_id" />
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
          <el-option label="库存" value="stock" />
          <el-option label="已下发" value="distributed" />
          <el-option label="已激活" value="activated" />
          <el-option label="已回拨" value="returned" />
        </el-select>
      </el-form-item>
      <el-form-item label="持有代理商">
        <AgentSelect v-model="searchForm.owner_agent_id" />
      </el-form-item>
      <el-form-item label="SN号">
        <el-input v-model="searchForm.sn" placeholder="请输入SN号" clearable />
      </el-form-item>
    </SearchForm>

    <!-- 表格 -->
    <ProTable
      ref="tableRef"
      :data="tableData"
      :loading="loading"
      :total="total"
      v-model:page="page"
      v-model:page-size="pageSize"
      :selection="true"
      @refresh="fetchData"
      @selection-change="handleSelectionChange"
    >
      <el-table-column prop="sn" label="SN号" width="150" />
      <el-table-column prop="channel_name" label="通道" width="100" />
      <el-table-column prop="owner_agent_name" label="持有代理商" width="120" />
      <el-table-column prop="merchant_name" label="绑定商户" width="100">
        <template #default="{ row }">
          {{ row.merchant_name || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">
            {{ getStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="activated_at" label="激活时间" width="170">
        <template #default="{ row }">
          {{ row.activated_at || '-' }}
        </template>
      </el-table-column>

      <template #action="{ row }">
        <el-button type="primary" link @click="handleView(row)">详情</el-button>
        <el-button
          v-if="row.status === 'stock'"
          type="success"
          link
          @click="handleDispatch([row])"
        >
          下发
        </el-button>
        <el-button
          v-if="row.status !== 'stock'"
          type="warning"
          link
          @click="handleRecall([row])"
        >
          回拨
        </el-button>
      </template>
    </ProTable>

    <!-- 下发弹窗 -->
    <el-dialog v-model="dispatchDialogVisible" title="终端下发" width="600px">
      <div class="dispatch-info">
        已选终端: {{ selectedTerminals.length }}台
      </div>
      <el-form :model="dispatchForm" label-width="120px">
        <el-form-item label="下发给代理商" required>
          <AgentSelect v-model="dispatchForm.to_agent_id" style="width: 100%" />
        </el-form-item>
        <el-form-item label="设置货款代扣">
          <el-switch v-model="dispatchForm.enable_deduction" />
        </el-form-item>
        <template v-if="dispatchForm.enable_deduction">
          <el-form-item label="单价">
            <el-input-number v-model="dispatchForm.unit_price" :min="0" />
            <span class="form-tip">元/台</span>
          </el-form-item>
          <el-form-item label="总金额">
            <span class="total-amount">
              ¥{{ (dispatchForm.unit_price * selectedTerminals.length).toFixed(2) }}
            </span>
          </el-form-item>
          <el-form-item label="扣款来源">
            <el-checkbox-group v-model="dispatchForm.wallet_sources">
              <el-checkbox label="profit">分润钱包</el-checkbox>
              <el-checkbox label="service">服务费钱包</el-checkbox>
              <el-checkbox label="reward">奖励钱包</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
        </template>
      </el-form>
      <div class="dispatch-warning">
        <el-icon><WarningFilled /></el-icon>
        注意: PC端支持跨级下发，APP端仅支持下发给直属下级
      </div>
      <template #footer>
        <el-button @click="dispatchDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitDispatch">确认下发</el-button>
      </template>
    </el-dialog>

    <!-- 入库弹窗 -->
    <el-dialog v-model="importDialogVisible" title="终端入库" width="650px">
      <el-form :model="importForm" label-width="100px">
        <el-form-item label="选择通道" required>
          <ChannelSelect v-model="importForm.channel_id" style="width: 100%" @change="handleChannelChange" />
        </el-form-item>
        <el-form-item label="终端类型">
          <el-select
            v-model="importForm.terminal_type_id"
            placeholder="请先选择通道"
            :disabled="!importForm.channel_id"
            style="width: 100%"
            clearable
            @change="handleTerminalTypeChange"
          >
            <el-option
              v-for="type in terminalTypes"
              :key="type.id"
              :label="type.full_name"
              :value="type.id"
            />
          </el-select>
          <div class="form-tip">选择终端类型后自动填充品牌和型号</div>
        </el-form-item>
        <el-form-item label="品牌编码">
          <el-input v-model="importForm.brand_code" placeholder="品牌编码（可手动填写）" />
        </el-form-item>
        <el-form-item label="型号编码">
          <el-input v-model="importForm.model_code" placeholder="型号编码（可手动填写）" />
        </el-form-item>
        <el-form-item label="入库方式">
          <el-radio-group v-model="importForm.import_mode">
            <el-radio label="range">号段区间</el-radio>
            <el-radio label="batch">批量输入</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 号段区间模式 -->
        <template v-if="importForm.import_mode === 'range'">
          <el-form-item label="起始SN" required>
            <el-input v-model="importForm.sn_start" placeholder="如：NL001 或 12345678" @input="calcRangePreview" />
          </el-form-item>
          <el-form-item label="结束SN" required>
            <el-input v-model="importForm.sn_end" placeholder="如：NL005 或 12345682" @input="calcRangePreview" />
          </el-form-item>
          <el-form-item label="预计数量">
            <span class="preview-count" :class="{ error: rangePreview.error }">
              {{ rangePreview.error || `${rangePreview.count} 台` }}
            </span>
          </el-form-item>
          <el-form-item v-if="rangePreview.count > 0 && rangePreview.count <= 10" label="预览">
            <div class="sn-preview">
              <el-tag v-for="sn in rangePreview.samples" :key="sn" size="small">{{ sn }}</el-tag>
            </div>
          </el-form-item>
        </template>

        <!-- 批量输入模式 -->
        <template v-else>
          <el-form-item label="终端SN号" required>
            <el-input
              v-model="importForm.sn_text"
              type="textarea"
              :rows="8"
              placeholder="请输入终端SN号，每行一个，或使用逗号/分号分隔&#10;最多支持1000个"
            />
          </el-form-item>
        </template>
      </el-form>

      <template v-if="importResult && importResult.failed_count > 0">
        <el-divider />
        <div class="import-result">
          <el-alert type="warning" :closable="false" show-icon>
            <template #title>
              入库完成：成功 {{ importResult.success_count }} 台，失败 {{ importResult.failed_count }} 台
            </template>
          </el-alert>
          <div v-if="importResult.errors.length > 0" class="error-list">
            <div v-for="(error, index) in importResult.errors.slice(0, 10)" :key="index" class="error-item">
              {{ error }}
            </div>
            <div v-if="importResult.errors.length > 10" class="error-more">
              ... 还有 {{ importResult.errors.length - 10 }} 条错误
            </div>
          </div>
        </div>
      </template>

      <template #footer>
        <el-button @click="importDialogVisible = false">关闭</el-button>
        <el-button type="primary" :loading="importLoading" @click="handleSubmitImport">确认入库</el-button>
      </template>
    </el-dialog>

    <!-- 回拨弹窗 -->
    <el-dialog v-model="recallDialogVisible" title="终端回拨" width="500px">
      <div class="recall-info">
        已选终端: {{ selectedTerminals.length }}台
      </div>
      <el-form :model="recallForm" label-width="120px">
        <el-form-item label="回拨给代理商" required>
          <AgentSelect v-model="recallForm.to_agent_id" style="width: 100%" placeholder="请选择上级代理商" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="recallForm.remark" type="textarea" :rows="2" placeholder="可选，填写回拨备注" />
        </el-form-item>
      </el-form>

      <div class="recall-warning">
        <el-icon><WarningFilled /></el-icon>
        注意: 只能回拨未激活的终端，PC端支持跨级回拨，APP端仅支持回拨给直属上级
      </div>

      <template v-if="recallResult && recallResult.failed_count > 0">
        <el-divider />
        <div class="recall-result">
          <el-alert type="warning" :closable="false" show-icon>
            <template #title>
              回拨完成：成功 {{ recallResult.success_count }} 台，失败 {{ recallResult.failed_count }} 台
            </template>
          </el-alert>
          <div v-if="recallResult.errors.length > 0" class="error-list">
            <div v-for="(error, index) in recallResult.errors.slice(0, 10)" :key="index" class="error-item">
              {{ error }}
            </div>
          </div>
        </div>
      </template>

      <template #footer>
        <el-button @click="recallDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="recallLoading" @click="handleSubmitRecall">确认回拨</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Upload, Bottom, Top, WarningFilled } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import ChannelSelect from '@/components/Common/ChannelSelect.vue'
import AgentSelect from '@/components/Common/AgentSelect.vue'
import { getTerminals, getTerminalStats, dispatchTerminals, importTerminals, recallTerminals } from '@/api/terminal'
import { getTerminalTypesByChannel } from '@/api/terminalType'
import type { Terminal, TerminalStatus, TerminalStats, TerminalImportResult, TerminalRecallResult } from '@/types'
import type { TerminalType } from '@/types/terminalType'

const router = useRouter()

// 统计数据
const stats = ref<TerminalStats>({
  total: 0,
  stock: 0,
  distributed: 0,
  activated: 0,
  returned: 0,
  yesterday_activated: 0,
  today_activated: 0,
  month_activated: 0,
})

// 搜索表单
const searchForm = reactive({
  channel_id: undefined as number | undefined,
  status: undefined as TerminalStatus | undefined,
  owner_agent_id: undefined as number | undefined,
  sn: '',
})

// 表格数据
const tableRef = ref()
const tableData = ref<Terminal[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const selectedTerminals = ref<Terminal[]>([])

// 下发弹窗
const dispatchDialogVisible = ref(false)
const dispatchForm = reactive({
  to_agent_id: undefined as number | undefined,
  enable_deduction: false,
  unit_price: 50,
  wallet_sources: ['profit'] as string[],
})

// 状态配置
function getStatusType(status: TerminalStatus) {
  const map: Record<TerminalStatus, string> = {
    stock: 'info',
    distributed: 'warning',
    activated: 'success',
    returned: '',
  }
  return map[status] || ''
}

function getStatusLabel(status: TerminalStatus) {
  const map: Record<TerminalStatus, string> = {
    stock: '库存',
    distributed: '已下发',
    activated: '已激活',
    returned: '已回拨',
  }
  return map[status] || status
}

// 获取统计数据
async function fetchStats() {
  try {
    stats.value = await getTerminalStats()
  } catch (error) {
    console.error('Fetch terminal stats error:', error)
  }
}

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getTerminals({
      ...searchForm,
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('Fetch terminals error:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
function handleSearch() {
  page.value = 1
  fetchData()
}

// 重置
function handleReset() {
  page.value = 1
  fetchData()
}

// 选择变化
function handleSelectionChange(selection: Terminal[]) {
  selectedTerminals.value = selection
}

// 查看详情
function handleView(row: Terminal) {
  router.push(`/terminals/${row.id}`)
}

// 入库弹窗
const importDialogVisible = ref(false)
const importForm = reactive({
  channel_id: undefined as number | undefined,
  channel_code: '',
  terminal_type_id: undefined as number | undefined,
  brand_code: '',
  model_code: '',
  import_mode: 'range' as 'range' | 'batch',
  sn_start: '',
  sn_end: '',
  sn_text: '', // 多行文本输入
})
const importLoading = ref(false)
const importResult = ref<TerminalImportResult | null>(null)
const terminalTypes = ref<TerminalType[]>([])
const rangePreview = reactive({
  count: 0,
  samples: [] as string[],
  error: '',
})

// 通道变更时加载终端类型
async function handleChannelChange(channelId: number | undefined) {
  importForm.terminal_type_id = undefined
  terminalTypes.value = []
  if (channelId) {
    try {
      terminalTypes.value = await getTerminalTypesByChannel(channelId)
    } catch (error) {
      console.error('获取终端类型失败:', error)
    }
  }
}

// 终端类型变更时自动填充品牌和型号
function handleTerminalTypeChange(typeId: number | undefined) {
  if (typeId) {
    const type = terminalTypes.value.find(t => t.id === typeId)
    if (type) {
      importForm.brand_code = type.brand_code
      importForm.model_code = type.model_code
    }
  }
}

// 计算号段区间预览
function calcRangePreview() {
  rangePreview.count = 0
  rangePreview.samples = []
  rangePreview.error = ''

  const start = importForm.sn_start.trim()
  const end = importForm.sn_end.trim()

  if (!start || !end) {
    return
  }

  // 解析号段
  const result = parseSnRange(start, end)
  if (result.error) {
    rangePreview.error = result.error
    return
  }

  rangePreview.count = result.count
  rangePreview.samples = result.samples
}

// 解析SN号段区间
function parseSnRange(start: string, end: string): { count: number; samples: string[]; snList: string[]; error: string } {
  // 找到共同前缀
  let prefixLen = 0
  while (prefixLen < start.length && prefixLen < end.length && start[prefixLen] === end[prefixLen]) {
    prefixLen++
  }

  // 找到数字部分的起始位置（从后向前找）
  let numStartIdx = start.length
  while (numStartIdx > 0 && /\d/.test(start[numStartIdx - 1])) {
    numStartIdx--
  }

  const prefix = start.substring(0, numStartIdx)
  const startNum = start.substring(numStartIdx)
  const endNum = end.substring(numStartIdx)

  // 验证前缀一致
  if (!end.startsWith(prefix)) {
    return { count: 0, samples: [], snList: [], error: '起始和结束SN前缀不一致' }
  }

  // 验证数字部分
  if (!/^\d+$/.test(startNum) || !/^\d+$/.test(endNum)) {
    return { count: 0, samples: [], snList: [], error: 'SN格式错误，尾部应为数字' }
  }

  const startVal = parseInt(startNum, 10)
  const endVal = parseInt(endNum, 10)

  if (startVal > endVal) {
    return { count: 0, samples: [], snList: [], error: '起始SN应小于或等于结束SN' }
  }

  const count = endVal - startVal + 1
  if (count > 1000) {
    return { count: 0, samples: [], snList: [], error: '单次入库不能超过1000台' }
  }

  // 生成SN列表
  const numLen = startNum.length
  const snList: string[] = []
  for (let i = startVal; i <= endVal; i++) {
    snList.push(prefix + i.toString().padStart(numLen, '0'))
  }

  // 只返回前10个作为预览
  const samples = snList.slice(0, 10)

  return { count, samples, snList, error: '' }
}

// 入库
function handleImport() {
  importForm.channel_id = undefined
  importForm.channel_code = ''
  importForm.terminal_type_id = undefined
  importForm.brand_code = ''
  importForm.model_code = ''
  importForm.import_mode = 'range'
  importForm.sn_start = ''
  importForm.sn_end = ''
  importForm.sn_text = ''
  importResult.value = null
  terminalTypes.value = []
  rangePreview.count = 0
  rangePreview.samples = []
  rangePreview.error = ''
  importDialogVisible.value = true
}

// 提交入库
async function handleSubmitImport() {
  if (!importForm.channel_id) {
    ElMessage.warning('请选择通道')
    return
  }

  let snList: string[] = []

  if (importForm.import_mode === 'range') {
    // 号段区间模式
    const start = importForm.sn_start.trim()
    const end = importForm.sn_end.trim()

    if (!start || !end) {
      ElMessage.warning('请输入起始SN和结束SN')
      return
    }

    const result = parseSnRange(start, end)
    if (result.error) {
      ElMessage.warning(result.error)
      return
    }

    snList = result.snList
  } else {
    // 批量输入模式
    snList = importForm.sn_text
      .split(/[\n,;]+/)
      .map(sn => sn.trim())
      .filter(sn => sn.length > 0)

    if (snList.length === 0) {
      ElMessage.warning('请输入终端SN号')
      return
    }

    if (snList.length > 1000) {
      ElMessage.warning('单次入库不能超过1000台')
      return
    }
  }

  importLoading.value = true
  try {
    const result = await importTerminals({
      channel_id: importForm.channel_id,
      channel_code: importForm.channel_code,
      brand_code: importForm.brand_code,
      model_code: importForm.model_code,
      sn_list: snList,
    })
    importResult.value = result

    if (result.failed_count === 0) {
      ElMessage.success(`入库成功: 共${result.success_count}台`)
      importDialogVisible.value = false
      fetchData()
      fetchStats()
    } else {
      ElMessage.warning(`入库完成: 成功${result.success_count}台，失败${result.failed_count}台`)
    }
  } catch (error) {
    console.error('Import terminals error:', error)
  } finally {
    importLoading.value = false
  }
}

// 回拨弹窗
const recallDialogVisible = ref(false)
const recallForm = reactive({
  to_agent_id: undefined as number | undefined,
  remark: '',
})
const recallLoading = ref(false)
const recallResult = ref<TerminalRecallResult | null>(null)

// 批量下发
function handleBatchDispatch() {
  if (selectedTerminals.value.length === 0) {
    ElMessage.warning('请先选择要下发的终端')
    return
  }
  handleDispatch(selectedTerminals.value)
}

// 下发
function handleDispatch(terminals: Terminal[]) {
  selectedTerminals.value = terminals
  dispatchForm.to_agent_id = undefined
  dispatchForm.enable_deduction = false
  dispatchDialogVisible.value = true
}

// 提交下发
async function handleSubmitDispatch() {
  if (!dispatchForm.to_agent_id) {
    ElMessage.warning('请选择下发的代理商')
    return
  }

  try {
    await dispatchTerminals({
      terminal_ids: selectedTerminals.value.map((t) => t.id),
      to_agent_id: dispatchForm.to_agent_id,
      cargo_deduction: dispatchForm.enable_deduction
        ? {
            unit_price: dispatchForm.unit_price,
            wallet_sources: dispatchForm.wallet_sources,
          }
        : undefined,
    })
    ElMessage.success('下发成功')
    dispatchDialogVisible.value = false
    tableRef.value?.clearSelection()
    fetchData()
    fetchStats()
  } catch (error) {
    console.error('Dispatch terminals error:', error)
  }
}

// 批量回拨
function handleBatchRecall() {
  if (selectedTerminals.value.length === 0) {
    ElMessage.warning('请先选择要回拨的终端')
    return
  }
  handleRecall(selectedTerminals.value)
}

// 回拨
function handleRecall(terminals: Terminal[]) {
  // 检查是否有已激活的终端
  const activatedTerminals = terminals.filter(t => t.status === 'activated')
  if (activatedTerminals.length > 0) {
    ElMessage.warning('已激活的终端不能回拨')
    return
  }

  selectedTerminals.value = terminals
  recallForm.to_agent_id = undefined
  recallForm.remark = ''
  recallResult.value = null
  recallDialogVisible.value = true
}

// 提交回拨
async function handleSubmitRecall() {
  if (!recallForm.to_agent_id) {
    ElMessage.warning('请选择回拨给的代理商')
    return
  }

  recallLoading.value = true
  try {
    const result = await recallTerminals({
      terminal_sns: selectedTerminals.value.map(t => t.sn),
      to_agent_id: recallForm.to_agent_id,
      remark: recallForm.remark,
    })
    recallResult.value = result

    if (result.failed_count === 0) {
      ElMessage.success(`回拨成功: 共${result.success_count}台`)
      recallDialogVisible.value = false
      tableRef.value?.clearSelection()
      fetchData()
      fetchStats()
    } else {
      ElMessage.warning(`回拨完成: 成功${result.success_count}台，失败${result.failed_count}台`)
    }
  } catch (error) {
    console.error('Recall terminals error:', error)
  } finally {
    recallLoading.value = false
  }
}

onMounted(() => {
  fetchStats()
  fetchData()
})
</script>

<style lang="scss" scoped>
.terminal-list-view {
  padding: 0;
}

.stats-row {
  margin-bottom: $spacing-md;

  .el-col {
    margin-bottom: $spacing-md;
  }
}

.stat-card {
  background: $bg-white;
  border-radius: $border-radius-md;
  padding: $spacing-md;
  text-align: center;
  box-shadow: $shadow-sm;

  .stat-value {
    font-size: 24px;
    font-weight: 600;
    color: $text-primary;
  }

  .stat-label {
    font-size: 12px;
    color: $text-secondary;
    margin-top: $spacing-xs;
  }

  &.activated .stat-value {
    color: $success-color;
  }

  &.stock .stat-value {
    color: $info-color;
  }

  &.today .stat-value {
    color: $primary-color;
  }
}

.dispatch-info {
  margin-bottom: $spacing-md;
  padding: $spacing-sm;
  background: $bg-color;
  border-radius: $border-radius-sm;
}

.form-tip {
  margin-left: $spacing-sm;
  color: $text-secondary;
}

.total-amount {
  font-size: 18px;
  font-weight: 600;
  color: $primary-color;
}

.dispatch-warning {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  margin-top: $spacing-md;
  padding: $spacing-sm;
  background: lighten($warning-color, 40%);
  border-radius: $border-radius-sm;
  color: $warning-color;
  font-size: 12px;
}

.recall-info {
  margin-bottom: $spacing-md;
  padding: $spacing-sm;
  background: $bg-color;
  border-radius: $border-radius-sm;
}

.recall-warning {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  margin-top: $spacing-md;
  padding: $spacing-sm;
  background: lighten($warning-color, 40%);
  border-radius: $border-radius-sm;
  color: $warning-color;
  font-size: 12px;
}

.import-result,
.recall-result {
  margin-top: $spacing-md;

  .error-list {
    margin-top: $spacing-sm;
    max-height: 200px;
    overflow-y: auto;

    .error-item {
      padding: $spacing-xs;
      font-size: 12px;
      color: $text-secondary;
      border-bottom: 1px dashed $border-color;

      &:last-child {
        border-bottom: none;
      }
    }

    .error-more {
      padding: $spacing-xs;
      font-size: 12px;
      color: $text-secondary;
      font-style: italic;
    }
  }
}

.preview-count {
  font-size: 16px;
  font-weight: 600;
  color: $primary-color;

  &.error {
    color: $danger-color;
    font-size: 13px;
    font-weight: normal;
  }
}

.sn-preview {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-xs;
}

.form-tip {
  font-size: 12px;
  color: $text-secondary;
  margin-top: 4px;
}
</style>
