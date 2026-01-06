package main

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows"
)

func (a *App) IgnoreClipboardContent(content string) {
	a.clipboardMonitorMutex.Lock()
	defer a.clipboardMonitorMutex.Unlock()
	if a.ignoredClipboard == nil {
		a.ignoredClipboard = make(map[string]bool)
	}
	a.ignoredClipboard[content] = true
}

func (a *App) CopyToClipboard(content string) error {
	a.IgnoreClipboardContent(content)

	user32 := windows.NewLazySystemDLL("user32.dll")
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")

	openClipboard := user32.NewProc("OpenClipboard")
	closeClipboard := user32.NewProc("CloseClipboard")
	emptyClipboard := user32.NewProc("EmptyClipboard")
	setClipboardData := user32.NewProc("SetClipboardData")
	globalAlloc := kernel32.NewProc("GlobalAlloc")
	globalLock := kernel32.NewProc("GlobalLock")
	globalUnlock := kernel32.NewProc("GlobalUnlock")

	ret, _, _ := openClipboard.Call(0)
	if ret == 0 {
		return fmt.Errorf("无法打开剪贴板")
	}
	defer closeClipboard.Call()

	emptyClipboard.Call()

	utf16, _ := windows.UTF16FromString(content)
	size := len(utf16) * 2

	const GMEM_MOVEABLE = 0x0002
	hMem, _, _ := globalAlloc.Call(GMEM_MOVEABLE, uintptr(size))
	if hMem == 0 {
		return fmt.Errorf("无法分配内存")
	}

	ptr, _, _ := globalLock.Call(hMem)
	if ptr == 0 {
		return fmt.Errorf("无法锁定内存")
	}

	for i, v := range utf16 {
		*(*uint16)(unsafe.Pointer(ptr + uintptr(i*2))) = v
	}
	globalUnlock.Call(hMem)

	const CF_UNICODETEXT = 13
	setClipboardData.Call(CF_UNICODETEXT, hMem)

	return nil
}

func (a *App) StartClipboardMonitor() {
	a.clipboardMonitorMutex.Lock()
	if a.clipboardMonitorStop != nil {
		a.clipboardMonitorMutex.Unlock()
		return
	}
	a.clipboardMonitorStop = make(chan struct{})
	a.ignoredClipboard = make(map[string]bool)
	a.clipboardMonitorMutex.Unlock()

	go func() {
		user32 := windows.NewLazySystemDLL("user32.dll")
		kernel32 := windows.NewLazySystemDLL("kernel32.dll")

		openClipboard := user32.NewProc("OpenClipboard")
		closeClipboard := user32.NewProc("CloseClipboard")
		getClipboardData := user32.NewProc("GetClipboardData")
		isClipboardFormatAvailable := user32.NewProc("IsClipboardFormatAvailable")
		globalLock := kernel32.NewProc("GlobalLock")
		globalUnlock := kernel32.NewProc("GlobalUnlock")

		const CF_UNICODETEXT = 13

		for {
			select {
			case <-a.clipboardMonitorStop:
				return
			default:
				time.Sleep(500 * time.Millisecond)

				ret, _, _ := isClipboardFormatAvailable.Call(CF_UNICODETEXT)
				if ret == 0 {
					continue
				}

				ret, _, _ = openClipboard.Call(0)
				if ret == 0 {
					continue
				}

				var content string
				hData, _, _ := getClipboardData.Call(CF_UNICODETEXT)
				if hData != 0 {
					ptr, _, _ := globalLock.Call(hData)
					if ptr != 0 {
						content = windows.UTF16PtrToString((*uint16)(unsafe.Pointer(ptr)))
						globalUnlock.Call(hData)
					}
				}
				closeClipboard.Call()
				if content == "" {
					continue
				}
				a.clipboardMonitorMutex.Lock()
				isIgnored := a.ignoredClipboard[content]
				lastContent := a.lastClipboardContent
				a.clipboardMonitorMutex.Unlock()
				if content != lastContent && !isIgnored {
					a.clipboardMonitorMutex.Lock()
					a.lastClipboardContent = content
					a.clipboardMonitorMutex.Unlock()
					now := time.Now().Format(time.RFC3339)
					capture := &ClipboardCapture{
						ID:        fmt.Sprintf("clip_%d", time.Now().UnixNano()),
						Content:   content,
						CreatedAt: now,
					}
					DBAddClipboardCapture(capture)
					DBCleanupOldClipboardCaptures(50)
					runtime.EventsEmit(a.ctx, "clipboard-captured", capture)
				}
			}
		}
	}()
}

func (a *App) StopClipboardMonitor() {
	a.clipboardMonitorMutex.Lock()
	defer a.clipboardMonitorMutex.Unlock()
	if a.clipboardMonitorStop != nil {
		close(a.clipboardMonitorStop)
		a.clipboardMonitorStop = nil
	}
}
