package taxcom

import (
	"log"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type taxcom struct {
	AgreementNumber string
	Login           string
	Password        string
	IdIntegrator    string
	r               *resty.Client
}

type Session struct {
	AgreementNumber string
	SessionToken    string `json:"sessionToken"`
}

type KKT struct {
	Address  string `json:"address"`
	Kktregid string `json:"kktregid"`
}

type Receipt struct {
	ID             int
	FP             string
	DocumentNumber int
	KktRegId       string
	FD             string
	Date           string
	Products       []Product
	Link           string
	Price          int
	VatPrice       int
}

type Product struct {
	Name       string
	Quantity   int
	Price      int
	Vat        int
	VatPrice   int
	TotalPrice int
	FP         string
	FD         string
	FN         string
	Time       string
}

type TPaginator struct {
	RecordCount           int `json:"recordCount"`
	RecordFilteredCount   int `json:"recordFilteredCount"`
	RecordInResponceCount int `json:"recordInResponceCount"`
}
type TResultOutlet struct {
	ReportDate string     `json:"reportDate"`
	Counts     TPaginator `json:"counts"`
	Records    []TOutlet  `json:"records"`
}
type TResultAccountList struct {
	ReportDate string         `json:"reportDate"`
	Counts     TPaginator     `json:"counts"`
	Records    []TAccountList `json:"records"`
}
type TResultKkt struct {
	ReportDate string     `json:"reportDate"`
	Counts     TPaginator `json:"counts"`
	Records    []TKktList `json:"records"`
}

type TResultShift struct {
	ReportDate string     `json:"reportDate"`
	Counts     TPaginator `json:"counts"`
	Records    []TShift   `json:"records"`
}

type TResultDocumentList struct {
	ReportDate string          `json:"reportDate"`
	Counts     TPaginator      `json:"counts"`
	Records    []TDocumentList `json:"records"`
}

type TAccountList struct {
	AgreementNumber string `json:"agreementNumber"`
	CompanyName     string `json:"companyName"`
	Inn             string `json:"inn"`
}

type TShift struct {
	FnFactoryNumber string `json:"fnFactoryNumber"`
	ShiftNumber     int    `json:"shiftNumber"`
	OpenDateTime    string `json:"openDateTime"`
	CloseDateTime   string `json:"closeDateTime"`
	ReceiptCount    int    `json:"receiptCount"`
}

type TKktList struct {
	Name             string `json:"name"`
	KktRegNumber     string `json:"kktRegNumber"`
	KktFactoryNumber string `json:"kktFactoryNumber"`
	FnFactoryNumber  string `json:"fnFactoryNumber"`
	CashdeskState    string `json:"cashdeskState"`
	ProblemIndicator string `json:"problemIndicator"`
}

type TOutlet struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Code             string `json:"code"`
	Address          string `json:"address"`
	ProblemIndicator string `json:"problemIndicator"`
	Department       string `json:"department"`
}

type TDocumentURL struct {
	TaxcomReceiptUrl string `json:"taxcomReceiptUrl"`
}

type TDocumentList struct {
	FnFactoryNumber string `json:"fnFactoryNumber"`
	FdNumber        int    `json:"fdNumber"`
}

type TDocument struct {
	Document struct {
		TransactionDate      string     `json:"1012"`
		FiscalDriveNumber    string     `json:"1041"`
		DocumentNumber       int        `json:"1042"`
		EcashTotalSum        int        `json:"1081"`
		FiscalDocumentNumber int        `json:"1040"`
		FP                   string     `json:"1077"`
		TaxationType         int        `json:"1055"`
		NdsNo                int        `json:"1105"`
		Nds0                 int        `json:"1104"`
		Nds10                int        `json:"1103"`
		Nds18                int        `json:"1102"`
		Nds20                int        `json:"nds20"`
		UserInn              string     `json:"1018"`
		KktRegId             string     `json:"1037"`
		CashTotalSum         int        `json:"1031"`
		TotalSum             int        `json:"1020"`
		OperationType        int        `json:"1054"`
		Items                []TProduct `json:"1059"`
	} `json:"document"`
}

type TProduct struct {
	Name     string `json:"1030"`
	Sum      int    `json:"1043"`
	Price    int    `json:"1079"`
	Quantity string `json:"1023"`
}

type TaxCom struct {
	Type         string
	BaseURL      string
	Login        string
	Password     string
	IdIntegrator string
}

var session Session

func Taxcom(login, password, idIntegrator, agreementNumber string) *taxcom {
	return &taxcom{
		r:               resty.New(),
		Login:           login,
		Password:        password,
		IdIntegrator:    idIntegrator,
		AgreementNumber: agreementNumber,
	}
}

func (ofd *taxcom) auth() {
	body := make(map[string]interface{})
	body["login"] = ofd.Login
	body["password"] = ofd.Password
	if ofd.AgreementNumber != "" {
		body["agreementNumber"] = ofd.AgreementNumber
		session.AgreementNumber = ofd.AgreementNumber
	}
	_, err := ofd.r.R().
		SetBody(body).
		SetHeader("Integrator-ID", ofd.IdIntegrator).
		SetResult(&session).
		Post("https://api-lk-ofd.taxcom.ru/API/v2/Login")

	if err != nil {
		log.Printf("[TaxCom] failed auth: %v", err)
	}
}

func (ofd *taxcom) startDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func (ofd *taxcom) endDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}

