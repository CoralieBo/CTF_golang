package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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
	endPort   = 8192 // Vérification des ports de 1025 à 8000
	timeout   = 2 * time.Second
)

type User struct {
	Username string `json:"User"`
}

type UserSecret struct {
	Username string `json:"User"`
	Secret   string `json:"Secret"`
}

type SubmitBody struct {
	Username string `json:"User"`
	Secret   string `json:"Secret"`
	Content  struct {
		Level     string `json:"Level"`
		Challenge struct {
			Username string `json:"Username"`
			Secret   string `json:"Secret"`
			Points   string `json:"Points"`
		} `json:"Challenge"`
		Protocol  string `json:"Protocol"`
		SecretKey string `json:"SecretKey"`
	} `json:"Content"`
}

func main() {
	ip := "10.49.122.144"
	ports := getPorts(ip)
	port := getPing(ip, ports)
	if port == 0 {
		log.Fatal("no port")
	}
	user := User{
		Username: "Coralie",
	}
	postSignUp(ip, port, user)
	postCheck(ip, port, user)
	userSecret := postGetUserSecret(ip, port, user)
	if userSecret.Secret == "" {
		log.Fatal("no secret")
	}
	level := getUserLevel(ip, port, userSecret)
	if level == "" {
		log.Fatal("no level")
	}
	points := getUserPoints(ip, port, userSecret)
	if points == "" {
		log.Fatal("no points")
	}
	fmt.Println("iNeedAHint....")
	for i := 0; i < 10; i++ {
		postINeedAHint(ip, port, userSecret)
	}
	protocol := 724490
	postEnterChallenge(ip, port, userSecret)
	var submit SubmitBody
	submit.Username = user.Username
	submit.Secret = userSecret.Secret
	submit.Content.Level = level
	submit.Content.Challenge.Username = user.Username
	submit.Content.Challenge.Secret = userSecret.Secret
	submit.Content.Challenge.Points = points
	submit.Content.Protocol = fmt.Sprintf("%d", protocol)
	submit.Content.SecretKey = ""
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
	fmt.Println("ping....")
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
			resp.Body.Close()
			for i := 0; i < 10; i++ {
				resp, err := http.Get(urlStr)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()
				response, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(response))
			}
			return port
		}
	}
	return 0
}

func postSignUp(ip string, port int, user User) {
	fmt.Println("signup....")
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
	fmt.Println("check....")
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

func postGetUserSecret(ip string, port int, user User) UserSecret {
	fmt.Println("getUserSecret....")
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", ip, port),
		Path:   "/getUserSecret",
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON :", err)
		return UserSecret{}
	}
	body := bytes.NewReader(userJSON)
	urlStr := u.String()
	resp, err := http.Post(urlStr, "application/json", body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Ici si on appel plusieurs fois la route elle nous renvoi un hash du user (donc j'ai direct hash)
	var userSecret UserSecret
	hasher := sha256.New()
	hasher.Write([]byte(user.Username))
	hashBytes := hasher.Sum(nil)
	hashHex := hex.EncodeToString(hashBytes)
	userSecret.Secret = hashHex
	userSecret.Username = user.Username
	return userSecret
}

func getUserLevel(ip string, port int, user UserSecret) string {
	fmt.Println("getUserLevel....")
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", ip, port),
		Path:   "/getUserLevel",
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON :", err)
		return ""
	}
	body := bytes.NewReader(userJSON)
	urlStr := u.String()
	resp, err := http.Post(urlStr, "application/json", body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	result := response[7 : len(response)-1]
	fmt.Println(string(result))
	return string(result)
}

func getUserPoints(ip string, port int, user UserSecret) string {
	fmt.Println("getUserPoints....")
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", ip, port),
		Path:   "/getUserPoints",
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON :", err)
		return ""
	}
	body := bytes.NewReader(userJSON)
	urlStr := u.String()
	resp, err := http.Post(urlStr, "application/json", body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	result := response[14+len(user.Username) : len(response)-1]
	fmt.Println(string(result))
	return string(result)
}

func postEnterChallenge(ip string, port int, user UserSecret) {
	fmt.Println("enterChallenge....")
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", ip, port),
		Path:   "/enterChallenge",
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
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(response))
	firstChallenge := response[57:89]
	dontForget := response[108 : len(response)-1]
	fmt.Println(string(firstChallenge))
	fmt.Println(string(dontForget))
}

func postINeedAHint(ip string, port int, user UserSecret) {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", ip, port),
		Path:   "/iNeedAHint",
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
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(response))
}

// func submitSolution(ip string, port int, user UserSecret) {
// 	u := url.URL{
// 		Scheme: "http",
// 		Host:   fmt.Sprintf("%s:%d", ip, port),
// 		Path:   "/submitSolution",
// 	}

// 	bodyJSON

// 	userJSON, err := json.Marshal(user)
// 	if err != nil {
// 		fmt.Println("Erreur lors de la conversion en JSON :", err)
// 		return
// 	}
// 	body := bytes.NewReader(userJSON)
// 	urlStr := u.String()
// 	resp, err := http.Post(urlStr, "application/json", body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()
// 	response, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(string(response))
// }
