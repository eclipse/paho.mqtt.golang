package paho

type (
	// Logger interface allows implementations to provide to this package any
	// object that implements the methods defined in it.
	Logger interface {
		Println(v ...interface{})
		Printf(format string, v ...interface{})
	}

	// NOOPLogger implements the logger that does not perform any operation
	// by default. This allows us to efficiently discard the unwanted messages.
	NOOPLogger struct{}
)

// Println is the library provided NOOPLogger's
// implementation of the required interface function()
func (NOOPLogger) Println(v ...interface{}) {}

// Printf is the library provided NOOPLogger's
// implementation of the required interface function(){}
func (NOOPLogger) Printf(format string, v ...interface{}) {}

// Internal levels of library output that are initialised to not print
// anything but can be overridden by programmer
var (
	errors Logger = NOOPLogger{}
	debug  Logger = NOOPLogger{}
)

// SetDebugLogger takes an instance of the paho Logger interface
// and sets it to be used by the debug log endpoint
func SetDebugLogger(l Logger) {
	debug = l
}

// SetErrorLogger takes an instance of the paho Logger interface
// and sets it to be used by the error log endpoint
func SetErrorLogger(l Logger) {
	errors = l
}
