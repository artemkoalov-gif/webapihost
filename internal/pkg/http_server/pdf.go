package http_server

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

func buildPDF(lines []string) ([]byte, error) {
	if len(lines) == 0 {
		lines = []string{"No data available"}
	}

	var buf bytes.Buffer
	offsets := []int{0} // slot for the free object

	writeObject := func(id int, body string) {
		for len(offsets) <= id {
			offsets = append(offsets, 0)
		}
		offsets[id] = buf.Len()
		buf.WriteString(body)
	}

	buf.WriteString("%PDF-1.4\n")

	writeObject(1, "1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	writeObject(2, "2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")
	writeObject(3, "3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>\nendobj\n")

	content := buildContentStream(lines)
	writeObject(4, fmt.Sprintf("4 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", len(content), content))

	writeObject(5, "5 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n")

	xrefStart := buf.Len()
	fmt.Fprintf(&buf, "xref\n0 %d\n0000000000 65535 f \n", len(offsets))
	for i := 1; i < len(offsets); i++ {
		fmt.Fprintf(&buf, "%010d 00000 n \n", offsets[i])
	}
	fmt.Fprintf(&buf, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF", len(offsets), xrefStart)

	return buf.Bytes(), nil
}

func buildLinesForPDF(data map[int]ResponseLinks, ids []int) []string {
	lines := make([]string, 0, len(data)*2)

	for _, id := range ids {
		links, ok := data[id]
		if !ok {
			continue
		}

		lines = append(lines, fmt.Sprintf("List #%d (links: %d)", id, links.LinksNum))

		urls := make([]string, 0, len(links.Links))
		for url := range links.Links {
			urls = append(urls, url)
		}
		sort.Strings(urls)

		for _, url := range urls {
			lines = append(lines, fmt.Sprintf("%s - %s", url, links.Links[url]))
		}
		lines = append(lines, "")
	}

	return lines
}

func buildContentStream(lines []string) string {
	var b strings.Builder
	b.WriteString("BT\n/F1 12 Tf\n72 720 Td\n")

	for i, line := range lines {
		if i > 0 {
			b.WriteString("0 -18 Td\n")
		}
		b.WriteString("(")
		b.WriteString(escapePDFString(line))
		b.WriteString(") Tj\n")
	}

	b.WriteString("ET")
	return b.String()
}

func escapePDFString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	return s
}
