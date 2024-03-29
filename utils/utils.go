package utils

import (
	"fmt"
	"sort"
	"time"
    "math/rand"

	"github.com/gen2brain/beeep"
)

// Calculating the mean of a integer array
func ArrayMean(in []int) float64 {
	var sum float64 = 0.0

	for _, item := range in {
		sum += float64(item)
	}
	return sum / float64(len(in))
}

// Calculate the median of an integer array
func ArrayMedian(in []int) int {
	sort.Ints(in)
	indexHalf := len(in) / 2
	if len(in)%2 == 0 {
		return (in[indexHalf-1] + in[indexHalf]) / 2
	} else {
		return in[indexHalf]
	}
}

// Get random item from string array
func RandArrayItem(in  []string) string {
    randomIndex := rand.Intn(len(in))
    return in[randomIndex]
}

// Check if all items in slice are the same in integer array
func ArrayAllItemsEqual(in []int) bool {
	first := in[0]
	for _, item := range in {
		if item != first {
			return false
		}
	}
	return true
}

// Duration to time string
// Format must have hours, minutes, seconds
//    Format example:   "%01d hours : %01d minutes : %01d seconds"
//    Format example:   "%01d hrs %01d minutes %01d seconds"
func DurationToString(duration time.Duration, format string) string {
	hrs := duration / time.Hour
	duration -= hrs * time.Hour
	mins := duration / time.Minute
	duration -= mins * time.Minute
	secs := duration / time.Second
	return fmt.Sprintf(format, hrs, mins, secs)
}

// Display text at a time interval
func DisplayTextWithTime(interval float32, text []string) {
	for _, word := range text {
		fmt.Println(word)
		time.Sleep(time.Duration(interval * float32(time.Second)))
	}
}

// Alert user via notifications
func ShowAlert(title string, message string) {
	err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
	if err != nil {
		panic(err)
	}
	err = beeep.Alert(title, message, "none")
	if err != nil {
		panic(err)
	}
}
