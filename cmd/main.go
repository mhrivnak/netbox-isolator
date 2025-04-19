package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mhrivnak/netbox-isolator/pkg/client"
	"github.com/mhrivnak/netbox-isolator/pkg/handlers"
)

func main() {
	for _, envvar := range []string{"NETBOX_URL", "NETBOX_TOKEN"} {
		if os.Getenv(envvar) == "" {
			fmt.Printf("%s environment variable is not set\n", envvar)
			os.Exit(1)
		}
	}

	url := os.Getenv("NETBOX_URL")
	token := os.Getenv("NETBOX_TOKEN")

	c, err := client.New(url, token)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	h := handlers.New(c)

	http.HandleFunc("/api/devices/", h.Device)
	fmt.Println("Listening on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
