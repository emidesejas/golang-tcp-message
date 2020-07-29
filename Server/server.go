package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// LoggedUser is the user with name and connection
type LoggedUser struct {
	username   string
	connection net.Conn
}

var connections []LoggedUser

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please, provide a port number")
		return
	}

	PORT := arguments[1]
	beginListen(":" + PORT)
}

func beginListen(port string) {
	// user ListenTCP
	listener, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		connection, err := listener.Accept()
		go serveConnection(connection, err)
	}
}

func serveConnection(connection net.Conn, err error) {
	if err != nil {
		fmt.Println("Incoming connection failed with ", connection.RemoteAddr())
		fmt.Println(err)
		return
	}

	var currentUser string

	defer connection.Close()

	for {
		netData, err := readConnection(connection)
		if err != nil {
			return
		}

		switch strings.TrimSpace(string(netData)) {
		case "Login":
			connection.Write([]byte("Give Me USER:\n"))

			currentUser, err = readConnection(connection)
			if err != nil {
				return
			}

			switch strings.TrimSpace(string(currentUser)) {
			case "Emiliano":
				connection.Write([]byte("Hi Emi\n"))
			case "Santiago":
				connection.Write([]byte("Hi Santi\n"))
			default:
				connection.Write([]byte("User not found\n"))
			}

			connections = append(connections, LoggedUser{currentUser, connection})
		case "Send Message":
			connection.Write([]byte("To who:\n"))
			to, err := readConnection(connection)
			if err != nil {
				return
			}

			connection.Write([]byte("What message:\n"))
			message, err := readConnection(connection)
			if err != nil {
				return
			}

			sendMessage(to, message, connection)
		case "Exit":
			connection.Write([]byte("Bye bye\n"))
			return
		}

		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		fmt.Print("-> ", string(netData))
		// t := time.Now()
		// myTime := t.Format(time.RFC3339) + "\n"
		// connection.Write([]byte(myTime))
	}
}

func readConnection(connection net.Conn) (string, error) {
	return bufio.NewReader(connection).ReadString('\n')
}

func sendMessage(user string, message string, connection net.Conn) {
	for _, loggedUser := range connections {
		if loggedUser.username == user {
			loggedUser.connection.Write([]byte(message + "\n"))
			return
		}
	}

	connection.Write([]byte("User not found\n"))
}
