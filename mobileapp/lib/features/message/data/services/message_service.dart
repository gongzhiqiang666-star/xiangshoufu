import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/network/api_client.dart';
import '../models/message_model.dart';

/// Message服务Provider
final messageServiceProvider = Provider<MessageService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return MessageService(apiClient);
});

/// 消息服务
class MessageService {
  final ApiClient _apiClient;

  MessageService(this._apiClient);

  /// 获取消息列表
  Future<PaginatedResponse<MessageModel>> getMessages({
    String? type,
    String? category,
    bool? isRead,
    int page = 1,
    int pageSize = 20,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (type != null) queryParams['type'] = type;
    if (category != null) queryParams['category'] = category;
    if (isRead != null) queryParams['is_read'] = isRead;

    final response = await _apiClient.get(
      '/api/v1/messages',
      queryParameters: queryParams,
    );
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return PaginatedResponse.fromJson(
      apiResponse.data,
      (json) => MessageModel.fromJson(json),
    );
  }

  /// 获取未读消息数量
  Future<UnreadCountModel> getUnreadCount() async {
    final response = await _apiClient.get('/api/v1/messages/unread-count');
    final apiResponse = ApiResponse.fromJson(
      response.data,
      (data) => UnreadCountModel.fromJson(data),
    );
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data!;
  }

  /// 获取消息统计
  Future<MessageStatsModel> getMessageStats() async {
    final response = await _apiClient.get('/api/v1/messages/stats');
    final apiResponse = ApiResponse.fromJson(
      response.data,
      (data) => MessageStatsModel.fromJson(data),
    );
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data!;
  }

  /// 获取消息类型和分类
  Future<({List<MessageTypeInfo> types, List<MessageCategoryInfo> categories})>
      getMessageTypes() async {
    final response = await _apiClient.get('/api/v1/messages/types');
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }

    final typesJson = apiResponse.data['types'] as List? ?? [];
    final categoriesJson = apiResponse.data['categories'] as List? ?? [];

    return (
      types: typesJson.map((e) => MessageTypeInfo.fromJson(e)).toList(),
      categories:
          categoriesJson.map((e) => MessageCategoryInfo.fromJson(e)).toList(),
    );
  }

  /// 获取消息详情
  Future<MessageModel> getMessageDetail(int id) async {
    final response = await _apiClient.get('/api/v1/messages/$id');
    final apiResponse = ApiResponse.fromJson(
      response.data,
      (data) => MessageModel.fromJson(data),
    );
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data!;
  }

  /// 标记消息为已读
  Future<void> markAsRead(int id) async {
    final response = await _apiClient.put('/api/v1/messages/$id/read');
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
  }

  /// 标记所有消息为已读
  Future<void> markAllAsRead() async {
    final response = await _apiClient.put('/api/v1/messages/read-all');
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
  }
}
