import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';

/// 消息分组头部
class MessageGroupHeader extends StatelessWidget {
  final String title;

  const MessageGroupHeader({
    super.key,
    required this.title,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      color: AppColors.background,
      child: Text(
        title,
        style: const TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.w500,
          color: AppColors.textTertiary,
        ),
      ),
    );
  }
}
