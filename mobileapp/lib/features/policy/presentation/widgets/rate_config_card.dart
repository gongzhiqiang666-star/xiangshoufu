import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';

/// 费率配置卡片（只读展示）
class RateConfigCard extends StatelessWidget {
  final double creditRate;
  final double debitRate;
  final double debitCap;
  final double unionpayRate;
  final double wechatRate;
  final double alipayRate;
  final bool readonly;

  const RateConfigCard({
    super.key,
    required this.creditRate,
    required this.debitRate,
    required this.debitCap,
    required this.unionpayRate,
    required this.wechatRate,
    required this.alipayRate,
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
            _buildRateGrid(),
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
            color: Colors.green.withOpacity(0.1),
            borderRadius: BorderRadius.circular(4),
          ),
          child: const Text(
            '分润钱包',
            style: TextStyle(
              fontSize: 12,
              color: Colors.green,
              fontWeight: FontWeight.w500,
            ),
          ),
        ),
        const SizedBox(width: 8),
        const Text(
          '成本费率配置',
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w600,
          ),
        ),
      ],
    );
  }

  Widget _buildRateGrid() {
    return Column(
      children: [
        Row(
          children: [
            Expanded(child: _buildRateItem('贷记卡', '$creditRate%', Icons.credit_card)),
            const SizedBox(width: 12),
            Expanded(child: _buildRateItem('借记卡', '$debitRate%', Icons.account_balance)),
            const SizedBox(width: 12),
            Expanded(child: _buildRateItem('借记卡封顶', '$debitCap元', Icons.vertical_align_top)),
          ],
        ),
        const SizedBox(height: 12),
        Row(
          children: [
            Expanded(child: _buildRateItem('云闪付', '$unionpayRate%', Icons.flash_on)),
            const SizedBox(width: 12),
            Expanded(child: _buildRateItem('微信扫码', '$wechatRate%', Icons.qr_code)),
            const SizedBox(width: 12),
            Expanded(child: _buildRateItem('支付宝', '$alipayRate%', Icons.payment)),
          ],
        ),
      ],
    );
  }

  Widget _buildRateItem(String label, String value, IconData icon) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.grey[50],
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: Colors.grey[200]!),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(icon, size: 16, color: Colors.grey[600]),
              const SizedBox(width: 4),
              Flexible(
                child: Text(
                  label,
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey[600],
                  ),
                  overflow: TextOverflow.ellipsis,
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            value,
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w600,
              color: AppColors.primary,
            ),
          ),
        ],
      ),
    );
  }
}
