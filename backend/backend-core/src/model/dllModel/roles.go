package dllModel

type Role struct {
	ID          uint32
	Label       string
	Permissions []Permission
}

type Permission struct {
	UID   string
	Label string
}
