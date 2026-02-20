 Eres un analista de meta
 
 dado el permiso y la descripcion que te voy a dar vas a dar una descripcion corta de en un oracion, label como se leeria el permiso en portugues y el key ejemplo:
```
whatsapp_business_management permite que tu app lea o administre los activos comerciales de WhatsApp que te pertenecen o a los que otros negocios te concedieron acceso. Estos activos comerciales incluyen cuentas de WhatsApp Business, números de teléfono de empresa, plantillas de mensajes, códigos QR y sus mensajes asociados, y suscripciones a webhooks. The allowed usage for this permission is to manage WhatsApp business assets and display WhatsApp Business Account analytics in your customer portal.
Administra activos comerciales de WhatsApp.
Muestra análisis de cuentas de WhatsApp Business en el portal para clientes.```

salida:
{
"whatsapp_business_management": {
 "Label": "Gerenciador WhatsApp Business",
 "Description": "Permite gerenciar ativos WhatsApp Business e mostrar análises no portal do cliente."
}
}

Checklist:
- [x] Debe quedar el contenido ordenado en orden alfabetico de por el Key.
- [x] La descripción es corta y clara.
- [x] El label es una traducción adecuada al portugués
- [x] El key es el mismo que el permiso dado.
- [x] No pueden haber Key repetidos.
- [x] Tienes que generar la misma cantidad de permisos que se pasa no permiso dependienes.
```json
{
  "whatsapp_business_management": {
    "Label": "Gerenciador WhatsApp Business",
    "Description": "Permite gerenciar ativos WhatsApp Business e mostrar análises no portal do cliente."
  }
}
```