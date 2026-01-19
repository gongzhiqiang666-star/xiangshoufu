import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';

/// 数据分析页面
class DataAnalysisPage extends StatelessWidget {
  const DataAnalysisPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('数据分析')),
      body: const Center(child: Text('数据分析页面 - 待实现')),
    );
  }
}
