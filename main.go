package main

import (
	"time"

	"github.com/gen2brain/beeep"
	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Configure logger
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	// log.SetLevel(log.DebugLevel)
	log.SetLevel(log.InfoLevel)
}

func show_alert() {
	err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
	if err != nil {
		panic(err)
	}
	// err = beeep.Notify("Notify", "Notify body", "assets/information.png")
	// if err != nil {
	// panic(err)
	// }
	err = beeep.Alert("Damn it!", "You touched the mouse!!", "assets/warning.png")
	if err != nil {
		panic(err)
	}

}

func main() {
	log.Info("START - Process ID: ", robotgo.GetPID())

	var x = make([]int, 5)
	var y = make([]int, 5)

	initGracePeriod := 2
	initPause := true

	pauseGracePeriod := 3

	mouseTouch := false
	postTouchTimerStart := time.Now()
	noMouseTouchTime := time.Duration(0)
	initTimerStart := time.Now()

	count := 0

	for {
		// Initial grace period
		if initPause {
			initTime := time.Now().Sub(initTimerStart)
			log.Debug("Initial delay: ", initTime, " / ", initGracePeriod)
			if initTime > time.Duration(initGracePeriod)*time.Second {
				initPause = false
			}
			time.Sleep(200 * time.Millisecond)
			continue
		}

		// Get the mouse coordinates
		xInst, yInst := robotgo.GetMousePos()

		// Change in x direction
		x = append(x, xInst)
		x = x[1:]
		x_change := x[len(x)-1] - x[len(x)-2]

		// Change in y direction
		y = append(y, yInst)
		y = y[1:]
		y_change := y[len(y)-1] - y[len(y)-2]

		// Detect movement
		if count > 5 && (x_change != 0 || y_change != 0) {
			if !mouseTouch {
				log.Errorf("You Touched the Mouse! (Duration: %.2f minutes)", noMouseTouchTime.Minutes())
				show_alert()
				postTouchTimerStart = time.Now()
				mouseTouch = true
			}
			log.Warning("Mouse is moving ...")
		}

		// Elapsed time
		noMouseTouchTime = time.Now().Sub(postTouchTimerStart)

		// Reset
		if mouseTouch && noMouseTouchTime > time.Duration(pauseGracePeriod)*time.Second {
			log.Debug("Reseting no touch timer ...")
			mouseTouch = false
		}

		log.Debug(xInst, yInst, mouseTouch, noMouseTouchTime)

		count++

		// Loop delay
		time.Sleep(200 * time.Millisecond)
	}
	log.Info("DONE")
}
