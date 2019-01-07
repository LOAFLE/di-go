package di

import (
	"log"
	"reflect"
	"testing"

	la "github.com/LOAFLE/annotation-go"
	_ "github.com/LOAFLE/di-go/annotation"
)

var InjectableServiceType = reflect.TypeOf((*InjectableService)(nil))

type InjectableService struct {
	la.TypeAnnotation `annotation:"@Injectable('name': 'InjectableService')"`
	Count             int
	Category          string
}

var InjectServiceType = reflect.TypeOf((*InjectService)(nil))

type InjectService struct {
	Service *InjectableService `annotation:"@Inject('name': 'InjectableService')"`

	_Init la.MethodAnnotation `annotation:"@PostConstruct()"`

	R string
}

func (s *InjectService) Init() {
	log.Print("Init")
}

func TestNew(t *testing.T) {
	type args struct {
		parent Registry
	}
	tests := []struct {
		name string
		args args
		want Registry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.parent); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterType(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "InjectableService",
			args: args{
				t: InjectableServiceType,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterType(tt.args.t)
		})
	}
}

func TestInstanceRegistry_RegisterType(t *testing.T) {
	type fields struct {
		parent           Registry
		definitionByType map[reflect.Type]*TypeDefinition
		definitionByName map[string]*TypeDefinition
		instanceByType   map[reflect.Type]interface{}
		instanceByName   map[string]interface{}
	}
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InstanceRegistry{
				parent:           tt.fields.parent,
				definitionByType: tt.fields.definitionByType,
				definitionByName: tt.fields.definitionByName,
				instanceByType:   tt.fields.instanceByType,
				instanceByName:   tt.fields.instanceByName,
			}
			r.RegisterType(tt.args.t)
		})
	}
}

func TestRegisterSingletonByName(t *testing.T) {
	type args struct {
		name     string
		resource interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterSingletonByName(tt.args.name, tt.args.resource); (err != nil) != tt.wantErr {
				t.Errorf("RegisterSingletonByName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceRegistry_RegisterSingletonByName(t *testing.T) {
	type fields struct {
		parent           Registry
		definitionByType map[reflect.Type]*TypeDefinition
		definitionByName map[string]*TypeDefinition
		instanceByType   map[reflect.Type]interface{}
		instanceByName   map[string]interface{}
	}
	type args struct {
		name     string
		resource interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InstanceRegistry{
				parent:           tt.fields.parent,
				definitionByType: tt.fields.definitionByType,
				definitionByName: tt.fields.definitionByName,
				instanceByType:   tt.fields.instanceByType,
				instanceByName:   tt.fields.instanceByName,
			}
			if err := r.RegisterSingletonByName(tt.args.name, tt.args.resource); (err != nil) != tt.wantErr {
				t.Errorf("InstanceRegistry.RegisterSingletonByName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceRegistry_buildDefinition(t *testing.T) {
	type fields struct {
		parent           Registry
		definitionByType map[reflect.Type]*TypeDefinition
		definitionByName map[string]*TypeDefinition
		instanceByType   map[reflect.Type]interface{}
		instanceByName   map[string]interface{}
	}
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TypeDefinition
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InstanceRegistry{
				parent:           tt.fields.parent,
				definitionByType: tt.fields.definitionByType,
				definitionByName: tt.fields.definitionByName,
				instanceByType:   tt.fields.instanceByType,
				instanceByName:   tt.fields.instanceByName,
			}
			got, err := r.buildDefinition(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstanceRegistry.buildDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstanceRegistry.buildDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInstance(t *testing.T) {
	RegisterType(InjectableServiceType)

	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "InjectService",
			args: args{
				t: InjectServiceType,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInstance(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}

func TestInstanceRegistry_GetInstance(t *testing.T) {
	type fields struct {
		parent           Registry
		definitionByType map[reflect.Type]*TypeDefinition
		definitionByName map[string]*TypeDefinition
		instanceByType   map[reflect.Type]interface{}
		instanceByName   map[string]interface{}
	}
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantInstance interface{}
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InstanceRegistry{
				parent:           tt.fields.parent,
				definitionByType: tt.fields.definitionByType,
				definitionByName: tt.fields.definitionByName,
				instanceByType:   tt.fields.instanceByType,
				instanceByName:   tt.fields.instanceByName,
			}
			gotInstance, err := r.GetInstance(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstanceRegistry.GetInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInstance, tt.wantInstance) {
				t.Errorf("InstanceRegistry.GetInstance() = %v, want %v", gotInstance, tt.wantInstance)
			}
		})
	}
}

func TestGetInstanceByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInstanceByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInstanceByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInstanceByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceRegistry_GetInstanceByName(t *testing.T) {
	type fields struct {
		parent           Registry
		definitionByType map[reflect.Type]*TypeDefinition
		definitionByName map[string]*TypeDefinition
		instanceByType   map[reflect.Type]interface{}
		instanceByName   map[string]interface{}
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InstanceRegistry{
				parent:           tt.fields.parent,
				definitionByType: tt.fields.definitionByType,
				definitionByName: tt.fields.definitionByName,
				instanceByType:   tt.fields.instanceByType,
				instanceByName:   tt.fields.instanceByName,
			}
			got, err := r.GetInstanceByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstanceRegistry.GetInstanceByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstanceRegistry.GetInstanceByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInstances(t *testing.T) {
	type args struct {
		ts []reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInstances(tt.args.ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInstances() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceRegistry_GetInstances(t *testing.T) {
	type fields struct {
		parent           Registry
		definitionByType map[reflect.Type]*TypeDefinition
		definitionByName map[string]*TypeDefinition
		instanceByType   map[reflect.Type]interface{}
		instanceByName   map[string]interface{}
	}
	type args struct {
		ts []reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InstanceRegistry{
				parent:           tt.fields.parent,
				definitionByType: tt.fields.definitionByType,
				definitionByName: tt.fields.definitionByName,
				instanceByType:   tt.fields.instanceByType,
				instanceByName:   tt.fields.instanceByName,
			}
			got, err := r.GetInstances(tt.args.ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstanceRegistry.GetInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstanceRegistry.GetInstances() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInstancesByAnnotationType(t *testing.T) {
	type args struct {
		at reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInstancesByAnnotationType(tt.args.at)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInstancesByAnnotationType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInstancesByAnnotationType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceRegistry_GetInstancesByAnnotationType(t *testing.T) {
	type fields struct {
		parent           Registry
		definitionByType map[reflect.Type]*TypeDefinition
		definitionByName map[string]*TypeDefinition
		instanceByType   map[reflect.Type]interface{}
		instanceByName   map[string]interface{}
	}
	type args struct {
		at reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InstanceRegistry{
				parent:           tt.fields.parent,
				definitionByType: tt.fields.definitionByType,
				definitionByName: tt.fields.definitionByName,
				instanceByType:   tt.fields.instanceByType,
				instanceByName:   tt.fields.instanceByName,
			}
			got, err := r.GetInstancesByAnnotationType(tt.args.at)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstanceRegistry.GetInstancesByAnnotationType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstanceRegistry.GetInstancesByAnnotationType() = %v, want %v", got, tt.want)
			}
		})
	}
}
