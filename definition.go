package di

import (
	"fmt"
	"reflect"
)

type TypeDefinition struct {
	FullName string
	PkgName  string
	TypeName string
	Type     reflect.Type
	RealType reflect.Type
}

func FullName(pkgName, typeName string) string {
	return fmt.Sprintf("%s/%s", pkgName, typeName)
}
