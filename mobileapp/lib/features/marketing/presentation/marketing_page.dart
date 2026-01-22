import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/poster_provider.dart';
import '../widgets/poster_card.dart';
import 'poster_detail_page.dart';
import '../../data/models/poster.dart';

/// 营销海报主页
class MarketingPage extends ConsumerStatefulWidget {
  const MarketingPage({super.key});

  @override
  ConsumerState<MarketingPage> createState() => _MarketingPageState();
}

class _MarketingPageState extends ConsumerState<MarketingPage>
    with SingleTickerProviderStateMixin {
  TabController? _tabController;
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    // 加载分类和海报
    Future.microtask(() {
      ref.read(posterCategoryProvider.notifier).loadCategories();
      ref.read(posterListProvider.notifier).loadPosters();
    });

    // 监听滚动加载更多
    _scrollController.addListener(_onScroll);
  }

  @override
  void dispose() {
    _tabController?.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >=
        _scrollController.position.maxScrollExtent - 200) {
      final categoryState = ref.read(posterCategoryProvider);
      ref
          .read(posterListProvider.notifier)
          .loadMore(categoryId: categoryState.selectedCategoryId);
    }
  }

  @override
  Widget build(BuildContext context) {
    final categoryState = ref.watch(posterCategoryProvider);
    final posterState = ref.watch(posterListProvider);

    // 初始化TabController
    if (categoryState.categories.isNotEmpty && _tabController == null) {
      _tabController = TabController(
        length: categoryState.categories.length + 1, // +1 for "全部"
        vsync: this,
      );
      _tabController!.addListener(_onTabChanged);
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('营销海报'),
        elevation: 0,
        bottom: categoryState.categories.isNotEmpty
            ? TabBar(
                controller: _tabController,
                isScrollable: true,
                labelColor: Theme.of(context).primaryColor,
                unselectedLabelColor: Colors.grey,
                indicatorColor: Theme.of(context).primaryColor,
                tabs: [
                  const Tab(text: '全部'),
                  ...categoryState.categories.map((cat) {
                    return Tab(text: '${cat.name} (${cat.posterCount})');
                  }),
                ],
              )
            : null,
      ),
      body: RefreshIndicator(
        onRefresh: _onRefresh,
        child: _buildContent(posterState),
      ),
    );
  }

  Widget _buildContent(PosterListState state) {
    if (state.isLoading && state.posters.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (state.error != null && state.posters.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.error_outline, size: 48, color: Colors.grey),
            const SizedBox(height: 16),
            Text(state.error!, style: const TextStyle(color: Colors.grey)),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: _onRefresh,
              child: const Text('重试'),
            ),
          ],
        ),
      );
    }

    if (state.posters.isEmpty) {
      return const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.image_not_supported, size: 48, color: Colors.grey),
            SizedBox(height: 16),
            Text('暂无海报', style: TextStyle(color: Colors.grey)),
          ],
        ),
      );
    }

    return GridView.builder(
      controller: _scrollController,
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        mainAxisSpacing: 16,
        crossAxisSpacing: 16,
        childAspectRatio: 0.65, // 海报比例
      ),
      itemCount: state.posters.length + (state.isLoadingMore ? 1 : 0),
      itemBuilder: (context, index) {
        if (index >= state.posters.length) {
          return const Center(
            child: Padding(
              padding: EdgeInsets.all(16),
              child: CircularProgressIndicator(strokeWidth: 2),
            ),
          );
        }

        final poster = state.posters[index];
        return PosterCard(
          poster: poster,
          onTap: () => _onPosterTap(poster),
        );
      },
    );
  }

  void _onTabChanged() {
    if (_tabController == null || !_tabController!.indexIsChanging) return;

    final categoryState = ref.read(posterCategoryProvider);
    int? categoryId;

    if (_tabController!.index > 0) {
      categoryId = categoryState.categories[_tabController!.index - 1].id;
    }

    ref.read(posterCategoryProvider.notifier).selectCategory(categoryId);
    ref
        .read(posterListProvider.notifier)
        .loadPosters(categoryId: categoryId, refresh: true);
  }

  Future<void> _onRefresh() async {
    final categoryState = ref.read(posterCategoryProvider);
    await ref.read(posterCategoryProvider.notifier).loadCategories();
    await ref
        .read(posterListProvider.notifier)
        .loadPosters(categoryId: categoryState.selectedCategoryId, refresh: true);
  }

  void _onPosterTap(Poster poster) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => PosterDetailPage(poster: poster),
      ),
    );
  }
}
