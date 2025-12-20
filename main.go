package main

import (
    "github.com/fardinabir/go-svc-boilerplate/cmd"
)

// @title         Users Service API
// @version       v1.0
// @description   User CRUD API for user-management service.
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
// @host          localhost:8082
// @BasePath      /api/v1
// @schemes       http
func main() {
	cmd.Execute()
}
