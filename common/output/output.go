package output

import (
	"bscan/common/utils"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
)

var (
	stringBuilderPool = &sync.Pool{New: func() interface{} { return new(strings.Builder) }}
	mutex             = &sync.Mutex{}
	WindowWidth       = utils.GetWindowWidth()
	out               = colorable.NewColorableStdout()
)

func getTime() string {
	return "[" + time.Now().Format("15:03:04") + "]"
}

func Error(msg string) {
	message := fmt.Sprint(aurora.Red("ERROR").String(), " ", msg)
	fmt.Println(message)
}

func Warning(msg string) {
	message := fmt.Sprint(aurora.Yellow("WARNING").String(), " ", msg)
	fmt.Println(message)
}

func PrintAlive(url string, status int32, length int32, title, app, webserver, desc, os, framework string) {
	message := "\r" + getTime()
	message += aurora.Bold("URL").String() + "[" + aurora.Cyan(url).String() + "] "
	message += aurora.Bold("Status").String() + "[" + aurora.Magenta(status).String() + "] "
	message += aurora.Bold("Size").String() + "[" + aurora.Yellow(length).String() + "] "
	if title != "" {
		message += aurora.Bold("Title").String() + "[" + aurora.Red(title).String() + "] "
	}
	if app != "" {
		message += aurora.Bold("App").String() + "[" + aurora.Green(app).String() + "] "
	}
	if webserver != "" {
		message += aurora.Bold("WebServer").String() + "[" + aurora.Green(webserver).String() + "] "
	}
	if framework != "" {
		message += aurora.Bold("Framework").String() + "[" + aurora.Green(framework).String() + "] "
	}
	if os != "" {
		message += aurora.Bold("Os").String() + "[" + aurora.Green(os).String() + "] "
	}
	if desc != "" {
		message += aurora.Bold("Desc").String() + "[" + aurora.Green(desc).String() + "] "
	}
	message += "% *s\n"
	fmt.Printf(message, WindowWidth-len(message)-1, "")
}

func PrintFound(url string, name string, status int32, length int) {
	message := "\r" + getTime()
	message += aurora.Bold("URL").String() + "[" + aurora.Cyan(url).String() + "] "
	message += aurora.Bold("Status").String() + "[" + aurora.Magenta(status).String() + "] "
	message += aurora.Bold("Size").String() + "[" + aurora.Yellow(length).String() + "] "
	message += aurora.Bold("Found").String() + "[" + aurora.Green(name).String() + "] "
	message += "% *s\n"
	fmt.Printf(message, WindowWidth-len(message)-1, "")
}

func PrintAliveConfig(threads, pnum, timeout int, total int) {
	split := " " + aurora.Magenta("|").String() + " "
	message := ""
	message += aurora.BrightWhite("[*] Starting AliveScan @ " + time.Now().Format("2006-01-02 15:04:05\n")).String()
	message += aurora.BrightWhite("[*] ").String()
	message += aurora.BrightWhite("Threads: ").String() + aurora.Cyan(threads).String() + split
	message += aurora.BrightWhite("Ports: ").String() + aurora.Cyan(pnum).String() + split
	message += aurora.BrightWhite("Timeout: ").String() + aurora.Cyan(timeout).String() + split
	message += aurora.BrightWhite("Total Requests: ").String() + aurora.Cyan(total).String()
	message += "\n"
	fmt.Print(message)
	fmt.Println()
}

func PrintPocConfig(threads, pnum int) {
	split := " " + aurora.Magenta("|").String() + " "
	message := ""
	message += aurora.BrightWhite("[*] Starting PocExploit @ " + time.Now().Format("2006-01-02 15:04:05\n")).String()
	message += aurora.BrightWhite("[*] ").String()
	message += aurora.BrightWhite("Threads: ").String() + aurora.Cyan(threads).String() + split
	message += aurora.BrightWhite("Total Pocs: ").String() + aurora.Cyan(pnum).String()
	message += "\n"
	fmt.Print(message)
	fmt.Println()
}

func Progress(index, total int) {
	sb := stringBuilderPool.Get().(*strings.Builder)
	message := "\r"
	message += "[*] Requested: %s | Progress: %.2f%s"
	message = fmt.Sprintf(message, aurora.BrightWhite(index).String(), aurora.BrightWhite(float32(index)/float32(total)*100), aurora.BrightWhite("%"))
	sb.WriteString(message)
	mutex.Lock()
	fmt.Fprint(out, sb.String())
	mutex.Unlock()
	sb.Reset()
	stringBuilderPool.Put(sb)
}
