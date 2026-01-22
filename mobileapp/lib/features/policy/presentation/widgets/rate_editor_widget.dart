import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../data/models/policy_model.dart';

/// 费率编辑器
class RateEditorWidget extends StatefulWidget {
  final RateConfig initialRates;
  final PolicyLimits limits;
  final ValueChanged<RateConfig> onChanged;

  const RateEditorWidget({
    super.key,
    required this.initialRates,
    required this.limits,
    required this.onChanged,
  });

  @override
  State<RateEditorWidget> createState() => _RateEditorWidgetState();
}

class _RateEditorWidgetState extends State<RateEditorWidget> {
  late TextEditingController _creditController;
  late TextEditingController _debitController;
  late TextEditingController _debitCapController;
  late TextEditingController _unionpayController;
  late TextEditingController _wechatController;
  late TextEditingController _alipayController;

  @override
  void initState() {
    super.initState();
    _creditController = TextEditingController(text: widget.initialRates.creditRate);
    _debitController = TextEditingController(text: widget.initialRates.debitRate);
    _debitCapController = TextEditingController(text: widget.initialRates.debitCap);
    _unionpayController = TextEditingController(text: widget.initialRates.unionpayRate);
    _wechatController = TextEditingController(text: widget.initialRates.wechatRate);
    _alipayController = TextEditingController(text: widget.initialRates.alipayRate);
  }

  @override
  void dispose() {
    _creditController.dispose();
    _debitController.dispose();
    _debitCapController.dispose();
    _unionpayController.dispose();
    _wechatController.dispose();
    _alipayController.dispose();
    super.dispose();
  }

  void _notifyChange() {
    widget.onChanged(RateConfig(
      creditRate: _creditController.text,
      debitRate: _debitController.text,
      debitCap: _debitCapController.text,
      unionpayRate: _unionpayController.text,
      wechatRate: _wechatController.text,
      alipayRate: _alipayController.text,
    ));
  }

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildInfoCard(),
          const SizedBox(height: 16),
          _buildRateField(
            label: '贷记卡费率',
            controller: _creditController,
            minValue: widget.limits.minCreditRate,
            suffix: '%',
            hint: '最低 ${widget.limits.minCreditRate}%',
          ),
          _buildRateField(
            label: '借记卡费率',
            controller: _debitController,
            minValue: widget.limits.minDebitRate,
            suffix: '%',
            hint: '最低 ${widget.limits.minDebitRate}%',
          ),
          _buildRateField(
            label: '借记卡封顶',
            controller: _debitCapController,
            suffix: '元',
            hint: '0表示不封顶',
            isAmount: true,
          ),
          _buildRateField(
            label: '云闪付费率',
            controller: _unionpayController,
            minValue: widget.limits.minUnionpayRate,
            suffix: '%',
            hint: '最低 ${widget.limits.minUnionpayRate}%',
          ),
          _buildRateField(
            label: '微信费率',
            controller: _wechatController,
            minValue: widget.limits.minWechatRate,
            suffix: '%',
            hint: '最低 ${widget.limits.minWechatRate}%',
          ),
          _buildRateField(
            label: '支付宝费率',
            controller: _alipayController,
            minValue: widget.limits.minAlipayRate,
            suffix: '%',
            hint: '最低 ${widget.limits.minAlipayRate}%',
          ),
        ],
      ),
    );
  }

  Widget _buildInfoCard() {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.primary.withOpacity(0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          const Icon(Icons.info_outline, color: AppColors.primary, size: 20),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              '给下级设置的费率不能低于您自己的费率，差价即为您的分润',
              style: TextStyle(
                fontSize: 13,
                color: AppColors.primary.withOpacity(0.8),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildRateField({
    required String label,
    required TextEditingController controller,
    String? minValue,
    required String suffix,
    required String hint,
    bool isAmount = false,
  }) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: TextField(
        controller: controller,
        keyboardType: const TextInputType.numberWithOptions(decimal: true),
        decoration: InputDecoration(
          labelText: label,
          hintText: hint,
          suffixText: suffix,
          border: const OutlineInputBorder(),
          filled: true,
          fillColor: Colors.white,
        ),
        onChanged: (_) => _notifyChange(),
      ),
    );
  }
}
