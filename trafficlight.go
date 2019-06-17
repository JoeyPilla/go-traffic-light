package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
)

var tog = make(chan bool)
var walk = make(chan string)

func main() {
	fmt.Println("opening gpio")
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer rpio.Close()

	green := rpio.Pin(23)
	yellow := rpio.Pin(15)
	red := rpio.Pin(18)
	pin := rpio.Pin(24)
	white := rpio.Pin(25)
	green.Output()
	green.High()
	yellow.Output()
	yellow.High()
	red.Output()
	red.High()
	white.Output()
	white.High()
	pin.Input()

	//go toggle(&green, 1000)
	//go toggle(&yellow, 2000)
	//go toggle(&red, 3000)
	go poi(&pin)
	light := [3]rpio.Pin{green, yellow, red}
	go toggle(&light)
	go pedLight(white)
	for {
		reader := bufio.NewReader(os.Stdin)
		reader.ReadRune()
		green.High()
		yellow.High()
		red.High()
		white.High()
		break
	}
}

func poi(pin *rpio.Pin) {
	pin.Input()
	for {
		start := time.Now()
		for pin.Read() == 1 {
			if time.Since(start) > time.Second*1 {
				tog <- true
			}
		}
	}
}

func pedLight(light rpio.Pin) {
	x := ""
	for {
		select {
		case tmp := <-walk:
			x = tmp
		default:
			if x == "green" {
				light.Low()
			} else if x == "yellow" {
				light.Toggle()
				time.Sleep(time.Second / 10)
			} else {
				light.High()
			}
		}
	}
}

func toggle(light *[3]rpio.Pin) {
	sel := "green"
	for {
		switch sel {
		case "green":
			light[0].Toggle()
			walk <- sel
			time.Sleep(time.Second * 5)
			light[0].Toggle()
			sel = "yellow"
		case "yellow":
			light[1].Toggle()
			walk <- sel
			time.Sleep(time.Second * 3)
			light[1].Toggle()
			sel = "red"

		case "red":
			light[2].Toggle()
			walk <- sel
			start := time.Now()
			stop := false
			for time.Since(start) < time.Second*5 {
				select {
				case <-tog:
					fmt.Println("toggled")
					stop = true
				default:
					continue
				}
				if stop {
					break
				}
			}
			sel = "green"
			light[2].Toggle()
		}
	}
}
