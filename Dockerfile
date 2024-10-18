FROM golang:1.21 AS build-env
WORKDIR /go/src/trash-levels/
COPY ./main.go /go/src/trash-levels/
RUN go mod init example.com/trash-levels
RUN go get github.com/gorilla/mux/
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN ["go","build","-tags","netgo"]

FROM scratch
LABEL maintainer="rx-m llc <info@rx-m.com>"
LABEL org.label-schema.name="trash levels"
LABEL org.label-schema.vendor="rx-m llc"
LABEL org.label-schema.schema-version="1.1"
COPY --from=build-env /go/src/trash-levels/trash-levels trash-levels
EXPOSE 8080
ENTRYPOINT ["./trash-levels"]

