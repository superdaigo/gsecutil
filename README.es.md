# gsecutil - Utilidad de Google Secret Manager

Un contenedor de línea de comandos simplificado para Google Secret Manager que funciona como un administrador de contraseñas por proyecto. Almacene, recupere y administre secretos con comandos intuitivos, integración con el portapapeles, control de versiones, archivos de configuración amigables para equipos y registros de auditoría.

## 🌍 Versiones de Idioma

- **English** - [README.md](README.md)
- **日本語** - [README.ja.md](README.ja.md)
- **中文** - [README.zh.md](README.zh.md)
- **Español** - [README.es.md](README.es.md)（actual）
- **हिंदी** - [README.hi.md](README.hi.md)
- **Português** - [README.pt.md](README.pt.md)

> **Nota**: Todas las versiones que no están en inglés son traducidas automáticamente. Para obtener la información más precisa, consulte la versión en inglés.

## Inicio Rápido

### Instalación

Descargue el binario más reciente para su plataforma desde la [página de lanzamientos](https://github.com/superdaigo/gsecutil/releases), o instale con Go:

```bash
go install github.com/superdaigo/gsecutil@latest
```

### Requisitos Previos

- [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install) instalado y autenticado
- Proyecto de Google Cloud con la API de Secret Manager habilitada

### Autenticación

```bash
# Autenticar con gcloud
gcloud auth login

# Establecer proyecto predeterminado
gcloud config set project YOUR_PROJECT_ID

# O establecer variable de entorno
export GSECUTIL_PROJECT=YOUR_PROJECT_ID
```

## Uso Básico

Cada proyecto suele tener su propio archivo de configuración que almacena el ID del proyecto, las convenciones de nomenclatura de secretos y los atributos de metadatos.

### 1. Crear un Archivo de Configuración

Ejecute la configuración interactiva para generar un archivo de configuración. Se le pedirá que ingrese el ID del proyecto de Google Cloud, el prefijo de nombre de secreto, los atributos de lista predeterminados y las credenciales de ejemplo opcionales. El archivo generado se guarda como `gsecutil.conf` en el directorio actual por defecto (use `--home` para guardar en `~/.config/gsecutil/gsecutil.conf`).

```bash
gsecutil config init
```

El archivo de configuración se busca en este orden:
1. Bandera `--config` (si se especifica)
2. Directorio actual: `gsecutil.conf`
3. Directorio de inicio: `~/.config/gsecutil/gsecutil.conf`

### 2. Administrar Secretos

```bash
# Crear un secreto
gsecutil create database-password

# Obtener la última versión
gsecutil get database-password

# Copiar al portapapeles
gsecutil get database-password --clipboard

# Listar todos los secretos
gsecutil list

# Actualizar un secreto
gsecutil update database-password

# Eliminar un secreto
gsecutil delete database-password
```

### Ejemplo de Configuración

```yaml
# ID del proyecto (opcional si se establece mediante entorno o gcloud)
project: "my-project-id"

# Prefijo de nombre de secreto para organización del equipo
prefix: "team-shared-"

# Atributos predeterminados para mostrar en el comando list
list:
  attributes:
    - title
    - owner
    - environment

# Metadatos de credenciales (los nombres son simples — el prefijo se añade automáticamente)
credentials:
  - name: "database-password"    # accede a "team-shared-database-password"
    title: "Production Database Password"
    environment: "production"
    owner: "backend-team"
```

> **El prefijo es transparente:** Cuando se configura un prefijo, siempre use nombres simples en comandos, configuración y archivos CSV. El prefijo se añade y elimina automáticamente.

Para opciones de configuración detalladas, consulte [docs/configuration.md](docs/configuration.md).

## Documentación

- **[Guía de Configuración](docs/configuration.md)** - Opciones de configuración detalladas y ejemplos
- **[Referencia de Comandos](docs/commands.md)** - Documentación completa de comandos
- **[Configuración de Registros de Auditoría](docs/audit-logging.md)** - Habilitar y usar registros de auditoría
- **[Guía de Solución de Problemas](docs/troubleshooting.md)** - Problemas comunes y soluciones
- **[Instrucciones de Compilación](BUILD.md)** - Compilar desde el código fuente
- **[Guía de Desarrollo](WARP.md)** - Desarrollo con WARP AI

## Licencia

Este proyecto está licenciado bajo la Licencia MIT - consulte el archivo LICENSE para más detalles.

## Relacionado

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Documentación de Secret Manager](https://cloud.google.com/secret-manager/docs)
