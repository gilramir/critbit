# Find all .dot files in the CWD
DOTFILES := $(wildcard *.dot)

# Calculate the .png filenames for all .dot files
PNGFILES := $(DOTFILES:.dot=.png)

# The default target will build PNG files for
# all .dot files.
.PHONY: all
all: $(PNGFILES)

# Convert a .dot to a .png
%.png : %.dot
	dot -Tpng -o $@ $<

# Delete the generated files
.PHONY: clean
clean:
	rm *.dot *.png
