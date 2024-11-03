package funcs

import (
	"context"
	"fmt"
	"github.com/KokoiRuby/k8s-copilot/cmd/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/restmapper"
)

func CreateResource(ctx context.Context, client *utils.OpenAI, input, kubeConfigPath string) (string, error) {
	sysPrompt := `
You're a K8s resource YAML manifest generator.
Please generate corresponding YAML manifest based on user input.
Don't include it into YAML code block.
`

	// generate YAML manifest given user input
	yaml, err := client.SendMessage(sysPrompt, input)
	if err != nil {
		return "", err
	}
	//return yamlContent, nil

	// client-go
	clientGo, err := utils.NewClientGo(kubeConfigPath)
	if err != nil {
		return "", err
	}

	// kubectl api-resources → REST mapper
	res, err := restmapper.GetAPIGroupResources(clientGo.DiscoveryClient)
	if err != nil {
		return "", err
	}
	mapper := restmapper.NewDiscoveryRESTMapper(res)

	// yaml to unstructured
	unstructuredObj := &unstructured.Unstructured{}
	_, _, err = scheme.Codecs.UniversalDeserializer().Decode([]byte(yaml), nil, unstructuredObj)
	if err != nil {
		return "", err
	}

	// get gvr from gvk of unstructured
	gvk := unstructuredObj.GroupVersionKind()
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return "", err
	}

	namespace := unstructuredObj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	// create unstructured gvr
	_, err = clientGo.DynamicClient.Resource(mapping.Resource).Namespace(namespace).Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Resource %s created successfully", unstructuredObj.GetName()), nil

}
