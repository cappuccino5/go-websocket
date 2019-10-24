package chat

import "strconv"

const (
	// send webSocket message error
	ERR01 = ERRCode(0x01)
	ERR02 = ERRCode(0x02)
	ERR03 = ERRCode(0x03)
	ERR04 = ERRCode(0x04)
)

var strings = map[ERRCode]string{
	ERR01: "send message timeout",
	ERR02: "invalidate user's id",
	ERR03: "connection not found",
	ERR04: "json marshal error",
}

type ERRCode int

func (e ERRCode) Error() string {
	if 0 <= int(e) {
		s := strings[e]
		if s != "" {
			return s
		}
	}
	return "error code " + strconv.Itoa(int(e))
}

func (e ERRCode) IsEqual(err error) bool {
	return e == err
}

func (e ERRCode) IsStrEqual(err error) bool {
	return e.Error() == err.Error()
}
