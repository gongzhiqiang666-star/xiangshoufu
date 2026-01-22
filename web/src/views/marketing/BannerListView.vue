<template>
  <div class="banner-list">
    <!-- 页面标题和操作栏 -->
    <div class="page-header">
      <h2>滚动图管理</h2>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新增滚动图
      </el-button>
    </div>

    <!-- 搜索筛选 -->
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="queryParams">
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never">
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="图片" width="150">
          <template #default="{ row }">
            <el-image
              :src="row.image_url"
              :preview-src-list="[row.image_url]"
              fit="cover"
              style="width: 120px; height: 60px; border-radius: 4px"
            />
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="150" show-overflow-tooltip />
        <el-table-column label="链接类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getLinkTypeTag(row.link_type)" size="small">
              {{ getLinkTypeLabel(row.link_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="展示时间" width="200">
          <template #default="{ row }">
            <div v-if="row.start_time || row.end_time">
              <div v-if="row.start_time">开始: {{ formatDateTime(row.start_time) }}</div>
              <div v-if="row.end_time">结束: {{ formatDateTime(row.end_time) }}</div>
            </div>
            <span v-else class="text-muted">长期有效</span>
          </template>
        </el-table-column>
        <el-table-column prop="click_count" label="点击次数" width="100" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.page_size"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { getBannerList, deleteBanner, updateBannerStatus } from '@/api/banner'
import type { Banner } from '@/types/banner'
import { LinkType, linkTypeOptions } from '@/types/banner'
import { formatDateTime } from '@/utils/format'

const router = useRouter()
const loading = ref(false)
const tableData = ref<Banner[]>([])
const total = ref(0)

const queryParams = reactive({
  page: 1,
  page_size: 20,
  status: undefined as number | undefined
})

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getBannerList(queryParams)
    tableData.value = res.data || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取Banner列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  queryParams.page = 1
  fetchData()
}

// 重置
const handleReset = () => {
  queryParams.page = 1
  queryParams.status = undefined
  fetchData()
}

// 新增
const handleCreate = () => {
  router.push('/marketing/banners/create')
}

// 编辑
const handleEdit = (row: Banner) => {
  router.push(`/marketing/banners/${row.id}/edit`)
}

// 删除
const handleDelete = async (row: Banner) => {
  try {
    await ElMessageBox.confirm(`确定要删除滚动图"${row.title}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteBanner(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 状态切换
const handleStatusChange = async (row: Banner) => {
  try {
    await updateBannerStatus(row.id, { status: row.status })
    ElMessage.success(row.status === 1 ? '已启用' : '已禁用')
  } catch (error) {
    // 恢复状态
    row.status = row.status === 1 ? 0 : 1
    ElMessage.error('状态更新失败')
  }
}

// 获取链接类型标签
const getLinkTypeLabel = (type: LinkType) => {
  const item = linkTypeOptions.find(opt => opt.value === type)
  return item?.label || '未知'
}

// 获取链接类型标签类型
const getLinkTypeTag = (type: LinkType) => {
  switch (type) {
    case LinkType.None:
      return 'info'
    case LinkType.Internal:
      return 'success'
    case LinkType.External:
      return 'warning'
    default:
      return 'info'
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.banner-list {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  h2 {
    margin: 0;
    font-size: 20px;
    font-weight: 500;
  }
}

.filter-card {
  margin-bottom: 20px;

  :deep(.el-card__body) {
    padding-bottom: 0;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.text-muted {
  color: #909399;
  font-size: 12px;
}
</style>
