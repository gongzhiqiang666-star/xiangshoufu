<template>
  <div class="messages-view">
    <!-- 顶部操作栏 -->
    <el-card class="filter-card">
      <div class="filter-bar">
        <div class="filter-left">
          <el-radio-group v-model="currentTab" @change="handleTabChange">
            <el-radio-button label="all">全部</el-radio-button>
            <el-radio-button label="unread">
              未读
              <el-badge v-if="unreadCount > 0" :value="unreadCount" :max="99" class="unread-badge" />
            </el-radio-button>
            <el-radio-button label="read">已读</el-radio-button>
          </el-radio-group>
        </div>
        <div class="filter-right">
          <el-button type="primary" @click="handleMarkAllRead" :disabled="unreadCount === 0">
            全部标为已读
          </el-button>
          <el-button @click="handleRefresh">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 消息列表 -->
    <el-card class="message-list-card" v-loading="loading">
      <div v-if="messages.length > 0" class="message-list">
        <div
          v-for="message in messages"
          :key="message.id"
          :class="['message-item', { unread: !message.is_read }]"
          @click="handleMessageClick(message)"
        >
          <div class="message-icon">
            <el-icon :size="24" :color="getMessageColor(message.type)">
              <component :is="getMessageIcon(message.type)" />
            </el-icon>
          </div>
          <div class="message-content">
            <div class="message-header">
              <span class="message-title">{{ message.title }}</span>
              <span class="message-time">{{ formatTime(message.created_at) }}</span>
            </div>
            <div class="message-body">{{ message.content }}</div>
            <div class="message-footer">
              <el-tag size="small" :type="getTagType(message.type)">
                {{ message.type_name || message.type }}
              </el-tag>
              <el-tag v-if="!message.is_read" size="small" type="danger">未读</el-tag>
            </div>
          </div>
        </div>
      </div>

      <el-empty v-else description="暂无消息" />

      <!-- 分页 -->
      <div v-if="total > 0" class="pagination-wrapper">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 消息详情弹窗 -->
    <el-dialog v-model="detailDialogVisible" :title="currentMessage?.title" width="600px">
      <div v-if="currentMessage" class="message-detail">
        <div class="detail-meta">
          <el-tag size="small" :type="getTagType(currentMessage.type)">
            {{ currentMessage.type_name || currentMessage.type }}
          </el-tag>
          <span class="detail-time">{{ formatTime(currentMessage.created_at) }}</span>
        </div>
        <div class="detail-content">
          {{ currentMessage.content }}
        </div>
      </div>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Refresh,
  Bell,
  Warning,
  InfoFilled,
  SuccessFilled,
  Message as MessageIcon,
} from '@element-plus/icons-vue'
import { getMessages, getUnreadCount, markAsRead, markAllAsRead } from '@/api/message'
import type { Message } from '@/types'

// 状态
const loading = ref(false)
const messages = ref<Message[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const currentTab = ref('all')
const unreadCount = ref(0)

// 消息详情弹窗
const detailDialogVisible = ref(false)
const currentMessage = ref<Message | null>(null)

// 加载消息列表
async function loadMessages() {
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (currentTab.value === 'unread') {
      params.is_read = false
    } else if (currentTab.value === 'read') {
      params.is_read = true
    }

    const response = await getMessages(params)
    messages.value = response.list || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to load messages:', error)
    ElMessage.error('加载消息失败')
  } finally {
    loading.value = false
  }
}

// 加载未读数
async function loadUnreadCount() {
  try {
    const data = await getUnreadCount()
    unreadCount.value = data.total || 0
  } catch (error) {
    console.error('Failed to load unread count:', error)
  }
}

// 点击消息
async function handleMessageClick(message: Message) {
  currentMessage.value = message
  detailDialogVisible.value = true

  // 标记为已读
  if (!message.is_read) {
    try {
      await markAsRead(message.id)
      message.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    } catch (error) {
      console.error('Failed to mark as read:', error)
    }
  }
}

