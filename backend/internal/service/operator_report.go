package service

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"example.com/german/backend/internal/repos"
)

var reportHeaders = []string{
	"created_at",
	"closed_at",
	"category",
	"applicant_type",
	"status",
	"priority",
	"return_count",
	"seconds_to_accept",
	"seconds_to_first_response",
	"seconds_to_close",
}

func buildCSVReport(records []repos.ReportTicketRecord) ([]byte, error) {
	var output bytes.Buffer
	output.Write([]byte{0xef, 0xbb, 0xbf})
	writer := csv.NewWriter(&output)
	if err := writer.Write(reportHeaders); err != nil {
		return nil, fmt.Errorf("write CSV header: %w", err)
	}
	for _, record := range records {
		if err := writer.Write(reportRow(record)); err != nil {
			return nil, fmt.Errorf("write CSV row: %w", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("finish CSV report: %w", err)
	}
	return output.Bytes(), nil
}

func buildXLSXReport(records []repos.ReportTicketRecord) ([]byte, error) {
	var output bytes.Buffer
	archive := zip.NewWriter(&output)
	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
  <Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>
  <Override PartName="/xl/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/>
</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
</Relationships>`,
		"xl/workbook.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <sheets><sheet name="Tickets" sheetId="1" r:id="rId1"/></sheets>
</workbook>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>
</Relationships>`,
		"xl/styles.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
  <fonts count="1"><font><sz val="11"/><name val="Calibri"/></font></fonts>
  <fills count="1"><fill><patternFill patternType="none"/></fill></fills>
  <borders count="1"><border/></borders>
  <cellStyleXfs count="1"><xf/></cellStyleXfs>
  <cellXfs count="1"><xf xfId="0"/></cellXfs>
</styleSheet>`,
	}
	for name, content := range files {
		entry, err := archive.Create(name)
		if err != nil {
			return nil, fmt.Errorf("create XLSX entry %s: %w", name, err)
		}
		if _, err := entry.Write([]byte(content)); err != nil {
			return nil, fmt.Errorf("write XLSX entry %s: %w", name, err)
		}
	}
	sheet, err := archive.Create("xl/worksheets/sheet1.xml")
	if err != nil {
		return nil, fmt.Errorf("create XLSX worksheet: %w", err)
	}
	if _, err := sheet.Write([]byte(xlsxSheet(records))); err != nil {
		return nil, fmt.Errorf("write XLSX worksheet: %w", err)
	}
	if err := archive.Close(); err != nil {
		return nil, fmt.Errorf("finish XLSX report: %w", err)
	}
	return output.Bytes(), nil
}

func xlsxSheet(records []repos.ReportTicketRecord) string {
	var output strings.Builder
	output.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	output.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)
	writeXLSXRow(&output, 1, reportHeaders)
	for index, record := range records {
		writeXLSXRow(&output, index+2, reportRow(record))
	}
	output.WriteString(`</sheetData></worksheet>`)
	return output.String()
}

func writeXLSXRow(output *strings.Builder, rowNumber int, values []string) {
	output.WriteString(`<row r="` + strconv.Itoa(rowNumber) + `">`)
	for column, value := range values {
		cell := excelColumn(column+1) + strconv.Itoa(rowNumber)
		output.WriteString(`<c r="` + cell + `" t="inlineStr"><is><t>`)
		var escaped bytes.Buffer
		_ = xml.EscapeText(&escaped, []byte(value))
		output.Write(escaped.Bytes())
		output.WriteString(`</t></is></c>`)
	}
	output.WriteString(`</row>`)
}

func excelColumn(number int) string {
	var result string
	for number > 0 {
		number--
		result = string(rune('A'+number%26)) + result
		number /= 26
	}
	return result
}

func reportRow(record repos.ReportTicketRecord) []string {
	return []string{
		record.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		nullableTimeString(record.ClosedAt),
		record.Category,
		record.ApplicantType.String(),
		record.Status.String(),
		record.Priority.String(),
		strconv.Itoa(record.ReturnCount),
		nullableNumberString(record.SecondsToAccept),
		nullableNumberString(record.SecondsToFirstResponse),
		nullableNumberString(record.SecondsToClose),
	}
}

func nullableTimeString(value sql.NullTime) string {
	if !value.Valid {
		return ""
	}
	return value.Time.UTC().Format(time.RFC3339)
}

func nullableNumberString(value sql.NullFloat64) string {
	if !value.Valid {
		return ""
	}
	return strconv.FormatFloat(value.Float64, 'f', 3, 64)
}
