package main

import "helpdesk/internals/services/api"

func main() {
	apiService := api.New("3001", true)

	apiService.Run()
}
