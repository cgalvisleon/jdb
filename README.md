# JDB - Go Database Library

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/Version-v0.1.21-orange.svg)](version.sh)

JDB es una librería de Go que proporciona una interfaz unificada y simplificada para trabajar con múltiples bases de datos. Ofrece soporte para PostgreSQL, MySQL y SQLite con una API consistente y fácil de usar.

## � Tabla de Contenidos

- [�🚀 Características](#-características)
- [📦 Instalación](#-instalación)
- [🔧 Configuración](#-configuración)
- [📖 Uso Básico](#-uso-básico)
- [🏗️ Estructura del Proyecto](#️-estructura-del-proyecto)
- [🔌 Drivers Soportados](#-drivers-soportados)
- [🎯 Ejemplos Avanzados](#-ejemplos-avanzados)
- [🚀 Compilación y Ejecución](#-compilación-y-ejecución)
- [📚 API Reference](#-api-reference)
- [🤝 Contribuir](#-contribuir)
- [📄 Licencia](#-licencia)

## 🚀 Características

- **Multi-driver**: Soporte para PostgreSQL, MySQL y SQLite
- **API Unificada**: Interfaz consistente independientemente del motor de base de datos
- **ORM Simplificado**: Definición de modelos y esquemas de manera declarativa
- **Transacciones**: Soporte completo para transacciones
- **Eventos**: Sistema de eventos para hooks antes y después de operaciones
- **Auditoría**: Sistema de auditoría automática
- **CQRS**: Soporte para Command Query Responsibility Segregation
- **Core System**: Sistema de metadatos y gestión de modelos
- **Debug Mode**: Modo de depuración para desarrollo
- **Gestión de Usuarios**: Creación y gestión de usuarios de base de datos
- **JavaScript VM**: Integración con Goja para scripts dinámicos
- **Sistema de Eventos Avanzado**: Emisión y manejo de eventos personalizados
- **Configuración Dinámica**: Configuración en tiempo de ejecución
- **Query Language**: Sistema de consultas avanzado con soporte para JOIN, agregaciones y filtros complejos

## 📦 Instalación

```bash
go get github.com/cgalvisleon/jdb
```

### Dependencias

```bash
go get github.com/cgalvisleon/et@v1.0.10
```

### Dependencias Principales

- **PostgreSQL**: `github.com/lib/pq v1.10.9`
- **MySQL**: `github.com/go-sql-driver/mysql v1.9.2`
- **SQLite**: `modernc.org/sqlite v1.37.1`
- **Utilidades**: `github.com/cgalvisleon/et v0.1.18`

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

## 🏗️ Estructura del Proyecto

```
jdb/
├── jdb/                    # Paquete principal (57 archivos)
│   ├── jdb.go             # Punto de entrada principal
│   ├── database.go        # Gestión de conexiones
│   ├── model.go           # Definición de modelos
│   ├── model-define.go    # Definición de campos especiales
│   ├── model-new.go       # Generación de datos
│   ├── command*.go        # Comandos CRUD (múltiples archivos)
│   ├── ql*.go            # Query Language (múltiples archivos)
│   ├── core*.go          # Sistema core (múltiples archivos)
│   ├── schema.go         # Gestión de esquemas
│   ├── tx.go             # Transacciones
│   ├── where.go          # Condiciones WHERE
│   ├── event.go          # Sistema de eventos
│   └── ...
├── drivers/               # Drivers de base de datos
│   ├── postgres/         # Driver PostgreSQL (22 archivos)
│   │   ├── database.go   # Conexión PostgreSQL
│   │   ├── users.go      # Gestión de usuarios
│   │   ├── schemas.go    # Gestión de esquemas
│   │   ├── ddl-*.go      # DDL operations
│   │   ├── sql-*.go      # SQL operations
│   │   └── ...
│   ├── mysql/            # Driver MySQL (21 archivos)
│   │   ├── database.go   # Conexión MySQL
│   │   ├── users.go      # Gestión de usuarios
│   │   ├── ddl-*.go      # DDL operations
│   │   ├── sql-*.go      # SQL operations
│   │   └── ...
│   └── sqlite/           # Driver SQLite (20 archivos)
│       ├── database.go   # Conexión SQLite
│       ├── ddl-*.go      # DDL operations
│       ├── sql-*.go      # SQL operations
│       └── ...
├── cmd/                  # Aplicación de ejemplo
│   └── main.go          # Ejemplo de uso
├── go.mod               # Dependencias del módulo
├── go.sum               # Checksums de dependencias
├── version.sh           # Script de gestión de versiones
├── .env.local           # Variables de entorno locales
└── README.md            # Este archivo
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
gofmt -w . && go run --race ./cmd
gofmt -w . && go run ./cmd
gofmt -w . && go run ./cmd/cli run --daemon
gofmt -w . && go run ./cmd/cli ps
gofmt -w . && go run ./cmd/cli send ping
gofmt -w . && go run ./cmd/cli send status
gofmt -w . && go run ./cmd/cli send status --tcp
gofmt -w . && go run ./cmd/cli send foo
gofmt -w . && go run ./cmd/cli send echo "hola mundo"
gofmt -w . && go run ./cmd/cli stop



grep -fl jdb
ps aux | grep jdb
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

### Módulo Principal

```bash
# Actualizar dependencias
go get -u github.com/cgalvisleon/jdb
go get -u github.com/cgalvisleon/et@v0.1.18

# Verificar versión
go list -m github.com/cgalvisleon/jdb
```

### Gestión de Versiones

```bash
# Incrementar versión de revisión
./version.sh --v

# Incrementar versión menor
./version.sh --n

# Incrementar versión mayor
./version.sh --m

# Ver ayuda
./version.sh --h
```

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
