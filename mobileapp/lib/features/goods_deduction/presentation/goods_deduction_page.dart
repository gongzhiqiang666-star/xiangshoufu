import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../router/app_router.dart';
import 'providers/goods_deduction_provider.dart';
import 'widgets/goods_deduction_card.dart';
import 'widgets/agreement_dialog.dart';

/// 货款代扣主页面
class GoodsDeductionPage extends ConsumerStatefulWidget {
  const GoodsDeductionPage({super.key});

  @override
  ConsumerState<GoodsDeductionPage> createState() => _GoodsDeductionPageState();
}

class _GoodsDeductionPageState extends ConsumerState<GoodsDeductionPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);

    // 加载数据
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(sentDeductionsProvider.notifier).loadDeductions(refresh: true);
      ref.read(receivedDeductionsProvider.notifier).loadDeductions(refresh: true);
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
        title: const Text('货款代扣'),
        centerTitle: true,
        bottom: TabBar(
          controller: _tabController,
          labelColor: AppColors.primary,
          unselectedLabelColor: AppColors.textSecondary,
          indicatorColor: AppColors.primary,
          tabs: const [
            Tab(text: '我发起的'),
            Tab(text: '我接收的'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _SentDeductionsTab(),
          _ReceivedDeductionsTab(),
        ],
      ),
    );
  }
}

/// 我发起的 Tab
class _SentDeductionsTab extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(sentDeductionsProvider);

    if (state.isLoading && state.list.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (state.error != null && state.list.isEmpty) {
      return _buildErrorWidget(
        context,
        state.error!,
        () => ref.read(sentDeductionsProvider.notifier).loadDeductions(refresh: true),
      );
    }

    if (state.list.isEmpty) {
      return _buildEmptyWidget('暂无发起的货款代扣');
    }

    return RefreshIndicator(
      onRefresh: () async {
        await ref.read(sentDeductionsProvider.notifier).loadDeductions(refresh: true);
      },
      child: NotificationListener<ScrollNotification>(
        onNotification: (notification) {
          if (notification is ScrollEndNotification &&
              notification.metrics.pixels >= notification.metrics.maxScrollExtent - 100 &&
              state.hasMore &&
              !state.isLoadingMore) {
            ref.read(sentDeductionsProvider.notifier).loadDeductions();
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

            final deduction = state.list[index];
            return GoodsDeductionCard(
              deduction: deduction,
              isSent: true,
              onTap: () => _navigateToDetail(context, deduction.id),
            );
          },
        ),
      ),
    );
  }

  void _navigateToDetail(BuildContext context, int id) {
    context.push(
      RoutePaths.goodsDeductionDetail.replaceFirst(':id', id.toString()),
      extra: {'isSent': true},
    );
  }
}

/// 我接收的 Tab
class _ReceivedDeductionsTab extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(receivedDeductionsProvider);

    if (state.isLoading && state.list.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (state.error != null && state.list.isEmpty) {
      return _buildErrorWidget(
        context,
        state.error!,
        () => ref.read(receivedDeductionsProvider.notifier).loadDeductions(refresh: true),
      );
    }

    if (state.list.isEmpty) {
      return _buildEmptyWidget('暂无接收的货款代扣');
    }

    return RefreshIndicator(
      onRefresh: () async {
        await ref.read(receivedDeductionsProvider.notifier).loadDeductions(refresh: true);
      },
      child: NotificationListener<ScrollNotification>(
        onNotification: (notification) {
          if (notification is ScrollEndNotification &&
              notification.metrics.pixels >= notification.metrics.maxScrollExtent - 100 &&
              state.hasMore &&
              !state.isLoadingMore) {
            ref.read(receivedDeductionsProvider.notifier).loadDeductions();
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

            final deduction = state.list[index];
            return GoodsDeductionCard(
              deduction: deduction,
              isSent: false,
              onTap: () => _navigateToDetail(context, ref, deduction.id),
              onAccept: deduction.status == 1
                  ? () => _handleAccept(context, ref, deduction)
                  : null,
              onReject: deduction.status == 1
                  ? () => _handleReject(context, ref, deduction)
                  : null,
            );
          },
        ),
      ),
    );
  }

  void _navigateToDetail(BuildContext context, WidgetRef ref, int id) {
    context.push(
      RoutePaths.goodsDeductionDetail.replaceFirst(':id', id.toString()),
      extra: {'isSent': false},
    );
  }

  Future<void> _handleAccept(
    BuildContext context,
    WidgetRef ref,
    dynamic deduction,
  ) async {
    // 显示协议弹窗
    final agreed = await showAgreementDialog(
      context: context,
      title: '代扣服务协议',
      content: getDefaultAgreementContent(
        fromAgentName: deduction.fromAgentName,
        toAgentName: deduction.toAgentName,
        totalAmount: deduction.totalAmountYuan,
        terminalCount: deduction.terminalCount,
      ),
    );

    if (agreed == true) {
      final success = await ref
          .read(receivedDeductionsProvider.notifier)
          .acceptDeduction(deduction.id);

      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(success ? '接收成功，代扣已开始' : '操作失败，请重试'),
            backgroundColor: success ? AppColors.success : AppColors.danger,
          ),
        );
      }
    }
  }

  Future<void> _handleReject(
    BuildContext context,
    WidgetRef ref,
    dynamic deduction,
  ) async {
    final reasonController = TextEditingController();

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('拒绝货款代扣'),
        content: TextField(
          controller: reasonController,
          maxLines: 3,
          decoration: const InputDecoration(
            hintText: '请输入拒绝原因',
            border: OutlineInputBorder(),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () {
              if (reasonController.text.isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('请输入拒绝原因')),
                );
                return;
              }
              Navigator.of(context).pop(true);
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.danger,
              foregroundColor: Colors.white,
            ),
            child: const Text('确认拒绝'),
          ),
        ],
      ),
    );

    if (confirmed == true && reasonController.text.isNotEmpty) {
      final success = await ref
          .read(receivedDeductionsProvider.notifier)
          .rejectDeduction(deduction.id, reasonController.text);

      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(success ? '已拒绝' : '操作失败，请重试'),
            backgroundColor: success ? AppColors.success : AppColors.danger,
          ),
        );
      }
    }

    reasonController.dispose();
  }
}

Widget _buildErrorWidget(BuildContext context, String error, VoidCallback onRetry) {
  return Center(
    child: Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        const Icon(Icons.error_outline, size: 48, color: AppColors.danger),
        const SizedBox(height: AppSpacing.md),
        Text(
          '加载失败',
          style: const TextStyle(
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
