.PHONY: all clean build

.DELETE_ON_ERROR:


protocols_files := $(shell find ./wayland/generate)

xml_protocols := $(shell find ./wayland/generate -name "*.xml")

generated_protocols := $(patsubst ./wayland/generate/resources/%,./wayland/protocols/%.go,$(xml_protocols))

generated_helpers := $(patsubst ./wayland/generate/resources/%,./wayland/%.helper.go,$(xml_protocols))

# TODO add term.everything to build
build: $(generated_protocols) $(generated_helpers)

# grouped target to generate all protocols and helpers in one go
# the & is what does this
$(generated_protocols) $(generated_helpers)&: $(protocols_files) ./wayland/generate.go
	go generate ./wayland

# term.everything: $(server_files) ./main.go wayland/wayland.xml.go
# 	go build -o term.everything .

clean:
	rm ./wayland/protocols/*.xml.go || true
	rm ./wayland/*.helper.go || true