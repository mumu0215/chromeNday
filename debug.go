// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
"context"
"log"
"time"

"github.com/chromedp/chromedp"
)

func main() {

	// 禁用chrome headless
	//--origin-trial-disabled-features=SecurePaymentConfirmation --flag-switches-begin --flag-switches-end
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("origin-trial-disabled-features","SecurePaymentConfirmatio"),
		chromedp.Flag("flag-switches-begin",true),
		chromedp.Flag("flag-switches-end",true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var example string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`chrome://version/`),
		chromedp.WaitVisible(`#version`),
		chromedp.OuterHTML(`document.querySelector("#version > span:nth-child(1)")`, &example, chromedp.ByJSPath),
		//chromedp.Value(`document.querySelector("#version")`,&example),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Go's time.After example:\n%s", example)
}

