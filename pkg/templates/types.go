package templates

// TemplateTypesGo is the template for the _types.go files
const TemplateTypesGo = `{{ if ne .License "" }}{{ .License }}

{{ end }}package {{ .CRD.Version }}

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// {{ .CRD.Names.Kind }} is the Schema for the {{ .CRD.Names.Plural }} API
type {{ .CRD.Names.Kind }} struct {
	metav1.TypeMeta   ` + "`" + `json:",inline"` + "`" + `
	metav1.ObjectMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `

	Spec   {{ .CRD.Names.Kind }}Spec   ` + "`" + `json:"spec,omitempty"` + "`" + `
	Status {{ .CRD.Names.Kind }}Status ` + "`" + `json:"status,omitempty"` + "`" + `
}

// +kubebuilder:object:root=true

// {{ .CRD.Names.Kind }}List contains a list of {{ .CRD.Names.Kind }}
type {{ .CRD.Names.Kind }}List struct {
	metav1.TypeMeta ` + "`" + `json:",inline"` + "`" + `
	metav1.ListMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `
	Items           []{{ .CRD.Names.Kind }} ` + "`" + `json:"items"` + "`" + `
}

func init() {
	SchemeBuilder.Register(&{{ .CRD.Names.Kind }}{}, &{{ .CRD.Names.Kind }}List{})
}

`

// TemplateSpecBlock is the template for the spec block of the _types.go files
const TemplateSpecBlock = `// {{ .Kind }}Spec defines the desired state of {{ .Kind }}
type {{ .Kind }}Spec struct { {{ range .Properties }}
	{{ .Name }}	{{ .Type }}	{{ .JSON }}{{ end }}
}
`

// TemplateStatusBlock is the template for the status block of the _types.go files
const TemplateStatusBlock = `// {{ .Kind }}Status defines the observed state of {{ .Kind }}
type {{ .Kind }}Status struct { {{ range .Properties }}
	{{ .Name }}	{{ .Type }}	{{ .JSON }}{{ end }}
}
`

// TemplateBlock is the template for the _types.go files
const TemplateBlock = `// {{ .Kind }} : {{ .Description }}
type {{ .Kind }} struct { {{ range .Properties }}
	{{ .Name }}	{{ .Type }}	{{ .JSON }}{{ end }}
}
`
