import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../data/models/message_model.dart';
import 'providers/message_provider.dart';
import 'widgets/message_tab_bar.dart';
import 'widgets/message_list_item.dart';
import 'widgets/message_group_header.dart';

/// 消息通知页面
class MessagePage extends ConsumerStatefulWidget {
  const MessagePage({super.key});

  @override
  ConsumerState<MessagePage> createState() => _MessagePageState();
}

class _MessagePageState extends ConsumerState<MessagePage> {
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
  }

  @override
  void dispose() {
    _scrollController.removeListener(_onScroll);
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >=
        _scrollController.position.maxScrollExtent - 200) {
      ref.read(messageListProvider.notifier).loadMore();
    }
  }

  @override
  Widget build(BuildContext context) {
    final selectedCategory = ref.watch(selectedCategoryProvider);
    final messageListState = ref.watch(messageListProvider);
    final unreadCountAsync = ref.watch(unreadCountProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('消息通知'),
        actions: [
          unreadCountAsync.when(
            data: (count) => count > 0
                ? TextButton(
                    onPressed: _handleMarkAllAsRead,
                    child: const Text('全部已读'),
                  )
                : const SizedBox.shrink(),
            loading: () => const SizedBox.shrink(),
            error: (_, __) => const SizedBox.shrink(),
          ),
        ],
      ),
      body: Column(
        children: [
          // 分类Tab
          MessageTabBar(
            selectedCategory: selectedCategory,
            onCategoryChanged: (category) {
              ref.read(selectedCategoryProvider.notifier).state = category;
            },
          ),
          const Divider(height: 1, color: AppColors.divider),
          // 消息列表
          Expanded(
            child: RefreshIndicator(
              onRefresh: () async {
                await ref.read(messageListProvider.notifier).refresh();
              },
              child: _buildMessageList(messageListState),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMessageList(MessageListState state) {
    if (state.isLoading && state.messages.isEmpty) {
      return const Center(
        child: CircularProgressIndicator(),
      );
    }

    if (state.error != null && state.messages.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(
              Icons.error_outline,
              size: 48,
              color: AppColors.textTertiary,
            ),
            const SizedBox(height: 16),
            Text(
              '加载失败',
              style: const TextStyle(
                fontSize: 14,
                color: AppColors.textSecondary,
              ),
            ),
            const SizedBox(height: 8),
            TextButton(
              onPressed: () {
                ref.read(messageListProvider.notifier).refresh();
              },
              child: const Text('点击重试'),
            ),
          ],
        ),
      );
    }

    if (state.messages.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.mail_outline,
              size: 64,
              color: AppColors.textDisabled,
            ),
            const SizedBox(height: 16),
            Text(
              '暂无消息',
              style: const TextStyle(
                fontSize: 14,
                color: AppColors.textSecondary,
              ),
            ),
          ],
        ),
      );
    }

    // 按日期分组
    final groupedMessages = _groupMessagesByDate(state.messages);

    return ListView.builder(
      controller: _scrollController,
      itemCount: groupedMessages.length + (state.isLoading ? 1 : 0),
      itemBuilder: (context, index) {
        if (index >= groupedMessages.length) {
          return const Padding(
            padding: EdgeInsets.all(16),
            child: Center(
              child: SizedBox(
                width: 20,
                height: 20,
                child: CircularProgressIndicator(strokeWidth: 2),
              ),
            ),
          );
        }

        final item = groupedMessages[index];
        if (item is String) {
          return MessageGroupHeader(title: item);
        } else if (item is MessageModel) {
          return MessageListItem(
            message: item,
            onTap: () => _handleMessageTap(item),
          );
        }
        return const SizedBox.shrink();
      },
    );
  }

  List<dynamic> _groupMessagesByDate(List<MessageModel> messages) {
    final result = <dynamic>[];
    String? lastGroup;

    for (final message in messages) {
      final group = message.dateGroupTitle;
      if (group != lastGroup) {
        result.add(group);
        lastGroup = group;
      }
      result.add(message);
    }

    return result;
  }

  Future<void> _handleMessageTap(MessageModel message) async {
    // 标记为已读
    if (!message.isRead) {
      await ref.read(messageListProvider.notifier).markAsRead(message.id);
      ref.invalidate(unreadCountProvider);
    }

    // 显示消息详情
    if (!mounted) return;
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: AppColors.cardBg,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (context) => _MessageDetailSheet(message: message),
    );
  }

  Future<void> _handleMarkAllAsRead() async {
    await ref.read(messageListProvider.notifier).markAllAsRead();
    ref.invalidate(unreadCountProvider);
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('已全部标记为已读')),
      );
    }
  }
}

