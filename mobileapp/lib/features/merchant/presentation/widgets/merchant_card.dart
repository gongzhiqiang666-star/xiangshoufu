import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import '../../data/models/merchant_model.dart';
import '../../../../core/theme/app_colors.dart';

/// 商户卡片组件
class MerchantCard extends StatelessWidget {
  final Merchant merchant;
  final VoidCallback? onTap;

  const MerchantCard({
    super.key,
    required this.merchant,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: EdgeInsets.symmetric(horizontal: 16.w, vertical: 6.h),
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12.r),
        side: BorderSide(color: Colors.grey.shade200),
      ),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12.r),
        child: Padding(
          padding: EdgeInsets.all(16.w),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // 顶部：商户名称和类型标签
              Row(
                children: [
                  Expanded(
                    child: Text(
                      merchant.merchantName,
                      style: TextStyle(
                        fontSize: 16.sp,
                        fontWeight: FontWeight.w600,
                        color: AppColors.textPrimary,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  SizedBox(width: 8.w),
                  _buildMerchantTypeTag(),
                ],
              ),
              SizedBox(height: 12.h),
              // 商户编号
              _buildInfoRow('商户编号', merchant.merchantNo),
              SizedBox(height: 6.h),
              // 终端SN
              _buildInfoRow('终端SN', merchant.terminalSn ?? '-'),
              SizedBox(height: 6.h),
              // 归属类型和状态
              Row(
                children: [
                  _buildOwnerTypeTag(),
                  SizedBox(width: 8.w),
                  _buildStatusTag(),
                  const Spacer(),
                  Icon(
                    Icons.arrow_forward_ios,
                    size: 14.sp,
                    color: Colors.grey,
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Row(
      children: [
        Text(
          '$label: ',
          style: TextStyle(
            fontSize: 13.sp,
            color: AppColors.textSecondary,
          ),
        ),
        Expanded(
          child: Text(
            value,
            style: TextStyle(
              fontSize: 13.sp,
              color: AppColors.textPrimary,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ],
    );
  }

  Widget _buildMerchantTypeTag() {
    Color bgColor;
    Color textColor;
    // 5档分类颜色映射：优质/中等/普通/预警/流失
    switch (merchant.merchantType) {
      case 'quality':
        bgColor = Colors.green.shade50;
        textColor = Colors.green.shade700;
        break;
      case 'medium':
        bgColor = Colors.blue.shade50;
        textColor = Colors.blue.shade700;
        break;
      case 'normal':
        bgColor = Colors.grey.shade50;
        textColor = Colors.grey.shade700;
        break;
      case 'warning':
        bgColor = Colors.orange.shade50;
        textColor = Colors.orange.shade700;
        break;
      case 'churned':
        bgColor = Colors.red.shade50;
        textColor = Colors.red.shade700;
        break;
      default:
        bgColor = Colors.grey.shade50;
        textColor = Colors.grey.shade700;
    }

    return Container(
      padding: EdgeInsets.symmetric(horizontal: 8.w, vertical: 4.h),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(4.r),
      ),
      child: Text(
        merchant.merchantTypeName,
        style: TextStyle(
          fontSize: 11.sp,
          color: textColor,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }

  Widget _buildOwnerTypeTag() {
    final isDirect = merchant.isDirect;
    return Container(
      padding: EdgeInsets.symmetric(horizontal: 8.w, vertical: 4.h),
      decoration: BoxDecoration(
        color: isDirect ? AppColors.primary.withValues(alpha: 0.1) : Colors.purple.shade50,
        borderRadius: BorderRadius.circular(4.r),
      ),
      child: Text(
        merchant.ownerTypeName,
        style: TextStyle(
          fontSize: 11.sp,
          color: isDirect ? AppColors.primary : Colors.purple.shade700,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }

  Widget _buildStatusTag() {
    final isActive = merchant.status == 1;
    return Container(
      padding: EdgeInsets.symmetric(horizontal: 8.w, vertical: 4.h),
      decoration: BoxDecoration(
        color: isActive ? Colors.green.shade50 : Colors.red.shade50,
        borderRadius: BorderRadius.circular(4.r),
      ),
      child: Text(
        merchant.statusName ?? (isActive ? '正常' : '禁用'),
        style: TextStyle(
          fontSize: 11.sp,
          color: isActive ? Colors.green.shade700 : Colors.red.shade700,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }
}
