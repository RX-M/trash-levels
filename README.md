# trash-levels

Simple service to report city trash can levels in golang


## Microservice component of CICD demo

The raw steps for the microservice container build portion of the RX-M Kubernetes CICD demo.


### Update box and install docker

```
user@ubuntu:~$ sudo apt-get update

Hit:1 http://us.archive.ubuntu.com/ubuntu xenial InRelease
Hit:2 http://us.archive.ubuntu.com/ubuntu xenial-updates InRelease                                   
Hit:3 http://security.ubuntu.com/ubuntu xenial-security InRelease                                    
Hit:4 http://us.archive.ubuntu.com/ubuntu xenial-backports InRelease    
Reading package lists... Done                     

user@ubuntu:~$ wget -O - https://get.docker.com | sh

--2020-01-09 12:56:19--  https://get.docker.com/
Resolving get.docker.com (get.docker.com)... 13.226.219.90, 13.226.219.89, 13.226.219.61, ...
Connecting to get.docker.com (get.docker.com)|13.226.219.90|:443... connected.
HTTP request sent, awaiting response... 200 OK
Length: 13216 (13K) [text/plain]
Saving to: ‘STDOUT’

-                            100%[==============================================>]  12.91K  --.-KB/s    in 0s      

2020-01-09 12:56:19 (229 MB/s) - written to stdout [13216/13216]

# Executing docker install script, commit: f45d7c11389849ff46a6b4d94e0dd1ffebca32c1
+ sudo -E sh -c apt-get update -qq >/dev/null
[sudo] password for user:
+ sudo -E sh -c DEBIAN_FRONTEND=noninteractive apt-get install -y -qq apt-transport-https ca-certificates curl >/dev/null
+ sudo -E sh -c curl -fsSL "https://download.docker.com/linux/ubuntu/gpg" | apt-key add -qq - >/dev/null
+ sudo -E sh -c echo "deb [arch=amd64] https://download.docker.com/linux/ubuntu xenial stable" > /etc/apt/sources.list.d/docker.list
+ sudo -E sh -c apt-get update -qq >/dev/null
+ [ -n  ]
+ sudo -E sh -c apt-get install -y -qq --no-install-recommends docker-ce >/dev/null
+ sudo -E sh -c docker version
Client: Docker Engine - Community
 Version:           19.03.5
 API version:       1.40
 Go version:        go1.12.12
 Git commit:        633a0ea838
 Built:             Wed Nov 13 07:50:12 2019
 OS/Arch:           linux/amd64
 Experimental:      false

Server: Docker Engine - Community
 Engine:
  Version:          19.03.5
  API version:      1.40 (minimum version 1.12)
  Go version:       go1.12.12
  Git commit:       633a0ea838
  Built:            Wed Nov 13 07:48:43 2019
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          1.2.10
  GitCommit:        b34a5c8af56e510852c35414db4c1f4fa6172339
 runc:
  Version:          1.0.0-rc8+dev
  GitCommit:        3e425f80a8c931f88e6d94a8c831b9d5aa481657
 docker-init:
  Version:          0.18.0
  GitCommit:        fec3683
If you would like to use Docker as a non-root user, you should now consider
adding your user to the "docker" group with something like:

  sudo usermod -aG docker user

Remember that you will have to log out and back in for this to take effect!

WARNING: Adding a user to the "docker" group will grant the ability to run
         containers which can be used to obtain root privileges on the
         docker host.
         Refer to https://docs.docker.com/engine/security/security/#docker-daemon-attack-surface
         for more information.

user@ubuntu:~$ sudo usermod -aG docker user

user@ubuntu:~$
```


### Setup go on the host

You could of course just run go in containers all the time, but its nice to be able to dev on the host if go is your
primary language and dev is your primary thing.

```
user@ubuntu:~$ curl -sLO https://dl.google.com/go/go1.13.linux-amd64.tar.gz

user@ubuntu:~$ tar zxf go1.13.linux-amd64.tar.gz

user@ubuntu:~$ sudo mv ~/go/ /usr/local/

user@ubuntu:~$ echo "export PATH=/usr/local/go/bin:$PATH" >> .bashrc

user@ubuntu:~$ echo "[[ -r ~/.bashrc ]] && . ~/.bashrc" >> ~/.bash_profile

user@ubuntu:~$ source ~/.bash_profile

user@ubuntu:~$ go version

go version go1.13 linux/amd64

user@ubuntu:~$
```


### Clone trash levels app

