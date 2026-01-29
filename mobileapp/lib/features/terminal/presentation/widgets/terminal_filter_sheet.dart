import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';
import '../../domain/models/terminal.dart';

/// 终端筛选条件
class TerminalFilterCondition {
  final int? channelId;
  final String? channelCode;
  final String? brandCode;
  final String? modelCode;
  final String? extraStatus; // 更多状态: unstock, stocked

  const TerminalFilterCondition({
    this.channelId,
    this.channelCode,
    this.brandCode,
    this.modelCode,
    this.extraStatus,
  });

  TerminalFilterCondition copyWith({
    int? channelId,
    String? channelCode,
    String? brandCode,
    String? modelCode,
    String? extraStatus,
    bool clearChannel = false,
    bool clearTerminalType = false,
    bool clearExtraStatus = false,
  }) {
    return TerminalFilterCondition(
      channelId: clearChannel ? null : (channelId ?? this.channelId),
      channelCode: clearChannel ? null : (channelCode ?? this.channelCode),
      brandCode: clearTerminalType ? null : (brandCode ?? this.brandCode),
      modelCode: clearTerminalType ? null : (modelCode ?? this.modelCode),
      extraStatus: clearExtraStatus ? null : (extraStatus ?? this.extraStatus),
    );
  }

  /// 是否有筛选条件
  bool get hasFilters =>
      channelId != null || brandCode != null || extraStatus != null;

  /// 获取筛选标签列表
  List<FilterTag> get filterTags {
    final tags = <FilterTag>[];
    if (channelCode != null) {
      tags.add(FilterTag(
        type: FilterTagType.channel,
        label: channelCode!,
        value: channelId.toString(),
      ));
    }
    if (brandCode != null) {
      final label = modelCode != null ? '$brandCode $modelCode' : brandCode!;
      tags.add(FilterTag(
        type: FilterTagType.terminalType,
        label: label,
        value: '$brandCode-$modelCode',
      ));
    }
    if (extraStatus != null) {
      final label = extraStatus == 'unstock' ? '未出库' : '已出库';
      tags.add(FilterTag(
        type: FilterTagType.extraStatus,
        label: label,
        value: extraStatus!,
      ));
    }
    return tags;
  }

  /// 清空所有筛选
  static const TerminalFilterCondition empty = TerminalFilterCondition();
}

/// 筛选标签类型
enum FilterTagType {
  channel,
  terminalType,
  extraStatus,
}

/// 筛选标签
class FilterTag {
  final FilterTagType type;
  final String label;
  final String value;

  const FilterTag({
    required this.type,
    required this.label,
    required this.value,
  });
}

/// 终端筛选底部Sheet
class TerminalFilterSheet extends StatefulWidget {
  final TerminalFilterOptions options;
  final TerminalFilterCondition initialCondition;
  final ValueChanged<TerminalFilterCondition> onApply;

  const TerminalFilterSheet({
    super.key,
    required this.options,
    required this.initialCondition,
    required this.onApply,
  });

  /// 显示筛选Sheet
  static Future<void> show({
    required BuildContext context,
    required TerminalFilterOptions options,
    required TerminalFilterCondition initialCondition,
    required ValueChanged<TerminalFilterCondition> onApply,
  }) {
    return showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => TerminalFilterSheet(
        options: options,
        initialCondition: initialCondition,
        onApply: onApply,
      ),
    );
  }

  @override
  State<TerminalFilterSheet> createState() => _TerminalFilterSheetState();
}

class _TerminalFilterSheetState extends State<TerminalFilterSheet> {
  late TerminalFilterCondition _condition;

  @override
  void initState() {
    super.initState();
    _condition = widget.initialCondition;
  }

  void _reset() {
    setState(() {
      _condition = TerminalFilterCondition.empty;
    });
  }

