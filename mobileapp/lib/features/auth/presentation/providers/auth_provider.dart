import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/network/api_client.dart';
import '../../data/models/auth_model.dart';
import '../../data/services/auth_service.dart';

// 导出 AuthState 供其他文件使用
export '../../data/models/auth_model.dart' show AuthState, UserInfo;

/// AuthService Provider
final authServiceProvider = Provider<AuthService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return AuthService(apiClient);
});

/// 认证状态 Provider
final authStateProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final authService = ref.watch(authServiceProvider);
  return AuthNotifier(authService);
});

/// 认证状态管理器
class AuthNotifier extends StateNotifier<AuthState> {
  final AuthService _authService;

  AuthNotifier(this._authService) : super(const AuthState()) {
    // 初始化时检查登录状态
    _checkAuthStatus();
  }

  /// 检查认证状态
  Future<void> _checkAuthStatus() async {
    state = state.copyWith(isLoading: true);
    try {
      final isLoggedIn = await _authService.isLoggedIn();
      if (isLoggedIn) {
        final token = await _authService.getToken();
        state = state.copyWith(
          isAuthenticated: true,
          isLoading: false,
          accessToken: token,
        );
      } else {
        state = state.copyWith(
          isAuthenticated: false,
          isLoading: false,
        );
      }
    } catch (e) {
      state = state.copyWith(
        isAuthenticated: false,
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// 登录
  Future<bool> login(String username, String password) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final response = await _authService.login(
        LoginRequest(username: username, password: password),
      );

      state = state.copyWith(
        isAuthenticated: true,
        isLoading: false,
        user: response.user,
        accessToken: response.accessToken,
        error: null,
      );
      return true;
    } catch (e) {
      String errorMessage = '登录失败';
      if (e.toString().contains('401')) {
        errorMessage = '用户名或密码错误';
      } else if (e.toString().contains('network')) {
        errorMessage = '网络连接失败';
      } else {
        errorMessage = e.toString().replaceAll('Exception: ', '');
      }

      state = state.copyWith(
        isAuthenticated: false,
        isLoading: false,
        error: errorMessage,
      );
      return false;
    }
  }

  /// 登出
  Future<void> logout() async {
    state = state.copyWith(isLoading: true);
    try {
      await _authService.logout();
    } catch (e) {
      // 忽略错误，继续清除状态
    }
    state = const AuthState(isAuthenticated: false, isLoading: false);
  }

  /// 刷新令牌
  Future<bool> refreshToken() async {
    try {
      final refreshToken = await _authService.getRefreshToken();
      if (refreshToken == null) return false;

      final response = await _authService.refreshToken(refreshToken);
      state = state.copyWith(
        accessToken: response.accessToken,
        user: response.user,
      );
      return true;
    } catch (e) {
      await logout();
      return false;
    }
  }

  /// 清除错误
  void clearError() {
    state = state.copyWith(error: null);
  }

  /// 修改密码
  Future<bool> changePassword({
    required String oldPassword,
    required String newPassword,
  }) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      await _authService.changePassword(
        oldPassword: oldPassword,
        newPassword: newPassword,
      );
      // 密码修改成功，需要重新登录
      state = const AuthState(isAuthenticated: false, isLoading: false);
      return true;
    } catch (e) {
      String errorMessage = '修改密码失败';
      if (e.toString().contains('原密码错误')) {
        errorMessage = '原密码错误';
      } else {
        errorMessage = e.toString().replaceAll('Exception: ', '');
      }
      state = state.copyWith(
        isLoading: false,
        error: errorMessage,
      );
      return false;
    }
  }
}

/// 是否已登录 Provider
final isLoggedInProvider = Provider<bool>((ref) {
  final authState = ref.watch(authStateProvider);
  return authState.isAuthenticated;
});
