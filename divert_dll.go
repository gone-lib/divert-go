//go:build windows && !divert_cgo && !divert_embedded

package divert

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/sys/windows"
)

var (
	winDivert                    = (*windows.DLL)(nil)
	winDivertOpen                = (*windows.Proc)(nil)
	winDivertClose               = (*windows.Proc)(nil)
	winDivertRecv                = (*windows.Proc)(nil)
	winDivertSend                = (*windows.Proc)(nil)
	winDivertHelperCalcChecksums = (*windows.Proc)(nil)
	winDivertHelperEvalFilter    = (*windows.Proc)(nil)
	//winDivertHelperCheckFilter    	= (*windows.Proc)(nil)
)

func Open(filter string, layer Layer, priority int16, flags uint64) (h *Handle, err error) {
	once.Do(func() {
		winDivert, err = windows.LoadDLL("WinDivert.dll")
		if err != nil {
			return
		}

		// WinDivertOpen
		winDivertOpen, err = winDivert.FindProc("WinDivertOpen")
		if err != nil {
			return
		}

		// WinDivertHelperCalcChecksums
		winDivertHelperCalcChecksums, err = winDivert.FindProc("WinDivertHelperCalcChecksums")
		if err != nil {
			return
		}

		vers := map[string]struct{}{
			"2.0": {},
			"2.1": {},
			"2.2": {},
		}

		h, err := open("false", LayerNetwork, PriorityDefault, FlagDefault)
		if err != nil {
			return
		}
		defer func() {
			err = h.Close()
		}()

		major, err := h.GetParam(VersionMajor)
		if err != nil {
			return
		}

		minor, err := h.GetParam(VersionMinor)
		if err != nil {
			return
		}

		ver := strings.Join([]string{strconv.Itoa(int(major)), strconv.Itoa(int(minor))}, ".")

		if _, ok := vers[ver]; !ok {
			err = fmt.Errorf("unsupported windivert version: %v", ver)
		}
	})
	if err != nil {
		return
	}

	return open(filter, layer, priority, flags)
}
