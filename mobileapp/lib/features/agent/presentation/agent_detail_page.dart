import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../data/models/agent_model.dart';
import 'providers/agent_provider.dart';

/// 代理商详情页面
class AgentDetailPage extends ConsumerWidget {
  final int agentId;

  const AgentDetailPage({super.key, required this.agentId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final detailAsync = ref.watch(agentDetailProvider(agentId));

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('代理商详情'),
        actions: [
          PopupMenuButton<String>(
            icon: const Icon(Icons.more_vert),
            onSelected: (value) => _handleAction(context, ref, value),
            itemBuilder: (context) => [
              const PopupMenuItem(
                value: 'policy',
                child: Row(
                  children: [
                    Icon(Icons.policy_outlined, size: 18, color: AppColors.primary),
                    SizedBox(width: 8),
                    Text('设置政策'),
                  ],
                ),
              ),
              const PopupMenuItem(
                value: 'channels',
                child: Row(
                  children: [
                    Icon(Icons.account_tree_outlined, size: 18, color: AppColors.info),
                    SizedBox(width: 8),
                    Text('通道管理'),
                  ],
                ),
              ),
            ],
          ),
        ],
      ),
      body: detailAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.error_outline, size: 48, color: AppColors.danger),
              const SizedBox(height: 16),
              Text('加载失败: $error', textAlign: TextAlign.center),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () => ref.invalidate(agentDetailProvider(agentId)),
                child: const Text('重试'),
              ),
            ],
          ),
        ),
        data: (detail) => _buildContent(context, ref, detail),
      ),
    );
  }

  Widget _buildContent(BuildContext context, WidgetRef ref, AgentDetail detail) {
    final isActive = detail.status == 1;

    return RefreshIndicator(
      onRefresh: () async {
        ref.invalidate(agentDetailProvider(agentId));
      },
      child: SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        padding: const EdgeInsets.all(AppSpacing.md),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // 头像和基本信息卡片
            _buildProfileCard(context, detail, isActive),
            const SizedBox(height: AppSpacing.md),

            // 统计数据卡片
            _buildStatsCard(detail),
            const SizedBox(height: AppSpacing.md),

            // 详细信息卡片
            _buildDetailCard(context, detail),
            const SizedBox(height: AppSpacing.md),

            // 结算信息卡片
            _buildSettlementCard(detail),
          ],
        ),
      ),
    );
  }

  /// 头像和基本信息卡片
  Widget _buildProfileCard(BuildContext context, AgentDetail detail, bool isActive) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.lg),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [
            AppColors.primary,
            AppColors.primary.withValues(alpha: 0.8),
          ],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        children: [
          // 头像
          Container(
            width: 72,
            height: 72,
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(20),
            ),
            child: Center(
              child: Text(
                detail.agentName.isNotEmpty ? detail.agentName.substring(0, 1) : '?',
                style: const TextStyle(
                  fontSize: 32,
                  fontWeight: FontWeight.bold,
                  color: AppColors.primary,
                ),
              ),
            ),
          ),
          const SizedBox(height: 12),
          // 名称
          Text(
            detail.agentName,
            style: const TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          const SizedBox(height: 4),
          // 编号
          Text(
            detail.agentNo,
            style: TextStyle(
              fontSize: 14,
              color: Colors.white.withValues(alpha: 0.8),
            ),
          ),
          const SizedBox(height: 12),
          // 状态和层级标签
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                decoration: BoxDecoration(
                  color: isActive
                      ? AppColors.success.withValues(alpha: 0.2)
                      : AppColors.danger.withValues(alpha: 0.2),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Text(
                  detail.statusName ?? (isActive ? '正常' : '禁用'),
                  style: TextStyle(
                    fontSize: 12,
                    color: isActive ? Colors.white : Colors.white70,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ),
              const SizedBox(width: 8),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                decoration: BoxDecoration(
                  color: Colors.white.withValues(alpha: 0.2),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Text(
                  '${detail.level}级代理',
                  style: const TextStyle(
                    fontSize: 12,
                    color: Colors.white,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  /// 统计数据卡片
  Widget _buildStatsCard(AgentDetail detail) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '团队数据',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: AppSpacing.md),
          GridView.count(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            crossAxisCount: 2,
            crossAxisSpacing: AppSpacing.cardGap,
            mainAxisSpacing: AppSpacing.cardGap,
            childAspectRatio: 2.2,
            children: [
              _buildStatItem(
                '直属代理',
                detail.directAgentCount.toString(),
                Icons.person,
                AppColors.primary,
              ),
              _buildStatItem(
                '团队代理',
                detail.teamAgentCount.toString(),
                Icons.groups,
                AppColors.profitReward,
              ),
              _buildStatItem(
                '直营商户',
                detail.directMerchantCount.toString(),
                Icons.store,
                AppColors.success,
              ),
              _buildStatItem(
                '团队商户',
                detail.teamMerchantCount.toString(),
                Icons.storefront,
                AppColors.warning,
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildStatItem(String label, String value, IconData icon, Color color) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.05),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Container(
            width: 36,
            height: 36,
            decoration: BoxDecoration(
              color: color.withValues(alpha: 0.15),
              borderRadius: BorderRadius.circular(10),
            ),
            child: Icon(icon, color: color, size: 20),
          ),
          const SizedBox(width: 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(
                  value,
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: color,
                  ),
                ),
                Text(
                  label,
                  style: const TextStyle(
                    fontSize: 11,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  /// 详细信息卡片
  Widget _buildDetailCard(BuildContext context, AgentDetail detail) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '基本信息',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: AppSpacing.md),
          _buildInfoRow('联系人', detail.contactName ?? '-'),
          _buildInfoRow('联系电话', detail.contactPhone, onTap: () {
            Clipboard.setData(ClipboardData(text: detail.contactPhone));
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('电话已复制')),
            );
          }),
          _buildInfoRow('身份证号', detail.idCardNo ?? '-'),
          _buildInfoRow('上级代理', detail.parentName ?? '无（顶级代理）'),
          _buildInfoRow('邀请码', detail.inviteCode ?? '-', onTap: detail.inviteCode != null ? () {
            Clipboard.setData(ClipboardData(text: detail.inviteCode!));
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('邀请码已复制')),
            );
          } : null),
          _buildInfoRow('注册时间', detail.registerTime ?? '-', showDivider: false),
        ],
      ),
    );
  }

  /// 结算信息卡片
  Widget _buildSettlementCard(AgentDetail detail) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '结算信息',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: AppSpacing.md),
          _buildInfoRow('开户行', detail.bankName ?? '-'),
          _buildInfoRow('开户名', detail.bankAccount ?? '-'),
          _buildInfoRow('银行卡号', detail.bankCardNo ?? '-', showDivider: false),
        ],
      ),
    );
  }

  Widget _buildInfoRow(String label, String value, {bool showDivider = true, VoidCallback? onTap}) {
    return Column(
      children: [
        InkWell(
          onTap: onTap,
          child: Padding(
            padding: const EdgeInsets.symmetric(vertical: 12),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  label,
                  style: const TextStyle(
                    fontSize: 14,
                    color: AppColors.textSecondary,
                  ),
                ),
                Expanded(
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      Flexible(
                        child: Text(
                          value,
                          style: const TextStyle(
                            fontSize: 14,
                            color: AppColors.textPrimary,
                          ),
                          textAlign: TextAlign.right,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      if (onTap != null) ...[
                        const SizedBox(width: 4),
                        const Icon(Icons.copy, size: 14, color: AppColors.textTertiary),
                      ],
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
        if (showDivider) const Divider(height: 1, color: AppColors.divider),
      ],
    );
  }

  void _handleAction(BuildContext context, WidgetRef ref, String action) {
    switch (action) {
      case 'policy':
        context.push('/agent/$agentId/policy');
        break;
      case 'channels':
        context.push('/agent/$agentId/channels');
        break;
    }
  }
}
