package metastatus

type META_STATUS string

const (
	META_CONNECTED    META_STATUS = "CONNECTED"
	META_FLAGGED      META_STATUS = "FLAGGED"
	META_RATE_LIMITED META_STATUS = "RATE_LIMITED"
	META_PENDING      META_STATUS = "PENDING"
	META_OFFLINE      META_STATUS = "OFFLINE"
	META_UNKNOWN      META_STATUS = "UNKNOWN"
	// Business
	META_BUSINESS_VERIFIED     META_STATUS = "VERIFIED"
	META_BUSINESS_NOT_VERIFIED META_STATUS = "NOT_VERIFIED"
	META_BUSINESS_PENDING      META_STATUS = "PENDING_SUBMISSION"
	META_BUSINESS_PENDING_INFO META_STATUS = "PENDING_NEED_MORE_INFO"
	META_BUSINESS_REJECTED     META_STATUS = "REJECTED"
	META_BUSINESS_FAILED       META_STATUS = "FAILED"
	META_BUSINESS_INELIGIBLE   META_STATUS = "INELIGIBLE"
	META_BUSINESS_REVOKED      META_STATUS = "REVOKED"
	META_BUSINESS_EXPIRED      META_STATUS = "EXPIRED"
	/*Waba
		REJECTED  (La revisión fue rechazada.)
	APPROVED (La WABA fue aprobada.)
	*/
	META_WABA_APPROVED META_STATUS = "APPROVED"
)

func (m META_STATUS) String() string {
	return string(m)
}

func (m META_STATUS) Code() int {
	code := map[META_STATUS]int{
		META_UNKNOWN:               1,
		META_PENDING:               2,
		META_CONNECTED:             3,
		META_FLAGGED:               4,
		META_RATE_LIMITED:          5,
		META_OFFLINE:               6,
		META_BUSINESS_VERIFIED:     7,
		META_BUSINESS_NOT_VERIFIED: 8,
		META_BUSINESS_PENDING:      9,
		META_BUSINESS_PENDING_INFO: 10,
		META_BUSINESS_REJECTED:     11,
		META_BUSINESS_FAILED:       12,
		META_BUSINESS_INELIGIBLE:   13,
		META_BUSINESS_REVOKED:      14,
		META_BUSINESS_EXPIRED:      15,
		META_WABA_APPROVED:         17,
	}

	if _, ok := code[m]; !ok {
		return code[META_UNKNOWN]
	}

	return code[m]
}

func (m META_STATUS) Description() string {
	description := map[META_STATUS]string{
		META_UNKNOWN:               "Status não reconhecido",
		META_CONNECTED:             "O número está registrado e pronto para enviar e receber mensagens.",
		META_FLAGGED:               "O número foi marcado pelo Meta e requer revisão ou apresenta algum problema.",
		META_RATE_LIMITED:          "O número tem restrições temporárias devido a limites de envio.",
		META_PENDING:               "O número ainda não completou o processo de ativação ou registro. Revisão em andamento.",
		META_OFFLINE:               "O número tem restrições que limitam determinadas operações.",
		META_BUSINESS_VERIFIED:     "A empresa concluiu com sucesso o processo de verificação da Meta.",
		META_BUSINESS_NOT_VERIFIED: "A empresa ainda não iniciou ou não concluiu o processo de verificação.",
		META_BUSINESS_PENDING:      "O processo de verificação foi iniciado, mas a documentação ou as informações ainda não foram enviadas para análise.",
		META_BUSINESS_PENDING_INFO: "A Meta solicitou informações ou documentos adicionais para continuar a análise da empresa.",
		META_BUSINESS_REJECTED:     "A Meta analisou a solicitação e a recusou. Geralmente é necessário corrigir as informações ou enviar uma nova solicitação.",
		META_BUSINESS_FAILED:       "O processo de verificação não pôde ser concluído devido a um erro ou ao não cumprimento dos requisitos.",
		META_BUSINESS_INELIGIBLE:   "A empresa não atende aos requisitos da Meta para realizar a verificação neste momento.",
		META_BUSINESS_REVOKED:      "A empresa já foi verificada anteriormente, mas a Meta revogou essa verificação.",
		META_BUSINESS_EXPIRED:      "A verificação perdeu a validade e precisa ser renovada ou refeita.",
		META_WABA_APPROVED:         "A WABA foi aprovada.",
	}

	if _, ok := description[m]; !ok {
		return description[META_UNKNOWN]
	}

	return description[m]
}
