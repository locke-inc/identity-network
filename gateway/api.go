package gateway

import "fmt"

// Basically, the ideal way to consume the Locke API would be to call:
// locke.get("desired_data_key")
// The package will already know who's authenticated and get the data from them
func get(key string) {
	fmt.Print(key)
}

func set(data map[string]string) {
	fmt.Println(data)
}
