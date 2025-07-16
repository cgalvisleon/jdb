# JDB - Go Database Library

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/Version-v0.1.18-orange.svg)](version.sh)

JDB es una librería de Go que proporciona una interfaz unificada y simplificada para trabajar con múltiples bases de datos. Ofrece soporte para PostgreSQL, MySQL, SQLite y Oracle con una API consistente y fácil de usar.

## 🚀 Características

- **Multi-driver**: Soporte para PostgreSQL, MySQL, SQLite y Oracle
- **API Unificada**: Interfaz consistente independientemente del motor de base de datos
- **ORM Simplificado**: Definición de modelos y esquemas de manera declarativa
- **Transacciones**: Soporte completo para transacciones
- **Eventos**: Sistema de eventos para hooks antes y después de operaciones
- **Auditoría**: Sistema de auditoría automática
- **CQRS**: Soporte para Command Query Responsibility Segregation
- **Core System**: Sistema de metadatos y gestión de modelos
- **Debug Mode**: Modo de depuración para desarrollo
- **Sistema de Daemon**: Gestión de servicios como daemon
- **Gestión de Usuarios**: Creación y gestión de usuarios de base de datos
- **JavaScript VM**: Integración con Goja para scripts dinámicos
- **Sistema de Eventos Avanzado**: Emisión y manejo de eventos personalizados
- **Gestión de PID**: Control de procesos con archivos PID
- **Configuración Dinámica**: Configuración en tiempo de ejecución

## 📦 Instalación

```bash
go get github.com/cgalvisleon/jdb
```

### Dependencias

```bash
go get github.com/cgalvisleon/et@v0.1.15
```

## 🔧 Configuración

### Variables de Entorno

```bash
# Configuración básica
NODEID=1
DB_NAME=myapp
DB_DRIVER=postgres  # postgres, mysql, sqlite, oracle
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
APP_NAME=myapp

# Configuración adicional
DB_SSL_MODE=disable
DB_TIMEZONE=UTC

# Configuración Oracle específica
ORA_DB_SERVICE_NAME_ORACLE=jdb
ORA_DB_SSL_ORACLE=false
ORA_DB_SSL_VERIFY_ORACLE=false
ORA_DB_VERSION_ORACLE=19
```

## 📖 Uso Básico

### Conexión a Base de Datos

```go
package main

import (
    "github.com/cgalvisleon/jdb"
)

func main() {
    // Configuración de conexión
    params := jdb.ConnectParams{
        Driver:   "postgres",
        Name:     "myapp",
        UserCore: true,
        NodeId:   1,
        Debug:    true,
        Params: jdb.Json{
            "host":     "localhost",
            "port":     5432,
            "user":     "postgres",
            "password": "password",
            "dbname":   "myapp",
        },
    }

    // Conectar a la base de datos
    db, err := jdb.ConnectTo(params)
    if err != nil {
        panic(err)
    }
    defer db.Disconected()

    fmt.Println("Conectado a:", db.Name)
}
```

### Definición de Modelos

```go
// Definir un esquema
schema := db.GetSchema("public")

// Definir un modelo
user := schema.DefineModel("users", "Usuarios del sistema")
user.DefineColumn("id", jdb.TypeDataKey, jdb.PrimaryKey)
user.DefineColumn("name", jdb.TypeDataText, jdb.Required)
user.DefineColumn("email", jdb.TypeDataText, jdb.Unique)
user.DefineColumn("age", jdb.TypeDataInt)
user.DefineColumn("active", jdb.TypeDataBool, jdb.Default(true))
user.DefineColumn("created_at", jdb.TypeDataTime, jdb.Default("NOW()"))

// Campos especiales del sistema
user.DefineCreatedAtField()    // Campo de fecha de creación
user.DefineUpdatedAtField()    // Campo de fecha de actualización
user.DefineStatusField()       // Campo de estado
user.DefineSystemKeyField()    // Campo de clave del sistema
user.DefineIndexField()        // Campo de índice
user.DefineSourceField()       // Campo de origen
user.DefineProjectField()      // Campo de proyecto

// Crear el modelo en la base de datos
err := db.LoadModel(user)
if err != nil {
    panic(err)
}
```

