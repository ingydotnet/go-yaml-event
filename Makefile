# Using the "Makes" Makefile setup - https://github.com/makeplus/makes
M := $(or $(MAKES_REPO_DIR),.cache/makes)
$(shell [ -d $M ] || git clone -q https://github.com/makeplus/makes $M)
include $M/init.mk
include $M/clean.mk
include $M/shell.mk
include $M/ys.mk

GO-YAML := go-yaml
GO-DEPS := $(GO-YAML)
include $M/go.mk

GO-VERSION := 1.23.4

GO-YAML-URL := https://github.com/yaml/go-yaml
GO-YAML-PATCH := go-yaml-patch


MAKES-REALCLEAN := $(GO-DEPS)


ys: $(YS)

$(GO-YAML): $(GO-YAML-PATCH)
	git clone --depth 1 -q https://github.com/yaml/go-yaml $@
	(cd $@ && ln -s ../$</events.go)
