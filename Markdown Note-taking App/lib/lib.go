package lib

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/RohithBN/types"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/google/uuid"
)

func GenerateUniqueFileName(originalName string) string {
	uid := uuid.New().String()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	return uid + "_" + timestamp + "_" + originalName
}

var headingRe = regexp.MustCompile(`^#{1,6}\s*`)

func CleanMarkdownLine(line string) string {
	cleaned := headingRe.ReplaceAllString(line, "")
	return strings.TrimSpace(cleaned)
}

func FormatGrammarResponse(GrammarAPIResponse *types.GrammarResponse) []map[string]interface{} {
	output := []map[string]interface{}{}
	for _, m := range GrammarAPIResponse.Matches {
		wrongText := m.Context.Text[m.Context.Offset : m.Context.Offset+m.Context.Length]

		suggestions := []string{}
		for _, s := range m.Replacements {
			suggestions = append(suggestions, s.Value)
		}
		item := map[string]interface{}{
			"issue":       m.Message,
			"errorText":   wrongText,
			"sentence":    m.Sentence,
			"suggestions": suggestions,
		}
		output = append(output, item)
	}
	return output
}

func MarkdownToHTML(md string) []byte {
	mardown_in_bytes := []byte(md)
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(mardown_in_bytes)
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	return markdown.Render(doc,renderer)
}
