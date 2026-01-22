/// 消息类型
enum MessageType {
  profit(1, '交易分润', 'profit'),
  activation(2, '激活奖励', 'profit'),
  deposit(3, '押金返现', 'profit'),
  simCashback(4, '流量返现', 'profit'),
  refund(5, '退款撤销', 'system'),
  announcement(6, '系统公告', 'system'),
  newAgent(7, '新代理注册', 'register'),
  transaction(8, '交易通知', 'consumption');

  final int value;
  final String label;
  final String category;

  const MessageType(this.value, this.label, this.category);

  static MessageType fromValue(int value) {
    return MessageType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => MessageType.announcement,
    );
  }
}

/// 消息分类
enum MessageCategory {
  all('all', '全部', [1, 2, 3, 4, 5, 6, 7, 8]),
  profit('profit', '分润', [1, 2, 3, 4]),
  register('register', '注册', [7]),
  consumption('consumption', '消费', [8]),
  system('system', '系统', [5, 6]);

  final String value;
  final String label;
  final List<int> types;

  const MessageCategory(this.value, this.label, this.types);

  static MessageCategory fromValue(String value) {
    return MessageCategory.values.firstWhere(
      (e) => e.value == value,
      orElse: () => MessageCategory.all,
    );
  }
}

/// 消息信息
class MessageModel {
  final int id;
  final int? agentId;
  final int messageType;
  final String? typeName;
  final String title;
  final String content;
  final bool isRead;
  final bool? isPushed;
  final int? relatedId;
  final String? relatedType;
  final String? expireAt;
  final String createdAt;

  MessageModel({
    required this.id,
    this.agentId,
    required this.messageType,
    this.typeName,
    required this.title,
    required this.content,
    required this.isRead,
    this.isPushed,
    this.relatedId,
    this.relatedType,
    this.expireAt,
    required this.createdAt,
  });

  factory MessageModel.fromJson(Map<String, dynamic> json) {
    return MessageModel(
      id: json['id'] ?? 0,
      agentId: json['agent_id'],
      messageType: json['message_type'] ?? 6,
      typeName: json['type_name'],
      title: json['title'] ?? '',
      content: json['content'] ?? '',
      isRead: json['is_read'] ?? false,
      isPushed: json['is_pushed'],
      relatedId: json['related_id'],
      relatedType: json['related_type'],
      expireAt: json['expire_at'],
      createdAt: json['created_at'] ?? '',
    );
  }

  MessageType get type => MessageType.fromValue(messageType);

  String get messageTypeName => typeName ?? type.label;

  /// 判断是否为今天的消息
  bool get isToday {
    final now = DateTime.now();
    final created = DateTime.tryParse(createdAt);
    if (created == null) return false;
    return now.year == created.year &&
        now.month == created.month &&
        now.day == created.day;
  }

  /// 判断是否为昨天的消息
  bool get isYesterday {
    final now = DateTime.now();
    final yesterday = now.subtract(const Duration(days: 1));
    final created = DateTime.tryParse(createdAt);
    if (created == null) return false;
    return yesterday.year == created.year &&
        yesterday.month == created.month &&
        yesterday.day == created.day;
  }

  /// 格式化时间显示
  String get formattedTime {
    final created = DateTime.tryParse(createdAt);
    if (created == null) return '';

    if (isToday) {
      return '${created.hour.toString().padLeft(2, '0')}:${created.minute.toString().padLeft(2, '0')}';
    } else if (isYesterday) {
      return '昨天 ${created.hour.toString().padLeft(2, '0')}:${created.minute.toString().padLeft(2, '0')}';
    } else {
      return '${created.month}/${created.day} ${created.hour.toString().padLeft(2, '0')}:${created.minute.toString().padLeft(2, '0')}';
    }
  }

  /// 获取日期分组标题
  String get dateGroupTitle {
    if (isToday) return '今天';
    if (isYesterday) return '昨天';
    return '更早';
  }
}

/// 未读消息统计
class UnreadCountModel {
  final int count;

  UnreadCountModel({required this.count});

  factory UnreadCountModel.fromJson(Map<String, dynamic> json) {
    return UnreadCountModel(count: json['count'] ?? 0);
  }
}

/// 消息统计
class MessageStatsModel {
  final int total;
  final int unreadTotal;
  final int profitCount;
  final int registerCount;
  final int consumptionCount;
  final int systemCount;

  MessageStatsModel({
    required this.total,
    required this.unreadTotal,
    required this.profitCount,
    required this.registerCount,
    required this.consumptionCount,
    required this.systemCount,
  });

  factory MessageStatsModel.fromJson(Map<String, dynamic> json) {
    return MessageStatsModel(
      total: json['total'] ?? 0,
      unreadTotal: json['unread_total'] ?? 0,
      profitCount: json['profit_count'] ?? 0,
      registerCount: json['register_count'] ?? 0,
      consumptionCount: json['consumption_count'] ?? 0,
      systemCount: json['system_count'] ?? 0,
    );
  }
}

/// 消息类型信息
class MessageTypeInfo {
  final int value;
  final String label;
  final String category;

  MessageTypeInfo({
    required this.value,
    required this.label,
    required this.category,
  });

  factory MessageTypeInfo.fromJson(Map<String, dynamic> json) {
    return MessageTypeInfo(
      value: json['value'] ?? 0,
      label: json['label'] ?? '',
      category: json['category'] ?? '',
    );
  }
}

/// 消息分类信息
class MessageCategoryInfo {
  final String value;
  final String label;

  MessageCategoryInfo({
    required this.value,
    required this.label,
  });

  factory MessageCategoryInfo.fromJson(Map<String, dynamic> json) {
    return MessageCategoryInfo(
      value: json['value'] ?? '',
      label: json['label'] ?? '',
    );
  }
}
