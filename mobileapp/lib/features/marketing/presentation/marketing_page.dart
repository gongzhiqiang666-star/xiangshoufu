import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';

/// 营销海报页面
class MarketingPage extends StatelessWidget {
  const MarketingPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('营销海报')),
      body: const Center(child: Text('营销海报页面 - 待实现')),
    );
  }
}
