import 'package:dio/dio.dart';
import '../models/poster.dart';
import '../models/poster_category.dart';

/// 海报服务 - 处理营销海报相关API调用
class PosterService {
  final Dio _dio;

  PosterService(this._dio);

  /// 获取海报分类列表
  Future<List<PosterCategory>> getCategories() async {
    try {
      final response = await _dio.get('/api/v1/posters/categories');

      if (response.statusCode == 200) {
        final data = response.data;
        if (data['code'] == 0 && data['data'] != null) {
          final List<dynamic> list = data['data'];
          return list.map((json) => PosterCategory.fromJson(json)).toList();
        }
      }
      return [];
    } catch (e) {
      print('获取海报分类失败: $e');
      return [];
    }
  }

  /// 获取海报列表
  Future<PosterListResult> getPosters({
    int? categoryId,
    int page = 1,
    int pageSize = 20,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'page_size': pageSize,
      };
      if (categoryId != null) {
        queryParams['category_id'] = categoryId;
      }

      final response = await _dio.get(
        '/api/v1/posters',
        queryParameters: queryParams,
      );

      if (response.statusCode == 200) {
        final data = response.data;
        if (data['code'] == 0) {
          final List<dynamic> list = data['data'] ?? [];
          final int total = data['total'] ?? 0;
          return PosterListResult(
            posters: list.map((json) => Poster.fromJson(json)).toList(),
            total: total,
          );
        }
      }
      return PosterListResult(posters: [], total: 0);
    } catch (e) {
      print('获取海报列表失败: $e');
      return PosterListResult(posters: [], total: 0);
    }
  }

  /// 获取海报详情
  Future<Poster?> getPosterDetail(int posterId) async {
    try {
      final response = await _dio.get('/api/v1/posters/$posterId');

      if (response.statusCode == 200) {
        final data = response.data;
        if (data['code'] == 0 && data['data'] != null) {
          return Poster.fromJson(data['data']);
        }
      }
      return null;
    } catch (e) {
      print('获取海报详情失败: $e');
      return null;
    }
  }

  /// 记录海报下载
  Future<void> recordDownload(int posterId) async {
    try {
      await _dio.post('/api/v1/posters/$posterId/download');
    } catch (e) {
      print('记录海报下载失败: $e');
    }
  }

  /// 记录海报分享
  Future<void> recordShare(int posterId) async {
    try {
      await _dio.post('/api/v1/posters/$posterId/share');
    } catch (e) {
      print('记录海报分享失败: $e');
    }
  }
}

/// 海报列表结果
class PosterListResult {
  final List<Poster> posters;
  final int total;

  PosterListResult({
    required this.posters,
    required this.total,
  });
}
