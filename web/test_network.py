#!/usr/bin/env python3
"""检查页面的 404 网络请求"""
from playwright.sync_api import sync_playwright

def main():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        # 收集所有请求和响应
        all_requests = []
        error_requests = []

        def handle_response(response):
            request = response.request
            status = response.status
            all_requests.append({
                "url": request.url,
                "method": request.method,
                "status": status,
                "status_text": response.status_text
            })
            # 记录 4xx 和 5xx 错误
            if status >= 400:
                error_requests.append({
                    "url": request.url,
                    "method": request.method,
                    "status": status,
                    "status_text": response.status_text
                })

        page.on("response", handle_response)

        try:
            print("正在访问 http://localhost:5173/dashboard ...")
            page.goto("http://localhost:5173/dashboard", timeout=30000)
            page.wait_for_load_state("networkidle", timeout=15000)

            # 等待让所有请求完成
            page.wait_for_timeout(5000)

            print(f"\n总请求数: {len(all_requests)}")

        except Exception as e:
            print(f"页面访问错误: {e}")

        # 输出 404 请求
        print("\n" + "="*70)
        print("404 请求:")
        print("="*70)
        not_found = [r for r in error_requests if r["status"] == 404]
        if not_found:
            for req in not_found:
                print(f"  [{req['method']}] {req['url']}")
        else:
            print("  无 404 请求")

        # 输出其他错误请求 (4xx, 5xx)
        other_errors = [r for r in error_requests if r["status"] != 404]
        if other_errors:
            print("\n" + "="*70)
            print("其他错误请求 (4xx/5xx):")
            print("="*70)
            for req in other_errors:
                print(f"  [{req['status']}] [{req['method']}] {req['url']}")

        # 输出所有 API 请求
        print("\n" + "="*70)
        print("所有 API 请求 (/v1/ 或 /api/):")
        print("="*70)
        api_requests = [r for r in all_requests if "/v1/" in r["url"] or "/api/" in r["url"]]
        if api_requests:
            for req in api_requests:
                status_icon = "✅" if req["status"] < 400 else "❌"
                print(f"  {status_icon} [{req['status']}] [{req['method']}] {req['url']}")
        else:
            print("  无 API 请求")

        browser.close()

if __name__ == "__main__":
    main()
