import 'package:shared_preferences/shared_preferences.dart';
import '../../../../core/network/api_client.dart';
import '../models/auth_model.dart';

/// 认证服务
class AuthService {
  final ApiClient _apiClient;
  static const String _tokenKey = 'token';
  static const String _refreshTokenKey = 'refresh_token';
  static const String _userKey = 'user_info';

  AuthService(this._apiClient);

  /// 登录
  Future<LoginResponse> login(LoginRequest request) async {
    final response = await _apiClient.post(
      '/api/v1/auth/login',
      data: request.toJson(),
    );

    final loginResponse = LoginResponse.fromJson(response.data['data']);

    // 保存 token
    await saveToken(loginResponse.accessToken);
    await saveRefreshToken(loginResponse.refreshToken);

    return loginResponse;
  }

  /// 刷新令牌
  Future<LoginResponse> refreshToken(String refreshToken) async {
    final response = await _apiClient.post(
      '/api/v1/auth/refresh',
      data: {'refresh_token': refreshToken},
    );

    final loginResponse = LoginResponse.fromJson(response.data['data']);

    // 更新 token
    await saveToken(loginResponse.accessToken);
    await saveRefreshToken(loginResponse.refreshToken);

    return loginResponse;
  }

  /// 登出
  Future<void> logout() async {
    try {
      await _apiClient.post('/api/v1/auth/logout');
    } catch (e) {
      // 忽略登出错误，继续清除本地数据
    }
    await clearTokens();
  }

  /// 保存访问令牌
  Future<void> saveToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_tokenKey, token);
  }

  /// 获取访问令牌
  Future<String?> getToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_tokenKey);
  }

  /// 保存刷新令牌
  Future<void> saveRefreshToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_refreshTokenKey, token);
  }

  /// 获取刷新令牌
  Future<String?> getRefreshToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_refreshTokenKey);
  }

  /// 清除所有令牌
  Future<void> clearTokens() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_tokenKey);
    await prefs.remove(_refreshTokenKey);
    await prefs.remove(_userKey);
  }

  /// 检查是否已登录
  Future<bool> isLoggedIn() async {
    final token = await getToken();
    return token != null && token.isNotEmpty;
  }

  /// 修改密码
  Future<void> changePassword({
    required String oldPassword,
    required String newPassword,
  }) async {
    await _apiClient.post('/api/v1/auth/change-password', data: {
      'old_password': oldPassword,
      'new_password': newPassword,
    });
  }
}
