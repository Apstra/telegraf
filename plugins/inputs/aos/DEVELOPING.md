# Building apstra-telegraf

## Install go
```
wget https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.12.9.linux-amd64.tar.gz
export GOPATH=/home/admin/go-projects
export GOROOT=/usr/local/go
```

Install dep
```
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```

## Get the source
```
go get -d github.com/Apstra/telegraf
```
Need to rename Apstra to influxdata to avoid: use of internal package not allowed

## Update aosstream.pb.go
```
protoc --go_out=. aosstream.proto
```

## Build telegraf container
```
cd $GOPATH/src/github.com/influxdata/telegraf
make build-for-docker
docker build -t apstra/telegraf:1.5.2_AOS_3.0.1 -t apstra/telegraf:latest .
```
The tag indicates that we are using telegraf 1.5.2 with the protobuf schema from AOS 3.0.1.