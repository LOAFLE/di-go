package annotation

// @Inject(name? string)

import (
	"reflect"

	la "github.com/LOAFLE/annotation-go"
)

func init() {
	la.Register(PostConstructAnnotationType)
}

var PostConstructAnnotationType = reflect.TypeOf((*PostConstructAnnotation)(nil))

type PostConstructAnnotation struct {
	la.Annotation `@name:"@PostConstruct"`
}
