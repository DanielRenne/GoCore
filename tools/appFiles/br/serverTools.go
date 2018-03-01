package br

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	sigar "github.com/cloudfoundry/gosigar"
	"github.com/pkg/errors"
)

var serverToolsSync *sync.RWMutex
var updateAvailable bool

func init() {
	Server.mem = sigar.Mem{}
	Server.swap = sigar.Swap{}
	serverToolsSync = &sync.RWMutex{}
}

type Server_Br struct {
	mem  sigar.Mem
	swap sigar.Swap
}

type Version struct {
	Version      string `json:"version"`
	Error        string `json:"error"`
	AlwaysUpdate bool   `json:"alwaysupdate"`
	ReleaseNotes string `json:"releaseNotes"`
}

func Reboot() (err error) {
	session_functions.FlushAllLogs()
	err = exec.Command("/usr/bin/sudo", "/sbin/shutdown", "-r", "now").Run()
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		session_functions.Log("Error->br->serverTools.go->Reboot", err.Error())
	}
	return
}

func Shutdown() (err error) {
	session_functions.FlushAllLogs()
	time.Sleep(time.Millisecond * 1000)
	err = exec.Command("/usr/bin/sudo", "/sbin/shutdown", "-h", "now").Run()
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		session_functions.Log("Error->br->serverTools.go->Shutdown", err.Error())
	}
	return
}

func SetTimeZone(zone string) (err error) {
	err = exec.Command("/usr/bin/sudo", "/usr/bin/timedatectl", "set-timezone", zone).Run()
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		session_functions.Log("Error->br->serverTools.go->SetTimeZone", err.Error())
	}
	return
}

func SetDate(date string) (err error) {
	session_functions.Log("Setting Linux Date", date)

	err = exec.Command("/usr/bin/sudo", "/usr/bin/timedatectl", "set-ntp", "0").Run()
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		session_functions.Log("Error->br->serverTools.go->SetDate", err.Error())
	}

	t := time.Now().In(Schedules.GetLocation())

	session_functions.Log("Setting Linux Date Exec", "/usr/bin/sudo"+" "+"/usr/bin/timedatectl"+" "+"set-time"+" "+"\""+date+" "+fmt.Sprintf("%+v", t.Format("15:04:05"))+"\"")
	err = exec.Command("/usr/bin/sudo", "/usr/bin/timedatectl", "set-time", date+" "+fmt.Sprintf("%+v", t.Format("15:04:05"))).Run()
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		session_functions.Log("Error->br->serverTools.go->SetDate", err.Error())
	}
	return
}

func SetTime(timeString string) (err error) {
	session_functions.Log("Setting Linux Time", timeString)

	err = exec.Command("/usr/bin/sudo", "/usr/bin/timedatectl", "set-ntp", "0").Run()
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		session_functions.Log("Error->br->serverTools.go->SetDate", err.Error())
	}

	format := "03:04:05 PM MST"
	tme, err := time.Parse(format, timeString)
	session_functions.Log("Setting Linux Time Exec", "/usr/bin/sudo"+" "+"/usr/bin/timedatectl"+" "+"set-time"+" "+fmt.Sprintf("%+v", tme.Format("15:04:05")))
	err = exec.Command("/usr/bin/sudo", "/usr/bin/timedatectl", "set-time", fmt.Sprintf("%+v", tme.Format("15:04:05"))).Run()
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		session_functions.Log("Error->br->serverTools.go->SetTime", err.Error())
	}
	return
}

func SetDateAndTime(date string, timeString string) (err error) {
	session_functions.Log("Setting Linux Date & Time", date+"  "+timeString)

	err = exec.Command("/usr/bin/sudo", "/usr/bin/timedatectl", "set-ntp", "0").Run()
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		session_functions.Log("Error->br->serverTools.go->SetDate", err.Error())
	}

	format := "03:04:05 PM MST"
	tme, err := time.Parse(format, timeString)
	session_functions.Log("Setting Linux Time and Date Exec", "/usr/bin/sudo"+" "+"/usr/bin/timedatectl"+" "+"set-time"+" "+"\""+date+" "+" "+fmt.Sprintf("%+v", tme.Format("15:04:05"))+"\"")
	err = exec.Command("/usr/bin/sudo", "/usr/bin/timedatectl", "set-time", date+" "+fmt.Sprintf("%+v", tme.Format("15:04:05"))).Run()
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		session_functions.Log("Error->br->serverTools.go->SetDate", err.Error())
	}
	return
}

func EnableNTPServer() (err error) {
	session_functions.Log("Enabling NTP Server", "")

	err = exec.Command("/usr/bin/sudo", "/usr/bin/timedatectl", "set-ntp", "true").Run()
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		session_functions.Log("Error->br->serverTools.go->SetDate", err.Error())
	}

	return
}

func IsUpdateAvailable() (available bool) {
	serverToolsSync.RLock()
	available = updateAvailable
	serverToolsSync.RUnlock()
	return
}

func setUpdateAvailable(available bool) {
	serverToolsSync.RLock()
	updateAvailable = available
	serverToolsSync.RUnlock()
	return
}

func (self Server_Br) formatSize(val uint64) uint64 {
	return val / 1024
}

func (self Server_Br) TotalMemory() uint64 {
	self.mem.Get()
	return self.formatSize(self.mem.Total)
}

func (self Server_Br) UsedMemeory() uint64 {
	self.mem.Get()
	return self.formatSize(self.mem.Total)
}
