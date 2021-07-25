protoc --proto_path=proto --go_out=glaccount --go_opt=paths=source_relative organization.proto 
protoc --proto_path=proto --go_out=glaccount --go_opt=paths=source_relative account.proto 
protoc --proto_path=proto --go_out=glaccount --go_opt=paths=source_relative journal.proto 