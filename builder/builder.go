package builder

import "kpackui/kpack"

type CustomClusterBuilderRepo interface {
	GetAllCustomClusterBuilders() ([]kpack.CustomClusterBuilder, error)
}

func NewCustomClusterGetter(repo CustomClusterBuilderRepo) *CustomClusterGetter {
	return &CustomClusterGetter{
		repo: repo,
	}
}

type CustomClusterGetter struct {
	repo CustomClusterBuilderRepo
}

func (c *CustomClusterGetter) GetAll() ([]kpack.CustomClusterBuilder, error) {
	return c.repo.GetAllCustomClusterBuilders()
}
