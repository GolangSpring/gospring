package application

import "reflect"

//goland:noinspection GoNameStartsWithPackageName
type ApplicationContext struct {
	Name        string
	Controllers []IController
	Models      []any
	Services    []IService
}

func (ctx *ApplicationContext) GetService(serviceType IService) IService {
	for _, service := range ctx.Services {
		if reflect.TypeOf(service) == reflect.TypeOf(serviceType) {
			return service
		}
	}
	return nil
}
