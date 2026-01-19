import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';

/// 代理拓展页面
class AgentPage extends StatelessWidget {
  const AgentPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('代理拓展')),
      body: const Center(child: Text('代理拓展页面 - 待实现')),
    );
  }
}
