import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';

/// 我的信息页面
class ProfilePage extends StatelessWidget {
  const ProfilePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('我的')),
      body: const Center(child: Text('我的信息页面 - 待实现')),
    );
  }
}
