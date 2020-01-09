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

