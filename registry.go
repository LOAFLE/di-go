package di

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	la "github.com/LOAFLE/annotation-go"
	"github.com/LOAFLE/di-go/annotation"
	lur "github.com/LOAFLE/util-go/reflect"
)

type Registry interface {
	RegisterType(t reflect.Type)
	RegisterSingleton(singleton interface{}) error
	RegisterSingletonByName(name string, singleton interface{}) error

	GetInstance(t reflect.Type) (interface{}, error)
	GetInstances(ts []reflect.Type) ([]interface{}, error)
	GetInstanceByName(name string) (interface{}, error)
	GetInstancesByAnnotationType(t reflect.Type) ([]interface{}, error)
}

func New(parent Registry) Registry {
	r := &InstanceRegistry{
		parent:           parent,
		definitionByType: make(map[reflect.Type]*TypeDefinition, 0),
		definitionByName: make(map[string]*TypeDefinition, 0),
		instanceByType:   make(map[reflect.Type]interface{}, 0),
		instanceByName:   make(map[string]interface{}, 0),
	}
	if nil == r.parent {
		r.parent = AppRegistry
	}
	return r
}

var AppRegistry = &InstanceRegistry{
	parent:           nil,
	definitionByType: make(map[reflect.Type]*TypeDefinition, 0),
	definitionByName: make(map[string]*TypeDefinition, 0),
	instanceByType:   make(map[reflect.Type]interface{}, 0),
	instanceByName:   make(map[string]interface{}, 0),
}

type InstanceRegistry struct {
	parent           Registry
	definitionByType map[reflect.Type]*TypeDefinition
	definitionByName map[string]*TypeDefinition
	instanceByType   map[reflect.Type]interface{}
	instanceByName   map[string]interface{}
}

func RegisterType(t reflect.Type) {
	AppRegistry.RegisterType(t)
}
func (r *InstanceRegistry) RegisterType(t reflect.Type) {
	if nil == t {
		log.Panicf("t[reflect.Type] is nil")
	}
	if !lur.IsTypeKind(t, reflect.Struct, true) {
		log.Panicf("t[reflect.Type] must be Struct but is %v", t)
	}

	td, err := r.buildDefinition(t)
	if nil != err {
		log.Panicf("DI: %v", err)
	}

	if _, ok := r.definitionByType[td.Type]; ok {
		log.Panicf("The type[%s] of Component is exist already", td.FullName)
	}
	r.definitionByType[td.Type] = td

	name := td.TypeName

	a, err := la.GetTypeAnnotation(td.Type, annotation.InjectableAnnotationType)
	if nil != err {
		log.Panicf("%v", err)
	}
	if nil != a {
		ia := a.(*annotation.InjectableAnnotation)
		if "" != strings.Trim(ia.Name, " ") {
			name = ia.Name
		}
	}

	if eTD, ok := r.definitionByName[name]; ok {
		log.Panicf("The name[%s] of Component is exist already type[%s]", name, eTD.FullName)
	}
	r.definitionByName[name] = td
}

func RegisterSingleton(singleton interface{}) error {
	return AppRegistry.RegisterSingleton(singleton)
}
func (r *InstanceRegistry) RegisterSingleton(singleton interface{}) error {
	t := reflect.TypeOf(singleton)

	if nil == t {
		log.Panicf("t[reflect.Type] is nil")
	}

	rt, _, tName := lur.GetTypeInfo(t)
	if !lur.IsTypeKind(rt, reflect.Struct, true) {
		log.Panicf("t[reflect.Type] must be Struct but is %v", t)
	}

	r.instanceByType[rt] = singleton

	r.RegisterSingletonByName(tName, singleton)

	return nil
}

func RegisterSingletonByName(name string, singleton interface{}) error {
	return AppRegistry.RegisterSingletonByName(name, singleton)
}
func (r *InstanceRegistry) RegisterSingletonByName(name string, singleton interface{}) error {
	t := reflect.TypeOf(singleton)

	if nil == t {
		log.Panicf("t[reflect.Type] is nil")
	}

	if _, ok := r.instanceByName[name]; ok {
		return fmt.Errorf("name[%s] of Singleton is already exist", name)
	}
	r.instanceByName[name] = singleton

	return nil
}

func (r *InstanceRegistry) buildDefinition(t reflect.Type) (*TypeDefinition, error) {
	if nil == t {
		return nil, fmt.Errorf("t[reflect.Type] is nil")
	}

	rt, pkgName, tName := lur.GetTypeInfo(t)
	td := &TypeDefinition{}
	td.FullName = FullName(pkgName, tName)
	td.PkgName = pkgName
	td.TypeName = tName
	td.Type = t
	td.RealType = rt

	return td, nil
}

