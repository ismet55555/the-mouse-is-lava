package lava

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/ismet55555/the-mouse-is-lava/assets/iconOFF"
	"github.com/ismet55555/the-mouse-is-lava/assets/iconON"
	"github.com/ismet55555/the-mouse-is-lava/utils"
	"github.com/spf13/viper"

	"github.com/fatih/color"
	"github.com/getlantern/systray"
	"github.com/go-vgo/robotgo"
	"github.com/sevlyar/go-daemon"
	log "github.com/sirupsen/logrus"
)

type Lava struct {
	Configs                Configs
	lavaOn                 bool
	totalNoTouchDuration   time.Duration
	totalNoTouchCount      int
	longestNoTouchDuration time.Duration
	systrayElapsedTime     *systray.MenuItem
}

type Configs struct {
	InitAnimation       bool
	InitGracePeriod     int
	InitPause           bool
	GracePeriodDuration int
	GracePeriod         bool
	Sensitivity         float64
	EnableSystray       bool
}

// Display intro title animation
func (lava *Lava) AnimateIntroTitle() {
	color.HiGreen("================================")
	titleText := []string{"    The", "       mouse", "            is"}
	for _, word := range titleText {
		color.HiGreen(word)
		time.Sleep(time.Duration(0.6 * float32(time.Second)))
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
		systray.SetTemplateIcon(iconON.Data, iconON.Data)
		systray.SetTitle("LAVA")
		systray.SetTooltip("The Mouse is LAVA")

		mChecked := systray.AddMenuItemCheckbox("Turn LAVA Off", "Turn On/Off", true)
		lava.systrayElapsedTime = systray.AddMenuItem("Off", "Mouse No-touch Time")
		lava.systrayElapsedTime.Disable()

		systray.AddSeparator()
		mQuit := systray.AddMenuItem("QUIT", "Quit the whole app")

		for {
			select {
			case <-mChecked.ClickedCh:
				if mChecked.Checked() {
					systray.SetIcon(iconOFF.Data)
					mChecked.Uncheck()
					mChecked.SetTitle("Turn LAVA On")
					lava.lavaOn = false
					log.Debug("Monitoring turned OFF via systray")
				} else {
					systray.SetIcon(iconON.Data)
					mChecked.Check()
					mChecked.SetTitle("Turn LAVA Off")
					lava.lavaOn = true
					log.Debug("Monitoring turned ON via systray")
				}
			case <-mQuit.ClickedCh:
				// TODO: Quit and clean up
				systray.Quit()
				os.Exit(0)
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

// Entrypoint to program execution
func (lava *Lava) Start() {
	// Load configurations (viper -> configs)
	configs := Configs{}
	viper.Unmarshal(&configs)
	lava.Configs = configs
	log.Debugln("Loaded configurations: ", lava.Configs)

    // Cannot have no systray and detached
	if viper.GetBool("noSystray") && viper.GetBool("detach") {
        color.HiRed("Unable to --detach with --nosystray menu")
        os.Exit(1)
    }

	// Open system tray
	if !viper.GetBool("noSystray") {
		go func() {
			systray.Run(lava.systrayOnReady, lava.systrayOnExit)
		}()
	} else {
		log.Debug("Systray is off")
	}

	// Animate Intro Title
	if !viper.GetBool("noIntro") {
        if lava.Configs.InitAnimation && !viper.GetBool("detach") {
            lava.AnimateIntroTitle()
        }
    }

	// Detach process to make daemon
    // NOTE: To force quit: kill `cat lava_mouse.pid`
	if viper.GetBool("detach") {
		log.Debug("Program detached. Running in background ... ")
		cntxt := &daemon.Context{
			PidFileName: "lava_mouse.pid",
			PidFilePerm: 0644,
			LogFileName: "lava_mouse.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
			Args:        []string{"[mouse lava detached]"},
		}

		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if d != nil {
			return
		}
		defer cntxt.Release()

		log.Print("*************************************************")
		log.Print("Mouse lava started as detached background process")
		log.Print("*************************************************")
	}

	// Initiate main loop
	lava.mainLoop()
}

// Main loop
func (lava *Lava) mainLoop() {
	lava.lavaOn = true

	// Handeling CTRL-c keyboard interrupt with go routine
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			color.HiGreen("\n\nProgram exited. No more lava.\n\n")
			os.Exit(1)
		}
	}()

	// Pre-allocate variables
	var (
		magPosSet                     = make([]int, 11)
		magPos                        = make([]int, 2)
		magPosMeans                   = make([]float64, 2)
		initTimerStart      time.Time = time.Now()
		postTouchTimerStart time.Time = time.Now()
		triggered           bool      = false
	)

	lava.totalNoTouchDuration = time.Duration(0)
	lava.totalNoTouchCount = 0
	lava.longestNoTouchDuration = time.Duration(0)

	// Main loop
	for {
		// Get the XY mouse pixel coordinates
		xPos, yPos := robotgo.GetMousePos()

		// Get the magnitude of mouse position - sqrt(x^2 + y^2)
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

		if lava.lavaOn && triggered && !lava.Configs.GracePeriod {
			lava.totalNoTouchCount++
			if lava.totalNoTouchDuration > lava.longestNoTouchDuration {
				lava.longestNoTouchDuration = lava.totalNoTouchDuration
			}
			message := fmt.Sprintf("No-Touch Duration: %s", utils.DurationToString(lava.totalNoTouchDuration, "%01d hours : %01d minutes : %01d seconds"))
			log.Debug("Triggered - ", message)
			utils.ShowAlert("Mouse LAVA!", message)
			color.Red("Mouse LAVA! - %s", message)
			postTouchTimerStart = time.Now()
			lava.Configs.GracePeriod = true
		}

		// Elapsed time
		lava.totalNoTouchDuration = time.Now().Sub(postTouchTimerStart)

		// Update systray
        if !viper.GetBool("noSystray") {
            if lava.lavaOn {
                lava.systrayElapsedTime.SetTitle(utils.DurationToString(lava.totalNoTouchDuration, "%01d hours : %01d minutes : %01d seconds"))
            } else {
                lava.systrayElapsedTime.SetTitle("Off")
            }
        }

		// Reset
		if lava.Configs.GracePeriod && lava.totalNoTouchDuration > time.Duration(lava.Configs.GracePeriodDuration)*time.Second {
			log.Debug("Reseting total no-touch timer ...")
			lava.Configs.GracePeriod = false
		}
		triggered = false

		// Loop delay
		time.Sleep(100 * time.Millisecond)
	}
}
