# K8sLogChecker

Typically, log collections are collected to a log center.

But if there are abnormal characters in the logs of the log collection pod that cause log collection exceptions, the abnormal characters will be generated cyclically.

Therefore, there must be an independent component to monitor the pod logs of the log collection component.

So, there is this project.

docker image: [lqbing/k8slogchecker](https://hub.docker.com/r/lqbing/k8slogchecker)

## useage

### run out of cluster

Create .env with command `cp .env.example .env` and modify it, and run below command.

```shell
docker run --rm -it --env-file .env -v /root/.kube/config:/root/.kube/config lqbing/k8slogchecker
```

PS: If you do not want send mail, remove environment var `KLC_SEND_MAIL` from file `.env`

### run in cluster

You will received error message run in a pod in cluster with default config.

Because program in pod have no permission to access cluster resource with default config.

So, you need create `ServiceAccount`, `Role` and `ClusterRoleBinding` for it.

- gen k8s.yaml

run below command and input values with tips to generate k8s.yaml

```shell
sh ./gen_k8s_yaml.sh
```

- apply k8s.yaml

```shell
kubectl apply -f k8s.yaml
```

-  try tirgger cronjob

```shell
kubectl -n <namespace> create job --from=cronjob/k8slogchecker k8slogchecker-001
```

## environment vars

- KLC_RESOURCE_TYPE

[required] The resource type(pod/deployment/statfulset)

- KLC_NAMESPACE

The namespace(default value is 'default')

- KLC_RESOURCE_NAME

[required] The resource name

- KLC_SCAN_LOG_COUNT

The log lines you want scan(default value is 20)

- KLC_ERR_CONTENT

[required] The error strings you want alert in logs(for example 'ElasticsearchError|[error]')

- KLC_IGNORE_CONTENT

The ignore strings you want alert in logs(for example 'deprecated')

- KLC_SEND_MAIL

[required] Do you want send email when alert

If exist and error content in logs, will send email.

Otherwise, will not send email and below environment vars will not effect.

- MAIL_FROM

[required] The email from(for example 'abc@mail.com')

- MAIL_PASSWORD

[required] The email host(for example 'smtp.office365.com')

- MAIL_HOST

[required] The email smtp server port (for example 587)

- MAIL_TO

[required] The email send to (for example rec@mail.com;cc@mail.com)

- MAIL_ADD_SUBJECT

A string in the end of mail subject
