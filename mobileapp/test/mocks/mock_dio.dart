/// Dio Mock 工具
///
/// 用于模拟网络请求的 Mock 类和工具函数

import 'package:dio/dio.dart';
import 'package:mocktail/mocktail.dart';

/// Mock Dio 实例
class MockDio extends Mock implements Dio {}

/// Mock Response
class MockResponse<T> extends Mock implements Response<T> {}

/// 创建成功的 Response
Response<T> createSuccessResponse<T>(T data, {int statusCode = 200}) {
  return Response<T>(
    data: data,
    statusCode: statusCode,
    requestOptions: RequestOptions(path: ''),
  );
}

/// 创建失败的 Response
Response<T> createErrorResponse<T>({
  int statusCode = 500,
  String message = 'Server Error',
}) {
  return Response<T>(
    data: {'message': message} as T,
    statusCode: statusCode,
    requestOptions: RequestOptions(path: ''),
  );
}

/// 创建 DioException
DioException createDioException({
  DioExceptionType type = DioExceptionType.badResponse,
  int? statusCode,
  String? message,
}) {
  return DioException(
    type: type,
    message: message ?? 'Request failed',
    requestOptions: RequestOptions(path: ''),
    response: statusCode != null
        ? Response(
            statusCode: statusCode,
            data: {'message': message ?? 'Error'},
            requestOptions: RequestOptions(path: ''),
          )
        : null,
  );
}

/// 模拟 API 成功响应
///
/// 用于设置 MockDio 的 get/post 等方法返回成功响应
void mockGetSuccess(MockDio dio, String path, dynamic data) {
  when(() => dio.get(
        path,
        queryParameters: any(named: 'queryParameters'),
        options: any(named: 'options'),
      )).thenAnswer((_) async => createSuccessResponse(data));
}

/// 模拟 API POST 成功响应
void mockPostSuccess(MockDio dio, String path, dynamic data) {
  when(() => dio.post(
        path,
        data: any(named: 'data'),
        queryParameters: any(named: 'queryParameters'),
        options: any(named: 'options'),
      )).thenAnswer((_) async => createSuccessResponse(data));
}

/// 模拟 API 失败响应
void mockGetError(MockDio dio, String path, {int statusCode = 500, String? message}) {
  when(() => dio.get(
        path,
        queryParameters: any(named: 'queryParameters'),
        options: any(named: 'options'),
      )).thenThrow(createDioException(
    statusCode: statusCode,
    message: message,
  ));
}

/// 模拟网络错误
void mockNetworkError(MockDio dio, String path) {
  when(() => dio.get(
        path,
        queryParameters: any(named: 'queryParameters'),
        options: any(named: 'options'),
      )).thenThrow(createDioException(
    type: DioExceptionType.connectionError,
    message: 'Network error',
  ));
}
