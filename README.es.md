# gsecutil - Utilidad de Google Secret Manager

> **Nota sobre la traducción**: Este archivo README ha sido traducido automáticamente. Para obtener la información más actualizada y precisa, consulte la versión en inglés [README.md](README.md).

🚀 **v1.1.0** - Un envoltorio simplificado de línea de comandos para Google Secret Manager con soporte para archivos de configuración. `gsecutil` proporciona comandos convenientes para operaciones comunes de secretos, facilitando que equipos pequeños gestionen contraseñas y credenciales usando Google Cloud Secret Manager sin necesidad de herramientas dedicadas de gestión de contraseñas.

**NUEVO en v1.1.0**: Soporte para archivos de configuración YAML, funcionalidad de prefijos, y comandos mejorados de lista y descripción con metadatos personalizados del equipo.

## ✨ Características

### 🔐 **Gestión Completa de Secretos**
- **Operaciones CRUD**: Crear, leer, actualizar, eliminar secretos con comandos simplificados
- **Gestión de versiones**: Acceso a cualquier versión, visualizar historial de versiones y metadatos
- **Soporte multiplataforma** (Linux, macOS, Windows con soporte ARM64)
- **Integración de portapapeles** - copiar valores de secretos directamente al portapapeles
- **Entrada interactiva y de archivos** - solicitudes seguras o carga de secretos basada en archivos

### 🛡️ **Gestión Avanzada de Acceso**
*(Introducido en v1.0.0)*
- **Análisis completo de políticas IAM** - ver quién tiene acceso a secretos en cualquier nivel
- **Verificación de permisos multinivel** - análisis de acceso a nivel de secreto y proyecto
- **Reconocimiento de condiciones IAM** - soporte completo para políticas de acceso condicional con expresiones CEL
- **Gestión de principales** - otorgar/revocar acceso para usuarios, grupos y cuentas de servicio
- **Análisis de todo el proyecto** - identificar roles de Editor/Propietario que proporcionan acceso a Secret Manager

### 📊 **Auditoría y Cumplimiento**
- **Registro de auditoría integral** - rastrear quién accedió a secretos, cuándo y qué operaciones
- **Filtrado basado en principales** - ver todos los secretos accesibles por usuarios/grupos específicos
- **Filtrado flexible** - por secreto, principal, tipo de operación, rango de tiempo
- **Evaluación de condiciones** - entender cuándo se aplica el acceso condicional

### 🎯 **Listo para Producción**
- **API consistente** - nomenclatura unificada de parámetros en todos los comandos
- **Características empresariales** - condiciones IAM, análisis a nivel de proyecto, auditoría de cumplimiento
- **Manejo robusto de errores** - manejo elegante de permisos faltantes y problemas de red
- **Salida flexible** - formatos JSON, YAML, tabla con formateo enriquecido

## Prerrequisitos

- [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install) instalado y autenticado
- Proyecto de Google Cloud con la API de Secret Manager habilitada
- Permisos IAM apropiados para operaciones de Secret Manager

## Instalación

### Binarios Precompilados

