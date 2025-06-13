SUBDIRS := $(wildcard */.)

$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)

dbs: $(SUBDIRS)

flyway: $(SUBDIRS)

.PHONY: dbs flyway $(SUBDIRS)
