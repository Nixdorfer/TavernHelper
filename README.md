[English](./readme/README_EN.md) | [Deutsch](./readme/README_DE.md) | [繁體中文](./readme/README_TC.md)

# 关于本项目

使用时遇到了问题？[点我吐槽](https://github.com/Nixdorfer/TavernHelper/issues/new)

不知道怎么用？[点我看说明书](./readme/MANUAL.md)

## 初衷

我是从风月开始玩酒馆的，他们家感觉做的东西还挺不错的，但是编辑实在过于恶臭了，而我又恰好是一个题词狂魔（单人物10kt+），因此该程序的目的旨在于改善提词体验。因此当然默认对接了风月的API，要对接其他的平台可以在上面↑给我发个issue，酌情开咕。

## 下载运行

如果您是一般用户，可以从[发布页](https://github.com/Nixdorfer/TavernHelper/releases)下载适用于您系统的可执行文件。

## 编译/调试启动

执行 Run.cmd 可用参数如下：

- -d 调试模式：该模式下会编译带有调试器的可执行文件，并且在运行时启动开发者工具，还会在sidebar中显示调试页

- -r Rust/Tauri模式：该项目目前使用Wails作为程序框架，使用该flag会启动Rust版本

- -m 移动端启动：这将在Android Studio模拟器中启动编译后的应用程序

- -b 仅构建：这将不会编译带有调试器的dev版本，会清空bin目录后构建一个release版本的可执行程序

- -c 清洁启动：这将会先清空所有编译产物和缓存（包括数据库，但不包括creations文件夹）后再编译启动

## TODO

- Tauri v2重构 -> 支持Android/iOS

- UDP内网多播，支持基于WS Mutex的文件修改

- 支持通过本地或远程SD直接生成图片

- 支持LLM官方接口

  - Claude

  - ChatGPT

  - DeepSeek

  - Google AntiGravity

  - Google Gemini

  - Grok

- 基于SQLite和自动端口映射的客户端MCP支持

- 支持魅魔岛接口

- 支持使用SessionKey + SSE通过账号订阅调用接口

- 发布创作剧本

  - 自动定时发布

  - 多平台发布

- 服务端 (有生之年系列)

  - 基于Goroutine的并行处理

  - 基于pgsql的持久化数据库

  - 基于Redis的缓存数据库

  - 基于Kafka的消息队列

  - 基于Elastic Search + IK Analyse的关键词检索

  - 支持MCP Connector

  - 基于FRP的自动端口映射

  - 基于Docker的容器化管理

  - 基于K8s Ingress的CA/HPA弹性扩容

  - RESTful接口设计

## 声明

该程序**从未**并且**在将来也不会**以任何方式将您的任何信息传送至除了您指定的平台的**官方接口**以外的其他位置。

这代表：如果您的信息泄露，在技术层面与该程序**完全无关**，因此请不要在issue中发布相关问题，而是请咨询相关平台客服。

该声明仅对于本程序的[主要开源发行版本](https://github.com/Nixdorfer/TavernHelper)有效。

## 技术栈

框架：Tauri v2 / Wails

后端：Rust / Go

前端：Vue3 + TypeScript @ClaudeCode

数据库：SQLite

MCP服务端：TypeScript

## 结语

- 总之！欢迎使用！欢迎提出意见！

- 该项目基于MIT协议开源，这意味着您可以自由用于商用，但是如果能附上我的名字我会很感谢的👍🏿

## 目前接口支持平台

- [风月](https://aify.pages.dev/)

## 联系方式

QQ群 `1074532978`

[Telegram](https://t.me/Nixdorfer) 或复制 `Nixdorfer`

[电子邮箱](mailto:admin@nixdorfer.com) 或复制 `admin@nixdorfer.com`