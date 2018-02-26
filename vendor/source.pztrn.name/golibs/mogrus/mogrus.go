package mogrus

type MogrusLogger struct {
    // Initialized loggers.
    // Key is a name of logger.
    loggers map[string]*LoggerHandler
}

// Creates new logger handler, adds it to list of known loggers and
// return it to caller.
// Note that logger handler will be "just initialized", to actually
// use it you should add output with LoggerHandler.CreateOutput().
func (ml *MogrusLogger) CreateLogger(name string) *LoggerHandler {
    lh := &LoggerHandler{}
    lh.Initialize()
    ml.loggers[name] = lh

    return lh
}

func (ml *MogrusLogger) Initialize() {
    ml.loggers = make(map[string]*LoggerHandler)
}
