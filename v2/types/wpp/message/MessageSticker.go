package message

/*
Un objeto media que contiene un sticker.
se admiten stickers de salida estáticos y animados de terceros,
además de todos los tipos de stickers de entrada. Un sticker estático debe
contener 512 x 512 píxeles y no puede exceder los 100 KB. Un sticker
animado debe contener 512 x 512 píxeles y no puede exceder los 500 KB.
*/
type MessageSticker struct {
	MessagerKernel
	Media `json:"sticker"`
}
