package model

type Field struct {
	Name string
	Type string
	Tag  string
}

type Package struct {
	Name       string
	Imports    []Import
	Structs    []Struct
	Interfaces []Interface
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
	Name   string
	Params []Field
	Return []Field
}

type Interface struct {
	Doc     string
	Name    string
	Methods []Method
}