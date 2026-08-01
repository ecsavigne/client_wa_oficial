package phonestatus

type NAME_STATUS string

const (
	NAME_UNKNOWN                  NAME_STATUS = "UNKNOWN"
	NAME_NON_EXISTS               NAME_STATUS = "NON_EXISTS"
	NAME_APPROVED                 NAME_STATUS = "APPROVED"
	NAME_AVAILABLE_WITHOUT_REVIEW NAME_STATUS = "AVAILABLE_WITHOUT_REVIEW"
	NAME_DECLINED                 NAME_STATUS = "DECLINED"
	NAME_EXPIRED                  NAME_STATUS = "EXPIRED"
	NAME_PENDING_REVIEW           NAME_STATUS = "PENDING_REVIEW"
)

func (n NAME_STATUS) String() string {
	return string(n)
}

func (n NAME_STATUS) Description() string {
	description := map[NAME_STATUS]string{
		NAME_NON_EXISTS:               "Não existe certificado nem solicitação de verificação.",
		NAME_APPROVED:                 "O nome comercial foi aprovado, o certificado está disponível para download.",
		NAME_AVAILABLE_WITHOUT_REVIEW: "O certificado está pronto sem revisão adicional.",
		NAME_DECLINED:                 "A verificação do nome comercial foi rejeitada.",
		NAME_EXPIRED:                  "O certificado expirou e precisa ser renovado.",
		NAME_PENDING_REVIEW:           "A verificação do nome está em revisão.",
	}

	if _, ok := description[n]; !ok {
		return description[NAME_UNKNOWN]
	}

	return description[n]
}

func (n NAME_STATUS) Code() int {
	code := map[NAME_STATUS]int{
		NAME_NON_EXISTS:               1,
		NAME_APPROVED:                 2,
		NAME_AVAILABLE_WITHOUT_REVIEW: 3,
		NAME_DECLINED:                 4,
		NAME_EXPIRED:                  5,
		NAME_PENDING_REVIEW:           6,
	}

	if _, ok := code[n]; !ok {
		return code[NAME_UNKNOWN]
	}

	return code[n]
}

type QUALITY_RATING string

const (
	QUALITY_UNKNOWN QUALITY_RATING = "UNKNOWN"
	QUALITY_GREEN   QUALITY_RATING = "GREEN"
	QUALITY_YELLOW  QUALITY_RATING = "YELLOW"
	QUALITY_RED     QUALITY_RATING = "RED"
)

func (q QUALITY_RATING) String() string {
	return string(q)
}

func (q QUALITY_RATING) Description() string {
	description := map[QUALITY_RATING]string{
		QUALITY_UNKNOWN: "A classificação de qualidade ainda não foi determinada (novos números de telefone).",
		QUALITY_GREEN:   "Alta qualidade - As mensagens estão sendo entregues e estão interagindo bem.",
		QUALITY_YELLOW:  "Qualidade média - Alguns problemas de entrega ou interação foram detectados.",
		QUALITY_RED:     "Baixa qualidade - Problemas significativos de entrega ou interação.",
	}

	if _, ok := description[q]; !ok {
		return description[QUALITY_UNKNOWN]
	}

	return description[q]
}

func (q QUALITY_RATING) Code() int {
	code := map[QUALITY_RATING]int{
		QUALITY_UNKNOWN: 1,
		QUALITY_GREEN:   2,
		QUALITY_YELLOW:  3,
		QUALITY_RED:     4,
	}

	if _, ok := code[q]; !ok {
		return code[QUALITY_UNKNOWN]
	}

	return code[q]
}

type CODE_VERIFICATION_STATUS string

const (
	CODE_VERIFICATION_UNKNOWN      CODE_VERIFICATION_STATUS = "UNKNOWN"
	CODE_VERIFICATION_VERIFIED     CODE_VERIFICATION_STATUS = "VERIFIED"
	CODE_VERIFICATION_NOT_VERIFIED CODE_VERIFICATION_STATUS = "NOT_VERIFIED"
	CODE_VERIFICATION_EXPIRED      CODE_VERIFICATION_STATUS = "EXPIRED"
)

func (c CODE_VERIFICATION_STATUS) String() string {
	return string(c)
}

func (c CODE_VERIFICATION_STATUS) Description() string {
	description := map[CODE_VERIFICATION_STATUS]string{
		CODE_VERIFICATION_UNKNOWN:      "A classificação de qualidade ainda não foi determinada (novos números de telefone).",
		CODE_VERIFICATION_VERIFIED:     "O número de telefone foi verificado com sucesso.",
		CODE_VERIFICATION_NOT_VERIFIED: "O número de telefone não foi verificado.",
		CODE_VERIFICATION_EXPIRED:      "O status de verificação do número de telefone expirou.",
	}

	if _, ok := description[c]; !ok {
		return description[CODE_VERIFICATION_UNKNOWN]
	}

	return description[c]
}

func (c CODE_VERIFICATION_STATUS) Code() int {
	code := map[CODE_VERIFICATION_STATUS]int{
		CODE_VERIFICATION_UNKNOWN:      1,
		CODE_VERIFICATION_VERIFIED:     2,
		CODE_VERIFICATION_NOT_VERIFIED: 3,
		CODE_VERIFICATION_EXPIRED:      4,
	}

	if _, ok := code[c]; !ok {
		return code[CODE_VERIFICATION_UNKNOWN]
	}

	return code[c]
}
