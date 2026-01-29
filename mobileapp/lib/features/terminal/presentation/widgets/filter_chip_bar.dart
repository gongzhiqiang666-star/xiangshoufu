import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';
import 'terminal_filter_sheet.dart';

/// 筛选标签条
/// 横向滚动显示已选筛选条件，每个标签带×可删除
/// 无筛选时隐藏
class FilterChipBar extends StatelessWidget {
  final List<FilterTag> tags;
  final ValueChanged<FilterTag> onRemove;
  final VoidCallback? onClearAll;

  const FilterChipBar({
    super.key,
    required this.tags,
    required this.onRemove,
    this.onClearAll,
  });

  @override
  Widget build(BuildContext context) {
    if (tags.isEmpty) {
      return const SizedBox.shrink();
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Row(
          children: [
            ...tags.map((tag) => Padding(
                  padding: const EdgeInsets.only(right: 8),
                  child: _buildChip(tag),
                )),
            if (onClearAll != null && tags.length > 1)
              GestureDetector(
                onTap: onClearAll,
                child: Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
                  decoration: BoxDecoration(
                    color: AppColors.danger.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: const Text(
                    '清空',
                    style: TextStyle(
                      fontSize: 12,
                      color: AppColors.danger,
                    ),
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildChip(FilterTag tag) {
    return Container(
      padding: const EdgeInsets.only(left: 12, right: 6, top: 6, bottom: 6),
      decoration: BoxDecoration(
        color: AppColors.primary.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: AppColors.primary.withValues(alpha: 0.3),
          width: 1,
        ),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            tag.label,
            style: const TextStyle(
              fontSize: 12,
              color: AppColors.primary,
              fontWeight: FontWeight.w500,
            ),
          ),
          const SizedBox(width: 4),
          GestureDetector(
            onTap: () => onRemove(tag),
            child: Container(
              padding: const EdgeInsets.all(2),
              decoration: BoxDecoration(
                color: AppColors.primary.withValues(alpha: 0.2),
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.close,
                size: 12,
                color: AppColors.primary,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
