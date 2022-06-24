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

- Save statistics in local file?
    - Per session

- Potentially a CLI tool with cobra
-
- Sound when touching mouse? Notification? Something annoying?


## Potential Tools

- Automation - robotgo
    - Example Project: https://github.com/ruanbekker/golang-autoclicker
- System dialogs - zenity
- System icon in notification area - systray
- System notificaiton - beeep
- Make CLI - cobra


## Plan of Attack

- Create a module - DONE
- Create main.py with simple print statment - DONE
- Install robotogo stuff - DONE
- Try sample code for mouse stuff online - DONE
- Adjust sample code neatly - DONE
- Research go routines, can we apply
- Find difference from current to previouis iteration (movement)



