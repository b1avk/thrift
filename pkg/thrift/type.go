package thrift

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
	LIST
	SET
)

type TMessageType = byte

const (
	CALL TMessageType = iota + 1
	REPLY
	EXCEPTION
	ONEWAY
)

type TMessageHeader struct {
	Name     string
	Type     TMessageType
	Identity int32
}

type TStructHeader struct {
	Name string
}

type TFieldHeader struct {
	Name     string
	Type     TType
	Identity int16
}

type TMapHeader struct {
	Key, Value TType
	Size       int
}

type TSetHeader struct {
	Element TType
	Size    int
}

type TListHeader struct {
	Element TType
	Size    int
}

type TStruct interface {
	Read(p TProtocol) (err error)

	Write(p TProtocol) (err error)
}
