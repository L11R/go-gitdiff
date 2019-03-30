	"io"
	const content = "the first line\nthe second line\nthe third line\n"
	t.Run("read", func(t *testing.T) {
		if err := p.Next(); err != nil {
			t.Fatalf("error advancing parser: %v", err)

		line := p.Line(0)
		if err := p.Next(); err != nil {
			t.Fatalf("error advancing parser: %v", err)

		line = p.Line(0)
		if err := p.Next(); err != nil {
			t.Fatalf("error advancing parser: %v", err)
		}
		line = p.Line(0)
		if line != "the third line\n" {
			t.Fatalf("incorrect third line: %s", line)

		if err := p.Next(); err != io.EOF {
			t.Fatalf("expected EOF, but got: %v", err)
	})
	t.Run("peek", func(t *testing.T) {
		p := newParser()

		if err := p.Next(); err != nil {
			t.Fatalf("error advancing parser: %v", err)

		line := p.Line(1)
		if line != "the second line\n" {
		if err := p.Next(); err != nil {
			t.Fatalf("error advancing parser: %v", err)

		line = p.Line(0)
		if line != "the second line\n" {
func TestParserAdvancment(t *testing.T) {
	tests := map[string]struct {
		Input   string
		Parse   func(p *parser) error
		EndLine string
	}{
		"ParseGitFileHeader": {
			Input: `diff --git a/dir/file.txt b/dir/file.txt
index 9540595..30e6333 100644
--- a/dir/file.txt
+++ b/dir/file.txt
@@ -1,2 +1,3 @@
context line
`,
			Parse: func(p *parser) error {
				_, err := p.ParseGitFileHeader()
				return err
			},
			EndLine: "@@ -1,2 +1,3 @@\n",
		},
		"ParseTraditionalFileHeader": {
			Input: `--- dir/file.txt
+++ dir/file.txt
@@ -1,2 +1,3 @@
context line
`,
			Parse: func(p *parser) error {
				_, err := p.ParseTraditionalFileHeader()
				return err
			},
			EndLine: "@@ -1,2 +1,3 @@\n",
		},
		"ParseFragmentHeader": {
			Input: `@@ -1,2 +1,3 @@
context line
`,
			Parse: func(p *parser) error {
				_, err := p.ParseFragmentHeader()
				return err
			},
			EndLine: "context line\n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p := &parser{r: bufio.NewReader(strings.NewReader(test.Input))}
			p.Next()

			if err := test.Parse(p); err != nil {
				t.Fatalf("unexpected error while parsing: %v", err)
			}

			if test.EndLine != p.Line(0) {
				t.Errorf("incorrect position after parsing\nexpected: %q\nactual: %q", test.EndLine, p.Line(0))
			}
		})
	}
}

		"trailingComment": {
			Input: "@@ -21,5 +28,9 @@ func test(n int) {\n",
				Comment:     "func test(n int) {",
			p := &parser{r: bufio.NewReader(strings.NewReader(test.Input))}
			p.Next()

			frag, err := p.ParseFragmentHeader()
			if !reflect.DeepEqual(test.Output, frag) {
				t.Fatalf("incorrect fragment\nexpected: %+v\nactual: %+v", test.Output, frag)