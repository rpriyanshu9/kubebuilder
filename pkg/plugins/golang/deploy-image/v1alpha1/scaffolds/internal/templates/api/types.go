/*
Copyright 2022 The Kubernetes Authors.

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

package api

import (
	"fmt"
	"path/filepath"

	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
)

var _ machinery.Template = &Types{}

// Types scaffolds the file that defines the schema for a CRD
// nolint:maligned
type Types struct {
	machinery.TemplateMixin
	machinery.MultiGroupMixin
	machinery.BoilerplateMixin
	machinery.ResourceMixin

	// Port if informed we will create the scaffold with this spec
	Port string
}

// SetTemplateDefaults implements file.Template
func (f *Types) SetTemplateDefaults() error {
	if f.Path == "" {
		if f.MultiGroup {
			if f.Resource.Group != "" {
				f.Path = filepath.Join("apis", "%[group]", "%[version]", "%[kind]_types.go")
			} else {
				f.Path = filepath.Join("apis", "%[version]", "%[kind]_types.go")
			}
		} else {
			f.Path = filepath.Join("api", "%[version]", "%[kind]_types.go")
		}
	}
	f.Path = f.Resource.Replacer().Replace(f.Path)
	fmt.Println(f.Path)

	f.TemplateBody = typesTemplate

	f.IfExistsAction = machinery.OverwriteFile

	return nil
}

const typesTemplate = `{{ .Boilerplate }}

package {{ .Resource.Version }}

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// {{ .Resource.Kind }}Spec defines the desired state of {{ .Resource.Kind }}
type {{ .Resource.Kind }}Spec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Size defines the number of {{ .Resource.Kind }} instances
	Size int32 ` + "`" + `json:"size,omitempty"` + "`" + `

	{{ if not (isEmptyStr .Port) -}}
	// Port defines the port that will be used to init the container with the image
	ContainerPort int32 ` + "`" + `json:"containerPort,omitempty"` + "`" + `
	{{- end }}
}

// {{ .Resource.Kind }}Status defines the observed state of {{ .Resource.Kind }}
type {{ .Resource.Kind }}Status struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
{{- if and (not .Resource.API.Namespaced) (not .Resource.IsRegularPlural) }}
//+kubebuilder:resource:path={{ .Resource.Plural }},scope=Cluster
{{- else if not .Resource.API.Namespaced }}
//+kubebuilder:resource:scope=Cluster
{{- else if not .Resource.IsRegularPlural }}
//+kubebuilder:resource:path={{ .Resource.Plural }}
{{- end }}

// {{ .Resource.Kind }} is the Schema for the {{ .Resource.Plural }} API
type {{ .Resource.Kind }} struct {
	metav1.TypeMeta   ` + "`" + `json:",inline"` + "`" + `
	metav1.ObjectMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `

	Spec   {{ .Resource.Kind }}Spec   ` + "`" + `json:"spec,omitempty"` + "`" + `
	Status {{ .Resource.Kind }}Status ` + "`" + `json:"status,omitempty"` + "`" + `
}

//+kubebuilder:object:root=true

// {{ .Resource.Kind }}List contains a list of {{ .Resource.Kind }}
type {{ .Resource.Kind }}List struct {
	metav1.TypeMeta ` + "`" + `json:",inline"` + "`" + `
	metav1.ListMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `
	Items           []{{ .Resource.Kind }} ` + "`" + `json:"items"` + "`" + `
}

func init() {
	SchemeBuilder.Register(&{{ .Resource.Kind }}{}, &{{ .Resource.Kind }}List{})
}
`
