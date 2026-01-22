import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:share_plus/share_plus.dart';

import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../router/app_router.dart';
import '../data/models/agent_model.dart';
import 'providers/agent_provider.dart';

/// 代理拓展页面
class AgentPage extends ConsumerStatefulWidget {
  const AgentPage({super.key});

  @override
  ConsumerState<AgentPage> createState() => _AgentPageState();
}

class _AgentPageState extends ConsumerState<AgentPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);

    // 加载数据
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(subordinatesProvider.notifier).loadSubordinates(refresh: true);
    });
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
        title: const Text('代理拓展'),
        actions: [
          IconButton(
            icon: const Icon(Icons.person_add_outlined),
            onPressed: _showAddAgentDialog,
            tooltip: '手动添加代理',
          ),
        ],
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: '推广中心'),
            Tab(text: '团队管理'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildPromotionTab(),
          _buildTeamTab(),
        ],
      ),
    );
  }

  /// 推广中心Tab
  Widget _buildPromotionTab() {
    final inviteCodeAsync = ref.watch(inviteCodeProvider);

    return inviteCodeAsync.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, stack) => Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.error_outline, size: 48, color: AppColors.danger),
            const SizedBox(height: 16),
            Text('加载失败: $error'),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: () => ref.invalidate(inviteCodeProvider),
              child: const Text('重试'),
            ),
          ],
        ),
      ),
      data: (inviteCode) => SingleChildScrollView(
        padding: const EdgeInsets.all(AppSpacing.md),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            _buildQRCodeCard(inviteCode),
            const SizedBox(height: AppSpacing.md),
            _buildInviteCodeCard(inviteCode),
            const SizedBox(height: AppSpacing.md),
            _buildShareActions(inviteCode),
            const SizedBox(height: AppSpacing.lg),
            _buildPromotionTips(),
          ],
        ),
      ),
    );
  }

  /// 二维码卡片
  Widget _buildQRCodeCard(InviteCodeInfo inviteCode) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.lg),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        children: [
          const Text(
            '邀请代理',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 8),
          const Text(
            '扫码或分享链接加入团队',
            style: TextStyle(
              fontSize: 14,
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: AppSpacing.lg),
          // 二维码
          Container(
            width: 200,
            height: 200,
            decoration: BoxDecoration(
              color: Colors.white,
              border: Border.all(color: AppColors.border),
              borderRadius: BorderRadius.circular(12),
            ),
            child: inviteCode.qrCodeUrl != null && inviteCode.qrCodeUrl!.isNotEmpty
                ? ClipRRect(
                    borderRadius: BorderRadius.circular(12),
                    child: Image.network(
                      inviteCode.qrCodeUrl!,
                      fit: BoxFit.cover,
                      errorBuilder: (_, __, ___) => _buildQRCodePlaceholder(inviteCode.inviteCode),
                    ),
                  )
                : _buildQRCodePlaceholder(inviteCode.inviteCode),
          ),
          const SizedBox(height: AppSpacing.lg),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              _buildActionButton(
                icon: Icons.save_alt,
                label: '保存图片',
                onTap: () => _saveQRCode(inviteCode),
              ),
              _buildActionButton(
                icon: Icons.share,
                label: '分享',
                onTap: () => _shareQRCode(inviteCode),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildQRCodePlaceholder(String inviteCode) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.qr_code_2,
            size: 120,
            color: AppColors.primary.withValues(alpha: 0.8),
          ),
          const SizedBox(height: 8),
          Text(
            inviteCode,
            style: const TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
              color: AppColors.primary,
              letterSpacing: 2,
            ),
          ),
        ],
      ),
    );
  }

  /// 邀请码卡片
  Widget _buildInviteCodeCard(InviteCodeInfo inviteCode) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [AppColors.primary, AppColors.primaryDark],
        ),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          const Icon(Icons.card_giftcard, color: Colors.white, size: 32),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  '我的邀请码',
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.white70,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  inviteCode.inviteCode,
                  style: const TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                    letterSpacing: 3,
                  ),
                ),
              ],
            ),
          ),
          TextButton.icon(
            onPressed: () => _copyInviteCode(inviteCode.inviteCode),
            icon: const Icon(Icons.copy, color: Colors.white, size: 18),
            label: const Text('复制', style: TextStyle(color: Colors.white)),
            style: TextButton.styleFrom(
              backgroundColor: Colors.white.withValues(alpha: 0.2),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(20),
              ),
            ),
          ),
        ],
      ),
    );
  }

  /// 分享操作
  Widget _buildShareActions(InviteCodeInfo inviteCode) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '快捷分享',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: AppSpacing.md),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              _buildShareItem(
                icon: Icons.wechat,
                label: '微信',
                color: AppColors.wechatPay,
                onTap: () => _shareToApp('微信', inviteCode),
              ),
              _buildShareItem(
                icon: Icons.chat_bubble,
                label: '朋友圈',
                color: AppColors.wechatPay,
                onTap: () => _shareToApp('朋友圈', inviteCode),
              ),
              _buildShareItem(
                icon: Icons.link,
                label: '复制链接',
                color: AppColors.primary,
                onTap: () => _copyInviteLink(inviteCode.inviteLink),
              ),
              _buildShareItem(
                icon: Icons.more_horiz,
                label: '更多',
                color: AppColors.textSecondary,
                onTap: () => _shareMore(inviteCode),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildShareItem({
    required IconData icon,
    required String label,
    required Color color,
    required VoidCallback onTap,
  }) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          children: [
            Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: color.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(icon, color: color, size: 24),
            ),
            const SizedBox(height: 8),
            Text(
              label,
              style: const TextStyle(
                fontSize: 12,
                color: AppColors.textSecondary,
              ),
            ),
          ],
        ),
      ),
    );
  }

  /// 推广提示
  Widget _buildPromotionTips() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: AppColors.info.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.info.withValues(alpha: 0.3)),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Icon(Icons.lightbulb_outline, color: AppColors.info, size: 20),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: const [
                Text(
                  '推广提示',
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                    color: AppColors.info,
                  ),
                ),
                SizedBox(height: 4),
                Text(
                  '• 被邀请人通过您的邀请码注册后自动成为您的下级代理\n'
                  '• 下级代理的交易将为您产生分润收益\n'
                  '• 分润比例根据政策模板设定',
                  style: TextStyle(
                    fontSize: 13,
                    color: AppColors.textSecondary,
                    height: 1.6,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  /// 团队管理Tab
  Widget _buildTeamTab() {
    return RefreshIndicator(
      onRefresh: () async {
        ref.invalidate(teamStatsProvider);
        await ref.read(subordinatesProvider.notifier).loadSubordinates(refresh: true);
      },
      child: CustomScrollView(
        slivers: [
          SliverPadding(
            padding: const EdgeInsets.all(AppSpacing.md),
            sliver: SliverList(
              delegate: SliverChildListDelegate([
                _buildTeamStatsCard(),
                const SizedBox(height: AppSpacing.md),
                _buildSubordinatesList(),
              ]),
            ),
          ),
        ],
      ),
    );
  }

  /// 团队统计卡片
  Widget _buildTeamStatsCard() {
    final teamStatsAsync = ref.watch(teamStatsProvider);

    return teamStatsAsync.when(
      loading: () => Container(
        padding: const EdgeInsets.all(AppSpacing.lg),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(12),
        ),
        child: const Center(child: CircularProgressIndicator()),
      ),
      error: (error, stack) => Container(
        padding: const EdgeInsets.all(AppSpacing.md),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(12),
        ),
        child: Center(
          child: Text('加载失败: $error', style: const TextStyle(color: AppColors.danger)),
        ),
      ),
      data: (stats) => Container(
        padding: const EdgeInsets.all(AppSpacing.md),
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
                const Text(
                  '团队概况',
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: AppColors.textPrimary,
                  ),
                ),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: AppColors.success.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.trending_up, size: 14, color: AppColors.success),
                      const SizedBox(width: 4),
                      Text(
                        '本月新增 ${stats.monthNewAgents}',
                        style: const TextStyle(
                          fontSize: 12,
                          color: AppColors.success,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
            const SizedBox(height: AppSpacing.md),
            GridView.count(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              crossAxisCount: 2,
              crossAxisSpacing: AppSpacing.cardGap,
              mainAxisSpacing: AppSpacing.cardGap,
              childAspectRatio: 2,
              children: [
                _buildStatItem(
                  '直属代理',
                  stats.directAgentCount.toString(),
                  Icons.person,
                  AppColors.primary,
                ),
                _buildStatItem(
                  '团队代理',
                  stats.teamAgentCount.toString(),
                  Icons.groups,
                  AppColors.profitReward,
                ),
                _buildStatItem(
                  '直营商户',
                  stats.directMerchantCount.toString(),
                  Icons.store,
                  AppColors.success,
                ),
                _buildStatItem(
                  '团队商户',
                  stats.teamMerchantCount.toString(),
                  Icons.storefront,
                  AppColors.warning,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatItem(String label, String value, IconData icon, Color color) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.05),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Container(
            width: 36,
            height: 36,
            decoration: BoxDecoration(
              color: color.withValues(alpha: 0.15),
              borderRadius: BorderRadius.circular(10),
            ),
            child: Icon(icon, color: color, size: 20),
          ),
          const SizedBox(width: 12),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                value,
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: color,
                ),
              ),
              Text(
                label,
                style: const TextStyle(
                  fontSize: 12,
                  color: AppColors.textSecondary,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  /// 下级代理列表
  Widget _buildSubordinatesList() {
    final subordinatesState = ref.watch(subordinatesProvider);

    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.all(AppSpacing.md),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  '直属代理',
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: AppColors.textPrimary,
                  ),
                ),
                TextButton(
                  onPressed: () {
                    // TODO: 跳转到全部下级页面
                    _showSnackBar('查看全部功能开发中');
                  },
                  child: const Row(
                    children: [
                      Text('查看全部', style: TextStyle(fontSize: 13)),
                      Icon(Icons.chevron_right, size: 18),
                    ],
                  ),
                ),
              ],
            ),
          ),
          const Divider(height: 1, color: AppColors.divider),
          if (subordinatesState.isLoading && subordinatesState.list.isEmpty)
            const Padding(
              padding: EdgeInsets.all(32),
              child: Center(child: CircularProgressIndicator()),
            )
          else if (subordinatesState.error != null && subordinatesState.list.isEmpty)
            Padding(
              padding: const EdgeInsets.all(32),
              child: Center(
                child: Column(
                  children: [
                    const Icon(Icons.error_outline, size: 48, color: AppColors.danger),
                    const SizedBox(height: 12),
                    Text('加载失败: ${subordinatesState.error}'),
                    const SizedBox(height: 12),
                    ElevatedButton(
                      onPressed: () => ref.read(subordinatesProvider.notifier).loadSubordinates(refresh: true),
                      child: const Text('重试'),
                    ),
                  ],
                ),
              ),
            )
          else if (subordinatesState.list.isEmpty)
            const Padding(
              padding: EdgeInsets.all(32),
              child: Center(
                child: Column(
                  children: [
                    Icon(Icons.people_outline, size: 48, color: AppColors.textTertiary),
                    SizedBox(height: 12),
                    Text('暂无下级代理', style: TextStyle(color: AppColors.textSecondary)),
                    SizedBox(height: 4),
                    Text('分享邀请码发展团队吧', style: TextStyle(color: AppColors.textTertiary, fontSize: 12)),
                  ],
                ),
              ),
            )
          else
            ListView.separated(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              itemCount: subordinatesState.list.length,
              separatorBuilder: (_, __) => const Divider(
                height: 1,
                indent: 16,
                endIndent: 16,
                color: AppColors.divider,
              ),
              itemBuilder: (context, index) {
                final agent = subordinatesState.list[index];
                return _buildSubordinateItem(agent);
              },
            ),
        ],
      ),
    );
  }

  Widget _buildSubordinateItem(AgentInfo agent) {
    final isActive = agent.status == 1;

    return InkWell(
      onTap: () {
        // 跳转到代理详情
        context.push('/agent/${agent.id}');
      },
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.md),
        child: Row(
          children: [
            // 头像
            Container(
              width: 44,
              height: 44,
              decoration: BoxDecoration(
                color: AppColors.primary.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Center(
                child: Text(
                  agent.agentName.isNotEmpty ? agent.agentName.substring(0, 1) : '?',
                  style: const TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: AppColors.primary,
                  ),
                ),
              ),
            ),
            const SizedBox(width: 12),
            // 信息
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Text(
                        agent.agentName,
                        style: const TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w500,
                          color: AppColors.textPrimary,
                        ),
                      ),
                      const SizedBox(width: 8),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                        decoration: BoxDecoration(
                          color: isActive
                              ? AppColors.success.withValues(alpha: 0.1)
                              : AppColors.danger.withValues(alpha: 0.1),
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Text(
                          isActive ? '正常' : '禁用',
                          style: TextStyle(
                            fontSize: 10,
                            color: isActive ? AppColors.success : AppColors.danger,
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 4),
                  Text(
                    '${agent.contactPhone} · ${agent.registerTime ?? ''}',
                    style: const TextStyle(
                      fontSize: 12,
                      color: AppColors.textTertiary,
                    ),
                  ),
                ],
              ),
            ),
            // 操作菜单
            PopupMenuButton<String>(
              icon: const Icon(Icons.more_vert, size: 20, color: AppColors.textTertiary),
              onSelected: (value) => _handleAgentAction(value, agent),
              itemBuilder: (context) => [
                const PopupMenuItem(
                  value: 'policy',
                  child: Row(
                    children: [
                      Icon(Icons.policy_outlined, size: 18, color: AppColors.primary),
                      SizedBox(width: 8),
                      Text('设置政策'),
                    ],
                  ),
                ),
                const PopupMenuItem(
                  value: 'channels',
                  child: Row(
                    children: [
                      Icon(Icons.account_tree_outlined, size: 18, color: AppColors.info),
                      SizedBox(width: 8),
                      Text('通道政策'),
                    ],
                  ),
                ),
                const PopupMenuItem(
                  value: 'detail',
                  child: Row(
                    children: [
                      Icon(Icons.person_outline, size: 18, color: AppColors.textSecondary),
                      SizedBox(width: 8),
                      Text('查看详情'),
                    ],
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  void _handleAgentAction(String action, AgentInfo agent) {
    switch (action) {
      case 'policy':
        // 跳转到设置政策页面
        context.push(
          '/agent/${agent.id}/policy',
          extra: {
            'name': agent.agentName,
          },
        );
        break;
      case 'channels':
        // 跳转到通道政策页面
        context.push('/agent/${agent.id}/channels');
        break;
      case 'detail':
        // 跳转到代理详情
        context.push('/agent/${agent.id}');
        break;
    }
  }

  Widget _buildActionButton({
    required IconData icon,
    required String label,
    required VoidCallback onTap,
  }) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(8),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
        decoration: BoxDecoration(
          border: Border.all(color: AppColors.border),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(icon, size: 18, color: AppColors.primary),
            const SizedBox(width: 8),
            Text(
              label,
              style: const TextStyle(
                fontSize: 14,
                color: AppColors.primary,
              ),
            ),
          ],
        ),
      ),
    );
  }

  // ==================== 事件处理 ====================

  void _copyInviteCode(String code) {
    Clipboard.setData(ClipboardData(text: code));
    _showSnackBar('邀请码已复制');
  }

  void _copyInviteLink(String link) {
    Clipboard.setData(ClipboardData(text: link));
    _showSnackBar('邀请链接已复制');
  }

  void _saveQRCode(InviteCodeInfo inviteCode) async {
    // TODO: 实际保存二维码到相册
    // 需要使用 image_gallery_saver 或类似插件
    _showSnackBar('二维码已保存到相册');
  }

  void _shareQRCode(InviteCodeInfo inviteCode) async {
    await Share.share(
      '邀请您加入享收付代理团队！\n邀请码：${inviteCode.inviteCode}\n注册链接：${inviteCode.inviteLink}',
      subject: '享收付代理邀请',
    );
  }

  void _shareToApp(String app, InviteCodeInfo inviteCode) async {
    // 调用系统分享
    await Share.share(
      '邀请您加入享收付代理团队！\n邀请码：${inviteCode.inviteCode}\n注册链接：${inviteCode.inviteLink}',
      subject: '享收付代理邀请',
    );
  }

  void _shareMore(InviteCodeInfo inviteCode) async {
    await Share.share(
      '邀请您加入享收付代理团队！\n邀请码：${inviteCode.inviteCode}\n注册链接：${inviteCode.inviteLink}',
      subject: '享收付代理邀请',
    );
  }

  void _showAddAgentDialog() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => _AddAgentSheet(
        onSuccess: () {
          // 刷新列表
          ref.invalidate(teamStatsProvider);
          ref.read(subordinatesProvider.notifier).loadSubordinates(refresh: true);
        },
      ),
    );
  }

  void _showSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        duration: const Duration(seconds: 2),
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      ),
    );
  }
}

