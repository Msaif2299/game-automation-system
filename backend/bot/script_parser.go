package bot

import (
	"fmt"
	"strconv"
	"strings"
)

type CurrentScript struct {
	ID                         string
	currentIndex               int
	script                     []*Message
	classAttackIndices         []int
	classScript                *classScript
	questEntryClickPositionMap [][]int32
}

func NewCurrentScript(ID string) *CurrentScript {
	return &CurrentScript{
		ID:                 ID,
		currentIndex:       0,
		script:             []*Message{},
		classAttackIndices: []int{},
		classScript:        nil,
		questEntryClickPositionMap: [][]int32{
			{245, 283},
			{204, 333},
		},
	}
}

// LoadClassScript loads the class into the script, or rewrites the old class used in the script, make sure to call this
// function before calling LoadFromFile for a new CurrentScript object
func (c *CurrentScript) LoadClassScript(cScript *classScript) {
	c.classScript = cScript
	// replace all class specific entries in the script with the new class entries
	if len(c.classAttackIndices) > 0 {
		for _, idx := range c.classAttackIndices {
			c.script[idx] = c.classScript.Next()
		}
	}
}

// LoadFromFile loads and parses the script from the file in the folder /data/scripts. Make sure to call LoadClassScript before
// calling this function if a new CurrentScript object is created
func (c *CurrentScript) LoadFromFile(fname string) {
	commands := readFile(fname)
	if len(commands) == 0 {
		return
	}
	if c.classScript == nil {
		fmt.Printf("[ERROR] No class script detected, cannot load script properly, exiting")
		return
	}
	c.script = []*Message{}
	c.classAttackIndices = []int{}
	// iterate each line and parse the commands
	for _, command := range commands {
		commandSplit := strings.Split(command, " ")
		instruction := commandSplit[0]
		events := []*Message{}
		var err error
		switch instruction {
		case "ATTACK":
			events, err = c.parseAttack(commandSplit[1:])
			// add the class specific events in the tracker, so we can swap them out when new class is loaded
			if len(events) > 0 && err != nil {
				for idx := len(c.script); idx < len(c.script)+len(events); idx++ {
					c.classAttackIndices = append(c.classAttackIndices, idx)
				}
			}
		case "CLICK":
			events, err = c.parseClick(commandSplit[1:])
		case "JOIN":
			events, err = c.parseJoinMap(commandSplit[1:])
		case "REST":
			events, err = c.parseRest()
		case "POTION":
			events, err = c.parsePotion()
		case "QUEST_TURNIN":
			events, err = c.parseQuestTurnIn(command)
		}

		if err != nil {
			fmt.Printf("[ERROR] Error encountered while parsing %s in %s, error: %s, skipping command", instruction, fname, err.Error())
			continue
		}
		if len(events) > 0 {
			c.script = append(c.script, events...)
		}
	}
}

// Next fetches the next message event in the script
func (c *CurrentScript) Next() *Message {
	if len(c.script) == 0 {
		fmt.Println("[WARN] Attempting to call Next() in CurrentScript when script is empty")
		return nil
	}
	c.currentIndex++
	if c.currentIndex == len(c.script) {
		c.currentIndex = 0
	}
	return c.script[c.currentIndex]
}

// Reset resets the current script
func (c *CurrentScript) Reset() {
	c.currentIndex = 0
}

// NextWaitTimeInSeconds tells how much time to sleep before fetching the next event via Next() call
func (c *CurrentScript) NextWaitTimeInSeconds() float32 {
	if c.script[c.currentIndex].delayInSeconds != 0 {
		return c.script[c.currentIndex].delayInSeconds
	}
	switch c.script[c.currentIndex].msg_type {
	case Click:
		return 3
	case KeyPress:
		return 0.1
	}
	return 0.1
}

// parseAttack adds skill button clicks according to the loaded class in classScript variable
func (c *CurrentScript) parseAttack(params []string) ([]*Message, error) {
	if len(params) != 1 {
		return []*Message{}, fmt.Errorf("too many or few params for ATTACK command, found: %d, need only 1", len(params))
	}
	count, err := strconv.Atoi(params[0])
	if err != nil {
		return []*Message{}, fmt.Errorf("unable to convert number of attacks to int, found: %s, need number", params[0])
	}
	events := make([]*Message, count)
	for i := 0; i < count; i++ {
		events[i] = c.classScript.Next()
	}
	return events, nil
}