```
user@ubuntu:~$ git clone git@github.com:RX-M/trash-levels.git

Cloning into 'trash-levels'...
Warning: Permanently added the RSA host key for IP address '192.30.255.112' to the list of known hosts.
Enter passphrase for key '/home/user/.ssh/id_rsa':
remote: Enumerating objects: 5, done.
remote: Counting objects: 100% (5/5), done.
remote: Compressing objects: 100% (5/5), done.
remote: Total 5 (delta 0), reused 0 (delta 0), pack-reused 0
Receiving objects: 100% (5/5), 4.77 KiB | 0 bytes/s, done.
Checking connectivity... done.

user@ubuntu:~$ cd trash-levels/

user@ubuntu:~/trash-levels$ ll


user@ubuntu:~/trash-levels$ go get github.com/gorilla/mux

```


### Test trash levels

```
user@ubuntu:~/trash-levels$ go run main.go
2020/01/09 12:47:06 Listening on: 8080

```

in another shell:

```
user@ubuntu:~/trash-levels$ curl localhost:8080/cans/10

{"id": "10","level": 80}

user@ubuntu:~/trash-levels$ curl -vv localhost:8080/cans/25

*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> GET /cans/25 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.47.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Thu, 09 Jan 2020 21:10:13 GMT
< Content-Length: 24
<
* Connection #0 to host localhost left intact
{"id": "25","level": 83}

user@ubuntu:~/trash-levels$
```

shut down the service:

```
user@ubuntu:~/trash-levels$ go run main.go
2020/01/09 12:47:06 Listening on: 8080
2020/01/09 12:47:10 Responding to request for level: 10
2020/01/09 12:47:30 Responding to request for level: 25
^Csignal: interrupt
user@ubuntu:~/trash-levels$
```


### Build and push the image


```
user@ubuntu:~/trash-levels$ make docker

docker build -t rxmllc/trash-levels .
Sending build context to Docker daemon   68.1kB
Step 1/13 : FROM golang:1.13 AS build-env
 ---> ed081345a3da
Step 2/13 : WORKDIR /go/src/trash-levels/
 ---> Using cache
 ---> 45da8be5c544
Step 3/13 : COPY ./main.go /go/src/trash-levels/
 ---> Using cache
 ---> d3e2d78ad9da
Step 4/13 : RUN go get github.com/gorilla/mux/
 ---> Using cache
 ---> 14cc2907ac45
Step 5/13 : RUN ["go","build","-tags","netgo"]
 ---> Using cache
 ---> bbecb2c91f4b
Step 6/13 : FROM scratch
 --->
Step 7/13 : LABEL maintainer="rx-m llc <info@rx-m.com>"
 ---> Using cache
 ---> 7dc527308ec4
Step 8/13 : LABEL org.label-schema.name="trash levels"
 ---> Using cache
 ---> 32319f00d523
Step 9/13 : LABEL org.label-schema.vendor="rx-m llc"
 ---> Using cache
 ---> 62c6051b163c
Step 10/13 : LABEL org.label-schema.schema-version="1.0"
 ---> Using cache
 ---> cd36122dd213
Step 11/13 : COPY --from=build-env /go/src/trash-levels/trash-levels trash-levels
 ---> 8f9ae93c6ba4
Step 12/13 : EXPOSE 8080
 ---> Running in d230677ddd0f
Removing intermediate container d230677ddd0f
 ---> a5fe308f4f5e
Step 13/13 : ENTRYPOINT ["./trash-levels"]
 ---> Running in 710f239a5ee8
Removing intermediate container 710f239a5ee8
 ---> c678b93278bf
Successfully built c678b93278bf
Successfully tagged rxmllc/trash-levels:latest

user@ubuntu:~/trash-levels$ docker push rxmllc/trash-levels

The push refers to repository [docker.io/rxmllc/trash-levels]
aef3e3776e46: Pushed
latest: digest: sha256:ebcb545f23c6e1c8e52e6ac5b77b223db4c5c396981654ee3c0ef9e7fc70c86e size: 528

user@ubuntu:~/trash-levels$
```


### Setup k8s access

Copy your admin config over.

Install kubectl:

