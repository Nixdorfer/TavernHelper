[简体中文](../README.md) | [Deutsch](./README_DE.md) | [繁體中文](./README_TC.md)

# About This Project

Having issues? [Click here to report](https://github.com/Nixdorfer/TavernHelper/issues/new)

Not sure how to use it? [Click here for the manual](./MANUAL_EN.md)

## Motivation

I started playing Tavern through 风月. Their products are pretty good, but the editing experience is quite frustrating. As someone who loves prompting extensively (10k+ tokens per character), this program aims to improve the prompting experience. Therefore, it naturally integrates with 风月's API by default. If you want to integrate with other platforms, you can submit an issue above ↑ and I'll consider it.

## Download & Run

If you're a regular user, you can download the executable for your system from the [Releases page](https://github.com/Nixdorfer/TavernHelper/releases).

## Compile / Debug Launch

Execute Run.cmd with the following available parameters:

- -d Debug mode: Compiles an executable with debugger, launches developer tools at runtime, and displays a debug page in the sidebar

- -m Mobile launch: This will launch the compiled application in Android Studio emulator

- -b Build only: This will not compile a dev version with debugger. It clears the bin directory and builds a release version executable

- -c Clean launch: This will clear all compiled outputs and cache (including database, but not the creations folder) before compiling and launching

## TODO

- Tauri v2 refactoring -> Android/iOS support

- UDP LAN multicast with Mutex file modification support

- Support direct image generation through local or network SD-WebUI

- MCP support

  - New server supporting MCP calls

  - Built-in automatic port mapping

- Meimodo API support

- Support using SessionKey + SSE for account subscription API calls

## Disclaimer

This program **HAS NEVER** and **WILL NEVER** transmit any of your information to any location other than the **OFFICIAL API** of the platform you specify.

This means: If your information is leaked, it is **COMPLETELY UNRELATED** to this program at the technical level. Therefore, please do not post related issues here, but instead consult the relevant platform's customer service.

This disclaimer is only valid for the [main open-source release](https://github.com/Nixdorfer/TavernHelper) of this program.

## Tech Stack

Framework: Tauri v2

Backend: Rust

Frontend: Vue3 + TypeScript @ClaudeCode

Database: SQLite

MCP Server: TypeScript

## Closing Words

- Welcome to use this tool! Feedback is always appreciated!

- This project is open-sourced under the MIT license, which means you can freely use it for commercial purposes. However, I would appreciate it if you could include my name 👍🏿

## Currently Supported Platforms

- [风月](https://aify.pages.dev/)

## Contact

QQ Group `1074532978`

[Telegram](https://t.me/Nixdorfer) or copy `Nixdorfer`

[Email](mailto:admin@nixdorfer.com) or copy `admin@nixdorfer.com`