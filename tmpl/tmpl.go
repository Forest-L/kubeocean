package tmpl

import (
	"fmt"
	"text/template"
)

func getTmpl(tmpl string) *template.Template {
	if len(tmpl) == 0 {
		fmt.Printf("There is no template named %s", tmpl)
	} else if tmpl == "kubelet" {
		return kubeletContainerTempl
	} else if tmpl == "kubeletService" {
		return kubeletServiceTempl
	}
	return nil
}
