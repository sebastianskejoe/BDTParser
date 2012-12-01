package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
    "log"
    "io"
    "strings"
    "strconv"
    "unicode"
)

/*const (
	TIME        = 0x01
	HR          = 0x02
	VENTILATION = 0x04
	VO2         = 0x08
	VCO2        = 0x10
	KONDITAL    = 0x20
	RER         = 0x40
)*/

var (
	time  []string
	hr    []int
	vent  []int
	vo2   []int
	vco2  []int
	kondi []float64
	rer   []float64

	// Flags
    timeflag    = flag.Bool("time", true, "Print time values")
    hrflag      = flag.Bool("hr", true, "Print heart-rate values")
    ventflag    = flag.Bool("ventilation", true, "Print ventilation values")
    vo2flag     = flag.Bool("vo2", true, "Print V'O2 values")
    vco2flag    = flag.Bool("vco2", true, "Print V'CO2 values")
    kondiflag   = flag.Bool("kondital", true, "Print kondital(V'O2/kg)")
    rerflag     = flag.Bool("rer", true, "Print respiratory exchange ratios")
)

func main() {
    flag.Parse()
    if flag.NArg() != 1 {
        printUsage()
        return
    }

    // Open file and create a reader
    file,err := os.Open(flag.Arg(0))
    if err != nil {
        log.Fatal(err)
        return
    }
    r := bufio.NewReader(file)
    line, err := r.ReadString('\r')
    for ; err != io.EOF ; line,err = r.ReadString('\r') {
        line = strings.TrimSpace(line)
        // If first character isn't digit, we don't want to parse it (for now)
        if len(line) < 1 || !unicode.IsDigit(rune(line[0])) {
            continue
        }

        spaces := 0
        parts := strings.Split(line, " ")
        for i,part := range parts {
            if part == "" {
                spaces++
                continue
            }
            switch i-spaces {
            case 0:
                time = append(time, part)
            case 1:
                h,_ := strconv.Atoi(part)
                hr = append(hr, h)
            case 2:
                v,_ := strconv.Atoi(part)
                vent = append(vent, v)
            case 3:
                v,_ := strconv.Atoi(part)
                vo2 = append(vo2, v)
            case 4:
                v,_ := strconv.Atoi(part)
                vco2 = append(vco2, v)
            case 5:
                k,_ := strconv.ParseFloat(part, 64)
                kondi = append(kondi, k)
            case 6:
                r,_ := strconv.ParseFloat(part, 64)
                rer = append(rer, r)
            }
        }
    }

    printData()
}

func printUsage() {
    fmt.Println("Usage:",os.Args[0], "[OPTION]... FILE")
    flag.PrintDefaults()
}

func printData() {
    for line,_ := range time {
        parts := make([]string,0)
        if *timeflag {
            parts = append(parts, time[line])
            parts = append(parts, strconv.Itoa(minsecToSec(time[line])))
        }
        if *hrflag {
            parts = append(parts, strconv.Itoa(hr[line]))
        }
        if *ventflag {
            parts = append(parts, strconv.Itoa(vent[line]))
        }
        if *vo2flag {
            parts = append(parts, strconv.Itoa(vo2[line]))
        }
        if *vco2flag {
            parts = append(parts, strconv.Itoa(vco2[line]))
        }
        if *kondiflag {
            parts = append(parts, fmt.Sprintf("%.1f", kondi[line]))
        }
        if *rerflag {
            parts = append(parts, fmt.Sprintf("%.2f", rer[line]))
        }
        fmt.Println(strings.Join(parts, "\t"))
    }
}

func minsecToSec(minsec string) (int) {
    parts := strings.Split(minsec, ":")
    min,_ := strconv.Atoi(parts[0])
    sec,_ := strconv.Atoi(parts[1])
    return min*60+sec
}
