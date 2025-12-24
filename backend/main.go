// @title           Kyora API
// @version         1.0
// @description     Kyora backend HTTP API.
// @BasePath        /
//
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Access JWT in the format: Bearer <token>
package main

import "github.com/abdelrahman146/kyora/cmd"

func main() {
	cmd.Execute()
}
