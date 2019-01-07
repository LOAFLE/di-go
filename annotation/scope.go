package annotation

// @Scope(name? string)

import (
	"encoding/json"
	"fmt"
	"reflect"

	la "github.com/LOAFLE/annotation-go"
)

func init() {
	la.Register(ScopeAnnotationType)
}

var ScopeAnnotationType = reflect.TypeOf((*ScopeAnnotation)(nil))

type ScopeAnnotation struct {
	la.Annotation `@name:"@Scope"`
	Value         ScopeType `json:"value" @default:"singleton"`
}

type ScopeType int8

const (
	ScopeTypeDefault ScopeType = iota + 1
	ScopeTypeSingleton
	ScopeTypeTransiant
)

var (
	scopeTypeID = map[ScopeType]string{
		ScopeTypeDefault:   "default",
		ScopeTypeSingleton: "singleton",
		ScopeTypeTransiant: "transiant",
	}

	scopeTypeKey = map[string]ScopeType{
		"default":   ScopeTypeDefault,
		"singleton": ScopeTypeSingleton,
		"transiant": ScopeTypeTransiant,
	}
)

func (st ScopeType) String() string {
	return scopeTypeID[st]
}

func (st ScopeType) MarshalJSON() ([]byte, error) {
	value, ok := scopeTypeID[st]
	if !ok {
		return nil, fmt.Errorf("Invalid EnumType[%s] value", st)
	}
	return json.Marshal(value)
}

func (st ScopeType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	value, ok := scopeTypeKey[s]
	if !ok {
		return fmt.Errorf("Invalid EnumType[%s] value", s)
	}
	st = value
	return nil
}
