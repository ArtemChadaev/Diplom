package pdf

import (
	"bytes"
	"fmt"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type QualityDocument struct {
	OrderNumber  string
	CustomerName string
	OrderDate    string
	Items        []QualityItem
}

type QualityItem struct {
	ProductName string
	SKU         string
	BatchSerial string
	ExpiryDate  string
	Status      string
}

func GenerateQualityRegistry(doc QualityDocument) ([]byte, error) {
	cfg := config.NewBuilder().WithPageNumber().Build()
	m := maroto.New(cfg)

	green := &props.Color{Red: 46, Green: 117, Blue: 89}
	gray := &props.Color{Red: 100, Green: 100, Blue: 100}
	lightGray := &props.Color{Red: 240, Green: 240, Blue: 240}

	m.AddRows(
		row.New(12).Add(
			col.New(12).Add(
				text.New("DIPLOM ERP — PHARMACEUTICAL QUALITY CONTROL", props.Text{
					Size:  8,
					Style: fontstyle.Bold,
					Align: align.Left,
					Color: gray,
				}),
			),
		),
		line.NewRow(1, props.Line{Color: green}),
		row.New(8).Add(col.New(12).Add(text.New("", props.Text{}))),
	)

	m.AddRows(
		row.New(20).Add(
			col.New(12).Add(
				text.New(fmt.Sprintf("QUALITY CERTIFICATE REGISTRY № %s", doc.OrderNumber), props.Text{
					Size:  16,
					Style: fontstyle.Bold,
					Align: align.Center,
					Color: green,
				}),
			),
		),
		row.New(8).Add(col.New(12).Add(text.New("", props.Text{}))),
	)

	m.AddRows(
		row.New(12).Add(
			col.New(6).Add(
				text.New(fmt.Sprintf("Recipient: %s", doc.CustomerName), props.Text{
					Size:  10,
					Style: fontstyle.Bold,
					Align: align.Left,
				}),
			),
			col.New(6).Add(
				text.New(fmt.Sprintf("Registry Date: %s", doc.OrderDate), props.Text{
					Size:  10,
					Style: fontstyle.Normal,
					Align: align.Right,
				}),
			),
		),
		row.New(15).Add(col.New(12).Add(text.New("", props.Text{}))),
	)

	m.AddRows(
		row.New(15).Add(
			col.New(12).Add(
				text.New("VERIFIED QUALITY ASSURANCE BATCHES", props.Text{
					Size:  12,
					Style: fontstyle.Bold,
					Align: align.Left,
					Color: green,
				}),
			),
		),
		row.New(5).Add(col.New(12).Add(text.New("", props.Text{}))),
	)

	m.AddRows(
		row.New(15).Add(
			col.New(4).Add(
				text.New("Product Name", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
					Align: align.Left,
				}),
			),
			col.New(2).Add(
				text.New("SKU", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
					Align: align.Center,
				}),
			),
			col.New(2).Add(
				text.New("Batch Serial", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
					Align: align.Center,
				}),
			),
			col.New(2).Add(
				text.New("Expiry Date", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
					Align: align.Center,
				}),
			),
			col.New(2).Add(
				text.New("Control Status", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
					Align: align.Right,
				}),
			),
		),
		line.NewRow(1, props.Line{Color: green}),
	)

	for _, item := range doc.Items {
		m.AddRows(
			row.New(14).Add(
				col.New(4).Add(
					text.New(item.ProductName, props.Text{
						Size:  9,
						Style: fontstyle.Normal,
						Align: align.Left,
					}),
				),
				col.New(2).Add(
					text.New(item.SKU, props.Text{
						Size:  9,
						Style: fontstyle.Normal,
						Align: align.Center,
					}),
				),
				col.New(2).Add(
					text.New(item.BatchSerial, props.Text{
						Size:  9,
						Style: fontstyle.Normal,
						Align: align.Center,
					}),
				),
				col.New(2).Add(
					text.New(item.ExpiryDate, props.Text{
						Size:  9,
						Style: fontstyle.Normal,
						Align: align.Center,
					}),
				),
				col.New(2).Add(
					text.New(item.Status, props.Text{
						Size:  9,
						Style: fontstyle.Bold,
						Align: align.Right,
						Color: green,
					}),
				),
			),
			line.NewRow(0.5, props.Line{Color: lightGray}),
		)
	}

	m.AddRows(
		row.New(40).Add(col.New(12).Add(text.New("", props.Text{}))),
		row.New(15).Add(
			col.New(8).Add(
				text.New("Quality Assurance Officer: _________________________", props.Text{
					Size:  10,
					Style: fontstyle.Normal,
					Align: align.Left,
				}),
			),
			col.New(4).Add(
				text.New("L.S. (М.П.)", props.Text{
					Size:  10,
					Style: fontstyle.Normal,
					Align: align.Center,
					Color: gray,
				}),
			),
		),
		row.New(10).Add(
			col.New(8).Add(
				text.New("Signature / Full Name", props.Text{
					Size:  8,
					Style: fontstyle.Normal,
					Align: align.Left,
					Color: gray,
				}),
			),
			col.New(4).Add(text.New("", props.Text{})),
		),
	)

	document, err := m.Generate()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = buf.Write(document.GetBytes())
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
