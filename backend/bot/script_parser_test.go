package bot

import (
	"testing"
)

func TestReadFiles(t *testing.T) {
	expectedOutput := []string{"ATTACK 5", "ATTACK 5"}
	generatedOutput := readFile("../../data/scripts/general_attack.txt")
	if len(generatedOutput) != len(expectedOutput) {
		t.Errorf("Expected output: %+v, Generated Output: %+v", expectedOutput, generatedOutput)
	}
	for idx := range expectedOutput {
		if generatedOutput[idx] != expectedOutput[idx] {
			t.Errorf("Expected output: |%s| => length(%d), Generated Output: |%s| => length(%d)\n",
				expectedOutput[idx],
				len(expectedOutput[idx]),
				generatedOutput[idx],
				len(generatedOutput[idx]))
		}
	}
}
