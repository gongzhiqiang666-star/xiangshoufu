import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';

/// 数据分析页面
class DataAnalysisPage extends StatefulWidget {
  const DataAnalysisPage({super.key});

  @override
  State<DataAnalysisPage> createState() => _DataAnalysisPageState();
}

class _DataAnalysisPageState extends State<DataAnalysisPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  String _selectedPeriod = '本月';

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('数据分析'),
        bottom: TabBar(
          controller: _tabController,
          labelColor: AppColors.primary,
          unselectedLabelColor: AppColors.textSecondary,
          indicatorColor: AppColors.primary,
          tabs: const [
            Tab(text: '交易分析'),
            Tab(text: '商户分析'),
            Tab(text: '代理分析'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildTransactionAnalysis(),
          _buildMerchantAnalysis(),
          _buildAgentAnalysis(),
        ],
      ),
    );
  }

  // 交易分析页
  Widget _buildTransactionAnalysis() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 时间筛选
          _buildPeriodSelector(),
          const SizedBox(height: 16),

          // 交易概览卡片
          _buildTransactionOverview(),
          const SizedBox(height: 16),

          // 交易趋势图
          _buildTransactionTrend(),
          const SizedBox(height: 16),

          // 交易类型占比
          _buildTransactionTypeRatio(),
        ],
      ),
    );
  }

  Widget _buildPeriodSelector() {
    final periods = ['今日', '本周', '本月', '近3月', '近6月'];
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Row(
        children: periods.map((period) {
          final isSelected = _selectedPeriod == period;
          return Padding(
            padding: const EdgeInsets.only(right: 8),
            child: ChoiceChip(
              label: Text(period),
              selected: isSelected,
              selectedColor: AppColors.primary.withOpacity(0.2),
              labelStyle: TextStyle(
                color: isSelected ? AppColors.primary : AppColors.textSecondary,
                fontWeight: isSelected ? FontWeight.w600 : FontWeight.normal,
              ),
              onSelected: (selected) {
                if (selected) {
                  setState(() {
                    _selectedPeriod = period;
                  });
                }
              },
            ),
          );
        }).toList(),
      ),
    );
  }

  Widget _buildTransactionOverview() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '交易概览',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: _buildStatItem('交易总额', '¥128,560.00', '+12.5%', true),
              ),
              Container(width: 1, height: 50, color: AppColors.divider),
              Expanded(
                child: _buildStatItem('交易笔数', '1,256笔', '+8.2%', true),
              ),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: _buildStatItem('日均交易额', '¥4,285.33', '-2.1%', false),
              ),
              Container(width: 1, height: 50, color: AppColors.divider),
              Expanded(
                child: _buildStatItem('笔均金额', '¥102.35', '+5.6%', true),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildStatItem(String label, String value, String change, bool isUp) {
    return Column(
      children: [
        Text(
          value,
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: TextStyle(
            fontSize: 12,
            color: AppColors.textSecondary,
          ),
        ),
        const SizedBox(height: 4),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              isUp ? Icons.arrow_upward : Icons.arrow_downward,
              size: 12,
              color: isUp ? AppColors.success : AppColors.error,
            ),
            Text(
              change,
              style: TextStyle(
                fontSize: 11,
                color: isUp ? AppColors.success : AppColors.error,
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildTransactionTrend() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '交易趋势',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 16),
          // 简化的趋势图（实际应使用fl_chart等图表库）
          SizedBox(
            height: 200,
            child: CustomPaint(
              size: const Size(double.infinity, 200),
              painter: TrendChartPainter(),
            ),
          ),
          const SizedBox(height: 12),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              _buildLegendItem('刷卡交易', AppColors.primary),
              const SizedBox(width: 24),
              _buildLegendItem('扫码交易', AppColors.success),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildLegendItem(String label, Color color) {
    return Row(
      children: [
        Container(
          width: 12,
          height: 12,
          decoration: BoxDecoration(
            color: color,
            borderRadius: BorderRadius.circular(2),
          ),
        ),
        const SizedBox(width: 4),
        Text(
          label,
          style: TextStyle(
            fontSize: 12,
            color: AppColors.textSecondary,
          ),
        ),
      ],
    );
  }

  Widget _buildTransactionTypeRatio() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '交易类型占比',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 16),
          _buildRatioBar('刷卡', 0.65, '¥83,564', AppColors.primary),
          const SizedBox(height: 12),
          _buildRatioBar('微信', 0.20, '¥25,712', const Color(0xFF07C160)),
          const SizedBox(height: 12),
          _buildRatioBar('支付宝', 0.15, '¥19,284', const Color(0xFF1677FF)),
        ],
      ),
    );
  }

  Widget _buildRatioBar(String label, double ratio, String amount, Color color) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              label,
              style: TextStyle(
                fontSize: 14,
                color: AppColors.textPrimary,
              ),
            ),
            Text(
              '$amount (${(ratio * 100).toStringAsFixed(0)}%)',
              style: TextStyle(
                fontSize: 14,
                fontWeight: FontWeight.w500,
                color: AppColors.textPrimary,
              ),
            ),
          ],
        ),
        const SizedBox(height: 6),
        ClipRRect(
          borderRadius: BorderRadius.circular(4),
          child: LinearProgressIndicator(
            value: ratio,
            minHeight: 8,
            backgroundColor: AppColors.background,
            valueColor: AlwaysStoppedAnimation<Color>(color),
          ),
        ),
      ],
    );
  }

  // 商户分析页
  Widget _buildMerchantAnalysis() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 商户分类统计
          _buildMerchantTypeStats(),
          const SizedBox(height: 16),

          // 商户交易排行
          _buildMerchantRanking(),
        ],
      ),
    );
  }

  Widget _buildMerchantTypeStats() {
    final types = [
      {'name': '忠诚商户', 'count': 12, 'color': AppColors.primary},
      {'name': '优质商户', 'count': 28, 'color': AppColors.success},
      {'name': '潜力商户', 'count': 45, 'color': AppColors.warning},
      {'name': '一般商户', 'count': 68, 'color': AppColors.textSecondary},
      {'name': '低活跃', 'count': 23, 'color': AppColors.error},
      {'name': '30天无交易', 'count': 15, 'color': Colors.grey},
    ];

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '商户分类统计',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 16),
          GridView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 3,
              childAspectRatio: 1.2,
              crossAxisSpacing: 12,
              mainAxisSpacing: 12,
            ),
            itemCount: types.length,
            itemBuilder: (context, index) {
              final type = types[index];
              return Container(
                decoration: BoxDecoration(
                  color: (type['color'] as Color).withOpacity(0.1),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                      '${type['count']}',
                      style: TextStyle(
                        fontSize: 24,
                        fontWeight: FontWeight.bold,
                        color: type['color'] as Color,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      type['name'] as String,
                      style: TextStyle(
                        fontSize: 12,
                        color: AppColors.textSecondary,
                      ),
                    ),
                  ],
                ),
              );
            },
          ),
        ],
      ),
    );
  }

  Widget _buildMerchantRanking() {
    final merchants = [
      {'name': '张三商户', 'amount': 25680, 'rank': 1},
      {'name': '李四商户', 'amount': 18920, 'rank': 2},
      {'name': '王五商户', 'amount': 15640, 'rank': 3},
      {'name': '赵六商户', 'amount': 12350, 'rank': 4},
      {'name': '钱七商户', 'amount': 9870, 'rank': 5},
    ];

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                '商户交易排行',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textPrimary,
                ),
              ),
              Text(
                '本月',
                style: TextStyle(
                  fontSize: 12,
                  color: AppColors.textSecondary,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          ...merchants.map((m) => _buildRankingItem(
                m['rank'] as int,
                m['name'] as String,
                m['amount'] as int,
              )),
        ],
      ),
    );
  }

  Widget _buildRankingItem(int rank, String name, int amount) {
    Color rankColor;
    if (rank == 1) {
      rankColor = const Color(0xFFFFD700);
    } else if (rank == 2) {
      rankColor = const Color(0xFFC0C0C0);
    } else if (rank == 3) {
      rankColor = const Color(0xFFCD7F32);
    } else {
      rankColor = AppColors.textTertiary;
    }

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Container(
            width: 24,
            height: 24,
            decoration: BoxDecoration(
              color: rankColor.withOpacity(0.2),
              borderRadius: BorderRadius.circular(4),
            ),
            child: Center(
              child: Text(
                '$rank',
                style: TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.bold,
                  color: rankColor,
                ),
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              name,
              style: TextStyle(
                fontSize: 14,
                color: AppColors.textPrimary,
              ),
            ),
          ),
          Text(
            '¥${(amount / 100).toStringAsFixed(2)}',
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w500,
              color: AppColors.textPrimary,
            ),
          ),
        ],
      ),
    );
  }

  // 代理分析页
  Widget _buildAgentAnalysis() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 代理概览
          _buildAgentOverview(),
          const SizedBox(height: 16),

          // 代理交易排行
          _buildAgentRanking(),
          const SizedBox(height: 16),

          // 终端激活统计
          _buildTerminalStats(),
        ],
      ),
    );
  }

  Widget _buildAgentOverview() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '团队概览',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(child: _buildOverviewItem('直属代理', '15', Icons.person)),
              Expanded(child: _buildOverviewItem('团队代理', '86', Icons.people)),
              Expanded(child: _buildOverviewItem('直营商户', '45', Icons.store)),
              Expanded(child: _buildOverviewItem('团队商户', '328', Icons.storefront)),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildOverviewItem(String label, String value, IconData icon) {
    return Column(
      children: [
        Icon(icon, color: AppColors.primary, size: 24),
        const SizedBox(height: 8),
        Text(
          value,
          style: TextStyle(
            fontSize: 20,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: TextStyle(
            fontSize: 12,
            color: AppColors.textSecondary,
          ),
        ),
      ],
    );
  }

  Widget _buildAgentRanking() {
    final agents = [
      {'name': '代理商A', 'amount': 358600, 'terminals': 45},
      {'name': '代理商B', 'amount': 286400, 'terminals': 38},
      {'name': '代理商C', 'amount': 195200, 'terminals': 28},
      {'name': '代理商D', 'amount': 142800, 'terminals': 22},
      {'name': '代理商E', 'amount': 98500, 'terminals': 15},
    ];

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '代理交易排行',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 12),
          ...agents.asMap().entries.map((entry) {
            final index = entry.key;
            final agent = entry.value;
            return Padding(
              padding: const EdgeInsets.symmetric(vertical: 8),
              child: Row(
                children: [
                  _buildRankBadge(index + 1),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          agent['name'] as String,
                          style: TextStyle(
                            fontSize: 14,
                            fontWeight: FontWeight.w500,
                            color: AppColors.textPrimary,
                          ),
                        ),
                        Text(
                          '${agent['terminals']}台终端',
                          style: TextStyle(
                            fontSize: 12,
                            color: AppColors.textSecondary,
                          ),
                        ),
                      ],
                    ),
                  ),
                  Text(
                    '¥${((agent['amount'] as int) / 100).toStringAsFixed(2)}',
                    style: TextStyle(
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                      color: AppColors.primary,
                    ),
                  ),
                ],
              ),
            );
          }),
        ],
      ),
    );
  }

  Widget _buildRankBadge(int rank) {
    Color color;
    if (rank == 1) {
      color = const Color(0xFFFFD700);
    } else if (rank == 2) {
      color = const Color(0xFFC0C0C0);
    } else if (rank == 3) {
      color = const Color(0xFFCD7F32);
    } else {
      color = AppColors.textTertiary;
    }

    return Container(
      width: 24,
      height: 24,
      decoration: BoxDecoration(
        color: color.withOpacity(0.2),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Center(
        child: Text(
          '$rank',
          style: TextStyle(
            fontSize: 12,
            fontWeight: FontWeight.bold,
            color: color,
          ),
        ),
      ),
    );
  }

  Widget _buildTerminalStats() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '终端激活统计',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: _buildTerminalStatItem('总数', '520', AppColors.textPrimary),
              ),
              Expanded(
                child: _buildTerminalStatItem('已激活', '385', AppColors.success),
              ),
              Expanded(
                child: _buildTerminalStatItem('未激活', '135', AppColors.warning),
              ),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: _buildTerminalStatItem('今日激活', '8', AppColors.primary),
              ),
              Expanded(
                child: _buildTerminalStatItem('本月激活', '45', AppColors.primary),
              ),
              Expanded(
                child: _buildTerminalStatItem('激活率', '74%', AppColors.success),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildTerminalStatItem(String label, String value, Color color) {
    return Column(
      children: [
        Text(
          value,
          style: TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
            color: color,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: TextStyle(
            fontSize: 12,
            color: AppColors.textSecondary,
          ),
        ),
      ],
    );
  }
}