// parseClick command parses a mouse click event, currently does not handle delays
func (c *CurrentScript) parseClick(params []string) ([]*Message, error) {
	if len(params) != 2 {
		return []*Message{}, fmt.Errorf("too many or few params for CLICK command, found: %d, need only 2, that is x and y", len(params))
	}
	cleanXStr := cleanStringOfSpaces(params[0])
	cleanYStr := cleanStringOfSpaces(params[1])
	if !validateStrForClickEvent(cleanXStr) || !validateStrForClickEvent(cleanYStr) {
		return []*Message{}, fmt.Errorf("invalid coordinates for CLICK command, found: %s and %s", cleanXStr, cleanYStr)
	}
	x, _ := strconv.Atoi(cleanXStr)
	y, _ := strconv.Atoi(cleanYStr)
	message, err := NewMessage(Click, int32(x), int32(y), '0')
	if err != nil {
		return []*Message{}, fmt.Errorf("unable to create message while parsing CLICK, error: %s", err.Error())
	}
	return []*Message{
		message,
	}, nil
}

// parseJoinMap parses the command to join a map, given in game as /join citadel-999999 becomes "JOIN citadel-999999"
func (c *CurrentScript) parseJoinMap(params []string) ([]*Message, error) {
	if len(params) != 1 {
		return []*Message{}, fmt.Errorf("too many or few params for JOIN command, found: %d, need only 1", len(params))
	}
	keyPressEvents := make([]*Message, len(params[0]))
	var err error
	for i, key := range params[0] {
		keyPressEvents[i], err = NewMessage(KeyPress, 0, 0, key)
		if err != nil {
			return []*Message{}, fmt.Errorf("unable to create message while parsing JOIN, error: %s", err.Error())
		}
	}
	enterMsg, err := NewMessage(KeyPress, 0, 0, '#') // "enter" key command
	if err != nil {
		return []*Message{}, fmt.Errorf("unable to create message while parsing JOIN, error: %s", err.Error())
	}
	keyPressEvents = append(keyPressEvents, enterMsg)
	return keyPressEvents, nil
}

// parseRest parses the clicking the rest button command
func (c *CurrentScript) parseRest() ([]*Message, error) {
	return c.parseSingleClickNoParamCmd(1275, 963)
}

// parsePotion parses the clicking the potion button command
func (c *CurrentScript) parsePotion() ([]*Message, error) {
	return c.parseSingleClickNoParamCmd(1162, 945)
}

// parseSingleClickNoParamCmd is a template function for commands that take no parameters, and only return 1 click event
func (c *CurrentScript) parseSingleClickNoParamCmd(x int32, y int32) ([]*Message, error) {
	msg, err := NewMessage(Click, x, y, '0')
	if msg == nil {
		return []*Message{}, err
	}
	return []*Message{msg}, nil
}

// parseQuestTurnInClickGenerator generates click events from a list of (x, y, delay) click descriptions. Delay is the time to sleep
// before the event is sent to the game.
func (c *CurrentScript) parseQuestTurnInClickGenerator(clicks [][]int32) ([]*Message, error) {
	events := []*Message{}
	for _, click := range clicks {
		msg, err := NewMessage(Click, click[0], click[1], '0')
		if err != nil {
			return []*Message{}, fmt.Errorf("error encountered while creating click events in QUEST_TURNIN, err: %s", err.Error())
		}
		if len(click) == 3 {
			msg.delayInSeconds = float32(click[2])
		}
		events = append(events, msg)
	}
	return events, nil
}

// parseQuestTurnInKeyPressGenerator generates key press events from words. It adds an "enter" key press at the end of each word.
func (c *CurrentScript) parseQuestTurnInKeyPressGenerator(words []string) ([]*Message, error) {
	events := []*Message{}
	for _, word := range words {
		for _, c := range word {
			msg, err := NewMessage(KeyPress, 0, 0, c)
			if err != nil {
				return []*Message{}, fmt.Errorf("error encountered while creating keypress events in QUEST_TURNIN, err: %s", err.Error())
			}
			events = append(events, msg)
		}
		msg, err := NewMessage(KeyPress, 0, 0, '#') // "enter" key
		if err != nil {
			return []*Message{}, fmt.Errorf("error encountered while creating keypress events in QUEST_TURNIN, err: %s", err.Error())
		}
		events = append(events, msg)
	}
	return events, nil
}

