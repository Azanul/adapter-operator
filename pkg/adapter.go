package meshsync

import (
	mesheryv1alpha1 "github.com/Azanul/adapter-operator/api/v1alpha1"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

const (
	ServerObject = "server-object"
)

type Object interface {
	runtime.Object
	metav1.Object
}

func GetObjects(m *mesheryv1alpha1.Adapter) map[string]Object {
	return map[string]Object{
		ServerObject: getServerObject(m.ObjectMeta.Namespace, m.ObjectMeta.Name, m.Spec.Image, m.Spec.Command, m.Spec.HostPort, m.Spec.ContainerPort),
	}
}

func getServerObject(namespace, name, image string, command []string, hostPort, containerPort int) Object {
	var obj = &v1.Deployment{}
	Deployment.DeepCopyInto(obj)
	obj.ObjectMeta.Namespace = namespace
	obj.ObjectMeta.Name = name
	obj.Spec.Template.Spec.Containers[0].Name = name
	obj.Spec.Template.Spec.Containers[0].Image = image
	obj.Spec.Template.Spec.Containers[0].Command = command
	obj.Spec.Template.Spec.Containers[0].Ports[0].HostPort = int32(hostPort)
	obj.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = int32(containerPort)
	return obj
}
