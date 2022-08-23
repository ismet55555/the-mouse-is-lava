# NOTES

## General idea
A program when started will track how often mouse has been used over time


## Potential Tools

- Automation - robotgo
    - Example Project: https://github.com/ruanbekker/golang-autoclicker
- System dialogs - zenity
- System icon in notification area - systray
- System notificaiton - beeep
- Make CLI - cobra
- Working with configurations - viper
- Detach process - go-daemon


## Sub-Ideas

- Name: The mouse is lava

- Stopping it
    - Ctrl-C on binary - if not detached
    - Keyboard combination (CTRL-ALT-SHIFT-q ?? or something rare like this)
    - CLI find process and stopps it
    - the-mouse-is-lava exit

- Pause / Resume Program
    - Keyboard combination?
    - CLI can do this

- Save statistics in local file?
    - Per session?

- Sound when touching mouse? Notification? Something annoying?
- Notification has dorky insults
    - Depending how long without mouse usage - May turn into praise
    - Scale from insult to parise
    - At beginning maybe random?

- Set a no-mouse touch goals

- Disable when using certain programs
    - Drawing/CAD
    - Slack
    - etc

- When not detached
    - Output json items to be consumed by another process


## System Tray

- Controls
    - Pause/resume
    - Quit

- System icon shows mouse usage
    - Number of times - keep a count
    - Amount of time used
    - No mouse time streaks
        - Contrasted to what? Keybaords time? System up time?


## Daemonize/Detach Go Program

- NOTE: In order for all sub-commands to make sense, has to be detached!

- Either need multiple binaries or program has to call itself with a subcommand?
    - subcommand will be visible?
    - maybe --quite flag?

- How to keep reference of background process
    - Name of process
    - Number of process
    - Local file that keeps track? naw
    - Windows?

- Solution: https://socketloop.com/tutorials/golang-daemonizing-a-simple-web-server-process-example


## CLI Menu

- ROOT
    - has hidden flag to start background process
        - https://pkg.go.dev/github.com/spf13/pflag?utm_source=godoc#FlagSet.MarkHidden
        -
- start - Setup and start monitor
    - CLI option for silent, no messages in terminal
    - Attach mode? Only one session running?
- status - Get all info and status of the program
    - Get process ID and Name
    - Version number
    - Similar look as systemctl status
    - stats
    - Show location of configuration file
    - Nicely formatted to show off
    - `--json` - output stats in JSON format
- pause
- resume
- exit - Stop and quit everything


## Plan of Attack

- Annoying Messages!

# DUMP

## Process stuff
```go
currentPID := robotgo.GetPID()
pidExists, error := robotgo.PidExists(currentPID)
currentProcessName, error := robotgo.FindName(currentPID)
fmt.Println(currentPID)
fmt.Println(pidExists)
fmt.Println(currentProcessName)
processIds, error := robotgo.FindIds("lava")
if error != nil {
    panic(error)
}
fmt.Println(processIds)
os.Exit(0)
```
