// До рефакторингу
if r.FormValue("title") == "" && r.FormValue("content") == "" && (r.FormValue("tags") == "" || r.FormValue("tags") == "[]") {
	if _, _, err := r.FormFile("headerImage"); err == http.ErrMissingFile {
		utils.HandleError(log, w, r, "no fields provided for update", nil, http.StatusBadRequest, "body", "At least one field must be provided")
		return
	}
}

// Після рефакторингу
if isTitleContentAndTagsEmpty(r) && isHeaderImageMissing(r) {
	utils.HandleError(log, w, r, "no fields provided for update", nil,
		http.StatusBadRequest, "body", "At least one field must be provided")
	return
}

func isTitleContentAndTagsEmpty(r *http.Request) bool {
    return r.FormValue("title") == "" &&
           r.FormValue("content") == "" &&
           (r.FormValue("tags") == "" || r.FormValue("tags") == "[]")
}

func isHeaderImageMissing(r *http.Request) bool {
    _, _, err := r.FormFile("headerImage")
    return err == http.ErrMissingFile
}

