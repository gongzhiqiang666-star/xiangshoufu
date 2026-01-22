/// 登录请求
class LoginRequest {
  final String username;
  final String password;

  LoginRequest({
    required this.username,
    required this.password,
  });

  Map<String, dynamic> toJson() => {
        'username': username,
        'password': password,
      };
}

/// 登录响应
class LoginResponse {
  final String accessToken;
  final String refreshToken;
  final int expiresIn;
  final UserInfo user;

  LoginResponse({
    required this.accessToken,
    required this.refreshToken,
    required this.expiresIn,
    required this.user,
  });

  factory LoginResponse.fromJson(Map<String, dynamic> json) {
    return LoginResponse(
      accessToken: json['access_token'] ?? '',
      refreshToken: json['refresh_token'] ?? '',
      expiresIn: json['expires_in'] ?? 0,
      user: UserInfo.fromJson(json['user'] ?? {}),
    );
  }
}

/// 用户信息
class UserInfo {
  final int id;
  final String username;
  final String? agentName;
  final String? phone;
  final int? agentId;
  final String role;

  UserInfo({
    required this.id,
    required this.username,
    this.agentName,
    this.phone,
    this.agentId,
    required this.role,
  });

  factory UserInfo.fromJson(Map<String, dynamic> json) {
    return UserInfo(
      id: json['id'] ?? 0,
      username: json['username'] ?? '',
      agentName: json['agent_name'],
      phone: json['phone'],
      agentId: json['agent_id'],
      role: json['role'] ?? 'agent',
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'username': username,
        'agent_name': agentName,
        'phone': phone,
        'agent_id': agentId,
        'role': role,
      };
}

/// 认证状态
class AuthState {
  final bool isAuthenticated;
  final bool isLoading;
  final String? error;
  final UserInfo? user;
  final String? accessToken;

  const AuthState({
    this.isAuthenticated = false,
    this.isLoading = false,
    this.error,
    this.user,
    this.accessToken,
  });

  AuthState copyWith({
    bool? isAuthenticated,
    bool? isLoading,
    String? error,
    UserInfo? user,
    String? accessToken,
  }) {
    return AuthState(
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      isLoading: isLoading ?? this.isLoading,
      error: error,
      user: user ?? this.user,
      accessToken: accessToken ?? this.accessToken,
    );
  }
}
