SUBDIRS := $(foreach dir,$(wildcard */),$(if $(wildcard $(dir)/Makefile),$(dir)))

all: $(SUBDIRS)
	sh ./gendoc.sh

$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all $(SUBDIRS)
