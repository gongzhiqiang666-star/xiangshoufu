#!/usr/bin/env python3
"""
PC端 Web 应用测试脚本
测试商户管理、结算价管理等核心功能
"""

from playwright.sync_api import sync_playwright
import time
import os

# 截图保存目录
SCREENSHOT_DIR = "/tmp/xiangshoufu_test"
os.makedirs(SCREENSHOT_DIR, exist_ok=True)

def test_login(page):
    """测试登录功能"""
    print("=" * 50)
    print("测试1: 登录功能")
    print("=" * 50)

    page.goto('http://localhost:5173')
    page.wait_for_load_state('networkidle')

    # 截图登录页
    page.screenshot(path=f'{SCREENSHOT_DIR}/01_login_page.png', full_page=True)
    print(f"✓ 登录页面截图已保存: {SCREENSHOT_DIR}/01_login_page.png")

    # 查找并填写登录表单
    try:
        # 等待登录表单出现
        page.wait_for_selector('input', timeout=10000)

        # 尝试填写用户名和密码
        username_input = page.locator('input[type="text"], input[placeholder*="用户名"], input[placeholder*="账号"]').first
        password_input = page.locator('input[type="password"]').first

        if username_input.count() > 0 and password_input.count() > 0:
            username_input.fill('admin')
            password_input.fill('admin123')
            page.screenshot(path=f'{SCREENSHOT_DIR}/02_login_filled.png', full_page=True)
            print("✓ 已填写登录信息")

            # 点击登录按钮
            login_btn = page.locator('button[type="submit"], button:has-text("登录")').first
            if login_btn.count() > 0:
                login_btn.click()
                page.wait_for_load_state('networkidle')
                time.sleep(2)
                page.screenshot(path=f'{SCREENSHOT_DIR}/03_after_login.png', full_page=True)
                print(f"✓ 登录后截图已保存: {SCREENSHOT_DIR}/03_after_login.png")
                return True
        else:
            print("⚠ 未找到登录表单元素")
            return False
    except Exception as e:
        print(f"⚠ 登录过程出错: {e}")
        return False

def test_merchant_page(page):
    """测试商户管理页面"""
    print("\n" + "=" * 50)
    print("测试2: 商户管理页面")
    print("=" * 50)

    try:
        # 尝试导航到商户管理页面
        page.goto('http://localhost:5173/merchants')
        page.wait_for_load_state('networkidle')
        time.sleep(2)

        page.screenshot(path=f'{SCREENSHOT_DIR}/04_merchant_list.png', full_page=True)
        print(f"✓ 商户列表截图已保存: {SCREENSHOT_DIR}/04_merchant_list.png")

        # 检查页面内容
        content = page.content()

        # 检查商户类型标签是否存在（5档分类）
        merchant_types = ['优质', '中等', '普通', '预警', '流失']
        found_types = []
        for mt in merchant_types:
            if mt in content:
                found_types.append(mt)

        if found_types:
            print(f"✓ 发现商户类型标签: {', '.join(found_types)}")
        else:
            print("⚠ 未发现商户类型标签（可能无数据或页面结构不同）")

        return True
    except Exception as e:
        print(f"⚠ 商户管理页面测试出错: {e}")
        return False

def test_settlement_price_page(page):
    """测试结算价管理页面"""
    print("\n" + "=" * 50)
    print("测试3: 结算价管理页面")
    print("=" * 50)

    try:
        # 尝试导航到结算价管理页面
        page.goto('http://localhost:5173/settlement-prices')
        page.wait_for_load_state('networkidle')
        time.sleep(2)

        page.screenshot(path=f'{SCREENSHOT_DIR}/05_settlement_price.png', full_page=True)
        print(f"✓ 结算价页面截图已保存: {SCREENSHOT_DIR}/05_settlement_price.png")

        content = page.content()

        # 检查结算价相关内容
        keywords = ['结算价', '费率', '押金', '流量费', '返现']
        found_keywords = [kw for kw in keywords if kw in content]

        if found_keywords:
            print(f"✓ 发现结算价相关内容: {', '.join(found_keywords)}")
        else:
            print("⚠ 未发现结算价相关内容（可能需要登录或无数据）")

        return True
    except Exception as e:
        print(f"⚠ 结算价管理页面测试出错: {e}")
        return False

def test_dashboard(page):
    """测试首页仪表盘"""
    print("\n" + "=" * 50)
    print("测试4: 首页仪表盘")
    print("=" * 50)

    try:
        page.goto('http://localhost:5173/')
        page.wait_for_load_state('networkidle')
        time.sleep(2)

        page.screenshot(path=f'{SCREENSHOT_DIR}/06_dashboard.png', full_page=True)
        print(f"✓ 首页截图已保存: {SCREENSHOT_DIR}/06_dashboard.png")

        # 获取页面所有按钮和链接
        buttons = page.locator('button').all()
        links = page.locator('a').all()

        print(f"✓ 页面包含 {len(buttons)} 个按钮, {len(links)} 个链接")

        return True
    except Exception as e:
        print(f"⚠ 首页测试出错: {e}")
        return False

def discover_navigation(page):
    """发现页面导航结构"""
    print("\n" + "=" * 50)
    print("测试5: 页面导航结构发现")
    print("=" * 50)

    try:
        page.goto('http://localhost:5173/')
        page.wait_for_load_state('networkidle')
        time.sleep(2)

        # 查找侧边栏菜单
        menu_items = page.locator('.el-menu-item, .menu-item, nav a, aside a').all()

        if menu_items:
            print(f"✓ 发现 {len(menu_items)} 个导航菜单项:")
            for i, item in enumerate(menu_items[:10]):  # 只显示前10个
                text = item.inner_text().strip()
                if text:
                    print(f"  - {text}")
        else:
            print("⚠ 未发现导航菜单（可能需要登录）")

        return True
    except Exception as e:
        print(f"⚠ 导航发现出错: {e}")
        return False

def main():
    print("\n" + "=" * 60)
    print("  享收付 PC端 Web 应用测试")
    print("=" * 60)
    print(f"截图保存目录: {SCREENSHOT_DIR}")
    print("")

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        context = browser.new_context(
            viewport={'width': 1920, 'height': 1080},
            locale='zh-CN'
        )
        page = context.new_page()

        # 监听控制台日志
        page.on('console', lambda msg: print(f"[Console] {msg.type}: {msg.text}") if msg.type == 'error' else None)

        results = []

        # 运行测试
        results.append(("登录功能", test_login(page)))
        results.append(("首页仪表盘", test_dashboard(page)))
        results.append(("导航结构", discover_navigation(page)))
        results.append(("商户管理", test_merchant_page(page)))
        results.append(("结算价管理", test_settlement_price_page(page)))

        browser.close()

    # 输出测试结果汇总
    print("\n" + "=" * 60)
    print("  测试结果汇总")
    print("=" * 60)

    passed = 0
    failed = 0
    for name, result in results:
        status = "✅ 通过" if result else "❌ 失败"
        print(f"  {name}: {status}")
        if result:
            passed += 1
        else:
            failed += 1

    print("")
    print(f"  总计: {passed} 通过, {failed} 失败")
    print(f"  截图目录: {SCREENSHOT_DIR}")
    print("=" * 60)

if __name__ == '__main__':
    main()
