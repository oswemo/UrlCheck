package utils

// Logging handling for the project.  Log in JSON format, a map of key value pairs.
// The caller sends any key/value pairs it desires, and we automatically add the function
// and package in which the log call was performed.

import (
    "strings"
    "runtime"
    log "github.com/Sirupsen/logrus"
)

// Stolen from logrus.  Allows use of the type in our classes without having
// to include logrus everywhere.  May be a better way, but not aware of it.
type LogFields map[string]interface{}

// PackageName returns the go package of the caller.
// Returns [package name] [method name]
func packageAndFunc(offset int) (string, string) {
    fpcs := make([]uintptr, 1)

    n := runtime.Callers(offset, fpcs)
    if n == 0 {
       return "", ""
    }

    // get the info of the actual function that's in the pointer
    fun := runtime.FuncForPC(fpcs[0]-1)
    if fun == nil {
      return "", ""
    }

    parts := strings.Split(fun.Name(), ".")
    pkgName := strings.Join(parts[:len(parts)-1], ".")
    funName := parts[len(parts)-1]
    return pkgName, funName
}

// PackageName returns the name of the calling package.
// Offset of 3 is used to return the package of the method calling *this* function.
func PackageName() (string) {
    pkg, _ := packageAndFunc(3)
    return pkg
}

// FunctionName returns the name of the calling function.
// Offset of 3 is used to return the function of the method calling *this* function.
func FunctionName() (string) {
    _,fun := packageAndFunc(3)
    return fun
}

// getFields is a private function that returns the fields that should be logged.
// When called, it takes the given fields and adds a field for the package name and
// a field for the function.
func getFields(fields map[string]interface{}) (log.Fields) {
    logFields := log.Fields{}
    for key := range fields {
        logFields[key] = fields[key]
    }

    pkg, fun := packageAndFunc(3)
    logFields["package"] = pkg
    logFields["function"] = fun
    return logFields
}

// SetDebug sets the log level to debug to show additional logging.
func SetDebug() {
    log.SetLevel(log.DebugLevel)
}

// LogInfo logs informational data
func LogInfo(fields map[string]interface{}, format string, vargs ...interface{}) {
    log.SetFormatter(&log.JSONFormatter{})
    log.WithFields(getFields(fields)).Infof(format, vargs...)
}

// LogError logs error data
func LogError(fields map[string]interface{}, err error, format string, vargs ...interface{}) {
    fields["error"] = err.Error()
    log.SetFormatter(&log.JSONFormatter{})
    log.WithFields(getFields(fields)).Errorf(format, vargs...)
}

// LogDebug logs debug data
func LogDebug(fields map[string]interface{}, format string, vargs ...interface{}) {
    log.SetFormatter(&log.JSONFormatter{})
    log.WithFields(getFields(fields)).Debugf(format, vargs...)
}
