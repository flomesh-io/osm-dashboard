package meshconfig

import (
	"context"
	"log"

	"github.com/kubernetes/dashboard/src/app/backend/api"
	osmconfigv1alph2 "github.com/openservicemesh/osm/pkg/apis/config/v1alpha2"
	osmconfigclientset "github.com/openservicemesh/osm/pkg/gen/client/config/clientset/versioned"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/kubernetes"
)

// MeshConfigDetail API resource provides mechanisms to inject containers with configuration data while keeping
// containers agnostic of Kubernetes
type MeshConfigDetail struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`
	// Spec is the MeshConfig specification.
	// +optional
	Spec     osmconfigv1alph2.MeshConfigSpec `json:"spec,omitempty"`
	MeshName string                          `json:"meshName"`
	Option   string                          `json:"option"`
}

// GetMeshConfigDetail returns detailed information about an meshconfig
func GetMeshConfigDetail(osmConfigClient osmconfigclientset.Interface, client client.Interface, namespace, name string) (*MeshConfigDetail, error) {
	log.Printf("Getting details of %s meshconfig in %s namespace", name, namespace)

	rawMeshConfig, err := osmConfigClient.ConfigV1alpha2().MeshConfigs(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return getMeshConfigDetail(rawMeshConfig, client), nil
}

func getMeshConfigDetail(meshConfig *osmconfigv1alph2.MeshConfig, client client.Interface) *MeshConfigDetail {
	meshName := ""
	deployment, err := client.AppsV1().Deployments(meshConfig.ObjectMeta.Namespace).Get(context.TODO(), "osm-controller", metaV1.GetOptions{})
	if err == nil {
		meshName = deployment.ObjectMeta.Labels["meshName"]
	}

	return &MeshConfigDetail{
		ObjectMeta: api.NewObjectMeta(meshConfig.ObjectMeta),
		TypeMeta:   api.NewTypeMeta(api.ResourceKindMeshConfig),
		Spec:       meshConfig.Spec,
		MeshName:   meshName,
	}
}
