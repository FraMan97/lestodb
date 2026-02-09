package commands

const (
	CMD_SET     string = "SET"
	CMD_GET     string = "GET"
	CMD_DEL     string = "DEL"
	CMD_BACKUP  string = "BACKUP"
	CMD_RESTORE string = "RESTORE"
)

var CMDS []string = []string{CMD_SET, CMD_GET, CMD_DEL, CMD_BACKUP, CMD_RESTORE}
var RULES map[string]int = map[string]int{
	CMD_SET:     4,
	CMD_GET:     2,
	CMD_DEL:     2,
	CMD_BACKUP:  2,
	CMD_RESTORE: 3,
}