```
user@ubuntu:~$ curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl

  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 41.4M  100 41.4M    0     0  10.5M      0  0:00:03  0:00:03 --:--:-- 10.5M

user@ubuntu:~$ chmod +x ./kubectl

user@ubuntu:~$ sudo mv ./kubectl /usr/local/bin/kubectl

user@ubuntu:~$ kubectl version

Client Version: version.Info{Major:"1", Minor:"17", GitVersion:"v1.17.0", GitCommit:"70132b0f130acc0bed193d9ba59dd186f0e634cf", GitTreeState:"clean", BuildDate:"2019-12-07T21:20:10Z", GoVersion:"go1.13.4", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"17", GitVersion:"v1.17.0", GitCommit:"70132b0f130acc0bed193d9ba59dd186f0e634cf", GitTreeState:"clean", BuildDate:"2019-12-07T21:12:17Z", GoVersion:"go1.13.4", Compiler:"gc", Platform:"linux/amd64"}

user@ubuntu:~$
```


### Run the app

Post the deployment:

```
user@ubuntu:~/trash-levels/k8s$ kubectl create -f deploy.yaml

deployment.apps/trash-levels created

user@ubuntu:~/trash-levels/k8s$
```

Test it:

```
user@ubuntu:~$ kubectl get pod -o wide

NAME                           READY   STATUS    RESTARTS   AGE   IP          NODE        NOMINATED NODE   READINESS GATES
trash-levels-9dbfd7765-5wdtg   1/1     Running   0          19m   10.44.0.1   k8sworker   <none>           <none>
trash-levels-9dbfd7765-9lv79   1/1     Running   0          19m   10.47.0.1   k8srouter   <none>           <none>

user@ubuntu:~/trash-levels$ kubectl run -it test --rm --image=busybox --generator=run-pod/v1

If you don't see a command prompt, try pressing enter.

/ # wget -q -O - 10.44.0.1:8080/cans/10

{"id": "10","level": 80}

/ # exit

Session ended, resume using 'kubectl attach test -c test -i -t' command when the pod is running
pod "test" deleted

user@ubuntu:~/trash-levels$
```

Create the service:

```
user@ubuntu:~/trash-levels/k8s$ kubectl create -f service.yaml

service/trash-levels created

user@ubuntu:~/trash-levels/k8s$ kubectl get service

NAME           TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
kubernetes     ClusterIP   10.96.0.1       <none>        443/TCP   84m
trash-levels   ClusterIP   10.96.174.113   <none>        80/TCP    2m12s

user@ubuntu:~/trash-levels/k8s$
```

test it:

```
user@ubuntu:~/trash-levels$ kubectl describe service trash-levels

Name:              trash-levels
Namespace:         default
Labels:            app=trash-levels
                   project=cicd-demo
                   vender=rx-m
Annotations:       <none>
Selector:          app=trash-levels
Type:              ClusterIP
IP:                10.96.254.212
Port:              <unset>  80/TCP
TargetPort:        8080/TCP
Endpoints:         10.44.0.1:8080,10.47.0.1:8080
Session Affinity:  None
Events:            <none>

user@ubuntu:~/trash-levels$ kubectl run -it test --rm --image=busybox --generator=run-pod/v1

If you don't see a command prompt, try pressing enter.

/ # wget -q -O - http://trash-levels/cans/10

{"id": "10","level": 80}

/ # exit

Session ended, resume using 'kubectl attach test -c test -i -t' command when the pod is running
pod "test" deleted

user@ubuntu:~/trash-levels$
```

Clean up:

```
user@ubuntu:~/trash-levels/trash-levels$ kubectl get deploy,service

NAME                           READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/trash-levels   2/2     2            2           83m

NAME                   TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
service/kubernetes     ClusterIP   10.96.0.1       <none>        443/TCP   162m
service/trash-levels   ClusterIP   10.96.254.212   <none>        80/TCP    57m

user@ubuntu:~/trash-levels/trash-levels$ kubectl delete service/trash-levels

service "trash-levels" deleted

user@ubuntu:~/trash-levels/trash-levels$ kubectl delete deployment.apps/trash-levels

deployment.apps "trash-levels" deleted

user@ubuntu:~/trash-levels/trash-levels$
```


### Use helm

Install helm:

