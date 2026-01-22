import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../core/utils/format_utils.dart';
import '../data/models/wallet_model.dart';
import 'providers/wallet_provider.dart';
import '../data/services/wallet_service.dart';

/// 下级代理商选择模型
class SubordinateAgentModel {
  final int id;
  final String name;
  final String phone;

  SubordinateAgentModel({
    required this.id,
    required this.name,
    required this.phone,
  });
}

/// 奖励发放页面
class IssueRewardPage extends ConsumerStatefulWidget {
  const IssueRewardPage({super.key});

  @override
  ConsumerState<IssueRewardPage> createState() => _IssueRewardPageState();
}

class _IssueRewardPageState extends ConsumerState<IssueRewardPage> {
  final _formKey = GlobalKey<FormState>();
  final _agentIdController = TextEditingController();
  final _amountController = TextEditingController();
  final _remarkController = TextEditingController();

  SubordinateAgentModel? _selectedAgent;
  bool _isLoading = false;

  @override
  void dispose() {
    _agentIdController.dispose();
    _amountController.dispose();
    _remarkController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final summaryAsync = ref.watch(chargingWalletSummaryProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('发放奖励'),
      ),
      body: summaryAsync.when(
        data: (summary) => _buildContent(summary),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.error_outline, size: 48, color: Colors.red),
              const SizedBox(height: 16),
              Text('加载失败: $e'),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () => ref.invalidate(chargingWalletSummaryProvider),
                child: const Text('重试'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildContent(ChargingWalletSummaryModel summary) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Form(
        key: _formKey,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 可用余额卡片
            _buildBalanceCard(summary),
            const SizedBox(height: 24),

            // 发放说明
            _buildInstructionsCard(),
            const SizedBox(height: 24),

            // 表单区域
            _buildFormSection(summary),
          ],
        ),
      ),
    );
  }

  Widget _buildBalanceCard(ChargingWalletSummaryModel summary) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.lg),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF11998e), Color(0xFF38ef7d)],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: const Color(0xFF11998e).withOpacity(0.3),
            blurRadius: 20,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  '充值钱包可用余额',
                  style: TextStyle(
                    color: Colors.white70,
                    fontSize: 14,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  FormatUtils.formatYuan(summary.balanceYuan),
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 28,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.2),
              borderRadius: BorderRadius.circular(12),
            ),
            child: const Icon(
              Icons.account_balance_wallet,
              color: Colors.white,
              size: 32,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildInstructionsCard() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.blue.shade50,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: Colors.blue.shade200),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(Icons.info_outline, color: Colors.blue.shade700, size: 20),
              const SizedBox(width: 8),
              Text(
                '奖励发放规则',
                style: TextStyle(
                  color: Colors.blue.shade700,
                  fontWeight: FontWeight.bold,
                  fontSize: 15,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          _buildInstructionItem('1', '只能给直属下级代理商发放奖励'),
          const SizedBox(height: 6),
          _buildInstructionItem('2', '发放金额将从您的充值钱包扣除'),
          const SizedBox(height: 6),
          _buildInstructionItem('3', '奖励将实时到账下级的奖励钱包'),
          const SizedBox(height: 6),
          _buildInstructionItem('4', '发放后不可撤销，请仔细核对'),
        ],
      ),
    );
  }

  Widget _buildInstructionItem(String number, String text) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Container(
          width: 20,
          height: 20,
          decoration: BoxDecoration(
            color: Colors.blue.shade700,
            shape: BoxShape.circle,
          ),
          child: Center(
            child: Text(
              number,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 12,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: Text(
            text,
            style: TextStyle(
              color: Colors.blue.shade700,
              fontSize: 13,
              height: 1.4,
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildFormSection(ChargingWalletSummaryModel summary) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '发放信息',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 20),

          // 选择下级代理商
          _buildAgentSelector(),
          const SizedBox(height: 16),

          // 发放金额
          TextFormField(
            controller: _amountController,
            keyboardType: const TextInputType.numberWithOptions(decimal: true),
            decoration: InputDecoration(
              labelText: '发放金额',
              hintText: '请输入发放金额',
              prefixText: '¥ ',
              prefixIcon: const Icon(Icons.monetization_on),
              border: const OutlineInputBorder(),
              suffixText: '元',
              helperText: '最大可发放: ${FormatUtils.formatYuan(summary.balanceYuan)}',
            ),
            validator: (value) {
              if (value == null || value.trim().isEmpty) {
                return '请输入发放金额';
              }
              final amount = double.tryParse(value);
              if (amount == null || amount <= 0) {
                return '请输入有效金额';
              }
              if (amount * 100 > summary.balance) {
                return '发放金额超过可用余额';
              }
              return null;
            },
          ),
          const SizedBox(height: 16),

          // 备注
          TextFormField(
            controller: _remarkController,
            decoration: const InputDecoration(
              labelText: '备注',
              hintText: '可选，填写发放原因',
              prefixIcon: Icon(Icons.note),
              border: OutlineInputBorder(),
            ),
            maxLines: 2,
            maxLength: 100,
          ),
          const SizedBox(height: 24),

          // 确认发放按钮
          SizedBox(
            width: double.infinity,
            height: 50,
            child: ElevatedButton(
              onPressed: _isLoading || summary.balance <= 0
                  ? null
                  : () => _submitReward(summary),
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.primary,
                foregroundColor: Colors.white,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              child: _isLoading
                  ? const SizedBox(
                      width: 24,
                      height: 24,
                      child: CircularProgressIndicator(
                        color: Colors.white,
                        strokeWidth: 2,
                      ),
                    )
                  : const Text(
                      '确认发放',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAgentSelector() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        TextFormField(
          controller: _agentIdController,
          keyboardType: TextInputType.number,
          decoration: InputDecoration(
            labelText: '下级代理商',
            hintText: '请输入下级代理商ID',
            prefixIcon: const Icon(Icons.person),
            border: const OutlineInputBorder(),
            suffixIcon: IconButton(
              icon: const Icon(Icons.search),
              onPressed: _searchAgent,
              tooltip: '查询代理商',
            ),
          ),
          validator: (value) {
            if (value == null || value.trim().isEmpty) {
              return '请输入下级代理商ID';
            }
            if (int.tryParse(value) == null) {
              return '请输入有效的代理商ID';
            }
            return null;
          },
        ),
        if (_selectedAgent != null) ...[
          const SizedBox(height: 8),
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.green.shade50,
              borderRadius: BorderRadius.circular(8),
              border: Border.all(color: Colors.green.shade200),
            ),
            child: Row(
              children: [
                Icon(Icons.check_circle, color: Colors.green.shade700, size: 20),
                const SizedBox(width: 8),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        _selectedAgent!.name,
                        style: TextStyle(
                          color: Colors.green.shade700,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      Text(
                        'ID: ${_selectedAgent!.id} | ${_selectedAgent!.phone}',
                        style: TextStyle(
                          color: Colors.green.shade600,
                          fontSize: 12,
                        ),
                      ),
                    ],
                  ),
                ),
                IconButton(
                  icon: Icon(Icons.close, color: Colors.green.shade700, size: 20),
                  onPressed: () {
                    setState(() {
                      _selectedAgent = null;
                      _agentIdController.clear();
                    });
                  },
                ),
              ],
            ),
          ),
        ],
      ],
    );
  }

  void _searchAgent() {
    final agentIdText = _agentIdController.text.trim();
    if (agentIdText.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入代理商ID')),
      );
      return;
    }

    final agentId = int.tryParse(agentIdText);
    if (agentId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入有效的代理商ID')),
      );
      return;
    }

    // TODO: 实际应调用API查询代理商信息并验证是否为直属下级
    // 这里暂时模拟
    setState(() {
      _selectedAgent = SubordinateAgentModel(
        id: agentId,
        name: '代理商$agentId',
        phone: '138****${agentId.toString().padLeft(4, '0').substring(0, 4)}',
      );
    });

    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('代理商信息已确认')),
    );
  }

  Future<void> _submitReward(ChargingWalletSummaryModel summary) async {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    if (_selectedAgent == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请先查询并确认下级代理商信息')),
      );
      return;
    }

    final amount = double.parse(_amountController.text.trim());

    // 确认对话框
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认发放'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('接收方: ${_selectedAgent!.name}'),
            const SizedBox(height: 8),
            Text('发放金额: ${FormatUtils.formatYuan(amount)}'),
            if (_remarkController.text.trim().isNotEmpty) ...[
              const SizedBox(height: 8),
              Text('备注: ${_remarkController.text.trim()}'),
            ],
            const SizedBox(height: 16),
            Container(
              padding: const EdgeInsets.all(8),
              decoration: BoxDecoration(
                color: Colors.orange.shade50,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Row(
                children: [
                  Icon(Icons.warning_amber, color: Colors.orange.shade700, size: 20),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      '发放后不可撤销',
                      style: TextStyle(
                        color: Colors.orange.shade700,
                        fontSize: 13,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context, true),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.primary,
              foregroundColor: Colors.white,
            ),
            child: const Text('确认发放'),
          ),
        ],
      ),
    );

    if (confirmed != true) return;

    setState(() => _isLoading = true);

    try {
      final walletService = ref.read(walletServiceProvider);
      final rewardNo = await walletService.issueChargingReward(
        toAgentId: _selectedAgent!.id,
        amount: (amount * 100).round(),
        remark: _remarkController.text.trim().isEmpty
            ? null
            : _remarkController.text.trim(),
      );

      if (mounted) {
        ref.invalidate(chargingWalletSummaryProvider);

        // 显示成功对话框
        await showDialog(
          context: context,
          barrierDismissible: false,
          builder: (context) => AlertDialog(
            title: Row(
              children: const [
                Icon(Icons.check_circle, color: Colors.green, size: 28),
                SizedBox(width: 8),
                Text('发放成功'),
              ],
            ),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('奖励已发放给 ${_selectedAgent!.name}'),
                const SizedBox(height: 8),
                Text('发放金额: ${FormatUtils.formatYuan(amount)}'),
                const SizedBox(height: 8),
                Text(
                  '单号: $rewardNo',
                  style: const TextStyle(
                    color: AppColors.textSecondary,
                    fontSize: 12,
                  ),
                ),
              ],
            ),
            actions: [
              ElevatedButton(
                onPressed: () {
                  Navigator.pop(context);
                  Navigator.pop(context); // 返回上一页
                },
                child: const Text('完成'),
              ),
            ],
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('发放失败: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }
}
