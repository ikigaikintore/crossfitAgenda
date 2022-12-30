package http

type StatusEnum int

const (
	Working StatusEnum = iota + 1
	Finished
	Failed
)

const (
	WorkingName  string = "working"
	FinishedName string = "finished"
	FailedName   string = "failed"
)

func (se StatusEnum) ToString() string {
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

type Status struct {
	ID       string     `json:"id"`
	Date     int64      `json:"date"`
	Detail   string     `json:"detail"`
	Status   StatusEnum `json:"status"`
	Complete bool       `json:"complete"`
}
