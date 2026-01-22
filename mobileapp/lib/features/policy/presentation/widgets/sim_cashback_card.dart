import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';
import '../../data/models/policy_model.dart';

/// 流量卡返现卡片（只读展示）
class SimCashbackCard extends StatelessWidget {
  final SimCashbackConfig config;
  final bool readonly;

  const SimCashbackCard({
    super.key,
    required this.config,
    this.readonly = true,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildHeader(),
            const SizedBox(height: 16),
            _buildContent(),
          ],
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Row(
      children: [
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
          decoration: BoxDecoration(
            color: Colors.blue.withOpacity(0.1),
            borderRadius: BorderRadius.circular(4),
          ),
          child: const Text(
            '服务费钱包',
            style: TextStyle(
              fontSize: 12,
              color: Colors.blue,
              fontWeight: FontWeight.w500,
            ),
          ),
        ),
        const SizedBox(width: 8),
        const Text(
          '流量卡返现配置',
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w600,
          ),
        ),
      ],
    );
  }

  Widget _buildContent() {
    return Column(
      children: [
        Row(
          children: [
            Expanded(
              child: _buildCashbackItem(
                '首次返现',
                config.firstTimeCashbackYuan,
                Icons.looks_one_outlined,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: _buildCashbackItem(
                '第2次返现',
                config.secondTimeCashbackYuan,
                Icons.looks_two_outlined,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: _buildCashbackItem(
                '第3次+返现',
                config.thirdPlusCashbackYuan,
                Icons.looks_3_outlined,
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        Container(
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: Colors.grey[50],
            borderRadius: BorderRadius.circular(8),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(Icons.sim_card_outlined, size: 16, color: Colors.grey[600]),
              const SizedBox(width: 8),
              Text(
                '流量费金额：${config.simFeeAmountYuan.toStringAsFixed(0)}元/年',
                style: TextStyle(
                  fontSize: 13,
                  color: Colors.grey[700],
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildCashbackItem(String label, double amount, IconData icon) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.grey[50],
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: Colors.grey[200]!),
      ),
      child: Column(
        children: [
          Icon(icon, size: 24, color: Colors.grey[500]),
          const SizedBox(height: 8),
          Text(
            label,
            style: TextStyle(
              fontSize: 12,
              color: Colors.grey[600],
            ),
          ),
          const SizedBox(height: 4),
          Text(
            '${amount.toStringAsFixed(2)}元',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.primary,
            ),
          ),
        ],
      ),
    );
  }
}
