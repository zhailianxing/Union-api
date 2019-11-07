package packet

type OkPacket struct {
	Code    uint32      `json:"code"`
	Data    interface{} `json:"data,,omitempty"`
	Message string      `json:"message,omitempty"`
}
