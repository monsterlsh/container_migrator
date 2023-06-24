package dirty

const (
	//checkpointPath = "/opt/migrator/client/myredis-noauto"
	container      = "mysql"
	checkpointPath = "/opt/container-migrator/client_repo/" + container
	pageSize       = 4096
	csvPath        = "/opt/migrator/csv/Trace/" + container
	allpagesPath   = "/opt/container-migrator/csv/Trace/" + container + "/allpage500.csv"
	dirtypagesPath = "/opt/container-migrator/csv/Trace/" + container + "/dirty.csv"
	timelens       = 500
)

type JsonEntry struct {
	CsvEntry []CsvEntry `json:"csventries"`
}
type CsvEntry struct {
	iteration int     `json:"iteration,omitempty"`
	Entries   []Entry `json:"entries"`
}
