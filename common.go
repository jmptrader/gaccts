/*

This is being referenced from the gaccts/common.go file

*/

package gaccts

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Money int64
func Float64ToMoney(f float64) Money { return Money(int64(f *100 + 0.5)) }

var rootDir string
var dataDir string

func Enfloat64(s string) (f float64, err error) {
	s1 := strings.Replace(s, ",", "", -1)
	f, err = strconv.ParseFloat(s1, 64)
	return
}

func AsFloat64(s string) (f float64) {
	f, err := Enfloat64(s)
	if err != nil { panic(err)}
	return
}

func Init() {

	// set root directory
	rootDir = "C:/cygwin/home/mcarter/"	

	trydir := func (d string) {
		_, err := os.Stat(d)
		if err == nil { rootDir = d }
		}

	trydir("/home/mcarter/")

	dataDir = rootDir + "redact/docs/accts2014/"
}


// Obtain user's home directory
// http://stackoverflow.com/questions/7922270/obtain-users-home-directory
func UserHomeDir() string {
    if runtime.GOOS == "windows" {
        home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
        if home == "" {
            home = os.Getenv("USERPROFILE")
        }
        return home
    }
    return os.Getenv("HOME")
}


func before() string { return "2013-04-05" } // TODO should be generalised
