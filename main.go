package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Cell interface {
	Execute() ([]string, int64)
}

type RatingS struct {
	Rating  float64
	CellObj Cell
}

func (r *RatingS) Valid(value, begin float64) bool {
	if value >= begin && value <= r.Rating+begin {
		return true
	}
	return false
}

type S struct {
	Data []RatingS
	name string
}

func (s *S) Execute() ([]string, int64) {
	value := rand.Float64()
	var begin float64
	begin = 0
	var data []string
	var outcome int64
	for _, rating := range s.Data {
		if rating.Valid(value, begin) {
			data, outcome = rating.CellObj.Execute()
			break
		}
		begin += rating.Rating
	}
	return append([]string{s.name}, data...), outcome
}

type Number struct {
	outcome int64
}

func (n *Number) Execute() ([]string, int64) {
	return []string{}, n.outcome
}

type Language struct {
	s        []*S
	outcomes []*Number
}

func (l *Language) ExecuteFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l.Execute(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (l *Language) HandleError(err error) {
	fmt.Println(err)
}

func (l *Language) init(line string) {
	if len(line) <= 8 {
		fmt.Println(line, "Line incorrect")
		return
	}
	min, err := strconv.ParseInt(string(line[5]), 10, 32)
	if err != nil {
		l.HandleError(err)
		return
	}
	max, err := strconv.ParseInt(string(line[8]), 10, 32)
	if err != nil {
		l.HandleError(err)
		return
	}
	index := int64(0)
	for i := min + 1; i <= max; i++ {
		if string(line[0]) == "s" {
			s := l.NewS(string(i))
			s.name = strconv.FormatInt(index, 10)
			l.s = append(l.s, s)
		} else {
			l.outcomes = append(l.outcomes, l.NewD(i))

		}
		index++
	}
}

func (l *Language) NewS(name string) *S {
	return &S{name: name}
}

func (l *Language) NewD(outcome int64) *Number {
	return &Number{outcome: outcome}
}

func (l *Language) handleAttr(s *S, content string, RatingValue float64) {
	in := 0
	next := true
	sp := strings.Split(content, "&")
	for next {
		if in > len(sp) {
			fmt.Println("Syntax not correct")
			return
		}
		value := sp[in]
		gotoIndex, err := strconv.ParseInt(string(value[7]), 10, 32)
		if err != nil {
			l.HandleError(err)
			return
		}
		if string(value[2]) == "s" {
			if gotoIndex < int64(len(l.s)) {
				var rating RatingS
				rating.CellObj = l.s[gotoIndex]
				rating.Rating = RatingValue
				s.Data = append(s.Data, rating)
				next = false
				return
			}
			in++
		} else {
			gotoIndex--
			if gotoIndex < int64(len(l.outcomes)) {
				var rating RatingS
				rating.CellObj = l.outcomes[gotoIndex]
				rating.Rating = RatingValue
				s.Data = append(s.Data, rating)
				next = false
				return
			}
			in++
		}
	}
}

func (l *Language) rowSetup(line string) {
	if !strings.Contains(line, ":") {
		return
	}
	index, err := strconv.ParseInt(string(line[2]), 10, 32)
	if err != nil {
		l.HandleError(err)
		return
	}
	s := l.s[index]
	// value = line
	data := strings.Split(line, "+")
	data[0] = strings.Split(data[0], "->")[1]
	for _, content := range data {
		newContent := strings.Split(content, ":")
		newContent[0] = strings.ReplaceAll(newContent[0], " ", "")
		rating, err := strconv.ParseFloat(newContent[0], 64)
		if err != nil {
			l.HandleError(err)
			return
		}
		l.handleAttr(s, newContent[1], rating)
	}
	var count float64
	for _, ratings := range s.Data {
		count += ratings.Rating
	}
	if count != 1.0 {
		fmt.Println("Rating doesn't count up to 1")
		os.Exit(0)
	}

}

func (l *Language) Execute(line string) {
	if len(line) <= 2 {
		return
	}
	if line[:2] == "//" {
		return
	}

	if strings.Contains(line, "init") {
		l.init(line)
		return
	}
	l.rowSetup(line)

}

func main() {
	// Setup
	var l Language
	l.ExecuteFile("2.28/script.txt")
	fmt.Println(l.s[0].Execute())
}
