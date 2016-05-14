package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

//Debugln calls Println
func Debugln(v ...interface{}) {
	debug.Println(v...)
}

//Debugf calls Printf
func Debugf(format string, v ...interface{}) {
	debug.Printf(format, v...)
}

//Infoln calls Println
func Infoln(v ...interface{}) {
	info.Println(v...)
}

//Infof calls Printf
func Infof(format string, v ...interface{}) {
	info.Printf(format, v...)
}

//Warningln calls Println
func Warningln(v ...interface{}) {
	warning.Println(v...)
}

//Warningf calls Printf
func Warningf(format string, v ...interface{}) {
	warning.Printf(format, v...)
}

//Errorln calls Println
func Errorln(v ...interface{}) {
	er.Println(v...)
}

//Errorf calls Printf
func Errorf(format string, v ...interface{}) {
	er.Printf(format, v...)
}

//Fatal calls Fatal
func Fatal(v ...interface{}) {
	fatal.Fatal(v...)
}

//Fatalln calls Fatalln
func Fatalln(v ...interface{}) {
	fatal.Fatalln(v...)
}

//Fatalf calls Fatalf
func Fatalf(format string, v ...interface{}) {
	fatal.Fatalf(format, v...)
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
func Init() {
	/*

		traceHandle io.Writer,
		infoHandle io.Writer,
		warningHandle io.Writer,
		errorHandle io.Writer
	*/
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
		log.Ldate|log.Ltime|log.Lshortfile)

	info = log.New(infoHandle,
		"[INFO] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	warning = log.New(warningHandle,
		"[WARN] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	er = log.New(errorHandle,
		"[ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	fatal = log.New(fatalHandle,
		"[ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
