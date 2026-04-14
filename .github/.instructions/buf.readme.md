# Preparando ambiente para buf

## 1. Instalar set de ferramentas pra buf-cli desde golang

```bash
go install github.com/bufbuild/buf/cmd/buf@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
```


## 4. Comandos para generar código

```bash
# Actualizar dependencias
$ buf dep update

# Validar archivos proto
$ buf lint

# Generar código
$ buf generate
```

## Estructura de directorios generados

- depende del tipo de datos que se quiera generar instagram(ig), facebook, messenger
```
v2/type/<instagram(ig), facebook, messenger, general>/gen/
├── <packages Go>/<version>
│   ├── *.pb.go           (Mensajes compilados)
```

## Checklist

- Cuando actualice dependencia de buf debe ser con `buf dep update`
- Crear archivos de generacion para cada caso.
- Ejecuta el template de generacion segun el caso de uso que se este analisando ej template de generacion para `messenger`: 
```bash
buf generate --template buf.gen.messenger.yaml
```