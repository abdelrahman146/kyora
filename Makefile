.PHONY: dev templates
templates:
	@templ generate
dev:
	@rm -rf tmp
	@yarn watch:css & air web