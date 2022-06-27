package main

import (
	"fmt"
	"math"
	"sort"
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
	log.SetLevel(log.InfoLevel)
}

// Alert the user
func show_alert(title string, message string) {
	err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
	if err != nil {
		panic(err)
	}
	err = beeep.Alert(title, message, "none")
	if err != nil {
		panic(err)
	}

}

// Calculating the mean of a array
func arrayMean(in []int) float64 {
	var sum float64 = 0.0
	for _, item := range in {
		sum += float64(item)
	}
	return sum / float64(len(in))
}

// Calculate the median of an array
func arrayMedian(in []int) int {
	sort.Ints(in)
	indexHalf := len(in) / 2
	if len(in)%2 == 0 {
		return (in[indexHalf-1] + in[indexHalf]) / 2
	} else {
		return in[indexHalf]
	}
}

// Check if all items in slice are the same
func arrayAllItemsEqual(in []int) bool {
	first := in[0]
	for _, item := range in {
		if item != first {
			return false
		}
	}
	return true
}

// Duration to time string
func durationToString(duration time.Duration) string {
	hrs := duration / time.Hour
	duration -= hrs * time.Hour
	mins := duration / time.Minute
	duration -= mins * time.Minute
	secs := duration / time.Second
	return fmt.Sprintf("%01d hours : %01d minutes : %01d seconds", hrs, mins, secs)
}

func main() {
	fmt.Println("╔╦╗┬ ┬┌─┐  ╔╦╗┌─┐┬ ┬┌─┐┌─┐  ┬┌─┐  ╦  ┌─┐┬  ┬┌─┐")
	fmt.Println(" ║ ├─┤├┤   ║║║│ ││ │└─┐├┤   │└─┐  ║  ├─┤└┐┌┘├─┤")
	fmt.Println(" ╩ ┴ ┴└─┘  ╩ ╩└─┘└─┘└─┘└─┘  ┴└─┘  ╩═╝┴ ┴ └┘ ┴ ┴")

	log.Debug("START - Process ID: ", robotgo.GetPID())

	// User defined settings
	var (
		initGracePeriod     int     = 2
		initPause           bool    = true
		gracePeriodDuration int     = 3
		gracePeriod         bool    = false
		sensitivity         float64 = 10.0
	)

	var (
		magPosSet                          = make([]int, 11)
		magPos                             = make([]int, 2)
		magPosMeans                        = make([]float64, 2)
		initTimerStart       time.Time     = time.Now()
		postTouchTimerStart  time.Time     = time.Now()
		triggered            bool          = false
		totalNoTouchDuration time.Duration = time.Duration(0)
	)

	for {
		// Get the XY mouse pixel coordinates
		xPos, yPos := robotgo.GetMousePos()

		// Get the magnitude - sqrt(x^2 + y^2)
		mag := int(math.Sqrt(math.Pow(float64(xPos), 2) + math.Pow(float64(yPos), 2)))
		magPos = append(magPos, mag)[1:]
		magPosSet = append(magPosSet, magPos[1])[1:]
		magPosMeans = append(magPosMeans, arrayMean(magPosSet))[1:]

		// Initial grace period with no mouse touch eveluation
		if initPause {
			initTime := time.Now().Sub(initTimerStart)
			log.Debug("Initial delay: ", initTime, " / ", initGracePeriod)
			if initTime > time.Duration(initGracePeriod)*time.Second {
				log.Info("Sensor active")
				initPause = false
			}
			time.Sleep(200 * time.Millisecond)
			continue
		}

		// Check for any mouse movement at all
		if magPos[0] != magPos[1] {
			log.Warning("Mouse is moving ...")
		}

		// Check for movement trigger - Difference of averages is above threshold
		if math.Abs(magPosMeans[0]-magPosMeans[1]) > sensitivity {
			triggered = true
		}

		if triggered && !gracePeriod {
			message := fmt.Sprintf("No-Touch Duration: %s", durationToString(totalNoTouchDuration))
			log.Errorln("Triggered - ", message)
			show_alert("You moved the mouse!", message)
			postTouchTimerStart = time.Now()
			gracePeriod = true
		}

		// Elapsed time
		totalNoTouchDuration = time.Now().Sub(postTouchTimerStart)

		// Reset
		if gracePeriod && totalNoTouchDuration > time.Duration(gracePeriodDuration)*time.Second {
			log.Debug("Reseting total no-touch timer ...")
			gracePeriod = false
		}

		triggered = false

		// Loop delay
		time.Sleep(100 * time.Millisecond)
	}
	log.Info("DONE")
}
