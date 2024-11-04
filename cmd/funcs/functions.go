package funcs

import (
	"context"
	"fmt"
	"github.com/KokoiRuby/k8s-copilot/cmd/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/restmapper"
	"strings"
)

func CreateResource(ctx context.Context, client *utils.OpenAI, input, kubeConfig string) (string, error) {
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
	clientGo, err := utils.NewClientGo(kubeConfig)
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
	_, err = clientGo.DynamicClient.Resource(mapping.Resource).Namespace(namespace).Create(ctx, unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Resource %s created successfully", unstructuredObj.GetName()), nil

}

func ListResource(ctx context.Context, namespace, resource, kubeConfig string) (string, error) {
	// client-go
	clientGo, err := utils.NewClientGo(kubeConfig)
	if err != nil {
		return "", err
	}

	var isNonNamespaced bool

	// build gvr based on resource
	// TODO: it gets heavy if we populate all resources, any improvements?
	var gvr schema.GroupVersionResource
	switch strings.ToLower(resource) {
	case "pods", "pod", "po":
		gvr = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	case "services", "service", "svc":
		gvr = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	case "deployments", "deployment":
		gvr = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	case "namespaces", "namespace", "ns":
		gvr = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}
		isNonNamespaced = true
	default:
		return "", fmt.Errorf("resource %s not supported", resource)
	}

	// list gvr
	var resList *unstructured.UnstructuredList
	if isNonNamespaced {
		resList, err = clientGo.DynamicClient.Resource(gvr).List(ctx, metav1.ListOptions{})
		if err != nil {
			return "", err
		}
	} else {
		resList, err = clientGo.DynamicClient.Resource(gvr).Namespace(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return "", err
		}
	}
	result := ""
	for _, item := range resList.Items {
		result += fmt.Sprintf("%s\n", item.GetName())
	}

	return result, nil
}
