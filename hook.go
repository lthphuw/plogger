package plogger

// Hook provide an interface that will call Run on specific Levels before the actual log runned.
type Hook interface {
	Levels() []Level
	Run(*Entry) error
}

// LevelHook defines a type for storing Hooks in each level
type LevelHook map[Level][]Hook

// Add adds hook to specific levels
func (h *LevelHook) Add(hook Hook) {
	for _, l := range hook.Levels() {
		(*h)[l] = append((*h)[l], hook)
	}
}

// Run runs all the hook in this level
func (h *LevelHook) Run(level Level, entry *Entry) error {
	for _, hook := range (*h)[level] {
		if err := hook.Run(entry); err != nil {
			return err
		}
	}
	return nil
}
