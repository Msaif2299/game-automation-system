# Automation System for AFK Grinding (AQW Bot) [UNDER DEVELOPMENT! STATUS: EARLY STAGES]

I am building this tool for my own entertainment purposes. There are better alternatives, but I am building this to improve my ability to code and learn new things, while also taking care of my paranoia of losing private data if I use other bots. This bot does not and will never use any data nor transmit it anywhere or to anyone, unless its been compromised.

## What does it provide?

1. Allows for farming of items in the background via scripts that contain clicks.
2. Designed for the PC client.
3. Provides the ability to write your own scripts using commands.

## Restrictions

1. Only designed for Windows.
2. Tested on Windows 10. May not work on other platforms.
3. Under development.
4. Designed for educational use. Use at your own discretion. I do not promote the usage of bots, and highly advise against using it, as it may lead to unforeseen consequences.

## How to use?

.exe file has not been uploaded yet. But to create it, clone the repository. Download wails. Run "wails build" in the root folder of the code. The output of the command will show where the .exe file is generated and stored.
Copy paste the data folder into the folder where the .exe file is created. Make a shortcut and paste it wherever you want.

## Notes

1. Under development, needs a lot of polishing and testing.
1. Does not require the usage of internet. No data is used or transmitted.
1. Copy paste the "data" folder into the same folder as the .exe file and make sure the program has read and write permissions in the data folder.
1. Put your custom scripts in the data folder, within any of the two script folders.
1. Still trying to figure out how to scale the x and y coordinates to the size of screen.
1. UI is still messed up. Trying to fix the dropdown issue at the moment.

## How to write a script

Scripts are just text files with commands. The following commands are available, and will be updated according to the development.
Script names should use underscore '\_' instead of spaces. Dropdown will appear with underscores replaces with spaces.
For example, 'general_attack.txt' will appear as 'General Attack'.
In the following commands, text inside the bracket is the description of the command. Do not include it in the actual text file.

### Commands for scripts:

- CLICK x y (Purpose: To click at specific coordinates. Replace x and y with the x and y coordinates of the click.)
- JOIN citadel-99999999 (Purpose: To join a map. Replace citadel-9999999 with the map name and instance number.)
- ATTACK 10 (Purpose: Class skills are clicked. 10 denotes the number of skill clicks. Can be changed to whatever number.)
- REST (Purpose: To press the rest button.)
- DELAY x (Purpose: To make the script wait for x amount of duration in seconds. x is a decimal value. It may or may not have a decimal point. Replace x with whatever number is required. Eg.: DELAY 0.3 will pause the script for 0.3 seconds. DELAY 5 will pause the script for 5 seconds.)
- POTION (Purpose: To press the potion button. Make sure the potion is equipped.)
- QUEST_TURNIN 2 (Purpose: To turn in the quest in the quest screen. 2 is the quest position, i.e., the second quest in the quest box. Currently only supports 1st and 2nd quests.)
- QUEST_TURNIN 2 CLICK x y (Purpose: Same as previous, but CLICK x y is used to define the position at which the quest reward is located, if it requires clicking.)

### Commands for class scripts:

- 0 (Purpose: Click autoattack skill)
- AUTOATTACK (Purpose: Same as above. Added to improve readability.)
- 1 (Purpose: Click the equipped class' 1st skill.)
- 2 (Purpose: Click the equipped class' 2nd skill.)
- 3 (Purpose: Click the equipped class' 3rd skill.)
- 4 (Purpose: Click the equipped class' 4th skill.)

## Things to do:

1. Scale the clicks to the window size.
2. Add more commands.
3. Add delay handling, especially after REST command.
4. Add more tests.
5. Add comments and clean up the code.
6. Fix the UI.
7. Add reset script button (Resets the script to index 0).
8. Add a drawer section with other apps such as ClickLocater, MapMaker, ScriptMaker.
9. Create those apps.
10. Add map handling functionality so its easier to use.
11. Add loops, small memory handling and subroutines in script and ScriptMaker.
12. Extend this to multi-account handling for soloing ultra bosses with 4 accounts via coordination.
13. Make this cross platform?
14. This is too much.
