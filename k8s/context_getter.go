package k8s

import "k8s.io/client-go/tools/clientcmd"

type contextGetter struct {
	configAccess clientcmd.ConfigAccess
}

func NewContextGetter() *contextGetter {
	return &contextGetter{
		configAccess: clientcmd.NewDefaultPathOptions(),
	}
}

func (c contextGetter) GetAll() ([]string, error) {
	config, err := c.configAccess.GetStartingConfig()
	if err != nil {
		return nil, err
	}
	allContextNames := []string{}
	for name := range config.Contexts {
		allContextNames = append(allContextNames, name)
	}

	return allContextNames, nil
}
