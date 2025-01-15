package service

import (
	"github.com/casbin/casbin/v2"
)

const ModelString = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[role_definition]
g = _, _

[matchers]
m = r.sub == p.sub && regexMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`

const CasbinPublicKey = "public"

type CasbinService struct {
	Enforcer *casbin.Enforcer
}

func NewCasbinService(enforcer *casbin.Enforcer) *CasbinService {
	return &CasbinService{Enforcer: enforcer}
}

func (service *CasbinService) PostConstruct() {}

func (service *CasbinService) HasPermission(subject, object, action string) (bool, error) {
	return service.Enforcer.Enforce(subject, object, action)
}