### Operaciones CRUD

```go
// Insertar datos
result, err := db.Command(&jdb.Command{
    Command: jdb.Insert,
    From:    user.GetFrom(),
    Values: []jdb.Json{
        {
            "name":  "Juan Pérez",
            "email": "juan@example.com",
            "age":   30,
        },
    },
})

// Consultar datos
items, err := db.Select(&jdb.Ql{
    From: user.GetFrom(),
    Where: &jdb.QlWhere{
        And: []*jdb.Where{
            {Field: "active", Op: jdb.Eq, Value: true},
        },
    },
})

// Actualizar datos
result, err := db.Command(&jdb.Command{
    Command: jdb.Update,
    From:    user.GetFrom(),
    Values: []jdb.Json{
        {"age": 31},
    },
    QlWhere: &jdb.QlWhere{
        And: []*jdb.Where{
            {Field: "id", Op: jdb.Eq, Value: "user123"},
        },
    },
})

// Eliminar datos
result, err := db.Command(&jdb.Command{
    Command: jdb.Delete,
    From:    user.GetFrom(),
    QlWhere: &jdb.QlWhere{
        And: []*jdb.Where{
            {Field: "id", Op: jdb.Eq, Value: "user123"},
        },
    },
})
```

### Bulk Insert

```go
// Inserción masiva
result, err := db.Command(&jdb.Command{
    Command: jdb.Bulk,
    From:    user.GetFrom(),
    Data: []jdb.Json{
        {"name": "Ana García", "email": "ana@example.com", "age": 25},
        {"name": "Carlos López", "email": "carlos@example.com", "age": 35},
        {"name": "María Rodríguez", "email": "maria@example.com", "age": 28},
    },
})
```

### Transacciones

```go
// Iniciar transacción
tx, err := db.Begin()
if err != nil {
    panic(err)
}
defer tx.Rollback()

// Operaciones en transacción
result, err := tx.Command(&jdb.Command{
    Command: jdb.Insert,
    From:    user.GetFrom(),
    Values: []jdb.Json{
        {"name": "Usuario Transaccional", "email": "tx@example.com"},
    },
})

// Commit de la transacción
err = tx.Commit()
if err != nil {
    panic(err)
}
```

## 🛠️ Sistema de Daemon

JDB incluye un sistema de daemon para gestionar servicios:

### Gestión del Servicio

```bash
# Mostrar ayuda
./jdb help

# Mostrar versión
./jdb version

# Mostrar estado del servicio
./jdb status

# Configurar el servicio
./jdb conf '{"port": 3500, "debug": true}'

# Iniciar el servicio
./jdb start

# Detener el servicio
./jdb stop

# Reiniciar el servicio
./jdb restart
```

### Configuración del Daemon

```go
// Configuración del daemon
config := et.Json{
    "port":  3500,
    "debug": true,
    "host":  "localhost",
}

// Aplicar configuración
daemon.SetConfig(config.ToString())
```

## 👥 Gestión de Usuarios

JDB proporciona funcionalidades para gestionar usuarios de base de datos:

### PostgreSQL

```go
// Crear usuario
err := db.CreateUser("nuevo_usuario", "password123", "password123")

// Cambiar contraseña
err := db.ChangePassword("nuevo_usuario", "nueva_password", "nueva_password")

// Otorgar privilegios
err := db.GrantPrivileges("nuevo_usuario", "myapp")

// Eliminar usuario
err := db.DeleteUser("nuevo_usuario")
```

### MySQL

```go
// Crear usuario
err := db.CreateUser("nuevo_usuario", "password123", "password123")

// Cambiar contraseña
err := db.ChangePassword("nuevo_usuario", "nueva_password", "nueva_password")

// Otorgar privilegios
err := db.GrantPrivileges("nuevo_usuario", "myapp")

// Eliminar usuario
err := db.DeleteUser("nuevo_usuario")
```

## 🎯 Nuevas Funcionalidades

### JavaScript VM Integration

