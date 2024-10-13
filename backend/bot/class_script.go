package bot

import (
	"fmt"
	"strconv"
	"strings"
)

type classScript struct {
	ID         string
	SkillMap   []*Message
	Script     []*Message
	currentIdx int
}

func NewClassScript(ID string) *classScript {
	return &classScript{
		ID:       ID,
		SkillMap: createSkillMap(),
	}
}

func createSkillMap() []*Message {
	skillMap := make([]*Message, 5)
	var err error
	positions := [][]int32{
		{752, 948},  // AUTOATTACK
		{836, 945},  // SKILL 1
		{917, 946},  // SKILL 2
		{999, 945},  // SKILL 3
		{1082, 947}, // SKILL 4
	}
	for i := 0; i < 5; i++ {
		skillMap[i], err = NewMessage(Click, positions[i][0], positions[i][1], '0')
		if err != nil {
			panic(err)
		}
		skillMap[i].delayInSeconds = 0.3
	}
	return skillMap
}

func (c *classScript) LoadFromFile(fname string) {
	commands := readFile(fname)
	if len(commands) == 0 {
		fmt.Printf("[ERROR] Empty file found: %s\n", fname)
		return
	}
	c.Script = []*Message{}
	for _, command := range commands {
		commandSplit := strings.Split(command, " ")
		var val *Message
		switch commandSplit[0] {
		case "0":
			fallthrough
		case "AUTOATTACK":
			val = c.SkillMap[0]
		case "1":
			fallthrough
		case "2":
			fallthrough
		case "3":
			fallthrough
		case "4":
			val = c.SkillMap[int(commandSplit[0][0])-int('0')]
		case "CLICK":
			if len(commandSplit) != 3 {
				fmt.Printf("[ERROR] Command \"%s\" ignored because too many parameters, correct format is CLICK x, y\n", command)
				continue
			}
			cleanXString := cleanStringOfSpaces(commandSplit[0])
			cleanYString := cleanStringOfSpaces(commandSplit[1])
			if validateStrForClickEvent(cleanXString) || validateStrForClickEvent(cleanYString) {
				fmt.Printf("[ERROR] Malformed parameters for CLICK event encountered while reading from %s", fname)
				fmt.Printf("X: %s, Y: %s, skipping command\n", commandSplit[0], commandSplit[1])
				continue
			}
			// skipping errors because they are handled in above validation functions
			x, _ := strconv.Atoi(cleanXString)
			y, _ := strconv.Atoi(cleanYString)
			var err error
			val, err = NewMessage(Click, int32(x), int32(y), '0')
			if err != nil {
				fmt.Printf("[ERROR] Error while creating new message during reading from %s, error: %s, params: %d, %d, skipping\n", fname, err.Error(), x, y)
				continue
			}
		}
		c.Script = append(c.Script, val)
		c.currentIdx = 0
	}
}

func (c *classScript) Next() *Message {
	if len(c.Script) == 0 {
		return nil
	}
	oldIdx := c.currentIdx
	c.currentIdx++
	if c.currentIdx == len(c.Script) {
		c.currentIdx = 0
	}
	return c.Script[oldIdx]
}
