import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../router/app_router.dart';
import 'providers/deduction_provider.dart';
import 'widgets/deduction_card.dart';

/// 代扣管理主页面
class DeductionPage extends ConsumerStatefulWidget {
  const DeductionPage({super.key});

  @override
  ConsumerState<DeductionPage> createState() => _DeductionPageState();
}

class _DeductionPageState extends ConsumerState<DeductionPage> {
  @override
  void initState() {
    super.initState();
    // 加载数据
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(deductionPlansProvider.notifier).loadPlans(refresh: true);
    });
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(deductionPlansProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('代扣管理'),
        centerTitle: true,
        actions: [
          IconButton(
            icon: const Icon(Icons.filter_list),
            onPressed: () => _showFilterSheet(context),
          ),
        ],
      ),
      body: _buildBody(state),
    );
  }

  Widget _buildBody(DeductionPlansState state) {
    if (state.isLoading && state.list.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (state.error != null && state.list.isEmpty) {
      return _buildErrorWidget(state.error!);
    }

    if (state.list.isEmpty) {
      return _buildEmptyWidget();
    }

    return RefreshIndicator(
      onRefresh: () async {
        await ref.read(deductionPlansProvider.notifier).loadPlans(refresh: true);
      },
      child: NotificationListener<ScrollNotification>(
        onNotification: (notification) {
          if (notification is ScrollEndNotification &&
              notification.metrics.pixels >= notification.metrics.maxScrollExtent - 100 &&
              state.hasMore &&
              !state.isLoadingMore) {
            ref.read(deductionPlansProvider.notifier).loadPlans();
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

            final plan = state.list[index];
            return DeductionCard(
              plan: plan,
              onTap: () => _navigateToDetail(plan.id),
              onPause: plan.status == 1 ? () => _handlePause(plan) : null,
              onResume: plan.status == 3 ? () => _handleResume(plan) : null,
              onCancel: (plan.status == 1 || plan.status == 3)
                  ? () => _handleCancel(plan)
                  : null,
            );
          },
        ),
      ),
    );
  }

  Widget _buildErrorWidget(String error) {
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
            onPressed: () =>
                ref.read(deductionPlansProvider.notifier).loadPlans(refresh: true),
            child: const Text('重试'),
          ),
        ],
      ),
    );
  }

  Widget _buildEmptyWidget() {
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
          const Text(
            '暂无代扣计划',
            style: TextStyle(
              fontSize: 16,
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }

  void _navigateToDetail(int id) {
    context.push(RoutePaths.deductionDetail.replaceFirst(':id', id.toString()));
  }

  Future<void> _handlePause(dynamic plan) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('暂停代扣'),
        content: Text('确定要暂停代扣计划 ${plan.planNo} 吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.warning,
              foregroundColor: Colors.white,
            ),
            child: const Text('确定暂停'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      final success =
          await ref.read(deductionPlansProvider.notifier).pausePlan(plan.id);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(success ? '暂停成功' : '操作失败'),
            backgroundColor: success ? AppColors.success : AppColors.danger,
          ),
        );
      }
    }
  }

  Future<void> _handleResume(dynamic plan) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('恢复代扣'),
        content: Text('确定要恢复代扣计划 ${plan.planNo} 吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.success,
              foregroundColor: Colors.white,
            ),
            child: const Text('确定恢复'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      final success =
          await ref.read(deductionPlansProvider.notifier).resumePlan(plan.id);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(success ? '恢复成功' : '操作失败'),
            backgroundColor: success ? AppColors.success : AppColors.danger,
          ),
        );
      }
    }
  }

  Future<void> _handleCancel(dynamic plan) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('取消代扣'),
        content: Text('确定要取消代扣计划 ${plan.planNo} 吗？取消后不可恢复。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('返回'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.danger,
              foregroundColor: Colors.white,
            ),
            child: const Text('确定取消'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      final success =
          await ref.read(deductionPlansProvider.notifier).cancelPlan(plan.id);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(success ? '已取消' : '操作失败'),
            backgroundColor: success ? AppColors.success : AppColors.danger,
          ),
        );
      }
    }
  }

  void _showFilterSheet(BuildContext context) {
    final state = ref.read(deductionPlansProvider);

    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (context) => Container(
        padding: const EdgeInsets.all(AppSpacing.md),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  '筛选',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                TextButton(
                  onPressed: () {
                    ref.read(deductionPlansProvider.notifier).setStatusFilter(null);
                    ref.read(deductionPlansProvider.notifier).setTypeFilter(null);
                    Navigator.of(context).pop();
                  },
                  child: const Text('重置'),
                ),
              ],
            ),
            const SizedBox(height: AppSpacing.md),
            const Text('状态', style: TextStyle(fontWeight: FontWeight.w500)),
            const SizedBox(height: AppSpacing.sm),
            Wrap(
              spacing: 8,
              children: [
                _buildFilterChip('全部', null, state.statusFilter, (v) {
                  ref.read(deductionPlansProvider.notifier).setStatusFilter(v);
                  Navigator.of(context).pop();
                }),
                _buildFilterChip('进行中', 1, state.statusFilter, (v) {
                  ref.read(deductionPlansProvider.notifier).setStatusFilter(v);
                  Navigator.of(context).pop();
                }),
                _buildFilterChip('已完成', 2, state.statusFilter, (v) {
                  ref.read(deductionPlansProvider.notifier).setStatusFilter(v);
                  Navigator.of(context).pop();
                }),
                _buildFilterChip('已暂停', 3, state.statusFilter, (v) {
                  ref.read(deductionPlansProvider.notifier).setStatusFilter(v);
                  Navigator.of(context).pop();
                }),
              ],
            ),
            const SizedBox(height: AppSpacing.md),
            const Text('类型', style: TextStyle(fontWeight: FontWeight.w500)),
            const SizedBox(height: AppSpacing.sm),
            Wrap(
              spacing: 8,
              children: [
                _buildFilterChip('全部', null, state.typeFilter, (v) {
                  ref.read(deductionPlansProvider.notifier).setTypeFilter(v);
                  Navigator.of(context).pop();
                }),
                _buildFilterChip('货款代扣', 1, state.typeFilter, (v) {
                  ref.read(deductionPlansProvider.notifier).setTypeFilter(v);
                  Navigator.of(context).pop();
                }),
                _buildFilterChip('伙伴代扣', 2, state.typeFilter, (v) {
                  ref.read(deductionPlansProvider.notifier).setTypeFilter(v);
                  Navigator.of(context).pop();
                }),
                _buildFilterChip('押金代扣', 3, state.typeFilter, (v) {
                  ref.read(deductionPlansProvider.notifier).setTypeFilter(v);
                  Navigator.of(context).pop();
                }),
              ],
            ),
            const SizedBox(height: AppSpacing.lg),
          ],
        ),
      ),
    );
  }

  Widget _buildFilterChip(
    String label,
    int? value,
    int? currentValue,
    Function(int?) onSelected,
  ) {
    final isSelected = value == currentValue;
    return ChoiceChip(
      label: Text(label),
      selected: isSelected,
      onSelected: (_) => onSelected(value),
      selectedColor: AppColors.primary.withOpacity(0.2),
      labelStyle: TextStyle(
        color: isSelected ? AppColors.primary : AppColors.textSecondary,
      ),
    );
  }
}
