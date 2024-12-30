package application

//goland:noinspection GoNameStartsWithPackageName
type ApplicationContext struct {
	Name        string
	Controllers []Controller
	Models      []any
	Services    []IService
}
