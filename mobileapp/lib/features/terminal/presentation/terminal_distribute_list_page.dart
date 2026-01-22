import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../domain/models/terminal.dart';
import 'providers/terminal_provider.dart';

/// 划拨记录列表页面
class TerminalDistributeListPage extends ConsumerStatefulWidget {
  const TerminalDistributeListPage({super.key});

  @override
  ConsumerState<TerminalDistributeListPage> createState() =>
      _TerminalDistributeListPageState();
}

class _TerminalDistributeListPageState
    extends ConsumerState<TerminalDistributeListPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);

    // 加载数据
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(sentDistributesProvider.notifier).loadList(refresh: true);
      ref.read(receivedDistributesProvider.notifier).loadList(refresh: true);
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('划拨记录'),
        centerTitle: true,
        bottom: TabBar(
          controller: _tabController,
          labelColor: AppColors.primary,
          unselectedLabelColor: AppColors.textSecondary,
          indicatorColor: AppColors.primary,
          tabs: const [
            Tab(text: '我下发的'),
            Tab(text: '下发给我的'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _SentDistributesTab(),
          _ReceivedDistributesTab(),
        ],
      ),
    );
  }
}

/// 我下发的 Tab
class _SentDistributesTab extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(sentDistributesProvider);

    if (state.isLoading && state.list.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (state.error != null && state.list.isEmpty) {
      return _buildErrorWidget(
        context,
        state.error!,
        () => ref.read(sentDistributesProvider.notifier).loadList(refresh: true),
      );
    }

    if (state.list.isEmpty) {
      return _buildEmptyWidget('暂无下发记录');
    }

    return RefreshIndicator(
      onRefresh: () async {
        await ref.read(sentDistributesProvider.notifier).loadList(refresh: true);
      },
      child: NotificationListener<ScrollNotification>(
        onNotification: (notification) {
          if (notification is ScrollEndNotification &&
              notification.metrics.pixels >=
                  notification.metrics.maxScrollExtent - 100 &&
              state.hasMore &&
              !state.isLoadingMore) {
            ref.read(sentDistributesProvider.notifier).loadList();
          }
          return false;
        },
        child: ListView.builder(
          padding: const EdgeInsets.all(AppSpacing.md),
          itemCount: state.list.length + (state.isLoadingMore ? 1 : 0),
          itemBuilder: (context, index) {
            if (index == state.list.length) {
              return const Center(
                child: Padding(
                  padding: EdgeInsets.all(AppSpacing.md),
                  child: CircularProgressIndicator(),
                ),
              );
            }

            final distribute = state.list[index];
            return _DistributeCard(
              distribute: distribute,
              isSent: true,
            );
          },
        ),
      ),
    );
  }
}

/// 下发给我的 Tab
class _ReceivedDistributesTab extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(receivedDistributesProvider);

    if (state.isLoading && state.list.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (state.error != null && state.list.isEmpty) {
      return _buildErrorWidget(
        context,
        state.error!,
        () => ref.read(receivedDistributesProvider.notifier).loadList(refresh: true),
      );
    }

    if (state.list.isEmpty) {
      return _buildEmptyWidget('暂无接收记录');
    }

    return RefreshIndicator(
      onRefresh: () async {
        await ref.read(receivedDistributesProvider.notifier).loadList(refresh: true);
      },
      child: NotificationListener<ScrollNotification>(
        onNotification: (notification) {
          if (notification is ScrollEndNotification &&
              notification.metrics.pixels >=
                  notification.metrics.maxScrollExtent - 100 &&
              state.hasMore &&
              !state.isLoadingMore) {
            ref.read(receivedDistributesProvider.notifier).loadList();
          }
          return false;
        },
        child: ListView.builder(
          padding: const EdgeInsets.all(AppSpacing.md),
          itemCount: state.list.length + (state.isLoadingMore ? 1 : 0),
          itemBuilder: (context, index) {
            if (index == state.list.length) {
              return const Center(
                child: Padding(
                  padding: EdgeInsets.all(AppSpacing.md),
                  child: CircularProgressIndicator(),
                ),
              );
            }

            final distribute = state.list[index];
            return _DistributeCard(
              distribute: distribute,
              isSent: false,
            );
          },
        ),
      ),
    );
  }
}

