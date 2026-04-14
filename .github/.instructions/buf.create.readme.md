## Reglas
1. Genera un archivo `.proto` válido (proto3).
2. Usa nombres de `message` en **PascalCase** y campos en **snake_case**.
3. Mantén el orden lógico de los campos y asigna tags consecutivos desde `1`.
4. Si un campo es lista/array, usa `repeated`.
5. Si detectas objetos anidados, crea `message` separados reutilizables.
6. Mapea tipos correctamente:
    - `string` -> `string`
    - `int/int32` -> `int32`
    - `int64` -> `int64`
    - `float64` -> `double`
    - `bool` -> `bool`
    - fecha/hora ISO8601 -> `google.protobuf.Timestamp` (importándolo)
 y agregando `import "google/protobuf/timestamp.proto";` cuando aplique).
10. El `.proto` debe seguir este formato base:
    - `edition = "2023";`
    - `package <dominio>pb.v1;`onal fuera de esos bloques.
    - Imports obligatorios:
      - `import "buf/validate/validate.proto";`
      - `import "google/api/field_behavior.proto";`
11. Aplica validaciones con `buf.validate` en campos de texto cuando sea posible (`min_len`, `max_len` o `cel`).
12. Marca campos obligatorios con `(google.api.field_behavior) = REQUIRED`.
13. Mantén mensajes limpios y consistentes con el estilo del ejemplo (anotaciones en línea por campo).

## Salida esperada (actualizada)
- Un bloque ` ```proto ` con el archivo completo.
- Un bloque ` ```json ` con un JSON equivalente al `message` principal (incluyendo objetos anidados y arrays cuando existan).
- Sin texto adicional fuera de esos bloques.
7. Si hay claves dinámicas tipo mapa, usa `map<key, value>`.
8. No inventes campos que no existan en la entrada.
9. Si hay ambigüedad de tipo, agrega un comentario breve `// revisar tipo`.

## Salida esperada
- Solo código `.proto` en un bloque ```proto```.
- Sin texto adicional fuera del bloque.

## Entrada
[Pega aquí el JSON o la estructura Go]

## Ejemplo
    product.id (string, required): The product ID of.
    product.name (string, required): The name of the product. length >= 2 && length <= 150
    product.description (string, optional): The description of the product. length >= 2 && length <= 150

    ```json
        Product = {
         id: string, 
         name: string, 
         description: string,
         price: float or double
    }
    ```
    ```proto
    edition="2023";

    package productpb.v1;

    import "buf/validate/validate.proto";
    import "google/api/field_behavior.proto";

    message Product {
        // The product ID of
        string id = 1;

        // The name of the product.
        string name = 2 [(buf.validate.field).cel = {
            id: "product.name"
            message: "name must has at least 2 and max 150 characters long"
            expression: "this.size() >= 2 && this.size() <= 150"
        }, 
        (google.api.field_behavior) = REQUIRED
        ]; 

        // The description of the product
        string description = 3 [(buf.validate.field).string = {
            min_len: 2
            max_len: 150
        }];

        double price = 4;
    }
    ```

## Checklist

- Las carpetas donde se guardan los .proto deben seguir la conversion siguiente `package.name/version/archivo.proto` para que satisfaga los nombres de los paquetes `package <dominio>pb.v1;`
- Los archivos .proto van a contener objetos simples y van a importar sus dependencias.