```
user@ubuntu:~$ wget https://get.helm.sh/helm-v3.0.2-linux-amd64.tar.gz

--2020-01-09 14:51:17--  https://get.helm.sh/helm-v3.0.2-linux-amd64.tar.gz
Resolving get.helm.sh (get.helm.sh)... 152.195.19.97, 2606:2800:11f:1cb7:261b:1f9c:2074:3c
Connecting to get.helm.sh (get.helm.sh)|152.195.19.97|:443... connected.
HTTP request sent, awaiting response... 200 OK
Length: 12101232 (12M) [application/x-tar]
Saving to: ‘helm-v3.0.2-linux-amd64.tar.gz’

helm-v3.0.2-linux-amd64.tar.gz  100%[====================================================>]  11.54M  11.5MB/s    in 1.0s    

2020-01-09 14:51:18 (11.5 MB/s) - ‘helm-v3.0.2-linux-amd64.tar.gz’ saved [12101232/12101232]

user@ubuntu:~$ tar -zxvf helm-v3.0.2-linux-amd64.tar.gz

linux-amd64/
linux-amd64/README.md
linux-amd64/LICENSE
linux-amd64/helm

user@ubuntu:~$ sudo mv linux-amd64/helm /usr/local/bin/helm

user@ubuntu:~$ helm version

version.BuildInfo{Version:"v3.0.2", GitCommit:"19e47ee3283ae98139d98460de796c1be1e3975f", GitTreeState:"clean", GoVersion:"go1.13.5"}

user@ubuntu:~$
```

Setup the chart:

```
user@ubuntu:~/trash-levels$ helm create trash-levels

Creating trash-levels

user@ubuntu:~/trash-levels$ cd trash-levels/

user@ubuntu:~/trash-levels/trash-levels$ ll

total 28
drwxr-xr-x 4 user user 4096 Jan  9 14:55 ./
drwxrwxr-x 5 user user 4096 Jan  9 14:55 ../
drwxr-xr-x 2 user user 4096 Jan  9 14:55 charts/
-rw-r--r-- 1 user user  910 Jan  9 14:55 Chart.yaml
-rw-r--r-- 1 user user  342 Jan  9 14:55 .helmignore
drwxr-xr-x 3 user user 4096 Jan  9 14:55 templates/
-rw-r--r-- 1 user user 1495 Jan  9 14:55 values.yaml

user@ubuntu:~/trash-levels/trash-levels$
```

Make edits until things look like this:


```
user@ubuntu:~/trash-levels/trash-levels$ cat *

cat: charts: Is a directory
apiVersion: v2
name: trash-levels
description: A Helm chart for Kubernetes

# A chart can be either an 'application' or a 'library' chart.
#
# Application charts are a collection of templates that can be packaged into versioned archives
# to be deployed.
#
# Library charts provide useful utilities or functions for the chart developer. They're included as
# a dependency of application charts to inject those utilities and functions into the rendering
# pipeline. Library charts do not define any templates and therefore cannot be deployed.
type: application

# This is the chart version. This version number should be incremented each time you make changes
# to the chart and its templates, including the app version.
version: 0.1.0

# This is the version number of the application being deployed. This version number should be
# incremented each time you make changes to the application.
appVersion: 0.1.0
cat: templates: Is a directory
# Default values for trash-levels.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

replicaCount: 2

containerPort: 8080

image:
  repository: rxmllc/trash-levels
  pullPolicy: IfNotPresent

imagePullSecrets: []
nameOverride: ""
fullnameOverride: ""

serviceAccount:
  # Specifies whether a service account should be created
  create: false
  # The name of the service account to use.
  # If not set and create is true, a name is generated using the fullname template
  name:

podSecurityContext: {}
  # fsGroup: 2000

securityContext: {}
  # capabilities:
  #   drop:
  #   - ALL
  # readOnlyRootFilesystem: true
  # runAsNonRoot: true
  # runAsUser: 1000

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: true
  annotations: {}
    # kubernetes.io/ingress.class: nginx
  hosts:
    - host: trashlevel.rx-m.com
  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

resources: {}
  # We usually recommend not to specify default resources and to leave this as a conscious
  # choice for the user. This also increases chances charts run on environments with little
  # resources, such as Minikube. If you do want to specify resources, uncomment the following
  # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
  # limits:
  #   cpu: 100m
  #   memory: 128Mi
  # requests:
  #   cpu: 100m
  #   memory: 128Mi

nodeSelector: {}

tolerations: []

affinity: {}

user@ubuntu:~/trash-levels/trash-levels$
```

... and the templates look like this:

