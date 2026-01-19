import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';

/// 终端划拨页面
class TerminalTransferPage extends StatefulWidget {
  final List<String> selectedSNs;

  const TerminalTransferPage({
    super.key,
    required this.selectedSNs,
  });

  @override
  State<TerminalTransferPage> createState() => _TerminalTransferPageState();
}

class _TerminalTransferPageState extends State<TerminalTransferPage> {
  String? _selectedAgentId;
  bool _enableCargoDeduction = false;
  double _unitPrice = 50.0;
  final Set<String> _selectedWallets = {'profit'};

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('终端划拨')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(AppSpacing.md),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 已选终端
            _buildSelectedTerminals(),
            const SizedBox(height: AppSpacing.md),

            // 划拨给
            _buildAgentSelector(),
            const SizedBox(height: AppSpacing.md),

            // 货款代扣设置
            _buildCargoDeductionSettings(),
            const SizedBox(height: AppSpacing.md),

            // 提示信息
            _buildWarningNotice(),
          ],
        ),
      ),
      bottomNavigationBar: _buildBottomBar(),
    );
  }

  Widget _buildSelectedTerminals() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '已选终端: ${widget.selectedSNs.length}台',
            style: const TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'SN: ${widget.selectedSNs.join(", ")}',
            style: const TextStyle(
              fontSize: 14,
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAgentSelector() {
    final agents = [
      {'id': 'A002', 'name': '李四', 'phone': '139****9999'},
      {'id': 'A003', 'name': '王五', 'phone': '137****7777'},
      {'id': 'A004', 'name': '赵六', 'phone': '136****6666'},
    ];

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
            '划拨给:',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 12),
          // 搜索框
          TextField(
            decoration: InputDecoration(
              hintText: '搜索直属下级代理商',
              prefixIcon: const Icon(Icons.search, color: AppColors.textTertiary),
              filled: true,
              fillColor: AppColors.background,
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(8),
                borderSide: BorderSide.none,
              ),
              contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            ),
          ),
          const SizedBox(height: 12),
          const Text(
            '直属下级代理商',
            style: TextStyle(
              fontSize: 14,
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: 8),
          ...agents.map((agent) => _buildAgentItem(agent)),
        ],
      ),
    );
  }

  Widget _buildAgentItem(Map<String, String> agent) {
    final isSelected = _selectedAgentId == agent['id'];

    return GestureDetector(
      onTap: () {
        setState(() {
          _selectedAgentId = agent['id'];
        });
      },
      child: Container(
        margin: const EdgeInsets.only(bottom: 8),
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: isSelected ? AppColors.primary.withOpacity(0.05) : AppColors.background,
          border: Border.all(
            color: isSelected ? AppColors.primary : Colors.transparent,
          ),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Row(
          children: [
            Icon(
              isSelected ? Icons.radio_button_checked : Icons.radio_button_off,
              color: isSelected ? AppColors.primary : AppColors.textTertiary,
              size: 20,
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    '${agent['name']} (${agent['id']})',
                    style: const TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w500,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    '手机: ${agent['phone']}',
                    style: const TextStyle(
                      fontSize: 13,
                      color: AppColors.textSecondary,
                    ),
                  ),
                ],
              ),
            ),
            if (isSelected)
              const Text(
                '已选择',
                style: TextStyle(
                  fontSize: 12,
                  color: AppColors.primary,
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildCargoDeductionSettings() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Checkbox(
                value: _enableCargoDeduction,
                onChanged: (value) {
                  setState(() {
                    _enableCargoDeduction = value ?? false;
                  });
                },
                activeColor: AppColors.primary,
              ),
              const Text(
                '设置货款代扣',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w500,
                  color: AppColors.textPrimary,
                ),
              ),
            ],
          ),
          if (_enableCargoDeduction) ...[
            const Divider(),
            const SizedBox(height: 8),
            Row(
              children: [
                const Text('单价: ¥', style: TextStyle(fontSize: 14)),
                const SizedBox(width: 8),
                SizedBox(
                  width: 80,
                  child: TextField(
                    keyboardType: TextInputType.number,
                    textAlign: TextAlign.center,
                    decoration: InputDecoration(
                      contentPadding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
                      border: OutlineInputBorder(borderRadius: BorderRadius.circular(6)),
                    ),
                    controller: TextEditingController(text: _unitPrice.toStringAsFixed(0)),
                    onChanged: (value) {
                      _unitPrice = double.tryParse(value) ?? 50.0;
                    },
                  ),
                ),
                const Text(' 元/台', style: TextStyle(fontSize: 14)),
              ],
            ),
            const SizedBox(height: 12),
            Text(
              '总金额: ¥${(_unitPrice * widget.selectedSNs.length).toStringAsFixed(2)}',
              style: const TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: AppColors.primary,
              ),
            ),
            const SizedBox(height: 16),
            const Text(
              '扣款来源:',
              style: TextStyle(fontSize: 14, color: AppColors.textSecondary),
            ),
            const SizedBox(height: 8),
            _buildWalletCheckbox('分润钱包', 'profit'),
            _buildWalletCheckbox('服务费钱包', 'service'),
            _buildWalletCheckbox('奖励钱包', 'reward'),
          ],
        ],
      ),
    );
  }

  Widget _buildWalletCheckbox(String label, String value) {
    return Row(
      children: [
        Checkbox(
          value: _selectedWallets.contains(value),
          onChanged: (checked) {
            setState(() {
              if (checked == true) {
                _selectedWallets.add(value);
              } else {
                _selectedWallets.remove(value);
              }
            });
          },
          activeColor: AppColors.primary,
        ),
        Text(label, style: const TextStyle(fontSize: 14)),
      ],
    );
  }

  Widget _buildWarningNotice() {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.warning.withOpacity(0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: const Row(
        children: [
          Icon(Icons.warning_amber, color: AppColors.warning, size: 20),
          SizedBox(width: 8),
          Expanded(
            child: Text(
              'APP仅支持划拨给直属下级',
              style: TextStyle(
                fontSize: 13,
                color: AppColors.warning,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBottomBar() {
    return Container(
      padding: EdgeInsets.only(
        left: AppSpacing.md,
        right: AppSpacing.md,
        top: 12,
        bottom: MediaQuery.of(context).padding.bottom + 12,
      ),
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: ElevatedButton(
        onPressed: _selectedAgentId != null ? _handleTransfer : null,
        child: const Text('确认划拨'),
      ),
    );
  }

  void _handleTransfer() {
    // TODO: 执行划拨操作
  }
}
