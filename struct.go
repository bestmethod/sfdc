package sfdc

import "time"

// Reporting - basic object to init, before calling .Login() on it
type Reporting struct {
	Sandbox    bool   // true == sfdc test
	User       string // username to login with
	Password   string // password to login with
	SecToken   string // security token to use to login
	sessionId  string // this is the token/session key that we use to communicate with sfdc post login
	serverHost string // the serverHost parameter returns the hostname to use for further sfdc communication post login
	baseUrl    string // the baseUrl of the sfdc host for future communications post login
}

// GetReportInput - fill with details required to GetReport()
type GetReportInput struct {
	Debug          bool        // if set, will print what it's doing
	ReportId       string      // sfdc report ID
	DateFieldName  string      // name of the field to use for date chunks
	IncrementDays  int         // number of days to do at once each time
	StartDate      time.Time   // date from which to get report
	EndDate        time.Time   // date to which to get report
	startDate      time.Time   // this is used by chunks and incremented for each chunk while generating the report
	endDate        time.Time   // this is used by chunks and incremented for each chunk while generating the report
	metadata       []byte      // this is where we will store report metadata once retrieved
	dateFieldName  string      // this field is filled with actual real data field name from metadata
	metadataParsed *reportMeta // metadata parsed so we can use it
}

// GetReportOutput - the return data from GetReport()
type GetReportOutput struct {
	Rows        []GetReportRow `json:"rows"`
	ColumnNames []string       `json:"columnNames"`
}

type GetReportRow map[string]string

type metadataDetailColumnInfo struct {
	ReportExtendedMetadata struct {
		DetailColumnInfo map[string]struct {
			DataType string `json:"dataType"`
			Label    string `json:"label"`
		} `json:"detailColumnInfo"`
	} `json:"reportExtendedMetadata"`
}
