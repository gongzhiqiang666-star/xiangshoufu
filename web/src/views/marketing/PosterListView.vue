<template>
  <div class="poster-list">
    <!-- 页面标题和操作栏 -->
    <div class="page-header">
      <h2>营销海报管理</h2>
      <div class="header-actions">
        <el-button @click="router.push('/marketing/poster-categories')">
          分类管理
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          新增海报
        </el-button>
      </div>
    </div>

    <!-- 搜索筛选 -->
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="queryParams">
        <el-form-item label="分类">
          <el-select v-model="queryParams.category_id" placeholder="全部分类" clearable style="width: 150px">
            <el-option
              v-for="cat in categories"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="queryParams.keyword" placeholder="搜索标题" clearable style="width: 200px" />
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

    <!-- 海报网格 -->
    <el-card shadow="never">
      <div v-loading="loading" class="poster-grid">
        <div v-for="poster in tableData" :key="poster.id" class="poster-item">
          <div class="poster-image">
            <el-image
              :src="poster.thumbnail_url || poster.image_url"
              :preview-src-list="[poster.image_url]"
              fit="cover"
            />
            <div class="poster-mask">
              <el-button type="primary" size="small" @click="handleEdit(poster)">编辑</el-button>
              <el-button type="danger" size="small" @click="handleDelete(poster)">删除</el-button>
            </div>
          </div>
          <div class="poster-info">
            <div class="poster-title">{{ poster.title }}</div>
            <div class="poster-meta">
              <span class="category">{{ poster.category_name || '未分类' }}</span>
              <el-tag :type="poster.status === 1 ? 'success' : 'info'" size="small">
                {{ poster.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </div>
            <div class="poster-stats">
              <span>下载: {{ poster.download_count }}</span>
              <span>分享: {{ poster.share_count }}</span>
            </div>
          </div>
        </div>
      </div>

      <el-empty v-if="!loading && tableData.length === 0" description="暂无海报" />

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.page_size"
          :total="total"
          :page-sizes="[12, 24, 48, 96]"
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
import { getPosterList, deletePoster, getPosterCategories } from '@/api/poster'
import type { Poster, PosterCategory } from '@/types/poster'

const router = useRouter()
const loading = ref(false)
const tableData = ref<Poster[]>([])
const categories = ref<PosterCategory[]>([])
const total = ref(0)

const queryParams = reactive({
  page: 1,
  page_size: 24,
  category_id: undefined as number | undefined,
  status: undefined as number | undefined,
  keyword: ''
})

// 获取分类
const fetchCategories = async () => {
  try {
    const res = await getPosterCategories({ status: 1 })
    categories.value = res || []
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getPosterList(queryParams)
    tableData.value = res.data || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取海报列表失败:', error)
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
  queryParams.category_id = undefined
  queryParams.status = undefined
  queryParams.keyword = ''
  fetchData()
}

// 新增
const handleCreate = () => {
  router.push('/marketing/posters/create')
}

// 编辑
const handleEdit = (row: Poster) => {
  router.push(`/marketing/posters/${row.id}/edit`)
}

// 删除
const handleDelete = async (row: Poster) => {
  try {
    await ElMessageBox.confirm(`确定要删除海报"${row.title}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deletePoster(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  fetchCategories()
  fetchData()
})
</script>

<style scoped lang="scss">
.poster-list {
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

  .header-actions {
    display: flex;
    gap: 10px;
  }
}

.filter-card {
  margin-bottom: 20px;

  :deep(.el-card__body) {
    padding-bottom: 0;
  }
}

.poster-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  min-height: 200px;
}

.poster-item {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);

    .poster-mask {
      opacity: 1;
    }
  }
}

.poster-image {
  position: relative;
  width: 100%;
  padding-top: 150%; // 2:3 比例

  :deep(.el-image) {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
  }
}

.poster-mask {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  opacity: 0;
  transition: opacity 0.3s;
}

.poster-info {
  padding: 12px;
}

.poster-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.poster-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;

  .category {
    font-size: 12px;
    color: #909399;
  }
}

.poster-stats {
  font-size: 12px;
  color: #909399;
  display: flex;
  gap: 12px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
