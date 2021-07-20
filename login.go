package sfdc

import (
	"crypto/tls"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// sfdc.Login - call on Reporting struct to get auth token from sfdc
func (sfdc *Reporting) Login() error {

	requestBody := loginGetRequestBody(html.EscapeString(sfdc.User), html.EscapeString(sfdc.Password), sfdc.SecToken)

	ua := "login"
	if sfdc.Sandbox {
		ua = "test"
	}

	req, err := http.NewRequest("POST", "https://"+ua+".salesforce.com/services/Soap/u/v38.0", strings.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("could not create http request: %s", err)
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
	req.Header.Set("content-type", "text/xml")
	req.Header.Set("charset", "UTF-8")
	req.Header.Set("SOAPAction", "login")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to login with sfdc: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to login with sfdc: response read failure: %s", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		//<faultcode>INVALID_LOGIN</faultcode><faultstring>INVALID_LOGIN: Invalid username, password, security token; or user locked out.</faultstring>
		re, err := regexp.Compile("<faultcode>(.*)</faultcode>")
		if err != nil {
			return fmt.Errorf("regex compile error: %s", err)
		}
		faultCode := ""
		answer := re.FindStringSubmatch(string(body))
		if len(answer) >= 2 {
			faultCode = answer[1]
		}
		re, err = regexp.Compile("<faultstring>(.*)</faultstring>")
		if err != nil {
			return fmt.Errorf("regex compile error: %s", err)
		}
		faultString := ""
		answer = re.FindStringSubmatch(string(body))
		if len(answer) >= 2 {
			faultString = answer[1]
		}
		if faultCode == "" && faultString == "" {
			faultString = string(body)
		}
		return fmt.Errorf("server returned %d: faultCode=%s faultString=%s", resp.StatusCode, faultCode, faultString)

	}

	re, err := regexp.Compile("<sessionId>(.*)</sessionId>")
	if err != nil {
		return fmt.Errorf("regex compile error: %s", err)
	}
	answer := re.FindStringSubmatch(string(body))
	if len(answer) < 2 {
		return fmt.Errorf("sessionid not found in response")
	}
	sfdc.sessionId = answer[1]

	re, err = regexp.Compile("<serverUrl>(.*)</serverUrl>")
	if err != nil {
		return fmt.Errorf("regex compile error: %s", err)
	}
	answer = re.FindStringSubmatch(string(body))
	if len(answer) < 2 {
		return fmt.Errorf("serverUrl not found in response")
	}
	u, err := url.Parse(answer[1])
	if err != nil {
		return fmt.Errorf("could not parse serverUrl `%s`: %s", answer[1], err)
	}
	sfdc.serverHost = u.Host

	sfdc.baseUrl = fmt.Sprintf("https://%s/services/data/v38.0/analytics", sfdc.serverHost)
	return nil
}

// loginGetRequestBody - generate xml request to login
func loginGetRequestBody(username, password, token string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8" ?>
	<env:Envelope
			xmlns:xsd="http://www.w3.org/2001/XMLSchema"
			xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
			xmlns:env="http://schemas.xmlsoap.org/soap/envelope/">
		<env:Body>
			<n1:login xmlns:n1="urn:partner.soap.sforce.com">
				<n1:username>%s</n1:username>
				<n1:password>%s%s</n1:password>
			</n1:login>
		</env:Body>
	</env:Envelope>`, username, password, token)
}

/*
<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns="urn:partner.soap.sforce.com" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<soapenv:Body>
		<loginResponse>
			<result>
				<passwordExpired>false</passwordExpired>
				<serverUrl>https://xxx.my.salesforce.com/services/Soap/u/7.0/xxxx</serverUrl>
				<sessionId>xxx</sessionId>
				<userId>xxx</userId>
				<userInfo>
					<accessibilityMode>false</accessibilityMode>
					<currencySymbol>$</currencySymbol>
					<organizationId>xxx</organizationId>
					<organizationMultiCurrency>false</organizationMultiCurrency>
					<organizationName>xxx</organizationName>
					<userDefaultCurrencyIsoCode xsi:nil="true"/>
					<userEmail>xxx</userEmail>
					<userFullName>xxx</userFullName>
					<userId>xxx</userId>
					<userLanguage>en_US</userLanguage>
					<userLocale>en_US</userLocale>
					<userTimeZone>Europe/Lisbon</userTimeZone>
					<userUiSkin>Theme3</userUiSkin>
				</userInfo>
			</result>
		</loginResponse>
	</soapenv:Body>
</soapenv:Envelope>
*/
