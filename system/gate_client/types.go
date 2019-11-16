package gate_client

type IWsCallback interface {
	onMessage(payload []byte)
	onConnected()
	onClosed()
}

const (
	ClientTypeServer = "server"
)

const (
	Request       = "request"
	Response      = "response"
	StatusSuccess = "success"
	StatusError   = "error"
)

const (
	GateStatusWait      = "wait"
	GateStatusDisabled  = "disabled"
	GateStatusConnected = "connected"
)