package bot

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Bot struct
type Bot struct {
	ctx           context.Context
	scripts       map[string]string
	classScripts  map[string]string
	isRunning     bool
	currentScript *CurrentScript
	eventChannel  chan *Message
	waitGroup     *sync.WaitGroup
}

// NewBot creates a new Bot application struct
func NewBot() *Bot {
	return &Bot{
		isRunning:    false,
		eventChannel: make(chan *Message, 500),
		waitGroup:    &sync.WaitGroup{},
		scripts:      map[string]string{},
		classScripts: map[string]string{},
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *Bot) Startup(ctx context.Context) {
	a.ctx = ctx
	a.waitGroup.Add(1)
	go MessageSenderConsumer(a.eventChannel, a.waitGroup)
}

func (a *Bot) Shutdown(ctx context.Context) {
	a.isRunning = false
	exitMessage, err := NewMessage(Exit, 0, 0, '0')
	if err != nil {
		panic(fmt.Sprintf("[PANIC] Could not shut down bot, err: %s", err.Error()))
	}
	a.eventChannel <- exitMessage
	a.waitGroup.Wait()
}

func (a *Bot) loadFiles(path string, isClassScript bool) []string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("[ERROR] Unable to create directory to store scripts. Please check file permissions. Error: %s\n", err.Error())
			return []string{}
		}
	}

	foundScripts, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("[ERROR] Error encountered: %+v", err)
		return []string{}
	}

	finalScripts := []string{}

	for _, script := range foundScripts {
		name := script.Name()
		splitName := strings.Split(name, "/")
		splitName = strings.Split(splitName[len(splitName)-1], "\\")
		name = strings.Split(splitName[len(splitName)-1], ".")[0]
		name = strings.Join(strings.Split(name, "_"), " ")
		finalScripts = append(finalScripts, name)
		if isClassScript {
			a.classScripts[name] = path + "/" + script.Name()
		} else {
			a.scripts[name] = path + "/" + script.Name()
		}
	}
	return finalScripts
}

func (a *Bot) LoadClasses() []string {
	return a.loadFiles("./data/class_scripts", true)
}

func (a *Bot) LoadScripts() []string {
	return a.loadFiles("./data/scripts", false)
}

func (a *Bot) StartBot(newClass string, newScript string) error {
	if _, ok := a.scripts[newScript]; !ok {
		return fmt.Errorf("unable to find \"%s\" script's file", newScript)
	}
	if _, ok := a.classScripts[newClass]; !ok {
		return fmt.Errorf("unable to find \"%s\" class' file", newClass)
	}
	if a.currentScript == nil || a.currentScript.ID != newScript {
		classScript := NewClassScript(newClass)
		classScript.LoadFromFile(a.classScripts[newClass])
		a.currentScript = NewCurrentScript(newScript)
		a.currentScript.LoadClassScript(classScript)
		a.currentScript.LoadFromFile(a.scripts[newScript])
	}
	if a.currentScript.classScript == nil || a.currentScript.classScript.ID != newClass {
		classScript := NewClassScript(newClass)
		classScript.LoadFromFile(a.classScripts[newClass])
		a.currentScript.LoadClassScript(classScript)
	}
	a.isRunning = true
	a.waitGroup.Add(1)
	go a.scriptLoop()
	return nil
}

func (a *Bot) StopBot() {
	a.isRunning = false
}

func (a *Bot) scriptLoop() {
	defer a.waitGroup.Done()
	for a.isRunning {
		if a.currentScript == nil {
			fmt.Printf("[WARN] No script loaded, exiting\n")
			return
		}
		timeToWait := int32(a.currentScript.NextWaitTimeInSeconds() * 1000) // convert to milliseconds
		time.Sleep(time.Duration(timeToWait) * time.Millisecond)
		event := a.currentScript.Next()
		// No need to send an event if its just for sleeping. TODO: Move this to the event consumer when making code clearer
		if event.msg_type == Delay {
			timeToWait = int32(event.delayInSeconds * 1000)
			time.Sleep(time.Duration(timeToWait) * time.Millisecond)
			continue
		}
		a.eventChannel <- event
	}
}
