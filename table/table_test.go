package table

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/jedib0t/go-pretty/text"
	"github.com/stretchr/testify/assert"
)

var (
	testAlign           = []text.Align{text.AlignDefault, text.AlignLeft, text.AlignLeft, text.AlignRight}
	testCaption         = "A Song of Ice and Fire"
	testColor           = text.Colors{text.FgGreen}
	testColorBoW        = text.Colors{text.FgBlack, text.BgWhite}
	testColorHiRedBold  = text.Colors{text.FgHiRed, text.Bold}
	testColorHiBlueBold = text.Colors{text.FgHiBlue, text.Bold}
	testColorWoB        = text.Colors{text.FgWhite, text.BgBlack}
	testColors          = []text.Colors{testColor, testColor, testColor, testColor, {text.FgCyan}}
	testColorsFooter    = []text.Colors{{}, {}, testColorHiBlueBold, testColorHiBlueBold}
	testColorsHeader    = []text.Colors{testColorHiRedBold, testColorHiRedBold, testColorHiRedBold, testColorHiRedBold}
	testCSSClass        = "test-css-class"
	testFooter          = Row{"", "", "Total", 10000}
	testFooterMultiLine = Row{"", "", "Total\nSalary", 10000}
	testHeader          = Row{"#", "First Name", "Last Name", "Salary"}
	testHeaderMultiLine = Row{"#", "First\nName", "Last\nName", "Salary"}
	testRows            = []Row{
		{1, "Arya", "Stark", 3000},
		{20, "Jon", "Snow", 2000, "You know nothing, Jon Snow!"},
		{300, "Tyrion", "Lannister", 5000},
	}
	testRowMultiLine = Row{0, "Winter", "Is", 0, "Coming.\r\nThe North Remembers!\nThis is known."}
	testRowNewLines  = Row{0, "Valar", "Morghulis", 0, "Faceless\nMen"}
	testRowPipes     = Row{0, "Valar", "Morghulis", 0, "Faceless|Men"}
	testRowTabs      = Row{0, "Valar", "Morghulis", 0, "Faceless\tMen"}
	testTitle1       = "Game of Thrones"
	testTitle2       = "When you play the Game of Thrones, you win or you die. There is no middle ground."
)

func init() {
	text.EnableColors()
}

type myMockOutputMirror struct {
	mirroredOutput string
}

func (t *myMockOutputMirror) Write(p []byte) (n int, err error) {
	t.mirroredOutput += string(p)
	return len(p), nil
}

func TestNewWriter(t *testing.T) {
	tw := NewWriter()
	assert.NotNil(t, tw.Style())
	assert.Equal(t, StyleDefault, *tw.Style())

	tw.SetStyle(StyleBold)
	assert.NotNil(t, tw.Style())
	assert.Equal(t, StyleBold, *tw.Style())
}

func TestTable_AppendFooter(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, len(table.rowsFooterRaw))

	table.AppendFooter([]interface{}{})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 1, len(table.rowsFooterRaw))
	assert.Equal(t, 0, len(table.rowsHeaderRaw))

	table.AppendFooter([]interface{}{})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 2, len(table.rowsFooterRaw))
	assert.Equal(t, 0, len(table.rowsHeaderRaw))
}

func TestTable_AppendHeader(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, len(table.rowsHeaderRaw))

	table.AppendHeader([]interface{}{})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 0, len(table.rowsFooterRaw))
	assert.Equal(t, 1, len(table.rowsHeaderRaw))

	table.AppendHeader([]interface{}{})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 0, len(table.rowsFooterRaw))
	assert.Equal(t, 2, len(table.rowsHeaderRaw))
}

func TestTable_AppendRow(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, table.Length())

	table.AppendRow([]interface{}{})
	assert.Equal(t, 1, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))

	table.AppendRow([]interface{}{})
	assert.Equal(t, 2, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))
}

func TestTable_AppendRows(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, table.Length())

	table.AppendRows([]Row{{}})
	assert.Equal(t, 1, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))

	table.AppendRows([]Row{{}})
	assert.Equal(t, 2, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))
}

func TestTable_Length(t *testing.T) {
	table := Table{}
	assert.Zero(t, table.Length())

	table.AppendRow(testRows[0])
	assert.Equal(t, 1, table.Length())
	table.AppendRow(testRows[1])
	assert.Equal(t, 2, table.Length())

	table.AppendHeader(testHeader)
	assert.Equal(t, 2, table.Length())
}

