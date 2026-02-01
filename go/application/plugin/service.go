package plugin

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
)

// Service is embedded in service structs to provide metadata.
// Tag format: `name:"service_name" desc:"Service description"`
type Service struct{}

// Op is a field type for declaring operations.
// Tag format: `desc:"Operation description"`
type Op struct{}

// Request contains the context for a handler invocation.
type Request struct {
	Client interface{} // Plugin-specific client (e.g., *AWSClient)
	Config interface{} // Parsed config struct
	Raw    []byte      // Raw config JSON
}

// HandlerFunc is the signature for operation handlers.
type HandlerFunc func(ctx context.Context, req *Request) (*entities.Result, error)

// MustRegisterService registers a service or panics.
// Use this in init() functions.
func MustRegisterService(plugin *PluginDefinition, svc interface{}) {
	if err := RegisterService(plugin, svc); err != nil {
		panic(fmt.Sprintf("failed to register service: %v", err))
	}
}

// RegisterService registers all operations from a service struct.
func RegisterService(plugin *PluginDefinition, svc interface{}) error {
	svcType := reflect.TypeOf(svc)
	svcValue := reflect.ValueOf(svc)

	// Must be a pointer to struct
	if svcType.Kind() != reflect.Ptr || svcType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("service must be a pointer to struct, got %T", svc)
	}

	structType := svcType.Elem()

	// Find embedded Service field and extract metadata
	serviceName, serviceDesc, err := extractServiceMetadata(structType)
	if err != nil {
		return err
	}

	// Find all Op fields and their descriptions
	ops, err := extractOperations(structType)
	if err != nil {
		return err
	}

	// Match operations to methods and register
	for _, op := range ops {
		method := svcValue.MethodByName(op.methodName)
		if !method.IsValid() {
			return fmt.Errorf("service %s: no method %s for operation %s (field %s)",
				serviceName, op.methodName, op.name, op.fieldName)
		}

		handler, err := wrapMethod(method)
		if err != nil {
			return fmt.Errorf("service %s, operation %s: %w",
				serviceName, op.name, err)
		}

		plugin.RegisterHandler(serviceName, serviceDesc, op.name, op.description, handler)
	}

	return nil
}

// extractServiceMetadata finds the embedded Service field and parses its tags.
func extractServiceMetadata(t reflect.Type) (name, desc string, err error) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type == reflect.TypeOf(Service{}) {
			tag := field.Tag
			name = tag.Get("name")
			desc = tag.Get("desc")
			if name == "" {
				return "", "", fmt.Errorf("Service field missing 'name' tag")
			}
			return name, desc, nil
		}
	}
	return "", "", fmt.Errorf("struct must embed plugin.Service")
}

// opInfo holds operation metadata extracted from struct fields.
type opInfo struct {
	fieldName   string // PascalCase field name
	methodName  string // Method name to invoke
	name        string // snake_case operation name
	description string
}

// extractOperations finds all Op fields and extracts their metadata.
func extractOperations(t reflect.Type) ([]opInfo, error) {
	var ops []opInfo

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type == reflect.TypeOf(Op{}) {
			methodName := field.Tag.Get("method")
			if methodName == "" {
				// Default to field name if no method tag, but this implies name collision
				// if method is defined on the same struct (which is invalid in Go).
				// We assume if no tag, the user might have defined the method on the pointer
				// and the field on the struct, but Go forbids name collision even then.
				// So we really expect the tag for valid Go code.
				methodName = field.Name
			}

			ops = append(ops, opInfo{
				fieldName:   field.Name,
				methodName:  methodName,
				name:        toSnakeCase(field.Name),
				description: field.Tag.Get("desc"),
			})
		}
	}

	if len(ops) == 0 {
		return nil, fmt.Errorf("service has no operations (no Op fields)")
	}

	return ops, nil
}

// wrapMethod wraps a reflected method as a HandlerFunc.
func wrapMethod(method reflect.Value) (HandlerFunc, error) {
	methodType := method.Type()

	// Expected signature: func(ctx context.Context, req *Request) (*entities.Result, error)
	if methodType.NumIn() != 2 || methodType.NumOut() != 2 {
		return nil, fmt.Errorf("method must have signature (context.Context, *Request) (*entities.Result, error)")
	}

	// Validate input types
	ctxType := reflect.TypeOf((*context.Context)(nil)).Elem()
	reqType := reflect.TypeOf((*Request)(nil))
	if !methodType.In(0).Implements(ctxType) {
		return nil, fmt.Errorf("first parameter must be context.Context")
	}
	if methodType.In(1) != reqType {
		return nil, fmt.Errorf("second parameter must be *plugin.Request")
	}

	// Validate output types
	resultType := reflect.TypeOf((*entities.Result)(nil))
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if methodType.Out(0) != resultType {
		return nil, fmt.Errorf("first return value must be *entities.Result")
	}
	if !methodType.Out(1).Implements(errorType) {
		return nil, fmt.Errorf("second return value must be error")
	}

	return func(ctx context.Context, req *Request) (*entities.Result, error) {
		args := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(req),
		}
		results := method.Call(args)

		var result *entities.Result
		if !results[0].IsNil() {
			result = results[0].Interface().(*entities.Result)
		}

		var err error
		if !results[1].IsNil() {
			err = results[1].Interface().(error)
		}

		return result, err
	}, nil
}

// toSnakeCase converts PascalCase to snake_case.
var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
