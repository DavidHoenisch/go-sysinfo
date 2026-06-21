package gosysinfo

type mockReader map[string]string

func (m mockReader) Read(path string) string {
	return m[path]
}
