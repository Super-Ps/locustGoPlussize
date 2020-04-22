// +build windows

package monitor

import "syscall"


func _init(){
	user32, _ = syscall.LoadLibrary("User32.dll")
	kernel32, _ = syscall.LoadLibrary("Kernel32.dll")
	psApi, _ = syscall.LoadLibrary("Psapi.dll")
	ntDll, _ = syscall.LoadLibrary("Ntdll.dll")

	getGuiResources, _ = syscall.GetProcAddress(user32, "GetGuiResources")
	getProcessCounters, _ = syscall.GetProcAddress(kernel32, "GetProcessIoCounters")
	getProcessHandleCount, _ = syscall.GetProcAddress(kernel32, "GetProcessHandleCount")
	getProcessImageFileNameW, _ = syscall.GetProcAddress(psApi, "GetProcessImageFileNameW")
	getProcessMemoryInfo, _ = syscall.GetProcAddress(psApi, "GetProcessMemoryInfo")
	rtlAdjustPrivilege, _ = syscall.GetProcAddress(ntDll, "RtlAdjustPrivilege")

	dosPath = getDOSPathMap()
	_ = setDebugPrivilege(true)
}