@REM 用于初始化proto
@echo protocol generate
@echo current dir: %~dp0
@set cDir=%~dp0
@set SRC_DIR=%cDir%
@set DST_DIR="%cDir%gen\"
@echo "SRC_DIR":%SRC_DIR%
@echo "DST_DIR":%DST_DIR%

@REM protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/address.proto
@REM protoc --doc_out=./doc --doc_opt=html,index.html, proto/*.proto

::protoc -I=./proto/ --go_out=./gen/ --doc_out=./doc --doc_opt=html,index.html, proto/*.proto

@mkdir doc

@mkdir csharp_out

@REM go install github.com/golang/protobuf/protoc-gen-go@latest
@REM go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
@REM go install github.com/envoyproxy/protoc-gen-validate@latest

@protoc --csharp_out=./csharp_out --csharp_opt=file_extension=.g.cs --doc_out=./doc --doc_opt=html,index.html proto/*.proto

@pause