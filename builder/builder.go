package builder

import "kpackui/kpack"

type ClusterBuilderRepo interface {
	GetAllCustomClusterBuilders() ([]kpack.ClusterBuilder, error)
	GetAllClusterBuilders() ([]kpack.ClusterBuilder, error)
}

func NewCustomClusterGetter(repo ClusterBuilderRepo) *CustomClusterGetter {
	return &CustomClusterGetter{
		repo: repo,
	}
}

type CustomClusterGetter struct {
	repo ClusterBuilderRepo
}

func (c *CustomClusterGetter) GetAll() ([]kpack.ClusterBuilder, error) {
	return c.repo.GetAllCustomClusterBuilders()
}

func NewClusterGetter(repo ClusterBuilderRepo) *ClusterGetter {
	return &ClusterGetter{
		repo: repo,
	}
}

type ClusterGetter struct {
	repo ClusterBuilderRepo
}

func (c *ClusterGetter) GetAll() ([]kpack.ClusterBuilder, error) {
	return c.repo.GetAllCustomClusterBuilders()
}
