import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';

/// 商户管理页面
class MerchantPage extends StatelessWidget {
  const MerchantPage({super.key});

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 2,
      child: Scaffold(
        backgroundColor: AppColors.background,
        appBar: AppBar(
          title: const Text('商户管理'),
          bottom: const TabBar(
            tabs: [
              Tab(text: '直营'),
              Tab(text: '团队'),
            ],
          ),
        ),
        body: const TabBarView(
          children: [
            Center(child: Text('直营商户列表')),
            Center(child: Text('团队商户列表')),
          ],
        ),
      ),
    );
  }
}
