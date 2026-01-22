/// Banner 滚动图模型
class Banner {
  final int id;
  final String title;
  final String imageUrl;
  final int linkType; // 0无链接 1内部页面 2外部链接
  final String? linkUrl;

  Banner({
    required this.id,
    required this.title,
    required this.imageUrl,
    required this.linkType,
    this.linkUrl,
  });

  factory Banner.fromJson(Map<String, dynamic> json) {
    return Banner(
      id: json['id'] as int,
      title: json['title'] as String,
      imageUrl: json['image_url'] as String,
      linkType: json['link_type'] as int? ?? 0,
      linkUrl: json['link_url'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'image_url': imageUrl,
      'link_type': linkType,
      'link_url': linkUrl,
    };
  }

  /// 是否有跳转链接
  bool get hasLink => linkType != 0 && linkUrl != null && linkUrl!.isNotEmpty;

  /// 是否为外部链接
  bool get isExternalLink => linkType == 2;

  /// 是否为内部页面
  bool get isInternalLink => linkType == 1;
}
