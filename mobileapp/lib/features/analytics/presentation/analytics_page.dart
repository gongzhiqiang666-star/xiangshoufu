import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';

import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../core/utils/format_utils.dart';
import '../../home/data/home_api.dart';

/// 数据分析页面
class AnalyticsPage extends ConsumerStatefulWidget {
  const AnalyticsPage({super.key});

  @override
  ConsumerState<AnalyticsPage> createState() => _AnalyticsPageState();
}

class _AnalyticsPageState extends ConsumerState<AnalyticsPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  String _selectedPeriod = 'month'; // day, week, month
  String _scope = 'direct'; // direct, team

  bool _isLoading = true;
  List<TrendPoint> _trendData = [];
  List<ChannelStats> _channelStats = [];
  List<MerchantDistribution> _merchantDistribution = [];

  int _totalTransAmount = 0;
  int _totalProfitTotal = 0;
  int _totalTransCount = 0;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    _tabController.addListener(_onTabChanged);
    _loadData();
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  void _onTabChanged() {
    if (_tabController.indexIsChanging) {
      setState(() {
        _scope = _tabController.index == 0 ? 'direct' : 'team';
      });
      _loadData();
    }
  }

  Future<void> _loadData() async {
    setState(() => _isLoading = true);

    try {
      final api = ref.read(homeApiProvider);

      int days = _selectedPeriod == 'day' ? 1 : (_selectedPeriod == 'week' ? 7 : 30);

      final results = await Future.wait([
        api.getTrendData(days: days, scope: _scope),
        api.getChannelStats(scope: _scope, period: _selectedPeriod),
        api.getMerchantDistribution(scope: _scope),
      ]);

      _trendData = results[0] as List<TrendPoint>;
      _channelStats = results[1] as List<ChannelStats>;
      _merchantDistribution = results[2] as List<MerchantDistribution>;

      // 计算汇总
      _totalTransAmount = 0;
      _totalProfitTotal = 0;
      _totalTransCount = 0;
      for (var t in _trendData) {
        _totalTransAmount += t.transAmount;
        _totalProfitTotal += t.profitTotal;
        _totalTransCount += t.transCount;
      }

      setState(() => _isLoading = false);
    } catch (e) {
      setState(() => _isLoading = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('加载失败: $e')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('数据分析'),
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: '直营'),
            Tab(text: '团队'),
          ],
        ),
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: _loadData,
              child: SingleChildScrollView(
                physics: const AlwaysScrollableScrollPhysics(),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _buildPeriodSelector(),
                    const SizedBox(height: AppSpacing.md),
                    _buildSummaryCards(),
                    const SizedBox(height: AppSpacing.md),
                    _buildTrendChart(),
                    const SizedBox(height: AppSpacing.md),
                    _buildChannelPieChart(),
                    const SizedBox(height: AppSpacing.md),
                    _buildMerchantDistributionChart(),
                    const SizedBox(height: AppSpacing.lg),
                  ],
                ),
              ),
            ),
    );
  }

  Widget _buildPeriodSelector() {
    return Container(
      margin: const EdgeInsets.all(AppSpacing.md),
      child: Row(
        children: [
          _buildPeriodButton('今日', 'day'),
          const SizedBox(width: 8),
          _buildPeriodButton('本周', 'week'),
          const SizedBox(width: 8),
          _buildPeriodButton('本月', 'month'),
        ],
      ),
    );
  }

  Widget _buildPeriodButton(String label, String period) {
    final isSelected = _selectedPeriod == period;
    return Expanded(
      child: GestureDetector(
        onTap: () {
          setState(() => _selectedPeriod = period);
          _loadData();
        },
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 10),
          decoration: BoxDecoration(
            color: isSelected ? AppColors.primary : Colors.white,
            borderRadius: BorderRadius.circular(8),
            border: Border.all(
              color: isSelected ? AppColors.primary : AppColors.border,
            ),
          ),
          child: Text(
            label,
            textAlign: TextAlign.center,
            style: TextStyle(
              color: isSelected ? Colors.white : AppColors.textSecondary,
              fontWeight: isSelected ? FontWeight.w600 : FontWeight.normal,
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildSummaryCards() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      child: Row(
        children: [
          Expanded(
            child: _buildSummaryCard(
              '交易统计',
              FormatUtils.formatYuan(_totalTransAmount / 100),
              '$_totalTransCount笔',
              AppColors.primary,
            ),
          ),
          const SizedBox(width: AppSpacing.cardGap),
          Expanded(
            child: _buildSummaryCard(
              '收益统计',
              FormatUtils.formatYuan(_totalProfitTotal / 100),
              '总收益',
              AppColors.success,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSummaryCard(String title, String value, String subtitle, Color color) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(title, style: TextStyle(fontSize: 13, color: AppColors.textSecondary)),
          const SizedBox(height: 8),
          Text(value, style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: color)),
          const SizedBox(height: 4),
          Text(subtitle, style: TextStyle(fontSize: 12, color: AppColors.textTertiary)),
        ],
      ),
    );
  }

  Widget _buildTrendChart() {
    if (_trendData.isEmpty) {
      return _buildEmptyChart('暂无趋势数据');
    }

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('交易趋势', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600)),
          const SizedBox(height: 16),
          SizedBox(
            height: 200,
            child: LineChart(
              LineChartData(
                gridData: FlGridData(show: false),
                titlesData: FlTitlesData(
                  leftTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  rightTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  topTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  bottomTitles: AxisTitles(
                    sideTitles: SideTitles(
                      showTitles: true,
                      getTitlesWidget: (value, meta) {
                        final index = value.toInt();
                        if (index >= 0 && index < _trendData.length) {
                          return Padding(
                            padding: const EdgeInsets.only(top: 8),
                            child: Text(
                              _trendData[index].date,
                              style: const TextStyle(fontSize: 10, color: AppColors.textTertiary),
                            ),
                          );
                        }
                        return const Text('');
                      },
                      reservedSize: 30,
                    ),
                  ),
                ),
                borderData: FlBorderData(show: false),
                lineBarsData: [
                  LineChartBarData(
                    spots: _trendData.asMap().entries.map((e) {
                      return FlSpot(e.key.toDouble(), e.value.transAmountYuan);
                    }).toList(),
                    isCurved: true,
                    color: AppColors.primary,
                    barWidth: 2,
                    dotData: FlDotData(show: false),
                    belowBarData: BarAreaData(
                      show: true,
                      color: AppColors.primary.withOpacity(0.1),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildChannelPieChart() {
    if (_channelStats.isEmpty) {
      return _buildEmptyChart('暂无通道数据');
    }

    final colors = [
      AppColors.primary,
      AppColors.success,
      AppColors.warning,
      AppColors.info,
      AppColors.profitTrade,
      AppColors.profitDeposit,
    ];

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('通道占比', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600)),
          const SizedBox(height: 16),
          SizedBox(
            height: 200,
            child: Row(
              children: [
                Expanded(
                  child: PieChart(
                    PieChartData(
                      sectionsSpace: 2,
                      centerSpaceRadius: 40,
                      sections: _channelStats.asMap().entries.map((e) {
                        final color = colors[e.key % colors.length];
                        return PieChartSectionData(
                          value: e.value.percentage,
                          color: color,
                          title: '${e.value.percentage.toStringAsFixed(1)}%',
                          titleStyle: const TextStyle(fontSize: 10, color: Colors.white, fontWeight: FontWeight.bold),
                          radius: 50,
                        );
                      }).toList(),
                    ),
                  ),
                ),
                const SizedBox(width: 16),
                Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: _channelStats.asMap().entries.map((e) {
                    final color = colors[e.key % colors.length];
                    return Padding(
                      padding: const EdgeInsets.symmetric(vertical: 4),
                      child: Row(
                        children: [
                          Container(
                            width: 12,
                            height: 12,
                            decoration: BoxDecoration(
                              color: color,
                              borderRadius: BorderRadius.circular(2),
                            ),
                          ),
                          const SizedBox(width: 8),
                          Text(
                            e.value.channelName,
                            style: const TextStyle(fontSize: 12, color: AppColors.textSecondary),
                          ),
                        ],
                      ),
                    );
                  }).toList(),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMerchantDistributionChart() {
    if (_merchantDistribution.isEmpty) {
      return _buildEmptyChart('暂无商户分布数据');
    }

    final colors = [
      AppColors.success,      // 忠诚
      AppColors.primary,      // 优质
      AppColors.info,         // 潜力
      AppColors.warning,      // 一般
      AppColors.textTertiary, // 低活跃
      AppColors.danger,       // 无交易
    ];

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('商户分布', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600)),
          const SizedBox(height: 16),
          SizedBox(
            height: 200,
            child: BarChart(
              BarChartData(
                alignment: BarChartAlignment.spaceAround,
                maxY: _merchantDistribution.map((e) => e.count.toDouble()).reduce((a, b) => a > b ? a : b) * 1.2,
                barTouchData: BarTouchData(enabled: false),
                titlesData: FlTitlesData(
                  leftTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  rightTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  topTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  bottomTitles: AxisTitles(
                    sideTitles: SideTitles(
                      showTitles: true,
                      getTitlesWidget: (value, meta) {
                        final index = value.toInt();
                        if (index >= 0 && index < _merchantDistribution.length) {
                          return Padding(
                            padding: const EdgeInsets.only(top: 8),
                            child: Text(
                              _merchantDistribution[index].typeName.length > 2
                                  ? _merchantDistribution[index].typeName.substring(0, 2)
                                  : _merchantDistribution[index].typeName,
                              style: const TextStyle(fontSize: 10, color: AppColors.textTertiary),
                            ),
                          );
                        }
                        return const Text('');
                      },
                      reservedSize: 30,
                    ),
                  ),
                ),
                gridData: FlGridData(show: false),
                borderData: FlBorderData(show: false),
                barGroups: _merchantDistribution.asMap().entries.map((e) {
                  final color = colors[e.key % colors.length];
                  return BarChartGroupData(
                    x: e.key,
                    barRods: [
                      BarChartRodData(
                        toY: e.value.count.toDouble(),
                        color: color,
                        width: 20,
                        borderRadius: const BorderRadius.only(
                          topLeft: Radius.circular(4),
                          topRight: Radius.circular(4),
                        ),
                      ),
                    ],
                  );
                }).toList(),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildEmptyChart(String message) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      padding: const EdgeInsets.all(32),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Center(
        child: Text(message, style: const TextStyle(color: AppColors.textTertiary)),
      ),
    );
  }
}
