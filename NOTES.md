# NOTES

## General idea
A program when started will track how often mouse has been used over time


## Sub-Ideas

- Name: The mouse is lava

- System icon shows mouse usage
    - Number of times
    - Amount of time used
    - Right Clicks
    - Left Clicks
    - Distance traveled
    - No mouse time streaks
        - Contrasted to what? Keybaords time? System up time?

- Stopping it
    - Ctrl-C on binary
    - Keyboard combination (CTRL-ALT-SHIFT-q ?? or something rare like this)

- Pause / Resume Program
    - Keyboard combination?

- Save statistics in local file?
    - Per session

- Potentially a CLI tool with cobra

- Sound when touching mouse? Notification? Something annoying?
- Notification has dorky insults
    - Need some kind of "database" of these
    - Depending how long without mouse usage - May turn into praise
    - Scale from insult to parise

- Set a no-mouse touch goals

- Disable when using certain programs
    - Drawing/CAD
    - Slack
    - etc

- CLI option for silent, no messages in terminall


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
- TOML file with encouragments and insults
- Organize everything within packages
    - Research what the concensus is
    - How are other people doing it?
- Make into CLI tool wiht cobra



