package templates

// TemplateGroupversionInfoGo is the template for the groupversion_info.go files
const TemplateGroupversionInfoGo = `{{ if ne .License "" }}{{ .License }}

{{ end }}// Package {{ .Version}} contains API Schema definitions for the {{ .Group }} {{ .Version }} API group
// +kubebuilder:object:generate=true
// +groupName={{ .Group }}
package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "{{ .Group }}", Version: "{{ .Version }}"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
`
