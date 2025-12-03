package wayland

import (
	"fmt"
	"net"
	"syscall"
)

func SendMessageAndFileDescriptors(conn *net.UnixConn, buf []byte, fds []int) error {
	if len(buf) == 0 {
		return nil
	}

	total := 0
	oobFirst := syscall.UnixRights(fds...)

	for total < len(buf) {
		chunk := buf[total:]
		var oob []byte
		if total == 0 {
			oob = oobFirst // send FDs only once
		}

		n, _, err := conn.WriteMsgUnix(chunk, oob, nil)
		if err != nil || n == -1 {
			var out_error error
			if err != nil {
				out_error = err
			} else {
				out_error = fmt.Errorf("N -1 on send message")
			}
			return out_error
		}
		total += n
	}

	return nil
}
