# gsecutil - Utilidad de Google Secret Manager

🚀 Un contenedor de línea de comandos simplificado para Google Secret Manager con soporte de archivos de configuración y funciones amigables para equipos.

## 🌍 Versiones de idioma

- **English** - [README.md](README.md)
- **日本語** - [README.ja.md](README.ja.md)
- **中文** - [README.zh.md](README.zh.md)
- **Español** - [README.es.md](README.es.md) (actual)
- **हिंदी** - [README.hi.md](README.hi.md)
- **Português** - [README.pt.md](README.pt.md)

> **Nota**: Todas las versiones que no están en inglés son traducidas automáticamente. Para obtener la información más precisa, consulte la versión en inglés.

## Inicio rápido

### Instalación

Descargue el binario más reciente para su plataforma desde la [página de lanzamientos](https://github.com/superdaigo/gsecutil/releases):

```bash
# macOS Apple Silicon
curl -L https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-darwin-arm64 -o gsecutil
chmod +x gsecutil
sudo mv gsecutil /usr/local/bin/

# macOS Intel
curl -L https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-darwin-amd64 -o gsecutil
chmod +x gsecutil
sudo mv gsecutil /usr/local/bin/

# Linux
curl -L https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-linux-amd64 -o gsecutil
chmod +x gsecutil
sudo mv gsecutil /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-windows-amd64.exe" -OutFile "gsecutil.exe"
# Move to a directory in your PATH, e.g., C:\Windows\System32
Move-Item gsecutil.exe C:\Windows\System32\gsecutil.exe
```

O instalar con Go:
```bash
go install github.com/superdaigo/gsecutil@latest
```

### Requisitos previos

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

## Uso básico

### Crear un secreto
```bash
# Entrada interactiva
gsecutil create database-password

# Desde línea de comandos
gsecutil create api-key -d "sk-1234567890"

# Desde archivo
gsecutil create config --data-file ./config.json
```

### Obtener un secreto
```bash
# Obtener última versión
gsecutil get database-password

# Copiar al portapapeles
gsecutil get api-key --clipboard

# Obtener versión específica
gsecutil get api-key --version 2
```

### Listar secretos
```bash
# Listar todos los secretos
gsecutil list

# Filtrar por etiqueta
gsecutil list --filter "labels.env=prod"
```

### Actualizar un secreto
```bash
# Entrada interactiva
gsecutil update database-password

# Desde línea de comandos
gsecutil update api-key -d "new-secret-value"
```

### Eliminar un secreto
```bash
gsecutil delete old-secret
```

## Configuración

gsecutil admite archivos de configuración para ajustes específicos del proyecto. Los archivos de configuración se buscan en este orden:

1. Bandera `--config` (si se especifica)
2. Directorio actual: `gsecutil.conf`
3. Directorio de inicio: `~/.config/gsecutil/gsecutil.conf`

### Ejemplo de configuración

```yaml
# ID del proyecto (opcional si se establece mediante variable de entorno o gcloud)
project: "my-project-id"

# Prefijo de nombre de secreto para organización de equipo
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

> **El prefijo es transparente:** Cuando se configura un prefijo, siempre se usan nombres simples (sin prefijo) en comandos, configuración y archivos CSV. El prefijo se añade y elimina automáticamente.

### Inicio rápido

```bash
# Generar configuración interactivamente
gsecutil config init

# O crear una configuración específica del proyecto
echo 'project: "my-project-123"' > gsecutil.conf
```

Para opciones de configuración detalladas, consulte [docs/configuration.md](docs/configuration.md).

## Características principales

- ✅ **Operaciones CRUD simples** - Comandos intuitivos para gestionar secretos
- ✅ **Integración con portapapeles** - Copiar secretos directamente al portapapeles
- ✅ **Gestión de versiones** - Acceder a versiones específicas y gestionar el ciclo de vida de versiones
- ✅ **Soporte de archivos de configuración** - Metadatos y organización amigables para equipos
- ✅ **Gestión de acceso** - Gestión básica de políticas IAM
- ✅ **Registros de auditoría** - Ver quién accedió a los secretos y cuándo
- ✅ **Múltiples métodos de entrada** - Interactivo, en línea o basado en archivos
- ✅ **Multiplataforma** - Linux, macOS, Windows (amd64/arm64)

## Documentación

- **[Guía de configuración](docs/configuration.md)** - Opciones de configuración detalladas y ejemplos
- **[Referencia de comandos](docs/commands.md)** - Documentación completa de comandos
- **[Configuración de registros de auditoría](docs/audit-logging.md)** - Habilitar y usar registros de auditoría
- **[Guía de solución de problemas](docs/troubleshooting.md)** - Problemas comunes y soluciones
- **[Instrucciones de compilación](BUILD.md)** - Compilar desde el código fuente
- **[Guía de desarrollo](WARP.md)** - Desarrollo con WARP AI

## Comandos comunes

```bash
# Mostrar detalles del secreto
gsecutil describe my-secret

# Mostrar historial de versiones
gsecutil describe my-secret --show-versions

# Ver registros de auditoría
gsecutil auditlog my-secret

# Gestionar acceso
gsecutil access list my-secret
gsecutil access grant my-secret --principal user:alice@example.com

# Validar configuración
gsecutil config validate

# Mostrar configuración
gsecutil config show
```

## Licencia

Este proyecto está licenciado bajo la Licencia MIT; consulte el archivo LICENSE para más detalles.

## Relacionado

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Documentación de Secret Manager](https://cloud.google.com/secret-manager/docs)
