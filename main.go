package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sort"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Echo server listening on port %s.\n", port)

	err := http.ListenAndServe(":"+port, http.HandlerFunc(handler))
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("URL:")
	fmt.Println(r.URL)

	fmt.Println("Method:")
	fmt.Println(r.Method)

	fmt.Println("Headers:")
	headers := printHeaders(r.Header)
	fmt.Println(headers)

	body := ""
	if r.Body != nil {
		fmt.Println("Body:")

		buf := &bytes.Buffer{}

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			fmt.Printf("Body reading error: %v", err)
		}
		body = buf.String()
		fmt.Println(body)
		fmt.Println()
	}

	err := r.Body.Close()
	if err != nil {
		fmt.Printf("Unable to close body: %v", err)
	}

	f, err := os.OpenFile(os.Getenv("FILE_NAME"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	if _, err = f.WriteString(fmt.Sprintf("%s,%s,\"%s\",%s\n", r.URL, r.Method, headers, body)); err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		fmt.Printf("Unable to close the file: %v", err)
	}
}

func printHeaders(h http.Header) string {
	sortedKeys := make([]string, 0, len(h))

	for key := range h {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Strings(sortedKeys)

	headers := ""
	for _, key := range sortedKeys {
		for _, value := range h[key] {
			headers = headers + fmt.Sprintf("%s: %s,", key, value)
		}
	}

	return headers
}
