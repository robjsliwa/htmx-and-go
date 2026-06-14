package web

import "html/template"

type Templates struct {
	T *template.Template
}

func MustLoadTemplates() *Templates {
	t := template.Must(template.ParseFiles(
		"templates/room.html",
		"templates/log_entry.html",
		"templates/partials.html",
		"templates/log_entry_oob.html",
	))
	return &Templates{T: t}
}