Descarga la última versión para tu plataforma desde la [página de versiones](https://github.com/superdaigo/gsecutil/releases):

| Plataforma | Arquitectura | Descarga |
|----------|--------------|----------|
| Linux | x64 | `gsecutil-linux-amd64-v{version}` |
| Linux | ARM64 | `gsecutil-linux-arm64-v{version}` |
| macOS | Intel | `gsecutil-darwin-amd64-v{version}` |
| macOS | Apple Silicon | `gsecutil-darwin-arm64-v{version}` |
| Windows | x64 | `gsecutil-windows-amd64-v{version}.exe` |

**Después de la descarga:** Renombra el binario para uso consistente:

```bash
# Ejemplo Linux/macOS:
mv gsecutil-linux-amd64-v1.1.0 gsecutil
chmod +x gsecutil

# Ejemplo Windows (PowerShell/Command Prompt):
ren gsecutil-windows-amd64-v1.1.0.exe gsecutil.exe
```

Esto te permite usar `gsecutil` de manera consistente independientemente de la versión.

### Instalar con Go

```bash
go install github.com/superdaigo/gsecutil@latest
```

### Compilar desde el Código Fuente

Para instrucciones de compilación completas, consulta [BUILD.md](BUILD.md).

**Compilación rápida:**
```bash
git clone https://github.com/superdaigo/gsecutil.git
cd gsecutil

# Compilar para la plataforma actual
make build
# O
./build.sh          # Linux/macOS
.\\build.ps1         # Windows

# Compilar para todas las plataformas
make build-all
# O
./build.sh all      # Linux/macOS
.\\build.ps1 all     # Windows
```

## Uso

### Opciones Globales

- `-p, --project`: ID del proyecto de Google Cloud (también se puede establecer mediante la variable de entorno `GOOGLE_CLOUD_PROJECT`)

### Comandos

#### Get Secret (Obtener Secreto)

Recupera un valor de secreto de Google Secret Manager. Por defecto, obtiene la versión más reciente, pero puedes especificar cualquier versión:

```bash
# Obtener la versión más reciente de un secreto
gsecutil get my-secret

# Obtener versión específica (útil para rollbacks, depuración o acceso a valores históricos)
gsecutil get my-secret --version 1
gsecutil get my-secret -v 3

# Obtener secreto y copiar al portapapeles
gsecutil get my-secret --clipboard

# Obtener versión específica con portapapeles
gsecutil get my-secret --version 2 --clipboard

# Obtener secreto con metadatos de versión (versión, tiempo de creación, estado)
gsecutil get my-secret --show-metadata
gsecutil get my-secret -v 1 --show-metadata    # Versión anterior con metadatos

# Obtener secreto de proyecto específico
gsecutil get my-secret --project my-gcp-project
```

**Soporte de Versiones:**
- 🔄 **Versión más reciente**: Comportamiento predeterminado cuando no se especifica `--version`
- 📅 **Versiones históricas**: Acceso a cualquier versión anterior por número (ej., `--version 1`, `--version 2`)
- 🔍 **Metadatos de versión**: Usar `--show-metadata` para ver detalles de versión (tiempo de creación, estado, ETag)
- 📋 **Soporte de portapapeles**: Funciona con cualquier versión usando `--clipboard`

## Configuración

### Variables de Entorno

- `GOOGLE_CLOUD_PROJECT`: ID del proyecto predeterminado (anulado por la bandera `--project`)

### Autenticación

`gsecutil` usa la misma autenticación que `gcloud`. Asegúrate de estar autenticado:

```bash
# Autenticar con gcloud
gcloud auth login

# Establecer proyecto predeterminado
gcloud config set project YOUR_PROJECT_ID

# Para cuentas de servicio (en CI/CD)
gcloud auth activate-service-account --key-file=service-account.json
```

## Seguridad y Mejores Prácticas

### Características de Seguridad

- **Sin almacenamiento persistente**: Los valores de secretos nunca son registrados o almacenados por `gsecutil`
- **Entrada segura**: Las solicitudes interactivas usan entrada de contraseña oculta
- **Portapapeles nativo del SO**: Las operaciones de portapapeles usan APIs nativas seguras del SO
- **Delegación gcloud**: Todas las operaciones se delegan al CLI `gcloud` autenticado

### Mejores Prácticas

- **Usar `--force` con cuidado**: Siempre revisar antes de usar `--force` en entornos automatizados
- **Variables de entorno**: Establecer `GOOGLE_CLOUD_PROJECT` para evitar banderas repetitivas `--project`
- **Control de versiones**: Usar versiones específicas de secretos en producción (`--version N`)
- **Auditar regularmente**: Monitorear acceso a secretos con `gsecutil auditlog secret-name`
- **Rotación de secretos**: Rotación regular de secretos usando `gsecutil update`

## Solución de Problemas

### Problemas Comunes

1. **"gcloud command not found"**
   - Asegurar que Google Cloud SDK esté instalado y `gcloud` esté en tu PATH

2. **Errores de autenticación**
   - Ejecutar `gcloud auth login` para autenticar
   - Verificar acceso al proyecto: `gcloud config get-value project`

3. **Errores de permisos denegados**
   - Asegurar que tu cuenta tenga los roles IAM necesarios:
     - `roles/secretmanager.admin` (para todas las operaciones)
     - `roles/secretmanager.secretAccessor` (para operaciones de lectura)
     - `roles/secretmanager.secretVersionManager` (para operaciones de creación/actualización)

4. **Portapapeles no funciona**
   - Asegurar que tengas un entorno gráfico (para Linux)
   - En servidores sin cabeza, las operaciones de portapapeles pueden fallar elegantemente

### Modo de Depuración

Añadir salida detallada a comandos gcloud estableciendo:

```bash
export CLOUDSDK_CORE_VERBOSITY=debug
```

## Documentación

- **[BUILD.md](BUILD.md)** - Instrucciones de compilación completas para todas las plataformas
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Pautas de contribución y flujo de trabajo de desarrollo
- **[WARP.md](WARP.md)** - Guía de desarrollo para integración con terminal WARP AI
- **README.md** - Este archivo, uso y descripción general

## Contribución

¡Las contribuciones son bienvenidas! Consulta [CONTRIBUTING.md](CONTRIBUTING.md) para pautas detalladas sobre cómo contribuir a este proyecto, incluyendo instrucciones de configuración para el entorno de desarrollo y ganchos de pre-commit.

## Licencia

Este proyecto está licenciado bajo la Licencia MIT - consulta el archivo LICENSE para más detalles.

## Proyectos Relacionados

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager Documentation](https://cloud.google.com/secret-manager/docs)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
