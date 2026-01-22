import 'package:flutter/material.dart' hide Banner;
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:carousel_slider/carousel_slider.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:url_launcher/url_launcher.dart';
import '../providers/banner_provider.dart';
import '../../data/models/banner.dart';

/// Banner轮播组件
class BannerCarousel extends ConsumerStatefulWidget {
  /// 轮播高度
  final double height;

  /// 自动播放间隔（毫秒）
  final int autoPlayInterval;

  /// 点击回调（用于内部页面跳转）
  final void Function(String route)? onInternalLinkTap;

  const BannerCarousel({
    super.key,
    this.height = 180,
    this.autoPlayInterval = 5000,
    this.onInternalLinkTap,
  });

  @override
  ConsumerState<BannerCarousel> createState() => _BannerCarouselState();
}

class _BannerCarouselState extends ConsumerState<BannerCarousel> {
  int _currentIndex = 0;

  @override
  void initState() {
    super.initState();
    // 加载Banner数据
    Future.microtask(() {
      ref.read(bannerListProvider.notifier).loadBanners();
    });
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(bannerListProvider);

    if (state.isLoading) {
      return _buildPlaceholder();
    }

    if (state.banners.isEmpty) {
      return _buildPlaceholder();
    }

    return Column(
      children: [
        CarouselSlider(
          options: CarouselOptions(
            height: widget.height,
            viewportFraction: 1.0,
            autoPlay: state.banners.length > 1,
            autoPlayInterval: Duration(milliseconds: widget.autoPlayInterval),
            autoPlayAnimationDuration: const Duration(milliseconds: 800),
            autoPlayCurve: Curves.fastOutSlowIn,
            onPageChanged: (index, reason) {
              setState(() {
                _currentIndex = index;
              });
            },
          ),
          items: state.banners.map((banner) {
            return _buildBannerItem(banner);
          }).toList(),
        ),
        if (state.banners.length > 1) _buildIndicator(state.banners.length),
      ],
    );
  }

  /// 构建Banner项
  Widget _buildBannerItem(Banner banner) {
    return GestureDetector(
      onTap: () => _handleBannerTap(banner),
      child: Container(
        width: double.infinity,
        margin: const EdgeInsets.symmetric(horizontal: 16),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(12),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.1),
              blurRadius: 8,
              offset: const Offset(0, 4),
            ),
          ],
        ),
        child: ClipRRect(
          borderRadius: BorderRadius.circular(12),
          child: CachedNetworkImage(
            imageUrl: banner.imageUrl,
            fit: BoxFit.cover,
            placeholder: (context, url) => Container(
              color: Colors.grey[200],
              child: const Center(
                child: CircularProgressIndicator(strokeWidth: 2),
              ),
            ),
            errorWidget: (context, url, error) => Container(
              color: Colors.grey[200],
              child: const Icon(Icons.image_not_supported, size: 40),
            ),
          ),
        ),
      ),
    );
  }

  /// 构建指示器
  Widget _buildIndicator(int count) {
    return Container(
      padding: const EdgeInsets.only(top: 12),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: List.generate(count, (index) {
          return AnimatedContainer(
            duration: const Duration(milliseconds: 300),
            margin: const EdgeInsets.symmetric(horizontal: 4),
            width: _currentIndex == index ? 20 : 8,
            height: 8,
            decoration: BoxDecoration(
              color: _currentIndex == index
                  ? Theme.of(context).primaryColor
                  : Colors.grey[300],
              borderRadius: BorderRadius.circular(4),
            ),
          );
        }),
      ),
    );
  }

  /// 构建占位图
  Widget _buildPlaceholder() {
    return Container(
      height: widget.height,
      margin: const EdgeInsets.symmetric(horizontal: 16),
      decoration: BoxDecoration(
        color: Colors.grey[200],
        borderRadius: BorderRadius.circular(12),
      ),
      child: const Center(
        child: Icon(
          Icons.image,
          size: 48,
          color: Colors.grey,
        ),
      ),
    );
  }

  /// 处理Banner点击
  Future<void> _handleBannerTap(Banner banner) async {
    // 记录点击
    ref.read(bannerListProvider.notifier).recordClick(banner.id);

    if (!banner.hasLink) return;

    if (banner.isInternalLink && widget.onInternalLinkTap != null) {
      // 内部页面跳转
      widget.onInternalLinkTap!(banner.linkUrl!);
    } else if (banner.isExternalLink) {
      // 外部链接跳转
      final uri = Uri.parse(banner.linkUrl!);
      if (await canLaunchUrl(uri)) {
        await launchUrl(uri, mode: LaunchMode.externalApplication);
      }
    }
  }
}
