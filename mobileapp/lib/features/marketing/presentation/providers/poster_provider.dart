import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:dio/dio.dart';
import '../data/models/poster.dart';
import '../data/models/poster_category.dart';
import '../data/services/poster_service.dart';

/// 海报服务Provider
final posterServiceProvider = Provider<PosterService>((ref) {
  final dio = Dio(BaseOptions(
    baseUrl: 'http://localhost:8080', // 开发环境
    connectTimeout: const Duration(seconds: 10),
    receiveTimeout: const Duration(seconds: 10),
  ));
  return PosterService(dio);
});

/// 海报分类状态
class PosterCategoryState {
  final List<PosterCategory> categories;
  final bool isLoading;
  final String? error;
  final int? selectedCategoryId;

  PosterCategoryState({
    this.categories = const [],
    this.isLoading = false,
    this.error,
    this.selectedCategoryId,
  });

  PosterCategoryState copyWith({
    List<PosterCategory>? categories,
    bool? isLoading,
    String? error,
    int? selectedCategoryId,
  }) {
    return PosterCategoryState(
      categories: categories ?? this.categories,
      isLoading: isLoading ?? this.isLoading,
      error: error,
      selectedCategoryId: selectedCategoryId ?? this.selectedCategoryId,
    );
  }
}

/// 海报分类Notifier
class PosterCategoryNotifier extends StateNotifier<PosterCategoryState> {
  final PosterService _service;

  PosterCategoryNotifier(this._service) : super(PosterCategoryState());

  /// 加载分类列表
  Future<void> loadCategories() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final categories = await _service.getCategories();
      state = state.copyWith(categories: categories, isLoading: false);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  /// 选择分类
  void selectCategory(int? categoryId) {
    state = state.copyWith(selectedCategoryId: categoryId);
  }
}

/// 海报列表状态
class PosterListState {
  final List<Poster> posters;
  final bool isLoading;
  final bool isLoadingMore;
  final String? error;
  final int total;
  final int page;
  final bool hasMore;

  PosterListState({
    this.posters = const [],
    this.isLoading = false,
    this.isLoadingMore = false,
    this.error,
    this.total = 0,
    this.page = 1,
    this.hasMore = true,
  });

  PosterListState copyWith({
    List<Poster>? posters,
    bool? isLoading,
    bool? isLoadingMore,
    String? error,
    int? total,
    int? page,
    bool? hasMore,
  }) {
    return PosterListState(
      posters: posters ?? this.posters,
      isLoading: isLoading ?? this.isLoading,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      error: error,
      total: total ?? this.total,
      page: page ?? this.page,
      hasMore: hasMore ?? this.hasMore,
    );
  }
}

/// 海报列表Notifier
class PosterListNotifier extends StateNotifier<PosterListState> {
  final PosterService _service;
  static const int _pageSize = 20;

  PosterListNotifier(this._service) : super(PosterListState());

  /// 加载海报列表
  Future<void> loadPosters({int? categoryId, bool refresh = false}) async {
    if (refresh) {
      state = state.copyWith(isLoading: true, error: null, page: 1);
    } else {
      state = state.copyWith(isLoading: true, error: null);
    }

    try {
      final result = await _service.getPosters(
        categoryId: categoryId,
        page: 1,
        pageSize: _pageSize,
      );
      state = state.copyWith(
        posters: result.posters,
        total: result.total,
        page: 1,
        hasMore: result.posters.length < result.total,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  /// 加载更多
  Future<void> loadMore({int? categoryId}) async {
    if (state.isLoadingMore || !state.hasMore) return;

    state = state.copyWith(isLoadingMore: true);

    try {
      final nextPage = state.page + 1;
      final result = await _service.getPosters(
        categoryId: categoryId,
        page: nextPage,
        pageSize: _pageSize,
      );
      final allPosters = [...state.posters, ...result.posters];
      state = state.copyWith(
        posters: allPosters,
        page: nextPage,
        hasMore: allPosters.length < result.total,
        isLoadingMore: false,
      );
    } catch (e) {
      state = state.copyWith(isLoadingMore: false, error: e.toString());
    }
  }

  /// 记录下载
  Future<void> recordDownload(int posterId) async {
    await _service.recordDownload(posterId);
  }

  /// 记录分享
  Future<void> recordShare(int posterId) async {
    await _service.recordShare(posterId);
  }
}

/// 海报分类Provider
final posterCategoryProvider =
    StateNotifierProvider<PosterCategoryNotifier, PosterCategoryState>((ref) {
  final service = ref.watch(posterServiceProvider);
  return PosterCategoryNotifier(service);
});

/// 海报列表Provider
final posterListProvider =
    StateNotifierProvider<PosterListNotifier, PosterListState>((ref) {
  final service = ref.watch(posterServiceProvider);
  return PosterListNotifier(service);
});
