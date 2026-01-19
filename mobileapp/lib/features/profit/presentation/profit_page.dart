import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';

/// 收益统计页面
class ProfitPage extends StatelessWidget {
  const ProfitPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('收益统计')),
      body: const Center(child: Text('收益统计页面 - 待实现')),
    );
  }
}
