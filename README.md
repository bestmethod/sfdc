# Download large salesforce reports

## Download large slaesforce reports by chunking them by date

The project runs existing salesforce reports, adding date filters to chunk the report into smaller portions. The data is then reassembled and returned.

## Prerequisites:

* create a report in salesforce, selecting which fields you want exported
* in the report specify your constraints (filters) - basically the report as you want it, without the date filters
* save the report and take it's report-id (easiest way is to use the URL)
* get a security token for API SOAP logins from salesforce for your user

## GetReport Input and Output structs

```go
type GetReportInput struct {
	Debug          bool        // if set, will print what it's doing
	ReportId       string      // sfdc report ID
	DateFieldName  string      // name of the field to use for date chunks
	IncrementDays  int         // number of days to do at once each time
	StartDate      time.Time   // date from which to get report
	EndDate        time.Time   // date to which to get report
}

type GetReportOutput struct {
	Rows        []GetReportRow `json:"rows"`        // a list of rows with their colums in each row
	ColumnNames []string       `json:"columnNames"` // a list of column names present; does not guarantee all columns exist in each row, user must check this for each row when parsing
}

type GetReportRow map[string]string // definition of a row: map[column-name]column-value
```

## Most basic usage:

Create reporting object with login details

```go
reports := &sfdc.Reporting{
	Sandbox:  false,
	User:     "...",
	Password: "...",
	SecToken: "...",
}
```

Login

```go
err := reports.Login()
if err != nil {
	log.Fatal(err)
}
fmt.Println(reports.String())
```