// 全部标为已读
async function handleMarkAllRead() {
  try {
    await markAllAsRead()
    ElMessage.success('已全部标为已读')
    loadMessages()
    loadUnreadCount()
  } catch (error) {
    console.error('Failed to mark all as read:', error)
    ElMessage.error('操作失败')
  }
}

// Tab切换
function handleTabChange() {
  page.value = 1
  loadMessages()
}

// 分页
function handleSizeChange() {
  page.value = 1
  loadMessages()
}

function handlePageChange() {
  loadMessages()
}

// 刷新
function handleRefresh() {
  loadMessages()
  loadUnreadCount()
}

// 格式化时间
function formatTime(dateStr: string): string {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  return date.toLocaleDateString('zh-CN')
}

// 获取消息图标
function getMessageIcon(type: string) {
  const iconMap: Record<string, any> = {
    system: InfoFilled,
    warning: Warning,
    success: SuccessFilled,
    notice: Bell,
  }
  return iconMap[type] || MessageIcon
}

// 获取消息颜色
function getMessageColor(type: string): string {
  const colorMap: Record<string, string> = {
    system: '#409eff',
    warning: '#e6a23c',
    success: '#67c23a',
    error: '#f56c6c',
    notice: '#909399',
  }
  return colorMap[type] || '#409eff'
}

// 获取标签类型
function getTagType(type: string): '' | 'success' | 'warning' | 'danger' | 'info' {
  const typeMap: Record<string, '' | 'success' | 'warning' | 'danger' | 'info'> = {
    system: 'info',
    warning: 'warning',
    success: 'success',
    error: 'danger',
    notice: '',
  }
  return typeMap[type] || 'info'
}

// 页面加载
onMounted(() => {
  loadMessages()
  loadUnreadCount()
})
</script>

<style lang="scss" scoped>
.messages-view {
  min-height: 100%;
}

.filter-card {
  margin-bottom: $spacing-lg;
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: $spacing-md;
}

.filter-left {
  .unread-badge {
    margin-left: $spacing-xs;
  }
}

.filter-right {
  display: flex;
  gap: $spacing-sm;
}

.message-list-card {
  min-height: 400px;
}

.message-list {
  .message-item {
    display: flex;
    gap: $spacing-md;
    padding: $spacing-lg;
    border-bottom: 1px solid $border-color-lighter;
    cursor: pointer;
    transition: background-color $transition-fast;

    &:hover {
      background-color: $bg-color;
    }

    &:last-child {
      border-bottom: none;
    }

    &.unread {
      background-color: rgba($primary-color, 0.03);

      .message-title {
        font-weight: 600;
      }
    }
  }

  .message-icon {
    flex-shrink: 0;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: $bg-color;
    border-radius: 50%;
  }

  .message-content {
    flex: 1;
    min-width: 0;
  }

  .message-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: $spacing-xs;
  }

  .message-title {
    font-size: 15px;
    color: $text-primary;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .message-time {
    font-size: 12px;
    color: $text-placeholder;
    flex-shrink: 0;
    margin-left: $spacing-md;
  }

  .message-body {
    font-size: 13px;
    color: $text-secondary;
    line-height: 1.5;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    margin-bottom: $spacing-sm;
  }

  .message-footer {
    display: flex;
    gap: $spacing-sm;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: $spacing-lg;
  padding-top: $spacing-lg;
  border-top: 1px solid $border-color-lighter;
}

.message-detail {
  .detail-meta {
    display: flex;
    align-items: center;
    gap: $spacing-md;
    margin-bottom: $spacing-lg;
    padding-bottom: $spacing-md;
    border-bottom: 1px solid $border-color-lighter;
  }

  .detail-time {
    font-size: 13px;
    color: $text-secondary;
  }

  .detail-content {
    font-size: 14px;
    color: $text-primary;
    line-height: 1.8;
    white-space: pre-wrap;
  }
}
</style>