  void _apply() {
    widget.onApply(_condition);
    Navigator.pop(context);
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      constraints: BoxConstraints(
        maxHeight: MediaQuery.of(context).size.height * 0.7,
      ),
      decoration: const BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // 标题栏
          _buildHeader(),
          const Divider(height: 1),
          // 筛选内容
          Flexible(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // 通道筛选
                  _buildSectionTitle('通道'),
                  const SizedBox(height: 8),
                  _buildChannelChips(),
                  const SizedBox(height: 20),
                  // 终端类型筛选
                  _buildSectionTitle('终端类型'),
                  const SizedBox(height: 8),
                  _buildTerminalTypeChips(),
                  const SizedBox(height: 20),
                  // 更多状态
                  _buildSectionTitle('更多状态'),
                  const SizedBox(height: 8),
                  _buildExtraStatusChips(),
                  const SizedBox(height: 20),
                ],
              ),
            ),
          ),
          // 确认按钮
          _buildFooter(),
        ],
      ),
    );
  }

  Widget _buildHeader() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          const Text(
            '筛选条件',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          TextButton(
            onPressed: _reset,
            child: const Text('重置'),
          ),
        ],
      ),
    );
  }

  Widget _buildSectionTitle(String title) {
    return Text(
      title,
      style: const TextStyle(
        fontSize: 14,
        fontWeight: FontWeight.w500,
        color: AppColors.textSecondary,
      ),
    );
  }

  Widget _buildChannelChips() {
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Row(
        children: [
          _buildFilterChip(
            label: '全部',
            isSelected: _condition.channelId == null,
            onTap: () {
              setState(() {
                _condition = _condition.copyWith(clearChannel: true);
              });
            },
          ),
          const SizedBox(width: 8),
          ...widget.options.channels.map((channel) => Padding(
                padding: const EdgeInsets.only(right: 8),
                child: _buildFilterChip(
                  label: channel.channelCode,
                  isSelected: _condition.channelId == channel.channelId,
                  onTap: () {
                    setState(() {
                      _condition = _condition.copyWith(
                        channelId: channel.channelId,
                        channelCode: channel.channelCode,
                        // 切换通道时清空终端类型
                        clearTerminalType: true,
                      );
                    });
                  },
                ),
              )),
        ],
      ),
    );
  }

  Widget _buildTerminalTypeChips() {
    // 根据选中的通道过滤终端类型
    final filteredTypes = _condition.channelId == null
        ? widget.options.terminalTypes
        : widget.options.terminalTypes
            .where((t) => t.channelId == _condition.channelId)
            .toList();

    // 构建唯一的终端类型列表
    final uniqueTypes = <String, TerminalTypeOption>{};
    for (final type in filteredTypes) {
      final key = '${type.brandCode}-${type.modelCode}';
      if (!uniqueTypes.containsKey(key)) {
        uniqueTypes[key] = type;
      }
    }

    if (uniqueTypes.isEmpty) {
      return const Text(
        '暂无终端类型',
        style: TextStyle(fontSize: 13, color: AppColors.textTertiary),
      );
    }

    return Wrap(
      spacing: 8,
      runSpacing: 8,
      children: [
        _buildFilterChip(
          label: '全部',
          isSelected: _condition.brandCode == null,
          onTap: () {
            setState(() {
              _condition = _condition.copyWith(clearTerminalType: true);
            });
          },
        ),
        ...uniqueTypes.entries.map((entry) {
          final type = entry.value;
          final isSelected = _condition.brandCode == type.brandCode &&
              _condition.modelCode == type.modelCode;
          return _buildFilterChip(
            label: type.displayName,
            isSelected: isSelected,
            onTap: () {
              setState(() {
                _condition = _condition.copyWith(
                  brandCode: type.brandCode,
                  modelCode: type.modelCode,
                );
              });
            },
          );
        }),
      ],
    );
  }

  Widget _buildExtraStatusChips() {
    return Row(
      children: [
        _buildFilterChip(
          label: '全部',
          isSelected: _condition.extraStatus == null,
          onTap: () {
            setState(() {
              _condition = _condition.copyWith(clearExtraStatus: true);
            });
          },
        ),
        const SizedBox(width: 8),
        _buildFilterChip(
          label: '未出库',
          isSelected: _condition.extraStatus == 'unstock',
          onTap: () {
            setState(() {
              _condition = _condition.copyWith(extraStatus: 'unstock');
            });
          },
        ),
        const SizedBox(width: 8),
        _buildFilterChip(
          label: '已出库',
          isSelected: _condition.extraStatus == 'stocked',
          onTap: () {
            setState(() {
              _condition = _condition.copyWith(extraStatus: 'stocked');
            });
          },
        ),
      ],
    );
  }

  Widget _buildFilterChip({
    required String label,
    required bool isSelected,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
        decoration: BoxDecoration(
          color: isSelected
              ? AppColors.primary.withValues(alpha: 0.1)
              : Colors.grey.shade100,
          borderRadius: BorderRadius.circular(20),
          border: Border.all(
            color: isSelected ? AppColors.primary : Colors.transparent,
            width: 1,
          ),
        ),
        child: Text(
          label,
          style: TextStyle(
            fontSize: 13,
            color: isSelected ? AppColors.primary : AppColors.textSecondary,
            fontWeight: isSelected ? FontWeight.w500 : FontWeight.normal,
          ),
        ),
      ),
    );
  }

  Widget _buildFooter() {
    return Container(
      padding: EdgeInsets.only(
        left: 16,
        right: 16,
        top: 12,
        bottom: MediaQuery.of(context).padding.bottom + 12,
      ),
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 10,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: SizedBox(
        width: double.infinity,
        child: ElevatedButton(
          onPressed: _apply,
          style: ElevatedButton.styleFrom(
            padding: const EdgeInsets.symmetric(vertical: 14),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(10),
            ),
          ),
          child: const Text('确认筛选'),
        ),
      ),
    );
  }
}