```
user@ubuntu:~/trash-levels/trash-levels$ ls -l templates/

total 24
-rw-r--r-- 1 user user  591 Jan  9 15:19 deployment.yaml
-rw-r--r-- 1 user user 1897 Jan  9 14:55 _helpers.tpl
-rw-r--r-- 1 user user 1040 Jan  9 14:55 ingress.yaml
-rw-r--r-- 1 user user 1601 Jan  9 14:55 NOTES.txt
-rw-r--r-- 1 user user  399 Jan  9 15:25 service.yaml
drwxr-xr-x 2 user user 4096 Jan  9 14:55 tests

user@ubuntu:~/trash-levels/trash-levels$ cat templates/*

apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "trash-levels.fullname" . }}
  labels:
    {{- include "trash-levels.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "trash-levels.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "trash-levels.selectorLabels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Chart.AppVersion }}"
          ports:
          - containerPort: 8080
{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "trash-levels.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "trash-levels.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "trash-levels.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "trash-levels.labels" -}}
helm.sh/chart: {{ include "trash-levels.chart" . }}
{{ include "trash-levels.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "trash-levels.selectorLabels" -}}
app.kubernetes.io/name: {{ include "trash-levels.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "trash-levels.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "trash-levels.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}
{{- if .Values.ingress.enabled -}}
{{- $fullName := include "trash-levels.fullname" . -}}
{{- $svcPort := .Values.service.port -}}
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: {{ $fullName }}
  labels:
    {{- include "trash-levels.labels" . | nindent 4 }}
spec:
  rules:
  {{- range .Values.ingress.hosts }}
    - host: {{ .host | quote }}
      http:
        paths:
          - backend:
              serviceName: {{ $fullName }}
              servicePort: {{ $svcPort }}
  {{- end }}
{{- end }}
1. Get the application URL by running these commands:
{{- if .Values.ingress.enabled }}
{{- range $host := .Values.ingress.hosts }}
  {{- range .paths }}
  http{{ if $.Values.ingress.tls }}s{{ end }}://{{ $host.host }}{{ . }}
  {{- end }}
{{- end }}
{{- else if contains "NodePort" .Values.service.type }}
  export NODE_PORT=$(kubectl get --namespace {{ .Release.Namespace }} -o jsonpath="{.spec.ports[0].nodePort}" services {{ include "trash-levels.fullname" . }})
  export NODE_IP=$(kubectl get nodes --namespace {{ .Release.Namespace }} -o jsonpath="{.items[0].status.addresses[0].address}")
  echo http://$NODE_IP:$NODE_PORT
{{- else if contains "LoadBalancer" .Values.service.type }}
     NOTE: It may take a few minutes for the LoadBalancer IP to be available.
           You can watch the status of by running 'kubectl get --namespace {{ .Release.Namespace }} svc -w {{ include "trash-levels.fullname" . }}'
  export SERVICE_IP=$(kubectl get svc --namespace {{ .Release.Namespace }} {{ include "trash-levels.fullname" . }} --template "{{"{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}"}}")
  echo http://$SERVICE_IP:{{ .Values.service.port }}
{{- else if contains "ClusterIP" .Values.service.type }}
  export POD_NAME=$(kubectl get pods --namespace {{ .Release.Namespace }} -l "app.kubernetes.io/name={{ include "trash-levels.name" . }},app.kubernetes.io/instance={{ .Release.Name }}" -o jsonpath="{.items[0].metadata.name}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl --namespace {{ .Release.Namespace }} port-forward $POD_NAME 8080:80
{{- end }}
apiVersion: v1
kind: Service
metadata:
  name: {{ include "trash-levels.fullname" . }}
  labels:
    {{- include "trash-levels.labels" . | nindent 4 }}
spec:
  type: {{ .Values.service.type }}
  ports:
    - port: {{ .Values.service.port }}
      targetPort: {{ .Values.containerPort }}
      protocol: TCP
      name: http
  selector:
    {{- include "trash-levels.selectorLabels" . | nindent 4 }}
cat: templates/tests: Is a directory

user@ubuntu:~/trash-levels/trash-levels$
```

... and the test looks like this:

```
user@ubuntu:~/trash-levels/trash-levels$ cat templates/tests/*

apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "trash-levels.fullname" . }}-test-connection"
  labels:
{{ include "trash-levels.labels" . | nindent 4 }}
  annotations:
    "helm.sh/hook": test-success
spec:
  containers:
    - name: wget
      image: busybox
      command: ['wget']
      args:  ['-q', '-O', '-', '{{ include "trash-levels.fullname" . }}:{{ .Values.service.port }}/cans/25']
  restartPolicy: Never

user@ubuntu:~/trash-levels/trash-levels$
```

... and the templates should render like this:

