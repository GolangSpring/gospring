package repository

type HostGroup struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Hosts []string `json:"hosts"`
}

type Inventory struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	HostGroups map[string][]string `json:"host_groups"`
}

func (inventory *Inventory) AsJson() map[string]map[string]any {
	converted := make(map[string]map[string]any)

	// Convert each host group into a structure Ansible expects
	for groupName, hosts := range inventory.HostGroups {
		hostEntries := make(map[string]any)
		for _, host := range hosts {
			hostEntries[host] = map[string]any{
				"ansible_connection": "local",
			}
		}
		converted[groupName] = map[string]any{
			"hosts": hostEntries,
		}
	}

	return converted
}

type IInventoryService interface {
	GetInventory(id string) (*Inventory, error)
	CreateInventory(inventory *Inventory) (*Inventory, error)
	GetAllInventory() ([]*Inventory, error)
	UpdateInventory(id string, inventory *Inventory) (*Inventory, error)
	DeleteInventory(id string) (any, error)
}
