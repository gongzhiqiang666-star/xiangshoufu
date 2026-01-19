/// 表单验证工具类
class Validators {
  Validators._();

  /// 手机号验证
  static String? phone(String? value) {
    if (value == null || value.isEmpty) {
      return '请输入手机号';
    }
    if (!RegExp(r'^1[3-9]\d{9}$').hasMatch(value)) {
      return '请输入正确的手机号';
    }
    return null;
  }

  /// 验证码验证
  static String? verifyCode(String? value, {int length = 6}) {
    if (value == null || value.isEmpty) {
      return '请输入验证码';
    }
    if (value.length != length) {
      return '请输入${length}位验证码';
    }
    if (!RegExp(r'^\d+$').hasMatch(value)) {
      return '验证码只能是数字';
    }
    return null;
  }

  /// 身份证号验证
  static String? idCard(String? value) {
    if (value == null || value.isEmpty) {
      return '请输入身份证号';
    }
    if (!RegExp(r'^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$')
        .hasMatch(value)) {
      return '请输入正确的身份证号';
    }
    return null;
  }

  /// 银行卡号验证
  static String? bankCard(String? value) {
    if (value == null || value.isEmpty) {
      return '请输入银行卡号';
    }
    // 移除空格
    final cardNo = value.replaceAll(' ', '');
    if (cardNo.length < 15 || cardNo.length > 19) {
      return '请输入正确的银行卡号';
    }
    if (!RegExp(r'^\d+$').hasMatch(cardNo)) {
      return '银行卡号只能是数字';
    }
    return null;
  }

  /// 金额验证
  static String? amount(String? value, {double? min, double? max}) {
    if (value == null || value.isEmpty) {
      return '请输入金额';
    }
    final amount = double.tryParse(value);
    if (amount == null) {
      return '请输入正确的金额';
    }
    if (amount <= 0) {
      return '金额必须大于0';
    }
    if (min != null && amount < min) {
      return '金额不能小于$min';
    }
    if (max != null && amount > max) {
      return '金额不能大于$max';
    }
    return null;
  }

  /// 费率验证（0.00% - 100.00%）
  static String? rate(String? value, {double? min, double? max}) {
    if (value == null || value.isEmpty) {
      return '请输入费率';
    }
    // 移除%符号
    final rateStr = value.replaceAll('%', '').trim();
    final rate = double.tryParse(rateStr);
    if (rate == null) {
      return '请输入正确的费率';
    }
    if (rate < 0 || rate > 100) {
      return '费率范围为0-100%';
    }
    if (min != null && rate < min) {
      return '费率不能小于$min%';
    }
    if (max != null && rate > max) {
      return '费率不能大于$max%';
    }
    return null;
  }

  /// 必填验证
  static String? required(String? value, {String? fieldName}) {
    if (value == null || value.trim().isEmpty) {
      return fieldName != null ? '请输入$fieldName' : '此项为必填';
    }
    return null;
  }

  /// 长度验证
  static String? length(String? value, {int? min, int? max, String? fieldName}) {
    if (value == null || value.isEmpty) {
      return null; // 空值交给required验证
    }
    if (min != null && value.length < min) {
      return '${fieldName ?? '内容'}不能少于$min个字符';
    }
    if (max != null && value.length > max) {
      return '${fieldName ?? '内容'}不能超过$max个字符';
    }
    return null;
  }

  /// SN号验证
  static String? snCode(String? value) {
    if (value == null || value.isEmpty) {
      return '请输入SN号';
    }
    // 移除空格
    final sn = value.replaceAll(' ', '');
    if (sn.length < 8 || sn.length > 20) {
      return '请输入正确的SN号';
    }
    if (!RegExp(r'^[A-Za-z0-9]+$').hasMatch(sn)) {
      return 'SN号只能包含字母和数字';
    }
    return null;
  }

  /// 邀请码验证
  static String? inviteCode(String? value) {
    if (value == null || value.isEmpty) {
      return '请输入邀请码';
    }
    if (value.length < 4 || value.length > 20) {
      return '邀请码长度为4-20位';
    }
    if (!RegExp(r'^[A-Za-z0-9]+$').hasMatch(value)) {
      return '邀请码只能包含字母和数字';
    }
    return null;
  }

  /// 密码验证
  static String? password(String? value, {int minLength = 6, int maxLength = 20}) {
    if (value == null || value.isEmpty) {
      return '请输入密码';
    }
    if (value.length < minLength) {
      return '密码不能少于$minLength位';
    }
    if (value.length > maxLength) {
      return '密码不能超过$maxLength位';
    }
    return null;
  }

  /// 确认密码验证
  static String? confirmPassword(String? value, String? password) {
    if (value == null || value.isEmpty) {
      return '请确认密码';
    }
    if (value != password) {
      return '两次输入的密码不一致';
    }
    return null;
  }

  /// 组合多个验证器
  static String? combine(String? value, List<String? Function(String?)> validators) {
    for (final validator in validators) {
      final result = validator(value);
      if (result != null) {
        return result;
      }
    }
    return null;
  }
}
