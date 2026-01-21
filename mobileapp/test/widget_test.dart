// Basic Flutter widget test for XiangshoufuApp

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

void main() {
  testWidgets('App smoke test', (WidgetTester tester) async {
    // Build a simple widget to verify the test framework works
    await tester.pumpWidget(
      const ProviderScope(
        child: MaterialApp(
          home: Scaffold(
            body: Center(
              child: Text('享收付'),
            ),
          ),
        ),
      ),
    );

    // Verify the text is displayed
    expect(find.text('享收付'), findsOneWidget);
  });
}
