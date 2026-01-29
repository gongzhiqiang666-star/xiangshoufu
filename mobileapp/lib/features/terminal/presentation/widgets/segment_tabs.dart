import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';

/// 终端状态分段配置
class TerminalSegment {
  final String key;
  final String label;
  final int count;

  const TerminalSegment({
    required this.key,
    required this.label,
    this.count = 0,
  });

  TerminalSegment copyWith({int? count}) {
    return TerminalSegment(
      key: key,
      label: label,
      count: count ?? this.count,
    );
  }
}

/// 终端状态分段控制器
/// 显示4个状态: 全部、已激活、未激活、未绑定
/// 每个状态显示对应数量
class TerminalSegmentTabs extends StatelessWidget {
  final List<TerminalSegment> segments;
  final int selectedIndex;
  final ValueChanged<int> onChanged;

  const TerminalSegmentTabs({
    super.key,
    required this.segments,
    required this.selectedIndex,
    required this.onChanged,
  });

  /// 默认的终端状态分段配置
  static List<TerminalSegment> get defaultSegments => const [
        TerminalSegment(key: 'all', label: '全部'),
        TerminalSegment(key: 'active', label: '已激活'),
        TerminalSegment(key: 'inactive', label: '未激活'),
        TerminalSegment(key: 'unbound', label: '未绑定'),
      ];

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      decoration: BoxDecoration(
        color: Colors.grey.shade100,
        borderRadius: BorderRadius.circular(10),
      ),
      padding: const EdgeInsets.all(4),
      child: Row(
        children: List.generate(segments.length, (index) {
          final segment = segments[index];
          final isSelected = index == selectedIndex;

          return Expanded(
            child: GestureDetector(
              onTap: () => onChanged(index),
              child: AnimatedContainer(
                duration: const Duration(milliseconds: 200),
                padding: const EdgeInsets.symmetric(vertical: 10),
                decoration: BoxDecoration(
                  color: isSelected ? Colors.white : Colors.transparent,
                  borderRadius: BorderRadius.circular(8),
                  boxShadow: isSelected
                      ? [
                          BoxShadow(
                            color: Colors.black.withValues(alpha: 0.08),
                            blurRadius: 4,
                            offset: const Offset(0, 2),
                          ),
                        ]
                      : null,
                ),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      segment.label,
                      style: TextStyle(
                        fontSize: 13,
                        fontWeight:
                            isSelected ? FontWeight.w600 : FontWeight.normal,
                        color: isSelected
                            ? AppColors.primary
                            : AppColors.textSecondary,
                      ),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      '${segment.count}',
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.bold,
                        color: isSelected
                            ? AppColors.primary
                            : AppColors.textTertiary,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          );
        }),
      ),
    );
  }
}
