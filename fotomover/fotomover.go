package fotomover

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	SOURCEDIR = "/home/juergen/imago/fotoarchiv"
	DESTDIR   = "/home/juergen/imago/fablbum"
)

type FotomoverData struct {
	SourceDir               string
	DestDir                 string
	CurCalWeek              string
	LastWeekFilenames       []string // list with filenames for pictures from "last week"
	LastMonthFilenames      []string
	LastYearFilenames       []string
	SecondLastYearFilenames []string
	FiveYearPastFilenames   []string
	TenYearPastFilenames    []string

	FileList     []string
	TestdateList []string
}

func InitFotomoverData() (*FotomoverData, error) {
	fmd := &FotomoverData{
		SourceDir: SOURCEDIR,
		DestDir:   DESTDIR,
	}

	// check if dirs exist
	fi, err := os.Stat(fmd.SourceDir)
	if err != nil {
		return nil, fmt.Errorf("could not init fotomover: source dir does not exist %v: %v", fmd.SourceDir, err)
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("could not init fotomover: source dir is not a dir %v", fmd.SourceDir)
	}

	// check if Destination Dir exists - and create it if not
	fi, err = os.Stat(fmd.DestDir)
	if err != nil {
		// dir does not exist - create it:
		err = os.MkdirAll(fmd.DestDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("could not create destination dir %v: %v", fmd.DestDir, err)
		}
	} else {
		if !fi.IsDir() {
			// dir does not exist - create it:
			if err != nil {
				return nil, fmt.Errorf("could not create destination dir %v: %v", fmd.DestDir, err)
			}
		}
	}

	return fmd, nil
}

// MakeFilenames calculates filenames for lastweek, lastmonth... and populates Fields in fmd
func (fmd *FotomoverData) MakeFilenames() {
	today := time.Now()
	// _, curIsoWeek := today.ISOWeek()

	lastWeeksDay := today.AddDate(0, 0, -7)
	secondLastWeeksDay := today.AddDate(0, 0, -14)

	lastMonthDate := today.AddDate(0, -1, 0)
	secondLastMonthDate := today.AddDate(0, -2, 0)

	lastYearDate := today.AddDate(-1, 0, 0)
	secondLastYearDate := today.AddDate(-2, 0, 0)
	fithLastYearDate := today.AddDate(-5, 0, 0)
	tenthLastYearDate := today.AddDate(-10, 0, 0)

	fmt.Printf("got the following dates for today: %s\n", today.Format("2006-01-02"))

	fmt.Printf("\tlast weeks day: %s\n", lastWeeksDay.Format("2006-01-02"))
	fmt.Printf("\t2nd last weeks day: %s\n", secondLastWeeksDay.Format("2006-01-02"))

	fmt.Printf("\tlast month day: %s\n", lastMonthDate.Format("2006-01-02"))
	fmt.Printf("\t2nd last month day: %s\n", secondLastMonthDate.Format("2006-01-02"))

	fmt.Printf("\tlast year day: %s\n", lastYearDate.Format("2006-01-02"))
	fmt.Printf("\t2nd last year day: %s\n", secondLastYearDate.Format("2006-01-02"))
	fmt.Printf("\t5th last year day: %s\n", fithLastYearDate.Format("2006-01-02"))
	fmt.Printf("\t10th last year day: %s\n", tenthLastYearDate.Format("2006-01-02"))

	fmd.TestdateList = append(fmd.TestdateList, today.Format("2006-01-02"))
	fmd.TestdateList = append(fmd.TestdateList, lastWeeksDay.Format("2006-01-02"))
	fmd.TestdateList = append(fmd.TestdateList, secondLastWeeksDay.Format("2006-01-02"))
	fmd.TestdateList = append(fmd.TestdateList, lastMonthDate.Format("2006-01-02"))
	fmd.TestdateList = append(fmd.TestdateList, secondLastMonthDate.Format("2006-01-02"))
	fmd.TestdateList = append(fmd.TestdateList, lastYearDate.Format("2006-01-02"))
	fmd.TestdateList = append(fmd.TestdateList, secondLastYearDate.Format("2006-01-02"))
	fmd.TestdateList = append(fmd.TestdateList, fithLastYearDate.Format("2006-01-02"))
	fmd.TestdateList = append(fmd.TestdateList, tenthLastYearDate.Format("2006-01-02"))

	/*
		lastIsoWeek := curIsoWeek - 1
		if lastIsoWeek == 0 {
			lastIsoWeek = 52 // TODO: could be 52 or 53 - is on year...
		}
		lastMonday := FirstDayOfISOWeek(2023, lastIsoWeek, time.Local) // gets lastSunday
		lastMonday = lastMonday.AddDate(0, 0, 1)
		lastMonthDate := lastMonday.AddDate(0, -1, 0)

		_, lastMonthIsoWeek := lastMonthDate.ISOWeek()
		fmt.Printf("Last Monday: %s\n", lastMonday.Format("02.01.2006"))
		fmt.Printf("Last Month Date: %s\n", lastMonthDate.Format("02.01.2006"))
		fmt.Printf("Last month iso week: %d\n", lastMonthIsoWeek)
	*/

}

func (fmd *FotomoverData) GetFiles() []string {
	files := []string{}

	fmt.Printf("search for patterns: %v\n", fmd.TestdateList)

	filepath.WalkDir(fmd.SourceDir, func(path string, d fs.DirEntry, err error) error {
		fi, err := os.Stat(path)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		file := strings.ToLower(filepath.Base(path))
		fmt.Printf("checking file %s\n", file)
		if !strings.HasSuffix(file, ".jpg") {
			return nil
		}
		for _, pat := range fmd.TestdateList {
			if strings.HasPrefix(file, pat) {
				fmd.FileList = append(fmd.FileList, path)
				files = append(files, path)
				return nil
			}
		}
		return nil
	})

	return files
}

func FirstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()

	// iterate back to Monday
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the first week
	for isoYear < year {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the given week
	for isoWeek < week {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	return date
}
