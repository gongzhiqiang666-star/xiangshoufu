import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../domain/models/terminal.dart';
import 'providers/terminal_provider.dart';

/// 终端详情页面
class TerminalDetailPage extends ConsumerStatefulWidget {
  final String terminalId; // 实际上是SN，为了路由参数命名统一暂时叫id

  const TerminalDetailPage({
    super.key,
    required this.terminalId,
  });

  @override
  ConsumerState<TerminalDetailPage> createState() => _TerminalDetailPageState();
}

class _TerminalDetailPageState extends ConsumerState<TerminalDetailPage> {
  @override
  Widget build(BuildContext context) {
    final terminalAsync = ref.watch(terminalDetailProvider(widget.terminalId));

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('终端详情')),
      body: terminalAsync.when(
        data: (terminal) => SingleChildScrollView(
          padding: const EdgeInsets.all(AppSpacing.md),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _buildHeaderCard(terminal),
              const SizedBox(height: AppSpacing.md),
              _buildInfoCard(terminal),
              const SizedBox(height: AppSpacing.md),
              _buildStatusCard(terminal),
              // 如果有商户信息，展示商户卡片
              if (terminal.merchantNo != null && terminal.merchantNo!.isNotEmpty) ...[
                const SizedBox(height: AppSpacing.md),
                _buildMerchantCard(terminal),
              ],
            ],
          ),
        ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.error_outline, size: 64, color: Colors.grey),
              const SizedBox(height: 16),
              Text('加载失败: $error'),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () => ref.refresh(terminalDetailProvider(widget.terminalId)),
                child: const Text('重试'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildHeaderCard(Terminal terminal) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.lg),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        children: [
          const Icon(Icons.point_of_sale, size: 48, color: AppColors.primary),
          const SizedBox(height: 12),
          Text(
            terminal.terminalSn,
            style: const TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 4),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: terminal.isActivated
                  ? AppColors.success.withValues(alpha: 0.1)
                  : AppColors.textTertiary.withValues(alpha: 0.1),
              borderRadius: BorderRadius.circular(4),
            ),
            child: Text(
              terminal.status.label,
              style: TextStyle(
                fontSize: 12,
                color: terminal.isActivated
                    ? AppColors.success
                    : AppColors.textTertiary,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildInfoCard(Terminal terminal) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '基本信息',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 16),
          _buildInfoItem('通道名称', terminal.channelCode),
          _buildInfoItem('品牌', terminal.brandCode ?? '-'),
          _buildInfoItem('型号', terminal.modelCode ?? '-'),
          _buildInfoItem('入库时间', terminal.createdAt.toString().substring(0, 16)),
        ],
      ),
    );
  }

  Widget _buildStatusCard(Terminal terminal) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '状态信息',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 16),
          _buildInfoItem('激活状态', terminal.isActivated ? '已激活' : '未激活'),
          if (terminal.activatedAt != null)
            _buildInfoItem('激活时间', terminal.activatedAt!.toLocal().toString().substring(0, 16)),
          _buildInfoItem('绑定状态', terminal.boundAt != null ? '已绑定' : '未绑定'),
          if (terminal.boundAt != null)
            _buildInfoItem('绑定时间', terminal.boundAt!.toLocal().toString().substring(0, 16)),
        ],
      ),
    );
  }

  Widget _buildMerchantCard(Terminal terminal) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '商户信息',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 16),
          _buildInfoItem('商户号', terminal.merchantNo ?? '-'),
        ],
      ),
    );
  }

  Widget _buildInfoItem(String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: const TextStyle(
              color: AppColors.textSecondary,
              fontSize: 14,
            ),
          ),
          Text(
            value,
            style: const TextStyle(
              color: AppColors.textPrimary,
              fontSize: 14,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }
}
