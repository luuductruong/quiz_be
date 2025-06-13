SUBDIRS := $(wildcard */.)

$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)

dbs: $(SUBDIRS)

.PHONY: dbs $(SUBDIRS)
