package k8s

import "k8s.io/client-go/tools/clientcmd"

type ContextGetter struct {
	configAccess clientcmd.ConfigAccess
}

func NewContextGetter() *ContextGetter {
	return &ContextGetter{
		configAccess: clientcmd.NewDefaultPathOptions(),
	}
}

func (c ContextGetter) GetAll() ([]string, error) {
	config, err := c.configAccess.GetStartingConfig()
	if err != nil {
		return nil, err
	}
	var allContextNames []string
	for name := range config.Contexts {
		allContextNames = append(allContextNames, name)
	}

	return allContextNames, nil
}
