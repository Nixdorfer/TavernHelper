# 关于本项目

[使用时遇到了问题? 点我吐槽](https://github.com/Nixdorfer/TavernHelper/issues/new)

## 初衷

我是从风月开始玩酒馆的 他们家感觉做的东西还挺不错的 但是编辑实在过于恶臭了 而我又恰好是一个题词狂魔(单人物10kt+) 因此该程序的目的旨在于改善提词体验 因此当然默认对接了风月的API 要对接其他的平台可以在上面↑给我发个issue 酌情开咕

## *下载*/编译/启动

如果您是一般用户 可以从[发布页](https://github.com/Nixdorfer/TavernHelper/releases)下载适用于您系统的可执行文件

执行 运行或编译.cmd 默认启动无控制台模式 可附加参数列表如下

-platform=pc/mobile 启动平台是PC还是移动端 默认: PC

注意: 如果要在移动端启动 需要安装Android SDK

-ver=rust/go 启动/编译版本是rust(tauri v2)还是go(wails) 默认: go(wails)

-mode=run/build 启动模式为仅编译还是直接运行 默认: run

-d 启用debug模式 默认: 关闭

交叉编译功能请自行实现

## 声明

该程序**从未**并且**在将来也不会**以任何方式将您的任何信息传送至除了您指定的平台的**官方接口**以外的其他位置

这代表: 如果您的信息泄露 在技术层面与该程序**完全无关** 因此请不要在issue中发布相关问题 而是请咨询相关平台客服

如果您不信任[发布页](https://github.com/Nixdorfer/TavernHelper/releases)中的可执行文件 可以审阅代码后使用自动脚本自行编译

该声明仅对于本程序的[主要开源发行版本](https://github.com/Nixdorfer/TavernHelper)有效

### TODO

- iOS/Android支持

  - 该项目使用Wails+Vue3实现 因此只可以在桌面端部署

  - 移动端暂时没有好主意 不如咱们内置虚幻5罢！/笑

- 压缩的非二进制文本

  - 剧本将不再显示可理解文字 而是以看似随机的字母数字和符号组成 以规避各种平台的审查

  - 剧本将以非二进制文件的形式导出 这可以使您在任何支持复制粘贴文本的地方传播您的剧本 而不需要上传文件或外链

- 可单独导出的最小单元

  - 您可以任意在层面导入导出 不受任何限制 这包括

    - 整个项目(也即整个世界树)

    - 某条从根节点开始的特定世界线

    - 某一个或几个事件

    - 前置/后置/全局提示词

    - 全部或一部分世界书

## 技术栈

Go作为后端

显示界面基于Wails

界面使用Vue3 @ClaudeCode

SQLite作为数据储存

非对称加密基于ssh-ed25519

## 结语

- 总之！欢迎使用！欢迎提出意见！

- 该项目基于MIT协议开源 这意味着您可以自由用于商用 但是如果能附上我的名字我会很感谢的👍🏿

## 目前接口支持平台

- [风月AI](https://aify.pages.dev/)

## 联系方式

QQ群`1074532978`

[Telegram](https://t.me/Nixdorfer)或复制`Nixdorfer`

[给我发邮件](mailto:admin@nixdorfer.com)或复制`admin@nixdorfer.com`
