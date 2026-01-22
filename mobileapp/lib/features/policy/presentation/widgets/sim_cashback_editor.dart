import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../data/models/policy_model.dart';

/// 流量卡返现编辑器
class SimCashbackEditor extends StatefulWidget {
  final SimCashbackConfig? initialConfig;
  final SimCashbackConfig? maxConfig;
  final ValueChanged<SimCashbackConfig> onChanged;

  const SimCashbackEditor({
    super.key,
    this.initialConfig,
    this.maxConfig,
    required this.onChanged,
  });

  @override
  State<SimCashbackEditor> createState() => _SimCashbackEditorState();
}

class _SimCashbackEditorState extends State<SimCashbackEditor> {
  late TextEditingController _firstController;
  late TextEditingController _secondController;
  late TextEditingController _thirdController;

  @override
  void initState() {
    super.initState();
    final config = widget.initialConfig;
    _firstController = TextEditingController(
      text: config?.firstTimeCashbackYuan.toStringAsFixed(2) ?? '0.00',
    );
    _secondController = TextEditingController(
      text: config?.secondTimeCashbackYuan.toStringAsFixed(2) ?? '0.00',
    );
    _thirdController = TextEditingController(
      text: config?.thirdPlusCashbackYuan.toStringAsFixed(2) ?? '0.00',
    );
  }

  @override
  void dispose() {
    _firstController.dispose();
    _secondController.dispose();
    _thirdController.dispose();
    super.dispose();
  }

  void _notifyChange() {
    widget.onChanged(SimCashbackConfig(
      firstTimeCashback: ((double.tryParse(_firstController.text) ?? 0) * 100).round(),
      secondTimeCashback: ((double.tryParse(_secondController.text) ?? 0) * 100).round(),
      thirdPlusCashback: ((double.tryParse(_thirdController.text) ?? 0) * 100).round(),
    ));
  }

  @override
  Widget build(BuildContext context) {
    final max = widget.maxConfig;

    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildInfoCard(),
          const SizedBox(height: 16),
          _buildCashbackField(
            label: '首次返现',
            controller: _firstController,
            maxValue: max?.firstTimeCashbackYuan,
            hint: '商户首次缴纳流量费时返现',
          ),
          _buildCashbackField(
            label: '二次返现',
            controller: _secondController,
            maxValue: max?.secondTimeCashbackYuan,
            hint: '商户第二次缴纳流量费时返现',
          ),
          _buildCashbackField(
            label: '后续返现',
            controller: _thirdController,
            maxValue: max?.thirdPlusCashbackYuan,
            hint: '商户第三次及以后返现',
          ),
        ],
      ),
    );
  }

  Widget _buildInfoCard() {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.warning.withOpacity(0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          const Icon(Icons.info_outline, color: AppColors.warning, size: 20),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              '商户缴纳流量费（99元/年）后，按次数给下级代理商返现。返现金额不能超过您的配置。',
              style: TextStyle(
                fontSize: 13,
                color: AppColors.warning.withOpacity(0.8),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCashbackField({
    required String label,
    required TextEditingController controller,
    double? maxValue,
    required String hint,
  }) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: TextField(
        controller: controller,
        keyboardType: const TextInputType.numberWithOptions(decimal: true),
        decoration: InputDecoration(
          labelText: label,
          hintText: maxValue != null ? '最高 ¥${maxValue.toStringAsFixed(2)}' : hint,
          helperText: hint,
          suffixText: '元',
          border: const OutlineInputBorder(),
          filled: true,
          fillColor: Colors.white,
        ),
        onChanged: (_) => _notifyChange(),
      ),
    );
  }
}
