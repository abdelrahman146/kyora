.PHONY: server.dev
server.dev:
	cd server && air serve
web.dev:
	cd web && yarn dev