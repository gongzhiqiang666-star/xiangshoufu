import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';
import '../../data/models/message_model.dart';

/// 消息列表项
class MessageListItem extends StatelessWidget {
  final MessageModel message;
  final VoidCallback? onTap;

  const MessageListItem({
    super.key,
    required this.message,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        decoration: BoxDecoration(
          color: message.isRead ? AppColors.cardBg : AppColors.cardBg,
          border: Border(
            bottom: BorderSide(color: AppColors.divider, width: 0.5),
          ),
        ),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 图标
            _buildIcon(),
            const SizedBox(width: 12),
            // 内容
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          message.title,
                          style: TextStyle(
                            fontSize: 15,
                            fontWeight: message.isRead ? FontWeight.normal : FontWeight.w600,
                            color: AppColors.textPrimary,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      const SizedBox(width: 8),
                      Text(
                        message.formattedTime,
                        style: const TextStyle(
                          fontSize: 12,
                          color: AppColors.textTertiary,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          message.content,
                          style: const TextStyle(
                            fontSize: 13,
                            color: AppColors.textSecondary,
                          ),
                          maxLines: 2,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      // 未读红点
                      if (!message.isRead) ...[
                        const SizedBox(width: 8),
                        Container(
                          width: 8,
                          height: 8,
                          decoration: const BoxDecoration(
                            color: AppColors.danger,
                            shape: BoxShape.circle,
                          ),
                        ),
                      ],
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildIcon() {
    IconData iconData;
    Color iconColor;
    Color bgColor;

    switch (message.type) {
      case MessageType.profit:
        iconData = Icons.trending_up;
        iconColor = AppColors.profitTrade;
        bgColor = AppColors.profitTrade.withOpacity(0.1);
        break;
      case MessageType.activation:
        iconData = Icons.star;
        iconColor = AppColors.profitReward;
        bgColor = AppColors.profitReward.withOpacity(0.1);
        break;
      case MessageType.deposit:
        iconData = Icons.account_balance_wallet;
        iconColor = AppColors.profitDeposit;
        bgColor = AppColors.profitDeposit.withOpacity(0.1);
        break;
      case MessageType.simCashback:
        iconData = Icons.sim_card;
        iconColor = AppColors.profitSim;
        bgColor = AppColors.profitSim.withOpacity(0.1);
        break;
      case MessageType.refund:
        iconData = Icons.undo;
        iconColor = AppColors.danger;
        bgColor = AppColors.danger.withOpacity(0.1);
        break;
      case MessageType.announcement:
        iconData = Icons.campaign;
        iconColor = AppColors.primary;
        bgColor = AppColors.primary.withOpacity(0.1);
        break;
      case MessageType.newAgent:
        iconData = Icons.person_add;
        iconColor = AppColors.success;
        bgColor = AppColors.success.withOpacity(0.1);
        break;
      case MessageType.transaction:
        iconData = Icons.receipt_long;
        iconColor = AppColors.warning;
        bgColor = AppColors.warning.withOpacity(0.1);
        break;
    }

    return Container(
      width: 40,
      height: 40,
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Icon(
        iconData,
        size: 20,
        color: iconColor,
      ),
    );
  }
}
