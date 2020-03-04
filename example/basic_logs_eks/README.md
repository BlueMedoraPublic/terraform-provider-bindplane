# logs eks example

## goal

1) define a destination
2) use EKS module to define an EKS source config
3) use Terraform outputs to retrieve the EKS log template id

The log template id can be used to deploy the log agent to EKS

## setup

Special care is required when handling the service account used to connect
BindPlane to Stackdriver.

Use GCP Secret Manager to store your `service_account.json` file.
Terraform retrieves the secret at runtime
```
data "google_secret_manager_secret_version" "bindplane_svc_act" {
  provider = google-beta
  project = var.project
  secret = "bindplane-service-account"
  version = 1
}
```

## usage

Set env
```
export BINDPLANE_API_KEY=<your api key>
export PROJECT=<your gcp project id>
```

Deploy
```
terraform init
terraform apply -var "project=${PROJECT}"
terraform refresh -var "project=${PROJECT}" > /dev/null

terraform output
```

Use bpcli to retrieve a K8s yaml deployment
```
bpcli logs agent install-cmd --platform kubernetes

apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: bindplane-log-agent
spec:
  selector:
    matchLabels:
      name: bindplane-log-agent
  template:
    metadata:
      creationTimestamp: null
      labels:
        name: bindplane-log-agent
    spec:
      serviceAccountName: default
      tolerations:
        - key: node-role.kubernetes.io/master
          effect: NoSchedule
      containers:
        - image: docker.io/bluemedora/bindplane-log-agent:0.8.0
          imagePullPolicy: Always
          name: bindplane-log-agent
          resources:
            limits:
              memory: "250Mi"
              cpu: 250m
            requests:
              memory: "250Mi"
              cpu: 100m
          volumeMounts:
            - mountPath: /var/log
              name: varlog
            - mountPath: /var/lib/docker/containers
              name: varlibdockercontainers
              readOnly: true
          env:
            - name: AGENT_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: COMPANY_ID
              value: '77f054bf-31fa-4ace-b1b3-1ad2432b7ac9'
            - name: SECRET_KEY
              value: 'e223ccf9-3e93-4f62-ba6e-2a9803b41bc2'
            - name: TEMPLATE_ID
              value: ''
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      volumes:
        - hostPath:
            path: /var/log
          name: varlog
        - hostPath:
            path: /var/lib/docker/containers
          name: varlibdockercontainers
```

Retrieve the template id and use it with the deployment yaml
```
terraform output --json eks | jq .template_id
```
