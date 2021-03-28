rwildcard = $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) $(filter $(subst *,%,$2),$d))
my-filter-out = $(foreach v,$(2),$(if $(findstring $(1),$(v)),,$(v)))
testfiles ?= $(sort $(call my-filter-out,node_modules/,$(call rwildcard,../..,*.md)))

all: test

test: info offensive-wording remark markdownlint spellcheck
	@echo ""
	@echo "--- OK. ---"
.PHONY: test

info: deps
	@echo ""
	@echo "==> Files to check:"
	@echo "  $(testfiles:../.%=%)"
.PHONY: info

offensive-wording: deps
	@echo ""
	@echo "==> Checking for inconsiderate/insensitive wording..."
	@node_modules/.bin/alex --diff $(testfiles) || true
.PHONY: offensive-wording

spellcheck: deps
	@echo ""
	@echo "==> Checking for spelling errors..."
	@node_modules/.bin/mdspell --report --en-us --ignore-numbers --ignore-acronyms $(testfiles)
.PHONY: spellcheck

remark: deps
	@echo ""
	@echo "==> Checking for messy formatting..."
	@node_modules/.bin/remark --frail $(testfiles)
.PHONY: remark

markdownlint: deps
	@echo ""
	@echo "==> Checking for good practices..."
	@node_modules/.bin/markdownlint $(testfiles)
.PHONY: markdownlint

deps: node_modules/
node_modules/: package.json
	npm install
	touch $@
.PHONY: deps
