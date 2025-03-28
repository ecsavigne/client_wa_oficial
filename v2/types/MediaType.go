package types

type Media struct {
	// Required when type is audio, document, image, sticker, or video and you are not using a link.
	Id string `json:"id,omitempty"`
	/*
		Required when type is audio, document, image, sticker, or video and you are not using an uploaded media ID (i.e. you are hosting the media asset on your public server).
		The protocol and URL of the media to be sent. Use only with HTTP/HTTPS URLs.
		Do not use this field when message type is set to text.
	*/
	Link string `json:"link,omitempty"`
	/*
			- Media asset caption. Do not use with audio or sticker media
		    - For v2.41.2 or newer, this field is is limited to 1024 characters.
		    - Captions are currently not supported for document media.

	*/
	Caption string `json:"caption,omitempty"`
	/*
		- Describes the filename for the specific document. Use only with document media
		- The extension of the filename will specify what format the document is displayed as in WhatsApp.
	*/
	Filename string `json:"filename,omitempty"`
	/*
		his path is optionally used with a link when the HTTP/HTTPS link is not directly
		accessible and requires additional configurations like a bearer token. For
		information on configuring providers, see the Media url:
		https://developers.facebook.com/docs/whatsapp/api/settings/media-providers
	*/
	Provider string `json:"provider,omitempty"`
}

type MediaType = string

const (
	AUDIO    MediaType = "audio"
	DOCUMENT MediaType = "document"
	IMAGE    MediaType = "image"
	STICKER  MediaType = "sticker"
	VIDEO    MediaType = "video"
)