func TestTable_SetAlign(t *testing.T) {
	table := Table{}
	assert.Nil(t, table.align)

	table.SetAlign([]text.Align{})
	assert.NotNil(t, table.align)

	table.AppendHeader(testHeader)
	table.SetAlignHeader([]text.Align{text.AlignRight, text.AlignRight, text.AlignLeft, text.AlignRight, text.AlignRight})
	table.AppendRows(testRows)
	table.AppendRow(testRowMultiLine)
	table.SetAlign([]text.Align{text.AlignDefault, text.AlignLeft, text.AlignLeft, text.AlignRight, text.AlignRight})
	table.AppendFooter(testFooter)
	table.SetAlignFooter([]text.Align{text.AlignRight, text.AlignRight, text.AlignRight, text.AlignRight, text.AlignRight})

	expectedOut := `+-----+------------+-----------+--------+-----------------------------+
|   # | FIRST NAME | LAST NAME | SALARY |                             |
+-----+------------+-----------+--------+-----------------------------+
|   1 | Arya       | Stark     |   3000 |                             |
|  20 | Jon        | Snow      |   2000 | You know nothing, Jon Snow! |
| 300 | Tyrion     | Lannister |   5000 |                             |
|   0 | Winter     | Is        |      0 |                     Coming. |
|     |            |           |        |        The North Remembers! |
|     |            |           |        |              This is known. |
+-----+------------+-----------+--------+-----------------------------+
|     |            |     TOTAL |  10000 |                             |
+-----+------------+-----------+--------+-----------------------------+`
	assert.Equal(t, expectedOut, table.Render())
}

func TestTable_SetAllowedColumnLengths(t *testing.T) {
	table := Table{}
	table.AppendRows(testRows)
	table.SetStyle(styleTest)

	expectedOut := `(-----^--------^-----------^------^-----------------------------)
[<  1>|<Arya  >|<Stark    >|<3000>|<                           >]
[< 20>|<Jon   >|<Snow     >|<2000>|<You know nothing, Jon Snow!>]
[<300>|<Tyrion>|<Lannister>|<5000>|<                           >]
\-----v--------v-----------v------v-----------------------------/`
	assert.Empty(t, table.allowedColumnLengths)
	assert.Equal(t, expectedOut, table.Render())

	table.SetAllowedColumnLengths([]int{0, 1, 2, 3, 7})
	expectedOut = `(-----^---^----^-----^---------)
[<  1>|<A>|<St>|<300>|<       >]
[<   >|<r>|<ar>|<  0>|<       >]
[<   >|<y>|<k >|<   >|<       >]
[<   >|<a>|<  >|<   >|<       >]
[< 20>|<J>|<Sn>|<200>|<You kno>]
[<   >|<o>|<ow>|<  0>|<w nothi>]
[<   >|<n>|<  >|<   >|<ng, Jon>]
[<   >|< >|<  >|<   >|< Snow! >]
[<300>|<T>|<La>|<500>|<       >]
[<   >|<y>|<nn>|<  0>|<       >]
[<   >|<r>|<is>|<   >|<       >]
[<   >|<i>|<te>|<   >|<       >]
[<   >|<o>|<r >|<   >|<       >]
[<   >|<n>|<  >|<   >|<       >]
\-----v---v----v-----v---------/`
	assert.Equal(t, []int{0, 1, 2, 3, 7}, table.allowedColumnLengths)
	assert.Equal(t, expectedOut, table.Render())

	table.SetAllowedColumnLengths([]int{100, 100, 100, 100, 7})
	expectedOut = `(-----^--------^-----------^------^---------)
[<  1>|<Arya  >|<Stark    >|<3000>|<       >]
[< 20>|<Jon   >|<Snow     >|<2000>|<You kno>]
[<   >|<      >|<         >|<    >|<w nothi>]
[<   >|<      >|<         >|<    >|<ng, Jon>]
[<   >|<      >|<         >|<    >|< Snow! >]
[<300>|<Tyrion>|<Lannister>|<5000>|<       >]
\-----v--------v-----------v------v---------/`
	assert.Equal(t, []int{100, 100, 100, 100, 7}, table.allowedColumnLengths)
	assert.Equal(t, expectedOut, table.Render())

	table.AppendRow(
		Row{20, "Jon", "Snow", 2000, text.FgHiRed.Sprint("You know nothing, Jon Snow!")},
	)
	expectedOut = "(-----^--------^-----------^------^---------)\n" +
		"[<  1>|<Arya  >|<Stark    >|<3000>|<       >]\n" +
		"[< 20>|<Jon   >|<Snow     >|<2000>|<You kno>]\n" +
		"[<   >|<      >|<         >|<    >|<w nothi>]\n" +
		"[<   >|<      >|<         >|<    >|<ng, Jon>]\n" +
		"[<   >|<      >|<         >|<    >|< Snow! >]\n" +
		"[<300>|<Tyrion>|<Lannister>|<5000>|<       >]\n" +
		"[< 20>|<Jon   >|<Snow     >|<2000>|<\x1b[91mYou kno\x1b[0m>]\n" +
		"[<   >|<      >|<         >|<    >|<\x1b[91mw nothi\x1b[0m>]\n" +
		"[<   >|<      >|<         >|<    >|<\x1b[91mng, Jon\x1b[0m>]\n" +
		"[<   >|<      >|<         >|<    >|<\x1b[91m Snow!\x1b[0m >]\n" +
		"\\-----v--------v-----------v------v---------/"
	assert.Equal(t, expectedOut, table.Render())
}

