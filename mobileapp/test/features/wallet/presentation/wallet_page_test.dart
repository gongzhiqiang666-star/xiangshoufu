import 'package:flutter_test/flutter_test.dart';
import 'package:xiangshoufu_app/features/wallet/data/models/wallet_model.dart';

/// 钱包页面测试
/// 测试场景：
/// 1. 默认钱包显示（API返回空时显示默认钱包）
/// 2. 通道筛选功能
void main() {
  group('WalletPage - 默认钱包显示', () {
    test('当API返回空列表时，应返回3个默认钱包', () {
      // Arrange
      final emptyWallets = <WalletModel>[];

      // Act
      final defaultWallets = _getDefaultWallets();

      // Assert
      expect(defaultWallets.length, 3);
      expect(defaultWallets[0].walletType, 1); // 分润钱包
      expect(defaultWallets[1].walletType, 2); // 服务费钱包
      expect(defaultWallets[2].walletType, 3); // 奖励钱包
    });

    test('默认钱包余额应为0', () {
      // Act
      final defaultWallets = _getDefaultWallets();

      // Assert
      for (final wallet in defaultWallets) {
        expect(wallet.balance, 0);
        expect(wallet.available, 0);
        expect(wallet.frozen, 0);
      }
    });

    test('默认钱包通道名称应为"通用"', () {
      // Act
      final defaultWallets = _getDefaultWallets();

      // Assert
      for (final wallet in defaultWallets) {
        expect(wallet.channelName, '通用');
      }
    });
  });

  group('WalletPage - 通道筛选', () {
    test('从钱包列表提取通道名称，应包含"全部"', () {
      // Arrange
      final wallets = [
        _createWallet(channelName: '恒信通'),
        _createWallet(channelName: '拉卡拉'),
        _createWallet(channelName: '恒信通'), // 重复
      ];

      // Act
      final channelSet = <String>{'全部'};
      for (final wallet in wallets) {
        if (wallet.channelName.isNotEmpty) {
          channelSet.add(wallet.channelName);
        }
      }
      final channels = channelSet.toList();

      // Assert
      expect(channels.contains('全部'), true);
      expect(channels.contains('恒信通'), true);
      expect(channels.contains('拉卡拉'), true);
      expect(channels.length, 3); // 全部 + 恒信通 + 拉卡拉（去重）
    });

    test('按通道筛选钱包列表', () {
      // Arrange
      final wallets = [
        _createWallet(channelName: '恒信通', walletType: 1),
        _createWallet(channelName: '拉卡拉', walletType: 1),
        _createWallet(channelName: '恒信通', walletType: 2),
      ];
      const selectedChannel = '恒信通';

      // Act
      final filteredWallets = selectedChannel == '全部'
          ? wallets
          : wallets.where((w) => w.channelName == selectedChannel).toList();

      // Assert
      expect(filteredWallets.length, 2);
      expect(filteredWallets.every((w) => w.channelName == '恒信通'), true);
    });

    test('选择"全部"时应返回所有钱包', () {
      // Arrange
      final wallets = [
        _createWallet(channelName: '恒信通'),
        _createWallet(channelName: '拉卡拉'),
      ];
      const selectedChannel = '全部';

      // Act
      final filteredWallets = selectedChannel == '全部'
          ? wallets
          : wallets.where((w) => w.channelName == selectedChannel).toList();

      // Assert
      expect(filteredWallets.length, 2);
    });

    test('当API返回空列表时，使用默认钱包进行筛选', () {
      // Arrange
      final apiWallets = <WalletModel>[];
      final displayWallets = apiWallets.isEmpty ? _getDefaultWallets() : apiWallets;

      // Act & Assert
      expect(displayWallets.length, 3);
      expect(displayWallets[0].channelName, '通用');
    });
  });

  group('WalletModel - 金额转换', () {
    test('分转元计算正确', () {
      // Arrange
      final wallet = _createWallet(balance: 12345);

      // Assert
      expect(wallet.balanceYuan, 123.45);
    });

    test('walletTypeName 返回正确名称', () {
      // Assert
      expect(_createWallet(walletType: 1).walletTypeName, '分润钱包');
      expect(_createWallet(walletType: 2).walletTypeName, '服务费钱包');
      expect(_createWallet(walletType: 3).walletTypeName, '奖励钱包');
      expect(_createWallet(walletType: 4).walletTypeName, '充值钱包');
      expect(_createWallet(walletType: 5).walletTypeName, '沉淀钱包');
      // 注意：walletType 默认值为1，所以无效值也返回分润钱包
    });
  });
}

/// 获取默认钱包列表（模拟 wallet_page.dart 中的 _getDefaultWallets）
List<WalletModel> _getDefaultWallets() {
  return [
    WalletModel(
      id: 0,
      agentId: 0,
      agentName: '',
      walletType: 1,
      channelId: 0,
      channelName: '通用',
      balance: 0,
      available: 0,
      frozen: 0,
      totalIncome: 0,
      totalWithdraw: 0,
      updatedAt: '',
    ),
    WalletModel(
      id: 0,
      agentId: 0,
      agentName: '',
      walletType: 2,
      channelId: 0,
      channelName: '通用',
      balance: 0,
      available: 0,
      frozen: 0,
      totalIncome: 0,
      totalWithdraw: 0,
      updatedAt: '',
    ),
    WalletModel(
      id: 0,
      agentId: 0,
      agentName: '',
      walletType: 3,
      channelId: 0,
      channelName: '通用',
      balance: 0,
      available: 0,
      frozen: 0,
      totalIncome: 0,
      totalWithdraw: 0,
      updatedAt: '',
    ),
  ];
}

/// 创建测试用钱包
WalletModel _createWallet({
  int walletType = 1,
  String channelName = '通用',
  int balance = 0,
}) {
  return WalletModel(
    id: 1,
    agentId: 1,
    agentName: 'test',
    walletType: walletType,
    channelId: 1,
    channelName: channelName,
    balance: balance,
    available: balance,
    frozen: 0,
    totalIncome: 0,
    totalWithdraw: 0,
    updatedAt: '2024-01-01',
  );
}
