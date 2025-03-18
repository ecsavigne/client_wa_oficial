package types

type Reaction struct {
	MessageId string `json:"message_id" valid:"required"`
	/*
		EMOJI OF THE REACTION:
		    Se admiten todos los emojis que se admiten en los dispositivos Android y iOS.
			Se admiten los emojis renderizados.
			Si se usan valores Unicode de emojis, deben estar codificados para Java o JavaScript.
			Se puede enviar solo un emoji en un mensaje de reacción.
			Usa una cadena vacía para eliminar un emoji que se envió con anterioridad.
	*/
	Emoji string `json:"emoji" valid:"required"`
}

type MessageReaction struct {
	Messager `json:"messager,omitempty"`
	Header
	Reaction `json:"reaction"`
}

/*
LA REACCIÓN NO SE ENVIARÁ EN LOS SIGUIENTES CASOS:

	Si el mensaje tiene más de 30 días
	Si se trata de un mensaje de reacción
	Si el mensaje se eliminó

Si el identificador es de un mensaje que se eliminó, este no se entregará.
*/
func NewMessageReaction(m *MessageReaction) Messager {
	mk := &messagerKernel{
		Type: MessageTypeReaction,
		m:    m,
	}

	m.Messager = mk
	return m
}
