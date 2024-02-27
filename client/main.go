package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/cloudmachinery/apps/tcp-guessgame/message"
)

func main() {
	addr, err := getServerAddr()
	failOnError(err, "cannot resolve addr")

	// DialTCP
	conn, err := net.DialTCP("tcp", nil, addr)
	failOnError(err, "can't connect to server")

	// Send "start" message and get the range
	err = message.Write(conn, message.Start)
	failOnError(err, "can't send start message")

	// fmt.Printf("min: %d, max: %d\n", min, max)
	m, err := message.Read(conn)
	failOnError(err, "can't read message")

	var min, max int
	n, err := fmt.Sscanf(m, message.MinMaxFormat, &min, &max)
	if err != nil {
		failOnError(err, "can't parse min max range")
	}
	if n != 2 {
		failOnError(err, "wrong number of elements")
	}

	var guessed bool
	for !guessed {
		//     use binary search to guess the number
		//     guess in the middle of the range
		mid := min + (max-min)/2
		log.Printf("guessing %d\n", mid)
		//     send the guess to the server
		err = message.Write(conn, strconv.Itoa(mid))
		failOnError(err, "can't send guess message")
		//     read the response from the server and log it
		m, err = message.Read(conn)
		failOnError(err, "can't read guess message")
		log.Println(m)
		//     adjust the range or exit if guess was correct
		switch m {
		case message.Higher:
			min = mid + 1
		case message.Lower:
			max = mid - 1
		default:
			log.Printf("number guessed %d", mid)
			guessed = true
		}
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getServerAddr() (*net.TCPAddr, error) {
	var (
		host string
		port int
	)
	flag.StringVar(&host, "host", "localhost", "host to listen on")
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()

	return net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
}
