[简体中文](../README.md) | [English](./README_EN.md) | [繁體中文](./README_TC.md)

# Über dieses Projekt

Probleme? [Klicken Sie hier, um zu berichten](https://github.com/Nixdorfer/TavernHelper/issues/new)

Sie wissen nicht, wie man es benutzt? [Klicken Sie hier für die Anleitung](./MANUAL_DE.md)

## Motivation

Ich habe angefangen, Tavern über 风月 zu nutzen. Ihre Produkte sind ziemlich gut, aber die Bearbeitungserfahrung ist recht frustrierend. Als jemand, der gerne ausführliche Prompts schreibt (10k+ Token pro Charakter), zielt dieses Programm darauf ab, das Prompting-Erlebnis zu verbessern. Daher ist es standardmäßig mit der 风月-API integriert. Wenn Sie eine Integration mit anderen Plattformen wünschen, können Sie oben ein Issue einreichen ↑ und ich werde es in Betracht ziehen.

## Download & Ausführung

Als normaler Benutzer können Sie die ausführbare Datei für Ihr System von der [Releases-Seite](https://github.com/Nixdorfer/TavernHelper/releases) herunterladen.

## Kompilieren / Debug-Start

Führen Sie Run.cmd mit den folgenden verfügbaren Parametern aus:

- -d Debug-Modus: Kompiliert eine ausführbare Datei mit Debugger, startet Entwicklertools zur Laufzeit und zeigt eine Debug-Seite in der Seitenleiste an

- -r Rust/Tauri-Modus: Dieses Projekt verwendet derzeit Wails als Programm-Framework. Dieses Flag startet die Rust-Version

- -m Mobiler Start: Dies startet die kompilierte Anwendung im Android Studio-Emulator

- -b Nur Build: Dies kompiliert keine Dev-Version mit Debugger. Es leert das bin-Verzeichnis und erstellt eine Release-Version der ausführbaren Datei

- -c Sauberer Start: Dies löscht alle kompilierten Ausgaben und den Cache (einschließlich Datenbank, aber nicht den creations-Ordner) vor dem Kompilieren und Starten

## TODO

- Tauri v2 Refactoring -> Android/iOS-Unterstützung

- UDP LAN-Multicast mit Mutex-Dateiänderungsunterstützung

- Unterstützung für direkte Bildgenerierung über lokales oder Netzwerk-SD-WebUI

- MCP-Unterstützung

  - Neuer Server mit MCP-Aufruf-Unterstützung

  - Integriertes automatisches Port-Mapping

- Meimodo API-Unterstützung

- Unterstützung für SessionKey + SSE für Konto-Abonnement-API-Aufrufe

## Haftungsausschluss

Dieses Programm hat **NIEMALS** und **WIRD NIEMALS** Ihre Informationen an einen anderen Ort als die **OFFIZIELLE API** der von Ihnen angegebenen Plattform übertragen.

Das bedeutet: Wenn Ihre Informationen durchsickern, ist dies auf technischer Ebene **VÖLLIG UNABHÄNGIG** von diesem Programm. Bitte posten Sie daher keine entsprechenden Issues hier, sondern wenden Sie sich an den Kundenservice der jeweiligen Plattform.

Dieser Haftungsausschluss gilt nur für die [Haupt-Open-Source-Version](https://github.com/Nixdorfer/TavernHelper) dieses Programms.

## Technologie-Stack

Framework: Tauri v2 / Wails

Backend: Rust / Go

Frontend: Vue3 + TypeScript @ClaudeCode

Datenbank: SQLite

MCP-Server: TypeScript

## Schlusswort

- Willkommen zur Nutzung dieses Tools! Feedback ist immer willkommen!

- Dieses Projekt ist unter der MIT-Lizenz Open Source, was bedeutet, dass Sie es frei für kommerzielle Zwecke nutzen können. Ich würde mich jedoch freuen, wenn Sie meinen Namen erwähnen könnten 👍🏿

## Derzeit unterstützte Plattformen

- [风月](https://aify.pages.dev/)

## Kontakt

QQ-Gruppe `1074532978`

[Telegram](https://t.me/Nixdorfer) oder kopieren Sie `Nixdorfer`

[E-Mail](mailto:admin@nixdorfer.com) oder kopieren Sie `admin@nixdorfer.com`