/// 划拨记录卡片
class _DistributeCard extends ConsumerWidget {
  final TerminalDistribute distribute;
  final bool isSent;

  const _DistributeCard({
    required this.distribute,
    required this.isSent,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Container(
      margin: const EdgeInsets.only(bottom: AppSpacing.md),
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 4,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 头部：单号和状态
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                '单号: ${distribute.distributeNo}',
                style: const TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: AppColors.textPrimary,
                ),
              ),
              _buildStatusChip(distribute.status),
            ],
          ),
          const SizedBox(height: AppSpacing.sm),
          const Divider(height: 1),
          const SizedBox(height: AppSpacing.sm),

          // 终端SN
          _buildInfoRow('终端SN', distribute.terminalSn),

          // 代理商信息
          if (isSent)
            _buildInfoRow('接收方', '代理商ID: ${distribute.toAgentId}')
          else
            _buildInfoRow('发起方', '代理商ID: ${distribute.fromAgentId}'),

          // 货款
          _buildInfoRow(
            '货款',
            '¥${(distribute.goodsPrice / 100).toStringAsFixed(2)}',
          ),

          // 扣款方式
          _buildInfoRow(
            '扣款方式',
            distribute.deductionType == 1 ? '一次性扣款' : '分期扣款',
          ),

          // 时间
          _buildInfoRow(
            '创建时间',
            distribute.createdAt.toLocal().toString().substring(0, 16),
          ),

          // 备注
          if (distribute.remark != null && distribute.remark!.isNotEmpty)
            _buildInfoRow('备注', distribute.remark!),

          // 操作按钮 - 只在待确认状态显示
          if (distribute.status == 1) ...[
            const SizedBox(height: AppSpacing.sm),
            const Divider(height: 1),
            const SizedBox(height: AppSpacing.sm),
            _buildActionButtons(context, ref),
          ],
        ],
      ),
    );
  }

  /// 构建操作按钮
  Widget _buildActionButtons(BuildContext context, WidgetRef ref) {
    if (isSent) {
      // 我下发的 - 显示取消按钮
      return Row(
        mainAxisAlignment: MainAxisAlignment.end,
        children: [
          OutlinedButton(
            onPressed: () => _showCancelConfirmDialog(context, ref),
            style: OutlinedButton.styleFrom(
              foregroundColor: AppColors.textSecondary,
              side: const BorderSide(color: AppColors.textTertiary),
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            ),
            child: const Text('取消'),
          ),
        ],
      );
    } else {
      // 下发给我的 - 显示确认和拒绝按钮
      return Row(
        mainAxisAlignment: MainAxisAlignment.end,
        children: [
          OutlinedButton(
            onPressed: () => _showRejectConfirmDialog(context, ref),
            style: OutlinedButton.styleFrom(
              foregroundColor: AppColors.danger,
              side: const BorderSide(color: AppColors.danger),
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            ),
            child: const Text('拒绝'),
          ),
          const SizedBox(width: AppSpacing.sm),
          ElevatedButton(
            onPressed: () => _showConfirmDialog(context, ref),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.primary,
              foregroundColor: Colors.white,
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            ),
            child: const Text('确认'),
          ),
        ],
      );
    }
  }

  /// 显示确认对话框
  void _showConfirmDialog(BuildContext context, WidgetRef ref) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认划拨'),
        content: Text('确认接收终端 ${distribute.terminalSn} 吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () async {
              Navigator.of(context).pop();
              await _handleConfirm(context, ref);
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.primary,
            ),
            child: const Text('确认'),
          ),
        ],
      ),
    );
  }

  /// 显示拒绝确认对话框
  void _showRejectConfirmDialog(BuildContext context, WidgetRef ref) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('拒绝划拨'),
        content: Text('确认拒绝接收终端 ${distribute.terminalSn} 吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () async {
              Navigator.of(context).pop();
              await _handleReject(context, ref);
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.danger,
            ),
            child: const Text('拒绝'),
          ),
        ],
      ),
    );
  }

  /// 显示取消确认对话框
  void _showCancelConfirmDialog(BuildContext context, WidgetRef ref) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('取消划拨'),
        content: Text('确认取消划拨终端 ${distribute.terminalSn} 吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('返回'),
          ),
          ElevatedButton(
            onPressed: () async {
              Navigator.of(context).pop();
              await _handleCancel(context, ref);
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.warning,
            ),
            child: const Text('确认取消'),
          ),
        ],
      ),
    );
  }

  /// 处理确认操作
  Future<void> _handleConfirm(BuildContext context, WidgetRef ref) async {
    // 显示loading
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => const Center(child: CircularProgressIndicator()),
    );

    final success = await ref
        .read(receivedDistributesProvider.notifier)
        .confirmDistribute(distribute.id);

    // 关闭loading
    if (context.mounted) {
      Navigator.of(context).pop();
    }

    // 显示结果
    if (context.mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(success ? '确认成功' : '确认失败，请重试'),
          backgroundColor: success ? AppColors.success : AppColors.danger,
        ),
      );
    }
  }

  /// 处理拒绝操作
  Future<void> _handleReject(BuildContext context, WidgetRef ref) async {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => const Center(child: CircularProgressIndicator()),
    );

    final success = await ref
        .read(receivedDistributesProvider.notifier)
        .rejectDistribute(distribute.id);

    if (context.mounted) {
      Navigator.of(context).pop();
    }

    if (context.mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(success ? '拒绝成功' : '拒绝失败，请重试'),
          backgroundColor: success ? AppColors.success : AppColors.danger,
        ),
      );
    }
  }

  /// 处理取消操作
  Future<void> _handleCancel(BuildContext context, WidgetRef ref) async {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => const Center(child: CircularProgressIndicator()),
    );

    final success = await ref
        .read(sentDistributesProvider.notifier)
        .cancelDistribute(distribute.id);

    if (context.mounted) {
      Navigator.of(context).pop();
    }

    if (context.mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(success ? '取消成功' : '取消失败，请重试'),
          backgroundColor: success ? AppColors.success : AppColors.danger,
        ),
      );
    }
  }

  Widget _buildStatusChip(int status) {
    Color bgColor;
    Color textColor;
    String label;

    switch (status) {
      case 1:
        bgColor = AppColors.warning.withValues(alpha: 0.1);
        textColor = AppColors.warning;
        label = '待确认';
        break;
      case 2:
        bgColor = AppColors.success.withValues(alpha: 0.1);
        textColor = AppColors.success;
        label = '已确认';
        break;
      case 3:
        bgColor = AppColors.danger.withValues(alpha: 0.1);
        textColor = AppColors.danger;
        label = '已拒绝';
        break;
      case 4:
        bgColor = AppColors.textTertiary.withValues(alpha: 0.1);
        textColor = AppColors.textTertiary;
        label = '已取消';
        break;
      default:
        bgColor = AppColors.textTertiary.withValues(alpha: 0.1);
        textColor = AppColors.textTertiary;
        label = '未知';
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        label,
        style: TextStyle(fontSize: 12, color: textColor),
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 6),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 70,
            child: Text(
              label,
              style: const TextStyle(
                fontSize: 13,
                color: AppColors.textSecondary,
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(
                fontSize: 13,
                color: AppColors.textPrimary,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

Widget _buildErrorWidget(BuildContext context, String error, VoidCallback onRetry) {
  return Center(
    child: Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        const Icon(Icons.error_outline, size: 48, color: AppColors.danger),
        const SizedBox(height: AppSpacing.md),
        const Text(
          '加载失败',
          style: TextStyle(
            fontSize: 16,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: AppSpacing.sm),
        Text(
          error,
          style: const TextStyle(
            fontSize: 14,
            color: AppColors.textSecondary,
          ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: AppSpacing.md),
        ElevatedButton(
          onPressed: onRetry,
          child: const Text('重试'),
        ),
      ],
    ),
  );
}

Widget _buildEmptyWidget(String message) {
  return Center(
    child: Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        const Icon(
          Icons.inbox_outlined,
          size: 64,
          color: AppColors.textTertiary,
        ),
        const SizedBox(height: AppSpacing.md),
        Text(
          message,
          style: const TextStyle(
            fontSize: 16,
            color: AppColors.textSecondary,
          ),
        ),
      ],
    ),
  );
}
