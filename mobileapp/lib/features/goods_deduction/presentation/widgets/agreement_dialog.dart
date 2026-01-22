import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';

/// 协议弹窗组件
class AgreementDialog extends StatefulWidget {
  final String title;
  final String content;
  final VoidCallback onAgree;
  final VoidCallback onReject;

  const AgreementDialog({
    super.key,
    required this.title,
    required this.content,
    required this.onAgree,
    required this.onReject,
  });

  @override
  State<AgreementDialog> createState() => _AgreementDialogState();
}

class _AgreementDialogState extends State<AgreementDialog> {
  final ScrollController _scrollController = ScrollController();
  bool _hasReadToBottom = false;

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_checkScrollPosition);
  }

  @override
  void dispose() {
    _scrollController.removeListener(_checkScrollPosition);
    _scrollController.dispose();
    super.dispose();
  }

  void _checkScrollPosition() {
    if (_scrollController.position.pixels >=
        _scrollController.position.maxScrollExtent - 50) {
      if (!_hasReadToBottom) {
        setState(() {
          _hasReadToBottom = true;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Dialog(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
      ),
      child: Container(
        width: double.maxFinite,
        constraints: BoxConstraints(
          maxHeight: MediaQuery.of(context).size.height * 0.7,
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            _buildHeader(),
            Flexible(child: _buildContent()),
            _buildActions(),
          ],
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: const BoxDecoration(
        border: Border(
          bottom: BorderSide(color: AppColors.divider),
        ),
      ),
      child: Row(
        children: [
          const Icon(Icons.description_outlined, color: AppColors.primary),
          const SizedBox(width: AppSpacing.sm),
          Expanded(
            child: Text(
              widget.title,
              style: const TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w600,
                color: AppColors.textPrimary,
              ),
            ),
          ),
          IconButton(
            icon: const Icon(Icons.close),
            onPressed: () => Navigator.of(context).pop(),
            padding: EdgeInsets.zero,
            constraints: const BoxConstraints(),
          ),
        ],
      ),
    );
  }

  Widget _buildContent() {
    return SingleChildScrollView(
      controller: _scrollController,
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Text(
        widget.content,
        style: const TextStyle(
          fontSize: 14,
          height: 1.8,
          color: AppColors.textSecondary,
        ),
      ),
    );
  }

  Widget _buildActions() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: const BoxDecoration(
        border: Border(
          top: BorderSide(color: AppColors.divider),
        ),
      ),
      child: Column(
        children: [
          if (!_hasReadToBottom)
            const Padding(
              padding: EdgeInsets.only(bottom: AppSpacing.sm),
              child: Text(
                '请阅读完协议内容后再进行操作',
                style: TextStyle(
                  fontSize: 12,
                  color: AppColors.warning,
                ),
              ),
            ),
          Row(
            children: [
              Expanded(
                child: OutlinedButton(
                  onPressed: widget.onReject,
                  style: OutlinedButton.styleFrom(
                    foregroundColor: AppColors.danger,
                    side: const BorderSide(color: AppColors.danger),
                    padding: const EdgeInsets.symmetric(vertical: 12),
                  ),
                  child: const Text('拒绝'),
                ),
              ),
              const SizedBox(width: AppSpacing.md),
              Expanded(
                child: ElevatedButton(
                  onPressed: _hasReadToBottom ? widget.onAgree : null,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppColors.success,
                    foregroundColor: Colors.white,
                    disabledBackgroundColor: AppColors.border,
                    padding: const EdgeInsets.symmetric(vertical: 12),
                  ),
                  child: Text(_hasReadToBottom ? '同意并接收' : '请阅读协议'),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

/// 显示协议弹窗
Future<bool?> showAgreementDialog({
  required BuildContext context,
  required String title,
  required String content,
}) {
  return showDialog<bool>(
    context: context,
    barrierDismissible: false,
    builder: (context) => AgreementDialog(
      title: title,
      content: content,
      onAgree: () => Navigator.of(context).pop(true),
      onReject: () => Navigator.of(context).pop(false),
    ),
  );
}

/// 获取默认协议内容
String getDefaultAgreementContent({
  required String fromAgentName,
  required String toAgentName,
  required double totalAmount,
  required int terminalCount,
}) {
  return '''
代扣服务协议

甲方（发起方）：$fromAgentName
乙方（接收方）：$toAgentName

第一条 服务内容
甲方根据与乙方的业务往来关系，按照本协议约定的方式和金额，从乙方的分润钱包和/或服务费钱包中代扣相应款项。

第二条 代扣金额
1. 代扣总金额：¥${totalAmount.toStringAsFixed(2)}
2. 关联终端数量：$terminalCount 台
3. 扣款频率：每日从乙方账户中自动扣款，直至扣完为止。

第三条 扣款规则
1. 扣款优先级：优先扣除分润钱包余额，分润余额不足时扣除服务费钱包余额；
2. 部分扣款：当乙方钱包余额不足时，系统将扣除全部可用余额，剩余部分继续扣除；
3. 扣款上限：每次扣款不设上限，有多少扣多少。

第四条 协议生效
本协议自乙方点击"同意并接收"按钮后生效，双方均应遵守协议约定。

第五条 违约责任
任何一方违反本协议约定的，应承担相应的违约责任。

第六条 争议解决
本协议履行过程中发生的争议，双方应友好协商解决；协商不成的，可向有管辖权的人民法院提起诉讼。

签署日期：${DateTime.now().toString().substring(0, 10)}
''';
}
