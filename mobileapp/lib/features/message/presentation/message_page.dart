import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';

/// 消息通知页面
class MessagePage extends StatelessWidget {
  const MessagePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('消息通知')),
      body: const Center(child: Text('消息通知页面 - 待实现')),
    );
  }
}
