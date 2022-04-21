default_target: pvscheck
INSTALL_PREFIX ?= /usr/local
#=============================================================================
.PHONY: help
help:
	@echo "  Usage: make"
	@echo "  Usage: make install"
	@echo "  Usage: make install [INSTALL_PREFIX=$(INSTALL_PREFIX)]"

.PHONY: pvscheck
pvscheck:
	@go build -ldflags "-s -w"

install: pvscheck
	@echo Installing into $(INSTALL_PREFIX)
	@install -m 755 pvscheck ${INSTALL_PREFIX}


