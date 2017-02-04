package interfaces

type MsgpackMarshaller interface {
	MarshalMsgpack() ([]byte, error)
	UnmarshalMsgpack(b []byte) error
}
