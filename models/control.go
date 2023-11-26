package models

type VmControlByUUIDInput struct {
	VmControlInput
	UUID UUID `json:"uuid"`
}

type VmControlInput struct {
	Host string `json:"host"`
}
