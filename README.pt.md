# gsecutil - Utilitário do Google Secret Manager

🚀 Um wrapper de linha de comando simplificado para o Google Secret Manager com suporte a arquivos de configuração e recursos amigáveis para equipes.

## 🌍 Versões de idioma

- **English** - [README.md](README.md)
- **日本語** - [README.ja.md](README.ja.md)
- **中文** - [README.zh.md](README.zh.md)
- **Español** - [README.es.md](README.es.md)
- **हिंदी** - [README.hi.md](README.hi.md)
- **Português** - [README.pt.md](README.pt.md) (atual)

> **Nota**: Todas as versões que não estão em inglês são traduzidas por máquina. Para obter as informações mais precisas, consulte a versão em inglês.

## Início Rápido

### Instalação

Baixe o binário mais recente para sua plataforma na [página de lançamentos](https://github.com/superdaigo/gsecutil/releases):

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

Ou instale com Go:
```bash
go install github.com/superdaigo/gsecutil@latest
```

### Pré-requisitos

- [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install) instalado e autenticado
- Projeto do Google Cloud com a API do Secret Manager ativada

### Autenticação

```bash
# Autenticar com gcloud
gcloud auth login

# Definir projeto padrão
gcloud config set project YOUR_PROJECT_ID

# Ou definir variável de ambiente
export GSECUTIL_PROJECT=YOUR_PROJECT_ID
```

## Uso Básico

### Criar um Segredo
```bash
# Entrada interativa
gsecutil create database-password

# Da linha de comando
gsecutil create api-key -d "sk-1234567890"

# De um arquivo
gsecutil create config --data-file ./config.json
```

### Obter um Segredo
```bash
# Obter versão mais recente
gsecutil get database-password

# Copiar para a área de transferência
gsecutil get api-key --clipboard

# Obter versão específica
gsecutil get api-key --version 2
```

### Listar Segredos
```bash
# Listar todos os segredos
gsecutil list

# Filtrar por rótulo
gsecutil list --filter "labels.env=prod"
```

### Atualizar um Segredo
```bash
# Entrada interativa
gsecutil update database-password

# Da linha de comando
gsecutil update api-key -d "new-secret-value"
```

### Excluir um Segredo
```bash
gsecutil delete old-secret
```

## Configuração

Crie um arquivo de configuração em `~/.config/gsecutil/gsecutil.conf`:

```yaml
# ID do projeto (opcional se definido via variável de ambiente ou gcloud)
project: "my-project-id"

# Prefixo de nome de segredo para organização de equipe
prefix: "team-shared-"

# Atributos padrão para exibir no comando list
list:
  attributes:
    - title
    - owner
    - environment

# Metadados de credenciais
credentials:
  - name: "database-password"
    title: "Production Database Password"
    environment: "production"
    owner: "backend-team"
```

Gerar configuração interativamente:
```bash
gsecutil config init
```

Para opções de configuração detalhadas, consulte [docs/configuration.md](docs/configuration.md).

## Recursos Principais

- ✅ **Operações CRUD Simples** - Comandos intuitivos para gerenciar segredos
- ✅ **Integração com Área de Transferência** - Copiar segredos diretamente para a área de transferência
- ✅ **Gerenciamento de Versões** - Acessar versões específicas e gerenciar o ciclo de vida das versões
- ✅ **Suporte a Arquivos de Configuração** - Metadados e organização amigáveis para equipes
- ✅ **Gerenciamento de Acesso** - Gerenciamento básico de políticas IAM
- ✅ **Logs de Auditoria** - Ver quem acessou os segredos e quando
- ✅ **Múltiplos Métodos de Entrada** - Interativo, inline ou baseado em arquivo
- ✅ **Multiplataforma** - Linux, macOS, Windows (amd64/arm64)

## Documentação

- **[Guia de Configuração](docs/configuration.md)** - Opções de configuração detalhadas e exemplos
- **[Referência de Comandos](docs/commands.md)** - Documentação completa de comandos
- **[Configuração de Logs de Auditoria](docs/audit-logging.md)** - Ativar e usar logs de auditoria
- **[Guia de Solução de Problemas](docs/troubleshooting.md)** - Problemas comuns e soluções
- **[Instruções de Compilação](BUILD.md)** - Compilar do código-fonte
- **[Guia de Desenvolvimento](WARP.md)** - Desenvolvimento com WARP AI

## Comandos Comuns

```bash
# Mostrar detalhes do segredo
gsecutil describe my-secret

# Mostrar histórico de versões
gsecutil describe my-secret --show-versions

# Ver logs de auditoria
gsecutil auditlog my-secret

# Gerenciar acesso
gsecutil access list my-secret
gsecutil access grant my-secret --principal user:alice@example.com

# Validar configuração
gsecutil config validate

# Mostrar configuração
gsecutil config show
```

## Licença

Este projeto está licenciado sob a Licença MIT - consulte o arquivo LICENSE para obter detalhes.

## Relacionado

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Documentação do Secret Manager](https://cloud.google.com/secret-manager/docs)
