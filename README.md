# Go Microservice Template

Este es un *template* para la creación rápida de microservicios en Go siguiendo el [Standard Go Project Layout](https://github.com/golang-standards/project-layout) y las mejores prácticas de la industria.

## 🚀 Stack Tecnológico

*   **Lenguaje:** Go ^1.24
*   **Enrutador HTTP:** [Chi](https://github.com/go-chi/chi) (Ligero e idiomático)
*   **Comunicación Interna:** gRPC & Protocol Buffers
*   **Base de Datos / ORM:** PostgreSQL & [GORM](https://gorm.io/)
*   **Configuración:** [Viper](https://github.com/spf13/viper) (.yaml y Variables de Entorno)
*   **Observabilidad/Logging:** [Zap](https://github.com/uber-go/zap) (Logs estructurados en JSON)

## 🏗️ Arquitectura y Buenas Prácticas

Este proyecto está construido con un fuerte énfasis en:
*   **Graceful Shutdown:** Intercepta `SIGTERM` y `SIGINT` para asegurar que las conexiones existenes (HTTP y gRPC) y las transacciones a bases de datos finalicen antes de que muera el proceso.
*   **Context Propagation:** Uso intensivo del `context.Context` desde la base (Handler) hasta el Access Data (DB) permitiendo control de Cancelaciones y Timeouts.
*   **Dependency Injection (DI):** Sin variables globales ni funciones `init()`. Todo se instancia en `cmd/server/main.go` a través de constructores (`NewX()`).
*   **Manejo de Errores Unificado:** Capa agnóstica de errores de dominio en `pkg/errors` que es mapeada automáticamente a Códigos HTTP específicos y Códigos de gRPC.
*   **Health Checks para K8s:** Endpoints expuestos (`/healthz` y `/ready`) diseñados para las *probes* de liveness y readiness de Kubernetes.

## 📂 Estructura de Directorios

```plaintext
.
├── api/             # Contratos de APIs y definición de Protocol Buffers
│   └── proto/       # .proto files (ej: user.proto) y sus auto-generados (*.pb.go)
├── cmd/
│   └── server/      # Entrypoint (func main) y armado/arranque (DI) de la app.
├── internal/        # (Private) Lógica exclusiva de este microservicio
│   ├── config/      # Carga y estructuración de Viper
│   ├── handler/     # Endpoints HTTP (Chi) y gRPC. Punto de entrada de requests.
│   ├── repository/  # Acceso a base de datos (Postgres, GORM, queries).
│   ├── server/      # Wrappers para instanciar Listeners HTTP y gRPC.
│   └── service/     # Lógica central del Negocio (Business logic de la Entidad).
├── pkg/             # Código de propósito general (exportable para reuso)
│   ├── errors/      # Encapsulador de Status y tipos de errores.
│   └── logger/      # Setup global del sistema de logging Zap.
├── Makefile         # Scripts para correr automatismos rápidos
└── config.yaml      # Variables iniciales del ambiente (puertos, DB, etc.)
```

## 🛠️ Cómo Utilizarlo (Guía Rápida)

Para iniciar, asegúrate de tener instalado `go`, `make`, y opcionalmente `docker` (para la base de datos local y compilación protoc).

### 1. Levantar Dependencias
El proyecto necesita PostgreSQL corriendo. Puedes levantar fácilmente una base de datos local temporal usando Docker:
```bash
docker run -d --rm --name go-template-db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=mydb -p 5432:5432 postgres:15
```

### 2. Sincronizar Módulos y Compilar Proto
```bash
# Sincroniza librerías del go.mod
make tidy

# Compila el código auto-generado gRPC en base al user.proto en la carpeta API
make proto
```
*(Nota: El comando `make proto` intentará interactuar con tu binario de protoc local. Si no lo tienes configurado, te sugerimos utilizar la ruta del contenedor docker).*

### 3. Ejecutar la Aplicación
Inicia el microservicio:
```bash
make run
```
> Esto iniciará tanto el servidor gRPC en el puerto `9090` como el servidor HTTP en el `8080` (configurables vía `config.yaml` o variables de entorno).

## 🧪 Puntos de Prueba (Endpoints)

Una vez corriendo, puedes probar cómo interactúan el Service, Handler y el Repositorio de base de datos abriendo otra terminal:

*   **Revisar Health:**
    ```bash
    curl -v http://localhost:8080/healthz
    ```
*   **Crear un Usuario (Post a JSON):**
    ```bash
    curl -X POST http://localhost:8080/users \
       -H "Content-Type: application/json" \
       -d '{"name": "Developer", "email": "dev@correo.com"}'
    ```
*   **Obtener Usuario (Sustituir ID):**
    ```bash
    curl http://localhost:8080/users/UUID-AQUI
    ```

## 📝 Próximos Pasos (Modificando del Template)
1.  **Cambia de Nombre del Módulo**: Sustituye `github.com/juanpblasi/go-template` a tu ruta real ejecutando un Reemplazo (Find & Replace).
2.  **Muta el Dominio:** Cambia la palabra `user` en los `internal/service/`, `handler`, y `repository` por la entidad real de este microservicio.
3.  **Agrega Modelos a gRPC:** Edita la carpeta `api/proto/` y ejecuta `make proto`.