```go
// Ejecutar scripts JavaScript en el modelo
user.vm.Set("customFunction", func(data et.Json) et.Json {
    // Lógica personalizada
    return data
})

// Ejecutar script
result, err := user.vm.RunString(`
    var data = {name: "Juan", age: 30};
    customFunction(data);
`)
```

### Sistema de Eventos Avanzado

```go
// Definir eventos personalizados
user.On("custom_event", func(message event.Message) {
    console.Log("Evento personalizado:", message)
})

// Emitir eventos
user.Emit("custom_event", event.Message{
    Type: "user_created",
    Data: et.Json{"user_id": "123"},
})
```

### Generación de Datos de Prueba

```go
// Generar datos de prueba para el modelo
testData := user.New("name", "email", "age")
// Resultado: {"name": "", "email": "", "age": 0}

// Generar datos con valores por defecto
testData := user.New()
// Resultado: {"id": "users:ulid", "name": "", "email": "", "age": 0, "active": true, "created_at": "2024-01-01T00:00:00Z"}
```

### Consultas Avanzadas

```go
// Consulta con campos ocultos
items, err := db.Select(&jdb.Ql{
    From: user.GetFrom(),
    Hidden: []string{"password", "secret_key"},
})

// Consulta con datos de origen
items, err := db.Select(&jdb.Ql{
    From: user.GetFrom(),
    TypeSelect: jdb.Source,
})
```

## 🏗️ Estructura del Proyecto

```
jdb/
├── jdb/                 # Paquete principal
│   ├── database.go      # Gestión de conexiones
│   ├── model.go         # Definición de modelos
│   ├── command.go       # Comandos CRUD
│   ├── ql.go           # Query Language
│   ├── model-new.go     # Generación de datos
│   ├── model-define.go  # Definición de campos especiales
│   └── ...
├── drivers/            # Drivers de base de datos
│   ├── postgres/       # Driver PostgreSQL
│   │   ├── users.go    # Gestión de usuarios
│   │   └── ...
│   ├── mysql/          # Driver MySQL
│   │   ├── users.go    # Gestión de usuarios
│   │   └── ...
│   ├── sqlite/         # Driver SQLite
│   ├── oracle/         # Driver Oracle
│   │   ├── users.go    # Gestión de usuarios
│   │   └── ...
├── cqrs/              # Patrón CQRS
└── cmd/               # Aplicación de ejemplo
    ├── jdb/           # Comando principal
    │   ├── main.go     # Punto de entrada
    │   ├── systemd.go  # Sistema de daemon
    │   ├── pid.go      # Gestión de PID
    │   └── msg.go      # Mensajes del sistema
    └── main.go         # Ejemplo de uso
```

## 🔌 Drivers Soportados

### PostgreSQL

```go
params := jdb.ConnectParams{
    Driver: "postgres",
    Params: jdb.Json{
        "host":     "localhost",
        "port":     5432,
        "user":     "postgres",
        "password": "password",
        "dbname":   "myapp",
        "sslmode":  "disable",
    },
}
```

### MySQL

```go
params := jdb.ConnectParams{
    Driver: "mysql",
    Params: jdb.Json{
        "host":     "localhost",
        "port":     3306,
        "user":     "root",
        "password": "password",
        "dbname":   "myapp",
    },
}
```

### SQLite

```go
params := jdb.ConnectParams{
    Driver: "sqlite",
    Params: jdb.Json{
        "file": "./data.db",
    },
}
```

### Oracle

```go
params := jdb.ConnectParams{
    Driver: "oracle",
    Params: jdb.Json{
        "host":         "localhost",
        "port":         1521,
        "username":     "system",
        "password":     "password",
        "service_name": "XE",
        "ssl":          false,
        "version":      19,
    },
}
```

## 🎯 Ejemplos Avanzados

### Consultas Complejas

