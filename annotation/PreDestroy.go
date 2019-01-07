package annotation

// @Inject(name? string)

import (
	"reflect"

	la "github.com/LOAFLE/annotation-go"
)

func init() {
	la.Register(PreDestroyAnnotationType)
}

var PreDestroyAnnotationType = reflect.TypeOf((*PreDestroyAnnotation)(nil))

type PreDestroyAnnotation struct {
	la.Annotation `@name:"@PreDestroy"`
}
