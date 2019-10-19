package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-cmd/cmd"
	"github.com/gorilla/websocket"
)

type msg struct {
	Status string
}

type IperfInfo struct {
	Speed float64
	Unit  string
}

var chSpeed chan IperfInfo
var server string

func init() {
	if len(os.Args) == 1 {
		log.Println("Pass the SERVER IP")
		server = "localhost"

	} else {
		server = os.Args[1]
	}

}

func main() {
	chSpeed = make(chan IperfInfo, 2)
	// 	ch <- v    // Send v to channel ch.
	// v := <-ch  // Receive from ch, and
	// assign value to v.
	http.Handle("/", http.FileServer(http.Dir(".")))
	// http.HandleFunc("/jquery.min.js", jsHandler)
	// http.HandleFunc("/guage.min.js", gsHandler)
	http.HandleFunc("/ws", wsHandler)
	// http.HandleFunc("/", rootHandler)

	panic(http.ListenAndServe(":8080", nil))
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requested path : ", r.RequestURI)
	http.ServeFile(w, r, "jquery.min.js")
}
func gsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requested path : ", r.RequestURI)
	http.ServeFile(w, r, "gauge.min.js")
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Println("Could not open file.", err)
	}
	fmt.Fprintf(w, "%s", content)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	// go echo(conn)

	go updateSpeed(conn, chSpeed)
}

func updateSpeed(conn *websocket.Conn, ch chan IperfInfo) {
	// var info IperfInfo
	go doIperf(ch)

	var err error
	for {
		// m := msg{}
		// err = conn.ReadJSON(&m)
		// if err != nil {
		// 	fmt.Println("Error reading json.", err)
		// 	break
		// } else {
		// 	if m.Status != "" {
		// 		fmt.Printf("Got message: %#v\n", m)
		// 		go doIperf(chSpeed)
		// 	}
		// }

		xinfo := <-ch
		log.Println("GO CH READ Success : ", xinfo.Speed, xinfo.Unit)
		if err = conn.WriteJSON(xinfo); err != nil {
			fmt.Println(err)
		}

		// go func() {
		// 	//Option 1
		// 	select {
		// 	case xinfo := <-ch:
		// 		log.Println("GO CH READ ", xinfo.Speed, xinfo.Unit)
		// 		if err = conn.WriteJSON(xinfo); err != nil {
		// 			fmt.Println(err)
		// 		}
		// 	default:

		// 	}

		// }()

		// Option 2 (Random )
		// time.Sleep(1000 * time.Millisecond)
		// info.Speed = 10000 + 5000*rand.Float64()
		// info.Unit = "Mbps"

		// if err = conn.WriteJSON(info); err != nil {
		// 	fmt.Println(err)
		// }
		// time.Sleep(1 * time.Second)

	}
}

func echo(conn *websocket.Conn) {
	for {
		m := msg{}

		err := conn.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error reading json.", err)
		}

		fmt.Printf("Got message: %#v\n", m)

		if err = conn.WriteJSON(m); err != nil {
			fmt.Println(err)
		}
	}
}

func doIperf(chan<- IperfInfo) {

	// Disable output buffering, enable streaming
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}

	args := []string{"/usr/local/bin/iperf3", "-c", server, "-b0M", "-fm", "-P10", "-i .5", "-t 20", "--forceflush"}
	fmt.Println("Executing .. : ", strings.Join(args, " "))
	envCmd := cmd.NewCmdOptions(cmdOptions, "/usr/local/bin/iperf3", args...)
	// Create Cmd with options
	// envCmd := cmd.NewCmdOptions(cmdOptions, "env")

	// Print STDOUT and STDERR lines streaming from Cmd
	count := 0
	go func() {
		for {
			select {
			case line := <-envCmd.Stdout:
				// fmt.Println("line ", line)
				// if strings.Contains(line, "[SUM]") && strings.Contains(line, "receiver") {
				if strings.Contains(line, "[") {

					xx := strings.Split(line, "  ")
					if len(xx) > 12 {

						if xx[0] == "[SUM]" {
							// fmt.Println("**")

							speedinfo := strings.Split(xx[4], " ")
							var spinfo IperfInfo
							var er error
							spinfo.Speed, er = strconv.ParseFloat(speedinfo[0], 64)
							if er != nil {
								log.Println("Error parsing speed", er)
							}
							spinfo.Unit = speedinfo[1]
							fmt.Println("SPINFO ", spinfo)

							select {
							case chSpeed <- spinfo:
								fmt.Printf("\nGO CH Write :Success %d  : SPEED (%v): %v \n", count, xx[13], xx[4])
								count++
							default:
								fmt.Printf("\nGO CH Missed :%d  : SPEED (%v): %v \n", count, xx[13], xx[4])
								count++
							}

							// for i, val := range xx {
							// 	fmt.Printf("\n %d  :  %v\n", i, val)

							// }

						} else {
							// if len(xx) == 15 {
							// 	fmt.Printf("\r ==> %v \t \t ", xx[5])
							// }
							// if len(xx) == 14 {
							// 	fmt.Printf("\r ==> %v \t \t ", xx[4])
							// }
							// for i, val := range xx {
							// 	fmt.Printf("\n %d  :  %v", i, val)

							// }
							// fmt.Printf("\n ** ==> %#v", xx[4])
						}

					}
				}

				// }
			case line := <-envCmd.Stderr:
				fmt.Fprintln(os.Stderr, "Error", line)
			}
		}
	}()

	// Run and wait for Cmd to return, discard Status
	<-envCmd.Start()

	// // Cmd has finished but wait for goroutine to print all lines
	// for len(envCmd.Stdout) > 0 || len(envCmd.Stderr) > 0 {
	// 	time.Sleep(10 * time.Millisecond)
	// }
}
