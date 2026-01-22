import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../data/models/message_model.dart';
import '../../data/services/message_service.dart';

/// 当前选中的分类
final selectedCategoryProvider = StateProvider<MessageCategory>((ref) {
  return MessageCategory.all;
});

/// 消息列表状态
class MessageListState {
  final List<MessageModel> messages;
  final bool isLoading;
  final bool hasMore;
  final int page;
  final String? error;

  MessageListState({
    this.messages = const [],
    this.isLoading = false,
    this.hasMore = true,
    this.page = 1,
    this.error,
  });

  MessageListState copyWith({
    List<MessageModel>? messages,
    bool? isLoading,
    bool? hasMore,
    int? page,
    String? error,
  }) {
    return MessageListState(
      messages: messages ?? this.messages,
      isLoading: isLoading ?? this.isLoading,
      hasMore: hasMore ?? this.hasMore,
      page: page ?? this.page,
      error: error,
    );
  }
}

/// 消息列表Provider
final messageListProvider =
    StateNotifierProvider<MessageListNotifier, MessageListState>((ref) {
  final messageService = ref.watch(messageServiceProvider);
  final category = ref.watch(selectedCategoryProvider);
  return MessageListNotifier(messageService, category);
});

/// 消息列表Notifier
class MessageListNotifier extends StateNotifier<MessageListState> {
  final MessageService _messageService;
  final MessageCategory _category;

  MessageListNotifier(this._messageService, this._category)
      : super(MessageListState()) {
    loadMessages();
  }

  /// 加载消息列表
  Future<void> loadMessages({bool refresh = true}) async {
    if (state.isLoading) return;

    final page = refresh ? 1 : state.page;
    state = state.copyWith(isLoading: true, error: null);

    try {
      final response = await _messageService.getMessages(
        category: _category == MessageCategory.all ? null : _category.value,
        page: page,
        pageSize: 20,
      );

      final newMessages = response.list;
      state = state.copyWith(
        messages: refresh ? newMessages : [...state.messages, ...newMessages],
        isLoading: false,
        hasMore: response.hasMore,
        page: page + 1,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// 加载更多
  Future<void> loadMore() async {
    if (!state.hasMore || state.isLoading) return;
    await loadMessages(refresh: false);
  }

  /// 刷新
  Future<void> refresh() async {
    await loadMessages(refresh: true);
  }

  /// 标记为已读
  Future<void> markAsRead(int messageId) async {
    try {
      await _messageService.markAsRead(messageId);
      // 更新本地状态
      final updatedMessages = state.messages.map((msg) {
        if (msg.id == messageId) {
          return MessageModel(
            id: msg.id,
            agentId: msg.agentId,
            messageType: msg.messageType,
            typeName: msg.typeName,
            title: msg.title,
            content: msg.content,
            isRead: true,
            isPushed: msg.isPushed,
            relatedId: msg.relatedId,
            relatedType: msg.relatedType,
            expireAt: msg.expireAt,
            createdAt: msg.createdAt,
          );
        }
        return msg;
      }).toList();
      state = state.copyWith(messages: updatedMessages);
    } catch (e) {
      debugPrint('Mark as read error: $e');
    }
  }

  /// 标记所有为已读
  Future<void> markAllAsRead() async {
    try {
      await _messageService.markAllAsRead();
      // 更新本地状态
      final updatedMessages = state.messages.map((msg) {
        return MessageModel(
          id: msg.id,
          agentId: msg.agentId,
          messageType: msg.messageType,
          typeName: msg.typeName,
          title: msg.title,
          content: msg.content,
          isRead: true,
          isPushed: msg.isPushed,
          relatedId: msg.relatedId,
          relatedType: msg.relatedType,
          expireAt: msg.expireAt,
          createdAt: msg.createdAt,
        );
      }).toList();
      state = state.copyWith(messages: updatedMessages);
    } catch (e) {
      debugPrint('Mark all as read error: $e');
    }
  }
}

/// 未读消息数量Provider
final unreadCountProvider = FutureProvider<int>((ref) async {
  final messageService = ref.watch(messageServiceProvider);
  final unreadCount = await messageService.getUnreadCount();
  return unreadCount.count;
});

/// 消息统计Provider
final messageStatsProvider = FutureProvider<MessageStatsModel>((ref) async {
  final messageService = ref.watch(messageServiceProvider);
  return await messageService.getMessageStats();
});

/// 消息详情Provider
final messageDetailProvider =
    FutureProvider.family<MessageModel, int>((ref, id) async {
  final messageService = ref.watch(messageServiceProvider);
  return await messageService.getMessageDetail(id);
});
