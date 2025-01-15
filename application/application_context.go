package application

import (
	"errors"
	"reflect"
)

//goland:noinspection GoNameStartsWithPackageName
type ApplicationContext struct {
	Name        string
	Controllers []IController
	Services    []IService
}

func GetServiceFromContext[T IService](ctx *ApplicationContext) (*T, error) {
	for _, service := range ctx.Services {
		// Check if the type matches T
		if serviceFound, ok := service.(T); ok {
			return &serviceFound, nil
		}
	}
	return nil, errors.New("service not found")
}

func (ctx *ApplicationContext) GetService(serviceType IService) (IService, error) {
	for _, service := range ctx.Services {
		if reflect.TypeOf(service) == reflect.TypeOf(serviceType) {
			return service, nil
		}
	}
	return nil, errors.New("service not found")
}
