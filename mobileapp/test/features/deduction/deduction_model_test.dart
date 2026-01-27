/// 代扣模块数据模型测试
import 'package:flutter_test/flutter_test.dart';
import 'package:xiangshoufu_app/features/deduction/data/models/deduction_model.dart';

void main() {
  group('DeductionPlanStatus', () {
    // ✅ 正常流程 - 状态值映射
    test('should return correct status from value', () {
      expect(DeductionPlanStatus.fromValue(0), equals(DeductionPlanStatus.pendingAccept));
      expect(DeductionPlanStatus.fromValue(1), equals(DeductionPlanStatus.active));
      expect(DeductionPlanStatus.fromValue(2), equals(DeductionPlanStatus.completed));
      expect(DeductionPlanStatus.fromValue(3), equals(DeductionPlanStatus.paused));
      expect(DeductionPlanStatus.fromValue(4), equals(DeductionPlanStatus.cancelled));
      expect(DeductionPlanStatus.fromValue(5), equals(DeductionPlanStatus.rejected));
    });

    // ✅ 边界情况 - 无效值
    test('should return default status for invalid value', () {
      expect(DeductionPlanStatus.fromValue(99), equals(DeductionPlanStatus.active));
      expect(DeductionPlanStatus.fromValue(-1), equals(DeductionPlanStatus.active));
    });

    // ✅ 状态标签
    test('should return correct labels', () {
      expect(DeductionPlanStatus.pendingAccept.label, equals('待接收'));
      expect(DeductionPlanStatus.active.label, equals('进行中'));
      expect(DeductionPlanStatus.completed.label, equals('已完成'));
      expect(DeductionPlanStatus.paused.label, equals('已暂停'));
      expect(DeductionPlanStatus.cancelled.label, equals('已取消'));
      expect(DeductionPlanStatus.rejected.label, equals('已拒绝'));
    });

    // ✅ 操作权限判断
    group('action permissions', () {
      test('pendingAccept status can accept and reject', () {
        const status = DeductionPlanStatus.pendingAccept;
        expect(status.canAccept, isTrue);
        expect(status.canReject, isTrue);
        expect(status.canPause, isFalse);
        expect(status.canResume, isFalse);
        expect(status.canCancel, isFalse);
      });

      test('active status can pause and cancel', () {
        const status = DeductionPlanStatus.active;
        expect(status.canAccept, isFalse);
        expect(status.canReject, isFalse);
        expect(status.canPause, isTrue);
        expect(status.canResume, isFalse);
        expect(status.canCancel, isTrue);
      });

      test('paused status can resume and cancel', () {
        const status = DeductionPlanStatus.paused;
        expect(status.canAccept, isFalse);
        expect(status.canReject, isFalse);
        expect(status.canPause, isFalse);
        expect(status.canResume, isTrue);
        expect(status.canCancel, isTrue);
      });

      test('completed status cannot do any action', () {
        const status = DeductionPlanStatus.completed;
        expect(status.canAccept, isFalse);
        expect(status.canReject, isFalse);
        expect(status.canPause, isFalse);
        expect(status.canResume, isFalse);
        expect(status.canCancel, isFalse);
      });

      test('rejected status cannot do any action', () {
        const status = DeductionPlanStatus.rejected;
        expect(status.canAccept, isFalse);
        expect(status.canReject, isFalse);
        expect(status.canPause, isFalse);
        expect(status.canResume, isFalse);
        expect(status.canCancel, isFalse);
      });
    });
  });

  group('DeductionSource', () {
    test('should return correct source from value', () {
      expect(DeductionSource.fromValue(1), equals(DeductionSource.profit));
      expect(DeductionSource.fromValue(2), equals(DeductionSource.serviceFee));
      expect(DeductionSource.fromValue(3), equals(DeductionSource.both));
    });

    test('should return default source for invalid value', () {
      expect(DeductionSource.fromValue(99), equals(DeductionSource.both));
      expect(DeductionSource.fromValue(0), equals(DeductionSource.both));
    });

    test('should return correct labels', () {
      expect(DeductionSource.profit.label, equals('分润'));
      expect(DeductionSource.serviceFee.label, equals('服务费'));
      expect(DeductionSource.both.label, equals('分润+服务费'));
    });
  });

  group('DeductionPlanType', () {
    test('should return correct type from value', () {
      expect(DeductionPlanType.fromValue(1), equals(DeductionPlanType.goods));
      expect(DeductionPlanType.fromValue(2), equals(DeductionPlanType.partner));
      expect(DeductionPlanType.fromValue(3), equals(DeductionPlanType.deposit));
    });

    test('should return default type for invalid value', () {
      expect(DeductionPlanType.fromValue(99), equals(DeductionPlanType.partner));
    });

    test('should return correct labels', () {
      expect(DeductionPlanType.goods.label, equals('货款代扣'));
      expect(DeductionPlanType.partner.label, equals('伙伴代扣'));
      expect(DeductionPlanType.deposit.label, equals('押金代扣'));
    });
  });

  group('DeductionRecordStatus', () {
    test('should return correct status from value', () {
      expect(DeductionRecordStatus.fromValue(0), equals(DeductionRecordStatus.pending));
      expect(DeductionRecordStatus.fromValue(1), equals(DeductionRecordStatus.success));
      expect(DeductionRecordStatus.fromValue(2), equals(DeductionRecordStatus.partialSuccess));
      expect(DeductionRecordStatus.fromValue(3), equals(DeductionRecordStatus.failed));
    });

    test('should return correct labels', () {
      expect(DeductionRecordStatus.pending.label, equals('待扣款'));
      expect(DeductionRecordStatus.success.label, equals('成功'));
      expect(DeductionRecordStatus.partialSuccess.label, equals('部分成功'));
      expect(DeductionRecordStatus.failed.label, equals('失败'));
    });
  });

  group('DeductionPlan', () {
    late DeductionPlan plan;

    setUp(() {
      plan = DeductionPlan(
        id: 1,
        planNo: 'DK202401001',
        deductorId: 100,
        deductorName: '上级代理',
        deducteeId: 200,
        deducteeName: '下级代理',
        planType: 2,
        totalAmount: 100000, // 1000元
        deductedAmount: 30000, // 300元
        remainingAmount: 70000, // 700元
        frozenAmount: 50000, // 500元
        totalPeriods: 10,
        currentPeriod: 3,
        periodAmount: 10000, // 100元
        status: 1,
        needAccept: true,
        acceptedAt: '2024-01-15 10:00:00',
        deductionSource: 3,
        createdBy: 1,
        createdAt: '2024-01-01 10:00:00',
        updatedAt: '2024-01-15 10:00:00',
      );
    });

    // ✅ 金额转换测试
    group('amount conversions', () {
      test('should convert total amount to yuan', () {
        expect(plan.totalAmountYuan, equals(1000.0));
      });

      test('should convert deducted amount to yuan', () {
        expect(plan.deductedAmountYuan, equals(300.0));
      });

      test('should convert remaining amount to yuan', () {
        expect(plan.remainingAmountYuan, equals(700.0));
      });

      test('should convert period amount to yuan', () {
        expect(plan.periodAmountYuan, equals(100.0));
      });

      test('should convert frozen amount to yuan', () {
        expect(plan.frozenAmountYuan, equals(500.0));
      });
    });

    // ✅ 进度计算测试
    group('progress calculation', () {
      test('should calculate progress correctly', () {
        expect(plan.progress, equals(30.0)); // 30000/100000 * 100 = 30%
      });

      test('should return 0 for zero total amount', () {
        final zeroPlan = DeductionPlan(
          id: 1,
          planNo: 'DK202401001',
          deductorId: 100,
          deductorName: '上级代理',
          deducteeId: 200,
          deducteeName: '下级代理',
          planType: 2,
          totalAmount: 0,
          deductedAmount: 0,
          remainingAmount: 0,
          totalPeriods: 0,
          currentPeriod: 0,
          periodAmount: 0,
          status: 1,
          createdBy: 1,
          createdAt: '2024-01-01',
          updatedAt: '2024-01-01',
        );
        expect(zeroPlan.progress, equals(0.0));
      });

      test('should calculate 100% progress for completed plan', () {
        final completedPlan = DeductionPlan(
          id: 1,
          planNo: 'DK202401001',
          deductorId: 100,
          deductorName: '上级代理',
          deducteeId: 200,
          deducteeName: '下级代理',
          planType: 2,
          totalAmount: 100000,
          deductedAmount: 100000,
          remainingAmount: 0,
          totalPeriods: 10,
          currentPeriod: 10,
          periodAmount: 10000,
          status: 2,
          createdBy: 1,
          createdAt: '2024-01-01',
          updatedAt: '2024-01-01',
        );
        expect(completedPlan.progress, equals(100.0));
      });
    });

    // ✅ 枚举获取测试
    group('enum getters', () {
      test('should return correct status enum', () {
        expect(plan.statusEnum, equals(DeductionPlanStatus.active));
      });

      test('should return correct type enum', () {
        expect(plan.typeEnum, equals(DeductionPlanType.partner));
      });

      test('should return correct deduction source enum', () {
        expect(plan.deductionSourceEnum, equals(DeductionSource.both));
      });
    });

    // ✅ 状态判断测试
    group('status checks', () {
      test('should identify pending accept status', () {
        final pendingPlan = DeductionPlan(
          id: 1,
          planNo: 'DK202401001',
          deductorId: 100,
          deductorName: '上级代理',
          deducteeId: 200,
          deducteeName: '下级代理',
          planType: 2,
          totalAmount: 100000,
          deductedAmount: 0,
          remainingAmount: 100000,
          totalPeriods: 10,
          currentPeriod: 0,
          periodAmount: 10000,
          status: 0,
          needAccept: true,
          createdBy: 1,
          createdAt: '2024-01-01',
          updatedAt: '2024-01-01',
        );
        expect(pendingPlan.isPendingAccept, isTrue);
        expect(pendingPlan.isRejected, isFalse);
      });

      test('should identify rejected status', () {
        final rejectedPlan = DeductionPlan(
          id: 1,
          planNo: 'DK202401001',
          deductorId: 100,
          deductorName: '上级代理',
          deducteeId: 200,
          deducteeName: '下级代理',
          planType: 2,
          totalAmount: 100000,
          deductedAmount: 0,
          remainingAmount: 100000,
          totalPeriods: 10,
          currentPeriod: 0,
          periodAmount: 10000,
          status: 5,
          createdBy: 1,
          createdAt: '2024-01-01',
          updatedAt: '2024-01-01',
        );
        expect(rejectedPlan.isRejected, isTrue);
        expect(rejectedPlan.isPendingAccept, isFalse);
      });
    });

    // ✅ fromJson 测试
    group('fromJson', () {
      test('should parse json correctly', () {
        final json = {
          'id': 1,
          'plan_no': 'DK202401001',
          'deductor_id': 100,
          'deductor_name': '上级代理',
          'deductee_id': 200,
          'deductee_name': '下级代理',
          'plan_type': 2,
          'total_amount': 100000,
          'deducted_amount': 30000,
          'remaining_amount': 70000,
          'frozen_amount': 50000,
          'total_periods': 10,
          'current_period': 3,
          'period_amount': 10000,
          'status': 1,
          'need_accept': true,
          'accepted_at': '2024-01-15 10:00:00',
          'deduction_source': 3,
          'created_by': 1,
          'created_at': '2024-01-01 10:00:00',
          'updated_at': '2024-01-15 10:00:00',
        };

        final parsedPlan = DeductionPlan.fromJson(json);

        expect(parsedPlan.id, equals(1));
        expect(parsedPlan.planNo, equals('DK202401001'));
        expect(parsedPlan.totalAmount, equals(100000));
        expect(parsedPlan.frozenAmount, equals(50000));
        expect(parsedPlan.needAccept, isTrue);
        expect(parsedPlan.deductionSource, equals(3));
      });

      test('should handle missing fields with defaults', () {
        final json = <String, dynamic>{};

        final parsedPlan = DeductionPlan.fromJson(json);

        expect(parsedPlan.id, equals(0));
        expect(parsedPlan.planNo, equals(''));
        expect(parsedPlan.frozenAmount, equals(0));
        expect(parsedPlan.needAccept, isFalse);
        expect(parsedPlan.deductionSource, equals(3));
      });
    });
  });

  group('DeductionPlanListResponse', () {
    test('should parse list response correctly', () {
      final json = {
        'list': [
          {
            'id': 1,
            'plan_no': 'DK001',
            'deductor_id': 100,
            'deductor_name': '代理A',
            'deductee_id': 200,
            'deductee_name': '代理B',
            'plan_type': 2,
            'total_amount': 10000,
            'deducted_amount': 5000,
            'remaining_amount': 5000,
            'total_periods': 10,
            'current_period': 5,
            'period_amount': 1000,
            'status': 1,
            'created_by': 1,
            'created_at': '2024-01-01',
            'updated_at': '2024-01-01',
          },
        ],
        'total': 1,
        'page': 1,
        'page_size': 10,
      };

      final response = DeductionPlanListResponse.fromJson(json);

      expect(response.list.length, equals(1));
      expect(response.total, equals(1));
      expect(response.page, equals(1));
      expect(response.pageSize, equals(10));
      expect(response.list[0].planNo, equals('DK001'));
    });

    test('should handle empty list', () {
      final json = {
        'list': null,
        'total': 0,
        'page': 1,
        'page_size': 10,
      };

      final response = DeductionPlanListResponse.fromJson(json);

      expect(response.list, isEmpty);
      expect(response.total, equals(0));
    });
  });

  group('DeductionPlanStats', () {
    test('should parse stats correctly', () {
      final json = {
        'total_count': 100,
        'pending_accept_count': 10,
        'active_count': 50,
        'completed_count': 30,
        'paused_count': 5,
        'rejected_count': 5,
        'total_amount': 10000000,
        'deducted_amount': 5000000,
        'remaining_amount': 5000000,
        'frozen_amount': 3000000,
      };

      final stats = DeductionPlanStats.fromJson(json);

      expect(stats.totalCount, equals(100));
      expect(stats.pendingAcceptCount, equals(10));
      expect(stats.activeCount, equals(50));
      expect(stats.rejectedCount, equals(5));
      expect(stats.frozenAmount, equals(3000000));
      expect(stats.totalAmountYuan, equals(100000.0));
      expect(stats.frozenAmountYuan, equals(30000.0));
    });
  });

  group('DeductionSummary', () {
    test('should parse summary correctly', () {
      final json = {
        'received_pending_count': 5,
        'received_active_count': 10,
        'received_total_amount': 500000,
        'sent_pending_count': 3,
        'sent_active_count': 8,
        'sent_total_amount': 300000,
      };

      final summary = DeductionSummary.fromJson(json);

      expect(summary.receivedPendingCount, equals(5));
      expect(summary.receivedActiveCount, equals(10));
      expect(summary.receivedTotalAmount, equals(500000));
      expect(summary.receivedTotalAmountYuan, equals(5000.0));
      expect(summary.sentPendingCount, equals(3));
      expect(summary.sentActiveCount, equals(8));
      expect(summary.sentTotalAmount, equals(300000));
      expect(summary.sentTotalAmountYuan, equals(3000.0));
    });

    test('should handle missing fields with defaults', () {
      final json = <String, dynamic>{};

      final summary = DeductionSummary.fromJson(json);

      expect(summary.receivedPendingCount, equals(0));
      expect(summary.receivedActiveCount, equals(0));
      expect(summary.receivedTotalAmount, equals(0));
      expect(summary.sentPendingCount, equals(0));
    });
  });

  group('DeductionRecord', () {
    test('should parse record correctly', () {
      final json = {
        'id': 1,
        'plan_id': 10,
        'plan_no': 'DK001',
        'deductor_id': 100,
        'deductee_id': 200,
        'period_num': 3,
        'amount': 10000,
        'actual_amount': 10000,
        'status': 1,
        'scheduled_at': '2024-01-15',
        'deducted_at': '2024-01-15 08:00:00',
        'created_at': '2024-01-01',
      };

      final record = DeductionRecord.fromJson(json);

      expect(record.id, equals(1));
      expect(record.planId, equals(10));
      expect(record.periodNum, equals(3));
      expect(record.amount, equals(10000));
      expect(record.amountYuan, equals(100.0));
      expect(record.actualAmount, equals(10000));
      expect(record.actualAmountYuan, equals(100.0));
      expect(record.statusEnum, equals(DeductionRecordStatus.success));
    });
  });
}
