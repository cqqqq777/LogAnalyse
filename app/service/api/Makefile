generate:
	hz new -idl ./../../idl/api/$(service).thrift -mod LogAnalyse/app/service/api -force
	rm go.mod .gitignore
	git checkout router.go .hz main.go
user:
	make service=user generate
all:
	make user