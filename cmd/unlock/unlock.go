package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/joeljunstrom/go-luhn"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const defaultImei         = 0000000000000000
const defaultOemCode      = 1000000000000000
const resumeFile          = "unlock.resume"
const unwrittenResumeData = 100

var imei int64
var oemCode int64
var resume bool

func main() {
	flag.Int64Var(&imei, "imei", defaultImei, "IMEI of your phone")
	flag.Int64Var(&oemCode, "oem-code", defaultOemCode, "the OEM code to start with")
	flag.BoolVar(&resume, "resume", false, "create/uses a resume file o resume brute force (tried OEM code frequently written)")
	flag.Parse()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if resume {
			fmt.Printf("\n\nBrute force stopped, you can continue with ./unlock --imei=%d --resume\n", imei)
			writeResume()
		} else {
			fmt.Printf("\n\nBrute force stopped, you can continue with ./unlock --imei=%d --oem-code=%d\n", imei, oemCode)
		}

		os.Exit(1)
	}()

	if luhn.Valid(string(imei)) != true || imei == defaultImei {
		panic(errors.New("invalid IMEI code provided"))
	}

	if resume { readResume() }

	binary, lookErr := exec.LookPath("fastboot")
	if lookErr != nil {
		panic(lookErr)
	}

	writeResumeCountdown := unwrittenResumeData

	for {
		fmt.Printf("=> Try %d\n", oemCode)

		var stdout, stderr bytes.Buffer
		cmd := exec.Command(binary, "oem", "unlock", string(oemCode))
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		writeResumeCountdown--
		if writeResumeCountdown == 0 {
			writeResumeCountdown = unwrittenResumeData
			writeResume()
		}

		err := cmd.Run()
		if err != nil {
			fmt.Println(strings.TrimSpace(string(stderr.Bytes())))
			oemCode += int64(math.Sqrt(float64(imei))*1024)
			continue
		}

		break
	}

	fmt.Printf("\n\nUnlock OEM code: %d\n", oemCode)
}

func readResume() {
	b, err := ioutil.ReadFile(resumeFile)
	if err != nil {
		fmt.Println("Could not read resume information from resume file!")
		return
	}

	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		fmt.Println("Invalid resume information from resume file!")
	}

	oemCode = i
}

func writeResume() {
	err := ioutil.WriteFile(resumeFile, []byte(strconv.FormatInt(oemCode, 10)), os.ModePerm)
	if err != nil {
		fmt.Println("Could not write resume information to resume file!")
	}
}
