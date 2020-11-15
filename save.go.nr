		OnTextChanged: func() {
			if (! assignTo.ReadOnly()) {
				assignTo.SetReadOnly(true)
				if (strings.ToLower(assignTo.Text()) == strings.ToLower(value)) {
					assignTo.SetTextColor(validColor)
				} else {
					assignTo.SetTextColor(errorColor)
					assignTo.AppendText(" --> " + value)
				}
			}
		},


func (mw *MyMainWindow) AfficherVerbSuivant() {
	mw.validButton.Visible() = true
	mw.nextButton.Visible() = false
}
