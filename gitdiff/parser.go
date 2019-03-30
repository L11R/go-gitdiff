func Parse(r io.Reader) ([]*File, error) {
	if err := p.Next(); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}
	var files []*File
		file, err := p.ParseNextFileHeader()
			return files, err
		if err = p.ParseFragments(file); err != nil {
			return files, err
// TODO(bkeyes): consider exporting the parser type with configuration
// this would enable OID validation, p-value guessing, and prefix stripping
// by allowing users to set or override defaults
// parser invariants:
// - methods that parse objects:
//     - start with the parser on the first line of the first object
//     - if returning nil, do not advance
//     - if returning an error, do not advance past the object
//     - if returning an object, advance to the first line after the object
// - any exported parsing methods must initialize the parser by calling Next()
type parser struct {
	r *bufio.Reader
	eof    bool
	lineno int64
	lines  [3]string
}
func (p *parser) ParseNextFileHeader() (*File, error) {
	var file *File
		frag, err := p.ParseFragmentHeader()
		if err != nil {
			// not a valid header, nothing to worry about
			goto NextLine
		}
		if frag != nil {
			return nil, p.Errorf(-1, "patch fragment without file header: %s", frag.Header())
		file, err = p.ParseGitFileHeader()
		if err != nil {
			return nil, err
		}
		if file != nil {
		// check for a "traditional" patch
		file, err = p.ParseTraditionalFileHeader()
		if file != nil {
			return file, nil
		}
	NextLine:
		if err := p.Next(); err != nil {
			if err == io.EOF {
				break
			return nil, err
	return nil, nil
// ParseFragments parses fragments until the next file header or the end of the
// stream and attaches them to the given file.
func (p *parser) ParseFragments(f *File) error {
		frag, err := p.ParseFragment()
		if frag == nil {
			// TODO(bkeyes): this could mean several things:
			//  - binary patch
			//  - reached the next file header
			//  - reached the end of the patch
			return nil
		}
		lines := int64(len(frag.Lines) + 1) // +1 for the header
		if f.IsNew && frag.OldLines > 0 {
			return p.Errorf(-lines, "new file depends on old contents")
		if f.IsDelete && frag.NewLines > 0 {
			return p.Errorf(-lines, "deleted file still has contents")

		f.Fragments = append(f.Fragments, frag)
	}
}

// Next advances the parser by one line. It returns any error encountered while
// reading the line, including io.EOF when the end of stream is reached.
func (p *parser) Next() error {
	if p.eof {
		p.lines[0] = ""
		return io.EOF
	if p.lineno == 0 {
		// on first call to next, need to shift in all lines
		for i := 0; i < len(p.lines)-1; i++ {
			if err := p.shiftLines(); err != nil && err != io.EOF {
				return err
			}
	err := p.shiftLines()
	if err == io.EOF {
		p.eof = p.lines[1] == ""
	} else if err != nil {
		return err
	p.lineno++
func (p *parser) shiftLines() (err error) {
	for i := 0; i < len(p.lines)-1; i++ {
		p.lines[i] = p.lines[i+1]
	p.lines[len(p.lines)-1], err = p.r.ReadString('\n')
// Line returns a line from the parser without advancing it. A delta of 0
// returns the current line, while higher deltas return read-ahead lines. It
// returns an empty string if the delta is higher than the available lines,
// either because of the buffer size or because the parser reached the end of
// the input. Valid lines always contain at least a newline character.
func (p *parser) Line(delta uint) string {
	return p.lines[delta]
func (p *parser) Errorf(delta int64, msg string, args ...interface{}) error {
	return fmt.Errorf("gitdiff: line %d: %s", p.lineno+delta, fmt.Sprintf(msg, args...))
func (p *parser) ParseFragmentHeader() (*Fragment, error) {
	const (
		startMark = "@@ -"
		endMark   = " @@"
	)
	if !strings.HasPrefix(p.Line(0), startMark) {
		return nil, nil
	parts := strings.SplitAfterN(p.Line(0), endMark, 2)
	if len(parts) < 2 {
		return nil, p.Errorf(0, "invalid fragment header")
	f := &Fragment{}
	f.Comment = strings.TrimSpace(parts[1])
	header := parts[0][len(startMark) : len(parts[0])-len(endMark)]
	ranges := strings.Split(header, " +")
	if len(ranges) != 2 {
		return nil, p.Errorf(0, "invalid fragment header")
	var err error
	if f.OldPosition, f.OldLines, err = parseRange(ranges[0]); err != nil {
		return nil, p.Errorf(0, "invalid fragment header: %v", err)
	if f.NewPosition, f.NewLines, err = parseRange(ranges[1]); err != nil {
		return nil, p.Errorf(0, "invalid fragment header: %v", err)
	if err := p.Next(); err != nil && err != io.EOF {
		return nil, err
	return f, nil
func (p *parser) ParseFragment() (*Fragment, error) {
	frag, err := p.ParseFragmentHeader()
	if err != nil {
		return nil, err
	}
	if frag == nil {
		return nil, nil
	}
	if p.Line(0) == "" {
		return nil, p.Errorf(0, "no content following fragment header")
	}

	oldLines, newLines := frag.OldLines, frag.NewLines
	for oldLines > 0 || newLines > 0 {
		line := p.Line(0)
		switch line[0] {
		case '\n':
			fallthrough // newer GNU diff versions create empty context lines
		case ' ':
			oldLines--
			newLines--
			if frag.LinesAdded == 0 && frag.LinesDeleted == 0 {
				frag.LeadingContext++
			} else {
				frag.TrailingContext++
			}
			frag.Lines = append(frag.Lines, FragmentLine{OpContext, line[1:]})
		case '-':
			frag.LinesDeleted++
			oldLines--
			frag.TrailingContext = 0
			frag.Lines = append(frag.Lines, FragmentLine{OpDelete, line[1:]})
		case '+':
			frag.LinesAdded++
			newLines--
			frag.TrailingContext = 0
			frag.Lines = append(frag.Lines, FragmentLine{OpAdd, line[1:]})
		default:
			// this could be "\ No newline at end of file", which we allow
			// only check the prefix because the text changes by locale
			// git also asserts that any translation is at least 12 characters
			if len(line) >= 12 && strings.HasPrefix("\\ ", line) {
				break
			return nil, p.Errorf(0, "invalid fragment line")
		if err := p.Next(); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	if oldLines != 0 || newLines != 0 {
		return nil, p.Errorf(0, "invalid fragment: remaining lines: %d old, %d new", oldLines, newLines)
	return frag, nil
func parseRange(s string) (start int64, end int64, err error) {
	parts := strings.SplitN(s, ",", 2)
	if start, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
		return 0, 0, fmt.Errorf("bad start of range: %s: %v", parts[0], nerr.Err)
		if end, err = strconv.ParseInt(parts[1], 10, 64); err != nil {
			nerr := err.(*strconv.NumError)
			return 0, 0, fmt.Errorf("bad end of range: %s: %v", parts[1], nerr.Err)
		end = 1
	return