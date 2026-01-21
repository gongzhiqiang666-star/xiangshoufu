import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../data/models/policy_model.dart';
import 'providers/policy_provider.dart';

/// 我的政策页面（只读）
class MyPolicyPage extends ConsumerWidget {
  const MyPolicyPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final policiesAsync = ref.watch(myPoliciesProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('我的政策'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => ref.invalidate(myPoliciesProvider),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async => ref.invalidate(myPoliciesProvider),
        child: policiesAsync.when(
          data: (policies) => policies.isEmpty
              ? const Center(child: Text('暂无政策配置'))
              : ListView.builder(
                  padding: const EdgeInsets.all(AppSpacing.md),
                  itemCount: policies.length,
                  itemBuilder: (context, index) => _PolicyCard(policy: policies[index]),
                ),
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (e, _) => Center(child: Text('加载失败: $e')),
        ),
      ),
    );
  }
}

class _PolicyCard extends StatelessWidget {
  final AgentPolicy policy;

  const _PolicyCard({required this.policy});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: AppSpacing.md),
      child: ExpansionTile(
        leading: CircleAvatar(
          backgroundColor: AppColors.primary.withOpacity(0.1),
          child: const Icon(Icons.policy, color: AppColors.primary),
        ),
        title: Text(
          policy.channelName ?? '通道${policy.channelId}',
          style: const TextStyle(fontWeight: FontWeight.bold),
        ),
        subtitle: policy.templateName != null
            ? Text('模版: ${policy.templateName}')
            : null,
        children: [
          Padding(
            padding: const EdgeInsets.all(AppSpacing.md),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildSectionTitle('费率配置'),
                _buildRateSection(),
                const Divider(height: 24),
                if (policy.depositCashbacks != null && policy.depositCashbacks!.isNotEmpty) ...[
                  _buildSectionTitle('押金返现'),
                  _buildDepositCashbackSection(),
                  const Divider(height: 24),
                ],
                if (policy.simCashback != null) ...[
                  _buildSectionTitle('流量卡返现'),
                  _buildSimCashbackSection(),
                  const Divider(height: 24),
                ],
                if (policy.activationRewards != null && policy.activationRewards!.isNotEmpty) ...[
                  _buildSectionTitle('激活奖励'),
                  _buildActivationRewardSection(),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSectionTitle(String title) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Text(
        title,
        style: const TextStyle(
          fontSize: 14,
          fontWeight: FontWeight.bold,
          color: AppColors.textSecondary,
        ),
      ),
    );
  }

  Widget _buildRateSection() {
    return Wrap(
      spacing: 16,
      runSpacing: 8,
      children: [
        _buildRateItem('贷记卡', policy.creditRate),
        _buildRateItem('借记卡', policy.debitRate),
        _buildRateItem('借记卡封顶', policy.debitCap, isCap: true),
        _buildRateItem('云闪付', policy.unionpayRate),
        _buildRateItem('微信', policy.wechatRate),
        _buildRateItem('支付宝', policy.alipayRate),
      ],
    );
  }

  Widget _buildRateItem(String label, String value, {bool isCap = false}) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: AppColors.primary.withOpacity(0.1),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Text(
        '$label: ${isCap ? '¥$value' : '$value%'}',
        style: const TextStyle(fontSize: 13),
      ),
    );
  }

  Widget _buildDepositCashbackSection() {
    return Wrap(
      spacing: 12,
      runSpacing: 8,
      children: policy.depositCashbacks!.map((item) {
        return Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
          decoration: BoxDecoration(
            color: AppColors.success.withOpacity(0.1),
            borderRadius: BorderRadius.circular(16),
          ),
          child: Text(
            '押金¥${item.depositAmountYuan.toInt()} → 返现¥${item.cashbackAmountYuan.toStringAsFixed(2)}',
            style: const TextStyle(fontSize: 13),
          ),
        );
      }).toList(),
    );
  }

  Widget _buildSimCashbackSection() {
    final sim = policy.simCashback!;
    return Wrap(
      spacing: 12,
      runSpacing: 8,
      children: [
        _buildCashbackChip('首次', sim.firstTimeCashbackYuan),
        _buildCashbackChip('二次', sim.secondTimeCashbackYuan),
        _buildCashbackChip('后续', sim.thirdPlusCashbackYuan),
      ],
    );
  }

  Widget _buildCashbackChip(String label, double amount) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: AppColors.warning.withOpacity(0.1),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Text(
        '$label: ¥${amount.toStringAsFixed(2)}',
        style: const TextStyle(fontSize: 13),
      ),
    );
  }

  Widget _buildActivationRewardSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: policy.activationRewards!.map((item) {
        return Padding(
          padding: const EdgeInsets.only(bottom: 8),
          child: Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.purple.withOpacity(0.1),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  item.rewardName,
                  style: const TextStyle(fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 4),
                Text(
                  item.summary,
                  style: const TextStyle(fontSize: 13, color: AppColors.textSecondary),
                ),
              ],
            ),
          ),
        );
      }).toList(),
    );
  }
}
