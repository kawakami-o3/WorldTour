TOOLS_DIR := tools
BIN_DIR := bin

# tools/ 以下の main.go を持つサブディレクトリを自動検出
TOOLS := $(notdir $(patsubst %/main.go,%,$(wildcard $(TOOLS_DIR)/*/main.go)))

.PHONY: all clean $(TOOLS)

all: $(TOOLS)

$(TOOLS):
	cd $(TOOLS_DIR) && go build -o ../$(BIN_DIR)/$@ ./$@/

clean:
	rm -f $(addprefix $(BIN_DIR)/,$(TOOLS))
