package controller

import (
	"crypto/md5"
	"fmt"
	v1 "github.com/yxwuxuanl/k8s-nginx-operator/api/v1"
	"github.com/yxwuxuanl/k8s-nginx-operator/internal/nginx"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"maps"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	NginxConfigKey = "nginx.conf"
	ConfigSumLabel = "nginx-operator/configsum"
)

type NgxObject interface {
	client.Object
	GetNginxSpec() v1.NginxSpec
	GetNginxConfig() nginx.Config
	GetNamePrefix() string
}

type ServiceProtocol interface {
	GetServiceProtocol() corev1.Protocol
}

func buildDeployment(ngx NgxObject) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getObjectName(ngx),
			Namespace: ngx.GetNamespace(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To(ngx.GetNginxSpec().Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: make(map[string]string),
			},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					EnableServiceLinks:           ptr.To(false),
					AutomountServiceAccountToken: ptr.To(false),
					NodeSelector:                 ngx.GetNginxSpec().NodeSelector,
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: getObjectName(ngx),
									},
								},
							},
						},
						{
							Name: "tmp",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: ngx.GetNginxSpec().Image,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: nginx.ConfigDirPath,
									SubPath:   NginxConfigKey,
								},
								{
									Name:      "tmp",
									MountPath: "/tmp",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "nginx",
									ContainerPort: nginx.ListenPort,
									Protocol:      getProtocol(ngx),
								},
								{
									Name:          "probe",
									ContainerPort: nginx.ProbePort,
								},
							},
							Env: ngx.GetNginxSpec().Env,
							SecurityContext: &corev1.SecurityContext{
								RunAsUser:                ptr.To[int64](nginx.UserID),
								RunAsGroup:               ptr.To[int64](nginx.GroupID),
								RunAsNonRoot:             ptr.To(true),
								Privileged:               ptr.To(false),
								AllowPrivilegeEscalation: ptr.To(false),
								ReadOnlyRootFilesystem:   ptr.To(true),
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{
										"ALL",
									},
								},
							},
							Command:   []string{"nginx"},
							Resources: ptr.Deref(ngx.GetNginxSpec().Resources, corev1.ResourceRequirements{}),
							LivenessProbe: &corev1.Probe{
								InitialDelaySeconds: 5,
								SuccessThreshold:    1,
								FailureThreshold:    10,
								TimeoutSeconds:      5,
								PeriodSeconds:       5,
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: nginx.ProbeURL,
										Port: intstr.Parse("probe"),
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								SuccessThreshold: 1,
								FailureThreshold: 10,
								TimeoutSeconds:   5,
								PeriodSeconds:    5,
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: nginx.ProbeURL,
										Port: intstr.Parse("probe"),
									},
								},
							},
						},
					},
					Affinity: ngx.GetNginxSpec().Affinity,
				},
			},
		},
	}

	setMetadata(&deployment.Spec.Template, ngx)
	maps.Copy(
		deployment.Spec.Selector.MatchLabels,
		deployment.Spec.Template.Labels,
	)

	if labels := ngx.GetNginxSpec().PodLabels; labels != nil {
		maps.Copy(labels, deployment.Spec.Template.Labels)
		deployment.Spec.Template.Labels = labels
	}

	return deployment
}

func buildConfigMap(ngx NgxObject) (*corev1.ConfigMap, string, error) {
	config, err := nginx.BuildConfig(ngx.GetNginxConfig())
	if err != nil {
		return nil, "", fmt.Errorf("build nginx config error: %w", err)
	}

	confsum := fmt.Sprintf("%x", md5.Sum(config))

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ngx.GetNamespace(),
			Name:      getObjectName(ngx),
			Labels: map[string]string{
				ConfigSumLabel: confsum,
			},
		},
		Data: map[string]string{
			NginxConfigKey: string(config),
		},
	}

	return cm, confsum, nil
}

func buildService(ngx NgxObject, selector map[string]string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getObjectName(ngx),
			Namespace: ngx.GetNamespace(),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "nginx",
					Port:       ngx.GetNginxSpec().ServicePort,
					TargetPort: intstr.Parse("nginx"),
					Protocol:   getProtocol(ngx),
				},
			},
			Selector: selector,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
}

func buildObjects(ngx NgxObject) (objects []client.Object, err error) {
	configMap, confsum, err := buildConfigMap(ngx)
	if err != nil {
		return nil, fmt.Errorf("failed to build configmap: %w", err)
	}

	deployment := buildDeployment(ngx)
	deployment.Spec.Template.Labels[ConfigSumLabel] = confsum

	service := buildService(ngx, deployment.Spec.Selector.MatchLabels)

	defer func() {
		for _, object := range objects {
			setMetadata(object, ngx)
		}
	}()

	return []client.Object{
		configMap, deployment, service,
	}, nil
}

func setMetadata(o metav1.Object, ngx NgxObject) {
	labels := o.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}

	labels["app.kubernetes.io/managed-by"] = "nginx-operator"
	labels["app.kubernetes.io/instance"] = getObjectName(ngx)

	o.SetLabels(labels)
}

func getObjectName(ngx NgxObject) string {
	return ngx.GetNamePrefix() + ngx.GetName()
}

func getProtocol(ngx NgxObject) corev1.Protocol {
	if v, ok := ngx.(ServiceProtocol); ok {
		return v.GetServiceProtocol()
	}

	return corev1.ProtocolTCP
}
