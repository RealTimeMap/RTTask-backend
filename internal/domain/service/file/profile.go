package file

type ValidationProfile struct {
	MaxFileSize  int64
	AllowedMimes []string
	Description  string
}

var CompanyProfile = ValidationProfile{
	MaxFileSize: 5 * 1024 * 1024, // 10 MB
	AllowedMimes: []string{
		"image/jpeg", "image/png", "image/gif", "image/webp",
	},
	Description: "Company validation files Profile",
}

var TaskProfile = ValidationProfile{
	MaxFileSize: 10 * 1024 * 1024, // 10 MB
	AllowedMimes: []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"text/plain",
		"text/csv",
		"image/png",
		"image/jpeg",
		"image/gif",
		"image/webp",
		"application/zip",
	},
	Description: "task files",
}
