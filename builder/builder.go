package builder

import "kpackui/kpack"

type ClusterBuilderRepo interface {
	GetAllCustomClusterBuilders() ([]*kpack.ClusterBuilder, error)
	GetAllClusterBuilders() ([]*kpack.ClusterBuilder, error)
}

type NamespacedBuilderRepo interface {
	GetAllCustomBuilders(namespace string) ([]*kpack.NamespacedBuilder, error)
	GetAllNamespacedBuilders(namespace string) ([]*kpack.NamespacedBuilder, error)
}

func NewCustomClusterGetter(repo ClusterBuilderRepo) *CustomClusterGetter {
	return &CustomClusterGetter{
		repo: repo,
	}
}

type CustomClusterGetter struct {
	repo ClusterBuilderRepo
}

func (c *CustomClusterGetter) GetAll() ([]*kpack.ClusterBuilder, error) {
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

func (c *ClusterGetter) GetAll() ([]*kpack.ClusterBuilder, error) {
	return c.repo.GetAllClusterBuilders()
}

func NewNamespacedGetter(repo NamespacedBuilderRepo) *NamespaceGetter {
	return &NamespaceGetter{
		repo: repo,
	}
}

type NamespaceGetter struct {
	repo NamespacedBuilderRepo
}

func (c *NamespaceGetter) GetAll(namespace string) ([]*kpack.NamespacedBuilder, error) {
	return c.repo.GetAllNamespacedBuilders(namespace)
}

func NewCustomNamespacedGetter(repo NamespacedBuilderRepo) *NamespaceGetter {
	return &NamespaceGetter{
		repo: repo,
	}
}

type CustomNamespaceGetter struct {
	repo NamespacedBuilderRepo
}

func (c *CustomNamespaceGetter) GetAll(namespace string) ([]*kpack.NamespacedBuilder, error) {
	return c.repo.GetAllCustomBuilders(namespace)
}
