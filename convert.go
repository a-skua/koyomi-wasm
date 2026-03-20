package main

import (
	"fmt"
	"time"

	"github.com/a-skua/koyomi/internal/gen/a-skua/koyomi/convert"
	"github.com/goark/koyomi/value"
	"go.bytecodealliance.org/cm"
)

var (
	westernDates = map[cm.Rep]value.DateJp{}
	warekiDates  = map[cm.Rep]value.DateJp{}
	nextRep      cm.Rep
)

func newRep() cm.Rep {
	nextRep++
	return nextRep
}

func init() {
	convert.Exports.WesternDate.Constructor = westernDateConstructor
	convert.Exports.WesternDate.Destructor = westernDateDestructor
	convert.Exports.WesternDate.Year = westernDateYear
	convert.Exports.WesternDate.Month = westernDateMonth
	convert.Exports.WesternDate.Day = westernDateDay
	convert.Exports.WesternDate.ToString = westernDateToString
	convert.Exports.WesternDate.ToWareki = westernDateToWareki

	convert.Exports.WarekiDate.ExportConstructor = warekiDateConstructor
	convert.Exports.WarekiDate.Destructor = warekiDateDestructor
	convert.Exports.WarekiDate.Era = warekiDateEra
	convert.Exports.WarekiDate.Year = warekiDateYear
	convert.Exports.WarekiDate.Month = warekiDateMonth
	convert.Exports.WarekiDate.Day = warekiDateDay
	convert.Exports.WarekiDate.ToString = warekiDateToString
	convert.Exports.WarekiDate.ToSeireki = warekiDateToSeireki
}

// western-date

func westernDateConstructor(year int32, month convert.Month, day uint8) convert.WesternDate {
	t := time.Date(int(year), monthToTime(month), int(day), 0, 0, 0, 0, value.JST)
	d := value.NewDate(t)
	rep := newRep()
	westernDates[rep] = d
	return convert.WesternDateResourceNew(rep)
}

func westernDateDestructor(self cm.Rep) {
	delete(westernDates, self)
}

func westernDateYear(self cm.Rep) int32 {
	return int32(westernDates[self].Year())
}

func westernDateMonth(self cm.Rep) convert.Month {
	return timeToMonth(westernDates[self].Month())
}

func westernDateDay(self cm.Rep) uint8 {
	return uint8(westernDates[self].Day())
}

func westernDateToString(self cm.Rep) string {
	d := westernDates[self]
	return d.Format("2006-01-02")
}

func westernDateToWareki(self cm.Rep) cm.Result[string, convert.WarekiDate, string] {
	d := westernDates[self]
	era, year := d.YearEraString()
	convertEra, ok := eraStringToConvert(era)
	if !ok {
		return cm.Err[cm.Result[string, convert.WarekiDate, string]](fmt.Sprintf("unsupported era: %s", era))
	}

	_ = convertEra
	_ = year

	warekiRep := newRep()
	warekiDates[warekiRep] = d
	return cm.OK[cm.Result[string, convert.WarekiDate, string]](
		convert.WarekiDateResourceNew(warekiRep),
	)
}

// wareki-date

func warekiDateConstructor(era convert.Era, year int32, month convert.Month, day uint8) convert.WarekiDate {
	eraName := value.EraName(eraToJapanese(era))
	d := value.NewDateEra(eraName, int(year), monthToTime(month), int(day))
	rep := newRep()
	warekiDates[rep] = d
	return convert.WarekiDateResourceNew(rep)
}

func warekiDateDestructor(self cm.Rep) {
	delete(warekiDates, self)
}

func warekiDateEra(self cm.Rep) convert.Era {
	d := warekiDates[self]
	era, _ := d.YearEraString()
	e, _ := eraStringToConvert(era)
	return e
}

func warekiDateYear(self cm.Rep) int32 {
	d := warekiDates[self]
	_, year := d.YearEraString()
	return int32(yearToInt(year))
}

func warekiDateMonth(self cm.Rep) convert.Month {
	return timeToMonth(warekiDates[self].Month())
}

func warekiDateDay(self cm.Rep) uint8 {
	return uint8(warekiDates[self].Day())
}

func warekiDateToString(self cm.Rep) string {
	d := warekiDates[self]
	era, year := d.YearEraString()
	return fmt.Sprintf("%s%s年%d月%d日", era, year, d.Month(), d.Day())
}

func warekiDateToSeireki(self cm.Rep) cm.Result[string, convert.WesternDate, string] {
	d := warekiDates[self]
	if d.IsZero() {
		return cm.Err[cm.Result[string, convert.WesternDate, string]](fmt.Sprintf("invalid wareki date: %v", d))
	}

	westernRep := newRep()
	westernDates[westernRep] = d
	return cm.OK[cm.Result[string, convert.WesternDate, string]](
		convert.WesternDateResourceNew(westernRep),
	)
}

// mapping helpers

func monthToTime(m convert.Month) time.Month {
	switch m {
	case convert.MonthJanuary:
		return time.January
	case convert.MonthFebruary:
		return time.February
	case convert.MonthMarch:
		return time.March
	case convert.MonthApril:
		return time.April
	case convert.MonthMay:
		return time.May
	case convert.MonthJune:
		return time.June
	case convert.MonthJuly:
		return time.July
	case convert.MonthAugust:
		return time.August
	case convert.MonthSeptember:
		return time.September
	case convert.MonthOctober:
		return time.October
	case convert.MonthNovember:
		return time.November
	case convert.MonthDecember:
		return time.December
	default:
		return time.January
	}
}

func timeToMonth(m time.Month) convert.Month {
	switch m {
	case time.January:
		return convert.MonthJanuary
	case time.February:
		return convert.MonthFebruary
	case time.March:
		return convert.MonthMarch
	case time.April:
		return convert.MonthApril
	case time.May:
		return convert.MonthMay
	case time.June:
		return convert.MonthJune
	case time.July:
		return convert.MonthJuly
	case time.August:
		return convert.MonthAugust
	case time.September:
		return convert.MonthSeptember
	case time.October:
		return convert.MonthOctober
	case time.November:
		return convert.MonthNovember
	case time.December:
		return convert.MonthDecember
	default:
		return convert.MonthJanuary
	}
}

func eraStringToConvert(s string) (convert.Era, bool) {
	switch s {
	case "明治":
		return convert.EraMeiji, true
	case "大正":
		return convert.EraTaisho, true
	case "昭和":
		return convert.EraShowa, true
	case "平成":
		return convert.EraHeisei, true
	case "令和":
		return convert.EraReiwa, true
	default:
		return 0, false
	}
}

func eraToJapanese(e convert.Era) string {
	switch e {
	case convert.EraMeiji:
		return "明治"
	case convert.EraTaisho:
		return "大正"
	case convert.EraShowa:
		return "昭和"
	case convert.EraHeisei:
		return "平成"
	case convert.EraReiwa:
		return "令和"
	default:
		return ""
	}
}

func yearToInt(s string) int {
	if s == "元" {
		return 1
	}
	n := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			n = n*10 + int(r-'0')
		}
	}
	return n
}
