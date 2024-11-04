## K8s-Copilot

### Overview

**K8s-Copilot** is a command-line tool based on [Golang](https://go.dev/) which allows you to create/list/update/delete Kubernetes built-in resources interactively powered by ChatGPT.

*Note: Model is currently hard-coded in `gpt-4o-mini`.*

### Use Cases

- Create a Kubernetes resource in a specific namespace.
- List Kubernetes either namespaced or non-namespaced resources.
- Update a Kubernetes resource given specific name & namespace.
- Delete a Kubernetes resource given specific name & namespace.

### Features

- [Cobra](https://github.com/spf13/cobra)
- [Function calling](https://platform.openai.com/docs/guides/function-calling)
- [client-go](https://github.com/kubernetes/client-go)

### Architecture

#### Cobra

```bash
k8s-copilot
├── ask
│   ├── chatgpt
└── analyze
    ├── event
```

### Demo

#### Prerequisite

An out-of-box Kubernetes cluster environment, try [kind](https://kind.sigs.k8s.io/). 👈

[Install](https://go.dev/doc/install) Golang.

#### Build

```bash
$ go build -o k8s-copilot
```

#### Setup ENV

Try [APIYI](https://www.apiyi.com/register/?aff_code=UFwG) 👈 if you had difficulty acquiring OpenAI API Key from mainland China.

```bash
# wsl env to win
$ export API_KEY="api_key"
$ export BASE_URL="base_url"
```

If you're using WSL, add below to `/etc/wsl.conf` then export the ENV.

```bash
[automount]
options = "metadata"
```

```bash
$ export WSLENV=API_KEY/w:BASE_URL/w
```

#### Run

Help

```bash
$ ./k8s-copilot -h
```

Ask

```bash
$ ./k8s-copilot ask chatgpt
```

A greeting prompt will show up.

```
Greetings, I'm a Copilot for Kubernetes, you require my assistant?
>
```

Type your queries:

*Note: open another terminal to run kubectl cmd for checking.*

```
> create a deploy named nginx, image is nginx:latest, replica is 2
```

```bash
# check
$ kubectl get deploy
```

```
> ls all pods
> ls all pods in kube-system
> ls all services in kube-system
> ls all namespaces
```

```bash
> update deploy named nginx replica to 3
```

```bash
# check
$ kubectl get deploy
```

```
> add label env=test to deploy named nginx
```

```bash
# check
$ kubectl get po --show-labels
```

```
> remove label env=test from deploy named nginx
```

```bash
# check
$ kubectl get po --show-labels
```

```
> update image of deploy named nginx to nginx:1.26.2
```

```bash
# check
$ kubectl get deploy nginx -o yaml | grep image:
```

```
> delete deploy named nginx
```

```bash
$ kubectl get deploy
```

```
> exit
```

### Implementation

See more in [UML](https://github.com/KokoiRuby/k8s-copilot/tree/main/uml). 👈

### Limitation

- Sometimes the update query is not strictly idempotent due to returned response from LLM, please try multiple times.
- The update query will create a new ReplicaSet if you try to add a label to Deployment.

### Operation and Maintenance

N/A

### Troubleshooting

N/A

### Q&A

N/A

### Reference

N/A

### TODO

- Add a flag to select LLM.
- Replace stdin with GNU-Readline.
- Add event analyzer feature.
- Fine tune system prompt to LLM to improve robustness & Idempotence.