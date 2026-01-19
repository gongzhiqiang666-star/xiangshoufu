import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';

/// 代扣管理页面
class DeductionPage extends StatelessWidget {
  const DeductionPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('代扣管理')),
      body: const Center(child: Text('代扣管理页面 - 待实现')),
    );
  }
}
