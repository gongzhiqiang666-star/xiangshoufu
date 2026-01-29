<template>
  <div class="pro-table">
    <!-- 表格工具栏 -->
    <div v-if="$slots.toolbar || showToolbar" class="table-toolbar">
      <div class="toolbar-left">
        <slot name="toolbar"></slot>
      </div>
      <div class="toolbar-right">
        <el-button v-if="showRefresh" :icon="Refresh" circle @click="handleRefresh" />
        <el-button v-if="showExport" :icon="Download" @click="handleExport">导出</el-button>
      </div>
    </div>

    <!-- 表格 -->
    <el-table
      ref="tableRef"
      v-loading="loading"
      :data="data"
      :border="border"
      :stripe="stripe"
      :row-key="rowKey"
      :default-expand-all="defaultExpandAll"
      @selection-change="handleSelectionChange"
      @sort-change="handleSortChange"
    >
      <!-- 选择列 -->
      <el-table-column v-if="selection" type="selection" width="55" align="center" />

      <!-- 序号列 -->
      <el-table-column v-if="showIndex" type="index" label="序号" width="60" align="center" />

      <!-- 动态列 -->
      <slot></slot>

      <!-- 操作列 -->
      <el-table-column
        v-if="$slots.action"
        label="操作"
        :width="actionWidth"
        :fixed="actionFixed"
        align="center"
      >
        <template #default="scope">
          <slot name="action" v-bind="scope"></slot>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div v-if="showPagination" class="table-pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="currentPageSize"
        :page-sizes="pageSizes"
        :total="total"
        :layout="paginationLayout"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Refresh, Download } from '@element-plus/icons-vue'

interface Props {
  data: any[]
  loading?: boolean
  total?: number
  page?: number
  pageSize?: number
  pageSizes?: number[]
  border?: boolean
  stripe?: boolean
  selection?: boolean
  showIndex?: boolean
  showToolbar?: boolean
  /** @deprecated 不再使用，刷新功能由页面自行处理 */
  showRefresh?: boolean
  showExport?: boolean
  showPagination?: boolean
  paginationLayout?: string
  rowKey?: string
  defaultExpandAll?: boolean
  actionWidth?: string | number
  actionFixed?: 'left' | 'right' | boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  total: 0,
  page: 1,
  pageSize: 10,
  pageSizes: () => [10, 20, 50, 100],
  border: true,
  stripe: true,
  selection: false,
  showIndex: false,
  showToolbar: false,
  showRefresh: false,
  showExport: false,
  showPagination: true,
  paginationLayout: 'total, sizes, prev, pager, next, jumper',
  rowKey: 'id',
  defaultExpandAll: false,
  actionWidth: 150,
  actionFixed: 'right',
})

const emit = defineEmits<{
  (e: 'update:page', value: number): void
  (e: 'update:pageSize', value: number): void
  (e: 'refresh'): void
  (e: 'export'): void
  (e: 'selection-change', selection: any[]): void
  (e: 'sort-change', sort: { prop: string; order: string }): void
}>()

const tableRef = ref()

const currentPage = computed({
  get: () => props.page,
  set: (val) => emit('update:page', val),
})

const currentPageSize = computed({
  get: () => props.pageSize,
  set: (val) => emit('update:pageSize', val),
})

function handleRefresh() {
  emit('refresh')
}

function handleExport() {
  emit('export')
}

function handleSelectionChange(selection: any[]) {
  emit('selection-change', selection)
}

function handleSortChange(sort: { prop: string; order: string }) {
  emit('sort-change', sort)
}

function handleSizeChange(size: number) {
  emit('update:pageSize', size)
  emit('refresh')
}

function handleCurrentChange(page: number) {
  emit('update:page', page)
  emit('refresh')
}

// 暴露方法
defineExpose({
  tableRef,
  clearSelection: () => tableRef.value?.clearSelection(),
  toggleRowSelection: (row: any, selected: boolean) =>
    tableRef.value?.toggleRowSelection(row, selected),
})
</script>

<style lang="scss" scoped>
.pro-table {
  background: $bg-white;
  border-radius: $border-radius-md;
  padding: $spacing-md;
}

.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $spacing-md;

  .toolbar-left {
    display: flex;
    gap: $spacing-sm;
  }

  .toolbar-right {
    display: flex;
    gap: $spacing-sm;
  }
}

.table-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: $spacing-md;
  padding-top: $spacing-md;
  border-top: 1px solid $border-light;
}
</style>
