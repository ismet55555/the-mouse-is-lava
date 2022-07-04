# NOTES

## General idea
A program when started will track how often mouse has been used over time


## Sub-Ideas

- Name: The mouse is lava

- System icon shows mouse usage
    - Number of times - keep a cound
    - Amount of time used
    - No mouse time streaks
        - Contrasted to what? Keybaords time? System up time?

- Stopping it
    - Ctrl-C on binary
    - Keyboard combination (CTRL-ALT-SHIFT-q ?? or something rare like this)
    - CLI find process and stopps it

- Pause / Resume Program
    - Keyboard combination?
    - CLI can do this

- Save statistics in local file?
    - Per session

- Potentially a CLI tool with cobra

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

- CLI option for silent, no messages in terminall

- Detach mode? Only one session running?
    - Ability to pause, stop though CLI


## Potential Tools

- Automation - robotgo
    - Example Project: https://github.com/ruanbekker/golang-autoclicker
- System dialogs - zenity
- System icon in notification area - systray
- System notificaiton - beeep
- Make CLI - cobra
- Working with configurations - viper


## Plan of Attack

- Systray icon as a red and grey volcano
    - Depending if on or off

