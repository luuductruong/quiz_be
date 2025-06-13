# =========[ Configuration ]=========
PROTOC         := protoc
IMPORT_PATHS   := -I. -I..
BASE_OUT       := ../../application/${SERVICE_NAME}
COMMON_OPTS    := --go_opt=paths=source_relative \
                  --go-grpc_opt=paths=source_relative \
                  --grpc-gateway_opt=logtostderr=true \
                  --grpc-gateway_opt=paths=source_relative

# =========[ Main Target ]=========
all:
	@if [ -f "$${SERVICE_NAME}.dto.proto" ]; then \
		echo ">> Generating DTO proto for $$SERVICE_NAME..."; \
		mkdir -p ../../application/$${SERVICE_NAME}/dto; \
		$(PROTOC) $(IMPORT_PATHS) \
			--go_out=../../application/$${SERVICE_NAME}/dto \
			--go-grpc_out=../../application/$${SERVICE_NAME}/dto \
			--grpc-gateway_out=../../application/$${SERVICE_NAME}/dto \
			$(COMMON_OPTS) \
			$${SERVICE_NAME}.dto.proto; \
		echo ">> Complete generating dto proto for $$SERVICE_NAME...\n"; \
	fi

	@if [ -f "$${SERVICE_NAME}.model.proto" ]; then \
    	echo ">> Generating MODEL proto for $$SERVICE_NAME..."; \
    	mkdir -p ../../application/$${SERVICE_NAME}/model; \
    	$(PROTOC) $(IMPORT_PATHS) \
    		--go_out=../../application/$${SERVICE_NAME}/model \
    		--go-grpc_out=../../application/$${SERVICE_NAME}/model \
   			--grpc-gateway_out=../../application/$${SERVICE_NAME}/model \
   			$(COMMON_OPTS) \
   			$${SERVICE_NAME}.model.proto; \
   		echo ">> Complete generating model proto for $$SERVICE_NAME...\n"; \
   	fi

	@if [ -f "$${SERVICE_NAME}.service.proto" ]; then \
		echo ">> Generating Service proto for $$SERVICE_NAME..."; \
		mkdir -p ../../application/$${SERVICE_NAME}/service; \
		$(PROTOC) $(IMPORT_PATHS) \
			--go_out=../../application/$${SERVICE_NAME}/service \
			--go-grpc_out=../../application/$${SERVICE_NAME}/service \
			--grpc-gateway_out=../../application/$${SERVICE_NAME}/service \
			$(COMMON_OPTS) \
			--grpc-gateway_opt=allow_delete_body=true \
			$${SERVICE_NAME}.service.proto; \
		echo ">> Complete generating service proto for $$SERVICE_NAME...\n"; \
	fi


# =========[ Clean Generated Files ]=========
clean:
	@echo ">> Cleaning generated files for ${SERVICE_NAME}..."
	@rm -f $(BASE_OUT)/dto/*.pb.go $(BASE_OUT)/dto/*.pb.gw.go 2>/dev/null || true
	@rm -f $(BASE_OUT)/model/*.pb.go $(BASE_OUT)/model/*.pb.gw.go 2>/dev/null || true
	@rm -f $(BASE_OUT)/service/*.pb.go $(BASE_OUT)/service/*.pb.gw.go 2>/dev/null || true

# =========[ Phony Targets ]=========
.PHONY: all clean
