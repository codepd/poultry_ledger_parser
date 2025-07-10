package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	eggTypes     = []string{"LARGE EGG", "MEDIUM EGG", "CORRECT EGG", "SMALL EGG"}
	mashTypes    = []string{"Layer Mash Bulk", "Grower Mash Bulk", "Pre Layer Mash Bulk"}
	medicineList = []string{"D3 Forte", "Vetmulin", "OXYCYCLINE", "Tiazin", "BPPS Forte", "CTC", "SHELL GRIT", "Rovimix", "Cholimarin", "Zagromin", "G Pro Naturo", "Necrovet", "Toxol"}
	paymentKey   = "Payment"
	billnoKey    = "BILLNO"
)

func extractAmount(line string) float64 {
	// Find the last field that looks like a number (possibly with decimals)
	parts := strings.Fields(line)
	for i := len(parts) - 1; i >= 0; i-- {
		amt := strings.ReplaceAll(parts[i], ",", "")
		if val, err := strconv.ParseFloat(amt, 64); err == nil {
			return val
		}
	}
	return 0
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	eggs := make(map[string][]string)
	mashes := make(map[string][]string)
	medicines := make(map[string][]string)
	payments := []string{}
	billnos := []string{}
	var lastDate string
	eggTotals := make(map[string]float64)
	mashTotals := make(map[string]float64)
	medicineTotals := make(map[string]float64)
	paymentTotal := 0.0
	eggCounts := make(map[string]int)
	mashCounts := make(map[string]int)
	medicineCounts := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, ",", "")
		fields := strings.Fields(line)
		// If the line starts with a date, update lastDate
		if len(fields) > 0 && strings.HasSuffix(fields[0], "-25") {
			lastDate = fields[0]
		}
		// BILLNO grouping with date: use date from this line if present, else from above
		if strings.Contains(line, billnoKey) {
			date := lastDate
			if len(fields) > 0 && strings.HasSuffix(fields[0], "-25") {
				date = fields[0]
			}
			billnos = append(billnos, date+" "+line)
			continue
		}
		// EGG grouping
		for _, egg := range eggTypes {
			if strings.Contains(line, egg) {
				eggs[egg] = append(eggs[egg], line)
				eggTotals[egg] += extractAmount(line)
				eggCounts[egg]++
			}
		}
		// Mash grouping
		for _, mash := range mashTypes {
			if mash == "Pre Layer Mash Bulk" {
				if strings.Contains(line, mash) {
					mashes[mash] = append(mashes[mash], line)
					mashTotals[mash] += extractAmount(line)
					mashCounts[mash]++
				}
			} else if mash == "Layer Mash Bulk" {
				if strings.Contains(line, mash) && !strings.Contains(line, "Pre Layer Mash Bulk") {
					mashes[mash] = append(mashes[mash], line)
					mashTotals[mash] += extractAmount(line)
					mashCounts[mash]++
				}
			} else {
				if strings.Contains(line, mash) {
					mashes[mash] = append(mashes[mash], line)
					mashTotals[mash] += extractAmount(line)
					mashCounts[mash]++
				}
			}
		}
		// Medicine grouping
		for _, med := range medicineList {
			if strings.Contains(line, med) {
				medicines[med] = append(medicines[med], line)
				medicineTotals[med] += extractAmount(line)
				medicineCounts[med]++
			}
		}
		// Payment grouping
		if strings.Contains(line, paymentKey) {
			payments = append(payments, line)
			paymentTotal += extractAmount(line)
		}
	}

	fmt.Println("EGG GROUPS:")
	for k, v := range eggs {
		fmt.Printf("%s (Count: %d, Total: %.2f):\n", k, eggCounts[k], eggTotals[k])
		for _, l := range v {
			fmt.Println("  ", l)
		}
	}
	fmt.Println("\nMASH GROUPS:")
	for k, v := range mashes {
		fmt.Printf("%s (Count: %d, Total: %.2f):\n", k, mashCounts[k], mashTotals[k])
		for _, l := range v {
			fmt.Println("  ", l)
		}
	}
	fmt.Println("\nMEDICINE GROUPS:")
	for k, v := range medicines {
		fmt.Printf("%s (Count: %d, Total: %.2f):\n", k, medicineCounts[k], medicineTotals[k])
		for _, l := range v {
			fmt.Println("  ", l)
		}
	}
	fmt.Printf("\nPAYMENTS (Total: %.2f):\n", paymentTotal)
	for _, l := range payments {
		fmt.Println("  ", l)
	}
	fmt.Println("\nBILLNOS:")
	for _, l := range billnos {
		fmt.Println("  ", l)
	}
}