// parseQuestTurnIn parses a QUEST_TURNIN command. Command is written in the two forms, "QUEST_TURNIN 2" or "QUEST_TURNIN 2 CLICK 234 345"
//
// First value is the command QUEST_TURNIN, second value is the quest position in the quest box, the CLICK command is optional.
// The CLICK command denotes the position of the Quest Reward to select. Only supports selecting one reward.
//
// NOTE: Always attempts to turn in 666 instances of the quest.
func (c *CurrentScript) parseQuestTurnIn(command string) ([]*Message, error) {
	// check if command has CLICK component
	commands := strings.Split(command, "CLICK")
	if len(commands) > 2 {
		return []*Message{}, fmt.Errorf("only one click command is allowed after QUEST_TURNIN, found multiple in \"%s\"", command)
	}
	// this event will be added into the middle of clicks if it exists, created from the CLICK command portion
	questRewardClickEvent := []*Message{}
	var err error
	if len(commands) == 2 {
		// remove the trailing spaces, else equality and conversions will be problematic
		clickCommandSplit := strings.Split(strings.TrimSpace(commands[1]), " ")
		if len(clickCommandSplit) != 2 {
			return []*Message{}, fmt.Errorf("invalid number of params in QUEST_TURNIN for CLICK command, x and y should be separated by space. found \"%s\"", commands[1])
		}
		// the parseClick function already does the validations, thats why we use the Message object instead of an array or two variables
		// to store the x and y variables
		questRewardClickEvent, err = c.parseClick(clickCommandSplit[1:])
		if err != nil {
			return []*Message{}, fmt.Errorf("error encountered in QUEST_TURNIN command, error: %s", err.Error())
		}
	}
	// handling the QUEST_TURNIN portion of the command
	commandSplit := strings.Split(strings.TrimSpace(commands[0]), " ")
	if len(commandSplit) != 2 {
		return []*Message{}, fmt.Errorf("one parameter denoting quest position in quest screen required for QUEST_TURNIN, found too few or many instead in \"%s\"", command)
	}
	questPosition, err := strconv.Atoi(commandSplit[1]) // will extend this to more than the first two positions later by expanding the map
	if err != nil {
		return []*Message{}, fmt.Errorf("invalid second param in QUEST_TURNIN, err: %s", err.Error())
	}
	if questPosition-1 > len(c.questEntryClickPositionMap) || questPosition-1 < 0 {
		return []*Message{}, fmt.Errorf("invalid value of second param in QUEST_TURNIN, value should lie between 1 and %d, found %d", len(c.questEntryClickPositionMap), questPosition)
	}
	events := []*Message{}
	// creating the list of commands to follow
	clicks := [][]int32{
		c.questEntryClickPositionMap[questPosition],
	}
	if len(questRewardClickEvent) == 1 {
		clicks = append(clicks, []int32{questRewardClickEvent[0].x, questRewardClickEvent[1].y, 5})
	}
	clicks = append(clicks,
		[]int32{410, 824, 5}, // "Turn In" button, add 5 seconds delay time after first click, since its slow to load the quest buttons
		[]int32{974, 522, 1}, // add 1 second delay time after clicking on amount of turn ins pop up box
	)
	clickEvents, err := c.parseQuestTurnInClickGenerator(clicks)
	if err != nil {
		return []*Message{}, err
	}
	events = append(events, clickEvents...)
	keyPresses := []string{
		"666",
	}
	keyPressEvents, err := c.parseQuestTurnInKeyPressGenerator(keyPresses)
	if err != nil {
		return []*Message{}, err
	}
	events = append(events, keyPressEvents...)
	clicks = [][]int32{
		{857, 622}, // Press yes on screen for multi turn in
		{292, 825}, // Quest "Back" button
	}
	clickEvents, err = c.parseQuestTurnInClickGenerator(clicks)
	if err != nil {
		return []*Message{}, err
	}
	events = append(events, clickEvents...)
	return events, nil
}
