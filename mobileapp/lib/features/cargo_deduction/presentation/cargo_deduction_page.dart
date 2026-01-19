import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';

/// 货款代扣页面
class CargoDeductionPage extends StatelessWidget {
  const CargoDeductionPage({super.key});

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 3,
      child: Scaffold(
        backgroundColor: AppColors.background,
        appBar: AppBar(
          title: const Text('货款代扣'),
          bottom: const TabBar(
            tabs: [
              Tab(text: '待接收'),
              Tab(text: '进行中'),
              Tab(text: '已完成'),
            ],
          ),
        ),
        body: const TabBarView(
          children: [
            Center(child: Text('待接收列表')),
            Center(child: Text('进行中列表')),
            Center(child: Text('已完成列表')),
          ],
        ),
      ),
    );
  }
}
