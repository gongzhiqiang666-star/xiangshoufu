import 'package:dio/dio.dart';
import '../models/banner.dart';

/// Banner服务 - 处理滚动图相关API调用
class BannerService {
  final Dio _dio;

  BannerService(this._dio);

  /// 获取有效的Banner列表
  Future<List<Banner>> getActiveBanners() async {
    try {
      final response = await _dio.get('/api/v1/banners');

      if (response.statusCode == 200) {
        final data = response.data;
        if (data['code'] == 0 && data['data'] != null) {
          final List<dynamic> list = data['data'];
          return list.map((json) => Banner.fromJson(json)).toList();
        }
      }
      return [];
    } catch (e) {
      print('获取Banner列表失败: $e');
      return [];
    }
  }

  /// 记录Banner点击
  Future<void> recordClick(int bannerId) async {
    try {
      await _dio.post('/api/v1/banners/$bannerId/click');
    } catch (e) {
      // 静默失败，不影响用户体验
      print('记录Banner点击失败: $e');
    }
  }
}
