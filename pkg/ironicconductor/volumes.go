package ironicconductor

import (
	"fmt"
	"strings"

	ironicv1 "github.com/openstack-k8s-operators/ironic-operator/api/v1beta1"
	"github.com/openstack-k8s-operators/ironic-operator/pkg/ironic"
	corev1 "k8s.io/api/core/v1"
)

// GetVolumes -
func GetVolumes(instance *ironicv1.IronicConductor) []corev1.Volume {
	var config0640AccessMode int32 = 0640
	pvcName := fmt.Sprintf("%s-%s", ironic.ServiceName, ironic.ConductorComponent)
	if instance.Spec.ConductorGroup != "" {
		pvcName = strings.ToLower(fmt.Sprintf("%s-%s", pvcName, instance.Spec.ConductorGroup))
	}
	conductorVolumes := []corev1.Volume{
		{
			Name: "config-data-custom",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: &config0640AccessMode,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: fmt.Sprintf("%s-config-data", instance.Name),
					},
				},
			},
		},
		{
			Name: "var-lib-ironic",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvcName,
				},
			},
		},
	}

	return append(ironic.GetVolumes(instance.Name), conductorVolumes...)
}

// GetInitVolumeMounts - Ironic Conductor init task VolumeMounts
func GetInitVolumeMounts() []corev1.VolumeMount {

	initVolumeMounts := []corev1.VolumeMount{
		{
			Name:      "config-data-custom",
			MountPath: "/var/lib/config-data/custom",
			ReadOnly:  true,
		},
		{
			Name:      "var-lib-ironic",
			MountPath: "/var/lib/ironic",
			ReadOnly:  false,
		},
	}

	return append(ironic.GetInitVolumeMounts(), initVolumeMounts...)
}

// GetVolumeMounts - Ironic Conductor VolumeMounts
func GetVolumeMounts(serviceName string) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "config-data-merged",
			MountPath: "/var/lib/kolla/config_files/config.json",
			SubPath:   serviceName + "-config.json",
			ReadOnly:  true,
		},
		{
			Name:      "var-lib-ironic",
			MountPath: "/var/lib/ironic",
			ReadOnly:  false,
		},
	}
	return append(ironic.GetVolumeMounts(), volumeMounts...)
}
