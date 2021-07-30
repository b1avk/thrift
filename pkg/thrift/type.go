package thrift

// TType type of value.
type TType = byte

const (
	STOP TType = iota
	VOID
	BOOL
	BYTE
	DOUBLE
	U16
	I16
	U32
	I32
	U64
	I64
	STRING
	STRUCT
	MAP
	SET
	LIST
)

// TMessageType type of message used in TMessageHeader.
type TMessageType = byte

const (
	CALL TMessageType = iota + 1
	REPLY
	EXCEPTION
	ONEWAY
)

// TMessageHeader message header.
type TMessageHeader struct {
	Name     string
	Type     TMessageType
	Identity int32
}

// TStructHeader struct header.
type TStructHeader struct {
	Name string
}

// TFieldHeader field header.
type TFieldHeader struct {
	Name     string
	Type     TType
	Identity int16
}

// TMapHeader map header.
type TMapHeader struct {
	Key, Value TType
	Size       int
}

// TSetHeader set header.
type TSetHeader struct {
	Element TType
	Size    int
}

// TListHeader list header.
type TListHeader struct {
	Element TType
	Size    int
}

// TStruct is interface that wraps Read and Write methods.
type TStruct interface {
	// Read reads value from p.
	Read(p TProtocol) (err error)

	// Write writes value to p.
	Write(p TProtocol) (err error)
}