func (ofd *taxcom) GetAccountList() []string {
	if session.SessionToken == "" {
		ofd.auth()
	}
	var agreementNumbers []string
	ot := TResultAccountList{}
	_, err := ofd.r.R().
		SetHeader("Session-Token", session.SessionToken).
		SetResult(&ot).
		Get("https://api-lk-ofd.taxcom.ru/API/v2/AccountList")
	for _, v := range ot.Records {
		agreementNumbers = append(agreementNumbers, v.AgreementNumber)
	}
	if err != nil {
		log.Printf("[TaxCom] AccountList: %v", err)
	}
	return agreementNumbers
}

func (ofd *taxcom) getOutletList() (outlet []string) {
	ot := TResultOutlet{}
	_, err := ofd.r.R().
		SetHeader("Session-Token", session.SessionToken).
		SetResult(&ot).
		Get("https://api-lk-ofd.taxcom.ru/API/v2/OutletList")
	for _, v := range ot.Records {
		outlet = append(outlet, v.Id)
	}
	if err != nil {
		log.Printf("[TaxCom] Login: %v", err)
	}
	return outlet
}

func (ofd *taxcom) getShiftList(fn string, begin time.Time, end time.Time) (shift TResultShift) {
	_, err := ofd.r.R().
		SetHeader("Session-Token", session.SessionToken).
		SetResult(&shift).
		Get("https://api-lk-ofd.taxcom.ru/API/v2/ShiftList?fn=" + fn + "&begin=" + begin.Format("2006-01-02T15:04:05") + "&end=" + end.Format("2006-01-02T15:04:05"))
	if err != nil {
		log.Printf("[TaxCom] ShiftList: %s", err.Error())
	}

	return shift
}

func (ofd *taxcom) getDocumentList(fn string, shift int) TResultDocumentList {
	documentList := TResultDocumentList{}
	_, err := ofd.r.R().
		SetHeader("Session-Token", session.SessionToken).
		SetResult(&documentList).
		Get("https://api-lk-ofd.taxcom.ru/API/v2/DocumentList?fn=" + fn + "&shift=" + strconv.Itoa(shift))

	if err != nil {
		log.Printf("[TaxCom] DocumentList: %s", err.Error())
	}

	return documentList
}

func (ofd *taxcom) getDocumentLink(fn string, fd int) string {
	documentUrl := TDocumentURL{}
	_, err := ofd.r.R().
		SetHeader("Session-Token", session.SessionToken).
		SetResult(&documentUrl).
		Get("https://api-lk-ofd.taxcom.ru/API/v2/DocumentURL?fn=" + fn + "&fd=" + strconv.Itoa(fd))

	if err != nil {
		log.Printf("[TaxCom] DocumentList: %s", err.Error())
	}
	return documentUrl.TaxcomReceiptUrl
}

func (ofd *taxcom) getKKT() (kkt []KKT, err error) {
	outlets := ofd.getOutletList()
	for _, v := range outlets {
		k := TResultKkt{}
		_, err = ofd.r.R().
			SetHeader("Session-Token", session.SessionToken).
			SetResult(&k).
			Get("https://api-lk-ofd.taxcom.ru/API/v2/KKTList?id=" + v)

		for _, v := range k.Records {
			kkt = append(kkt, KKT{
				Address:  v.CashdeskState,
				Kktregid: v.FnFactoryNumber,
			})
		}
		if err != nil {
			log.Printf("[TaxCom] KKTList: %s", err.Error())
		}
	}

	return kkt, err
}

func (ofd *taxcom) GetReceipts(date time.Time) (receipts []Receipt, err error) {
	if session.SessionToken == "" || session.AgreementNumber != ofd.AgreementNumber {
		ofd.auth()
	}
	kkts, err := ofd.getKKT()
	if err != nil {
		return receipts, err
	}
	for _, v := range kkts {
		r, err := ofd.getDocuments(v.Kktregid, date)
		if err != nil {
			return receipts, err
		}
		receipts = append(receipts, r...)
	}
	return receipts, err
}

func (ofd *taxcom) getDocuments(kkt string, date time.Time) (documents []Receipt, err error) {
	shift := ofd.getShiftList(kkt, ofd.startDay(date), ofd.endDay(date))
	for _, v := range shift.Records {
		documentList := ofd.getDocumentList(kkt, v.ShiftNumber)
		for _, dn := range documentList.Records {
			docs := TDocument{}
			_, err := ofd.r.R().
				SetHeader("Session-Token", session.SessionToken).
				SetResult(&docs).
				Get("https://api-lk-ofd.taxcom.ru/API/v2/DocumentInfo?fn=" + kkt + "&fd=" + strconv.Itoa(dn.FdNumber))

			if err != nil {
				log.Printf("[TaxCom] DocumentInfo: %s", err.Error())
			}
			if docs.Document.FP == "" {
				continue
			}
			transactionDate, _ := time.Parse("2006-01-02T15:04:05", docs.Document.TransactionDate)
			doc := Receipt{}
			doc.Date = transactionDate.Format(time.RFC3339)
			doc.FD = strconv.Itoa(docs.Document.FiscalDocumentNumber)
			doc.DocumentNumber = docs.Document.DocumentNumber
			doc.FP = docs.Document.FP
			doc.Price = docs.Document.TotalSum
			doc.KktRegId = docs.Document.KktRegId
			var products []Product
			for _, v := range docs.Document.Items {
				q, _ := strconv.Atoi(v.Quantity)
				product := Product{}
				product.Name = v.Name
				product.Quantity = q
				product.Price = v.Price
				product.TotalPrice = v.Sum
				products = append(products, product)
			}
			doc.Products = products
			doc.Link = ofd.getDocumentLink(docs.Document.FiscalDriveNumber, docs.Document.FiscalDocumentNumber)

			documents = append(documents, doc)
		}
	}

	return documents, err
}