// 简化的趋势图绘制器
class TrendChartPainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..strokeWidth = 2
      ..style = PaintingStyle.stroke;

    // 绘制网格线
    final gridPaint = Paint()
      ..color = Colors.grey.withOpacity(0.2)
      ..strokeWidth = 1;

    for (var i = 0; i <= 4; i++) {
      final y = size.height * i / 4;
      canvas.drawLine(Offset(0, y), Offset(size.width, y), gridPaint);
    }

    // 绘制刷卡交易曲线
    paint.color = AppColors.primary;
    final path1 = Path();
    final points1 = [0.3, 0.5, 0.4, 0.6, 0.55, 0.7, 0.65];
    for (var i = 0; i < points1.length; i++) {
      final x = size.width * i / (points1.length - 1);
      final y = size.height * (1 - points1[i]);
      if (i == 0) {
        path1.moveTo(x, y);
      } else {
        path1.lineTo(x, y);
      }
    }
    canvas.drawPath(path1, paint);

    // 绘制扫码交易曲线
    paint.color = AppColors.success;
    final path2 = Path();
    final points2 = [0.2, 0.25, 0.3, 0.35, 0.4, 0.38, 0.45];
    for (var i = 0; i < points2.length; i++) {
      final x = size.width * i / (points2.length - 1);
      final y = size.height * (1 - points2[i]);
      if (i == 0) {
        path2.moveTo(x, y);
      } else {
        path2.lineTo(x, y);
      }
    }
    canvas.drawPath(path2, paint);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
