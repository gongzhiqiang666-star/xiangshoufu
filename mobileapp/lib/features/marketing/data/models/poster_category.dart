/// 海报分类模型
class PosterCategory {
  final int id;
  final String name;
  final int posterCount;

  PosterCategory({
    required this.id,
    required this.name,
    required this.posterCount,
  });

  factory PosterCategory.fromJson(Map<String, dynamic> json) {
    return PosterCategory(
      id: json['id'] as int,
      name: json['name'] as String,
      posterCount: json['poster_count'] as int? ?? 0,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'poster_count': posterCount,
    };
  }
}