/// 消息详情弹窗
class _MessageDetailSheet extends StatelessWidget {
  final MessageModel message;

  const _MessageDetailSheet({required this.message});

  @override
  Widget build(BuildContext context) {
    return Container(
      constraints: BoxConstraints(
        maxHeight: MediaQuery.of(context).size.height * 0.7,
      ),
      padding: const EdgeInsets.all(20),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 拖动条
          Center(
            child: Container(
              width: 40,
              height: 4,
              margin: const EdgeInsets.only(bottom: 20),
              decoration: BoxDecoration(
                color: AppColors.textDisabled,
                borderRadius: BorderRadius.circular(2),
              ),
            ),
          ),
          // 标题
          Row(
            children: [
              _buildIcon(),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      message.title,
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                        color: AppColors.textPrimary,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      message.messageTypeName,
                      style: const TextStyle(
                        fontSize: 12,
                        color: AppColors.textTertiary,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          const Divider(color: AppColors.divider),
          const SizedBox(height: 16),
          // 内容
          Flexible(
            child: SingleChildScrollView(
              child: Text(
                message.content,
                style: const TextStyle(
                  fontSize: 14,
                  color: AppColors.textSecondary,
                  height: 1.6,
                ),
              ),
            ),
          ),
          const SizedBox(height: 16),
          // 时间
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                message.createdAt,
                style: const TextStyle(
                  fontSize: 12,
                  color: AppColors.textTertiary,
                ),
              ),
              if (message.expireAt != null)
                Text(
                  '有效期至: ${message.expireAt}',
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.textTertiary,
                  ),
                ),
            ],
          ),
          const SizedBox(height: 20),
          // 关闭按钮
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: () => Navigator.of(context).pop(),
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.primary,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 12),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
              child: const Text('关闭'),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildIcon() {
    IconData iconData;
    Color iconColor;
    Color bgColor;

    switch (message.type) {
      case MessageType.profit:
        iconData = Icons.trending_up;
        iconColor = AppColors.profitTrade;
        bgColor = AppColors.profitTrade.withOpacity(0.1);
        break;
      case MessageType.activation:
        iconData = Icons.star;
        iconColor = AppColors.profitReward;
        bgColor = AppColors.profitReward.withOpacity(0.1);
        break;
      case MessageType.deposit:
        iconData = Icons.account_balance_wallet;
        iconColor = AppColors.profitDeposit;
        bgColor = AppColors.profitDeposit.withOpacity(0.1);
        break;
      case MessageType.simCashback:
        iconData = Icons.sim_card;
        iconColor = AppColors.profitSim;
        bgColor = AppColors.profitSim.withOpacity(0.1);
        break;
      case MessageType.refund:
        iconData = Icons.undo;
        iconColor = AppColors.danger;
        bgColor = AppColors.danger.withOpacity(0.1);
        break;
      case MessageType.announcement:
        iconData = Icons.campaign;
        iconColor = AppColors.primary;
        bgColor = AppColors.primary.withOpacity(0.1);
        break;
      case MessageType.newAgent:
        iconData = Icons.person_add;
        iconColor = AppColors.success;
        bgColor = AppColors.success.withOpacity(0.1);
        break;
      case MessageType.transaction:
        iconData = Icons.receipt_long;
        iconColor = AppColors.warning;
        bgColor = AppColors.warning.withOpacity(0.1);
        break;
    }

    return Container(
      width: 48,
      height: 48,
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Icon(
        iconData,
        size: 24,
        color: iconColor,
      ),
    );
  }
}
