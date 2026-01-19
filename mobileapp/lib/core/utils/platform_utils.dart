import 'dart:io';

/// 平台工具类
/// 用于判断当前运行平台，适配iOS/Android/HarmonyOS
class PlatformUtils {
  PlatformUtils._();

  /// 是否iOS
  static bool get isIOS => Platform.isIOS;

  /// 是否Android
  static bool get isAndroid => Platform.isAndroid;

  /// 是否HarmonyOS（鸿蒙）
  /// 注：HarmonyOS基于Android，需要额外判断
  static bool get isHarmonyOS {
    if (!Platform.isAndroid) return false;
    // 通过环境变量或系统属性判断
    // 实际项目中需要通过原生代码获取
    return Platform.environment.containsKey('HARMONYOS_VERSION') ||
        Platform.operatingSystemVersion.toLowerCase().contains('harmony');
  }

  /// 获取平台名称
  static String get platformName {
    if (isIOS) return 'iOS';
    if (isHarmonyOS) return 'HarmonyOS';
    if (isAndroid) return 'Android';
    return 'Unknown';
  }

  /// 是否移动端
  static bool get isMobile => isIOS || isAndroid;

  /// 是否桌面端
  static bool get isDesktop =>
      Platform.isWindows || Platform.isMacOS || Platform.isLinux;

  /// 获取操作系统版本
  static String get osVersion => Platform.operatingSystemVersion;

  /// 获取Dart版本
  static String get dartVersion => Platform.version;
}

/// 设备信息工具类
class DeviceUtils {
  DeviceUtils._();

  /// 是否是刘海屏/挖孔屏
  /// 需要配合MediaQuery判断
  static bool hasNotch(double topPadding) {
    return topPadding > 24;
  }

  /// 是否有底部安全区域（如iPhone X系列）
  static bool hasBottomSafeArea(double bottomPadding) {
    return bottomPadding > 0;
  }
}
