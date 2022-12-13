package modules

type Module struct {
	Namespace string
	System    string
	Name      string
	Version   string
}

type Modules struct {
	Modules []Module `json:"modules"`
}

func (modules *Modules) exists(rhs Module) bool {
	for _, lhs := range modules.Modules {
		if lhs.Namespace == rhs.Namespace && lhs.System == rhs.Namespace && lhs.Name == rhs.Namespace && lhs.Version == rhs.Version {
			return true
		}
	}
	return false
}

func (modules *Modules) Add(m Module) {
	if !modules.exists(m) {
		modules.Modules = append(modules.Modules, m)
	}
}

func (module *Module) Validate() error {
	return nil
}
