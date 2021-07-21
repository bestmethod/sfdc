package sfdc

type jsonOutput struct {
	Attributes struct {
		DescribeURL  string `json:"describeUrl"`
		InstancesURL string `json:"instancesUrl"`
		ReportID     string `json:"reportId"`
		ReportName   string `json:"reportName"`
		Type         string `json:"type"`
	} `json:"attributes"`
	AllData bool `json:"allData"`
	FactMap struct {
		TT struct {
			Aggregates []struct {
				Label string `json:"label"`
				Value int    `json:"value"`
			} `json:"aggregates"`
			Rows []struct {
				DataCells []struct {
					Label string      `json:"label"`
					Value interface{} `json:"value"`
				} `json:"dataCells"`
			} `json:"rows"`
		} `json:"T!T"`
	} `json:"factMap"`
	GroupingsAcross struct {
		Groupings []interface{} `json:"groupings"`
	} `json:"groupingsAcross"`
	GroupingsDown struct {
		Groupings []interface{} `json:"groupings"`
	} `json:"groupingsDown"`
	HasDetailRows  bool `json:"hasDetailRows"`
	PicklistColors struct {
	} `json:"picklistColors"`
	ReportExtendedMetadata struct {
		AggregateColumnInfo struct {
			RowCount struct {
				DataType string `json:"dataType"`
				Label    string `json:"label"`
			} `json:"RowCount"`
		} `json:"aggregateColumnInfo"`
		DetailColumnInfo map[string]struct {
			DataType string `json:"dataType"`
			Label    string `json:"label"`
		} `json:"detailColumnInfo"`
		GroupingColumnInfo struct {
		} `json:"groupingColumnInfo"`
	} `json:"reportExtendedMetadata"`
	ReportMetadata struct {
		Aggregates              []string      `json:"aggregates"`
		Chart                   interface{}   `json:"chart"`
		CrossFilters            []interface{} `json:"crossFilters"`
		Currency                interface{}   `json:"currency"`
		Description             interface{}   `json:"description"`
		DetailColumns           []string      `json:"detailColumns"`
		DeveloperName           string        `json:"developerName"`
		Division                interface{}   `json:"division"`
		FolderID                string        `json:"folderId"`
		GroupingsAcross         []interface{} `json:"groupingsAcross"`
		GroupingsDown           []interface{} `json:"groupingsDown"`
		HasDetailRows           bool          `json:"hasDetailRows"`
		HasRecordCount          bool          `json:"hasRecordCount"`
		HistoricalSnapshotDates []interface{} `json:"historicalSnapshotDates"`
		ID                      string        `json:"id"`
		Name                    string        `json:"name"`
		ReportBooleanFilter     interface{}   `json:"reportBooleanFilter"`
		ReportFilters           []struct {
			Column            string `json:"column"`
			IsRunPageEditable bool   `json:"isRunPageEditable"`
			Operator          string `json:"operator"`
			Value             string `json:"value"`
		} `json:"reportFilters"`
		ReportFormat string `json:"reportFormat"`
		ReportType   struct {
			Label string `json:"label"`
			Type  string `json:"type"`
		} `json:"reportType"`
		Scope          string `json:"scope"`
		ShowGrandTotal bool   `json:"showGrandTotal"`
		ShowSubtotals  bool   `json:"showSubtotals"`
		SortBy         []struct {
			SortColumn string `json:"sortColumn"`
			SortOrder  string `json:"sortOrder"`
		} `json:"sortBy"`
		StandardDateFilter struct {
			Column        string      `json:"column"`
			DurationValue string      `json:"durationValue"`
			EndDate       interface{} `json:"endDate"`
			StartDate     interface{} `json:"startDate"`
		} `json:"standardDateFilter"`
		StandardFilters []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"standardFilters"`
		SupportsRoleHierarchy   bool        `json:"supportsRoleHierarchy"`
		UserOrHierarchyFilterID interface{} `json:"userOrHierarchyFilterId"`
	} `json:"reportMetadata"`
}
