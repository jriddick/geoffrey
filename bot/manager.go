package bot

import (
	"os"
	"sync"
	"syscall"

	"os/signal"
)

// Manager is a bot manager for multiple bots
type Manager struct {
	bots     map[string]*Bot
	signals  chan os.Signal
	wait     sync.WaitGroup
	running  bool
	handlers map[string]Handler
}

// NewManager creates and returns a new manager
func NewManager() *Manager {
	manager := &Manager{
		running:  false,
		signals:  make(chan os.Signal, 1),
		bots:     make(map[string]*Bot),
		handlers: make(map[string]Handler),
	}

	// Notify when program receives certain signals
	signal.Notify(manager.signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	return manager
}

// Add adds a bot to the manager and starts it if the
// manager is currently running.
func (m *Manager) Add(name string, bot *Bot) error {
	// Check if key already exists
	if _, ok := m.bots[name]; ok {
		return ErrBotExists
	}

	// Check if we should start the bot
	if m.running {
		if err := bot.Connect(); err != nil {
			return err
		}
	}

	// Add the bot to the manager
	m.bots[name] = bot

	return nil
}

// Remove removes the bot from the manager. It will also
// stop the bot if its currently running.
func (m *Manager) Remove(name string) error {
	if _, ok := m.bots[name]; ok {
		// Stop bot if manager is running
		if m.running {
			m.bots[name].Close()
		}

		delete(m.bots, name)
	} else {
		return ErrBotNotFound
	}

	return nil
}

// Start will start the manager and all added bots. It will
// only return successfully if *ALL* bots start. If one bot
// fails it will stop all other bots.
func (m *Manager) Start() error {
	for _, bot := range m.bots {
		if err := bot.Connect(); err != nil {
			// Run through all bots and stop them
			for _, bot := range m.bots {
				bot.Close()
			}

			// Continue with the error handling
			return err
		}

		// Run the bot in another thread
		go bot.Run()
	}

	m.running = true
	return nil
}

// Stop will stop all running bots
func (m *Manager) Stop() {
	for _, bot := range m.bots {
		bot.Close()
	}

	m.running = false
}

// Run will start all bots and start to handle all messages
func (m *Manager) Run() error {
	if m.running == false {
		if err := m.Start(); err != nil {
			return err
		}
	}

	go func() {
		// Wait until we get shutdown signal
		<-m.signals

		// Stop the waiter
		m.wait.Done()
	}()

	// Make sure we wait until we get signals
	m.wait.Add(1)
	m.wait.Wait()

	// Stop all bots and return
	m.Stop()
	return nil
}
