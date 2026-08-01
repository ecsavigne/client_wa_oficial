package metastatus

type META_STATUS string

const (
	META_CONNECTED    META_STATUS = "CONNECTED"
	META_FLAGGED      META_STATUS = "FLAGGED"
	META_RATE_LIMITED META_STATUS = "RATE_LIMITED"
	META_PENDING      META_STATUS = "PENDING"
	META_OFFLINE      META_STATUS = "OFFLINE"
	META_UNKNOWN      META_STATUS = "UNKNOWN"
)

func (m META_STATUS) Code() int {
	code := map[META_STATUS]int{
		META_UNKNOWN:      1,
		META_PENDING:      2,
		META_CONNECTED:    3,
		META_FLAGGED:      4,
		META_RATE_LIMITED: 5,
		META_OFFLINE:      6,
	}

	if _, ok := code[m]; !ok {
		return code[META_UNKNOWN]
	}

	return code[m]
}

func (m META_STATUS) Description() string {
	description := map[META_STATUS]string{
		META_UNKNOWN:      "Status não reconhecido",
		META_CONNECTED:    "O número está registrado e pronto para enviar e receber mensagens.",
		META_FLAGGED:      "O número foi marcado pelo Meta e requer revisão ou apresenta algum problema.",
		META_RATE_LIMITED: "O número tem restrições temporárias devido a limites de envio.",
		META_PENDING:      "O número ainda não completou o processo de ativação ou registro.",
		META_OFFLINE:      "O número tem restrições que limitam determinadas operações.",
	}

	if _, ok := description[m]; !ok {
		return description[META_UNKNOWN]
	}

	return description[m]
}