func TestTable_SetAllowedRowLength(t *testing.T) {
	table := Table{}
	table.AppendRows(testRows)
	table.SetStyle(styleTest)

	expectedOutWithNoRowLimit := `(-----^--------^-----------^------^-----------------------------)
[<  1>|<Arya  >|<Stark    >|<3000>|<                           >]
[< 20>|<Jon   >|<Snow     >|<2000>|<You know nothing, Jon Snow!>]
[<300>|<Tyrion>|<Lannister>|<5000>|<                           >]
\-----v--------v-----------v------v-----------------------------/`
	assert.Zero(t, table.allowedRowLength)
	assert.Equal(t, expectedOutWithNoRowLimit, table.Render())

	table.SetAllowedRowLength(utf8.RuneCountInString(table.style.Box.UnfinishedRow))
	assert.Equal(t, utf8.RuneCountInString(table.style.Box.UnfinishedRow), table.allowedRowLength)
	assert.Equal(t, "", table.Render())

	table.SetAllowedRowLength(5)
	expectedOutWithRowLimit := `( ~~~
[ ~~~
[ ~~~
[ ~~~
\ ~~~`
	assert.Equal(t, 5, table.allowedRowLength)
	assert.Equal(t, expectedOutWithRowLimit, table.Render())

	table.SetAllowedRowLength(30)
	expectedOutWithRowLimit = `(-----^--------^---------- ~~~
[<  1>|<Arya  >|<Stark     ~~~
[< 20>|<Jon   >|<Snow      ~~~
[<300>|<Tyrion>|<Lannister ~~~
\-----v--------v---------- ~~~`
	assert.Equal(t, 30, table.allowedRowLength)
	assert.Equal(t, expectedOutWithRowLimit, table.Render())

	table.SetAllowedRowLength(300)
	assert.Equal(t, 300, table.allowedRowLength)
	assert.Equal(t, expectedOutWithNoRowLimit, table.Render())
}

