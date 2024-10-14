package bot

import (
	"fmt"
	"sync"
	"testing"
)

func TestSenderConsumer(t *testing.T) {
	input := make(chan *Message, 500)
	var err error
	defer func() {
		close(input)
		if err != nil {
			t.Errorf(err.Error())
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go MessageSenderConsumer(input, &wg)
	testClick1, err := NewMessage(Click, 532, 548, ' ')
	if err != nil {
		return
	}
	input <- testClick1
	fmt.Println("[INFO] Pushed testClick1")
	testClick2, err := NewMessage(Click, 545, 520, ' ')
	if err != nil {
		return
	}
	input <- testClick2
	fmt.Println("[INFO] Pushed testClick2")
	testText, err := NewMessage(KeyPress, 0, 0, 'i')
	if err != nil {
		return
	}
	input <- testText
	fmt.Println("[INFO] Pushed testText(1)")
	input <- testText
	fmt.Println("[INFO] Pushed testClick(2)")
	testExit, err := NewMessage(Exit, 0, 0, 'a')
	if err != nil {
		return
	}
	input <- testExit
	fmt.Println("[INFO] Pushed testExit")
	wg.Wait()
}
