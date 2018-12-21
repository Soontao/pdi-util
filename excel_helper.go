package main

import (
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/spreadsheet"
)

func addSheetTo(workbook *spreadsheet.Workbook, sheetName string, headers []string, data [][]string) spreadsheet.Sheet {
	columnMaxWidth := make([]int, len(headers))

	sheet := workbook.AddSheet()

	fs := workbook.StyleSheet.AddFont()
	fs.SetBold(true)
	cs := workbook.StyleSheet.AddCellStyle()
	cs.SetFont(fs)

	sheet.SetName(sheetName)

	header := sheet.AddRow()

	for columnIdx, cellData := range headers {
		c := header.AddCell()
		c.SetString(cellData)
		c.SetStyle(cs)
		columnMaxWidth[columnIdx] = len(cellData)
	}
	for _, rowData := range data {
		row := sheet.AddRow()
		row.SetHeightAuto()
		for columnIdx, cellData := range rowData {

			row.AddCell().SetString(cellData)

			cellLen := len(cellData)

			if columnMaxWidth[columnIdx] < cellLen {
				columnMaxWidth[columnIdx] = cellLen
			}
		}
	}
	for columnIdx, stringLen := range columnMaxWidth {
		sheet.Column(uint32(columnIdx) + 1).SetWidth((measurement.Distance(measurement.Pixel72 * 8 * (stringLen + 2))))
	}

	return sheet
}
