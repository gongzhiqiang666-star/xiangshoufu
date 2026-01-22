import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_colors.dart';
import '../../agent/data/models/agent_model.dart';
import '../../agent/presentation/providers/agent_provider.dart';
import 'providers/terminal_provider.dart';

/// 终端划拨页面
class TerminalTransferPage extends ConsumerStatefulWidget {
  final List<String> selectedSNs;

  const TerminalTransferPage({
    super.key,
    required this.selectedSNs,
  });

  @override
  ConsumerState<TerminalTransferPage> createState() => _TerminalTransferPageState();
}

class _TerminalTransferPageState extends ConsumerState<TerminalTransferPage> {
  int? _selectedAgentId;
  String? _selectedAgentName;
  bool _enableCargoDeduction = false;
  double _unitPrice = 50.0;
  final Set<String> _selectedWallets = {'profit'};
  String _searchKeyword = '';
  bool _isSubmitting = false;

  @override
  void initState() {
    super.initState();
    // 加载下级代理商列表
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(subordinatesProvider.notifier).loadSubordinates(refresh: true);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('终端划拨')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 已选终端
            _buildSelectedTerminals(),
            const SizedBox(height: 16),

            // 划拨给
            _buildAgentSelector(),
            const SizedBox(height: 16),

            // 货款代扣设置
            _buildCargoDeductionSettings(),
            const SizedBox(height: 16),

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
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '已选终端: ${widget.selectedSNs.length}台',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'SN: ${widget.selectedSNs.join(", ")}',
            style: TextStyle(
              fontSize: 14,
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAgentSelector() {
    final subordinatesState = ref.watch(subordinatesProvider);

    // 根据搜索关键词过滤代理商列表
    final filteredAgents = _searchKeyword.isEmpty
        ? subordinatesState.list
        : subordinatesState.list.where((agent) {
            final keyword = _searchKeyword.toLowerCase();
            return agent.agentName.toLowerCase().contains(keyword) ||
                agent.agentNo.toLowerCase().contains(keyword) ||
                agent.contactPhone.contains(keyword);
          }).toList();

    return Container(
      padding: const EdgeInsets.all(16),
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
            ),
          ),
          const SizedBox(height: 12),
          // 搜索框
          TextField(
            decoration: InputDecoration(
              hintText: '搜索直属下级代理商',
              prefixIcon: Icon(Icons.search, color: AppColors.textTertiary),
              filled: true,
              fillColor: AppColors.background,
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(8),
                borderSide: BorderSide.none,
              ),
              contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            ),
            onChanged: (value) {
              setState(() {
                _searchKeyword = value;
              });
            },
          ),
          const SizedBox(height: 12),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                '直属下级代理商',
                style: TextStyle(
                  fontSize: 14,
                  color: AppColors.textSecondary,
                ),
              ),
              if (subordinatesState.isLoading)
                const SizedBox(
                  width: 16,
                  height: 16,
                  child: CircularProgressIndicator(strokeWidth: 2),
                ),
            ],
          ),
          const SizedBox(height: 8),
          if (subordinatesState.error != null)
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Colors.red.shade50,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Row(
                children: [
                  Icon(Icons.error_outline, color: Colors.red.shade700, size: 20),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      '加载失败: ${subordinatesState.error}',
                      style: TextStyle(color: Colors.red.shade700, fontSize: 13),
                    ),
                  ),
                  TextButton(
                    onPressed: () {
                      ref.read(subordinatesProvider.notifier).loadSubordinates(refresh: true);
                    },
                    child: const Text('重试'),
                  ),
                ],
              ),
            )
          else if (filteredAgents.isEmpty && !subordinatesState.isLoading)
            Container(
              padding: const EdgeInsets.all(24),
              child: Center(
                child: Text(
                  _searchKeyword.isEmpty ? '暂无直属下级代理商' : '未找到匹配的代理商',
                  style: TextStyle(color: AppColors.textTertiary),
                ),
              ),
            )
          else
            ...filteredAgents.map((agent) => _buildAgentItem(agent)),
        ],
      ),
    );
  }

  Widget _buildAgentItem(AgentInfo agent) {
    final isSelected = _selectedAgentId == agent.id;

    return GestureDetector(
      onTap: () {
        setState(() {
          _selectedAgentId = agent.id;
          _selectedAgentName = agent.agentName;
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
                    '${agent.agentName} (${agent.agentNo})',
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w500,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    '手机: ${_maskPhone(agent.contactPhone)}',
                    style: TextStyle(
                      fontSize: 13,
                      color: AppColors.textSecondary,
                    ),
                  ),
                ],
              ),
            ),
            if (isSelected)
              Text(
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

  /// 手机号脱敏显示
  String _maskPhone(String phone) {
    if (phone.length >= 11) {
      return '${phone.substring(0, 3)}****${phone.substring(7)}';
    }
    return phone;
  }

  Widget _buildCargoDeductionSettings() {
    return Container(
      padding: const EdgeInsets.all(16),
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
                Text(
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
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.primary,
            ),
          ),
          const SizedBox(height: 16),
          Text(
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
      child: Row(
        children: [
          Icon(Icons.warning_amber, color: AppColors.warning, size: 20),
          const SizedBox(width: 8),
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
        left: 16,
        right: 16,
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
        onPressed: (_selectedAgentId != null && !_isSubmitting) ? _handleTransfer : null,
        child: _isSubmitting
            ? const SizedBox(
                width: 20,
                height: 20,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                ),
              )
            : const Text('确认划拨'),
      ),
    );
  }

  void _handleTransfer() async {
    if (_selectedAgentId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请选择下级代理商')),
      );
      return;
    }

    // 确认对话框
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认划拨'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('确定要将 ${widget.selectedSNs.length} 台终端划拨给 $_selectedAgentName 吗？'),
            if (_enableCargoDeduction) ...[
              const SizedBox(height: 8),
              Text(
                '货款代扣: ¥${(_unitPrice * widget.selectedSNs.length).toStringAsFixed(2)}',
                style: TextStyle(
                  fontSize: 14,
                  color: AppColors.primary,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ],
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('确认'),
          ),
        ],
      ),
    );

    if (confirmed != true) return;

    setState(() => _isSubmitting = true);

    try {
      // 计算货款金额（分）
      final goodsPrice = _enableCargoDeduction ? (_unitPrice * 100).toInt() : 0;
      // 代扣类型: 1=一次性, 2=分期, 3=货款代扣
      final deductionType = _enableCargoDeduction ? 3 : 1;

      final successCount = await ref.read(terminalDistributeProvider.notifier).batchDistribute(
        toAgentId: _selectedAgentId!,
        terminalSns: widget.selectedSNs,
        channelId: 1, // 默认通道ID，后续可从终端信息获取
        goodsPrice: goodsPrice,
        deductionType: deductionType,
      );

      if (mounted) {
        if (successCount > 0) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('划拨成功，共${successCount}台终端'),
              backgroundColor: AppColors.success,
            ),
          );
          // 清空选中状态
          ref.read(selectedTerminalsProvider.notifier).state = [];
          // 返回上一页
          context.pop(true);
        } else {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('划拨失败，请稍后重试'),
              backgroundColor: Colors.red,
            ),
          );
        }
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('划拨失败: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isSubmitting = false);
      }
    }
  }
}
