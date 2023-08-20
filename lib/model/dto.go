package model

type InterfaceTransferObject struct {
	Path      string    `json:"path"`
	Imports   []Import  `json:"imports"`
	Interface Interface `json:"interface"`
}

type MethodTransferObject struct {
	Path    string   `json:"path"`
	Imports []Import `json:"imports"`
	Method  Method   `json:"method"`
}
