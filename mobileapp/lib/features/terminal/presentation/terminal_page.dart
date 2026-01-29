import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_colors.dart';
import '../../../router/app_router.dart';
import '../../agent/presentation/providers/agent_provider.dart';
import '../domain/models/terminal.dart';
import 'providers/terminal_provider.dart';
import 'widgets/segment_tabs.dart';
import 'widgets/terminal_filter_sheet.dart';
import 'widgets/filter_chip_bar.dart';

/// 终端管理页面
/// 采用分段控制器 + 筛选标签模式
class TerminalPage extends ConsumerStatefulWidget {
  const TerminalPage({super.key});

  @override
  ConsumerState<TerminalPage> createState() => _TerminalPageState();
}

class _TerminalPageState extends ConsumerState<TerminalPage> {
  final TextEditingController _searchController = TextEditingController();
  Timer? _debounceTimer;

  // 当前选中的分段索引（0:全部, 1:已激活, 2:未激活, 3:未绑定）
  int _selectedSegmentIndex = 0;

  // 分段配置
  List<TerminalSegment> _segments = TerminalSegmentTabs.defaultSegments;

  // 当前筛选条件
  TerminalFilterCondition _filterCondition = TerminalFilterCondition.empty;

  @override
  void initState() {
    super.initState();
    // 初始化加载数据
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(terminalListProvider.notifier).loadTerminals();
      _updateSegmentCounts();
    });
  }

  @override
  void dispose() {
    _searchController.dispose();
    _debounceTimer?.cancel();
    super.dispose();
  }

  /// 更新分段数量统计
  void _updateSegmentCounts() {
    final statsAsync = ref.read(terminalStatsProvider);
    statsAsync.whenData((stats) {
      setState(() {
        _segments = [
          TerminalSegment(key: 'all', label: '全部', count: stats.total),
          TerminalSegment(
              key: 'active', label: '已激活', count: stats.activatedCount),
          TerminalSegment(
              key: 'inactive', label: '未激活', count: stats.inactiveCount),
          TerminalSegment(
              key: 'unbound', label: '未绑定', count: stats.stockCount),
        ];
      });
    });
  }

  /// 处理分段切换
  void _onSegmentChanged(int index) {
    setState(() {
      _selectedSegmentIndex = index;
    });

    final segmentKey = _segments[index].key;
    final statusGroup = segmentKey == 'all' ? null : segmentKey;

    ref.read(terminalListProvider.notifier).loadTerminals(
          statusGroup: statusGroup,
          channelId: _filterCondition.channelId,
          brandCode: _filterCondition.brandCode,
          modelCode: _filterCondition.modelCode,
          keyword:
              _searchController.text.isNotEmpty ? _searchController.text : null,
        );
  }

  /// 处理搜索输入（防抖300ms）
  void _onSearchChanged(String value) {
    _debounceTimer?.cancel();
    _debounceTimer = Timer(const Duration(milliseconds: 300), () {
      _applyFilters();
    });
  }

  /// 应用筛选条件
  void _applyFilters() {
    final segmentKey = _segments[_selectedSegmentIndex].key;
    String? statusGroup = segmentKey == 'all' ? null : segmentKey;

    // 如果有更多状态筛选，使用它覆盖分段状态
    if (_filterCondition.extraStatus != null) {
      statusGroup = _filterCondition.extraStatus;
    }

    ref.read(terminalListProvider.notifier).loadTerminals(
          statusGroup: statusGroup,
          channelId: _filterCondition.channelId,
          brandCode: _filterCondition.brandCode,
          modelCode: _filterCondition.modelCode,
          keyword:
              _searchController.text.isNotEmpty ? _searchController.text : null,
        );
  }

  /// 显示筛选Sheet
  void _showFilterSheet() {
    final filterOptionsAsync = ref.read(terminalFilterOptionsProvider);
    filterOptionsAsync.whenData((options) {
      TerminalFilterSheet.show(
        context: context,
        options: options,
        initialCondition: _filterCondition,
        onApply: (condition) {
          setState(() {
            _filterCondition = condition;
          });
          _applyFilters();
        },
      );
    });
  }

  /// 移除筛选标签
  void _onRemoveFilterTag(FilterTag tag) {
    setState(() {
      switch (tag.type) {
        case FilterTagType.channel:
          _filterCondition = _filterCondition.copyWith(clearChannel: true);
          break;
        case FilterTagType.terminalType:
          _filterCondition = _filterCondition.copyWith(clearTerminalType: true);
          break;
        case FilterTagType.extraStatus:
          _filterCondition = _filterCondition.copyWith(clearExtraStatus: true);
          break;
      }
    });
    _applyFilters();
  }

  /// 清空所有筛选
  void _clearAllFilters() {
    setState(() {
      _filterCondition = TerminalFilterCondition.empty;
      _searchController.clear();
      _selectedSegmentIndex = 0;
    });
    ref.read(terminalListProvider.notifier).resetFilters();
  }

  @override
  Widget build(BuildContext context) {
    // 监听统计数据变化，更新分段数量
    ref.listen(terminalStatsProvider, (previous, next) {
      next.whenData((stats) {
        setState(() {
          _segments = [
            TerminalSegment(key: 'all', label: '全部', count: stats.total),
            TerminalSegment(
                key: 'active', label: '已激活', count: stats.activatedCount),
            TerminalSegment(
                key: 'inactive', label: '未激活', count: stats.inactiveCount),
            TerminalSegment(
                key: 'unbound', label: '未绑定', count: stats.stockCount),
          ];
        });
      });
    });

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: _buildAppBar(),
      body: Column(
        children: [
          // 搜索框（常驻）
          _buildSearchBar(),
          // 分段控制器
          TerminalSegmentTabs(
            segments: _segments,
            selectedIndex: _selectedSegmentIndex,
            onChanged: _onSegmentChanged,
          ),
          // 筛选标签条（有筛选时显示）
          FilterChipBar(
            tags: _filterCondition.filterTags,
            onRemove: _onRemoveFilterTag,
            onClearAll: _clearAllFilters,
          ),
          // 快捷入口
          _buildQuickActions(),
          // 终端列表
          Expanded(child: _buildTerminalList()),
        ],
      ),
      bottomNavigationBar: _buildBottomBar(),
    );
  }

  /// 构建AppBar
  PreferredSizeWidget _buildAppBar() {
    return AppBar(
      title: const Text('终端管理'),
      actions: [
        // 漏斗筛选按钮
        Stack(
          children: [
            IconButton(
              onPressed: _showFilterSheet,
              icon: const Icon(Icons.filter_list),
              tooltip: '筛选',
            ),
            // 有筛选时显示小红点
            if (_filterCondition.hasFilters)
              Positioned(
                right: 8,
                top: 8,
                child: Container(
                  width: 8,
                  height: 8,
                  decoration: const BoxDecoration(
                    color: AppColors.danger,
                    shape: BoxShape.circle,
                  ),
                ),
              ),
          ],
        ),
      ],
    );
  }

  /// 构建搜索框
  Widget _buildSearchBar() {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      color: Colors.white,
      child: TextField(
        controller: _searchController,
        onChanged: _onSearchChanged,
        decoration: InputDecoration(
          hintText: '搜索终端SN或商户号...',
          hintStyle: const TextStyle(
            fontSize: 14,
            color: AppColors.textTertiary,
          ),
          prefixIcon: const Icon(
            Icons.search,
            color: AppColors.textTertiary,
            size: 20,
          ),
          suffixIcon: _searchController.text.isNotEmpty
              ? IconButton(
                  onPressed: () {
                    _searchController.clear();
                    _applyFilters();
                  },
                  icon: const Icon(
                    Icons.close,
                    color: AppColors.textTertiary,
                    size: 18,
                  ),
                )
              : null,
          filled: true,
          fillColor: Colors.grey.shade100,
          contentPadding:
              const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(10),
            borderSide: BorderSide.none,
          ),
        ),
        style: const TextStyle(fontSize: 14),
      ),
    );
  }

  /// 快捷入口：划拨记录、回拨记录
  Widget _buildQuickActions() {
    return Container(
      margin: const EdgeInsets.only(left: 16, right: 16, top: 8, bottom: 8),
      child: Row(
        children: [
          Expanded(
            child: _buildQuickActionCard(
              icon: Icons.send_outlined,
              title: '划拨记录',
              subtitle: '查看划拨历史',
              color: const Color(0xFF4CAF50),
              onTap: () => context.push(RoutePaths.terminalDistributeList),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: _buildQuickActionCard(
              icon: Icons.undo_outlined,
              title: '回拨记录',
              subtitle: '查看回拨历史',
              color: const Color(0xFFFF9800),
              onTap: () => context.push(RoutePaths.terminalRecallList),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildQuickActionCard({
    required IconData icon,
    required String title,
    required String subtitle,
    required Color color,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(10),
          border: Border.all(color: color.withValues(alpha: 0.2)),
        ),
        child: Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: color.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(icon, color: color, size: 22),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: const TextStyle(
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    subtitle,
                    style: const TextStyle(
                      fontSize: 11,
                      color: AppColors.textTertiary,
                    ),
                  ),
                ],
              ),
            ),
            Icon(Icons.chevron_right,
                color: color.withValues(alpha: 0.5), size: 20),
          ],
        ),
      ),
    );
  }

  Widget _buildTerminalList() {
    final listState = ref.watch(terminalListProvider);
    final selectedTerminals = ref.watch(selectedTerminalsProvider);

    // 显示当前筛选条件
    final hasFilters = _filterCondition.hasFilters ||
        _searchController.text.isNotEmpty ||
        _selectedSegmentIndex != 0;

    if (listState.error != null && listState.terminals.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline, size: 48, color: Colors.grey),
            const SizedBox(height: 12),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24),
              child: Text(
                '加载失败: ${listState.error}',
                textAlign: TextAlign.center,
                style: const TextStyle(fontSize: 14),
              ),
            ),
            const SizedBox(height: 12),
            ElevatedButton(
              onPressed: () =>
                  ref.read(terminalListProvider.notifier).refresh(),
              child: const Text('重试'),
            ),
          ],
        ),
      );
    }

    if (listState.terminals.isEmpty && !listState.isLoading) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.inbox_outlined, size: 48, color: Colors.grey),
            const SizedBox(height: 12),
            Text(
              hasFilters ? '没有符合条件的终端' : '暂无终端数据',
              style: const TextStyle(color: Colors.grey),
            ),
            if (hasFilters) ...[
              const SizedBox(height: 12),
              TextButton(
                onPressed: _clearAllFilters,
                child: const Text('清空筛选条件'),
              ),
            ],
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: () async {
        await ref.read(terminalListProvider.notifier).refresh();
        ref.invalidate(terminalStatsProvider);
      },
      child: NotificationListener<ScrollNotification>(
        onNotification: (notification) {
          if (notification is ScrollEndNotification &&
              notification.metrics.pixels >=
                  notification.metrics.maxScrollExtent - 100) {
            ref.read(terminalListProvider.notifier).loadMore();
          }
          return false;
        },
        child: ListView.builder(
          padding: const EdgeInsets.symmetric(horizontal: 16),
          itemCount: listState.terminals.length + (listState.isLoading ? 1 : 0),
          itemBuilder: (context, index) {
            if (index >= listState.terminals.length) {
              return const Center(
                child: Padding(
                  padding: EdgeInsets.all(16),
                  child: CircularProgressIndicator(),
                ),
              );
            }

            final terminal = listState.terminals[index];
            final isSelected = selectedTerminals.contains(terminal);

            return _buildTerminalCard(terminal, isSelected);
          },
        ),
      ),
    );
  }

  Widget _buildTerminalCard(Terminal terminal, bool isSelected) {
    return GestureDetector(
      onTap: () => _toggleSelection(terminal),
      onLongPress: () => _showTerminalActions(terminal),
      child: Container(
        margin: const EdgeInsets.only(bottom: 12),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(12),
          border: isSelected
              ? Border.all(color: AppColors.primary, width: 2)
              : null,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Expanded(
                  child: Row(
                    children: [
                      if (isSelected)
                        const Padding(
                          padding: EdgeInsets.only(right: 8),
                          child: Icon(Icons.check_circle,
                              color: AppColors.primary, size: 20),
                        ),
                      Expanded(
                        child: Text('SN: ${terminal.terminalSn}',
                            style: const TextStyle(
                                fontSize: 16, fontWeight: FontWeight.w600),
                            overflow: TextOverflow.ellipsis),
                      ),
                    ],
                  ),
                ),
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: terminal.isActivated
                        ? AppColors.success.withValues(alpha: 0.1)
                        : AppColors.textTertiary.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(terminal.status.label,
                      style: TextStyle(
                          fontSize: 12,
                          color: terminal.isActivated
                              ? AppColors.success
                              : AppColors.textTertiary)),
                ),
              ],
            ),
            const SizedBox(height: 8),
            // 显示通道和终端类型
            Row(
              children: [
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                  decoration: BoxDecoration(
                    color: AppColors.primary.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    terminal.channelCode,
                    style: const TextStyle(
                      fontSize: 11,
                      color: AppColors.primary,
                    ),
                  ),
                ),
                if (terminal.brandCode != null &&
                    terminal.brandCode!.isNotEmpty) ...[
                  const SizedBox(width: 6),
                  Text(
                    '${terminal.brandCode ?? ""} ${terminal.modelCode ?? ""}',
                    style: const TextStyle(
                      fontSize: 12,
                      color: AppColors.textSecondary,
                    ),
                  ),
                ],
              ],
            ),
            const SizedBox(height: 4),
            Text(
                '商户: ${terminal.merchantNo != null && terminal.merchantNo!.isNotEmpty ? terminal.merchantNo : "-"}',
                style: const TextStyle(
                    fontSize: 14, color: AppColors.textSecondary)),
            if (terminal.activatedAt != null)
              Text(
                  '激活时间: ${terminal.activatedAt!.toLocal().toString().substring(0, 16)}',
                  style: const TextStyle(
                      fontSize: 12, color: AppColors.textTertiary)),
          ],
        ),
      ),
    );
  }

  void _toggleSelection(Terminal terminal) {
    final currentSelection = ref.read(selectedTerminalsProvider);
    if (currentSelection.contains(terminal)) {
      ref.read(selectedTerminalsProvider.notifier).state =
          currentSelection.where((t) => t.id != terminal.id).toList();
    } else {
      ref.read(selectedTerminalsProvider.notifier).state = [
        ...currentSelection,
        terminal
      ];
    }
  }

  void _showTerminalActions(Terminal terminal) {
    showModalBottomSheet(
      context: context,
      builder: (context) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              leading: const Icon(Icons.info_outline),
              title: const Text('查看详情'),
              onTap: () {
                Navigator.pop(context);
                context.push('/terminal/${terminal.terminalSn}');
              },
            ),
            ListTile(
              leading: const Icon(Icons.history),
              title: const Text('查看流动记录'),
              onTap: () {
                Navigator.pop(context);
                context.push('/terminal/${terminal.terminalSn}/flow-logs');
              },
            ),
            if (terminal.canDistribute)
              ListTile(
                leading: const Icon(Icons.send),
                title: const Text('下发'),
                onTap: () {
                  Navigator.pop(context);
                  _navigateToTransfer([terminal]);
                },
              ),
            if (terminal.canRecall)
              ListTile(
                leading: const Icon(Icons.undo),
                title: const Text('回拨'),
                onTap: () {
                  Navigator.pop(context);
                  _showRecallDialog([terminal]);
                },
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildBottomBar() {
    final selectedTerminals = ref.watch(selectedTerminalsProvider);

    return Container(
      padding: EdgeInsets.only(
          left: 16,
          right: 16,
          top: 12,
          bottom: MediaQuery.of(context).padding.bottom + 12),
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(
              color: Colors.black.withValues(alpha: 0.05),
              blurRadius: 10,
              offset: const Offset(0, -2))
        ],
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (selectedTerminals.isNotEmpty)
            Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    '已选${selectedTerminals.length}台',
                    style: const TextStyle(
                        color: AppColors.primary, fontWeight: FontWeight.w500),
                  ),
                  TextButton(
                    onPressed: () =>
                        ref.read(selectedTerminalsProvider.notifier).state = [],
                    child: const Text('清空选择'),
                  ),
                ],
              ),
            ),
          Row(
            children: [
              // 批量设置按钮
              IconButton(
                onPressed: selectedTerminals.isEmpty
                    ? null
                    : () => _showBatchSetMenu(selectedTerminals),
                icon: Icon(
                  Icons.settings,
                  color: selectedTerminals.isEmpty
                      ? AppColors.textTertiary
                      : AppColors.primary,
                ),
                tooltip: '批量设置',
              ),
              const SizedBox(width: 8),
              Expanded(
                  child: OutlinedButton(
                onPressed: selectedTerminals.isEmpty
                    ? null
                    : () => _handleBatchRecall(selectedTerminals),
                child: const Text('批量回拨'),
              )),
              const SizedBox(width: 12),
              Expanded(
                  child: ElevatedButton(
                onPressed: selectedTerminals.isEmpty
                    ? null
                    : () => _navigateToTransfer(selectedTerminals),
                child: const Text('批量划拨'),
              )),
            ],
          ),
        ],
      ),
    );
  }

  void _navigateToTransfer(List<Terminal> terminals) {
    final snList = terminals.map((t) => t.terminalSn).toList();
    context.push('/terminal/transfer', extra: snList);
  }

  void _handleBatchRecall(List<Terminal> terminals) {
    // 检查是否有已激活的终端
    final activatedTerminals = terminals.where((t) => t.isActivated).toList();
    if (activatedTerminals.isNotEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('已激活的终端不能回拨')),
      );
      return;
    }

    _showRecallDialog(terminals);
  }

  void _showRecallDialog(List<Terminal> terminals) {
    showDialog(
      context: context,
      builder: (dialogContext) => AlertDialog(
        title: const Text('终端回拨'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('确定要回拨 ${terminals.length} 台终端吗？'),
            const SizedBox(height: 8),
            const Text(
              '回拨后终端将归还给上级代理商',
              style: TextStyle(fontSize: 12, color: AppColors.textSecondary),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(dialogContext),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () async {
              Navigator.pop(dialogContext);
              await _executeRecall(terminals);
            },
            child: const Text('确认回拨'),
          ),
        ],
      ),
    );
  }

  /// 执行回拨操作
  Future<void> _executeRecall(List<Terminal> terminals) async {
    // 获取当前代理商信息，获取上级代理商ID
    final myProfileAsync = ref.read(myProfileProvider);

    await myProfileAsync.when(
      data: (agentDetail) async {
        // 检查是否有上级代理商
        if (agentDetail.parentId == null || agentDetail.parentId == 0) {
          if (mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Text('您是顶级代理商，无法回拨'),
                backgroundColor: Colors.orange,
              ),
            );
          }
          return;
        }

        // 执行回拨
        final success =
            await ref.read(terminalRecallProvider.notifier).batchRecall(
                  toAgentId: agentDetail.parentId!,
                  terminalSns: terminals.map((t) => t.terminalSn).toList(),
                );

        if (mounted) {
          final recallState = ref.read(terminalRecallProvider);
          if (success) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text('回拨成功，共${recallState.successCount}台终端'),
                backgroundColor: AppColors.success,
              ),
            );
            // 清空选中状态
            ref.read(selectedTerminalsProvider.notifier).state = [];
          } else if (recallState.failedCount > 0) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(
                    '部分回拨失败: 成功${recallState.successCount}台, 失败${recallState.failedCount}台'),
                backgroundColor: Colors.orange,
              ),
            );
          } else {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text('回拨失败: ${recallState.error ?? "未知错误"}'),
                backgroundColor: Colors.red,
              ),
            );
          }
        }
      },
      loading: () {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('正在获取代理商信息...')),
        );
      },
      error: (error, stack) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('获取代理商信息失败: $error'),
            backgroundColor: Colors.red,
          ),
        );
      },
    );
  }

  /// 显示批量设置菜单
  void _showBatchSetMenu(List<Terminal> terminals) {
    final snList = terminals.map((t) => t.terminalSn).toList();

    showModalBottomSheet(
      context: context,
      builder: (context) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              padding: const EdgeInsets.all(16),
              child: const Text(
                '批量设置',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
            const Divider(height: 1),
            ListTile(
              leading: const Icon(Icons.percent),
              title: const Text('设置费率'),
              subtitle: const Text('设置终端信用卡费率'),
              onTap: () {
                Navigator.pop(context);
                context.push(RoutePaths.terminalBatchSetRate, extra: snList);
              },
            ),
            ListTile(
              leading: const Icon(Icons.sim_card),
              title: const Text('设置流量费'),
              subtitle: const Text('设置流量费金额和收费周期'),
              onTap: () {
                Navigator.pop(context);
                context.push(RoutePaths.terminalBatchSetSim, extra: snList);
              },
            ),
            ListTile(
              leading: const Icon(Icons.account_balance_wallet),
              title: const Text('设置押金'),
              subtitle: const Text('设置终端激活押金'),
              onTap: () {
                Navigator.pop(context);
                context.push(RoutePaths.terminalBatchSetDeposit, extra: snList);
              },
            ),
            ListTile(
              leading: const Icon(Icons.history),
              title: const Text('查看流动记录'),
              subtitle: const Text('查看选中终端的流动历史'),
              enabled: terminals.length == 1,
              onTap: terminals.length == 1
                  ? () {
                      Navigator.pop(context);
                      context.push(
                          '/terminal/${terminals.first.terminalSn}/flow-logs');
                    }
                  : null,
            ),
            const SizedBox(height: 8),
          ],
        ),
      ),
    );
  }
}
