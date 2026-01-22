import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import '../../../core/theme/app_colors.dart';
import '../data/models/merchant_model.dart';
import 'providers/merchant_provider.dart';

/// 商户详情页面
class MerchantDetailPage extends ConsumerWidget {
  final int merchantId;

  const MerchantDetailPage({super.key, required this.merchantId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final detailAsync = ref.watch(merchantDetailProvider(merchantId));

    return Scaffold(
      backgroundColor: Colors.grey.shade50,
      appBar: AppBar(
        title: const Text('商户详情'),
        centerTitle: true,
        elevation: 0,
      ),
      body: detailAsync.when(
        data: (detail) => _buildContent(context, detail),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, _) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(Icons.error_outline, size: 48.sp, color: Colors.grey),
              SizedBox(height: 16.h),
              Text('加载失败', style: TextStyle(fontSize: 14.sp, color: Colors.grey)),
              SizedBox(height: 16.h),
              ElevatedButton(
                onPressed: () => ref.refresh(merchantDetailProvider(merchantId)),
                child: const Text('重试'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildContent(BuildContext context, MerchantDetail detail) {
    return SingleChildScrollView(
      padding: EdgeInsets.all(16.w),
      child: Column(
        children: [
          // 商户头部信息
          _buildHeaderCard(detail),
          SizedBox(height: 16.h),
          // 基本信息
          _buildInfoCard(detail),
          SizedBox(height: 16.h),
          // 费率信息
          _buildRateCard(detail),
          SizedBox(height: 16.h),
          // 统计信息
          _buildStatsCard(detail),
        ],
      ),
    );
  }

  Widget _buildHeaderCard(MerchantDetail detail) {
    return Container(
      padding: EdgeInsets.all(20.w),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [AppColors.primary, AppColors.primary.withValues(alpha: 0.8)],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(16.r),
        boxShadow: [
          BoxShadow(
            color: AppColors.primary.withValues(alpha: 0.3),
            blurRadius: 10,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Column(
        children: [
          // 商户头像
          Container(
            width: 64.w,
            height: 64.w,
            decoration: BoxDecoration(
              color: Colors.white.withValues(alpha: 0.2),
              shape: BoxShape.circle,
            ),
            child: Icon(Icons.store, size: 32.sp, color: Colors.white),
          ),
          SizedBox(height: 12.h),
          // 商户名称
          Text(
            detail.merchantName,
            style: TextStyle(
              fontSize: 20.sp,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          SizedBox(height: 4.h),
          // 商户编号
          Text(
            detail.merchantNo,
            style: TextStyle(
              fontSize: 14.sp,
              color: Colors.white.withValues(alpha: 0.8),
            ),
          ),
          SizedBox(height: 12.h),
          // 状态标签
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              _buildHeaderTag(
                detail.status == 1 ? '正常' : '禁用',
                detail.status == 1 ? Colors.green : Colors.red,
              ),
              SizedBox(width: 8.w),
              _buildHeaderTag(
                detail.isDirect ? '直营' : '团队',
                detail.isDirect ? Colors.blue : Colors.purple,
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildHeaderTag(String text, Color color) {
    return Container(
      padding: EdgeInsets.symmetric(horizontal: 12.w, vertical: 4.h),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12.r),
      ),
      child: Text(
        text,
        style: TextStyle(
          fontSize: 12.sp,
          fontWeight: FontWeight.w500,
          color: color,
        ),
      ),
    );
  }

  Widget _buildInfoCard(MerchantDetail detail) {
    return _buildCard(
      title: '基本信息',
      child: Column(
        children: [
          _buildInfoRow('所属代理', detail.agentName ?? '-'),
          _buildInfoRow('代理层级', detail.agentLevel != null ? '${detail.agentLevel}级' : '-'),
          _buildInfoRow('所属通道', detail.channelName ?? '-'),
          _buildInfoRow('终端SN', detail.terminalSn ?? '-'),
          _buildInfoRow('MCC码', detail.mcc ?? '-'),
          _buildInfoRow('法人姓名', detail.legalName ?? '-'),
          _buildInfoRow('激活时间', detail.activatedAt ?? '-'),
          _buildInfoRow('创建时间', detail.createdAt ?? '-'),
        ],
      ),
    );
  }

  Widget _buildRateCard(MerchantDetail detail) {
    return _buildCard(
      title: '费率信息',
      child: Row(
        children: [
          Expanded(
            child: _buildRateItem(
              '贷记卡费率',
              _formatRate(detail.creditRate),
              Colors.blue,
            ),
          ),
          SizedBox(width: 16.w),
          Expanded(
            child: _buildRateItem(
              '借记卡费率',
              _formatRate(detail.debitRate),
              Colors.green,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildRateItem(String label, String value, Color color) {
    return Container(
      padding: EdgeInsets.all(16.w),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(12.r),
      ),
      child: Column(
        children: [
          Text(
            value,
            style: TextStyle(
              fontSize: 24.sp,
              fontWeight: FontWeight.bold,
              color: color,
            ),
          ),
          SizedBox(height: 4.h),
          Text(
            label,
            style: TextStyle(
              fontSize: 12.sp,
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatsCard(MerchantDetail detail) {
    return _buildCard(
      title: '交易统计',
      child: Row(
        children: [
          Expanded(
            child: _buildStatItem(
              '本月交易额',
              '¥${detail.monthAmountFormatted}',
              Icons.payments_outlined,
              Colors.orange,
            ),
          ),
          SizedBox(width: 12.w),
          Expanded(
            child: _buildStatItem(
              '本月笔数',
              '${detail.monthCount ?? 0}',
              Icons.receipt_long_outlined,
              Colors.purple,
            ),
          ),
          SizedBox(width: 12.w),
          Expanded(
            child: _buildStatItem(
              '终端数',
              '${detail.terminalCount ?? 0}',
              Icons.devices_outlined,
              Colors.teal,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatItem(String label, String value, IconData icon, Color color) {
    return Container(
      padding: EdgeInsets.all(12.w),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(12.r),
      ),
      child: Column(
        children: [
          Icon(icon, size: 24.sp, color: color),
          SizedBox(height: 8.h),
          Text(
            value,
            style: TextStyle(
              fontSize: 16.sp,
              fontWeight: FontWeight.bold,
              color: AppColors.textPrimary,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          SizedBox(height: 4.h),
          Text(
            label,
            style: TextStyle(
              fontSize: 11.sp,
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCard({required String title, required Widget child}) {
    return Container(
      width: double.infinity,
      padding: EdgeInsets.all(16.w),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12.r),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: TextStyle(
              fontSize: 16.sp,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          SizedBox(height: 16.h),
          child,
        ],
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: EdgeInsets.symmetric(vertical: 8.h),
      child: Row(
        children: [
          SizedBox(
            width: 80.w,
            child: Text(
              label,
              style: TextStyle(
                fontSize: 14.sp,
                color: AppColors.textSecondary,
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: TextStyle(
                fontSize: 14.sp,
                color: AppColors.textPrimary,
              ),
            ),
          ),
        ],
      ),
    );
  }

  String _formatRate(String? rate) {
    if (rate == null || rate.isEmpty) return '0.00%';
    try {
      final num = double.parse(rate);
      return '${(num * 100).toStringAsFixed(2)}%';
    } catch (e) {
      return rate;
    }
  }
}
