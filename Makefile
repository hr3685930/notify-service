PROTO_DIR ?= notify
K8S_SVC ?= app-notifyservice
K8S_SVC_PORT ?= 8080

PROTOUSER=$(u)
PROTOPWD=$(p)
NAMESPACE=$(ns)

# make protodep u=xxx p=xxx
protodep:
	protodep up -f --basic-auth-username $(PROTOUSER)  --basic-auth-password $(PROTOPWD)

exchange:
	sudo ktctl -n $(NAMESPACE) exchange $(K8S_SVC) --expose $(K8S_SVC_PORT)

connect:
	sudo ktctl connect

proto:
	protoc --go_out=plugins=grpc:. api/proto/$(PROTO_DIR)/*.proto
