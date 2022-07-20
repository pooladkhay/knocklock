package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

/*
	Ideas for improving the KnockLock logic:
		- Listening on more ports than the unlock sequence to make brute-force attacks harder.

		- After a user successfully unlocks the lock, they should no longer knock on any of the ports defined in the sequence.
		  Detecting and black-listing that users can help making brute-force attacks harder.

		- Logic to detect number of complete knocks and maybe limiting it.
		  If we have 4 ports in our sequence, That would be how many time a user knocked 4 times in a row.

		- Implementing some sort of timeout for each user. so they have a defined time frame to enter the knock sequence.
*/

// "knocked" represents a single request (knock)
type knocked struct {
	// IP address of the machine who knocked
	ip string
	// knocked port
	port string
}

func main() {
	// port sequence we expect
	portSeq := []string{"2002", "6006", "3003", "1001"}

	// maping IPs to the sequence of ports they knocked
	knockDb := make(map[string][]string)

	// channel to communicate data from listeners to main routine.
	// data is of type "knocked"
	comm := make(chan knocked)

	// Spawning listeners to listen on all defined ports
	for _, port := range portSeq {
		go listener(port, comm)
	}

	// waiting for knocks to arrive
	for knock := range comm {
		fmt.Printf("%s knocked on port %s\n", knock.ip, knock.port)

		knockDb[knock.ip] = append(knockDb[knock.ip], knock.port)

		// comparing knocked sequence to defined port sequence,
		// only when the number of knocked ports reaches the number of defined ports
		if len(knockDb[knock.ip]) == len(portSeq) {
			if isCorrect := compareSeq(portSeq, knockDb[knock.ip]); isCorrect {
				// logic to handle allowed IP addresses:
				fmt.Printf("%s is ALLOWED to enter :)\n", knock.ip)
				knockDb[knock.ip] = []string{}
			} else {
				// logic to handle denied IP addresses:
				fmt.Printf("%s is NOT_ALLOWED to enter :(\n", knock.ip)
				knockDb[knock.ip] = []string{}
			}
		}
	}
}

// listener listens on a specific raw TCP port.
func listener(port string, comm chan knocked) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalln(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()
	fmt.Printf("listening for knocks on port %s\n", port)

	for {
		c, err := l.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			return
		}

		// building and sending an instance of "knocked" data back to main routine
		data := knocked{
			ip:   strings.Split(c.RemoteAddr().String(), ":")[0],
			port: port,
		}
		comm <- data

		// closing the connection immediately after processing the data
		c.Close()
	}
}

// helper function for comparing two slices (both equality and sequence of items)
func compareSeq(portSeq, knockSeq []string) bool {
	for i, v := range portSeq {
		if v != knockSeq[i] {
			return false
		}
	}
	return true
}
