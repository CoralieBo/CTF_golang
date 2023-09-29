package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

const (
	startPort = 1025
	endPort   = 8000 // Vérification des ports de 1025 à 8000
	timeout   = 2 * time.Second
)

type User struct {
	Username string `json:"User"`
}

func main() {
	ip := "10.49.122.144"
	ports := getPorts(ip)
	port := getPing(ip, ports)
	if port == 0 {
		fmt.Println("no port")
		return
	}
	user := User{
		Username: "Coralie",
	}
	postSignUp(ip, port, user)
	postCheck(ip, port, user)
	postGetUserSecret(ip, port, user)
}

func scanPort(ip string, port int, wg *sync.WaitGroup, c chan int) int {
	defer wg.Done()

	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, timeout)

	if err != nil {
		return 0
	}

	defer conn.Close()
	fmt.Printf("Port %d est ouvert\n", port)
	c <- port
	return port
}

func getPorts(ip string) []int {
	c := make(chan int, 10)
	defer close(c)
	var wg sync.WaitGroup

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go scanPort(ip, port, &wg, c)
	}

	wg.Wait()

	var openPorts []int
	for i := 0; i <= len(c); i++ {
		port := <-c
		openPorts = append(openPorts, port)
	}
	return openPorts
}

func getPing(ip string, ports []int) int {
	for i := 0; i < len(ports); i++ {
		port := ports[i]
		u := url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", ip, port),
			Path:   "/ping",
		}
		urlStr := u.String()
		resp, err := http.Get(urlStr)
		if err != nil {
			fmt.Println("not right port")
		} else {
			defer resp.Body.Close()
			// dump, err := httputil.DumpResponse(resp, true)
			// if err != nil {
			// 	fmt.Println(err)
			// } else {
			// 	fmt.Println(string(dump))
			// }
			return port
		}
	}
	return 0
}

func postSignUp(ip string, port int, user User) {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", ip, port),
		Path:   "/signup",
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON :", err)
		return
	}
	body := bytes.NewReader(userJSON)
	urlStr := u.String()
	resp, err := http.Post(urlStr, "application/json", body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(dump))
	}
}

func postCheck(ip string, port int, user User) {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", ip, port),
		Path:   "/check",
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON :", err)
		return
	}
	body := bytes.NewReader(userJSON)
	urlStr := u.String()
	resp, err := http.Post(urlStr, "application/json", body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(dump))
	}
}

func postGetUserSecret(ip string, port int, user User) {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", ip, port),
		Path:   "/getUserSecret",
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON :", err)
		return
	}
	body := bytes.NewReader(userJSON)
	urlStr := u.String()
	resp, err := http.Post(urlStr, "application/json", body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(dump))
	}
}
