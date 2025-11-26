package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppConfigSpec struct {
	TargetDeployment string            `json:"targetDeployment"`
	ConfigData       map[string]string `json:"configData"`
}

type AppConfigStatus struct {
	LastApplied time.Time `json:"lastApplied,omitempty"`
}

type AppConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppConfigSpec   `json:"spec,omitempty"`
	Status AppConfigStatus `json:"status,omitempty"`
}

type AppConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppConfig{}, &AppConfigList{})
}
