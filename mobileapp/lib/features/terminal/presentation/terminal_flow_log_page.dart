import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import '../domain/models/terminal.dart';
import 'providers/terminal_provider.dart';

/// 终端流动记录页面
class TerminalFlowLogPage extends ConsumerStatefulWidget {
  final String terminalSn;

  const TerminalFlowLogPage({
    super.key,
    required this.terminalSn,
  });

  @override
  ConsumerState<TerminalFlowLogPage> createState() => _TerminalFlowLogPageState();
}

class _TerminalFlowLogPageState extends ConsumerState<TerminalFlowLogPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final List<Map<String, String>> _tabs = [
    {'key': 'all', 'label': '全部'},
    {'key': 'distribute', 'label': '下发'},
    {'key': 'recall', 'label': '回拨'},
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: _tabs.length, vsync: this);
    _tabController.addListener(_onTabChanged);
    // 初始加载数据
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(flowLogListProvider(widget.terminalSn).notifier).loadList();
    });
  }

  @override
  void dispose() {
    _tabController.removeListener(_onTabChanged);
    _tabController.dispose();
    super.dispose();
  }

  void _onTabChanged() {
    if (!_tabController.indexIsChanging) {
      final logType = _tabs[_tabController.index]['key']!;
      ref.read(flowLogListProvider(widget.terminalSn).notifier).setLogType(logType);
    }
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(flowLogListProvider(widget.terminalSn));

    return Scaffold(
      appBar: AppBar(
        title: const Text('流动记录'),
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(80),
          child: Column(
            children: [
              // 终端信息
              if (state.terminal != null)
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                  child: Row(
                    children: [
                      Text(
                        'SN: ${state.terminal!.terminalSn}',
                        style: const TextStyle(
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                      const SizedBox(width: 16),
                      Text(
                        '${state.terminal!.channelCode} ${state.terminal!.brandCode} ${state.terminal!.modelCode}',
                        style: TextStyle(
                          color: Colors.white.withOpacity(0.8),
                          fontSize: 12,
                        ),
                      ),
                    ],
                  ),
                ),
              // Tab栏
              TabBar(
                controller: _tabController,
                tabs: _tabs.map((tab) => Tab(text: tab['label'])).toList(),
                indicatorColor: Colors.white,
                labelColor: Colors.white,
                unselectedLabelColor: Colors.white60,
              ),
            ],
          ),
        ),
      ),
      body: _buildBody(state),
    );
  }

  Widget _buildBody(FlowLogListState state) {
    if (state.isLoading && state.list.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (state.error != null && state.list.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(state.error!, style: const TextStyle(color: Colors.red)),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: () {
                ref.read(flowLogListProvider(widget.terminalSn).notifier).refresh();
              },
              child: const Text('重试'),
            ),
          ],
        ),
      );
    }

    if (state.list.isEmpty) {
      return const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.history, size: 64, color: Colors.grey),
            SizedBox(height: 16),
            Text('暂无流动记录', style: TextStyle(color: Colors.grey)),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: () => ref.read(flowLogListProvider(widget.terminalSn).notifier).refresh(),
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: state.list.length + (state.hasMore ? 1 : 0),
        itemBuilder: (context, index) {
          if (index == state.list.length) {
            // 加载更多
            if (!state.isLoadingMore) {
              ref.read(flowLogListProvider(widget.terminalSn).notifier).loadList();
            }
            return const Center(
              child: Padding(
                padding: EdgeInsets.all(16),
                child: CircularProgressIndicator(),
              ),
            );
          }

          final log = state.list[index];
          final isLast = index == state.list.length - 1;

          return _buildLogItem(log, isLast);
        },
      ),
    );
  }

  Widget _buildLogItem(TerminalFlowLog log, bool isLast) {
    final dateFormat = DateFormat('yyyy-MM-dd HH:mm');

    // 根据日志类型设置图标和颜色
    IconData icon;
    Color color;
    switch (log.logType) {
      case 'distribute':
        icon = Icons.arrow_forward;
        color = Colors.blue;
        break;
      case 'recall':
        icon = Icons.arrow_back;
        color = Colors.orange;
        break;
      case 'bind':
        icon = Icons.link;
        color = Colors.green;
        break;
      case 'unbind':
        icon = Icons.link_off;
        color = Colors.red;
        break;
      case 'activate':
        icon = Icons.check_circle;
        color = Colors.teal;
        break;
      default:
        icon = Icons.history;
        color = Colors.grey;
    }

    return IntrinsicHeight(
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 时间线
          SizedBox(
            width: 40,
            child: Column(
              children: [
                Container(
                  width: 32,
                  height: 32,
                  decoration: BoxDecoration(
                    color: color.withOpacity(0.1),
                    shape: BoxShape.circle,
                  ),
                  child: Icon(icon, size: 16, color: color),
                ),
                if (!isLast)
                  Expanded(
                    child: Container(
                      width: 2,
                      color: Colors.grey.shade300,
                    ),
                  ),
              ],
            ),
          ),
          const SizedBox(width: 12),
          // 内容
          Expanded(
            child: Container(
              margin: const EdgeInsets.only(bottom: 16),
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(8),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.05),
                    blurRadius: 4,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // 标题行
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Row(
                        children: [
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 8,
                              vertical: 2,
                            ),
                            decoration: BoxDecoration(
                              color: color.withOpacity(0.1),
                              borderRadius: BorderRadius.circular(4),
                            ),
                            child: Text(
                              log.logTypeName.isNotEmpty ? log.logTypeName : _getLogTypeName(log.logType),
                              style: TextStyle(
                                color: color,
                                fontSize: 12,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ),
                          const SizedBox(width: 8),
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 6,
                              vertical: 2,
                            ),
                            decoration: BoxDecoration(
                              color: _getStatusColor(log.status).withOpacity(0.1),
                              borderRadius: BorderRadius.circular(4),
                            ),
                            child: Text(
                              log.statusName.isNotEmpty ? log.statusName : _getStatusName(log.status),
                              style: TextStyle(
                                color: _getStatusColor(log.status),
                                fontSize: 11,
                              ),
                            ),
                          ),
                        ],
                      ),
                      Text(
                        dateFormat.format(log.createdAt),
                        style: TextStyle(
                          color: Colors.grey.shade600,
                          fontSize: 12,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  // 详情内容
                  if (log.logType == 'distribute' || log.logType == 'recall') ...[
                    Row(
                      children: [
                        Text(
                          log.fromAgentName.isNotEmpty ? log.fromAgentName : '代理商${log.fromAgentId}',
                          style: const TextStyle(fontSize: 13),
                        ),
                        const Padding(
                          padding: EdgeInsets.symmetric(horizontal: 8),
                          child: Icon(Icons.arrow_forward, size: 14, color: Colors.grey),
                        ),
                        Text(
                          log.toAgentName.isNotEmpty ? log.toAgentName : '代理商${log.toAgentId}',
                          style: const TextStyle(fontSize: 13),
                        ),
                      ],
                    ),
                  ],
                  if (log.merchantNo.isNotEmpty) ...[
                    const SizedBox(height: 4),
                    Text(
                      '商户: ${log.merchantNo}',
                      style: TextStyle(
                        color: Colors.grey.shade600,
                        fontSize: 12,
                      ),
                    ),
                  ],
                  if (log.remark.isNotEmpty) ...[
                    const SizedBox(height: 4),
                    Text(
                      '备注: ${log.remark}',
                      style: TextStyle(
                        color: Colors.grey.shade600,
                        fontSize: 12,
                      ),
                    ),
                  ],
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  String _getLogTypeName(String logType) {
    switch (logType) {
      case 'distribute':
        return '下发';
      case 'recall':
        return '回拨';
      case 'bind':
        return '绑定';
      case 'unbind':
        return '解绑';
      case 'activate':
        return '激活';
      default:
        return '未知';
    }
  }

  String _getStatusName(int status) {
    switch (status) {
      case 1:
        return '待确认';
      case 2:
        return '已确认';
      case 3:
        return '已拒绝';
      case 4:
        return '已取消';
      default:
        return '未知';
    }
  }

  Color _getStatusColor(int status) {
    switch (status) {
      case 1:
        return Colors.orange;
      case 2:
        return Colors.green;
      case 3:
        return Colors.red;
      case 4:
        return Colors.grey;
      default:
        return Colors.grey;
    }
  }
}
