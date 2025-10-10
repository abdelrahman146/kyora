.PHONY: dev templates
templates:
	@templ generate
dev:
	@rm -rf tmp
	@yarn css:watch & air web