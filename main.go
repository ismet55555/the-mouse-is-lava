package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/go-vgo/robotgo"
	"github.com/ismet55555/mouse-lava/icon"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Configure logger
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)

	// Open system tray
	go func() {
		systray.Run(onReady, onExit)
	}()

	if runtime.GOOS == "windows" {
		color.HiYellow("Sorry. Windows is not fully supported yet :/")
		os.Exit(1)
	}
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

// Display text at a time interval
func displayTextWithTime(interval float32, text []string) {
	for _, word := range text {
		fmt.Println(word)
		time.Sleep(time.Duration(interval * float32(time.Second)))
	}
}

// Display intro title animation
func showIntroTitle() {
	color.HiGreen("================================")
	titleText := []string{"    The", "       mouse", "            is"}
	for _, word := range titleText {
		color.HiGreen(word)
		time.Sleep(time.Duration(1.0 * float32(time.Second)))
	}
	color.HiRed("               ╦  ┌─┐┬  ┬┌─┐")
	color.HiRed("               ║  ├─┤└┐┌┘├─┤")
	color.HiRed("               ╩═╝┴ ┴ └┘ ┴ ┴")
	color.HiGreen("================================")
	color.HiGreen("    (Press CTRL-c to exit)\n\n")
	color.HiGreen("================================")
}

func main() {
	// Handeling CTRL-c keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			color.HiRed("\n\n Program exited. No more lava.\n\n")
			os.Exit(1)
		}
	}()

	showIntroTitle()
	log.Debug("START - Process ID: ", robotgo.GetPID())

	// User defined settings
	var (
		initGracePeriod     int     = 2
		initPause           bool    = true
		gracePeriodDuration int     = 3
		gracePeriod         bool    = false
		sensitivity         float64 = 8.0
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
			initElapsedTime := time.Now().Sub(initTimerStart)
			log.Debug("Initial delay: ", initElapsedTime, " / ", initGracePeriod)

			if initElapsedTime > time.Duration(initGracePeriod)*time.Second {
				log.Debug("Mouse movement sensor is active ...")
				initPause = false
			}
			time.Sleep(200 * time.Millisecond)
			continue
		}

		// Check for any mouse movement at all
		if magPos[0] != magPos[1] {
			log.Debug("Mouse is moving ...")
		}

		// Check for movement trigger - Difference of averages is above threshold
		if math.Abs(magPosMeans[0]-magPosMeans[1]) > sensitivity {
			triggered = true
		}

		if triggered && !gracePeriod {
			message := fmt.Sprintf("No-Touch Duration: %s", durationToString(totalNoTouchDuration))
			log.Debug("Triggered - ", message)
			show_alert("You moved the mouse!", message)
			color.Red("Oh no! Lava!")
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
}

// System tray menu
func onExit() {
	fmt.Println("Clean up")
	now := time.Now()
	ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
}
func onReady() {
	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("LAVA[00:08:45]")
		systray.SetTooltip("The Mouse is LAVA")

		mChecked := systray.AddMenuItemCheckbox("Turn LAVA Off", "Turn On/Off", true)
		mElapsedTime := systray.AddMenuItem("0 hours 8 minutes 45 seconds", "No-touch Time")
		mElapsedTime.Disable()

		systray.AddSeparator()
		mQuit := systray.AddMenuItem("QUIT", "Quit the whole app")

		for {
			select {
			case <-mChecked.ClickedCh:
				if mChecked.Checked() {
					mChecked.Uncheck()
					mChecked.SetTitle("Turn LAVA On")
				} else {
					mChecked.Check()
					mChecked.SetTitle("Turn LAVA Off")
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
				// TODO: Quite and clean up
				return
			}
		}
	}()
}
