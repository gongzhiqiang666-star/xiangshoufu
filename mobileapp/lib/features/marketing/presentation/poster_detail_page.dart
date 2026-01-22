import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:image_gallery_saver/image_gallery_saver.dart';
import 'package:share_plus/share_plus.dart';
import 'package:dio/dio.dart';
import 'package:path_provider/path_provider.dart';
import 'dart:io';
import '../../data/models/poster.dart';
import '../providers/poster_provider.dart';

/// 海报详情页
class PosterDetailPage extends ConsumerStatefulWidget {
  final Poster poster;

  const PosterDetailPage({
    super.key,
    required this.poster,
  });

  @override
  ConsumerState<PosterDetailPage> createState() => _PosterDetailPageState();
}

class _PosterDetailPageState extends ConsumerState<PosterDetailPage> {
  bool _isSaving = false;
  bool _isSharing = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        iconTheme: const IconThemeData(color: Colors.white),
        title: Text(
          widget.poster.title,
          style: const TextStyle(color: Colors.white),
        ),
      ),
      body: Column(
        children: [
          // 图片预览区域（支持缩放）
          Expanded(
            child: InteractiveViewer(
              minScale: 0.5,
              maxScale: 4.0,
              child: Center(
                child: CachedNetworkImage(
                  imageUrl: widget.poster.imageUrl,
                  fit: BoxFit.contain,
                  placeholder: (context, url) => const Center(
                    child: CircularProgressIndicator(
                      color: Colors.white,
                    ),
                  ),
                  errorWidget: (context, url, error) => const Icon(
                    Icons.error,
                    color: Colors.white,
                    size: 48,
                  ),
                ),
              ),
            ),
          ),
          // 底部操作栏
          _buildBottomBar(),
        ],
      ),
    );
  }

  Widget _buildBottomBar() {
    return Container(
      padding: EdgeInsets.only(
        left: 24,
        right: 24,
        top: 16,
        bottom: MediaQuery.of(context).padding.bottom + 16,
      ),
      decoration: BoxDecoration(
        color: Colors.black.withOpacity(0.8),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          // 保存到相册
          _buildActionButton(
            icon: Icons.download,
            label: '保存到相册',
            isLoading: _isSaving,
            onTap: _saveToGallery,
          ),
          const SizedBox(width: 32),
          // 分享
          _buildActionButton(
            icon: Icons.share,
            label: '分享',
            isLoading: _isSharing,
            onTap: _shareImage,
          ),
        ],
      ),
    );
  }

  Widget _buildActionButton({
    required IconData icon,
    required String label,
    required bool isLoading,
    required VoidCallback onTap,
  }) {
    return InkWell(
      onTap: isLoading ? null : onTap,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 56,
            height: 56,
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.1),
              shape: BoxShape.circle,
            ),
            child: isLoading
                ? const Center(
                    child: SizedBox(
                      width: 24,
                      height: 24,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        color: Colors.white,
                      ),
                    ),
                  )
                : Icon(icon, color: Colors.white, size: 28),
          ),
          const SizedBox(height: 8),
          Text(
            label,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 12,
            ),
          ),
        ],
      ),
    );
  }

  /// 保存到相册
  Future<void> _saveToGallery() async {
    if (_isSaving) return;

    setState(() => _isSaving = true);

    try {
      // 下载图片
      final dio = Dio();
      final response = await dio.get(
        widget.poster.imageUrl,
        options: Options(responseType: ResponseType.bytes),
      );

      // 保存到相册
      final result = await ImageGallerySaver.saveImage(
        Uint8List.fromList(response.data),
        quality: 100,
        name: 'poster_${widget.poster.id}_${DateTime.now().millisecondsSinceEpoch}',
      );

      if (result['isSuccess'] == true) {
        // 记录下载
        ref.read(posterListProvider.notifier).recordDownload(widget.poster.id);

        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('已保存到相册'),
              backgroundColor: Colors.green,
            ),
          );
        }
      } else {
        throw Exception('保存失败');
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('保存失败: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isSaving = false);
      }
    }
  }

  /// 分享图片
  Future<void> _shareImage() async {
    if (_isSharing) return;

    setState(() => _isSharing = true);

    try {
      // 下载图片到临时目录
      final dio = Dio();
      final tempDir = await getTemporaryDirectory();
      final fileName = 'poster_${widget.poster.id}.jpg';
      final filePath = '${tempDir.path}/$fileName';

      await dio.download(widget.poster.imageUrl, filePath);

      // 分享
      await Share.shareXFiles(
        [XFile(filePath)],
        text: widget.poster.title,
      );

      // 记录分享
      ref.read(posterListProvider.notifier).recordShare(widget.poster.id);
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('分享失败: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isSharing = false);
      }
    }
  }
}
