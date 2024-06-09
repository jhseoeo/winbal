//go:build cgo
// +build cgo

package main

/*
#include <windows.h>
#include <stdio.h>

void mouseMove(int x, int y)
{
	int nScreenWidth = GetSystemMetrics(SM_CXVIRTUALSCREEN);
	int nScreenHeight = GetSystemMetrics(SM_CYVIRTUALSCREEN);
	int nScreenLeft = GetSystemMetrics(SM_XVIRTUALSCREEN);
	int nScreenTop = GetSystemMetrics(SM_YVIRTUALSCREEN);
    int nX = (int)((((double)(x)-nScreenLeft) * 65536) / nScreenWidth + 65536 / (nScreenWidth));
    int nY = (int)((((double)(y)-nScreenTop) * 65536) / nScreenHeight + 65536 / (nScreenHeight));
    printf("nScreenWidth: %d, nScreenHeight: %d, nScreenLeft: %d, nScreenTop: %d, nX: %d, nY: %d\n", nScreenWidth, nScreenHeight, nScreenLeft, nScreenTop, nX, nY);
    fflush(stdout);

	INPUT ip = { 0 };
	ip.type = INPUT_MOUSE;
	ip.mi.dwFlags = MOUSEEVENTF_ABSOLUTE | MOUSEEVENTF_MOVE;
	ip.mi.dx = nX;
	ip.mi.dy = nY;
	ip.mi.time = GetTickCount();
	ip.mi.dwExtraInfo = GetMessageExtraInfo();
	SendInput(1, &ip, sizeof(INPUT));
}

void mouseDown(int btn) {
    INPUT ip = {0};
    ip.type = INPUT_MOUSE;
    if (btn == 0) {
        ip.mi.dwFlags = MOUSEEVENTF_LEFTDOWN;
    } else if (btn == 1) {
        ip.mi.dwFlags = MOUSEEVENTF_MIDDLEDOWN;
    } else if (btn == 2) {
        ip.mi.dwFlags = MOUSEEVENTF_RIGHTDOWN;
    } else {
        return;
    }
    UINT sent = SendInput(1, &ip, sizeof(INPUT));
    if (sent != 1) {
        printf("failed to send mouse down event: 0x%x\n", HRESULT_FROM_WIN32(GetLastError()));
    }
}

void mouseUp(int btn) {
    INPUT ip = {0};
    ip.type = INPUT_MOUSE;
    if (btn == 0) {
        ip.mi.dwFlags = MOUSEEVENTF_LEFTUP;
    } else if (btn == 1) {
        ip.mi.dwFlags = MOUSEEVENTF_MIDDLEUP;
    } else if (btn == 2) {
        ip.mi.dwFlags = MOUSEEVENTF_RIGHTUP;
    } else {
        return;
    }
    UINT sent = SendInput(1, &ip, sizeof(INPUT));
    if (sent != 1) {
        printf("failed to send mouse down event: 0x%x\n", HRESULT_FROM_WIN32(GetLastError()));
    }
}

void keyDown(WORD key) {
    UINT ScanCode = MapVirtualKey(key & 0xff, 0);
    INPUT ip;
    ip.type = INPUT_KEYBOARD;
    ip.ki.wScan = ScanCode;
    ip.ki.dwFlags = KEYEVENTF_SCANCODE;
    // ip.ki.time = 0;
    // ip.ki.wVk = 0;
    // ip.ki.dwExtraInfo = 0;
    UINT sent = SendInput(1, &ip, sizeof(INPUT));
    if (sent != 1) {
        printf("failed to send mouse down event: 0x%x\n", HRESULT_FROM_WIN32(GetLastError()));
    }
}

void keyUp(WORD key) {
    UINT ScanCode = MapVirtualKey(key & 0xff, 0) | 0x80;
    INPUT ip;
    ip.type = INPUT_KEYBOARD;
    ip.ki.wScan = ScanCode;
    ip.ki.dwFlags = KEYEVENTF_SCANCODE | KEYEVENTF_KEYUP;
    // ip.ki.time = 0;
    // ip.ki.wVk = 0;
    // ip.ki.dwExtraInfo = 0;
    UINT sent = SendInput(1, &ip, sizeof(INPUT));
    if (sent != 1) {
        printf("failed to send mouse down event: 0x%x\n", HRESULT_FROM_WIN32(GetLastError()));
    }
}
*/
import "C"
import "fmt"

func keyDown(keyCode int) {
	C.keyDown(C.WORD(keyCode))
	fmt.Println("keyDown for windows", keyCode)
}

func keyUp(keyCode int) {
	C.keyUp(C.WORD(keyCode))
	fmt.Println("keyUp for windows", keyCode)
}

func mouseMove(x, y int) {
	C.mouseMove(C.int(x), C.int(y))
	fmt.Println("mouseMove for windows", x, y)
}

func mouseDown(btn MouseBtnType) {
	C.mouseDown(C.int(btn))
	fmt.Println("mouseDown for windows", btn)
}

func mouseUp(btn MouseBtnType) {
	C.mouseUp(C.int(btn))
	fmt.Println("mouseUp for windows", btn)
}
