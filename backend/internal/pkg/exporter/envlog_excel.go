package exporter

import (
	"bytes"
	"fmt"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/xuri/excelize/v2"
)

// ExportEnvLogsToExcel generates an Excel sheet with climate logs.
func ExportEnvLogsToExcel(logs []domain.EnvironmentLog) ([]byte, error) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	sheetName := "Climate Logs"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	_ = f.DeleteSheet("Sheet1")

	// Title Block
	_ = f.SetCellValue(sheetName, "A1", "Environment Climate Monitoring Report")
	_ = f.MergeCell(sheetName, "A1", "G1")

	// Table Headers
	headers := []string{"Date & Time", "Zone ID", "Shift", "Temperature (°C)", "Humidity (%)", "Recorded By (User ID)", "Notes"}
	for colIdx, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 3)
		_ = f.SetCellValue(sheetName, cell, h)
	}

	// Stylings
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 16, Color: "003366"},
	})
	_ = f.SetCellStyle(sheetName, "A1", "G1", titleStyle)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"003366"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	_ = f.SetCellStyle(sheetName, "A3", "G3", headerStyle)

	// Populate Data
	for i, log := range logs {
		rowIdx := i + 4
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx), log.RecordedAt.Format("2006-01-02 15:04:05"))
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx), log.ZoneID)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIdx), log.Shift)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIdx), log.Temperature)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIdx), log.Humidity)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIdx), log.RecordedBy)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowIdx), log.Notes)
	}

	// Set widths for a readable clean grid layout
	for colIdx := 1; colIdx <= 7; colIdx++ {
		colName, _ := excelize.ColumnNumberToName(colIdx)
		_ = f.SetColWidth(sheetName, colName, colName, 22)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
