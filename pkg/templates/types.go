package templates

// TemplateTypesGo is the template for the _types.go files
const TemplateTypesGo = `/*
Copyright 2020 Datadrivers GmbH.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package {{ .Version }}

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// {{ .Names.Kind }} is the Schema for the {{ .Names.Plural }} API
type {{ .Names.Kind }} struct {
	metav1.TypeMeta   ` + "`" + `json:",inline"` + "`" + `
	metav1.ObjectMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `

	Spec   {{ .Names.Kind }}Spec   ` + "`" + `json:"spec,omitempty"` + "`" + `
	Status {{ .Names.Kind }}Status ` + "`" + `json:"status,omitempty"` + "`" + `
}

// +kubebuilder:object:root=true

// {{ .Names.Kind }}List contains a list of {{ .Names.Kind }}
type {{ .Names.Kind }}List struct {
	metav1.TypeMeta ` + "`" + `json:",inline"` + "`" + `
	metav1.ListMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `
	Items           []{{ .Names.Kind }} ` + "`" + `json:"items"` + "`" + `
}

func init() {
	SchemeBuilder.Register(&{{ .Names.Kind }}{}, &{{ .Names.Kind }}List{})
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
