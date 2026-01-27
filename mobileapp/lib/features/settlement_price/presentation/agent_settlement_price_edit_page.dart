import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../data/models/settlement_price_model.dart';
import 'providers/settlement_price_provider.dart';

/// 结算价编辑页面
/// 入口：下级结算价列表 → 点击编辑
/// 支持编辑费率、押金返现、流量费返现
class AgentSettlementPriceEditPage extends ConsumerStatefulWidget {
  final int agentId;
  final int priceId;
  final String? agentName;
  final String? channelName;

  const AgentSettlementPriceEditPage({
    super.key,
    required this.agentId,
    required this.priceId,
    this.agentName,
    this.channelName,
  });

  @override
  ConsumerState<AgentSettlementPriceEditPage> createState() => _AgentSettlementPriceEditPageState();
}

class _AgentSettlementPriceEditPageState extends ConsumerState<AgentSettlementPriceEditPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  bool _isLoading = false;

  // 费率控制器
  final _creditRateController = TextEditingController();
  final _debitRateController = TextEditingController();
  final _debitCapController = TextEditingController();
  final _unionpayRateController = TextEditingController();
  final _wechatRateController = TextEditingController();
  final _alipayRateController = TextEditingController();

  // 流量费返现控制器
  final _simFirstController = TextEditingController();
  final _simSecondController = TextEditingController();
  final _simThirdPlusController = TextEditingController();

  // 押金返现列表
  List<DepositCashbackItem> _depositCashbacks = [];

  // 当前结算价数据
  SettlementPriceModel? _currentPrice;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    _creditRateController.dispose();
    _debitRateController.dispose();
    _debitCapController.dispose();
    _unionpayRateController.dispose();
    _wechatRateController.dispose();
    _alipayRateController.dispose();
    _simFirstController.dispose();
    _simSecondController.dispose();
    _simThirdPlusController.dispose();
    super.dispose();
  }

  /// 初始化表单数据
  void _initFormData(SettlementPriceModel price) {
    if (_currentPrice?.id == price.id) return;
    _currentPrice = price;

    _creditRateController.text = price.creditRate ?? '';
    _debitRateController.text = price.debitRate ?? '';
    _debitCapController.text = price.debitCap ?? '';
    _unionpayRateController.text = price.unionpayRate ?? '';
    _wechatRateController.text = price.wechatRate ?? '';
    _alipayRateController.text = price.alipayRate ?? '';

    _simFirstController.text = price.simFirstCashbackYuan.toStringAsFixed(0);
    _simSecondController.text = price.simSecondCashbackYuan.toStringAsFixed(0);
    _simThirdPlusController.text = price.simThirdPlusCashbackYuan.toStringAsFixed(0);

    _depositCashbacks = List.from(price.depositCashbacks);
  }

  @override
  Widget build(BuildContext context) {
    final priceAsync = ref.watch(settlementPriceDetailProvider(widget.priceId));

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: Text(widget.channelName != null ? '编辑${widget.channelName}结算价' : '编辑结算价'),
        actions: [
          TextButton(
            onPressed: _isLoading ? null : _saveAll,
            child: _isLoading
                ? const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                  )
                : const Text('保存', style: TextStyle(color: Colors.white, fontWeight: FontWeight.w600)),
          ),
        ],
        bottom: TabBar(
          controller: _tabController,
          labelColor: Colors.white,
          unselectedLabelColor: Colors.white70,
          indicatorColor: Colors.white,
          tabs: const [
            Tab(text: '费率'),
            Tab(text: '押金返现'),
            Tab(text: '流量返现'),
          ],
        ),
      ),
      body: priceAsync.when(
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
                onPressed: () => ref.invalidate(settlementPriceDetailProvider(widget.priceId)),
                child: const Text('重试'),
              ),
            ],
          ),
        ),
        data: (price) {
          _initFormData(price);
          return TabBarView(
            controller: _tabController,
            children: [
              _buildRateTab(price),
              _buildDepositTab(price),
              _buildSimTab(price),
            ],
          );
        },
      ),
    );
  }

  /// 费率设置Tab
  Widget _buildRateTab(SettlementPriceModel price) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildSectionCard(
            title: '费率配置',
            subtitle: '设置该下级代理商的交易费率',
            children: [
              _buildRateField(
                label: '贷记卡费率',
                controller: _creditRateController,
                suffix: '%',
                hint: '例如: 0.55',
              ),
              const SizedBox(height: AppSpacing.md),
              _buildRateField(
                label: '借记卡费率',
                controller: _debitRateController,
                suffix: '%',
                hint: '例如: 0.50',
              ),
              const SizedBox(height: AppSpacing.md),
              _buildRateField(
                label: '借记卡封顶',
                controller: _debitCapController,
                suffix: '元',
                hint: '例如: 20',
              ),
              const SizedBox(height: AppSpacing.md),
              _buildRateField(
                label: '云闪付费率',
                controller: _unionpayRateController,
                suffix: '%',
                hint: '例如: 0.38',
              ),
              const SizedBox(height: AppSpacing.md),
              _buildRateField(
                label: '微信费率',
                controller: _wechatRateController,
                suffix: '%',
                hint: '例如: 0.38',
              ),
              const SizedBox(height: AppSpacing.md),
              _buildRateField(
                label: '支付宝费率',
                controller: _alipayRateController,
                suffix: '%',
                hint: '例如: 0.38',
              ),
            ],
          ),
        ],
      ),
    );
  }

  /// 押金返现设置Tab
  Widget _buildDepositTab(SettlementPriceModel price) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildSectionCard(
            title: '押金返现配置',
            subtitle: '设置不同押金档位的返现金额',
            trailing: IconButton(
              icon: const Icon(Icons.add_circle_outline, color: AppColors.primary),
              onPressed: _addDepositItem,
            ),
            children: [
              if (_depositCashbacks.isEmpty)
                Container(
                  padding: const EdgeInsets.all(AppSpacing.lg),
                  child: Center(
                    child: Column(
                      children: [
                        Icon(Icons.inbox_outlined, size: 48, color: Colors.grey[400]),
                        const SizedBox(height: 8),
                        Text('暂无押金返现配置', style: TextStyle(color: Colors.grey[600])),
                        const SizedBox(height: 8),
                        TextButton.icon(
                          onPressed: _addDepositItem,
                          icon: const Icon(Icons.add),
                          label: const Text('添加配置'),
                        ),
                      ],
                    ),
                  ),
                )
              else
                ..._depositCashbacks.asMap().entries.map((entry) {
                  final index = entry.key;
                  final item = entry.value;
                  return _buildDepositItemCard(index, item);
                }),
            ],
          ),
        ],
      ),
    );
  }

  /// 流量费返现设置Tab
  Widget _buildSimTab(SettlementPriceModel price) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildSectionCard(
            title: '流量费返现配置',
            subtitle: '设置流量费返现金额（单位：元）',
            children: [
              _buildRateField(
                label: '首次返现',
                controller: _simFirstController,
                suffix: '元',
                hint: '例如: 49',
                keyboardType: TextInputType.number,
              ),
              const SizedBox(height: AppSpacing.md),
              _buildRateField(
                label: '第2次返现',
                controller: _simSecondController,
                suffix: '元',
                hint: '例如: 39',
                keyboardType: TextInputType.number,
              ),
              const SizedBox(height: AppSpacing.md),
              _buildRateField(
                label: '第3次及以后返现',
                controller: _simThirdPlusController,
                suffix: '元',
                hint: '例如: 29',
                keyboardType: TextInputType.number,
              ),
            ],
          ),
        ],
      ),
    );
  }

  /// 构建区域卡片
  Widget _buildSectionCard({
    required String title,
    String? subtitle,
    Widget? trailing,
    required List<Widget> children,
  }) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.all(AppSpacing.md),
            child: Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        title,
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: AppColors.textPrimary,
                        ),
                      ),
                      if (subtitle != null) ...[
                        const SizedBox(height: 4),
                        Text(
                          subtitle,
                          style: const TextStyle(
                            fontSize: 12,
                            color: AppColors.textSecondary,
                          ),
                        ),
                      ],
                    ],
                  ),
                ),
                if (trailing != null) trailing,
              ],
            ),
          ),
          const Divider(height: 1, color: AppColors.divider),
          Padding(
            padding: const EdgeInsets.all(AppSpacing.md),
            child: Column(children: children),
          ),
        ],
      ),
    );
  }

  /// 构建费率输入字段
  Widget _buildRateField({
    required String label,
    required TextEditingController controller,
    required String suffix,
    String? hint,
    TextInputType keyboardType = const TextInputType.numberWithOptions(decimal: true),
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w500,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: 8),
        TextField(
          controller: controller,
          keyboardType: keyboardType,
          inputFormatters: [
            FilteringTextInputFormatter.allow(RegExp(r'[\d.]')),
          ],
          decoration: InputDecoration(
            hintText: hint,
            suffixText: suffix,
            filled: true,
            fillColor: AppColors.background,
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(8),
              borderSide: BorderSide.none,
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(8),
              borderSide: const BorderSide(color: AppColors.primary),
            ),
            contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          ),
        ),
      ],
    );
  }

  /// 构建押金返现项卡片
  Widget _buildDepositItemCard(int index, DepositCashbackItem item) {
    return Container(
      margin: const EdgeInsets.only(bottom: AppSpacing.sm),
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: AppColors.background,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  '押金档位 ${index + 1}',
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.textSecondary,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  '¥${item.depositAmountYuan.toStringAsFixed(0)} → 返¥${item.cashbackAmountYuan.toStringAsFixed(0)}',
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: AppColors.textPrimary,
                  ),
                ),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(Icons.edit_outlined, color: AppColors.primary, size: 20),
            onPressed: () => _editDepositItem(index, item),
          ),
          IconButton(
            icon: const Icon(Icons.delete_outline, color: AppColors.danger, size: 20),
            onPressed: () => _removeDepositItem(index),
          ),
        ],
      ),
    );
  }

  /// 添加押金返现项
  void _addDepositItem() {
    _showDepositDialog(null, null);
  }

  /// 编辑押金返现项
  void _editDepositItem(int index, DepositCashbackItem item) {
    _showDepositDialog(index, item);
  }

  /// 删除押金返现项
  void _removeDepositItem(int index) {
    setState(() {
      _depositCashbacks.removeAt(index);
    });
  }

  /// 显示押金返现编辑对话框
  void _showDepositDialog(int? index, DepositCashbackItem? item) {
    final depositController = TextEditingController(
      text: item?.depositAmountYuan.toStringAsFixed(0) ?? '',
    );
    final cashbackController = TextEditingController(
      text: item?.cashbackAmountYuan.toStringAsFixed(0) ?? '',
    );

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(index == null ? '添加押金返现' : '编辑押金返现'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: depositController,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(
                labelText: '押金金额（元）',
                hintText: '例如: 99',
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: cashbackController,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(
                labelText: '返现金额（元）',
                hintText: '例如: 69',
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () {
              final deposit = int.tryParse(depositController.text) ?? 0;
              final cashback = int.tryParse(cashbackController.text) ?? 0;

              if (deposit <= 0 || cashback <= 0) {
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('请输入有效的金额')),
                );
                return;
              }

              final newItem = DepositCashbackItem(
                depositAmount: deposit * 100,
                cashbackAmount: cashback * 100,
              );

              setState(() {
                if (index == null) {
                  _depositCashbacks.add(newItem);
                } else {
                  _depositCashbacks[index] = newItem;
                }
              });

              Navigator.pop(context);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }

  /// 保存所有修改
  Future<void> _saveAll() async {
    setState(() => _isLoading = true);

    try {
      final service = ref.read(settlementPriceServiceProvider);

      // 保存费率
      await service.updateRate(widget.priceId, {
        'credit_rate': _creditRateController.text,
        'debit_rate': _debitRateController.text,
        'debit_cap': _debitCapController.text,
        'unionpay_rate': _unionpayRateController.text,
        'wechat_rate': _wechatRateController.text,
        'alipay_rate': _alipayRateController.text,
      });

      // 保存押金返现
      await service.updateDepositCashback(widget.priceId, {
        'deposit_cashbacks': _depositCashbacks.map((e) => e.toJson()).toList(),
      });

      // 保存流量费返现
      final simFirst = int.tryParse(_simFirstController.text) ?? 0;
      final simSecond = int.tryParse(_simSecondController.text) ?? 0;
      final simThirdPlus = int.tryParse(_simThirdPlusController.text) ?? 0;
      await service.updateSimCashback(widget.priceId, {
        'sim_first_cashback': simFirst * 100,
        'sim_second_cashback': simSecond * 100,
        'sim_third_plus_cashback': simThirdPlus * 100,
      });

      // 刷新列表
      ref.invalidate(agentSettlementPriceListProvider(widget.agentId));
      ref.invalidate(settlementPriceDetailProvider(widget.priceId));

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('保存成功'), backgroundColor: AppColors.success),
        );
        context.pop();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('保存失败: $e'), backgroundColor: AppColors.danger),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }
}
