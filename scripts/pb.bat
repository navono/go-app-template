@echo off

REM set PROTOC_GEN_TS_PATH=%cd%\\react-client\\node_modules\\.bin\\protoc-gen-ts.cmd
REM set JS_DIR=.\\react-client\\src

REM --js_out=import_style=commonjs,binary:%JS_DIR%  ^
REM --ts_out=%JS_DIR%  ^

protoc.exe -I. ^
	--plugin="protoc-gen-ts=%PROTOC_GEN_TS_PATH%" ^
	--go_out=. ^
	../api/greeter.proto
