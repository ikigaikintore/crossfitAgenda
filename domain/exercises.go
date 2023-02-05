package domain

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ExerciseType string
type ExercisesTypes []ExerciseType

const (
	LegsExercise      ExerciseType = "legs"
	ShouldersExercise              = "shoulders"
	CoreExercise                   = "core"
	CardioExercise                 = "cardio"
	ArmsExercise                   = "arms"
	PecsExercise                   = "pecs"
)

type MonthWod string
type MonthsWod []MonthWod

const (
	January   MonthWod = "January"
	February           = "February"
	March              = "March"
	April              = "April"
	May                = "May"
	June               = "June"
	July               = "July"
	August             = "August"
	September          = "September"
	October            = "October"
	November           = "November"
	Nobember           = "Nobember"
	December           = "December"
)

var months = MonthsWod{January, February, March, April, May, June, July, August, September, October, November, Nobember, December}

func (mns MonthsWod) toMap() map[MonthWod]struct{} {
	m := make(map[MonthWod]struct{})
	for _, month := range months {
		m[month] = struct{}{}
	}
	return m
}

func (mn MonthWod) toMonth() time.Month {
	switch mn {
	case January:
		return time.January
	case February:
		return time.February
	case March:
		return time.March
	case April:
		return time.April
	case May:
		return time.May
	case June:
		return time.June
	case July:
		return time.July
	case August:
		return time.August
	case September:
		return time.September
	case October:
		return time.October
	case November, Nobember:
		return time.November
	case December:
		return time.December
	default:
		return -1
	}
}

type ExerciseName string

func (e ExerciseName) String() string {
	return string(e)
}

const (
	HandStand       ExerciseName = "Handstand"
	CleanAndJerk                 = "C&J"
	OHS                          = "OHS"
	BUS                          = "BUS"
	HandStandPushUp              = "HSPU"
	DeadLift                     = "DL"
	FrontSquat                   = "Front Squat"
	StrictPress                  = "Strict Press"
	Lunge                        = "Lunge"
	BackStrength                 = "Back Strength"
	BackSquat                    = "Back Squat"
	Snatch                       = "Snatch"
	BenchPress                   = "Bench Press"
	StrictPullUp                 = "Strict Pull Up"
	OTM                          = "OTM"
	PowerSnatch                  = "Power Snatch"
	HIIT                         = "HIIT"
	Clean                        = "Clean"
	Core                         = "Core"
	Ring                         = "Ring"
	PushPress                    = "Push Press"
	Jerk                         = "Jerk"
	Gymnastics                   = "Gymnastics"
	UpperBody                    = "Upper Body"
	Thruster                     = "Thruster"
)

var listExercises = []ExerciseName{
	HandStand,
	Clean,
	CleanAndJerk,
	OTM,
	BUS,
	BackSquat,
	BackStrength,
	Snatch,
	BenchPress,
	PowerSnatch,
	PushPress,
	HIIT,
	Core,
	Ring,
	Lunge,
	StrictPress,
	FrontSquat,
	DeadLift,
	HandStandPushUp,
	StrictPullUp,
	Jerk,
	Gymnastics,
	UpperBody,
	OHS,
	Thruster,
}

var patterns = []string{`^([1-9]|[12]\d|3[01])$`}

type raw struct {
	data        []byte
	len         uint64
	rgxPatterns []*regexp.Regexp
}

type RawProcessor interface {
	Convert() MonthWodExercises
}

func NewRawProcessor(text string) RawProcessor {
	rgx := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		rgx[i] = regexp.MustCompile(pattern)
	}

	return &raw{
		data:        []byte(text),
		len:         uint64(len(text)),
		rgxPatterns: rgx,
	}
}

func (r raw) prepareMonthWod(monthTweet MonthWod) MonthWodExercises {
	month := make(MonthWodExercises)
	var firstDayMonth time.Time
	var firstDayNextMonth time.Time

	{
		loc, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Now()
		diff := monthTweet.toMonth() - now.Month()
		firstDayMonth = time.Date(now.Year(), now.Month()+diff, 1, 0, 0, 0, 0, loc)
		firstDayNextMonth = time.Date(firstDayMonth.Year(), firstDayMonth.Month()+1, 1, 0, 0, 0, 0, loc)
	}

	for {
		if firstDayMonth.Month() == firstDayNextMonth.Month() {
			break
		}
		month[firstDayMonth.Day()] = &wod{
			day:     firstDayMonth.Day(),
			rawDate: firstDayMonth,
			month:   firstDayMonth.Month(),
			year:    firstDayMonth.Year(),
		}
		firstDayMonth = firstDayMonth.AddDate(0, 0, 1)
	}
	return month
}

func (r raw) Convert() MonthWodExercises {
	/*
		send workers from last line to beginning until it finds the weekend days
	*/
	sentences := strings.Split(string(r.data), "\n")
	mapMonths := months.toMap()
	var monthTweet MonthWod
	for _, sentence := range sentences {
		if _, exists := mapMonths[MonthWod(sentence)]; exists {
			monthTweet = MonthWod(sentence)
			break
		}
	}
	weeks := r.prepareMonthWod(monthTweet)

	exercisesExistsMap := make(map[ExerciseName]struct{})
	for _, v := range listExercises {
		if _, exists := exercisesExistsMap[v]; exists {
			continue
		}
		exercisesExistsMap[v] = struct{}{}
	}

	listNodeExercises := newListNodes()
	for i := 0; i < len(sentences)-1; i++ {
		if r.omitIfContains(sentences[i+1]) {
			sentences = append(sentences[:i+1], sentences[i+2:]...)
		}
		if r.wodValidName(sentences[i+1], exercisesExistsMap) && r.wodValidDay(sentences[i]) {
			listNodeExercises.Insert(strings.TrimSpace(sentences[i+1]), strings.TrimSpace(sentences[i]))
		}
	}

	for {
		if listNodeExercises.IsEmpty() {
			break
		}

		wodExercise := listNodeExercises.Get()

		day, err := strconv.Atoi(wodExercise.element.day)
		if err != nil {
			log.Printf("incorrect day: %s", err.Error())
		}

		name := strings.TrimSpace(wodExercise.element.name)
		if _, exist := exercisesExistsMap[ExerciseName(name)]; !exist {
			continue
		}

		weeks[day].name = ExerciseName(name)
	}
	return weeks
}

func (r raw) wodValidDay(day string) bool {
	for _, p := range r.rgxPatterns {
		if p.MatchString(day) {
			return true
		}
	}
	return false
}

func (r raw) wodValidName(name string, exercisesExistsMap map[ExerciseName]struct{}) bool {
	if _, exists := exercisesExistsMap[ExerciseName(name)]; !exists {
		return false
	}
	return true
}

func (r raw) omitIfContains(str string) bool {
	omitList := map[string]struct{}{
		"RM":   {},
		"PAMA": {},
		"BA":   {},
	}
	if _, exists := omitList[str]; exists {
		return true
	}
	return false
}

type wod struct {
	style   ExercisesTypes
	name    ExerciseName
	day     int
	month   time.Month
	year    int
	rawDate time.Time
}

func (w *wod) ExerciseName() ExerciseName {
	return w.name
}

func (w *wod) Day() int {
	return w.day
}

func (w *wod) Month() time.Month {
	return w.month
}

type MonthWodExercises map[int]*wod
