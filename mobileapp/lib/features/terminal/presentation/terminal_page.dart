import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';

/// 终端管理页面
class TerminalPage extends StatefulWidget {
  const TerminalPage({super.key});

  @override
  State<TerminalPage> createState() => _TerminalPageState();
}

class _TerminalPageState extends State<TerminalPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final List<String> _tabs = ['全部', '已激活', '未激活', '库存'];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: _tabs.length, vsync: this);
  }

  @override
  void dispose() {
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
              tabs: _tabs.map((e) => Tab(text: e)).toList(),
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
    return Container(
      margin: const EdgeInsets.all(16),
      child: Row(
        children: [
          _buildStatCard('终端总数', '200', AppColors.primary),
          const SizedBox(width: 12),
          _buildStatCard('已激活', '180', AppColors.success),
          const SizedBox(width: 12),
          _buildStatCard('未激活', '20', AppColors.warning),
          const SizedBox(width: 12),
          _buildStatCard('今日激活', '5', AppColors.profitReward),
        ],
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
            Text(value, style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: color)),
            const SizedBox(height: 4),
            Text(title, style: const TextStyle(fontSize: 12, color: AppColors.textSecondary)),
          ],
        ),
      ),
    );
  }

  Widget _buildTerminalList() {
    return ListView.builder(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      itemCount: 10,
      itemBuilder: (context, index) {
        final isActivated = index % 3 != 0;
        return Container(
          margin: const EdgeInsets.only(bottom: 12),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(12),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text('SN: 1234567${index.toString().padLeft(2, '0')}',
                      style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w600)),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: isActivated ? AppColors.success.withOpacity(0.1) : AppColors.textTertiary.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(isActivated ? '已激活' : '未激活',
                        style: TextStyle(fontSize: 12, color: isActivated ? AppColors.success : AppColors.textTertiary)),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              Text('商户: ${isActivated ? "张三商店" : "-"}', style: const TextStyle(fontSize: 14, color: AppColors.textSecondary)),
            ],
          ),
        );
      },
    );
  }

  Widget _buildBottomBar() {
    return Container(
      padding: EdgeInsets.only(left: 16, right: 16, top: 12, bottom: MediaQuery.of(context).padding.bottom + 12),
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [BoxShadow(color: Colors.black.withOpacity(0.05), blurRadius: 10, offset: const Offset(0, -2))],
      ),
      child: Row(
        children: [
          Expanded(child: OutlinedButton(onPressed: () {}, child: const Text('批量回拨'))),
          const SizedBox(width: 12),
          Expanded(child: ElevatedButton(onPressed: () {}, child: const Text('批量划拨'))),
        ],
      ),
    );
  }
}
