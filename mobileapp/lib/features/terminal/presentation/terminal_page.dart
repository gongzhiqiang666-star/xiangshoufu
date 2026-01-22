import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_colors.dart';
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
  final List<Map<String, dynamic>> _tabs = [
    {'label': '全部', 'status': null},
    {'label': '已激活', 'status': TerminalStatus.activated.value},
    {'label': '未激活', 'status': TerminalStatus.bound.value},
    {'label': '库存', 'status': TerminalStatus.pending.value},
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: _tabs.length, vsync: this);
    _tabController.addListener(_onTabChanged);

    // 初始化加载数据
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(terminalListProvider.notifier).loadTerminals();
    });
  }

  void _onTabChanged() {
    if (!_tabController.indexIsChanging) {
      final status = _tabs[_tabController.index]['status'] as int?;
      ref.read(terminalListProvider.notifier).setStatusFilter(status);
    }
  }

  @override
  void dispose() {
    _tabController.removeListener(_onTabChanged);
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('终端管理'),
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(48),
          child: Container(
            color: Colors.white,
            child: TabBar(
              controller: _tabController,
              tabs: _tabs.map((e) => Tab(text: e['label'] as String)).toList(),
            ),
          ),
        ),
      ),
      body: Column(
        children: [
          _buildStatistics(),
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: _tabs.map((tab) => _buildTerminalList()).toList(),
            ),
          ),
        ],
      ),
      bottomNavigationBar: _buildBottomBar(),
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
      return const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.inbox_outlined, size: 48, color: Colors.grey),
            SizedBox(height: 12),
            Text('暂无终端数据', style: TextStyle(color: Colors.grey)),
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
                        ? AppColors.success.withOpacity(0.1)
                        : AppColors.textTertiary.withOpacity(0.1),
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
                // TODO: 跳转到详情页
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('详情页开发中...')),
                );
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
              color: Colors.black.withOpacity(0.05),
              blurRadius: 10,
              offset: const Offset(0, -2))
        ],
      ),
      child: Row(
        children: [
          if (selectedTerminals.isNotEmpty)
            Padding(
              padding: const EdgeInsets.only(right: 12),
              child: Text(
                '已选${selectedTerminals.length}台',
                style: const TextStyle(
                    color: AppColors.primary, fontWeight: FontWeight.w500),
              ),
            ),
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
      builder: (context) => AlertDialog(
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
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              // TODO: 实现回拨功能，需要选择上级代理商
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('回拨功能需要选择上级代理商，暂未完整实现')),
              );
            },
            child: const Text('确认回拨'),
          ),
        ],
      ),
    );
  }
}
