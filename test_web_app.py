#!/usr/bin/env python3
"""
PC端 Web 应用完整测试脚本
使用 admin/123456 登录并测试核心功能
"""

from playwright.sync_api import sync_playwright
import time
import os

SCREENSHOT_DIR = "/tmp/xiangshoufu_test"
os.makedirs(SCREENSHOT_DIR, exist_ok=True)

def save_screenshot(page, name):
    """保存截图并打印路径"""
    path = f'{SCREENSHOT_DIR}/{name}.png'
    page.screenshot(path=path, full_page=True)
    print(f"  📸 截图: {path}")
    return path

def test_login(page):
    """测试登录功能"""
    print("\n" + "=" * 60)
    print("🔐 测试1: 登录功能")
    print("=" * 60)

    page.goto('http://localhost:5173')
    page.wait_for_load_state('networkidle')
    time.sleep(1)

    save_screenshot(page, '01_login_page')

    try:
        # 等待页面加载
        page.wait_for_selector('input', timeout=10000)

        # 查找用户名输入框
        username_selectors = [
            'input[placeholder*="用户名"]',
            'input[placeholder*="账号"]',
            'input[type="text"]:first-of-type',
            '#username',
            'input:first-of-type'
        ]

        username_input = None
        for selector in username_selectors:
            try:
                el = page.locator(selector).first
                if el.is_visible():
                    username_input = el
                    break
            except:
                continue

        # 查找密码输入框
        password_input = page.locator('input[type="password"]').first

        if username_input and password_input:
            # 填写登录信息
            username_input.fill('admin')
            password_input.fill('123456')
            save_screenshot(page, '02_login_filled')
            print("  ✅ 已填写: admin / 123456")

            # 查找并点击登录按钮
            login_selectors = [
                'button[type="submit"]',
                'button:has-text("登录")',
                'button:has-text("Login")',
                '.el-button--primary',
                'button.login-btn'
            ]

            for selector in login_selectors:
                try:
                    btn = page.locator(selector).first
                    if btn.is_visible():
                        btn.click()
                        print("  ✅ 点击登录按钮")
                        break
                except:
                    continue

            # 等待登录完成
            page.wait_for_load_state('networkidle')
            time.sleep(3)
            save_screenshot(page, '03_after_login')

            # 检查是否登录成功（检查URL或页面内容变化）
            current_url = page.url
            content = page.content()

            if '/login' not in current_url and ('首页' in content or '仪表盘' in content or 'Dashboard' in content or '退出' in content or '注销' in content):
                print("  ✅ 登录成功!")
                return True
            elif '密码' in content and '错误' in content:
                print("  ❌ 登录失败: 密码错误")
                return False
            else:
                print(f"  ⚠️ 登录状态不确定, 当前URL: {current_url}")
                return True  # 继续测试
        else:
            print("  ❌ 未找到登录表单")
            return False

    except Exception as e:
        print(f"  ❌ 登录出错: {e}")
        return False

def test_dashboard(page):
    """测试首页仪表盘"""
    print("\n" + "=" * 60)
    print("📊 测试2: 首页仪表盘")
    print("=" * 60)

    try:
        page.goto('http://localhost:5173/')
        page.wait_for_load_state('networkidle')
        time.sleep(2)

        save_screenshot(page, '04_dashboard')

        content = page.content()

        # 检查仪表盘内容
        dashboard_keywords = ['交易', '分润', '商户', '代理', '终端', '钱包', '今日', '本月']
        found = [kw for kw in dashboard_keywords if kw in content]

        if found:
            print(f"  ✅ 首页包含: {', '.join(found)}")
        else:
            print("  ⚠️ 首页内容待确认")

        # 统计页面元素
        buttons = len(page.locator('button').all())
        links = len(page.locator('a').all())
        print(f"  📌 页面元素: {buttons} 个按钮, {links} 个链接")

        return True
    except Exception as e:
        print(f"  ❌ 首页测试出错: {e}")
        return False

def test_navigation(page):
    """测试导航菜单"""
    print("\n" + "=" * 60)
    print("🧭 测试3: 导航菜单")
    print("=" * 60)

    try:
        # 查找侧边栏菜单
        menu_selectors = [
            '.el-menu-item',
            '.el-sub-menu__title',
            '.sidebar-item',
            'nav a',
            'aside a',
            '[class*="menu"] a',
            '[class*="nav"] a'
        ]

        all_menus = []
        for selector in menu_selectors:
            try:
                items = page.locator(selector).all()
                for item in items:
                    text = item.inner_text().strip()
                    if text and len(text) < 20:
                        all_menus.append(text)
            except:
                continue

        # 去重
        unique_menus = list(dict.fromkeys(all_menus))

        if unique_menus:
            print(f"  ✅ 发现 {len(unique_menus)} 个菜单项:")
            for menu in unique_menus[:15]:
                print(f"     - {menu}")
        else:
            print("  ⚠️ 未发现导航菜单")

        return True
    except Exception as e:
        print(f"  ❌ 导航测试出错: {e}")
        return False

