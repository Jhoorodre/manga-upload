package models

// UploadResult representa o resultado do upload de um único arquivo.
type UploadResult struct {
	URL      string
	Filename string
	Success  bool
	Error    string
}

// ChapterUploadResult representa o resultado do upload de um capítulo inteiro.
type ChapterUploadResult struct {
	ChapterName   string
	AlbumURL      string
	ImageURLs     []string
	FailedUploads []string
	Success       bool
}

// Chapter representa os metadados de um capítulo no reader.json.
type Chapter struct {
	Title       string              `json:"title"`
	Volume      string              `json:"volume"`
	LastUpdated string              `json:"last_updated"`
	Groups      map[string][]string `json:"groups"`
}

// ReaderJSON representa a estrutura principal do arquivo JSON do mangá.
type ReaderJSON struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Artist      string             `json:"artist"`
	Author      string             `json:"author"`
	Cover       string             `json:"cover"`
	Status      string             `json:"status"`
	Chapters    map[string]Chapter `json:"chapters"`
}

// OrderedReaderJSON is an alternative representation utilized strictly for custom marshaling
// to ensure the chapters map keys are printed in descending order in the final JSON string.
// Go maps are unordered, so we use a map alias with custom marshaling locally when saving.
type OrderedReaderJSON struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Artist      string      `json:"artist"`
	Author      string      `json:"author"`
	Cover       string      `json:"cover"`
	Status      string      `json:"status"`
	Chapters    interface{} `json:"chapters"` // Will be injected as an ordered representation
}
