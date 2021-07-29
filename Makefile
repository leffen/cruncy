# Made by Leif Terje Fonnes, authumn 2017
# requires that bump version is installed	go get github.com/Shyp/bump_version


.PHONY: test

deps:
	go get go.elastic.co/go-licence-detector
	go get github.com/Shyp/bump_version

bump:
	bump_version patch cruncy.go

push:
	git push origin --tags -f

license:
	go list -m -json all | go-licence-detector \
		-includeIndirect  \
		-depsOut=dependencies.asciidoc \
		-rules templates/rules.json \
		-noticeOut=NOTICE.txt \
		-overrides templates/overrides.json \
		-depsTemplate templates/dependencies.asciidoc.tmpl \
		-noticeTemplate templates/NOTICE.txt.tmpl

release: test bump push

test:
	go test ./... -v -cover -race