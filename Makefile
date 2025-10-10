.PHONY: dev templates
templates:
	@templ generate
dev:
	@rm -rf tmp
	air web