func TestTable_SetAutoIndex(t *testing.T) {
	table := Table{}
	table.AppendRows(testRows)
	table.SetStyle(styleTest)

	expectedOut := `(-----^--------^-----------^------^-----------------------------)
[<  1>|<Arya  >|<Stark    >|<3000>|<                           >]
[< 20>|<Jon   >|<Snow     >|<2000>|<You know nothing, Jon Snow!>]
[<300>|<Tyrion>|<Lannister>|<5000>|<                           >]
\-----v--------v-----------v------v-----------------------------/`
	assert.False(t, table.autoIndex)
	assert.Equal(t, expectedOut, table.Render())

	table.SetAutoIndex(true)
	expectedOut = `(---^-----^--------^-----------^------^-----------------------------)
[< >|< A >|<   B  >|<    C    >|<  D >|<             E             >]
{---+-----+--------+-----------+------+-----------------------------}
[<1>|<  1>|<Arya  >|<Stark    >|<3000>|<                           >]
[<2>|< 20>|<Jon   >|<Snow     >|<2000>|<You know nothing, Jon Snow!>]
[<3>|<300>|<Tyrion>|<Lannister>|<5000>|<                           >]
\---v-----v--------v-----------v------v-----------------------------/`
	assert.True(t, table.autoIndex)
	assert.Equal(t, expectedOut, table.Render())

	table.AppendHeader(testHeader)
	expectedOut = `(---^-----^------------^-----------^--------^-----------------------------)
[< >|<  #>|<FIRST NAME>|<LAST NAME>|<SALARY>|<                           >]
{---+-----+------------+-----------+--------+-----------------------------}
[<1>|<  1>|<Arya      >|<Stark    >|<  3000>|<                           >]
[<2>|< 20>|<Jon       >|<Snow     >|<  2000>|<You know nothing, Jon Snow!>]
[<3>|<300>|<Tyrion    >|<Lannister>|<  5000>|<                           >]
\---v-----v------------v-----------v--------v-----------------------------/`
	assert.True(t, table.autoIndex)
	assert.Equal(t, expectedOut, table.Render())

	table.AppendRow(testRowMultiLine)
	expectedOut = `(---^-----^------------^-----------^--------^-----------------------------)
[< >|<  #>|<FIRST NAME>|<LAST NAME>|<SALARY>|<                           >]
{---+-----+------------+-----------+--------+-----------------------------}
[<1>|<  1>|<Arya      >|<Stark    >|<  3000>|<                           >]
[<2>|< 20>|<Jon       >|<Snow     >|<  2000>|<You know nothing, Jon Snow!>]
[<3>|<300>|<Tyrion    >|<Lannister>|<  5000>|<                           >]
[<4>|<  0>|<Winter    >|<Is       >|<     0>|<Coming.                    >]
[< >|<   >|<          >|<         >|<      >|<The North Remembers!       >]
[< >|<   >|<          >|<         >|<      >|<This is known.             >]
\---v-----v------------v-----------v--------v-----------------------------/`
	assert.Equal(t, expectedOut, table.Render())

	table.SetStyle(StyleLight)
	expectedOut = `┌───┬─────┬────────────┬───────────┬────────┬─────────────────────────────┐
│   │   # │ FIRST NAME │ LAST NAME │ SALARY │                             │
├───┼─────┼────────────┼───────────┼────────┼─────────────────────────────┤
│ 1 │   1 │ Arya       │ Stark     │   3000 │                             │
│ 2 │  20 │ Jon        │ Snow      │   2000 │ You know nothing, Jon Snow! │
│ 3 │ 300 │ Tyrion     │ Lannister │   5000 │                             │
│ 4 │   0 │ Winter     │ Is        │      0 │ Coming.                     │
│   │     │            │           │        │ The North Remembers!        │
│   │     │            │           │        │ This is known.              │
└───┴─────┴────────────┴───────────┴────────┴─────────────────────────────┘`
	assert.Equal(t, expectedOut, table.Render())
}

func TestTable_SetCaption(t *testing.T) {
	table := Table{}
	assert.Empty(t, table.caption)

	table.SetCaption(testCaption)
	assert.NotEmpty(t, table.caption)
	assert.Equal(t, testCaption, table.caption)
}

func TestTable_SetColors(t *testing.T) {
	table := Table{}
	assert.Empty(t, table.colors)
	assert.Empty(t, table.colorsFooter)
	assert.Empty(t, table.colorsHeader)

	table.SetColors([]text.Colors{testColorWoB, testColorBoW})
	assert.NotEmpty(t, table.colors)
	assert.Empty(t, table.colorsFooter)
	assert.Empty(t, table.colorsHeader)
	assert.Equal(t, 2, len(table.colors))
}

func TestTable_SetColorsFooter(t *testing.T) {
	table := Table{}
	assert.Empty(t, table.colors)
	assert.Empty(t, table.colorsFooter)
	assert.Empty(t, table.colorsHeader)

	table.SetColorsFooter([]text.Colors{testColorWoB, testColorBoW})
	assert.Empty(t, table.colors)
	assert.NotEmpty(t, table.colorsFooter)
	assert.Empty(t, table.colorsHeader)
	assert.Equal(t, 2, len(table.colorsFooter))
}

func TestTable_SetColorsHeader(t *testing.T) {
	table := Table{}
	assert.Empty(t, table.colors)
	assert.Empty(t, table.colorsFooter)
	assert.Empty(t, table.colorsHeader)

	table.SetColorsHeader([]text.Colors{testColorWoB, testColorBoW})
	assert.Empty(t, table.colors)
	assert.Empty(t, table.colorsFooter)
	assert.NotEmpty(t, table.colorsHeader)
	assert.Equal(t, 2, len(table.colorsHeader))
}

