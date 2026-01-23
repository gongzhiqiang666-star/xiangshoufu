#!/usr/bin/env python3
"""检查 Dashboard 页面的错误"""
from playwright.sync_api import sync_playwright

def main():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        # 收集控制台消息
        console_messages = []
        page.on("console", lambda msg: console_messages.append({
            "type": msg.type,
            "text": msg.text,
            "location": msg.location
        }))

        # 收集页面错误
        page_errors = []
        page.on("pageerror", lambda err: page_errors.append(str(err)))

        # 收集请求失败
        failed_requests = []
        page.on("requestfailed", lambda req: failed_requests.append({
            "url": req.url,
            "failure": req.failure
        }))

        try:
            print("正在访问 http://localhost:5173/ ...")
            page.goto("http://localhost:5173/", timeout=30000)
            page.wait_for_load_state("networkidle", timeout=15000)

            # 等待一下让所有请求完成
            page.wait_for_timeout(3000)

            # 截图
            page.screenshot(path="/tmp/dashboard_screenshot.png", full_page=True)
            print(f"截图已保存到 /tmp/dashboard_screenshot.png")

        except Exception as e:
            print(f"页面访问错误: {e}")
            page.screenshot(path="/tmp/dashboard_error.png", full_page=True)

        # 输出控制台错误
        print("\n" + "="*60)
        print("控制台错误 (error/warning):")
        print("="*60)
        for msg in console_messages:
            if msg["type"] in ["error", "warning"]:
                print(f"[{msg['type'].upper()}] {msg['text']}")
                if msg.get("location"):
                    print(f"  位置: {msg['location']}")

        # 输出页面错误
        if page_errors:
            print("\n" + "="*60)
            print("页面错误 (JavaScript 异常):")
            print("="*60)
            for err in page_errors:
                print(f"  {err}")

        # 输出请求失败
        if failed_requests:
            print("\n" + "="*60)
            print("请求失败:")
            print("="*60)
            for req in failed_requests:
                print(f"  URL: {req['url']}")
                print(f"  原因: {req['failure']}")

        # 输出所有控制台日志
        print("\n" + "="*60)
        print("所有控制台日志:")
        print("="*60)
        for msg in console_messages:
            print(f"[{msg['type']}] {msg['text']}")

        browser.close()

if __name__ == "__main__":
    main()
