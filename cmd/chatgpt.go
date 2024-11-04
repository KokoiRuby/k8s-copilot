/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/KokoiRuby/k8s-copilot/cmd/funcs"
	"github.com/KokoiRuby/k8s-copilot/cmd/utils"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"os"

	"github.com/spf13/cobra"
)

// chatgptCmd represents the chatgpt command
var chatgptCmd = &cobra.Command{
	Use:   "chatgpt",
	Short: "ChatGPT",
	Long: `Start an interactive window where you can input the queries.
Type [exit|quit|q|bye] and press "Enter" to exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		startToChat()
	},
}

func init() {
	askCmd.AddCommand(chatgptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatgptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatgptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// 1. startToChat retrieves user input from stdin & prepares to process it.
func startToChat() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Greetings, I'm a Copilot for Kubernetes, you require my assistant?")

	for {
		ctx := context.Background()
		fmt.Print("> ")
		if scanner.Scan() {
			input := scanner.Text()
			if input == "exit" || input == "quit" || input == "q" || input == "bye" {
				fmt.Println("Have a good day, Bye!;)")
				break
			}
			if input == "" {
				continue
			}
			//fmt.Println("Your query is:", input)
			fmt.Println(processInput(ctx, input))
		}
	}
}

// 2. processInput processes user input by function calling.
func processInput(ctx context.Context, input string) string {
	client, err := utils.NewOpenAI()
	if err != nil {
		return err.Error()
	}
	resp := funcCalling(ctx, input, client)
	return resp
}

// 3. funcCalling defines the functions & prepares to invoke.
func funcCalling(ctx context.Context, input string, client *utils.OpenAI) string {
	f1 := openai.FunctionDefinition{
		Name:        "createResource",
		Description: "Create Kubernetes resource YAML manifest",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"input": {
					Type:        jsonschema.String,
					Description: "Extract verb, resource and necessary flags",
				},
			},
			Required: []string{"input"},
		},
	}
	t1 := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f1,
	}

	f2 := openai.FunctionDefinition{
		Name:        "listResource",
		Description: "List Kubernetes resources",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"namespace": {
					Type: jsonschema.String,
					Description: `The namespace where resource is. 
For non-namespaced resources, such as namespaces, persistentvolumes, 
this field shall not be set.`,
				},
				"resource": {
					Type:        jsonschema.String,
					Description: "K8s built-in resource, for example: pods, deployments, services, you can also use singular or short name (if had)",
				},
			},
			Required: []string{"namespace", "resource"},
		},
	}
	t2 := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f2,
	}

	f3 := openai.FunctionDefinition{
		Name:        "updateResource",
		Description: "Update Kubernetes resources",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"namespace": {
					Type:        jsonschema.String,
					Description: "The namespace where resource is.",
				},
				"resource": {
					Type:        jsonschema.String,
					Description: "Kubernetes built-in resource, for example: pods, deployments, services, you can also use singular or short name (if had)",
				},
				"resource_name": {
					Type:        jsonschema.String,
					Description: "Name of the resource to be deleted",
				},
				"delta": {
					Type:        jsonschema.String,
					Description: "The delta to update the resource.",
				},
			},
			Required: []string{"namespace", "resource", "resource_name", "delta"},
		},
	}
	t3 := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f3,
	}

	f4 := openai.FunctionDefinition{
		Name:        "deleteResource",
		Description: "Delete Kubernetes resources",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"namespace": {
					Type:        jsonschema.String,
					Description: "The namespace where resource is.",
				},
				"resource": {
					Type:        jsonschema.String,
					Description: "Kubernetes built-in resource, for example: pods, deployments, services, you can also use singular or short name (if had)",
				},
				"resource_name": {
					Type:        jsonschema.String,
					Description: "Name of the resource to be deleted",
				},
			},
			Required: []string{"namespace", "resource", "resource_name"},
		},
	}
	t4 := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f4,
	}

	dialogue := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		},
	}

	resp, err := client.Client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT4oMini,
			Messages: dialogue,
			Tools:    []openai.Tool{t1, t2, t3, t4},
		},
	)
	if err != nil {
		return err.Error()
	}

	msg := resp.Choices[0].Message
	if len(msg.ToolCalls) != 1 {
		return fmt.Sprintf("No appropriate tool is found, %v", len(msg.ToolCalls))
	}

	// build chat history
	dialogue = append(dialogue, msg)
	//return fmt.Sprintf("Function to call: %s, arg: %s", msg.ToolCalls[0].Function.Name, msg.ToolCalls[0].Function.Arguments)
	fmt.Printf("Function to call: %s, arg: %s\n", msg.ToolCalls[0].Function.Name, msg.ToolCalls[0].Function.Arguments)
	result, err := invokeFunc(ctx, client, msg.ToolCalls[0].Function.Name, msg.ToolCalls[0].Function.Arguments)
	if err != nil {
		return err.Error()
	}
	return result
}

// 4. invokeFunc invokes the function
func invokeFunc(ctx context.Context, client *utils.OpenAI, name, args string) (string, error) {
	switch name {
	case "createResource":
		params := struct {
			Input string `json:"input"`
		}{}
		if err := json.Unmarshal([]byte(args), &params); err != nil {
			return "", err
		}
		return funcs.CreateResource(ctx, client, params.Input, kubeconfig)
	case "listResource":
		params := struct {
			Namespace string `json:"namespace"`
			Resource  string `json:"resource"`
		}{}
		if err := json.Unmarshal([]byte(args), &params); err != nil {
			return "", err
		}
		return funcs.ListResource(ctx, params.Namespace, params.Resource, kubeconfig)
	default:
		return "", fmt.Errorf("unknown function %s", name)
	}
}
