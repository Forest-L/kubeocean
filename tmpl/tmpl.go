package tmpl

import (
	"fmt"
	"text/template"
)

func GetTmpl(tmpl string) *template.Template {
	if len(tmpl) == 0 {
		fmt.Printf("There is no template named %s", tmpl)
	} else if tmpl == "kubelet" {
		return kubeletTempl
	} else if tmpl == "kubeletService" {
		return kubeletServiceTempl
	} else if tmpl == "kubeletContainer" {
		return kubeletContainerTempl
	}
	return nil
}
