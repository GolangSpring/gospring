package service

import (
	"github.com/casbin/casbin/v2"
)

const (
	ModelString = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[role_definition]
g = _, _

[matchers]
m = g(r.sub, p.sub) && regexMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`
	DefaultPublic = "public"
)

const CasbinPublicKey = "public"

type CasbinService struct {
	Enforcer *casbin.Enforcer
}

func NewCasbinService(enforcer *casbin.Enforcer) *CasbinService {
	service := &CasbinService{Enforcer: enforcer}
	return service
}

func (service *CasbinService) PostConstruct() {}

func (service *CasbinService) HasPermission(subject, object, action string) (bool, error) {
	return service.Enforcer.Enforce(subject, object, action)
}

func (service *CasbinService) GetAllUsedRoles() ([]string, error) {
	roles := []string{}
	policyRoles, err := service.Enforcer.GetAllSubjects()
	if err != nil {
		return roles, err
	}

	for _, role := range policyRoles {
		if role == DefaultPublic {
			continue
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (service *CasbinService) AddRoleForUser(user string, role string) (bool, error) {
	return service.Enforcer.AddRoleForUser(user, role)
}
