import uiautomator2 as u2

# 自动连接（USB）
d = u2.connect()

# 启动大麦 App（包名可用 adb shell dumpsys package 找到）
d.app_start("cn.damai")



# 查找“演出”关键词并点击
d(textContains="演出").click_exists(timeout=5)

# 模拟滑动找按钮
d.swipe(500, 1500, 500, 500)

# 点击“立即预定”
if d(textContains="立即预定").exists:
    d(textContains="立即预定").click()
else:
    print("没找到")
