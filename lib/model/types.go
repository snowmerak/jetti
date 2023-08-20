package model

type Field struct {
	Name     string    `json:"name,omitempty"`
	Type     string    `json:"type,omitempty"`
	FuncType *FuncType `json:"funcType,omitempty"`
	Tag      string    `json:"tag,omitempty"`
}

type FuncType struct {
	Params []Field
	Return []Field
}

type Package struct {
	Name            string
	Imports         []Import
	Structs         []Struct
	Interfaces      []Interface
	Functions       []Function
	Methods         []Method
	Aliases         []Alias
	GlobalVariables []GlobalVariable
}

type Import struct {
	Alias string
	Path  string
}

type Struct struct {
	Doc     string
	Name    string
	Fields  []Field
	Methods []Method
}

type Method struct {
	Receiver Field    `json:"receiver"`
	Name     string   `json:"name"`
	Params   []Field  `json:"params"`
	Return   []Field  `json:"return"`
	Code     []string `json:"code"`
}

type Interface struct {
	Doc     string   `json:"doc"`
	Name    string   `json:"name"`
	Methods []Method `json:"methods"`
}

type Function struct {
	Doc      string
	Name     string
	Receiver string
	Params   []Field
	Return   []Field
	Code     []string
}

type Alias struct {
	Doc  string
	Name string
	Type string
}

type GlobalVariable struct {
	Doc   string
	Name  string
	Type  string
	Value string
}
