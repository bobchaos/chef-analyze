# This Makefile only supports Unix systems
# TODO @afiune make it compatible for Windows
UNAME_S:=$(shell uname -s)
ifeq ($(UNAME_S),Linux)
  PLATFORM:=linux_amd64
endif
ifeq ($(UNAME_S),Darwin)
  PLATFORM:=darwin_amd64
endif
BINARY:=chef-analyze_$(PLATFORM)

default: patch_local_workstation_install

patch_local_workstation_install: build_cross_platform override_binary

build:
	hab pkg build .

override_binary:
	ln -sf $(PWD)/bin/$(BINARY) /usr/local/bin/chef-analyze

build_cross_platform:
	hab studio run 'source .studiorc; build_cross_platform'

clean:
	@echo "Removing local artifacts... "
	rm -f bin/
	rm -f results/
	rm -f coverage/

edit:
	$(EDITOR) Makefile

view:
	$(PAGER) Makefile || cat Makefile
