import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../router/app_router.dart';
import '../data/models/deduction_model.dart';
import 'providers/deduction_provider.dart';
import 'widgets/deduction_card.dart';

/// 代扣管理主页面
class DeductionPage extends ConsumerStatefulWidget {
  const DeductionPage({super.key});

  @override
  ConsumerState<DeductionPage> createState() => _DeductionPageState();
}

class _DeductionPageState extends ConsumerState<DeductionPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    _tabController.addListener(_onTabChanged);
    // 加载数据
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(deductionPlansProvider.notifier).loadPlans(refresh: true);
    });
  }

  @override
  void dispose() {
    _tabController.removeListener(_onTabChanged);
    _tabController.dispose();
    super.dispose();
  }

  void _onTabChanged() {
    if (!_tabController.indexIsChanging) {
      final mode = _tabController.index == 0
          ? DeductionListMode.received
          : DeductionListMode.sent;
      ref.read(deductionPlansProvider.notifier).setListMode(mode);
    }
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
        bottom: TabBar(
          controller: _tabController,
          labelColor: AppColors.primary,
          unselectedLabelColor: AppColors.textSecondary,
          indicatorColor: AppColors.primary,
          tabs: const [
            Tab(text: '我接收的'),
            Tab(text: '我发起的'),
          ],
        ),
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
            final isReceived = state.listMode == DeductionListMode.received;
            return DeductionCard(
              plan: plan,
              showDeductor: isReceived,
              onTap: () => _navigateToDetail(plan.id),
              onAccept: plan.statusEnum.canAccept ? () => _handleAccept(plan) : null,
              onReject: plan.statusEnum.canReject ? () => _handleReject(plan) : null,
              onPause: plan.statusEnum.canPause ? () => _handlePause(plan) : null,
              onResume: plan.statusEnum.canResume ? () => _handleResume(plan) : null,
              onCancel: plan.statusEnum.canCancel ? () => _handleCancel(plan) : null,
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

  Future<void> _handleAccept(DeductionPlan plan) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('接收确认'),
        content: Text('确认接收代扣计划 ${plan.planNo}？\n\n接收后，系统将自动冻结您的相应余额用于扣款。'),
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
            child: const Text('确认接收'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      final success =
          await ref.read(deductionPlansProvider.notifier).acceptPlan(plan.id);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(success ? '已接收，开始冻结余额' : '操作失败'),
            backgroundColor: success ? AppColors.success : AppColors.danger,
          ),
        );
      }
    }
  }

  Future<void> _handleReject(DeductionPlan plan) async {
    final reasonController = TextEditingController();
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('拒绝代扣'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('确定要拒绝代扣计划 ${plan.planNo} 吗？'),
            const SizedBox(height: 16),
            TextField(
              controller: reasonController,
              decoration: const InputDecoration(
                labelText: '拒绝原因（可选）',
                border: OutlineInputBorder(),
              ),
              maxLines: 2,
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.danger,
              foregroundColor: Colors.white,
            ),
            child: const Text('确定拒绝'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      final reason = reasonController.text.trim();
      final success = await ref.read(deductionPlansProvider.notifier).rejectPlan(
            plan.id,
            reason: reason.isNotEmpty ? reason : null,
          );
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(success ? '已拒绝' : '操作失败'),
            backgroundColor: success ? AppColors.success : AppColors.danger,
          ),
        );
      }
    }
    reasonController.dispose();
  }

  Future<void> _handlePause(DeductionPlan plan) async {
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

  Future<void> _handleResume(DeductionPlan plan) async {
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

  Future<void> _handleCancel(DeductionPlan plan) async {
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
                _buildStatusFilterChip('全部', null, state.statusFilter),
                _buildStatusFilterChip('待接收', '0', state.statusFilter),
                _buildStatusFilterChip('进行中', '1', state.statusFilter),
                _buildStatusFilterChip('已完成', '2', state.statusFilter),
                _buildStatusFilterChip('已暂停', '3', state.statusFilter),
                _buildStatusFilterChip('已拒绝', '5', state.statusFilter),
              ],
            ),
            const SizedBox(height: AppSpacing.md),
            const Text('类型', style: TextStyle(fontWeight: FontWeight.w500)),
            const SizedBox(height: AppSpacing.sm),
            Wrap(
              spacing: 8,
              children: [
                _buildTypeFilterChip('全部', null, state.typeFilter),
                _buildTypeFilterChip('货款代扣', 1, state.typeFilter),
                _buildTypeFilterChip('伙伴代扣', 2, state.typeFilter),
                _buildTypeFilterChip('押金代扣', 3, state.typeFilter),
              ],
            ),
            const SizedBox(height: AppSpacing.lg),
          ],
        ),
      ),
    );
  }

  Widget _buildStatusFilterChip(String label, String? value, String? currentValue) {
    final isSelected = value == currentValue;
    return ChoiceChip(
      label: Text(label),
      selected: isSelected,
      onSelected: (_) {
        ref.read(deductionPlansProvider.notifier).setStatusFilter(value);
        Navigator.of(context).pop();
      },
      selectedColor: AppColors.primary.withOpacity(0.2),
      labelStyle: TextStyle(
        color: isSelected ? AppColors.primary : AppColors.textSecondary,
      ),
    );
  }

  Widget _buildTypeFilterChip(String label, int? value, int? currentValue) {
    final isSelected = value == currentValue;
    return ChoiceChip(
      label: Text(label),
      selected: isSelected,
      onSelected: (_) {
        ref.read(deductionPlansProvider.notifier).setTypeFilter(value);
        Navigator.of(context).pop();
      },
      selectedColor: AppColors.primary.withOpacity(0.2),
      labelStyle: TextStyle(
        color: isSelected ? AppColors.primary : AppColors.textSecondary,
      ),
    );
  }
}
