package pdiutil

import (
	"strconv"

	"github.com/unidoc/unioffice/color"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/spreadsheet"
)

// OverviewStatus enum
type OverviewStatus string

const (
	// Successful status
	Successful OverviewStatus = "Successful"
	// Warning status
	Warning = "Warning"
	// FatalError status
	FatalError = "Fatal Error"
)

// OverviewItem type
type OverviewItem struct {
	ItemName        string
	ItemDescription string
	ItemStatus      OverviewStatus
}

func addOverviewSheetTo(workbook *spreadsheet.Workbook, items []OverviewItem, solution Solution) spreadsheet.Sheet {

	sheet := workbook.AddSheet()
	sheet.SetName("Overview")

	headerFontStyle := workbook.StyleSheet.AddFont()
	headerFontStyle.SetBold(true)
	headerFontStyle.SetSize(16)
	headerStyle := workbook.StyleSheet.AddCellStyle()
	headerStyle.SetFont(headerFontStyle)

	greyFill := workbook.StyleSheet.Fills().AddFill()
	greyFill.SetPatternFill().SetFgColor(color.LightGray)
	defaultStyle := workbook.StyleSheet.AddCellStyle()
	defaultStyle.SetFill(greyFill)

	greenFill := workbook.StyleSheet.Fills().AddFill()
	greenFill.SetPatternFill().SetFgColor(color.LightGreen)
	successStyle := workbook.StyleSheet.AddCellStyle()
	successStyle.SetFill(greenFill)

	yelloFill := workbook.StyleSheet.Fills().AddFill()
	yelloFill.SetPatternFill().SetFgColor(color.Yellow)
	warningStyle := workbook.StyleSheet.AddCellStyle()
	warningStyle.SetFill(yelloFill)

	redFill := workbook.StyleSheet.Fills().AddFill()
	redFill.SetPatternFill().SetFgColor(color.Red)
	errorStyle := workbook.StyleSheet.AddCellStyle()
	errorStyle.SetFill(redFill)

	headerRow := sheet.AddRow()
	headerRow.SetHeightAuto()

	headerCell := headerRow.AddCell()
	headerCell.SetString("Solution Overview")
	headerCell.SetStyle(headerStyle)

	solutionInfo := [][]string{
		[]string{"Name", solution.Name},
		[]string{"Description", solution.Description},
		[]string{"IsPatch", strconv.FormatBool(solution.PatchSolution)},
		[]string{"Status", solution.Status},
		[]string{"Contact", solution.Contact},
		[]string{"Email", solution.Email},
	}

	for _, solutionLine := range solutionInfo {
		_row := sheet.AddRow()
		_row.AddCell().SetString(solutionLine[0])
		_row.AddCell().SetString(solutionLine[1])
	}

	sheet.AddRow() // empty row

	statusRow := sheet.AddRow()
	statusCell := statusRow.AddCell()
	statusCell.SetString("Status Overview")
	statusCell.SetStyle(headerStyle)

	for _, item := range items {
		row := sheet.AddRow()
		row.SetHeightAuto()
		row.AddCell().SetString(item.ItemName)
		row.AddCell().SetString(item.ItemDescription)

		statusCell := row.AddCell()

		statusCell.SetString(string(item.ItemStatus))

		switch item.ItemStatus {
		case Successful:
			statusCell.SetStyle(successStyle)
		case Warning:
			statusCell.SetStyle(warningStyle)
		case FatalError:
			statusCell.SetStyle(errorStyle)
		default:
			statusCell.SetStyle(defaultStyle)
		}

	}

	sheet.Column(1).SetWidth(measurement.Pixel72 * 280)
	sheet.Column(2).SetWidth(measurement.Pixel72 * 300)
	sheet.Column(3).SetWidth(measurement.Pixel72 * 120)

	return sheet
}

func addSheetTo(workbook *spreadsheet.Workbook, sheetName string, headers []string, data [][]string) spreadsheet.Sheet {
	columnMaxWidth := make([]int, len(headers))

	sheet := workbook.AddSheet()

	headerFontStyle := workbook.StyleSheet.AddFont()
	headerFontStyle.SetBold(true)
	headerFontStyle.SetSize(13)

	cs := workbook.StyleSheet.AddCellStyle()
	cs.SetFont(headerFontStyle)

	sheet.SetName(sheetName)

	header := sheet.AddRow()

	header.SetHeightAuto()

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

// AddOverviewSheetTo internal usage
var AddOverviewSheetTo = addOverviewSheetTo

// AddSheetTo internal usage
var AddSheetTo = addSheetTo
