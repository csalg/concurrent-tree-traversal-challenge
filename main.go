package main

func main(){
	apiClient := NewApiClient()
	apiClient.GetToken()
	apiClient.GetPage(1)
}

// TODO
// Make it work with login but without network timeouts or random errors.
