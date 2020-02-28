package builder

import "kpackui/kpack"

type CustomClusterBuilderRepo interface {
	GetAllCustomClusterBuilders() ([]kpack.CustomClusterBuilder, error)
}

func NewCustomClusterGetter(repo CustomClusterBuilderRepo) *customClusterGetter {
	return &customClusterGetter{
		repo: repo,
	}
}

type customClusterGetter struct {
	repo CustomClusterBuilderRepo
}

func (c *customClusterGetter) GetAll() ([]kpack.CustomClusterBuilder, error) {
	return c.repo.GetAllCustomClusterBuilders()
}
