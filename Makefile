update:
	@git submodule update --init --recursive
	@echo "Submodules updated."

.PHONY: update

credits:
	@go install github.com/Songmu/gocredits/cmd/gocredits@latest
	@gocredits -skip-missing . > CREDITS

.PHONY: credits
