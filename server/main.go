package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/cloudmachinery/apps/tcp-guessgame/message"
)

const (
	RandSeedEnvVar = "RAND_SEED"
)

func main() {
	cfg, err := getConfig()
	failOnError(err, "cannot get config")

	// Listen for fmt.Sprintf(":%d", cfg.Port) tcp port
	fmt.Printf("Starting server on :%d\n", cfg.Port)
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.Port)) // Запуск сервера
	failOnError(err, "Error starting server")
	defer l.Close()

	guesser := Guesser{
		Min: cfg.Min,
		Max: cfg.Max,
		Gen: rand.New(rand.NewSource(cfg.Seed)),
	}

	for {
		// accept connections
		conn, err := l.Accept()
		failOnError(err, "Error accept connection")

		go func() {
			defer conn.Close()
			err := guesser.Play(conn)
			if err != nil {
				log.Printf("error playing game: %s", err)
			}
		}()
	}
}

type Guesser struct {
	Min, Max int
	Gen      *rand.Rand
}

func (g *Guesser) Play(conn net.Conn) error {
	// read message. it should be start, otherwise return error
	messageStr, err := message.Read(conn)
	if err != nil {
		return err
	}
	if messageStr != message.Start {
		return errors.New("game not started")
	}
	// write configured min and max
	minMaxRange := fmt.Sprintf(message.MinMaxFormat, g.Min, g.Max)
	err = message.Write(conn, minMaxRange)
	if err != nil {
		return err
	}
	// make a guess and log it
	guessNum := g.Gen.Intn(g.Max-g.Min+1) + g.Min
	log.Printf("guessed num is %d", guessNum)

	for {
		// read client guess
		messageStr, err := message.Read(conn)
		if err != nil {
			return err
		}
		// convert to int
		mNum, err := strconv.Atoi(messageStr)
		if err != nil {
			return err
		}
		// use switch to compare numbers and return appropriate message
		switch {
		case mNum > guessNum:
			err = message.Write(conn, message.Lower)
			if err != nil {
				return err
			}
		case mNum < guessNum:
			err = message.Write(conn, message.Higher)
			if err != nil {
				return err
			}
		default:
			err = message.Write(conn, message.Correct)
			if err != nil {
				return err
			}
			// in case of correct guess end the game and return nil
			return nil
		}
	}
}

type config struct {
	Port int
	Min  int
	Max  int
	Seed int64
}

func getConfig() (config, error) {
	var c config
	flag.IntVar(&c.Port, "port", 8080, "port to listen on")
	flag.IntVar(&c.Min, "min", 0, "minimum number to guess")
	flag.IntVar(&c.Max, "max", 100, "maximum number to guess")
	flag.Parse()

	if c.Min >= c.Max {
		return c, fmt.Errorf("min must be less than max")
	}

	seedEnv, ok := os.LookupEnv(RandSeedEnvVar)
	if ok {
		seed, env := strconv.Atoi(seedEnv)
		if env != nil {
			return c, fmt.Errorf("invalid seed value: %s", seedEnv)
		}

		c.Seed = int64(seed)
	} else {
		c.Seed = time.Now().UnixNano()
	}

	return c, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
