import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../router/app_router.dart';

/// 设置页面
class SettingsPage extends ConsumerWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('设置'),
      ),
      body: ListView(
        children: [
          const SizedBox(height: AppSpacing.md),

          // 账户安全
          _buildSectionTitle('账户安全'),
          _buildSettingCard([
            _buildSettingItem(
              context,
              icon: Icons.lock_outline,
              title: '修改密码',
              onTap: () => _showChangePasswordDialog(context),
            ),
            _buildSettingItem(
              context,
              icon: Icons.phone_android,
              title: '修改手机号',
              subtitle: '需要验证原手机号',
              onTap: () => _showToast(context, '功能开发中...'),
            ),
          ]),

          const SizedBox(height: AppSpacing.md),

          // 结算设置
          _buildSectionTitle('结算设置'),
          _buildSettingCard([
            _buildSettingItem(
              context,
              icon: Icons.credit_card,
              title: '结算卡管理',
              onTap: () => context.push(RoutePaths.bankCard),
            ),
          ]),

          const SizedBox(height: AppSpacing.md),

          // 通用设置
          _buildSectionTitle('通用设置'),
          _buildSettingCard([
            _buildSettingItem(
              context,
              icon: Icons.notifications_outlined,
              title: '消息通知',
              onTap: () => _showToast(context, '功能开发中...'),
            ),
            _buildSettingItem(
              context,
              icon: Icons.language,
              title: '语言设置',
              subtitle: '简体中文',
              onTap: () => _showToast(context, '功能开发中...'),
            ),
            _buildSettingItem(
              context,
              icon: Icons.cleaning_services_outlined,
              title: '清除缓存',
              subtitle: '12.5MB',
              onTap: () => _showClearCacheDialog(context),
            ),
          ]),

          const SizedBox(height: AppSpacing.md),

          // 关于
          _buildSectionTitle('关于'),
          _buildSettingCard([
            _buildSettingItem(
              context,
              icon: Icons.info_outline,
              title: '关于我们',
              onTap: () => _showToast(context, '功能开发中...'),
            ),
            _buildSettingItem(
              context,
              icon: Icons.article_outlined,
              title: '用户协议',
              onTap: () => _showToast(context, '功能开发中...'),
            ),
            _buildSettingItem(
              context,
              icon: Icons.privacy_tip_outlined,
              title: '隐私政策',
              onTap: () => _showToast(context, '功能开发中...'),
            ),
            _buildSettingItem(
              context,
              icon: Icons.system_update_outlined,
              title: '版本更新',
              subtitle: 'v1.0.0',
              onTap: () => _showToast(context, '已是最新版本'),
              showArrow: false,
            ),
          ]),

          const SizedBox(height: AppSpacing.xl),

          // 退出登录按钮
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
            child: ElevatedButton(
              onPressed: () => _showLogoutDialog(context, ref),
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.white,
                foregroundColor: AppColors.danger,
                padding: const EdgeInsets.symmetric(vertical: AppSpacing.md),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                  side: const BorderSide(color: AppColors.danger),
                ),
              ),
              child: const Text(
                '退出登录',
                style: TextStyle(fontSize: 16, fontWeight: FontWeight.w500),
              ),
            ),
          ),

          const SizedBox(height: AppSpacing.xl),
        ],
      ),
    );
  }

  Widget _buildSectionTitle(String title) {
    return Padding(
      padding: const EdgeInsets.symmetric(
        horizontal: AppSpacing.md,
        vertical: AppSpacing.sm,
      ),
      child: Text(
        title,
        style: const TextStyle(
          fontSize: 13,
          color: AppColors.textSecondary,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }

  Widget _buildSettingCard(List<Widget> children) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        children: children,
      ),
    );
  }

  Widget _buildSettingItem(
    BuildContext context, {
    required IconData icon,
    required String title,
    String? subtitle,
    required VoidCallback onTap,
    bool showArrow = true,
  }) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Padding(
        padding: const EdgeInsets.symmetric(
          horizontal: AppSpacing.md,
          vertical: AppSpacing.md,
        ),
        child: Row(
          children: [
            Container(
              width: 36,
              height: 36,
              decoration: BoxDecoration(
                color: AppColors.background,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Icon(icon, color: AppColors.primary, size: 20),
            ),
            const SizedBox(width: AppSpacing.md),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: const TextStyle(
                      fontSize: 15,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  if (subtitle != null)
                    Padding(
                      padding: const EdgeInsets.only(top: 2),
                      child: Text(
                        subtitle,
                        style: const TextStyle(
                          fontSize: 12,
                          color: AppColors.textTertiary,
                        ),
                      ),
                    ),
                ],
              ),
            ),
            if (showArrow)
              const Icon(
                Icons.chevron_right,
                color: AppColors.textTertiary,
                size: 20,
              ),
          ],
        ),
      ),
    );
  }

  void _showToast(BuildContext context, String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message)),
    );
  }

  void _showChangePasswordDialog(BuildContext context) {
    final oldPasswordController = TextEditingController();
    final newPasswordController = TextEditingController();
    final confirmPasswordController = TextEditingController();

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('修改密码'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: oldPasswordController,
              obscureText: true,
              decoration: const InputDecoration(
                labelText: '当前密码',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: AppSpacing.md),
            TextField(
              controller: newPasswordController,
              obscureText: true,
              decoration: const InputDecoration(
                labelText: '新密码',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: AppSpacing.md),
            TextField(
              controller: confirmPasswordController,
              obscureText: true,
              decoration: const InputDecoration(
                labelText: '确认新密码',
                border: OutlineInputBorder(),
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
              // TODO: 实现修改密码逻辑
              Navigator.pop(context);
              _showToast(context, '密码修改成功');
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }

  void _showClearCacheDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('清除缓存'),
        content: const Text('确定要清除本地缓存吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              _showToast(context, '缓存已清除');
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }

  void _showLogoutDialog(BuildContext context, WidgetRef ref) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('退出登录'),
        content: const Text('确定要退出登录吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              // TODO: 清除登录状态
              context.go(RoutePaths.login);
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.danger,
            ),
            child: const Text('退出'),
          ),
        ],
      ),
    );
  }
}
