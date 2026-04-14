# Preparando ambiente pra buf para instagram 

## 2. Configuración buf.yaml

Para más detalles, visita: https://buf.build/docs/configuration/v2/buf-yaml

**Ubicación actual:** `./buf.yaml`

```yaml
version: v2
modules:
  - path: v2/types/ig/proto
deps: 
  - buf.build/bufbuild/protovalidate
  - buf.build/googleapis/googleapis
lint:
  use:
    - STANDARD
breaking:
  use:
    - FILE
```

### Parámetros principales:

- **modules**: Directorio donde se encuentran los archivos proto (`v2/types/ig/proto`)
- **deps**: Dependencias de protobuf que se utilizan
- **lint**: Linter estándar para validar los archivos proto
- **breaking**: Verificación de cambios incompatibles en la API
## 3. Configuración buf.gen.yaml

Para más detalles, visita: https://buf.build/plugins/protobuf

**Ubicación actual:** `./buf.gen.yaml`

```yaml
version: v2
plugins:
  # Plugin: Protocol Buffers para Go
  # Genera definiciones de mensajes proto para Go
  - remote: buf.build/protocolbuffers/go
    out: v2/types/ig/gen
    opt:
      - paths=source_relative
      - default_api_level=API_OPAQUE
      - Mprotoc-gen-openapiv2/options/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options
    
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/ecsavigne/client_wa_oficial/v2/types
  disable:
    - file_option: go_package 
      module: buf.build/bufbuild/protovalidate
    - module: buf.build/googleapis/googleapis
```

### Parámetros principales:

- **plugins**: Lista de generadores a utilizar
  - **protocolbuffers/go**: Genera código Go para mensajes

- **managed**: Configuración de opciones de archivo administradas
  - **go_package_prefix**: Prefijo del paquete Go para los archivos generados
  - **disable**: Deshabilita opciones para módulos específicos
  - **Si un archivo proto importa algun mensaje que esta en otro `.proto`**: debe incluir la `opt` del plugin de generacion M ejem:
  ```
Mmessengerpb/v1/messenger_flow.proto=github.com/ecsavigne/client_wa_oficial/v2/types/messenger/gen/messengerpb/v1
  ``` 
  para importar las definiciones en el archivo `.go` 
  
 ``` v1 "github.com/ecsavigne/client_wa_oficial/v2/types/messenger/gen/messengerpb/v1"
```


## CheckList

- asegurarte que los archivos generados y los `.proto` sean para estructuras que se encuentran en folder `v2/types/ig`
-  solo si es contexto de tipos de instagram
-  mostrar cuales son los mensajes a generar y las estructuras de golang antes de crear, luego usuario **decide**
-  buf.gen.yaml no puede tener `include_imports: true`
-  los archivos `buf.yaml y buf.gen.yaml` deben seguir estrictamente las instruciones dada aqui.
-  Si se importa algun mensaje dentro de otro asegurase de que sea creadas correctamente las opciones del plugin de generacion para hacer la importacion correctamente.
-  Asegurate de que cada checklist se cumpla antes de generar y crear codigos.
-  la session de 
  ```manager  override:
    - file_option: go_package_prefix
    value: ...
   ``` 
   incluya el path obsoluto del modulo