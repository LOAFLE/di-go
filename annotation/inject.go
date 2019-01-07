package annotation

// @Inject(name? string)

import (
	"reflect"

	la "github.com/LOAFLE/annotation-go"
)

func init() {
	la.Register(InjectAnnotationType)
}

var InjectAnnotationType = reflect.TypeOf((*InjectAnnotation)(nil))

type InjectAnnotation struct {
	la.Annotation `@name:"@Inject"`
	Name          string `json:"name" @default:""`
	Required      bool   `json:"required" @default:"true"`
}
