package sfdc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bestmethod/inslice"
)

func (input *GetReportInput) logDebug(format string, values ...interface{}) {
	if !input.Debug {
		return
	}
	if len(values) > 0 {
		log.Printf(format, values...)
	} else {
		log.Print(format)
	}
}

// Get Report - get sfdc pre-made report with extra constraint details
func (sfdc *Reporting) GetReport(input *GetReportInput) (output *GetReportOutput, err error) {
	if err := sfdc.getReportSanitise(input); err != nil {
		return nil, err
	}

	output = new(GetReportOutput)

	input.logDebug("Getting Report Metadata")
	input.metadata, err = sfdc.getReportMetadata(input)
	if err != nil {
		return nil, err
	}

	input.logDebug("Getting DateField real name from metadata")
	err = sfdc.getReportMetadataDateFieldName(input)
	if err != nil {
		return nil, err
	}
	input.logDebug("Value dataFieldName=%s", input.dateFieldName)

	input.logDebug("Parsing metadata to prepare for chunk reporting")
	err = sfdc.getReportMetadataParse(input)
	if err != nil {
		return nil, err
	}

	input.startDate = input.StartDate
	chunks := 0
	for !input.startDate.After(input.EndDate) {
		chunks++
		input.endDate = input.startDate.AddDate(0, 0, input.IncrementDays)
		if input.endDate.After(input.EndDate) {
			input.endDate = input.EndDate
		}
		input.startDate = input.endDate.AddDate(0, 0, 1)
	}

	input.logDebug("Chunks=%d StartDate=%s EndDate=%s", chunks, input.StartDate.Format("2006-01-02"), input.EndDate.Format("2006-01-02"))
	input.startDate = input.StartDate
	chunk := 0
	for !input.startDate.After(input.EndDate) {
		chunk++
		input.endDate = input.startDate.AddDate(0, 0, input.IncrementDays)
		if input.endDate.After(input.EndDate) {
			input.endDate = input.EndDate
		}
		input.logDebug("Chunk=%d StartDate=%s EndDate=%s", chunk, input.startDate.Format("2006-01-02"), input.endDate.Format("2006-01-02"))
		rows, err := sfdc.getReportChunk(input)
		if err != nil {
			return output, err
		}
		output.Rows = append(output.Rows, rows.Rows...)
		for _, colName := range rows.ColumnNames {
			if !inslice.HasString(output.ColumnNames, colName) {
				output.ColumnNames = append(output.ColumnNames, colName)
			}
		}

		input.startDate = input.endDate.AddDate(0, 0, 1)
	}
	input.logDebug("Done")
	return output, nil
}

func (sfdc *Reporting) getReportSanitise(input *GetReportInput) (err error) {
	if input == nil {
		return errors.New("invalid request, GetReportInput is nil")
	}
	if input.ReportId == "" || input.IncrementDays == 0 || input.DateFieldName == "" {
		return errors.New("report ID, increment days and date field name are mandatory")
	}
	if input.StartDate.After(input.EndDate) {
		return errors.New("start date cannot be after end date")
	}
	return nil
}

func (sfdc *Reporting) getReportMetadata(input *GetReportInput) (metadata []byte, err error) {
	metadataUrl := fmt.Sprintf("%s/reports/%s", sfdc.baseUrl, input.ReportId)
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("OAuth %s", sfdc.sessionId)
	resp, err := httpGet(metadataUrl+"/describe", nil, headers)
	if err != nil {
		return nil, fmt.Errorf("could not get metadata: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("metadata-get server returned %d: %s", resp.StatusCode, string(body))
	}
	if err != nil {
		return nil, fmt.Errorf("metadata-get read-body: %s", err)
	}
	return body, nil
}

func (sfdc *Reporting) getReportMetadataDateFieldName(input *GetReportInput) (err error) {
	meta := new(metadataDetailColumnInfo)
	err = json.Unmarshal(input.metadata, meta)
	if err != nil {
		return err
	}
	for n, v := range meta.ReportExtendedMetadata.DetailColumnInfo {
		if v.Label == input.DateFieldName {
			input.dateFieldName = n
			return nil
		}
	}
	return errors.New("date field name not found in report metadata")
}

func (sfdc *Reporting) getReportMetadataParse(input *GetReportInput) (err error) {
	input.metadataParsed = new(reportMeta)
	err = json.Unmarshal(input.metadata, input.metadataParsed)
	if err != nil {
		return err
	}
	return nil
}

func (sfdc *Reporting) getReportChunk(input *GetReportInput) (output *GetReportOutput, err error) {
	reportUrl := fmt.Sprintf("%s/reports/%s?includeDetails=true", sfdc.baseUrl, input.ReportId)
	input.metadataParsed.ReportMetadata.StandardDateFilter.Column = input.dateFieldName
	input.metadataParsed.ReportMetadata.StandardDateFilter.DurationValue = "CUSTOM"
	input.metadataParsed.ReportMetadata.StandardDateFilter.EndDate = input.endDate.Format("2006-01-02")
	input.metadataParsed.ReportMetadata.StandardDateFilter.StartDate = input.startDate.Format("2006-01-02")
	meta, err := json.Marshal(input.metadataParsed)
	if err != nil {
		return nil, err
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("OAuth %s", sfdc.sessionId)
	headers["Content-Type"] = "application/json"
	resp, err := httpPost(reportUrl, meta, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("chunk-get server returned %d: %s", resp.StatusCode, string(body))
	}
	if err != nil {
		return nil, fmt.Errorf("chunk-get read-body: %s", err)
	}
	result := new(jsonOutput)
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	if !result.AllData {
		return nil, errors.New("report returned partial data - too many results; lower the date-increment")
	}
	output = new(GetReportOutput)
	for _, colName := range result.ReportMetadata.DetailColumns {
		output.ColumnNames = append(output.ColumnNames, result.ReportExtendedMetadata.DetailColumnInfo[colName].Label)
	}
	for _, row := range result.FactMap.TT.Rows {
		rowData := make(GetReportRow)
		for colNo, col := range row.DataCells {
			rowData[output.ColumnNames[colNo]] = col.Label
		}
		output.Rows = append(output.Rows, rowData)
	}
	return output, nil
}

func httpGet(baseUrl string, requestBody []byte, headers map[string]string) (resp *http.Response, err error) {
	return httpRequest(baseUrl, "GET", requestBody, headers)
}

func httpPost(baseUrl string, requestBody []byte, headers map[string]string) (resp *http.Response, err error) {
	return httpRequest(baseUrl, "POST", requestBody, headers)
}

// requestType == POST||GET
func httpRequest(baseUrl string, requestType string, requestBody []byte, headers map[string]string) (body *http.Response, err error) {
	var req *http.Request
	if len(requestBody) != 0 {
		req, err = http.NewRequest(requestType, baseUrl, bytes.NewReader(requestBody))
	} else {
		req, err = http.NewRequest(requestType, baseUrl, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("could not create http request: %s", err)
	}
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: 5 * time.Second,
		IdleConnTimeout:     5 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}
	for hKey, hVal := range headers {
		req.Header.Set(hKey, hVal)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to login with sfdc: %s", err)
	}
	return resp, nil
}
