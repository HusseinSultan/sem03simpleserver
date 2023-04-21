
package main

import (
	"io"
	"log"
	"net"
      	"sync"
        "github.com/HusseinSultan/is105sem03/mycrypt"
	"github.com/HusseinSultan/funtemps/conv"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.2:8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
      					}
				        dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))
 					switch msg := string(dekryptertMelding); msg {
  				        case "ping":
						kryptertMelding := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, 4)
                                        log.Println("Kryptert melding: ", string(kryptertMelding))
                                        _, err = conn.Write([]byte(string(kryptertMelding)))
					case "kjevik":
						celsiusStr := "20" // bare for testing, erstatt med din egen logikk for å hente celsiusverdien
						celsius, err := conv.ParseFloat(celsiusStr)
						if err != nil {
							kryptertMelding := mycrypt.Krypter([]rune("Invalid temperature value"), mycrypt.ALF_SEM03, 26)
							_, err = conn.Write([]byte(string(kryptertMelding)))
						} else {
							fahrenheit := conv.CelsiusToFahrenheit(celsius)
							response := fmt.Sprintf("%.2f degrees Fahrenheit", fahrenheit)
							kryptertMelding := mycrypt.Krypter([]rune(response), mycrypt.ALF_SEM03, len(response))
							_, err = conn.Write([]byte(string(kryptertMelding)))
						}       					default:
						_, err = c.Write(buf[:n])
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
