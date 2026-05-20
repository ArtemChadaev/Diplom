package pdf

import (
	"bytes"
	"fmt"
	"strconv"

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

type TTNDocument struct {
	OrderNumber  string
	CustomerName string
	OrderDate    string
	OrderStatus  string
	OrderType    string
	Priority     string
	Items        []TTNItem
}

type TTNItem struct {
	ProductName string
	SKU         string
	BatchSerial string
	Quantity    int
}

func GenerateTTN(doc TTNDocument) ([]byte, error) {
	cfg := config.NewBuilder().WithPageNumber().Build()
	m := maroto.New(cfg)

	blue := &props.Color{Red: 0, Green: 51, Blue: 102}
	gray := &props.Color{Red: 100, Green: 100, Blue: 100}
	lightGray := &props.Color{Red: 240, Green: 240, Blue: 240}

	m.AddRows(
		row.New(12).Add(
			col.New(12).Add(
				text.New("DIPLOM ERP — PHARMACEUTICAL WAREHOUSE SYSTEM", props.Text{
					Size:  8,
					Style: fontstyle.Bold,
					Align: align.Left,
					Color: gray,
				}),
			),
		),
		line.NewRow(1, props.Line{Color: blue}),
		row.New(8).Add(col.New(12).Add(text.New("", props.Text{}))),
	)

	m.AddRows(
		row.New(20).Add(
			col.New(12).Add(
				text.New(fmt.Sprintf("CONSIGNMENT NOTE (ТТН) № %s", doc.OrderNumber), props.Text{
					Size:  16,
					Style: fontstyle.Bold,
					Align: align.Center,
					Color: blue,
				}),
			),
		),
		row.New(8).Add(col.New(12).Add(text.New("", props.Text{}))),
	)

	m.AddRows(
		row.New(12).Add(
			col.New(6).Add(
				text.New(fmt.Sprintf("Customer: %s", doc.CustomerName), props.Text{
					Size:  10,
					Style: fontstyle.Bold,
					Align: align.Left,
				}),
			),
			col.New(6).Add(
				text.New(fmt.Sprintf("Order Type: %s", doc.OrderType), props.Text{
					Size:  10,
					Style: fontstyle.Normal,
					Align: align.Right,
				}),
			),
		),
		row.New(12).Add(
			col.New(6).Add(
				text.New(fmt.Sprintf("Date: %s", doc.OrderDate), props.Text{
					Size:  10,
					Style: fontstyle.Normal,
					Align: align.Left,
				}),
			),
			col.New(6).Add(
				text.New(fmt.Sprintf("Priority: %s", doc.Priority), props.Text{
					Size:  10,
					Style: fontstyle.Normal,
					Align: align.Right,
				}),
			),
		),
		row.New(12).Add(
			col.New(6).Add(
				text.New(fmt.Sprintf("Status: %s", doc.OrderStatus), props.Text{
					Size:  10,
					Style: fontstyle.Normal,
					Align: align.Left,
				}),
			),
			col.New(6).Add(text.New("", props.Text{})),
		),
		row.New(15).Add(col.New(12).Add(text.New("", props.Text{}))),
	)

	m.AddRows(
		row.New(15).Add(
			col.New(12).Add(
				text.New("SHIPPED PRODUCTS LIST", props.Text{
					Size:  12,
					Style: fontstyle.Bold,
					Align: align.Left,
					Color: blue,
				}),
			),
		),
		row.New(5).Add(col.New(12).Add(text.New("", props.Text{}))),
	)

	m.AddRows(
		row.New(15).Add(
			col.New(5).Add(
				text.New("Product Name", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
					Align: align.Left,
				}),
			),
			col.New(3).Add(
				text.New("SKU", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
					Align: align.Center,
				}),
			),
			col.New(2).Add(
				text.New("Batch / Serial", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
					Align: align.Center,
				}),
			),
			col.New(2).Add(
				text.New("Quantity", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
					Align: align.Right,
				}),
			),
		),
		line.NewRow(1, props.Line{Color: blue}),
	)

	for _, item := range doc.Items {
		m.AddRows(
			row.New(14).Add(
				col.New(5).Add(
					text.New(item.ProductName, props.Text{
						Size:  9,
						Style: fontstyle.Normal,
						Align: align.Left,
					}),
				),
				col.New(3).Add(
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
					text.New(strconv.Itoa(item.Quantity), props.Text{
						Size:  9,
						Style: fontstyle.Normal,
						Align: align.Right,
					}),
				),
			),
			line.NewRow(0.5, props.Line{Color: lightGray}),
		)
	}

	m.AddRows(
		row.New(30).Add(col.New(12).Add(text.New("", props.Text{}))),
		row.New(15).Add(
			col.New(6).Add(
				text.New("Released by: _________________________", props.Text{
					Size:  10,
					Style: fontstyle.Normal,
					Align: align.Left,
				}),
			),
			col.New(6).Add(
				text.New("Received by: _________________________", props.Text{
					Size:  10,
					Style: fontstyle.Normal,
					Align: align.Right,
				}),
			),
		),
		row.New(10).Add(
			col.New(6).Add(
				text.New("Signature / Timestamp", props.Text{
					Size:  8,
					Style: fontstyle.Normal,
					Align: align.Left,
					Color: gray,
				}),
			),
			col.New(6).Add(
				text.New("Signature / Timestamp", props.Text{
					Size:  8,
					Style: fontstyle.Normal,
					Align: align.Right,
					Color: gray,
				}),
			),
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
