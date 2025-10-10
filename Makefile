.PHONY: dev.web dev.css templates
templates:
	@templ generate
dev.web:
	@rm -rf tmp
	@air web
dev.css:
	@yarn css:watch