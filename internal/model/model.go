package model

type Structs struct {
	PackageName string
	StructNames []string
}

type Field struct {
	Name string
	Type string
	Tags map[string]string
}

type Struct struct {
	PackageName string
	StructName  string
	Fields      []Field
}
