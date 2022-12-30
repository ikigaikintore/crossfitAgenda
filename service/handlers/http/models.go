package http

const (
	WorkingName  string = "working"
	FinishedName string = "finished"
	FailedName   string = "failed"
)

func (se ProcessStatusId) ToString() string {
	switch se {
	case Working:
		return WorkingName
	case Finished:
		return FinishedName
	case Failed:
		return FailedName
	}
	return ""
}
