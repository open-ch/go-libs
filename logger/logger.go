package logger

// logger interface for reuse, to respect the spirit of
// SCRUBBED-URL

// For some details about having logger work with both new things like Zap and more traditional ones like logrus,
// see https://www.mountedthoughts.com/golang-logger-interface/

//Logger is our contract for the logger
type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Panicf(format string, args ...interface{})

	// WithFields(keyValues Fields) Logger
}
