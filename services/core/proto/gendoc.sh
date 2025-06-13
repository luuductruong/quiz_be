# Tạo thư mục đầu ra nếu chưa có, phục vụ documenting về sau.
mkdir -p ../swagger

PROTOS=$(find . -type f -name "*service.proto" -print)
protoc \
    -I. -I.. \
    --openapiv2_out ../swagger \
    --openapiv2_opt allow_merge=true \
    --openapiv2_opt merge_file_name=quiz \
    --openapiv2_opt json_names_for_fields=false \
    --openapiv2_opt output_format=json \
    --openapiv2_opt logtostderr=true \
    --openapiv2_opt allow_delete_body=true \
    --openapiv2_opt omit_enum_default_value=false \
    --openapiv2_opt disable_default_errors=true \
    $PROTOS