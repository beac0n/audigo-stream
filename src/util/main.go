package util

import (
	"regexp"
	"remote-audio-streamer/src/commands"
	"strconv"
)

func GetPulseAudioInputSinkIndexes() ([]int64, error) {
	regex, regexErr := regexp.Compile("index: (\\d+)")
	if regexErr != nil {
		return nil, regexErr
	}

	output, pacmdErr := commands.ListSinkInputs()
	if pacmdErr != nil {
		return nil, pacmdErr
	}

	return toIntList(regex, output), nil
}

func Difference(slice1, slice2 []int64) []int64 {
	var diff []int64

	mDiff := map[int64]byte{}
	m1 := map[int64]byte{}
	m2 := map[int64]byte{}

	for _, val := range slice1 {
		m1[val] = 1
	}

	for _, val := range slice2 {
		m2[val] = 1
	}

	// we only care about elements which are in m2 but NOT in m1
	for key, val := range m2 {
		if val != m1[key] {
			mDiff[key] = 1
		}
	}

	for key := range mDiff {
		diff = append(diff, key)
	}

	return diff
}

func toIntList(regex *regexp.Regexp, output string) []int64 {
	matches := regex.FindAllStringSubmatch(output, -1)
	integers := make([]int64, len(matches))

	for i, match := range matches {
		if integer, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			integers[i] = integer
		}
	}

	return integers
}
