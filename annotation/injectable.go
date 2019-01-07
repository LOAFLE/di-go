package annotation

// @Inject(name? string)

import (
	"reflect"

	la "github.com/LOAFLE/annotation-go"
)

func init() {
	la.Register(InjectableAnnotationType)
}

var InjectableAnnotationType = reflect.TypeOf((*InjectableAnnotation)(nil))

type InjectableAnnotation struct {
	la.Annotation `@name:"@Injectable"`
	Name          string `json:"name" @default:""`
}
