/// 营销海报模型
class Poster {
  final int id;
  final String title;
  final int categoryId;
  final String imageUrl;
  final String? thumbnailUrl;
  final String? description;
  final int width;
  final int height;

  Poster({
    required this.id,
    required this.title,
    required this.categoryId,
    required this.imageUrl,
    this.thumbnailUrl,
    this.description,
    this.width = 0,
    this.height = 0,
  });

  factory Poster.fromJson(Map<String, dynamic> json) {
    return Poster(
      id: json['id'] as int,
      title: json['title'] as String,
      categoryId: json['category_id'] as int,
      imageUrl: json['image_url'] as String,
      thumbnailUrl: json['thumbnail_url'] as String?,
      description: json['description'] as String?,
      width: json['width'] as int? ?? 0,
      height: json['height'] as int? ?? 0,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'category_id': categoryId,
      'image_url': imageUrl,
      'thumbnail_url': thumbnailUrl,
      'description': description,
      'width': width,
      'height': height,
    };
  }

  /// 获取显示用的图片URL（优先使用缩略图）
  String get displayImageUrl => thumbnailUrl ?? imageUrl;

  /// 计算宽高比
  double get aspectRatio {
    if (width > 0 && height > 0) {
      return width / height;
    }
    return 0.67; // 默认2:3比例
  }
}
