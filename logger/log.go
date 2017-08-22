package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

func ShortFile(file string) (short string) {
	short = file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return short
}
func prepareFileAndLine(v ...interface{}) (slice []interface{}) {
	_, fn, line, _ := runtime.Caller(2)
	fn = ShortFile(fn)
	s := []interface{}{fn, line}
	slice = append(s, v...)
	return
}

//Debugln calls Println
func Debugln(v ...interface{}) {
	s := prepareFileAndLine(v...)
	debug.Printf("%s:%d: %s", s...)
}

//Debugf calls Printf
func Debugf(format string, v ...interface{}) {
	s := prepareFileAndLine(v...)
	debug.Printf("%s:%d: "+format, s...)
}

//Infoln calls Println
func Infoln(v ...interface{}) {
	s := prepareFileAndLine(v...)
	info.Printf("%s:%d: %s", s...)
}

//Infof calls Printf
func Infof(format string, v ...interface{}) {
	s := prepareFileAndLine(v...)
	info.Printf("%s:%d: "+format, s...)
}

//Warningln calls Println
func Warningln(v ...interface{}) {
	s := prepareFileAndLine(v...)
	warning.Printf("%s:%d: %s", s...)
}

//Warningf calls Printf
func Warningf(format string, v ...interface{}) {
	s := prepareFileAndLine(v...)
	warning.Printf("%s:%d: "+format, s...)
}

//Errorln calls Println
func Errorln(v ...interface{}) {
	s := prepareFileAndLine(v...)
	er.Printf("%s:%d: %s", s...)
}

//Errorf calls Printf
func Errorf(format string, v ...interface{}) {
	s := prepareFileAndLine(v...)
	er.Printf("%s:%d: "+format, s...)
}

//Fatalln calls Fatalln
func Fatalln(v ...interface{}) {
	s := prepareFileAndLine(v...)
	fatal.Fatalf("%s:%d: %s", s...)
}

//Fatalf calls Fatalf
func Fatalf(format string, v ...interface{}) {
	s := prepareFileAndLine(v...)
	fatal.Fatalf("%s:%d: "+format, s...)
}

var (
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	er      *log.Logger
	fatal   *log.Logger
)

// Init all the different levels of log to use on the program depending on the environment.
// For a non-production environment, all logs will be on os.Stdout, as well as a file.
// For a production environment, all logs will be stored on an external file.
func init() {
	//TODO: Stop execution if environment is not set
	env := os.Getenv("ENVIRONMENT")
	var debugHandle, infoHandle, warningHandle, errorHandle, fatalHandle io.Writer
	//	var fp, fd *os.File

	switch env {
	case "production":
		fp, err := os.OpenFile("./gsmc.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file gsmc.log: ", err)
		}
		debugHandle = ioutil.Discard
		infoHandle = fp
		warningHandle = fp
		errorHandle = io.MultiWriter(fp, os.Stdout)
		fatalHandle = io.MultiWriter(fp, os.Stdout)
	case "testing":
		debugHandle = os.Stdout
		infoHandle = os.Stdout
		warningHandle = os.Stdout
		errorHandle = os.Stdout
		fatalHandle = os.Stdout
	default:
		fd, err := os.OpenFile("./gsmc-develop.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file gsmc-develop.log: ", err)
		}
		debugHandle = io.MultiWriter(fd, os.Stdout)
		infoHandle = io.MultiWriter(fd, os.Stdout)
		warningHandle = io.MultiWriter(fd, os.Stdout)
		errorHandle = io.MultiWriter(fd, os.Stdout)
		fatalHandle = io.MultiWriter(fd, os.Stdout)
	}

	debug = log.New(debugHandle,
		"[DEBUG] ",
		log.Ldate|log.Ltime)

	info = log.New(infoHandle,
		"[INFO] ",
		log.Ldate|log.Ltime)

	warning = log.New(warningHandle,
		"[WARN] ",
		log.Ldate|log.Ltime)

	er = log.New(errorHandle,
		"[ERROR] ",
		log.Ldate|log.Ltime)

	fatal = log.New(fatalHandle,
		"[ERROR] ",
		log.Ldate|log.Ltime)
}
