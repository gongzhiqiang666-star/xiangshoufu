import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../data/models/policy_model.dart';
import '../data/services/policy_service.dart';
import 'providers/policy_provider.dart';
import 'widgets/rate_editor_widget.dart';
import 'widgets/deposit_cashback_editor.dart';
import 'widgets/sim_cashback_editor.dart';
import 'widgets/activation_reward_editor.dart';

/// 下级政策调整页面（可编辑）
class SubordinatePolicyPage extends ConsumerStatefulWidget {
  final int subordinateId;
  final String subordinateName;
  final int? initialChannelId;

  const SubordinatePolicyPage({
    super.key,
    required this.subordinateId,
    required this.subordinateName,
    this.initialChannelId,
  });

  @override
  ConsumerState<SubordinatePolicyPage> createState() => _SubordinatePolicyPageState();
}

class _SubordinatePolicyPageState extends ConsumerState<SubordinatePolicyPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  int? _selectedChannelId;
  bool _isLoading = false;

  // 编辑状态
  RateConfig? _editedRates;
  List<DepositCashbackItem>? _editedDepositCashbacks;
  SimCashbackConfig? _editedSimCashback;
  List<ActivationRewardItem>? _editedActivationRewards;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
    _selectedChannelId = widget.initialChannelId;
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final channelsAsync = ref.watch(availableChannelsProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: Text('调整政策 - ${widget.subordinateName}'),
        bottom: TabBar(
          controller: _tabController,
          isScrollable: true,
          tabs: const [
            Tab(text: '成本费率'),
            Tab(text: '押金返现'),
            Tab(text: '流量返现'),
            Tab(text: '激活奖励'),
          ],
        ),
        actions: [
          TextButton.icon(
            onPressed: _isLoading ? null : _savePolicy,
            icon: _isLoading
                ? const SizedBox(
                    width: 16,
                    height: 16,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Icon(Icons.save),
            label: const Text('保存'),
            style: TextButton.styleFrom(foregroundColor: Colors.white),
          ),
        ],
      ),
      body: Column(
        children: [
          // 通道选择器
          channelsAsync.when(
            data: (channels) => _buildChannelSelector(channels),
            loading: () => const LinearProgressIndicator(),
            error: (e, _) => Padding(
              padding: const EdgeInsets.all(8),
              child: Text('加载通道失败: $e'),
            ),
          ),
          // Tab内容
          Expanded(
            child: _selectedChannelId == null
                ? const Center(child: Text('请选择通道'))
                : _buildTabContent(),
          ),
        ],
      ),
    );
  }

  Widget _buildChannelSelector(List<ChannelInfo> channels) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      color: Colors.white,
      child: DropdownButtonFormField<int>(
        value: _selectedChannelId,
        decoration: const InputDecoration(
          labelText: '选择通道',
          border: OutlineInputBorder(),
          contentPadding: EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        ),
        items: channels.map((c) {
          return DropdownMenuItem(
            value: c.id,
            child: Text(c.channelName),
          );
        }).toList(),
        onChanged: (value) {
          setState(() {
            _selectedChannelId = value;
            _resetEditState();
          });
        },
      ),
    );
  }

  void _resetEditState() {
    _editedRates = null;
    _editedDepositCashbacks = null;
    _editedSimCashback = null;
    _editedActivationRewards = null;
  }

  Widget _buildTabContent() {
    final policyParams = SubordinatePolicyParams(
      subordinateId: widget.subordinateId,
      channelId: _selectedChannelId!,
    );
    final policyAsync = ref.watch(subordinatePolicyProvider(policyParams));
    final limitsAsync = ref.watch(policyLimitsProvider(_selectedChannelId!));

    return policyAsync.when(
      data: (policy) => limitsAsync.when(
        data: (limits) => TabBarView(
          controller: _tabController,
          children: [
            // 成本费率
            RateEditorWidget(
              initialRates: _editedRates ?? policy.rateConfig,
              limits: limits,
              onChanged: (rates) => _editedRates = rates,
            ),
            // 押金返现
            DepositCashbackEditor(
              initialItems: _editedDepositCashbacks ?? policy.depositCashbacks ?? [],
              maxItems: limits.maxDepositCashbacks ?? [],
              onChanged: (items) => _editedDepositCashbacks = items,
            ),
            // 流量返现
            SimCashbackEditor(
              initialConfig: _editedSimCashback ?? policy.simCashback,
              maxConfig: limits.maxSimCashback,
              onChanged: (config) => _editedSimCashback = config,
            ),
            // 激活奖励
            ActivationRewardEditor(
              initialItems: _editedActivationRewards ?? policy.activationRewards ?? [],
              maxItems: limits.maxActivationRewards ?? [],
              onChanged: (items) => _editedActivationRewards = items,
            ),
          ],
        ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('加载限制失败: $e')),
      ),
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (e, _) => Center(child: Text('加载政策失败: $e')),
    );
  }

  Future<void> _savePolicy() async {
    if (_selectedChannelId == null) return;

    setState(() => _isLoading = true);

    try {
      final request = UpdateSubordinatePolicyRequest(
        channelId: _selectedChannelId!,
        creditRate: _editedRates?.creditRate,
        debitRate: _editedRates?.debitRate,
        debitCap: _editedRates?.debitCap,
        unionpayRate: _editedRates?.unionpayRate,
        wechatRate: _editedRates?.wechatRate,
        alipayRate: _editedRates?.alipayRate,
        depositCashbacks: _editedDepositCashbacks,
        simCashback: _editedSimCashback,
        activationRewards: _editedActivationRewards,
      );

      final service = ref.read(policyServiceProvider);
      await service.updateSubordinatePolicy(widget.subordinateId, request);

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('保存成功')),
        );
        // 刷新数据
        ref.invalidate(subordinatePolicyProvider(SubordinatePolicyParams(
          subordinateId: widget.subordinateId,
          channelId: _selectedChannelId!,
        )));
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('保存失败: $e')),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }
}
