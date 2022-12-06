package faareg

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

type Registration struct {
	Aircraft        Aircraft        `json:"aircraft"`
	RegisteredOwner RegisteredOwner `json:"registered_owner"`
	Airworthiness   Airworthiness   `json:"airworthiness"`
}

type Aircraft struct {
	SerialNumber         string `json:"serial_number"`
	Status               string `json:"status"`
	ManufacturerName     string `json:"manufacturer_name"`
	CertificateIssueDate string `json:"certificate_issue_date"`
	Model                string `json:"model"`
	ExpirationDate       string `json:"expiration_date"`
	AircraftType         string `json:"aircraft_type"`
	EngineType           string `json:"engine_type"`
	PendingNumberChange  string `json:"pending_number_change"`
	Dealer               string `json:"dealer"`
	DateChangeAuthorized string `json:"date_change_authorized"`
	ModeSCodeOct         string `json:"mode_s_code_oct"`
	MfrYear              string `json:"mfr_year"`
	ModeSCodeHex         string `json:"mode_s_code_hex"`
	TypeRegistration     string `json:"type_registration"`
	FractionalOwner      string `json:"fractional_owner"`
}

type RegisteredOwner struct {
	Name    string `json:"name"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	County  string `json:"county"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

type Airworthiness struct {
	TypeCertificateDataSheet string `json:"type_certificate_data_sheet"`
	TypeCertificateHolder    string `json:"type_certificate_holder"`
	EngineManufacturer       string `json:"engine_manufacturer"`
	Classification           string `json:"classification"`
	EngineModel              string `json:"engine_model"`
	Category                 string `json:"category"`
	Date                     string `json:"date"`
	ExceptionCode            string `json:"exception_code"`
}

var (
	ErrUnableToQueryFAA = fmt.Errorf("couldn't query the FAA registration page")
	ErrNotAssigned      = fmt.Errorf("n number not assigned")
)

func GetRegistration(n string) (*Registration, error) {

	baseURL := "https://registry.faa.gov/AircraftInquiry/Search/NNumberResult"

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}

	q.Add("nNumberTxt", n)
	u.RawQuery = q.Encode()

	regStatus := ""

	ac := Aircraft{}
	ro := RegisteredOwner{}
	aw := Airworthiness{}

	c := colly.NewCollector()

	c.OnHTML("#mainDiv", func(e *colly.HTMLElement) {
		status := e.ChildText(".noprint > p:nth-child(4)")
		if strings.Contains(status, "Not Assigned/Reserved") {
			regStatus = "not assigned/reserved"
			return
		}

		regStatus = "assigned"
	})

	// AIRCRAFT
	c.OnHTML("#mainDiv div:nth-child(5) table tbody td[data-label]", func(e *colly.HTMLElement) {

		switch e.Attr("data-label") {
		// Aircraft
		case "Serial Number":
			fmt.Printf("%s\n", e.Text)
			ac.SerialNumber = clean(e.Text)
		case "Status":
			ac.Status = clean(e.Text)
		case "Manufacturer Name":
			ac.ManufacturerName = clean(e.Text)
		case "Certificate Issue Date":
			ac.CertificateIssueDate = clean(e.Text)
		case "Model":
			ac.Model = clean(e.Text)
		case "Expiration Date":
			ac.ExpirationDate = clean(e.Text)
		case "Aircraft Type":
			ac.AircraftType = clean(e.Text)
		case "Engine Type":
			ac.EngineType = clean(e.Text)
		case "Pending Number Change":
			ac.PendingNumberChange = clean(e.Text)
		case "Dealer":
			ac.Dealer = clean(e.Text)
		case "Date Change Authorized":
			ac.DateChangeAuthorized = clean(e.Text)
		case "Mode S Code (Base 8 / oct)":
			ac.ModeSCodeOct = clean(e.Text)
		case "Mfr Year":
			ac.MfrYear = clean(e.Text)
		case "Mode S Code (Base 16 / Hex)":
			ac.ModeSCodeHex = clean(e.Text)
		case "Type Registration":
			ac.TypeRegistration = clean(e.Text)
		case "Fractional Owner":
			ac.FractionalOwner = clean(e.Text)
		}
	})

	// REGISTERED OWNER
	c.OnHTML("#mainDiv div:nth-child(6) table tbody td[data-label]", func(e *colly.HTMLElement) {

		switch e.Attr("data-label") {
		case "Name":
			ro.Name = clean(e.Text)
		case "Street":
			ro.Street = clean(e.Text)
		case "City":
			ro.City = clean(e.Text)
		case "State":
			ro.State = clean(e.Text)
		case "County":
			ro.County = clean(e.Text)
		case "Zip Code":
			ro.ZipCode = clean(e.Text)
		case "Country":
			ro.Country = clean(e.Text)
		}
	})

	// AIRWORTHINESS
	c.OnHTML("#mainDiv div:nth-child(7) table tbody td[data-label]", func(e *colly.HTMLElement) {

		switch e.Attr("data-label") {
		case "Type Certificate Data Sheet":
			aw.TypeCertificateDataSheet = clean(e.Text)
		case "Type Certificate Holder":
			aw.TypeCertificateHolder = clean(e.Text)
		case "Engine Manufacturer":
			aw.EngineManufacturer = clean(e.Text)
		case "Classification":
			aw.Classification = clean(e.Text)
		case "Engine Model":
			aw.EngineModel = clean(e.Text)
		case "Category":
			aw.Category = clean(e.Text)
		case "A/W Date":
			aw.Date = clean(e.Text)
		case "Exception Code":
			aw.ExceptionCode = clean(e.Text)
		}
	})

	err = c.Visit(u.String())
	if err != nil {
		return nil, ErrUnableToQueryFAA
	}

	if regStatus != "assigned" {
		return nil, ErrNotAssigned
	}

	reg := Registration{}
	reg.Aircraft = ac
	reg.RegisteredOwner = ro
	reg.Airworthiness = aw

	return &reg, nil
}

func clean(s string) string {
	return strings.TrimSpace(s)
}
