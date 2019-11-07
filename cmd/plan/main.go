package main

import (
	"bytes"
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
)

// will be replaced with git hash
var build = "develop"

const (
	subject = iota
	startDate
	startTime
	endDate
	endTime
	allDayEvent
	description
	location
	private
)

type recordReader func([]string)

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}

func run() error {
	// =========================================================================
	// Logging
	log := log.New(os.Stdout, "PLAN : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	var cfg struct {
		CSVPath    string `conf:"default:input.csv"`
		Output     string `conf:"default:plan.ics"`
		Timezone   string `conf:"default:NZDT"`
		CalendarID string `conf:"default://HungryLegs//HungryLegs Plan App//EN"`
		Calendar   string `conf:"default:HungryLegs"`
	}

	if err := conf.Parse(os.Args[1:], "PLAN", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("PLAN", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	inFile := cfg.CSVPath
	outFile := cfg.Output

	log.Printf("Using version %v\n", build)
	log.Printf("Using file: %v\n", inFile)

	var ics bytes.Buffer

	log.Printf("Creating prolog")
	prolog(&ics, cfg.Calendar, cfg.Timezone, cfg.CalendarID)

	log.Printf("Parsing file")
	foo := func(record []string) {
		err := formatRecord(record, &ics, cfg.Calendar)
		if err != nil {
			log.Printf("Couldn't format record: %v", err)
		}
	}
	err := parseFile(inFile, foo)
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}

	log.Printf("Writing epilog")
	epilog(&ics)

	log.Printf("Creating .ics %v\n", outFile)
	writeICS(outFile, &ics)

	return nil
}

func formatRecord(record []string, event *bytes.Buffer, calendar string) error {
	uuid, err := newid()
	if err != nil {
		panic("Bad id gen")
	}

	if len(record)-1 < private {
		fmt.Printf("%v %v", len(record), private)
		return errors.New("Bad record length")
	}

	formatDate := strings.Replace(record[startDate], "-", "", -1)
	if formatDate == "" {
		return errors.New("Event missing date")
	}

	event.WriteString("BEGIN:VEVENT\r\n")
	fmt.Fprintf(event, "DTSTAMP:%vT000000Z\r\n", formatDate)
	fmt.Fprintf(event, "UID:ROHAN-%v\r\n", uuid)
	fmt.Fprintf(event, "DTSTART;VALUE=DATE:%v\r\n", formatDate)
	fmt.Fprintf(event, "DTEND;VALUE=DATE:%v\r\n", formatDate)
	fmt.Fprintf(event, "SUMMARY:%v\r\n", record[subject])
	fmt.Fprintf(event, "DESCRIPTION:%v\r\n", record[description])
	fmt.Fprintf(event, "CATEGORIES:%v\r\n", calendar)
	event.WriteString("END:VEVENT\r\n")

	return nil
}

func prolog(prolog *bytes.Buffer, calendar string, timeZone string, prodID string) {
	prolog.WriteString("BEGIN:VCALENDAR\r\n")
	prolog.WriteString("VERSION:2.0\r\n")
	fmt.Fprintf(prolog, "X-WR-CALNAME:%v\r\n", calendar)
	fmt.Fprintf(prolog, "PRODID:%v\r\n", prodID)
	fmt.Fprintf(prolog, "X-WR-TIMEZONE:%v\r\n", timeZone)
	fmt.Fprintf(prolog, "X-WR-CALDESC:%v\r\n", calendar)
	prolog.WriteString("CALSCALE:GREGORIAN\r\n")
}

func epilog(epilog *bytes.Buffer) {
	epilog.WriteString("END:VCALENDAR\r\n")
}

func writeICS(path string, buffer *bytes.Buffer) {
	err := ioutil.WriteFile(path, buffer.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func parseFile(path string, fn recordReader) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		// ignore header row
		if record[subject] != "Subject" {
			fn(record)
		}
	}

	return nil
}

// Thanks internet
// https://stackoverflow.com/questions/15130321/is-there-a-method-to-generate-a-uuid-with-go-language#15134490
// Close enough for jazz.
func newid() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		return "", err
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	return hex.EncodeToString(u), nil
}
