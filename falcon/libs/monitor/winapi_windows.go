// +build windows

package monitor

import (
	"errors"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

const (
	pathBufferSize = 8192
	dllPath  = "winapi.dll"
)

var (
	isUnZip bool
	winApi syscall.Handle
	getProcessWorkDirByHandle uintptr
	getProcessWorkDirByPid uintptr
	getProcessCommandLineByHandle uintptr
	getProcessCommandLineByPid uintptr
	getProcessorCount uintptr
)


func LoadWinDll() {
	if !isUnZip {
		isUnZip = true
		_ = RestoreAssets("./", dllPath)
	}

	if winApi == 0 {
		winApi, _ = syscall.LoadLibrary(dllPath)
	}

	if getProcessWorkDirByHandle == 0 {
		getProcessWorkDirByHandle, _ = syscall.GetProcAddress(winApi, "GetProcessWorkDirByHandle")
	}

	if getProcessWorkDirByPid == 0 {
		getProcessWorkDirByPid, _ = syscall.GetProcAddress(winApi, "GetProcessWorkDirByPid")
	}

	if getProcessCommandLineByHandle == 0 {
		getProcessCommandLineByHandle, _ = syscall.GetProcAddress(winApi, "GetProcessCommandLineByHandle")
	}

	if getProcessCommandLineByPid == 0 {
		getProcessCommandLineByPid, _ = syscall.GetProcAddress(winApi, "GetProcessCommandLineByPid")
	}

	if getProcessorCount == 0 {
		getProcessorCount, _ = syscall.GetProcAddress(winApi, "GetProcessorCount")
	}
}

func GetProcessWorkDirByPid(pid int) (string, error) {
	LoadWinDll()
	if winApi == 0 || getProcessWorkDirByPid == 0 {
		return "", errors.New("Load api fail")
	}

	buffer := make([]uint16, pathBufferSize)
	bufferSize := pathBufferSize
	r1, _, _ := syscall.Syscall(uintptr(getProcessWorkDirByPid), 3, uintptr(pid), uintptr(unsafe.Pointer(&buffer[0])), uintptr(unsafe.Pointer(&bufferSize)))
	if r1 == 0 || bufferSize == 0 {
		return "", errors.New("Get work dir fail")
	}

	return windows.UTF16ToString(buffer[:bufferSize]), nil
}

func GetProcessWorkDirByHandle(handle windows.Handle) (string, error) {
	LoadWinDll()
	if winApi == 0 || getProcessWorkDirByPid == 0 {
		return "", errors.New("Load api fail")
	}

	buffer := make([]uint16, pathBufferSize)
	bufferSize := pathBufferSize
	r1, _, _ := syscall.Syscall(uintptr(getProcessWorkDirByHandle), 3, uintptr(handle), uintptr(unsafe.Pointer(&buffer[0])), uintptr(unsafe.Pointer(&bufferSize)))
	if r1 == 0 || bufferSize == 0 {
		return "", errors.New("Get work dir fail")
	}

	return windows.UTF16ToString(buffer[:bufferSize]), nil
}

func GetProcessCommandLineByPid(pid int) (string, error) {
	LoadWinDll()
	if winApi == 0 || getProcessWorkDirByPid == 0 {
		return "", errors.New("Load api fail")
	}

	buffer := make([]uint16, pathBufferSize)
	bufferSize := pathBufferSize
	r1, _, _ := syscall.Syscall(uintptr(getProcessCommandLineByPid), 3, uintptr(pid), uintptr(unsafe.Pointer(&buffer[0])), uintptr(unsafe.Pointer(&bufferSize)))
	if r1 == 0 || bufferSize == 0 {
		return "", errors.New("Get cmd line fail")
	}

	return windows.UTF16ToString(buffer[:bufferSize]), nil
}

func GetProcessCommandLineByHandle(handle windows.Handle) (string, error) {
	LoadWinDll()
	if winApi == 0 || getProcessWorkDirByPid == 0 {
		return "", errors.New("Load api fail")
	}

	buffer := make([]uint16, pathBufferSize)
	bufferSize := pathBufferSize
	r1, _, _ := syscall.Syscall(uintptr(getProcessCommandLineByHandle), 3, uintptr(handle), uintptr(unsafe.Pointer(&buffer[0])), uintptr(unsafe.Pointer(&bufferSize)))
	if r1 == 0 || bufferSize == 0 {
		return "", errors.New("Get cmd line fail")
	}

	return windows.UTF16ToString(buffer[:bufferSize]), nil
}

func GetCpuCountFromWinApi() (int, int, error){
	LoadWinDll()
	if winApi == 0 || getProcessorCount == 0 {
		return 0, 0, errors.New("Load api fail")
	}
	logicalCount := 0
	physicsCount := 0
	r1, _, _ := syscall.Syscall(uintptr(getProcessorCount), 2, uintptr(unsafe.Pointer(&logicalCount)), uintptr(unsafe.Pointer(&physicsCount)), 0)
	if r1 == 0 {
		return 0, 0, errors.New("Get cpu count fail")
	}
	return logicalCount, physicsCount, nil
}