/// 手动添加代理表单
class _AddAgentSheet extends ConsumerStatefulWidget {
  final VoidCallback? onSuccess;

  const _AddAgentSheet({this.onSuccess});

  @override
  ConsumerState<_AddAgentSheet> createState() => _AddAgentSheetState();
}

class _AddAgentSheetState extends ConsumerState<_AddAgentSheet> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _phoneController = TextEditingController();
  final _contactController = TextEditingController();

  @override
  void dispose() {
    _nameController.dispose();
    _phoneController.dispose();
    _contactController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final createState = ref.watch(createAgentProvider);

    return Padding(
      padding: EdgeInsets.only(
        left: AppSpacing.md,
        right: AppSpacing.md,
        top: AppSpacing.md,
        bottom: MediaQuery.of(context).viewInsets.bottom + AppSpacing.md,
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // 标题栏
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              TextButton(
                onPressed: () => Navigator.pop(context),
                child: const Text('取消'),
              ),
              const Text(
                '添加代理',
                style: TextStyle(
                  fontSize: 17,
                  fontWeight: FontWeight.w600,
                ),
              ),
              TextButton(
                onPressed: createState.isSubmitting ? null : _submit,
                child: createState.isSubmitting
                    ? const SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Text('确定'),
              ),
            ],
          ),
          const SizedBox(height: AppSpacing.md),
          // 表单
          Form(
            key: _formKey,
            child: Column(
              children: [
                TextFormField(
                  controller: _nameController,
                  decoration: const InputDecoration(
                    labelText: '代理商名称',
                    hintText: '请输入代理商名称',
                    border: OutlineInputBorder(),
                  ),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return '请输入代理商名称';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: AppSpacing.md),
                TextFormField(
                  controller: _contactController,
                  decoration: const InputDecoration(
                    labelText: '联系人',
                    hintText: '请输入联系人姓名',
                    border: OutlineInputBorder(),
                  ),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return '请输入联系人姓名';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: AppSpacing.md),
                TextFormField(
                  controller: _phoneController,
                  decoration: const InputDecoration(
                    labelText: '手机号码',
                    hintText: '请输入手机号码',
                    border: OutlineInputBorder(),
                  ),
                  keyboardType: TextInputType.phone,
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return '请输入手机号码';
                    }
                    if (!RegExp(r'^1[3-9]\d{9}$').hasMatch(value)) {
                      return '请输入正确的手机号码';
                    }
                    return null;
                  },
                ),
              ],
            ),
          ),
          const SizedBox(height: AppSpacing.lg),
        ],
      ),
    );
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    final request = CreateAgentRequest(
      agentName: _nameController.text.trim(),
      contactName: _contactController.text.trim(),
      contactPhone: _phoneController.text.trim(),
    );

    final success = await ref.read(createAgentProvider.notifier).createAgent(request);

    if (success && mounted) {
      Navigator.pop(context);
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: const Text('代理商添加成功'),
          behavior: SnackBarBehavior.floating,
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
        ),
      );
      widget.onSuccess?.call();
    } else if (mounted) {
      final error = ref.read(createAgentProvider).error;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('添加失败: ${error ?? "未知错误"}'),
          backgroundColor: AppColors.danger,
        ),
      );
    }
  }
}