func TestTable_SetColumnConfigs(t *testing.T) {
	table := Table{}
	assert.Empty(t, table.columnConfigs)

	table.SetColumnConfigs([]ColumnConfig{{}, {}, {}})
	assert.NotEmpty(t, table.columnConfigs)
	assert.Equal(t, 3, len(table.columnConfigs))
}

func TestTable_SetHTMLCSSClass(t *testing.T) {
	table := Table{}
	table.AppendRow(testRows[0])
	expectedHTML := `<table class="` + DefaultHTMLCSSClass + `">
  <tbody>
  <tr>
    <td align="right">1</td>
    <td>Arya</td>
    <td>Stark</td>
    <td align="right">3000</td>
  </tr>
  </tbody>
</table>`
	assert.Equal(t, "", table.htmlCSSClass)
	assert.Equal(t, expectedHTML, table.RenderHTML())

	table.SetHTMLCSSClass(testCSSClass)
	assert.Equal(t, testCSSClass, table.htmlCSSClass)
	assert.Equal(t, strings.Replace(expectedHTML, DefaultHTMLCSSClass, testCSSClass, -1), table.RenderHTML())
}

func TestTable_SetOutputMirror(t *testing.T) {
	table := Table{}
	table.AppendRow(testRows[0])
	expectedOut := `+---+------+-------+------+
| 1 | Arya | Stark | 3000 |
+---+------+-------+------+`
	assert.Equal(t, nil, table.outputMirror)
	assert.Equal(t, expectedOut, table.Render())

	mockOutputMirror := &myMockOutputMirror{}
	table.SetOutputMirror(mockOutputMirror)
	assert.Equal(t, mockOutputMirror, table.outputMirror)
	assert.Equal(t, expectedOut, table.Render())
	assert.Equal(t, expectedOut+"\n", mockOutputMirror.mirroredOutput)
}

func TestTable_SePageSize(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, table.pageSize)

	table.SetPageSize(13)
	assert.Equal(t, 13, table.pageSize)
}

func TestTable_SortByColumn(t *testing.T) {
	table := Table{}
	assert.Empty(t, table.sortBy)

	table.SortBy([]SortBy{{Name: "#", Mode: Asc}})
	assert.Equal(t, 1, len(table.sortBy))

	table.SortBy([]SortBy{{Name: "First Name", Mode: Dsc}, {Name: "Last Name", Mode: Asc}})
	assert.Equal(t, 2, len(table.sortBy))
}

func TestTable_SetVAlign(t *testing.T) {
	table := Table{}
	assert.Nil(t, table.vAlign)

	table.SetVAlign([]text.VAlign{})
	assert.NotNil(t, table.vAlign)

	table.AppendHeader(testRowMultiLine)
	table.SetVAlignHeader([]text.VAlign{text.VAlignBottom, text.VAlignBottom, text.VAlignBottom, text.VAlignBottom})
	table.AppendRow(testRowMultiLine)
	table.SetVAlign([]text.VAlign{text.VAlignMiddle, text.VAlignMiddle, text.VAlignMiddle, text.VAlignMiddle})
	table.AppendFooter(testRowMultiLine)
	table.SetVAlignFooter([]text.VAlign{text.VAlignTop, text.VAlignTop, text.VAlignTop, text.VAlignTop})

	expectedOut := `+---+--------+----+---+----------------------+
|   |        |    |   | COMING.              |
|   |        |    |   | THE NORTH REMEMBERS! |
| 0 | WINTER | IS | 0 | THIS IS KNOWN.       |
+---+--------+----+---+----------------------+
|   |        |    |   | Coming.              |
| 0 | Winter | Is | 0 | The North Remembers! |
|   |        |    |   | This is known.       |
+---+--------+----+---+----------------------+
| 0 | WINTER | IS | 0 | COMING.              |
|   |        |    |   | THE NORTH REMEMBERS! |
|   |        |    |   | THIS IS KNOWN.       |
+---+--------+----+---+----------------------+`
	assert.Equal(t, expectedOut, table.Render())
}

func TestTable_SetStyle(t *testing.T) {
	table := Table{}
	assert.NotNil(t, table.Style())
	assert.Equal(t, StyleDefault, *table.Style())

	table.SetStyle(StyleDefault)
	assert.NotNil(t, table.Style())
	assert.Equal(t, StyleDefault, *table.Style())
}