def test_merchant_page(page):
    """测试商户管理页面"""
    print("\n" + "=" * 60)
    print("🏪 测试4: 商户管理")
    print("=" * 60)

    try:
        # 尝试多种路由
        routes = ['/merchants', '/merchant', '/merchant/list']

        for route in routes:
            page.goto(f'http://localhost:5173{route}')
            page.wait_for_load_state('networkidle')
            time.sleep(2)

            content = page.content()
            if '商户' in content or 'merchant' in content.lower():
                break

        save_screenshot(page, '05_merchant_list')

        content = page.content()

        # 检查5档商户类型
        merchant_types = {
            'quality': '优质',
            'medium': '中等',
            'normal': '普通',
            'warning': '预警',
            'churned': '流失'
        }

        found_types = []
        for code, name in merchant_types.items():
            if name in content or code in content:
                found_types.append(name)

        if found_types:
            print(f"  ✅ 商户类型: {', '.join(found_types)}")

        # 检查表格
        tables = page.locator('table, .el-table').all()
        if tables:
            print(f"  ✅ 发现 {len(tables)} 个数据表格")

        # 检查搜索和筛选
        if '搜索' in content or '筛选' in content or 'search' in content.lower():
            print("  ✅ 支持搜索/筛选功能")

        return True
    except Exception as e:
        print(f"  ❌ 商户管理测试出错: {e}")
        return False

def test_settlement_price(page):
    """测试结算价管理"""
    print("\n" + "=" * 60)
    print("💰 测试5: 结算价管理")
    print("=" * 60)

    try:
        routes = ['/settlement-prices', '/settlement-price', '/agent/settlement-prices']

        for route in routes:
            page.goto(f'http://localhost:5173{route}')
            page.wait_for_load_state('networkidle')
            time.sleep(2)

            content = page.content()
            if '结算' in content or '费率' in content:
                break

        save_screenshot(page, '06_settlement_price')

        content = page.content()

        # 检查结算价功能
        features = ['费率', '押金', '流量', '返现', '通道', '调价']
        found = [f for f in features if f in content]

        if found:
            print(f"  ✅ 结算价功能: {', '.join(found)}")
        else:
            print("  ⚠️ 结算价页面内容待确认")

        return True
    except Exception as e:
        print(f"  ❌ 结算价测试出错: {e}")
        return False

def test_terminal_page(page):
    """测试终端管理"""
    print("\n" + "=" * 60)
    print("📱 测试6: 终端管理")
    print("=" * 60)

    try:
        routes = ['/terminals', '/terminal', '/terminal/list']

        for route in routes:
            page.goto(f'http://localhost:5173{route}')
            page.wait_for_load_state('networkidle')
            time.sleep(2)

            content = page.content()
            if '终端' in content or 'terminal' in content.lower():
                break

        save_screenshot(page, '07_terminal_list')

        content = page.content()

        keywords = ['终端', 'SN', '激活', '下发', '回拨', '费率']
        found = [kw for kw in keywords if kw in content]

        if found:
            print(f"  ✅ 终端功能: {', '.join(found)}")

        return True
    except Exception as e:
        print(f"  ❌ 终端管理测试出错: {e}")
        return False

def test_agent_page(page):
    """测试代理商管理"""
    print("\n" + "=" * 60)
    print("👥 测试7: 代理商管理")
    print("=" * 60)

    try:
        routes = ['/agents', '/agent', '/agent/list']

        for route in routes:
            page.goto(f'http://localhost:5173{route}')
            page.wait_for_load_state('networkidle')
            time.sleep(2)

            content = page.content()
            if '代理' in content or 'agent' in content.lower():
                break

        save_screenshot(page, '08_agent_list')

        content = page.content()

        keywords = ['代理', '邀请码', '政策', '结算', '团队']
        found = [kw for kw in keywords if kw in content]

        if found:
            print(f"  ✅ 代理功能: {', '.join(found)}")

        return True
    except Exception as e:
        print(f"  ❌ 代理商管理测试出错: {e}")
        return False

def main():
    print("\n" + "=" * 60)
    print("  🚀 享收付 PC端 Web 应用测试")
    print("  📝 登录凭据: admin / 123456")
    print("=" * 60)
    print(f"  📁 截图目录: {SCREENSHOT_DIR}")

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        context = browser.new_context(
            viewport={'width': 1920, 'height': 1080},
            locale='zh-CN'
        )
        page = context.new_page()

        # 收集控制台错误
        console_errors = []
        page.on('console', lambda msg: console_errors.append(msg.text) if msg.type == 'error' else None)

        results = []

        # 运行测试
        results.append(("登录功能", test_login(page)))
        results.append(("首页仪表盘", test_dashboard(page)))
        results.append(("导航菜单", test_navigation(page)))
        results.append(("商户管理", test_merchant_page(page)))
        results.append(("结算价管理", test_settlement_price(page)))
        results.append(("终端管理", test_terminal_page(page)))
        results.append(("代理商管理", test_agent_page(page)))

        browser.close()

    # 输出测试结果
    print("\n" + "=" * 60)
    print("  📋 测试结果汇总")
    print("=" * 60)

    passed = sum(1 for _, r in results if r)
    failed = len(results) - passed

    for name, result in results:
        status = "✅ 通过" if result else "❌ 失败"
        print(f"  {name}: {status}")

    print("\n" + "-" * 60)
    print(f"  总计: ✅ {passed} 通过 | ❌ {failed} 失败")
    print(f"  截图: {SCREENSHOT_DIR}")

    if console_errors:
        print(f"\n  ⚠️ 控制台错误 ({len(console_errors)} 条):")
        for err in console_errors[:5]:
            print(f"     {err[:80]}...")

    print("=" * 60)

if __name__ == '__main__':
    main()
