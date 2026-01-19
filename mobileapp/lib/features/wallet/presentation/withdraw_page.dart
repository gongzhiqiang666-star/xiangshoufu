import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';

/// 提现页面
class WithdrawPage extends StatelessWidget {
  final String walletId;
  
  const WithdrawPage({super.key, required this.walletId});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('申请提现')),
      body: const Center(child: Text('提现页面 - 待实现')),
    );
  }
}