```
user@ubuntu:~/trash-levels/trash-levels$ helm template test-release .

---
# Source: trash-levels/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: test-release-trash-levels
  labels:
    helm.sh/chart: trash-levels-0.1.0
    app.kubernetes.io/name: trash-levels
    app.kubernetes.io/instance: test-release
    app.kubernetes.io/version: "0.1.0"
    app.kubernetes.io/managed-by: Helm
spec:
  type: ClusterIP
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
      name: http
  selector:
    app.kubernetes.io/name: trash-levels
    app.kubernetes.io/instance: test-release
---
# Source: trash-levels/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-release-trash-levels
  labels:
    helm.sh/chart: trash-levels-0.1.0
    app.kubernetes.io/name: trash-levels
    app.kubernetes.io/instance: test-release
    app.kubernetes.io/version: "0.1.0"
    app.kubernetes.io/managed-by: Helm
spec:
  replicas: 2
  selector:
    matchLabels:
      app.kubernetes.io/name: trash-levels
      app.kubernetes.io/instance: test-release
  template:
    metadata:
      labels:
        app.kubernetes.io/name: trash-levels
        app.kubernetes.io/instance: test-release
    spec:
      containers:
        - name: trash-levels
          image: "rxmllc/trash-levels:0.1.0"
          ports:
          - containerPort: 8080
---
# Source: trash-levels/templates/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: test-release-trash-levels
  labels:
    helm.sh/chart: trash-levels-0.1.0
    app.kubernetes.io/name: trash-levels
    app.kubernetes.io/instance: test-release
    app.kubernetes.io/version: "0.1.0"
    app.kubernetes.io/managed-by: Helm
spec:
  rules:
    - host: "trashlevel.rx-m.com"
      http:
        paths:
          - backend:
              serviceName: test-release-trash-levels
              servicePort: 80
---
# Source: trash-levels/templates/tests/test-connection.yaml
apiVersion: v1
kind: Pod
metadata:
  name: "test-release-trash-levels-test-connection"
  labels:

    helm.sh/chart: trash-levels-0.1.0
    app.kubernetes.io/name: trash-levels
    app.kubernetes.io/instance: test-release
    app.kubernetes.io/version: "0.1.0"
    app.kubernetes.io/managed-by: Helm
  annotations:
    "helm.sh/hook": test-success
spec:
  containers:
    - name: wget
      image: busybox
      command: ['wget']
      args:  ['-q', '-O', '-', 'test-release-trash-levels:80/cans/25']
  restartPolicy: Never

user@ubuntu:~/trash-levels/trash-levels$
```

Create the resources:

```
user@ubuntu:~/trash-levels/trash-levels$ helm template test-release . | kubectl create -f -
service/test-release-trash-levels created
deployment.apps/test-release-trash-levels created
ingress.networking.k8s.io/test-release-trash-levels created
pod/test-release-trash-levels-test-connection created

user@ubuntu:~/trash-levels/trash-levels$
```

check the test:

```
user@ubuntu:~/trash-levels/trash-levels$ kubectl get deploy,service,ing

NAME                                        READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/ing-trash-level             1/1     1            1           38m
deployment.apps/test-release-trash-levels   2/2     2            2           12s

NAME                                TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
service/ing-trash-level             ClusterIP   10.96.248.184   <none>        80/TCP    21m
service/kubernetes                  ClusterIP   10.96.0.1       <none>        443/TCP   4h7m
service/test-release-trash-levels   ClusterIP   10.96.226.16    <none>        80/TCP    12s

NAME                                           HOSTS                 ADDRESS   PORTS   AGE
ingress.extensions/test-release-trash-levels   trashlevel.rx-m.com             80      12s
ingress.extensions/trash-level-ingress         www.example.com                 80      26m

user@ubuntu:~/trash-levels/trash-levels$ kubectl get pod

NAME                                         READY   STATUS      RESTARTS   AGE
ing-trash-level-6784f7d4b-j4k27              1/1     Running     0          38m
test-release-trash-levels-5df9f7c967-87zwk   1/1     Running     0          23s
test-release-trash-levels-5df9f7c967-hncvj   1/1     Running     0          23s
test-release-trash-levels-test-connection    0/1     Completed   0          23s

user@ubuntu:~/trash-levels/trash-levels$ kubectl logs test-release-trash-levels-test-connection

{"id": "25","level": 83}

user@ubuntu:~/trash-levels/trash-levels$
```

Clean up:

```
user@ubuntu:~/trash-levels/trash-levels$ helm template test-release . | kubectl delete -f -

service "test-release-trash-levels" deleted
deployment.apps "test-release-trash-levels" deleted
ingress.networking.k8s.io/test-release-trash-levels deleted
pod "test-release-trash-levels-test-connection" deleted

user@ubuntu:~/trash-levels/trash-levels$
```

