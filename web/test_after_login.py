#!/usr/bin/env python3
"""登录后检查 Dashboard 页面的错误和 404 请求"""
from playwright.sync_api import sync_playwright

def main():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        # 收集控制台消息
        console_messages = []
        page.on("console", lambda msg: console_messages.append({
            "type": msg.type,
            "text": msg.text
        }))

        # 收集页面错误
        page_errors = []
        page.on("pageerror", lambda err: page_errors.append(str(err)))

        # 收集所有请求响应
        all_requests = []
        error_requests = []

        def handle_response(response):
            request = response.request
            status = response.status
            all_requests.append({
                "url": request.url,
                "method": request.method,
                "status": status
            })
            if status >= 400:
                error_requests.append({
                    "url": request.url,
                    "method": request.method,
                    "status": status
                })

        page.on("response", handle_response)

        try:
            # 1. 访问登录页面
            print("=" * 70)
            print("步骤 1: 访问登录页面...")
            print("=" * 70)
            page.goto("http://localhost:5173/login", timeout=30000)
            page.wait_for_load_state("networkidle", timeout=15000)
            page.screenshot(path="/tmp/step1_login_page.png")
            print("  ✅ 登录页面已加载")

            # 2. 输入用户名和密码
            print("\n步骤 2: 输入登录凭据...")
            # 等待输入框加载
            page.wait_for_selector('.el-input__inner', timeout=10000)

            # 定位用户名和密码输入框
            inputs = page.locator('.el-input__inner').all()
            if len(inputs) >= 2:
                inputs[0].fill("admin")
                inputs[1].fill("123456")
                print("  ✅ 已输入用户名: admin")
                print("  ✅ 已输入密码: ******")
            else:
                print(f"  ❌ 只找到 {len(inputs)} 个输入框")

            # 3. 点击登录按钮
            print("\n步骤 3: 点击登录...")
            # 使用 class 定位登录按钮
            login_button = page.locator('.login-btn')
            login_button.click()

            # 等待登录完成和跳转
            page.wait_for_timeout(3000)
            page.wait_for_load_state("networkidle", timeout=15000)

            current_url = page.url
            print(f"  当前 URL: {current_url}")

            if "/login" not in current_url:
                print("  ✅ 登录成功，已跳转")
            else:
                print("  ⚠️ 可能登录失败，仍在登录页")
                page.screenshot(path="/tmp/login_failed.png")

            # 4. 访问 Dashboard
            print("\n步骤 4: 访问 Dashboard...")
            page.goto("http://localhost:5173/dashboard", timeout=30000)
            page.wait_for_load_state("networkidle", timeout=15000)

            # 等待数据加载
            page.wait_for_timeout(5000)

            page.screenshot(path="/tmp/step4_dashboard.png", full_page=True)
            print(f"  当前 URL: {page.url}")
            print("  ✅ Dashboard 页面已加载")

        except Exception as e:
            print(f"\n❌ 执行错误: {e}")
            page.screenshot(path="/tmp/error_screenshot.png", full_page=True)

        # 输出结果
        print("\n" + "=" * 70)
        print("检测结果")
        print("=" * 70)

        # 404 请求
        print("\n【404 请求】")
        not_found = [r for r in error_requests if r["status"] == 404]
        if not_found:
            for req in not_found:
                print(f"  ❌ [{req['method']}] {req['url']}")
        else:
            print("  ✅ 无 404 请求")

        # 其他错误请求
        other_errors = [r for r in error_requests if r["status"] != 404]
        if other_errors:
            print("\n【其他 HTTP 错误 (4xx/5xx)】")
            for req in other_errors:
                print(f"  ❌ [{req['status']}] [{req['method']}] {req['url']}")

        # 控制台错误
        print("\n【控制台错误】")
        errors = [m for m in console_messages if m["type"] == "error"]
        if errors:
            for msg in errors:
                text = msg["text"][:200] + "..." if len(msg["text"]) > 200 else msg["text"]
                print(f"  ❌ {text}")
        else:
            print("  ✅ 无控制台错误")

        # 控制台警告
        print("\n【控制台警告】")
        warnings = [m for m in console_messages if m["type"] == "warning"]
        if warnings:
            for msg in warnings:
                text = msg["text"][:200] + "..." if len(msg["text"]) > 200 else msg["text"]
                print(f"  ⚠️ {text}")
        else:
            print("  ✅ 无控制台警告")

        # 页面错误
        if page_errors:
            print("\n【JavaScript 异常】")
            for err in page_errors:
                print(f"  ❌ {err[:200]}...")

        # API 请求统计
        print("\n【API 请求统计 (/v1/)】")
        api_requests = [r for r in all_requests if "/v1/" in r["url"]]
        success_count = len([r for r in api_requests if r["status"] < 400])
        fail_count = len([r for r in api_requests if r["status"] >= 400])
        print(f"  成功: {success_count}  失败: {fail_count}")

        if api_requests:
            print("\n  详细列表:")
            for req in api_requests:
                icon = "✅" if req["status"] < 400 else "❌"
                url_short = req["url"].replace("http://localhost:5173", "").replace("http://localhost:8080", "")
                print(f"    {icon} [{req['status']}] {url_short}")

        print("\n" + "=" * 70)
        print("截图已保存:")
        print("  /tmp/step1_login_page.png")
        print("  /tmp/step4_dashboard.png")
        print("=" * 70)

        browser.close()

if __name__ == "__main__":
    main()
