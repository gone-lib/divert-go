//go:build windows && !divert_cgo && !divert_embedded && (amd64 || arm64)
// +build windows
// +build !divert_cgo
// +build !divert_embedded
// +build amd64 arm64

package divert

import (
	"runtime"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

func open(filter string, layer Layer, priority int16, flags uint64) (h *Handle, err error) {
	if priority < PriorityLowest || priority > PriorityHighest {
		return nil, errPriority
	}

	filterPtr, err := windows.BytePtrFromString(filter)
	if err != nil {
		return nil, err
	}

	runtime.LockOSThread()
	hd, _, err := winDivertOpen.Call(uintptr(unsafe.Pointer(filterPtr)), uintptr(layer), uintptr(priority), uintptr(flags))
	runtime.UnlockOSThread()

	if windows.Handle(hd) == windows.InvalidHandle {
		return nil, Error(err.(windows.Errno))
	}

	rEvent, _ := windows.CreateEvent(nil, 0, 0, nil)
	wEvent, _ := windows.CreateEvent(nil, 0, 0, nil)

	return &Handle{
		Mutex:  sync.Mutex{},
		Handle: windows.Handle(hd),
		rOverlapped: windows.Overlapped{
			HEvent: rEvent,
		},
		wOverlapped: windows.Overlapped{
			HEvent: wEvent,
		},
		Open: HandleOpen,
	}, nil
}

func HelperCalcChecksum(packet *Packet, flags ChecksumFlag) bool {
	result, _, err := winDivertHelperCalcChecksums.Call(
		uintptr(unsafe.Pointer(&packet.Content[0])),
		uintptr(len(packet.Content)),
		uintptr(unsafe.Pointer(&packet.Addr)),
		uintptr(flags))

	if err != nil {
		return false
	}

	return result == 1
}
