import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:dio/dio.dart';
import '../../data/models/banner.dart';
import '../../data/services/banner_service.dart';

/// Banner服务Provider
final bannerServiceProvider = Provider<BannerService>((ref) {
  // 这里需要从全局获取Dio实例
  // 假设有一个全局的dioProvider
  final dio = Dio(BaseOptions(
    baseUrl: 'http://localhost:8080', // 开发环境
    connectTimeout: const Duration(seconds: 10),
    receiveTimeout: const Duration(seconds: 10),
  ));
  return BannerService(dio);
});

/// Banner列表状态
class BannerListState {
  final List<Banner> banners;
  final bool isLoading;
  final String? error;

  BannerListState({
    this.banners = const [],
    this.isLoading = false,
    this.error,
  });

  BannerListState copyWith({
    List<Banner>? banners,
    bool? isLoading,
    String? error,
  }) {
    return BannerListState(
      banners: banners ?? this.banners,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

/// Banner列表Notifier
class BannerListNotifier extends StateNotifier<BannerListState> {
  final BannerService _service;

  BannerListNotifier(this._service) : super(BannerListState());

  /// 加载Banner列表
  Future<void> loadBanners() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final banners = await _service.getActiveBanners();
      state = state.copyWith(banners: banners, isLoading: false);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  /// 记录点击
  Future<void> recordClick(int bannerId) async {
    await _service.recordClick(bannerId);
  }
}

/// Banner列表Provider
final bannerListProvider =
    StateNotifierProvider<BannerListNotifier, BannerListState>((ref) {
  final service = ref.watch(bannerServiceProvider);
  return BannerListNotifier(service);
});
