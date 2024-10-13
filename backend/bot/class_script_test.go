package bot

import (
	"testing"
)

func TestClassScript(t *testing.T) {
	c := NewClassScript("")
	c.LoadFromFile("../../data/class_scripts/general_class.txt")
	expectedOutput := []*Message{}
	for i := 0; i < 5; i++ {
		msg, err := NewMessage(Click, c.SkillMap[i].x, c.SkillMap[i].y, '0')
		if err != nil {
			t.Errorf("Failed to create expected output, error: %s", err.Error())
		}
		expectedOutput = append(expectedOutput, msg)
	}
	expectedOutput = append(expectedOutput, expectedOutput...)
	for i := 0; i < 10; i++ {
		output := c.Next()
		if output == nil || !output.Equals(*expectedOutput[i]) {
			t.Errorf("Invalid output. Expected: %+v, found: %+v", *expectedOutput[i], output)
		}
	}
	// TODO: Test empty file name, empty file, file not found, only one value in file
}
