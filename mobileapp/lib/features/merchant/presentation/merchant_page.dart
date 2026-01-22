import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import '../../../core/theme/app_colors.dart';
import '../data/models/merchant_model.dart';
import 'providers/merchant_provider.dart';
import 'widgets/merchant_card.dart';
import 'merchant_detail_page.dart';

/// 商户管理页面
class MerchantPage extends ConsumerStatefulWidget {
  const MerchantPage({super.key});

  @override
  ConsumerState<MerchantPage> createState() => _MerchantPageState();
}

class _MerchantPageState extends ConsumerState<MerchantPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final statsAsync = ref.watch(merchantStatsProvider);

    return Scaffold(
      backgroundColor: Colors.grey.shade50,
      appBar: AppBar(
        title: const Text('商户管理'),
        centerTitle: true,
        elevation: 0,
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            onPressed: () => _showSearchDialog(context),
          ),
        ],
        bottom: TabBar(
          controller: _tabController,
          labelColor: AppColors.primary,
          unselectedLabelColor: AppColors.textSecondary,
          indicatorColor: AppColors.primary,
          indicatorWeight: 3,
          tabs: const [
            Tab(text: '直营商户'),
            Tab(text: '团队商户'),
          ],
        ),
      ),
      body: Column(
        children: [
          // 统计卡片
          statsAsync.when(
            data: (stats) => _buildStatsCard(stats),
            loading: () => _buildStatsCardLoading(),
            error: (_, __) => const SizedBox.shrink(),
          ),
          // 商户列表
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: const [
                _MerchantListTab(isDirect: true),
                _MerchantListTab(isDirect: false),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatsCard(MerchantStats stats) {
    return Container(
      margin: EdgeInsets.all(16.w),
      padding: EdgeInsets.symmetric(vertical: 16.h, horizontal: 12.w),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12.r),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceAround,
        children: [
          _buildStatItem('商户总数', stats.totalCount, Colors.blue),
          _buildStatItem('直营', stats.directCount, Colors.green),
          _buildStatItem('团队', stats.teamCount, Colors.purple),
          _buildStatItem('今日新增', stats.todayNewCount, Colors.orange),
        ],
      ),
    );
  }

  Widget _buildStatsCardLoading() {
    return Container(
      margin: EdgeInsets.all(16.w),
      padding: EdgeInsets.symmetric(vertical: 20.h),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12.r),
      ),
      child: const Center(
        child: CircularProgressIndicator(strokeWidth: 2),
      ),
    );
  }

  Widget _buildStatItem(String label, int value, Color color) {
    return Column(
      children: [
        Text(
          value.toString(),
          style: TextStyle(
            fontSize: 20.sp,
            fontWeight: FontWeight.bold,
            color: color,
          ),
        ),
        SizedBox(height: 4.h),
        Text(
          label,
          style: TextStyle(
            fontSize: 12.sp,
            color: AppColors.textSecondary,
          ),
        ),
      ],
    );
  }

  void _showSearchDialog(BuildContext context) {
    final controller = TextEditingController();
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('搜索商户'),
        content: TextField(
          controller: controller,
          decoration: const InputDecoration(
            hintText: '请输入商户名称/编号/机具号',
            border: OutlineInputBorder(),
          ),
          autofocus: true,
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () {
              final keyword = controller.text.trim();
              if (_tabController.index == 0) {
                ref.read(directMerchantListProvider.notifier).setKeyword(keyword);
              } else {
                ref.read(teamMerchantListProvider.notifier).setKeyword(keyword);
              }
              Navigator.pop(context);
            },
            child: const Text('搜索'),
          ),
        ],
      ),
    );
  }
}

/// 商户列表 Tab
class _MerchantListTab extends ConsumerWidget {
  final bool isDirect;

  const _MerchantListTab({required this.isDirect});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = isDirect
        ? ref.watch(directMerchantListProvider)
        : ref.watch(teamMerchantListProvider);
    final notifier = isDirect
        ? ref.read(directMerchantListProvider.notifier)
        : ref.read(teamMerchantListProvider.notifier);

    if (state.error != null && state.merchants.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.error_outline, size: 48.sp, color: Colors.grey),
            SizedBox(height: 16.h),
            Text('加载失败', style: TextStyle(fontSize: 14.sp, color: Colors.grey)),
            SizedBox(height: 16.h),
            ElevatedButton(
              onPressed: () => notifier.refresh(),
              child: const Text('重试'),
            ),
          ],
        ),
      );
    }

    if (state.merchants.isEmpty && !state.isLoading) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.store_mall_directory_outlined, size: 64.sp, color: Colors.grey.shade300),
            SizedBox(height: 16.h),
            Text(
              '暂无${isDirect ? '直营' : '团队'}商户',
              style: TextStyle(fontSize: 14.sp, color: Colors.grey),
            ),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: () => notifier.refresh(),
      child: ListView.builder(
        padding: EdgeInsets.only(top: 8.h, bottom: 16.h),
        itemCount: state.merchants.length + (state.hasMore ? 1 : 0),
        itemBuilder: (context, index) {
          if (index >= state.merchants.length) {
            // 加载更多
            if (!state.isLoading) {
              WidgetsBinding.instance.addPostFrameCallback((_) {
                notifier.loadMerchants();
              });
            }
            return Padding(
              padding: EdgeInsets.all(16.w),
              child: const Center(
                child: CircularProgressIndicator(strokeWidth: 2),
              ),
            );
          }

          final merchant = state.merchants[index];
          return MerchantCard(
            merchant: merchant,
            onTap: () => _navigateToDetail(context, merchant),
          );
        },
      ),
    );
  }

  void _navigateToDetail(BuildContext context, Merchant merchant) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => MerchantDetailPage(merchantId: merchant.id),
      ),
    );
  }
}
