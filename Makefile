.PHONY: extract-data build-extractor clean-extracted

GAME_PATH ?= $(HOME)/.steam/debian-installation/steamapps/common/Path of Exile 2
EXTRACTED_DIR ?= $(PWD)/extracted

build-extractor:
	docker build --no-cache -t poe-export .

extract-data: build-extractor
	mkdir -p $(EXTRACTED_DIR)
	@echo "Scanning game files..."
	docker run --rm -v "$(GAME_PATH)":/game poe-export bun_extract_file list-files /game > "$(EXTRACTED_DIR)/all_files.txt"
	@echo "Generating extraction list..."
	grep -E "^data/.*\.datc64$$|^metadata/" "$(EXTRACTED_DIR)/all_files.txt" > "$(EXTRACTED_DIR)/files_to_extract_dynamic.txt"
	@echo "Extracting files..."
	cat "$(EXTRACTED_DIR)/files_to_extract_dynamic.txt" | docker run --rm -i \
		-v "$(GAME_PATH)":/game \
		-v "$(EXTRACTED_DIR)":/output \
		poe-export bun_extract_file extract-files /game /output
	@echo "Fixing permissions..."
	docker run --rm -v "$(EXTRACTED_DIR)":/output poe-export chown -R $(shell id -u):$(shell id -g) /output
	rm "$(EXTRACTED_DIR)/all_files.txt" "$(EXTRACTED_DIR)/files_to_extract_dynamic.txt"
	@echo "Extraction complete."

clean-extracted:
	rm -rf $(EXTRACTED_DIR)
