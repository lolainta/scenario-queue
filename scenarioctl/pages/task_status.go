package pages

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

var allTaskStatuses = []TaskStatus{
	TaskStatusPending,
	TaskStatusInProgress,
	TaskStatusCompleted,
	TaskStatusFailed,
}

func (s TaskStatus) Label() string {
	switch s {
	case TaskStatusPending:
		return "Pending"
	case TaskStatusInProgress:
		return "In Progress"
	case TaskStatusCompleted:
		return "Completed"
	case TaskStatusFailed:
		return "Failed"
	default:
		return string(s)
	}
}

func TaskStatusLabel(value string) string {
	return TaskStatus(value).Label()
}

func TaskStatusOptions() []SelectOption {
	options := make([]SelectOption, 0, len(allTaskStatuses))
	for _, status := range allTaskStatuses {
		options = append(options, SelectOption{Label: status.Label(), Value: string(status)})
	}
	return options
}

func TaskStatusIndex(value string) int {
	for i, status := range allTaskStatuses {
		if string(status) == value {
			return i
		}
	}
	return 0
}
