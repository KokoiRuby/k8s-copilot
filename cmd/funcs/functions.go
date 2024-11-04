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
)

type Resource struct {
	GVR        schema.GroupVersionResource
	Namespaced bool
}

var resourceMap = map[string]Resource{
	"pods":                              {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}, Namespaced: true},
	"services":                          {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}, Namespaced: true},
	"deployments":                       {GVR: schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}, Namespaced: true},
	"namespaces":                        {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}, Namespaced: false},
	"bindings":                          {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "bindings"}, Namespaced: true},
	"componentstatuses":                 {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "componentstatuses"}, Namespaced: false},
	"configmaps":                        {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}, Namespaced: true},
	"endpoints":                         {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "endpoints"}, Namespaced: true},
	"events":                            {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "events"}, Namespaced: true},
	"limitranges":                       {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "limitranges"}, Namespaced: true},
	"nodes":                             {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "nodes"}, Namespaced: false},
	"persistentvolumeclaims":            {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "persistentvolumeclaims"}, Namespaced: true},
	"persistentvolumes":                 {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "persistentvolumes"}, Namespaced: false},
	"podtemplates":                      {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "podtemplates"}, Namespaced: true},
	"replicationcontrollers":            {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "replicationcontrollers"}, Namespaced: true},
	"resourcequotas":                    {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "resourcequotas"}, Namespaced: true},
	"secrets":                           {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}, Namespaced: true},
	"serviceaccounts":                   {GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "serviceaccounts"}, Namespaced: true},
	"challenges":                        {GVR: schema.GroupVersionResource{Group: "acme.cert-manager.io", Version: "v1", Resource: "challenges"}, Namespaced: true},
	"orders":                            {GVR: schema.GroupVersionResource{Group: "acme.cert-manager.io", Version: "v1", Resource: "orders"}, Namespaced: true},
	"mutatingwebhookconfigurations":     {GVR: schema.GroupVersionResource{Group: "admissionregistration.k8s.io", Version: "v1", Resource: "mutatingwebhookconfigurations"}, Namespaced: false},
	"validatingadmissionpolicies":       {GVR: schema.GroupVersionResource{Group: "admissionregistration.k8s.io", Version: "v1", Resource: "validatingadmissionpolicies"}, Namespaced: false},
	"validatingadmissionpolicybindings": {GVR: schema.GroupVersionResource{Group: "admissionregistration.k8s.io", Version: "v1", Resource: "validatingadmissionpolicybindings"}, Namespaced: false},
	"validatingwebhookconfigurations":   {GVR: schema.GroupVersionResource{Group: "admissionregistration.k8s.io", Version: "v1", Resource: "validatingwebhookconfigurations"}, Namespaced: false},
	"customresourcedefinitions":         {GVR: schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}, Namespaced: false},
	"apiservices":                       {GVR: schema.GroupVersionResource{Group: "apiregistration.k8s.io", Version: "v1", Resource: "apiservices"}, Namespaced: false},
	"applications":                      {GVR: schema.GroupVersionResource{Group: "application.aiops.com", Version: "v1", Resource: "applications"}, Namespaced: true},
	"controllerrevisions":               {GVR: schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "controllerrevisions"}, Namespaced: true},
	"daemonsets":                        {GVR: schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}, Namespaced: true},
	"replicasets":                       {GVR: schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"}, Namespaced: true},
	"statefulsets":                      {GVR: schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}, Namespaced: true},
	"selfsubjectreviews":                {GVR: schema.GroupVersionResource{Group: "authentication.k8s.io", Version: "v1", Resource: "selfsubjectreviews"}, Namespaced: false},
	"tokenreviews":                      {GVR: schema.GroupVersionResource{Group: "authentication.k8s.io", Version: "v1", Resource: "tokenreviews"}, Namespaced: false},
	"localsubjectaccessreviews":         {GVR: schema.GroupVersionResource{Group: "authorization.k8s.io", Version: "v1", Resource: "localsubjectaccessreviews"}, Namespaced: true},
	"selfsubjectaccessreviews":          {GVR: schema.GroupVersionResource{Group: "authorization.k8s.io", Version: "v1", Resource: "selfsubjectaccessreviews"}, Namespaced: false},
	"selfsubjectrulesreviews":           {GVR: schema.GroupVersionResource{Group: "authorization.k8s.io", Version: "v1", Resource: "selfsubjectrulesreviews"}, Namespaced: false},
	"subjectaccessreviews":              {GVR: schema.GroupVersionResource{Group: "authorization.k8s.io", Version: "v1", Resource: "subjectaccessreviews"}, Namespaced: false},
	"horizontalpodautoscalers":          {GVR: schema.GroupVersionResource{Group: "autoscaling", Version: "v2", Resource: "horizontalpodautoscalers"}, Namespaced: true},
	"cronjobs":                          {GVR: schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "cronjobs"}, Namespaced: true},
	"jobs":                              {GVR: schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"}, Namespaced: true},
	"certificaterequests":               {GVR: schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "certificaterequests"}, Namespaced: true},
	"certificates":                      {GVR: schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "certificates"}, Namespaced: true},
	"clusterissuers":                    {GVR: schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "clusterissuers"}, Namespaced: false},
	"issuers":                           {GVR: schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "issuers"}, Namespaced: true},
	"certificatesigningrequests":        {GVR: schema.GroupVersionResource{Group: "certificates.k8s.io", Version: "v1", Resource: "certificatesigningrequests"}, Namespaced: false},
	"leases":                            {GVR: schema.GroupVersionResource{Group: "coordination.k8s.io", Version: "v1", Resource: "leases"}, Namespaced: true},
	"endpointslices":                    {GVR: schema.GroupVersionResource{Group: "discovery.k8s.io", Version: "v1", Resource: "endpointslices"}, Namespaced: true},
	"flowschemas":                       {GVR: schema.GroupVersionResource{Group: "flowcontrol.apiserver.k8s.io", Version: "v1", Resource: "flowschemas"}, Namespaced: false},
	"prioritylevelconfigurations":       {GVR: schema.GroupVersionResource{Group: "flowcontrol.apiserver.k8s.io", Version: "v1", Resource: "prioritylevelconfigurations"}, Namespaced: false},
	"ingressclasses":                    {GVR: schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "ingressclasses"}, Namespaced: false},
	"ingresses":                         {GVR: schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"}, Namespaced: true},
	"networkpolicies":                   {GVR: schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"}, Namespaced: true},
	"runtimeclasses":                    {GVR: schema.GroupVersionResource{Group: "node.k8s.io", Version: "v1", Resource: "runtimeclasses"}, Namespaced: false},
	"catalogsources":                    {GVR: schema.GroupVersionResource{Group: "operators.coreos.com", Version: "v1alpha1", Resource: "catalogsources"}, Namespaced: true},
	"clusterserviceversions":            {GVR: schema.GroupVersionResource{Group: "operators.coreos.com", Version: "v1alpha1", Resource: "clusterserviceversions"}, Namespaced: true},
	"installplans":                      {GVR: schema.GroupVersionResource{Group: "operators.coreos.com", Version: "v1alpha1", Resource: "installplans"}, Namespaced: true},
	"olmconfigs":                        {GVR: schema.GroupVersionResource{Group: "operators.coreos.com", Version: "v1", Resource: "olmconfigs"}, Namespaced: false},
	"operatorconditions":                {GVR: schema.GroupVersionResource{Group: "operators.coreos.com", Version: "v2", Resource: "operatorconditions"}, Namespaced: true},
	"operatorgroups":                    {GVR: schema.GroupVersionResource{Group: "operators.coreos.com", Version: "v1", Resource: "operatorgroups"}, Namespaced: true},
	"operators":                         {GVR: schema.GroupVersionResource{Group: "operators.coreos.com", Version: "v1", Resource: "operators"}, Namespaced: false},
	"subscriptions":                     {GVR: schema.GroupVersionResource{Group: "operators.coreos.com", Version: "v1alpha1", Resource: "subscriptions"}, Namespaced: true},
	"packagemanifests":                  {GVR: schema.GroupVersionResource{Group: "packages.operators.coreos.com", Version: "v1", Resource: "packagemanifests"}, Namespaced: true},
	"poddisruptionbudgets":              {GVR: schema.GroupVersionResource{Group: "policy", Version: "v1", Resource: "poddisruptionbudgets"}, Namespaced: true},
	"clusterrolebindings":               {GVR: schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"}, Namespaced: false},
	"clusterroles":                      {GVR: schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"}, Namespaced: false},
	"rolebindings":                      {GVR: schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"}, Namespaced: true},
	"roles":                             {GVR: schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"}, Namespaced: true},
	"priorityclasses":                   {GVR: schema.GroupVersionResource{Group: "scheduling.k8s.io", Version: "v1", Resource: "priorityclasses"}, Namespaced: false},
	"csidrivers":                        {GVR: schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "csidrivers"}, Namespaced: false},
	"csinodes":                          {GVR: schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "csinodes"}, Namespaced: false},
	"csistoragecapacities":              {GVR: schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "csistoragecapacities"}, Namespaced: true},
	"storageclasses":                    {GVR: schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"}, Namespaced: false},
	"volumeattachments":                 {GVR: schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "volumeattachments"}, Namespaced: false},
}

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
	return fmt.Sprintf("Resource [%s] created successfully", unstructuredObj.GetName()), nil

}

func ListResource(ctx context.Context, namespace, resource, kubeConfig string) (string, error) {
	// client-go
	clientGo, err := utils.NewClientGo(kubeConfig)
	if err != nil {
		return "", err
	}

	var resList *unstructured.UnstructuredList
	if res, ok := resourceMap[resource]; !ok {
		return "", fmt.Errorf("resource [%s] not supported", resource)
	} else {
		if res.Namespaced {
			resList, err = clientGo.DynamicClient.Resource(res.GVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				return "", err
			}
		} else {
			resList, err = clientGo.DynamicClient.Resource(res.GVR).List(ctx, metav1.ListOptions{})
			if err != nil {
				return "", err
			}
		}
	}
	result := ""
	for _, item := range resList.Items {
		result += fmt.Sprintf("%s\n", item.GetName())
	}

	return result, nil
}

func DeleteResource(ctx context.Context, namespace, resource, resourceName, kubeConfig string) (string, error) {
	// client-go
	clientGo, err := utils.NewClientGo(kubeConfig)
	if err != nil {
		return "", err
	}

	if res, ok := resourceMap[resource]; !ok {
		return "", fmt.Errorf("resource [%s] not supported", resource)
	} else {
		fmt.Printf("Are you sure that you want to delete the resource [%s] in namespace [%s]? (yes/no): ", resourceName, namespace)
		var confirm string
		_, err := fmt.Scanln(&confirm)
		if err != nil {
			return "", err
		}

		if confirm != "yes" {
			return "Deletion aborted by user.", nil
		}

		if res.Namespaced {
			err := clientGo.DynamicClient.Resource(res.GVR).Namespace(namespace).Delete(ctx, resourceName, metav1.DeleteOptions{})
			if err != nil {
				return "", err
			}
		} else {
			err := clientGo.DynamicClient.Resource(res.GVR).Delete(ctx, resourceName, metav1.DeleteOptions{})
			if err != nil {
				return "", err
			}
		}
	}

	return fmt.Sprintf("Resource [%s] deleted successfully", resourceName), nil
}
