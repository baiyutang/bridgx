package ecloud

import "github.com/galaxy-future/BridgX/pkg/cloud"

func (p *ECloud) ContainerInstanceList(region string, pageNumber, pageSize int) ([]cloud.RegistryInstance, int, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) EnterpriseNamespaceList(region, instanceId string, pageNumber, pageSize int) ([]cloud.Namespace, int, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) PersonalNamespaceList(region string) ([]cloud.Namespace, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) EnterpriseRepositoryList(region, instanceId, namespace string, pageNumber, pageSize int) ([]cloud.Repository, int, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) PersonalRepositoryList(region, namespace string, pageNumber, pageSize int) ([]cloud.Repository, int, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) EnterpriseImageList(region, instanceId, repoId, namespace, repoName string, pageNumber, pageSize int) ([]cloud.DockerArtifact, int, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) PersonalImageList(region, repoNamespace, repoName string, pageNum, pageSize int) ([]cloud.DockerArtifact, int, error) {
	// TODO implement me
	panic("implement me")
}
