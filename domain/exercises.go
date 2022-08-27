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
	ShouldersExercise ExerciseType = "shoulders"
	CoreExercise      ExerciseType = "core"
	CardioExercise    ExerciseType = "cardio"
	ArmsExercise      ExerciseType = "arms"
	PecsExercise      ExerciseType = "pecs"
)

type ExerciseName string

func (e ExerciseName) String() string {
	return string(e)
}

const (
	CleanAndJerk    ExerciseName = "C&J"
	OverHeadSquat   ExerciseName = "OHS"
	BUS             ExerciseName = "BUS"
	HandStandPushUp ExerciseName = "HSPU"
	DeadLift        ExerciseName = "DL"
	FrontSquat      ExerciseName = "Front Squat"
	StrictPress     ExerciseName = "Strict Press"
	Lunge           ExerciseName = "Lunge"
	BackStrength    ExerciseName = "Back Strength"
	BackSquat       ExerciseName = "Back Squat"
	Snatch          ExerciseName = "Snatch"
	BenchPress      ExerciseName = "Bench Press"
	StrictPullUp    ExerciseName = "Strict Pull Up"
	OTM             ExerciseName = "OTM"
	PowerSnatch     ExerciseName = "Power Snatch"
	HIIT            ExerciseName = "HIIT"
	Clean           ExerciseName = "Clean"
	Core            ExerciseName = "Core"
	Ring            ExerciseName = "Ring"
	PushPress       ExerciseName = "Push Press"
	Jerk            ExerciseName = "Jerk"
	Gymnastics      ExerciseName = "Gymnastics"
)

var listExercises = []ExerciseName{
	Clean,
	CleanAndJerk,
	OTM,
	OverHeadSquat,
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

func (r raw) prepareMonthWod() MonthWodExercises {
	month := make(MonthWodExercises)
	var firstDayMonth time.Time
	var actualMonth time.Month

	{
		loc, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Now()
		firstDayMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		actualMonth = firstDayMonth.Month()
	}

	for {
		if firstDayMonth.Month() != actualMonth {
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
	weeks := r.prepareMonthWod()

	exercisesExistsMap := make(map[ExerciseName]struct{})
	for _, v := range listExercises {
		if _, exists := exercisesExistsMap[v]; exists {
			continue
		}
		exercisesExistsMap[v] = struct{}{}
	}

	sentences := strings.Split(string(r.data), "\n")
	listNodeExercises := newListNodes()
	for i := 0; i < len(sentences)-1; i++ {
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

		/*
			1 2 3 4 5 6 7
			0 1 2 3 4 5 6

			8 9 10 11 12 13 14
			0 1 2  3  4  5  6

			15 16 17 18 19 20 21
			0  1  2  3  4  5  6
		*/

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

type MonthWodExercises map[int]*wod
