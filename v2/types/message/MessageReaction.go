package message

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

/*
LA REACCIÓN NO SE ENVIARÁ EN LOS SIGUIENTES CASOS:

	Si el mensaje tiene más de 30 días
	Si se trata de un mensaje de reacción
	Si el mensaje se eliminó

Si el identificador es de un mensaje que se eliminó, este no se entregará.
*/
type MessageReaction struct {
	MessagerKernel
	Reaction `json:"reaction"`
}

func (*MessageReaction) NewReactionMessage(config Messager) *MessageReaction {
	switch v := any(config).(type) {
	case *MessageReaction:
		v.MessagerKernel.parent = v
		return v
	}
	panic("Invalid protocol type, expected *MessageReaction")
}
