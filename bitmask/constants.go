package bitmask

type Command uint32

const (
	C_VERSION Command = iota
)

type Version uint32

const (
	V_0_1_0 Version = iota
)

type Subnet uint32

const (
	S_MAIN Subnet = iota
)

type Mode uint32

const (
	M_NORMAL Mode = iota
)