// GetInstance returns instance of type
// t must be pointer of struct
func GetInstance(t reflect.Type) (interface{}, error) {
	return AppRegistry.GetInstance(t)
}
func (r *InstanceRegistry) GetInstance(t reflect.Type) (instance interface{}, err error) {
	if nil == t {
		return nil, fmt.Errorf("t[reflect.Type] is nil")
	}

	if reflect.Ptr != t.Kind() {
		return nil, fmt.Errorf("t[reflect.Type] must be pointer of struct")
	}

	if reflect.Struct != t.Elem().Kind() {
		return nil, fmt.Errorf("t[reflect.Type] must be pointer of struct")
	}

	var (
		td          *TypeDefinition
		injectableV interface{}
		ok          bool
		name        string
	)

	rt, _, _ := lur.GetTypeInfo(t)
	if td, ok = r.definitionByType[t]; !ok {
		if td, err = r.buildDefinition(t); nil != err {
			return nil, fmt.Errorf("DI: %v", err)
		}
	}

	if injectableV, ok = r.instanceByType[td.RealType]; ok {
		return injectableV, nil
	}

	v := reflect.New(rt)
	rv := v.Elem()

	instance = v.Interface()
	r.instanceByType[td.RealType] = instance

	if a, err := la.GetTypeAnnotation(td.Type, annotation.InjectableAnnotationType); nil == err && nil != a {
		_a := a.(*annotation.InjectableAnnotation)
		if "" != _a.Name {
			name = _a.Name
			r.instanceByName[name] = instance
		}
	}

	err = nil
	defer func() {
		if nil != err {
			instance = nil
			if _, ok := r.instanceByType[td.RealType]; ok {
				delete(r.instanceByType, td.RealType)
			}
			if _, ok := r.instanceByName[name]; ok {
				delete(r.instanceByName, name)
			}
		}
	}()

	var (
		a  la.Annotation
		fV interface{}
	)

	ass, err := la.GetAllFieldAnnotations(td.Type)
	if nil != err {
		return nil, fmt.Errorf("%v", err)
	}
	if nil != ass {
	LOOP:
		for n, as := range ass {
			f := rv.FieldByName(n)

			if !f.IsValid() {
				err = fmt.Errorf("Field[%s] is not valid", n)
				return
			}
			if !f.CanSet() {
				err = fmt.Errorf("Field[%s] can not set", n)
				return
			}

			a, ok = as[annotation.InjectAnnotationType]
			if ok {
				_a := a.(*annotation.InjectAnnotation)
				if "" == _a.Name {
					if fV, err = r.GetInstance(f.Type()); nil == err {
						log.Printf("%s of %s injected by type[%s]", n, td.RealType.Name(), reflect.TypeOf(fV))
						f.Set(reflect.ValueOf(fV))
						continue LOOP
					} else {
						err = fmt.Errorf("cannot find instance for %s[Type:%s]", n, f.Type())
					}

					if fV, err = r.GetInstanceByName(n); nil == err {
						log.Printf("%s of %s injected by name[%s]", n, td.RealType.Name(), n)
						f.Set(reflect.ValueOf(fV))
						continue LOOP
					} else {
						err = fmt.Errorf("cannot find instance for %s[Name:%s]", n, n)
					}
				} else {
					if fV, err = r.GetInstanceByName(_a.Name); nil == err {
						log.Printf("%s of %s injected by name[%s]", n, td.RealType.Name(), _a.Name)
						f.Set(reflect.ValueOf(fV))
						continue LOOP
					} else {
						err = fmt.Errorf("cannot find instance for %s[Name:%s]", n, _a.Name)
					}

					if fV, err = r.GetInstance(f.Type()); nil == err {
						log.Printf("%s of %s injected by type[%s]", n, td.RealType.Name(), reflect.TypeOf(fV))
						f.Set(reflect.ValueOf(fV))
						continue LOOP
					} else {
						err = fmt.Errorf("cannot find instance for %s[Type:%s]", n, f.Type())
					}
				}

				if nil != err {
					if _a.Required {
						return
					}
				}
			}
		}
	}

	if as, err := la.GetMethodAnnotationsByType(td.Type, annotation.PostConstructAnnotationType); nil == err && nil != as {
		for k := range as {
			ins := make([]reflect.Value, 0)
			reflect.ValueOf(instance).MethodByName(k).Call(ins)
		}
	}

	return
}

func GetInstanceByName(name string) (interface{}, error) {
	return AppRegistry.GetInstanceByName(name)
}
func (r *InstanceRegistry) GetInstanceByName(name string) (interface{}, error) {
	if td, ok := r.definitionByName[name]; ok {
		v, err := r.GetInstance(td.Type)
		if nil == err {
			return v, nil
		}
	}

	v, ok := r.instanceByName[name]
	if ok {
		return v, nil
	}
	return nil, fmt.Errorf("Instance[%s] is not exist", name)
}

// GetInstances returns instance of annotated
// n must be name of registered annotation
func GetInstances(ts []reflect.Type) ([]interface{}, error) {
	return AppRegistry.GetInstances(ts)
}
func (r *InstanceRegistry) GetInstances(ts []reflect.Type) ([]interface{}, error) {
	var (
		i   interface{}
		err error
	)
	instances := make([]interface{}, 0)
	for _, t := range ts {
		if i, err = r.GetInstance(t); nil != err {
			return nil, err
		}
		instances = append(instances, i)
	}

	return instances, nil
}

func GetInstancesByAnnotationType(at reflect.Type) ([]interface{}, error) {
	return AppRegistry.GetInstancesByAnnotationType(at)
}
func (r *InstanceRegistry) GetInstancesByAnnotationType(at reflect.Type) ([]interface{}, error) {
	var (
		a   la.Annotation
		i   interface{}
		err error
	)
	instances := make([]interface{}, 0)

	for _, td := range r.definitionByType {
		a, err = la.GetTypeAnnotation(td.Type, at)
		if nil != err {
			return nil, err
		}
		if nil != a {
			i, err = r.GetInstance(td.Type)
			if nil != err {
				return nil, err
			}
			instances = append(instances, i)
		}
	}

	return instances, nil
}
