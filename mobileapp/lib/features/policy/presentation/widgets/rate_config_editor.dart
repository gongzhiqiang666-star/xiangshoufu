import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';

/// 费率配置编辑器（可编辑）
class RateConfigEditor extends StatelessWidget {
  final double creditRate;
  final double debitRate;
  final double debitCap;
  final double unionpayRate;
  final double wechatRate;
  final double alipayRate;

  // 最小值限制（自己的费率）
  final double minCreditRate;
  final double minDebitRate;
  final double minUnionpayRate;
  final double minWechatRate;
  final double minAlipayRate;

  final ValueChanged<double>? onCreditRateChanged;
  final ValueChanged<double>? onDebitRateChanged;
  final ValueChanged<double>? onDebitCapChanged;
  final ValueChanged<double>? onUnionpayRateChanged;
  final ValueChanged<double>? onWechatRateChanged;
  final ValueChanged<double>? onAlipayRateChanged;

  const RateConfigEditor({
    super.key,
    required this.creditRate,
    required this.debitRate,
    required this.debitCap,
    required this.unionpayRate,
    required this.wechatRate,
    required this.alipayRate,
    this.minCreditRate = 0,
    this.minDebitRate = 0,
    this.minUnionpayRate = 0,
    this.minWechatRate = 0,
    this.minAlipayRate = 0,
    this.onCreditRateChanged,
    this.onDebitRateChanged,
    this.onDebitCapChanged,
    this.onUnionpayRateChanged,
    this.onWechatRateChanged,
    this.onAlipayRateChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        _buildHeader(),
        const SizedBox(height: 16),

        // 贷记卡费率
        _buildRateField(
          label: '贷记卡费率',
          value: creditRate,
          minValue: minCreditRate,
          maxValue: 5.0,
          unit: '%',
          icon: Icons.credit_card,
          onChanged: onCreditRateChanged,
        ),
        const SizedBox(height: 12),

        // 借记卡费率
        _buildRateField(
          label: '借记卡费率',
          value: debitRate,
          minValue: minDebitRate,
          maxValue: 5.0,
          unit: '%',
          icon: Icons.account_balance,
          onChanged: onDebitRateChanged,
        ),
        const SizedBox(height: 12),

        // 借记卡封顶
        _buildRateField(
          label: '借记卡封顶',
          value: debitCap,
          minValue: 0,
          maxValue: 100,
          unit: '元',
          icon: Icons.vertical_align_top,
          step: 1,
          precision: 0,
          onChanged: onDebitCapChanged,
        ),
        const SizedBox(height: 12),

        // 云闪付费率
        _buildRateField(
          label: '云闪付费率',
          value: unionpayRate,
          minValue: minUnionpayRate,
          maxValue: 5.0,
          unit: '%',
          icon: Icons.flash_on,
          onChanged: onUnionpayRateChanged,
        ),
        const SizedBox(height: 12),

        // 微信扫码费率
        _buildRateField(
          label: '微信扫码费率',
          value: wechatRate,
          minValue: minWechatRate,
          maxValue: 5.0,
          unit: '%',
          icon: Icons.qr_code,
          onChanged: onWechatRateChanged,
        ),
        const SizedBox(height: 12),

        // 支付宝费率
        _buildRateField(
          label: '支付宝费率',
          value: alipayRate,
          minValue: minAlipayRate,
          maxValue: 5.0,
          unit: '%',
          icon: Icons.payment,
          onChanged: onAlipayRateChanged,
        ),

        const SizedBox(height: 16),
        _buildHint(),
      ],
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

  Widget _buildRateField({
    required String label,
    required double value,
    required double minValue,
    required double maxValue,
    required String unit,
    required IconData icon,
    double step = 0.01,
    int precision = 4,
    ValueChanged<double>? onChanged,
  }) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: Colors.grey[300]!),
      ),
      child: Row(
        children: [
          Icon(icon, size: 20, color: Colors.grey[600]),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  label,
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.grey[700],
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  '最低 $minValue$unit',
                  style: TextStyle(
                    fontSize: 11,
                    color: Colors.grey[500],
                  ),
                ),
              ],
            ),
          ),
          SizedBox(
            width: 120,
            child: Row(
              children: [
                InkWell(
                  onTap: () {
                    final newValue = (value - step).clamp(minValue, maxValue);
                    onChanged?.call(double.parse(newValue.toStringAsFixed(precision)));
                  },
                  child: Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: Colors.grey[100],
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: const Icon(Icons.remove, size: 16),
                  ),
                ),
                Expanded(
                  child: Text(
                    precision == 0 ? value.toInt().toString() : value.toStringAsFixed(2),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: AppColors.primary,
                    ),
                  ),
                ),
                InkWell(
                  onTap: () {
                    final newValue = (value + step).clamp(minValue, maxValue);
                    onChanged?.call(double.parse(newValue.toStringAsFixed(precision)));
                  },
                  child: Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: AppColors.primary.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Icon(Icons.add, size: 16, color: AppColors.primary),
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(width: 8),
          Text(
            unit,
            style: TextStyle(
              fontSize: 14,
              color: Colors.grey[600],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildHint() {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.blue.withOpacity(0.05),
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: Colors.blue.withOpacity(0.2)),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(Icons.info_outline, size: 16, color: Colors.blue[600]),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              '设置的费率不能低于您自己的费率，下级代理商的分润 = 下级费率 - 您的费率',
              style: TextStyle(
                fontSize: 12,
                color: Colors.blue[700],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