```go
// Consulta con JOIN
items, err := db.Select(&jdb.Ql{
    From: user.GetFrom(),
    Joins: []*jdb.QlJoin{
        {
            Type:  jdb.InnerJoin,
            Table: "profiles",
            On: &jdb.QlWhere{
                And: []*jdb.Where{
                    {Field: "users.id", Op: jdb.Eq, Value: "profiles.user_id"},
                },
            },
        },
    },
    Where: &jdb.QlWhere{
        And: []*jdb.Where{
            {Field: "users.active", Op: jdb.Eq, Value: true},
            {Field: "profiles.verified", Op: jdb.Eq, Value: true},
        },
    },
    OrderBy: &jdb.QlOrder{
        Asc: []*jdb.Field{{Name: "users.created_at"}},
    },
    Limit: 10,
})
```

### Eventos y Hooks

```go
// Evento antes de insertar
user.EventsInsert = append(user.EventsInsert, func(model *jdb.Model, before, after jdb.Json) error {
    fmt.Println("Insertando usuario:", after)
    return nil
})

// Evento después de actualizar
user.EventsUpdate = append(user.EventsUpdate, func(model *jdb.Model, before, after jdb.Json) error {
    fmt.Println("Usuario actualizado:", after)
    return nil
})
```

### Campos Especiales

```go
// Definir campo de texto completo
user.DefineFullText("spanish", []string{"name", "description"})

// Definir relación
user.DefineRelation("profile", "profiles", map[string]string{
    "user_id": "id",
}, 1)

// Definir rollup
user.DefineRollup("total_orders", "orders", map[string]string{
    "user_id": "id",
}, "amount")

// Definir objeto
user.DefineObject("address", "addresses", map[string]string{
    "user_id": "id",
}, []string{"street", "city", "country"})
```

## 🚀 Compilación y Ejecución

### Ejecutar en modo desarrollo

```bash
gofmt -w . && go run --race ./cmd/jdb -port 3500
```

### Compilar para producción

```bash
gofmt -w . && go build --race -a -o ./jdb ./cmd/jdb
```

### Ejecutar con configuración personalizada

```bash
gofmt -w . && go run --race ./cmd/jdb -port 3600 -rpc 4600
```

### Gestión de Versiones

```bash
# Incrementar versión de revisión
./version.sh --v

# Incrementar versión menor
./version.sh --n

# Incrementar versión mayor
./version.sh --m
```

## 📚 API Reference

### Version

```bash
git add .
git commit -m 'Set new version'
git push -u origin
git tag v0.1.18
git push origin --tags
```

### Tipos de Datos Soportados

- `TypeDataText` - VARCHAR(250)
- `TypeDataShortText` - VARCHAR(80)
- `TypeDataMemo` - TEXT
- `TypeDataInt` - INTEGER
- `TypeDataNumber` - DECIMAL(18,2)
- `TypeDataBool` - BOOLEAN
- `TypeDataTime` - TIMESTAMP
- `TypeDataObject` - JSONB
- `TypeDataArray` - JSONB
- `TypeDataKey` - VARCHAR(80)
- `TypeDataState` - VARCHAR(20)
- `TypeDataSerie` - BIGINT
- `TypeDataPrecision` - DOUBLE PRECISION
- `TypeDataBytes` - BYTEA
- `TypeDataGeometry` - JSONB
- `TypeDataFullText` - TSVECTOR

### Tipos de ID Soportados

- `TpNodeId` - ID de nodo
- `TpUUId` - UUID
- `TpULId` - ULID
- `TpXId` - XID

### Operadores de Consulta

- `Eq` - Igual
- `Ne` - No igual
- `Gt` - Mayor que
- `Gte` - Mayor o igual que
- `Lt` - Menor que
- `Lte` - Menor o igual que
- `Like` - Como
- `ILike` - Como (case insensitive)
- `In` - En
- `NotIn` - No en
- `IsNull` - Es nulo
- `IsNotNull` - No es nulo

### Comandos del Sistema

- `Insert` - Insertar
- `Update` - Actualizar
- `Delete` - Eliminar
- `Bulk` - Inserción masiva
- `Upsert` - Insertar o actualizar
- `Delsert` - Eliminar e insertar

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

## 🆘 Soporte

Si tienes alguna pregunta o necesitas ayuda, por favor:

1. Revisa la documentación
2. Busca en los issues existentes
3. Crea un nuevo issue con detalles del problema

---

**JDB** - Simplificando el acceso a bases de datos en Go 🚀