Get the report - example - get last month, but do not get today (isn't full day finished yet), chunk by 7-day increments
NOTE: the start/end dates are INCLUSIVE, so the day present in endDate will also be exported

```go
out, err := reports.GetReport(&sfdc.GetReportInput{
	ReportId:      "...",
	DateFieldName: "...",
	IncrementDays: 7,
	StartDate:     time.Now().AddDate(0, -1, -1),
	EndDate:       time.Now().AddDate(0, 0, -1),
	Debug:         true,
})
if err != nil {
	log.Fatal(err)
}
```

Just print the result

```go
fmt.Println(out.String())
```

## Basic CSV example:

Basic example for getting a report for a selected date range, chunking it into 7-day increments for each API call, and storing the result in a CSV file.

```go
package main

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/bestmethod/sfdc"
)

func main() {
	// create object with login details
	reports := &sfdc.Reporting{
		Sandbox:  false,
		User:     "...",
		Password: "...",
		SecToken: "...",
	}

	// login
	err := reports.Login()
	if err != nil {
		log.Fatalf("sfdc login failed: %s", err)
	}

	// create csv file for writing
	fd, err := os.OpenFile("out.csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed to open csv for writing: %s", err)
	}
	defer fd.Close()

	// open csv writer
	csvWriter := csv.NewWriter(fd)
	defer csvWriter.Flush()

	// get the report - example - get last month - do not get today (isn't full day), chunk by 7-day increments
	// NOTE: the start/end dates are INCLUSIVE, so the day present in endDate will also be exported
	out, err := reports.GetReport(&sfdc.GetReportInput{
		ReportId:      "...",
		DateFieldName: "...",
		IncrementDays: 7,
		StartDate:     time.Now().AddDate(0, -1, -1),
		EndDate:       time.Now().AddDate(0, 0, -1),
		Debug:         true,
	})
	if err != nil {
		log.Fatalf("failed to get report: %s", err)
	}

	// write out csv header row
	err = csvWriter.Write(out.ColumnNames)
	if err != nil {
		log.Fatalf("failed to write csv header row: %s", err)
	}

	// write data rows
	for _, row := range out.Rows {
		col := []string{}
		for _, colName := range out.ColumnNames {
			colVal, ok := row[colName]
			if !ok {
				colVal = ""
			}
			col = append(col, colVal)
		}
		err = csvWriter.Write(col)
		if err != nil {
			log.Fatalf("failed to write csv rows: %s", err)
		}
	}
}
```

## Full Usage Example:

Example includes a yaml configuration file, support for multiple reports and date ranges, storing to csv files.

### config.yml

```yaml
connect:
  user: "some@email.com"
  password: "user-password"
  secToken: "wefhhuih934"
  sandbox: false
reports:
  - fileName: "cases.csv"
    reportId: "00O1x000001z1xDIEF"
    dateFieldName: "Last Modified"
    incrementDays: 28
    startDate: "2021-07-01"
    endDate: "2021-07-31"
    debug: true
  - fileName: "caseHistory.csv"
    reportId: "00O1x000001z1xEIGE"
    dateFieldName: "Edit Date"
    incrementDays: 28
    startDate: "2021-07-01"
    endDate: "2021-07-31"
    debug: true
```

### main code

```go
package main

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/bestmethod/sfdc"
)

type Salesforce struct {
	Connect struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SecToken string `yaml:"secToken"`
		Sandbox  bool   `yaml:"sandbox"`
	} `yaml:"connect"`
	Reports []struct {
		FileName      string `yaml:"fileName"`
		ReportID      string `yaml:"reportId"`
		DateFieldName string `yaml:"dateFieldName"`
		IncrementDays int    `yaml:"incrementDays"`
		StartDate     string `yaml:"startDate"`
		EndDate       string `yaml:"endDate"`
		Debug         bool   `yaml:"debug"`
		startDate     time.Time
		endDate       time.Time
	} `yaml:"reports"`
	reports *sfdc.Reporting
}

func main() {
	s := new(Salesforce)
	conf, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Printf("ERROR opening config.yml: %s", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(conf, s)
	if err != nil {
		log.Printf("ERROR loading config.yml: %s", err)
		os.Exit(1)
	}

	s.reports = &sfdc.Reporting{
		Sandbox:  false,
		User:     s.Connect.User,
		Password: s.Connect.Password,
		SecToken: s.Connect.SecToken,
	}
	err = s.reports.Login()
	if err != nil {
		log.Printf("sfdc login failed: %s", err)
		os.Exit(1)
	}
	for i := range s.Reports {
		s.Reports[i].startDate, err = time.Parse("2006-01-02", s.Reports[i].StartDate)
		if err != nil {
			log.Printf("start date wrong format: %s", err)
			os.Exit(1)
		}
		s.Reports[i].endDate, err = time.Parse("2006-01-02", s.Reports[i].EndDate)
		if err != nil {
			log.Printf("end date wrong format: %s", err)
			os.Exit(1)
		}
		err = s.getReport(i)
		if err != nil {
			log.Printf("sfdc reporting failed: %s", err)
			os.Exit(1)
		}
	}
}

func (s *Salesforce) getReport(reportNo int) error {
	fd, err := os.OpenFile(s.Reports[reportNo].FileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()
	csvWriter := csv.NewWriter(fd)
	defer csvWriter.Flush()

	out, err := s.reports.GetReport(&sfdc.GetReportInput{
		ReportId:      s.Reports[reportNo].ReportID,
		DateFieldName: s.Reports[reportNo].DateFieldName,
		IncrementDays: s.Reports[reportNo].IncrementDays,
		StartDate:     s.Reports[reportNo].startDate,
		EndDate:       s.Reports[reportNo].endDate,
		Debug:         s.Reports[reportNo].Debug,
	})
	if err != nil {
		return err
	}

	err = csvWriter.Write(out.ColumnNames)
	if err != nil {
		return err
	}

	for _, row := range out.Rows {
		col := []string{}
		for _, colName := range out.ColumnNames {
			colVal, ok := row[colName]
			if !ok {
				colVal = ""
			}
			col = append(col, colVal)
		}
		err = csvWriter.Write(col)
		if err != nil {
			return err
		}
	}

	return nil
}
```
