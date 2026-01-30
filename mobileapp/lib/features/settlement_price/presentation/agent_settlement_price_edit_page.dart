import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../channel/channel.dart';
import '../data/models/settlement_price_model.dart';
import 'providers/settlement_price_provider.dart';

/// 结算价编辑页面
/// 入口：下级结算价列表 → 点击编辑
/// 支持编辑费率、押金返现、流量费返现
/// 添加通道配置校验：费率范围、返现上限
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

  // 动态费率控制器 Map: rateCode -> controller
  final Map<String, TextEditingController> _rateControllers = {};

  // 动态流量费返现控制器 Map: tierOrder -> controller
  final Map<int, TextEditingController> _simControllers = {};

  // 押金返现列表
  List<DepositCashbackItem> _depositCashbacks = [];

  // 当前结算价数据
  SettlementPriceModel? _currentPrice;

  // 通道配置（用于校验）
  ChannelFullConfig? _channelConfig;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    for (final controller in _rateControllers.values) {
      controller.dispose();
    }
    for (final controller in _simControllers.values) {
      controller.dispose();
    }
    super.dispose();
  }

  /// 初始化表单数据
  void _initFormData(SettlementPriceModel price, ChannelFullConfig? config) {
    if (_currentPrice?.id == price.id) return;
    _currentPrice = price;

    // 清理旧控制器
    for (final controller in _rateControllers.values) {
      controller.dispose();
    }
    _rateControllers.clear();

    for (final controller in _simControllers.values) {
      controller.dispose();
    }
    _simControllers.clear();

    // 初始化费率控制器（从通道配置获取费率类型）
    if (config != null) {
      for (final rateConfig in config.rateConfigs) {
        final rateCode = rateConfig.rateCode;
        final existingRate = price.rateConfigs[rateCode]?.rate ?? '';
        _rateControllers[rateCode] = TextEditingController(text: existingRate);
      }
    }

    // 初始化流量费返现控制器（从通道配置获取档位）
    if (config != null) {
      for (final tier in config.simCashbackTiers) {
        final existingTier = price.simCashbackTiers.firstWhere(
          (t) => t.tierOrder == tier.tierOrder,
          orElse: () => SimCashbackTier(tierOrder: tier.tierOrder, tierName: tier.tierName, cashbackAmount: 0),
        );
        _simControllers[tier.tierOrder] = TextEditingController(
          text: existingTier.cashbackAmountYuan.toStringAsFixed(0),
        );
      }
    }

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
          // 加载通道配置
          final configAsync = ref.watch(channelFullConfigProvider(price.channelId));

          return configAsync.when(
            loading: () => const Center(child: CircularProgressIndicator()),
            error: (e, s) {
              // 通道配置加载失败时，仍允许编辑但不显示限制提示
              _channelConfig = null;
              _initFormData(price, null);
              return TabBarView(
                controller: _tabController,
                children: [
                  _buildRateTab(price),
                  _buildDepositTab(price),
                  _buildSimTab(price),
                ],
              );
            },
            data: (config) {
              _channelConfig = config;
              _initFormData(price, config);
              return TabBarView(
                controller: _tabController,
                children: [
                  _buildRateTab(price),
                  _buildDepositTab(price),
                  _buildSimTab(price),
                ],
              );
            },
          );
        },
      ),
    );
  }

  /// 费率设置Tab - 动态渲染
  Widget _buildRateTab(SettlementPriceModel price) {
    final rateConfigs = _channelConfig?.rateConfigs ?? [];

    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildSectionCard(
            title: '费率配置',
            subtitle: '设置该下级代理商的交易费率',
            children: rateConfigs.isEmpty
                ? [
                    Container(
                      padding: const EdgeInsets.all(AppSpacing.lg),
                      child: Center(
                        child: Text('该通道未配置费率类型', style: TextStyle(color: Colors.grey[600])),
                      ),
                    ),
                  ]
                : rateConfigs.map((rateConfig) {
                    return Padding(
                      padding: const EdgeInsets.only(bottom: AppSpacing.md),
                      child: _buildRateField(
                        label: rateConfig.rateName,
                        controller: _rateControllers[rateConfig.rateCode] ?? TextEditingController(),
                        suffix: '%',
                        hint: '范围: ${rateConfig.minRate}% ~ ${rateConfig.maxRate}%',
                        rateConfig: rateConfig,
                      ),
                    );
                  }).toList(),
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

  /// 流量费返现设置Tab - 动态渲染
  Widget _buildSimTab(SettlementPriceModel price) {
    final simTiers = _channelConfig?.simCashbackTiers ?? [];

    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildSectionCard(
            title: '流量费返现配置',
            subtitle: '设置流量费返现金额（单位：元）',
            children: simTiers.isEmpty
                ? [
                    Container(
                      padding: const EdgeInsets.all(AppSpacing.lg),
                      child: Center(
                        child: Text('该通道未配置流量费档位', style: TextStyle(color: Colors.grey[600])),
                      ),
                    ),
                  ]
                : simTiers.map((tier) {
                    return Padding(
                      padding: const EdgeInsets.only(bottom: AppSpacing.md),
                      child: _buildSimCashbackField(
                        label: tier.tierName,
                        controller: _simControllers[tier.tierOrder] ?? TextEditingController(),
                        tier: tier,
                      ),
                    );
                  }).toList(),
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

  /// 构建费率输入字段（带通道配置范围提示）
  Widget _buildRateField({
    required String label,
    required TextEditingController controller,
    required String suffix,
    String? hint,
    ChannelRateConfig? rateConfig,
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
          keyboardType: const TextInputType.numberWithOptions(decimal: true),
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
        // 显示费率范围提示
        if (rateConfig != null) ...[
          const SizedBox(height: 4),
          Text(
            '范围: ${rateConfig.minRate}% ~ ${rateConfig.maxRate}%',
            style: TextStyle(
              fontSize: 12,
              color: Colors.grey[600],
            ),
          ),
        ],
      ],
    );
  }

  /// 构建流量费返现输入字段（带通道配置上限提示）
  Widget _buildSimCashbackField({
    required String label,
    required TextEditingController controller,
    required ChannelSimCashbackTier tier,
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
          keyboardType: TextInputType.number,
          inputFormatters: [
            FilteringTextInputFormatter.digitsOnly,
          ],
          decoration: InputDecoration(
            hintText: '例如: 49',
            suffixText: '元',
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
        // 显示返现上限提示
        if (tier.maxCashbackAmount > 0) ...[
          const SizedBox(height: 4),
          Text(
            '最高 ${tier.maxCashbackAmountYuan.toStringAsFixed(0)} 元',
            style: TextStyle(
              fontSize: 12,
              color: Colors.grey[600],
            ),
          ),
        ],
      ],
    );
  }

  /// 构建押金返现项卡片
  Widget _buildDepositItemCard(int index, DepositCashbackItem item) {
    // 获取通道押金档位上限
    ChannelDepositTier? tier;
    if (_channelConfig != null) {
      tier = _channelConfig!.getDepositTierByAmount(item.depositAmount);
    }

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
                // 显示返现上限提示
                if (tier != null && tier.maxCashbackAmount > 0) ...[
                  const SizedBox(height: 2),
                  Text(
                    '最高可返 ${tier.maxCashbackAmountYuan.toStringAsFixed(0)} 元',
                    style: TextStyle(
                      fontSize: 11,
                      color: Colors.grey[500],
                    ),
                  ),
                ],
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

    // 获取通道押金档位上限（用于对话框提示）
    int? maxCashback;
    if (_channelConfig != null && item != null) {
      final tier = _channelConfig!.getDepositTierByAmount(item.depositAmount);
      maxCashback = tier?.maxCashbackAmount;
    }

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
              onChanged: (value) {
                // 当押金金额变化时，更新返现上限提示
                final deposit = int.tryParse(value) ?? 0;
                if (_channelConfig != null && deposit > 0) {
                  final tier = _channelConfig!.getDepositTierByAmount(deposit * 100);
                  maxCashback = tier?.maxCashbackAmount;
                }
              },
            ),
            const SizedBox(height: 16),
            TextField(
              controller: cashbackController,
              keyboardType: TextInputType.number,
              decoration: InputDecoration(
                labelText: '返现金额（元）',
                hintText: '例如: 69',
                helperText: maxCashback != null && maxCashback! > 0
                    ? '最高 ${(maxCashback! / 100).toStringAsFixed(0)} 元'
                    : null,
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

              // 验证返现金额是否超过通道上限
              if (_channelConfig != null) {
                final tier = _channelConfig!.getDepositTierByAmount(deposit * 100);
                if (tier != null && cashback * 100 > tier.maxCashbackAmount) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text('返现金额不能超过${tier.maxCashbackAmountYuan.toStringAsFixed(0)}元'),
                      backgroundColor: AppColors.danger,
                    ),
                  );
                  return;
                }
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

  /// 显示错误提示
  void _showError(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message), backgroundColor: AppColors.danger),
    );
  }

  /// 验证费率是否在通道允许范围内
  bool _validateRates() {
    if (_channelConfig == null) return true;

    for (final rateConfig in _channelConfig!.rateConfigs) {
      final controller = _rateControllers[rateConfig.rateCode];
      if (controller == null || controller.text.isEmpty) continue;

      final rate = double.tryParse(controller.text) ?? 0;
      final minRate = rateConfig.minRateValue;
      final maxRate = rateConfig.maxRateValue;

      if (rate < minRate || rate > maxRate) {
        _showError('${rateConfig.rateName}必须在 $minRate% ~ $maxRate% 范围内');
        return false;
      }
    }
    return true;
  }

  /// 验证押金返现是否在通道允许范围内
  bool _validateDepositCashbacks() {
    if (_channelConfig == null) return true;

    for (final item in _depositCashbacks) {
      final tier = _channelConfig!.getDepositTierByAmount(item.depositAmount);
      if (tier != null && item.cashbackAmount > tier.maxCashbackAmount) {
        _showError('押金${item.depositAmountYuan.toStringAsFixed(0)}元的返现不能超过${tier.maxCashbackAmountYuan.toStringAsFixed(0)}元');
        return false;
      }
    }
    return true;
  }

  /// 验证流量费返现是否在通道允许范围内
  bool _validateSimCashbacks() {
    if (_channelConfig == null) return true;

    for (final tier in _channelConfig!.simCashbackTiers) {
      final controller = _simControllers[tier.tierOrder];
      if (controller == null || controller.text.isEmpty) continue;

      final cashback = int.tryParse(controller.text) ?? 0;
      final maxCashback = tier.maxCashbackAmount;

      if (cashback * 100 > maxCashback) {
        _showError('${tier.tierName}不能超过${tier.maxCashbackAmountYuan.toStringAsFixed(0)}元');
        return false;
      }
    }
    return true;
  }

  /// 保存所有修改
  Future<void> _saveAll() async {
    // 1. 验证费率范围
    if (!_validateRates()) return;

    // 2. 验证押金返现上限
    if (!_validateDepositCashbacks()) return;

    // 3. 验证流量费返现上限
    if (!_validateSimCashbacks()) return;

    setState(() => _isLoading = true);

    try {
      final service = ref.read(settlementPriceServiceProvider);

      // 保存费率 - 使用动态 rate_configs
      final rateConfigs = <String, Map<String, String>>{};
      for (final entry in _rateControllers.entries) {
        if (entry.value.text.isNotEmpty) {
          rateConfigs[entry.key] = {'rate': entry.value.text};
        }
      }
      await service.updateRate(widget.priceId, {
        'rate_configs': rateConfigs,
      });

      // 保存押金返现
      await service.updateDepositCashback(widget.priceId, {
        'deposit_cashbacks': _depositCashbacks.map((e) => e.toJson()).toList(),
      });

      // 保存流量费返现 - 使用动态 sim_cashback_tiers
      final simCashbackTiers = <Map<String, dynamic>>[];
      for (final tier in _channelConfig?.simCashbackTiers ?? []) {
        final controller = _simControllers[tier.tierOrder];
        final cashback = int.tryParse(controller?.text ?? '0') ?? 0;
        simCashbackTiers.add({
          'tier_order': tier.tierOrder,
          'tier_name': tier.tierName,
          'cashback_amount': cashback * 100,
        });
      }
      await service.updateSimCashback(widget.priceId, {
        'sim_cashback_tiers': simCashbackTiers,
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
