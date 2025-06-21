# JDB - Go Database Library

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/Version-v0.1.4-orange.svg)](version.sh)

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

## 📦 Instalación

```bash
go get github.com/cgalvisleon/jdb
```

### Dependencias

```bash
go get github.com/cgalvisleon/et@v0.1.4
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

## 🏗️ Estructura del Proyecto

```
jdb/
├── jdb/                 # Paquete principal
│   ├── database.go      # Gestión de conexiones
│   ├── model.go         # Definición de modelos
│   ├── command.go       # Comandos CRUD
│   ├── ql.go           # Query Language
│   └── ...
├── drivers/            # Drivers de base de datos
│   ├── postgres/       # Driver PostgreSQL
│   ├── mysql/          # Driver MySQL
│   ├── sqlite/         # Driver SQLite
│   └── oracle/         # Driver Oracle
├── cqrs/              # Patrón CQRS
└── cmd/               # Aplicación de ejemplo
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

## 📚 API Reference

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
