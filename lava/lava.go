package lava

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/ismet55555/the-mouse-is-lava/utils"
	"github.com/spf13/viper"

	"github.com/fatih/color"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
)

type Lava struct {
	Configs Configs
}
type Configs struct {
	InitAnimation       bool
	InitGracePeriod     int
	InitPause           bool
	GracePeriodDuration int
	GracePeriod         bool
	Sensitivity         float64
}

func (lava *Lava) PrintSomething(something string) {
	fmt.Println(something)
}

// Display intro title animation
func (lava *Lava) AnimateIntroTitle() {
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
	color.HiGreen("    (Press CTRL-c to exit)")
	color.HiGreen("================================\n\n")
}

// Systray setup
func (lava *Lava) systrayOnReady() {
	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		// systray.SetTitle("LAVA[00:08:45]")
		systray.SetTitle("LAVA")
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
				// TODO: Quit and clean up
				return
			}
		}
	}()
}

// Systray cleanup
func (lava *Lava) systrayOnExit() {
	now := time.Now()
	ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
}

func (lava *Lava) Start() {
	// Configurations viper -> configs
	configs := Configs{}
	viper.Unmarshal(&configs)
	lava.Configs = configs
	log.Debugln("Loaded configurations: ", lava.Configs)

	// Open system tray
	go func() {
		systray.Run(lava.systrayOnReady, lava.systrayOnExit)
	}()

	// Animate Intro Title
	if lava.Configs.InitAnimation {
		lava.AnimateIntroTitle()
	}

	// Handeling CTRL-c keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			color.HiGreen("\n\n Program exited. No more lava.\n\n")
			os.Exit(1)
		}
	}()

	// Pre-allocate variables
	var (
		magPosSet                            = make([]int, 11)
		magPos                               = make([]int, 2)
		magPosMeans                          = make([]float64, 2)
		initTimerStart         time.Time     = time.Now()
		postTouchTimerStart    time.Time     = time.Now()
		triggered              bool          = false
		totalNoTouchDuration   time.Duration = time.Duration(0)
		totalNoTouchCount      int           = 0
		longestNoTouchDuration time.Duration = time.Duration(0)
	)

	for {
		// Get the XY mouse pixel coordinates
		xPos, yPos := robotgo.GetMousePos()

		// Get the magnitude - sqrt(x^2 + y^2)
		mag := int(math.Sqrt(math.Pow(float64(xPos), 2) + math.Pow(float64(yPos), 2)))
		magPos = append(magPos, mag)[1:]
		magPosSet = append(magPosSet, magPos[1])[1:]
		magPosMeans = append(magPosMeans, utils.ArrayMean(magPosSet))[1:]

		// Initial grace period with no mouse touch eveluation
		if lava.Configs.InitPause {
			initElapsedTime := time.Now().Sub(initTimerStart)
			log.Debug("Initial delay: ", initElapsedTime, " / ", lava.Configs.InitGracePeriod)

			if initElapsedTime > time.Duration(lava.Configs.InitGracePeriod)*time.Second {
				log.Debug("Mouse movement sensor is active ...")
				lava.Configs.InitPause = false
			}
			time.Sleep(200 * time.Millisecond)
			continue
		}

		// Check for any mouse movement at all
		if magPos[0] != magPos[1] {
			log.Debug("Mouse is moving ...")
		}

		// Check for movement trigger - Difference of averages is above threshold
		if math.Abs(magPosMeans[0]-magPosMeans[1]) > lava.Configs.Sensitivity {
			triggered = true
		}

		if triggered && !lava.Configs.GracePeriod {
			totalNoTouchCount++
			if totalNoTouchDuration > longestNoTouchDuration {
				longestNoTouchDuration = totalNoTouchDuration
			}
			message := fmt.Sprintf("No-Touch Duration: %s", utils.DurationToString(totalNoTouchDuration))
			log.Debug("Triggered - ", message)
			utils.ShowAlert("Mouse LAVA!", message)
			color.Red("Mouse LAVA! - %s", message)
			postTouchTimerStart = time.Now()
			lava.Configs.GracePeriod = true
		}

		// Elapsed time
		totalNoTouchDuration = time.Now().Sub(postTouchTimerStart)

		// Reset
		if lava.Configs.GracePeriod && totalNoTouchDuration > time.Duration(lava.Configs.GracePeriodDuration)*time.Second {
			log.Debug("Reseting total no-touch timer ...")
			lava.Configs.GracePeriod = false
		}

		triggered = false

		// Loop delay
		time.Sleep(100 * time.Millisecond)
	}
}
