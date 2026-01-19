import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';

/// 商户详情页面
class MerchantDetailPage extends StatelessWidget {
  final String merchantId;
  
  const MerchantDetailPage({super.key, required this.merchantId});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('商户详情')),
      body: const Center(child: Text('商户详情页面')),
    );
  }
}
