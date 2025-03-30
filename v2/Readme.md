# Al crear un cliente la propiedad 'Config.EnvFilePath' debe ser igual al path del archivo .env. Mas detalles en ayuda de propiedad
## Ejemplo de archivo .env
```
## Configuration pra api oficial
# The base URL to send all SDK requests to (default graph.facebook.com).
# This variable should not be used unless necessary for development or special routing needs.
WA_BASE_URL=https://graph.facebook.com

# Your Meta for Developers app Id.
M4D_APP_ID=

# Your Meta for Developers Business app secret.
M4D_APP_SECRET=

# Your WhatsApp phone number Id (sender).
WA_PHONE_NUMBER_ID=XXXXXXXXXXXXXX

# Your WhatsApp business account Id.
WA_BUSINESS_ACCOUNT_ID=XXXXXXXXXXXXXX

# System user access token. Recommended: Do not use a temporary access token.
CLOUD_API_ACCESS_TOKEN=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

# Cloud API version number.
CLOUD_API_VERSION=v22.0

# Customize your incoming webhook listener endpoint. Path should be
# https://{host}/{WEBHOOK_ENDPOINT}. A trailing slash is not added by default,
# so the variable should include that if it's required by your API gateway.
WEBHOOK_ENDPOINT=

# A custom verification token string to validate incoming webhook payloads.
# Needs to match webhook configuration.
WEBHOOK_VERIFICATION_TOKEN=

# Override the default app listener port (port 3000).
LISTENER_PORT=

# Turn on global debug logging
DEBUG=

# The number of request retries after waiting (default 30 retries)
MAX_RETRIES_AFTER_WAIT=

# The timeout period in milliseconds for a request to wait for a response (default 20000ms)
REQUEST_TIMEOUT=
```
# Instalando paquete
```
go get github.com/ecsavigne/client_wa_oficial/v2
ó
go get github.com/ecsavigne/client_wa_oficial/v2@latest

Version 2

go get github.com/ecsavigne/client_wa_oficial/v2/v2@latest
go mod tidy
```