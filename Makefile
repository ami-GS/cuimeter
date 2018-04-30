EXAMPLEDIRS  := $(shell find ./examples -type d)
EXAMPLESRCS := $(foreach dir, $(EXAMPLEDIRS), $(wildcard $(dir)/*.go))
EXAMPLEBINS := $(basename $(EXAMPLESRCS))

default:

example:
	@echo $(EXAMPLESRCS)
	@for gofile in $(EXAMPLESRCS); do \
	echo $$gofile; \
	go build $$gofile;\
	done

clean:
