package model

type ObjectMetadata struct {
	Name         string `json:"name"`
	SourceURI    string `json:"sourceURI"`
	LastModified int64  `json:"lastModified"`
	Size         int64  `json:"size"`
	Type         string `json:"type"`
	ThumbnailURI string `json:"thumbnailURI"`
}
