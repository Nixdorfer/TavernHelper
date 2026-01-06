[简体中文](../README.md) | [English](./README_EN.md) | [Deutsch](./README_DE.md)

# 關於本專案

使用時遇到了問題？[點我吐槽](https://github.com/Nixdorfer/TavernHelper/issues/new)

不知道怎麼用？[點我看說明書](./MANUAL_TC.md)

## 初衷

我是從風月開始玩酒館的，他們家感覺做的東西還挺不錯的，但是編輯實在過於惡臭了，而我又恰好是一個題詞狂魔（單人物10kt+），因此該程式的目的旨在於改善提詞體驗。因此當然預設對接了風月的API，要對接其他的平台可以在上面↑給我發個issue，酌情開咕。

## 下載執行

如果您是一般使用者，可以從[發布頁](https://github.com/Nixdorfer/TavernHelper/releases)下載適用於您系統的可執行檔。

## 編譯/除錯啟動

執行 Run.cmd 可用參數如下：

- -d 除錯模式：該模式下會編譯帶有除錯器的可執行檔，並且在執行時啟動開發者工具，還會在sidebar中顯示除錯頁

- -m 行動端啟動：這將在Android Studio模擬器中啟動編譯後的應用程式

- -b 僅構建：這將不會編譯帶有除錯器的dev版本，會清空bin目錄後構建一個release版本的可執行程式

- -c 清潔啟動：這將會先清空所有編譯產物和快取（包括資料庫，但不包括creations資料夾）後再編譯啟動

## TODO

- Tauri v2重構 -> 支援Android/iOS

- UDP內網多播，支援Mutex的檔案修改

- 支援通過本地或網路SD-WebUI直接生成圖片

- MCP支援

  - 新建支援MCP呼叫的伺服器端

  - 內建自動埠映射

- 支援魅魔島介面

- 支援使用SessionKey + SSE通過帳號訂閱呼叫介面

## 聲明

該程式**從未**並且**在將來也不會**以任何方式將您的任何資訊傳送至除了您指定的平台的**官方介面**以外的其他位置。

這代表：如果您的資訊洩露，在技術層面與該程式**完全無關**，因此請不要在issue中發布相關問題，而是請諮詢相關平台客服。

該聲明僅對於本程式的[主要開源發行版本](https://github.com/Nixdorfer/TavernHelper)有效。

## 技術棧

框架：Tauri v2

後端：Rust

前端：Vue3 + TypeScript @ClaudeCode

資料庫：SQLite

MCP伺服器端：TypeScript

## 結語

- 總之！歡迎使用！歡迎提出意見！

- 該專案基於MIT協議開源，這意味著您可以自由用於商用，但是如果能附上我的名字我會很感謝的👍🏿

## 目前介面支援平台

- [風月](https://aify.pages.dev/)

## 聯絡方式

QQ群 `1074532978`

[Telegram](https://t.me/Nixdorfer) 或複製 `Nixdorfer`

[電子郵箱](mailto:admin@nixdorfer.com) 或複製 `admin@nixdorfer.com`