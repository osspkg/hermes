package addons

type ManifestModel struct {
	API         int    `json:"api"`
	Version     string `json:"version"`
	PkgName     string `json:"pkgname"`
	Author      string `json:"author"`
	Support     string `json:"support"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icons       Icons  `json:"icons"`
	Filename    string `json:"filename"`
}
type Icons struct {
	X64  string `json:"x64"`
	X128 string `json:"x128"`
	X256 string `json:"x256"`
}
