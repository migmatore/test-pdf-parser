package main

import (
	"fmt"
	"github.com/ledongthuc/pdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"log"
)

type BankType string

const (
	Sber    BankType = "Сбербанк"
	Tinkoff          = "Тинькофф"
	Alpha            = "Альфабанк"
)

type BankInfo struct {
	BankName      string
	DateTime      string
	RecipientName string
	Amount        string
	Status        string
}

func main() {
	pdf.DebugOn = true

	//content, err := parseBankPdf("sample.pdf", Sber)
	//content, err := parseBankPdf("receipt_16.07.2023.pdf", Tinkoff)
	content, err := parseBankPdf("document16.07.23.pdf", Alpha)

	//content, err := readPdf("sample.pdf") // Read local pdf file
	//content, err := readPdf("document16.07.23.pdf") // Read local pdf file
	//content, err := readPdf("receipt_16.07.2023.pdf") // Read local pdf file
	if err != nil {
		panic(err)
	}
	fmt.Println(content)

	if err := api.ExtractImagesFile("sample.pdf", "./", nil, nil); err != nil {
		log.Fatalf("%s", err.Error())
	}
}

func parseBankPdf(path string, bankType BankType) (BankInfo, error) {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return BankInfo{}, err
	}

	totalPage := r.NumPage()

	var bankInfo BankInfo

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		var flag bool

		var currentField string

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			for _, word := range row.Content {
				switch bankType {
				case Sber:
					bankInfo.BankName = string(Sber)

					if word.S == "Чек по операции" || word.S == "Сумма перевода" || word.S == "ФИО получателя" {
						currentField = word.S
						flag = true
						continue
					}

					if flag {
						flag = false

						switch currentField {
						case "Чек по операции":
							bankInfo.DateTime = word.S
						case "Сумма перевода":
							bankInfo.Amount = word.S
						case "ФИО получателя":
							bankInfo.RecipientName = word.S
						}
					}
				case Tinkoff:
					bankInfo.BankName = Tinkoff

					if word.S == "Сумма" || word.S == "Получатель" || word.S == "Статус" {
						currentField = word.S
						flag = true
						continue
					}

					if flag {
						flag = false

						switch currentField {
						case "Сумма":
							bankInfo.Amount = word.S
						case "Получатель":
							bankInfo.RecipientName = word.S
						case "Статус":
							bankInfo.Status = word.S
						}
					}
				case Alpha:
					bankInfo.BankName = Alpha

					if word.S == "Сумма:" || word.S == "Получатель:" || word.S == "Дата отправки перевода:" || word.S == "Статус:" {
						currentField = word.S
						flag = true
						continue
					}

					if flag {
						flag = false

						switch currentField {
						case "Дата отправки перевода:":
							bankInfo.DateTime = word.S
						case "Сумма:":
							bankInfo.Amount = word.S
						case "Получатель:":
							bankInfo.RecipientName = word.S
						case "Статус:":
							bankInfo.Status = word.S
						}
					}
				}
			}
		}
	}
	return bankInfo, nil
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return "", err
	}

	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			for _, word := range row.Content {
				fmt.Println(word.S)
			}
		}
	}
	return "", nil
}
