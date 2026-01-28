import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_colors.dart';
import '../../../router/app_router.dart';
import '../../agent/presentation/providers/agent_provider.dart';
import '../domain/models/terminal.dart';
import 'providers/terminal_provider.dart';

/// 终端管理页面
class TerminalPage extends ConsumerStatefulWidget {
  const TerminalPage({super.key});

  @override
  ConsumerState<TerminalPage> createState() => _TerminalPageState();
}

class _TerminalPageState extends ConsumerState<TerminalPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final TextEditingController _searchController = TextEditingController();

  // 状态分组Tab配置
  final List<Map<String, String>> _statusTabs = [
    {'key': 'all', 'label': '全部'},
    {'key': 'unstock', 'label': '未出库'},
    {'key': 'stocked', 'label': '已出库'},
    {'key': 'unbound', 'label': '未绑定'},
    {'key': 'inactive', 'label': '未激活'},
    {'key': 'active', 'label': '已激活'},
  ];

  // 当前筛选条件
  int? _selectedChannelId;
  String? _selectedBrandCode;
  String? _selectedModelCode;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: _statusTabs.length, vsync: this);
    _tabController.addListener(_onTabChanged);

    // 初始化加载数据
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(terminalListProvider.notifier).loadTerminals();
    });
  }

  void _onTabChanged() {
    if (!_tabController.indexIsChanging) {
      final statusGroup = _statusTabs[_tabController.index]['key']!;
      ref.read(terminalListProvider.notifier).setStatusGroup(
        statusGroup == 'all' ? null : statusGroup,
      );
    }
  }

  @override
  void dispose() {
    _tabController.removeListener(_onTabChanged);
    _tabController.dispose();
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('终端管理'),
        actions: [
          PopupMenuButton<String>(
            icon: const Icon(Icons.more_vert),
            onSelected: (value) => _handleMenuAction(value),
            itemBuilder: (context) => [
              const PopupMenuItem(
                value: 'distribute_list',
                child: Row(
                  children: [
                    Icon(Icons.send, size: 20, color: AppColors.textSecondary),
                    SizedBox(width: 8),
                    Text('划拨记录'),
                  ],
                ),
              ),
              const PopupMenuItem(
                value: 'recall_list',
                child: Row(
                  children: [
                    Icon(Icons.undo, size: 20, color: AppColors.textSecondary),
                    SizedBox(width: 8),
                    Text('回拨记录'),
                  ],
                ),
              ),
            ],
          ),
        ],
      ),
      body: Column(
        children: [
          // 筛选栏
          _buildFilterBar(),
          // 状态Tab栏
          _buildStatusTabBar(),
          // 统计卡片
          _buildStatistics(),
          // 快捷入口
          _buildQuickActions(),
          // 终端列表
          Expanded(child: _buildTerminalList()),
        ],
      ),
      bottomNavigationBar: _buildBottomBar(),
    );
  }

  /// 构建筛选栏
  Widget _buildFilterBar() {
    final filterOptionsAsync = ref.watch(terminalFilterOptionsProvider);

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      color: Colors.white,
      child: filterOptionsAsync.when(
        data: (options) => Row(
          children: [
            // 通道筛选下拉
            Expanded(
              child: _buildDropdown<int?>(
                value: _selectedChannelId,
                hint: '全部通道',
                items: [
                  const DropdownMenuItem<int?>(
                    value: null,
                    child: Text('全部通道'),
                  ),
                  ...options.channels.map((channel) => DropdownMenuItem<int?>(
                    value: channel.channelId,
                    child: Text(channel.channelCode),
                  )),
                ],
                onChanged: (value) {
                  setState(() {
                    _selectedChannelId = value;
                    _selectedBrandCode = null;
                    _selectedModelCode = null;
                  });
                  _applyFilters();
                },
              ),
            ),
            const SizedBox(width: 8),
            // 终端类型筛选下拉
            Expanded(
              child: _buildTerminalTypeDropdown(options),
            ),
            const SizedBox(width: 8),
            // 搜索按钮
            SizedBox(
              width: 40,
              height: 40,
              child: IconButton(
                onPressed: _showSearchDialog,
                icon: const Icon(Icons.search, color: AppColors.primary),
                style: IconButton.styleFrom(
                  backgroundColor: AppColors.primary.withValues(alpha: 0.1),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
              ),
            ),
          ],
        ),
        loading: () => Row(
          children: [
            Expanded(
              child: Container(
                height: 40,
                decoration: BoxDecoration(
                  color: Colors.grey.shade100,
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
            ),
            const SizedBox(width: 8),
            Expanded(
              child: Container(
                height: 40,
                decoration: BoxDecoration(
                  color: Colors.grey.shade100,
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
            ),
            const SizedBox(width: 8),
            const SizedBox(width: 40, height: 40),
          ],
        ),
        error: (_, __) => const SizedBox.shrink(),
      ),
    );
  }

  /// 构建下拉选择器
  Widget _buildDropdown<T>({
    required T value,
    required String hint,
    required List<DropdownMenuItem<T>> items,
    required ValueChanged<T?> onChanged,
  }) {
    return Container(
      height: 40,
      padding: const EdgeInsets.symmetric(horizontal: 12),
      decoration: BoxDecoration(
        color: Colors.grey.shade50,
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: Colors.grey.shade200),
      ),
      child: DropdownButtonHideUnderline(
        child: DropdownButton<T>(
          value: value,
          hint: Text(hint, style: const TextStyle(fontSize: 14)),
          isExpanded: true,
          icon: const Icon(Icons.arrow_drop_down, size: 20),
          style: const TextStyle(fontSize: 14, color: AppColors.textPrimary),
          items: items,
          onChanged: onChanged,
        ),
      ),
    );
  }

  /// 构建终端类型下拉
  Widget _buildTerminalTypeDropdown(TerminalFilterOptions options) {
    // 根据选中的通道过滤终端类型
    final filteredTypes = _selectedChannelId == null
        ? options.terminalTypes
        : options.terminalTypes.where((t) => t.channelId == _selectedChannelId).toList();

    // 构建唯一的终端类型列表
    final uniqueTypes = <String, TerminalTypeOption>{};
    for (final type in filteredTypes) {
      final key = '${type.brandCode}-${type.modelCode}';
      if (!uniqueTypes.containsKey(key)) {
        uniqueTypes[key] = type;
      }
    }

    return Container(
      height: 40,
      padding: const EdgeInsets.symmetric(horizontal: 12),
      decoration: BoxDecoration(
        color: Colors.grey.shade50,
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: Colors.grey.shade200),
      ),
      child: DropdownButtonHideUnderline(
        child: DropdownButton<String?>(
          value: _selectedBrandCode != null && _selectedModelCode != null
              ? '$_selectedBrandCode-$_selectedModelCode'
              : null,
          hint: const Text('全部类型', style: TextStyle(fontSize: 14)),
          isExpanded: true,
          icon: const Icon(Icons.arrow_drop_down, size: 20),
          style: const TextStyle(fontSize: 14, color: AppColors.textPrimary),
          items: [
            const DropdownMenuItem<String?>(
              value: null,
              child: Text('全部类型'),
            ),
            ...uniqueTypes.entries.map((entry) => DropdownMenuItem<String?>(
              value: entry.key,
              child: Text(entry.value.displayName, overflow: TextOverflow.ellipsis),
            )),
          ],
          onChanged: (value) {
            if (value == null) {
              setState(() {
                _selectedBrandCode = null;
                _selectedModelCode = null;
              });
            } else {
              final parts = value.split('-');
              setState(() {
                _selectedBrandCode = parts[0];
                _selectedModelCode = parts.length > 1 ? parts[1] : null;
              });
            }
            _applyFilters();
          },
        ),
      ),
    );
  }

  /// 应用筛选条件
  void _applyFilters() {
    final statusGroup = _statusTabs[_tabController.index]['key']!;
    ref.read(terminalListProvider.notifier).loadTerminals(
      channelId: _selectedChannelId,
      brandCode: _selectedBrandCode,
      modelCode: _selectedModelCode,
      statusGroup: statusGroup == 'all' ? null : statusGroup,
      keyword: _searchController.text.isNotEmpty ? _searchController.text : null,
    );
  }

  /// 显示搜索对话框
  void _showSearchDialog() {
    showDialog(
      context: context,
      builder: (dialogContext) => AlertDialog(
        title: const Text('搜索终端'),
        content: TextField(
          controller: _searchController,
          decoration: const InputDecoration(
            hintText: '输入SN号或商户号搜索',
            prefixIcon: Icon(Icons.search),
          ),
          autofocus: true,
          onSubmitted: (_) {
            Navigator.pop(dialogContext);
            _applyFilters();
          },
        ),
        actions: [
          TextButton(
            onPressed: () {
              _searchController.clear();
              Navigator.pop(dialogContext);
              _applyFilters();
            },
            child: const Text('清空'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(dialogContext);
              _applyFilters();
            },
            child: const Text('搜索'),
          ),
        ],
      ),
    );
  }

  /// 构建状态Tab栏
  Widget _buildStatusTabBar() {
    return Container(
      color: Colors.white,
      child: TabBar(
        controller: _tabController,
        isScrollable: true,
        tabAlignment: TabAlignment.start,
        tabs: _statusTabs.map((tab) => Tab(text: tab['label'])).toList(),
        labelColor: AppColors.primary,
        unselectedLabelColor: AppColors.textSecondary,
        indicatorColor: AppColors.primary,
        indicatorSize: TabBarIndicatorSize.label,
      ),
    );
  }

  Widget _buildStatistics() {
    final statsAsync = ref.watch(terminalStatsProvider);

    return statsAsync.when(
      data: (stats) => Container(
        margin: const EdgeInsets.all(16),
        child: Row(
          children: [
            _buildStatCard('终端总数', stats.total.toString(), AppColors.primary),
            const SizedBox(width: 12),
            _buildStatCard(
                '已激活', stats.activatedCount.toString(), AppColors.success),
            const SizedBox(width: 12),
            _buildStatCard(
                '未激活', stats.inactiveCount.toString(), AppColors.warning),
            const SizedBox(width: 12),
            _buildStatCard('今日激活', stats.todayActivated.toString(),
                AppColors.profitReward),
          ],
        ),
      ),
      loading: () => Container(
        margin: const EdgeInsets.all(16),
        child: Row(
          children: [
            _buildStatCard('终端总数', '-', AppColors.primary),
            const SizedBox(width: 12),
            _buildStatCard('已激活', '-', AppColors.success),
            const SizedBox(width: 12),
            _buildStatCard('未激活', '-', AppColors.warning),
            const SizedBox(width: 12),
            _buildStatCard('今日激活', '-', AppColors.profitReward),
          ],
        ),
      ),
      error: (error, stack) => Container(
        margin: const EdgeInsets.all(16),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: Colors.red.shade50,
          borderRadius: BorderRadius.circular(8),
        ),
        child: Text('加载统计失败: $error',
            style: TextStyle(color: Colors.red.shade700)),
      ),
    );
  }

  /// 快捷入口：划拨记录、回拨记录
  Widget _buildQuickActions() {
    return Container(
      margin: const EdgeInsets.only(left: 16, right: 16, bottom: 12),
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
            Icon(Icons.chevron_right, color: color.withValues(alpha: 0.5), size: 20),
          ],
        ),
      ),
    );
  }

  Widget _buildStatCard(String title, String value, Color color) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(8),
        ),
        child: Column(
          children: [
            Text(value,
                style: TextStyle(
                    fontSize: 20, fontWeight: FontWeight.bold, color: color)),
            const SizedBox(height: 4),
            Text(title,
                style: const TextStyle(
                    fontSize: 12, color: AppColors.textSecondary)),
          ],
        ),
      ),
    );
  }

  Widget _buildTerminalList() {
    final listState = ref.watch(terminalListProvider);
    final selectedTerminals = ref.watch(selectedTerminalsProvider);

    // 显示当前筛选条件
    final hasFilters = _selectedChannelId != null ||
        _selectedBrandCode != null ||
        _searchController.text.isNotEmpty;

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
                onPressed: _clearFilters,
                child: const Text('清空筛选条件'),
              ),
            ],
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: () => ref.read(terminalListProvider.notifier).refresh(),
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

  /// 清空筛选条件
  void _clearFilters() {
    setState(() {
      _selectedChannelId = null;
      _selectedBrandCode = null;
      _selectedModelCode = null;
      _searchController.clear();
    });
    _tabController.animateTo(0);
    ref.read(terminalListProvider.notifier).resetFilters();
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
                  padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
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
                if (terminal.brandCode != null && terminal.brandCode!.isNotEmpty) ...[
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
        final success = await ref.read(terminalRecallProvider.notifier).batchRecall(
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
                content: Text('部分回拨失败: 成功${recallState.successCount}台, 失败${recallState.failedCount}台'),
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

  /// 处理菜单操作
  void _handleMenuAction(String action) {
    switch (action) {
      case 'distribute_list':
        context.push(RoutePaths.terminalDistributeList);
        break;
      case 'recall_list':
        context.push(RoutePaths.terminalRecallList);
        break;
    }
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
                      context.push('/terminal/${terminals.first.terminalSn}/flow-logs');
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
