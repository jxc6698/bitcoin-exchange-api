package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"runtime/debug"
	"fmt"
	"strconv"
)

func init() {
	fmt.Print("")
}


func copyStructValueNested(s, d reflect.Value, ignoreType bool) {
	st := s.Type()
	for i := 0; i < s.NumField(); i++ {
		sf := s.Field(i)
		switch sf.Kind() {
		case reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
			if sf.IsNil() {
				continue
			}
		}

		if s.Type().Field(i).Anonymous {
			copyStructValueNested(sf, d, ignoreType)
			continue
		}
		df := d.FieldByName(st.Field(i).Name)
		if df.IsValid() && (ignoreType || df.Type() == sf.Type()) {
			df.Set(sf)
		}
	}
}

/* Copy the same-name members from one struct to another */
func CopyStructPartial(from interface{}, to interface{}, ignoreType bool) {
	s := reflect.ValueOf(from).Elem()
	d := reflect.ValueOf(to).Elem()
	copyStructValueNested(s, d, ignoreType)
}




func LoadConfig(file string, configPtr interface{}, configJsonPtr interface{}) error {
	//configFile := flag.String("c", file, "config file path")
	//flag.Parse()
	configFile := file

	fileContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileContent, &configJsonPtr)
	if err != nil {
		return err
	}

	CopyStructPartial(configJsonPtr, configPtr, false)
	return nil
}

func LoadConfigJson(file string, configJsonPtr interface{}) error {
	configFile := file

	fileContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileContent, &configJsonPtr)
	if err != nil {
		return err
	}

	return nil
}


/**
 *  read conf file to map[string]
 */
func LoadConfigToMap(file string, submap *map[string]interface{}) error {
	configFile := file

	fileContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileContent, &submap)

	return err
}

type Xlogger struct {
	file     *os.File
	logger   *log.Logger
	stdout   bool
	logLevel int // 0 - only panics, 1 - information, 2 - debugging
}

func NewXlogger(filepath, prefix string, stdoutOutput bool, loglevel int) (*Xlogger, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	logger := log.New(file, prefix, log.LstdFlags | log.Llongfile)

	return &Xlogger{logger: logger, stdout: stdoutOutput, logLevel: loglevel, file:file}, nil
}

func (xl *Xlogger) Fatal(v ...interface{}) {
	stack := debug.Stack()
	xl.Print(v...)
	xl.Print(string(stack[:]))
	os.Exit(1)
}

func (xl *Xlogger) Trace(v ...interface{}) {
	if xl.logLevel == 0 {
		return
	}

	stack := debug.Stack()
	xl.Print(v...)
	xl.Print(string(stack[:]))
}

func (xl *Xlogger) Panic(v ...interface{}) {
	stack := debug.Stack()
	xl.Print("Panic:")
	xl.Print(v...)
	xl.Print(string(stack[:]))
	panic(v)
}

func (xl *Xlogger) Print(v ...interface{}) {
	xl.logger.Print(v...)
	if xl.stdout {
		log.Print(v...)
	}
}

func (xl *Xlogger) Println(v ...interface{}) {
	xl.logger.Println(v...)
	if xl.stdout {
		log.Println(v...)
	}
}

func (xl *Xlogger) Info(v ...interface{}) {
	if xl.logLevel == 0 {
		return
	}

	xl.logger.Println(v...)
	if xl.stdout {
		log.Println(v...)
	}
}

func (xl *Xlogger) Debug(v ...interface{}) {
	if xl.logLevel != 2 {
		return
	}

	xl.logger.Println(v...)
	if xl.stdout {
		log.Println(v...)
	}
}

// Make a process' stderr/stdout pipe to this logger.
func (xl *Xlogger) Setout(writerp *io.Writer) {
	if xl.stdout {
		*writerp = io.MultiWriter(xl.file, os.Stdout)
	} else {
		*writerp = xl.file
	}
}

/* returns true if any panics catched.
*  IMPORTANT: this can only be called by defer statement, not even in inner nested calls in defer!
*/
func PanicCatcher(logger *Xlogger) {
	if err := recover(); err != nil {
		logger.Print(err)
		logger.Print(string(debug.Stack()[:]))
	}
}

func IfErrExit(logger *Xlogger, msg string, err error) {
	if err != nil {
		if logger != nil {
			logger.Fatal("Fatal," + msg + " " + err.Error())
		} else {
			log.Fatal("Fatal," + msg + " " + err.Error())
		}
	}
}


/** restrict {bits} bit behind decimal point  */
func Float64ToString(f float64, bits int) string {
	return strconv.FormatFloat(f, 'f', bits, 64)
}

func Float32ToString(f float32, bits int) string {
	return strconv.FormatFloat(float64(f), 'f', bits, 64)
}
