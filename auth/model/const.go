package model

const (
	AuthTokenName      = "X-Webitel-Access"
	RequestContextName = "grpc_ctx"
)

type AccessMode uint8

func (a AccessMode) Value() uint8 {
	return uint8(a)
}

const (
	Delete AccessMode = 1 << iota
	Edit
	Read
	Add

	NONE AccessMode = 0
	FULL            = Add | Read | Edit | Delete
)

const (
	SuperSelectPermission = "read"
	SuperEditPermission   = "write"
	SuperCreatePermission = "add"
	SuperDeletePermission = "delete"
)
