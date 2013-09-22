/*
 * Author: Sidney Kochman (sidney@kochman.org)
 * Created: 2013-09-21
 * Description: Ciphers or deciphers a string.
 *
 * The cipher only encodes characters defined in the CHARS constant.
 * The input string is broken up into segments with maximum lengths
 * of 10 characters (this is defined in the SEGMENT_LENGTH constant).
 * Each character in every segment is then converted to a different
 * character based on its segment number. This process is lossless,
 * deterministic, and reversible.
 *
 * Usage:
 * 	To cipher:
 *		cipher -input "Input text"
 *	To decipher:
 *		cipher -input "Ciphered text" -decipher
 *	To launch the web interface:
 *		cipher -web
 */

package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

const (
	//CHARS          = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890@#%^&*"
	CHARS          = "yF0%U2kBap4@KbzSh3imc5Or9EsVefdjGAIgxW#vtLM6R&X7Tw1oDC*PnZ8quYHlQNJ^" // this is the same as above but shuffled for randomness
	SEGMENT_LENGTH = 10
)

func main() {
	inputText := ""
	var decipher bool
	var web bool
	flag.StringVar(&inputText, "input", "", "input text to cipher")
	flag.BoolVar(&decipher, "decipher", false, "decipher input")
	flag.BoolVar(&web, "web", false, "launch web interface")
	flag.Parse()

	if web {
		http.HandleFunc("/", cipherWeb)
		if os.Getenv("PORT") != "" {
			fmt.Println("Listening on 0.0.0.0:" + os.Getenv("PORT") + "...")
			http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		} else {
			fmt.Println("Listening on 0.0.0.0:7777...")
			http.ListenAndServe(":7777", nil)
		}
	} else if decipher {
		fmt.Println(cipherInput(inputText, false))
	} else {
		fmt.Println(cipherInput(inputText, true))
	}
}

func cipherInput(input string, cipher bool) string {
	// split the input into segments with a max length of SEGMENT_LENGTH
	s := strings.Split(input, "")
	segments := make([][]string, 0)
	for i := 0; i < numberOfSegments(s); i++ {
		segments = append(segments, getSegment(s, i))
	}

	charList := strings.Split(CHARS, "")
	output := make([]string, 0)
	// loop through segments
	for segNum, segment := range segments {
		// through letters
		for _, letter := range segment {
			// and determine what output char each input letter should get
			if isInSlice(charList, letter) {
				for i, char := range charList {
					index := 0
					if cipher {
						// we are ciphering
						index = i + segNum + 1
						for index > len(charList)-1 {
							distance := index - len(charList)
							index = distance
						}
					} else {
						// we are deciphering
						index = i - segNum - 1
						for index < 0 {
							distance := -index
							index = len(charList) - distance
						}
					}
					if letter == char {
						output = append(output, charList[index])
					}
				}
			} else {
				output = append(output, letter)
			}
		}
	}

	return strings.Join(output, "")
}

func cipherWeb(w http.ResponseWriter, req *http.Request) {
	template, _ := template.ParseFiles("index.html")
	if req.Method == "POST" {
		req.ParseForm()
		input := req.FormValue("input")
		output := ""
		if _, exists := req.Form["cipher"]; exists {
			output = cipherInput(input, true)
		} else if _, exists := req.Form["decipher"]; exists {
			output = cipherInput(input, false)
		}
		template.Execute(w, output)
	} else {
		template.Execute(w, nil)
	}
}

// misc util stuff

func getSegment(s []string, num int) []string {
	curSegment := 0
	seg := make([]string, 0)
	for i, x := range s {
		if i%SEGMENT_LENGTH == 0 && i != 0 {
			curSegment++
		}

		if curSegment == num {
			seg = append(seg, x)
		} else if curSegment > num {
			break
		}
	}
	return seg
}

func isInSlice(slice []string, x string) bool {
	for _, y := range slice {
		if x == y {
			return true
		}
	}
	return false
}

func numberOfSegments(s []string) int {
	curSegment := 0
	for i, _ := range s {
		if i%SEGMENT_LENGTH == 0 {
			curSegment++
		}
	}
	return curSegment
}
