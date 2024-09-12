#!/bin/sh
#set -x

src=pkg/proto/cases
dst=$src

protos="\
$src/openapiv2.proto \
$src/appeal.proto \
$src/close_reason.proto \
$src/reason.proto \
$src/lookup.proto \
$src/status_condition.proto \
$src/status.proto \
$src/priority.proto \
$src/sla.proto \
$src/sla_condition.proto \
$src/catalog.proto \
$src/service.proto \
"

#openapiv2_format=yaml
openapiv2_format=json
openapiv2_filename=portal
openapiv2_file_ext=.swagger.$openapiv2_format

openapiv2_options="\
allow_merge=true\
,merge_file_name=$openapiv2_filename\
,openapi_naming_strategy=fqn\
,json_names_for_fields=false\
,disable_default_errors=true\
,repeated_path_param_separator=csv\
,allow_delete_body=true\
,logtostderr=true\
"

protoc -I pkg/proto  \
 --openapiv2_out=$openapiv2_options:$dst \
 $protos

res=$? # last command execution
#echo $res
if [ $res -ne 0 ]; then
  >&2 echo "[ERR]: protoc: failed with exit code ${res}."
  exit $res
fi

echo "[INFO]: swagger: successfully generated Swagger file."
