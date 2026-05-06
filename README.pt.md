# gsecutil - Utilitário do Google Secret Manager

Um wrapper de linha de comando simplificado para o Google Secret Manager que funciona como um gerenciador de senhas por projeto. Armazene, recupere e gerencie segredos com comandos intuitivos, integração com a área de transferência, controle de versão, arquivos de configuração amigáveis para equipes e registros de auditoria.

## 🌍 Versões de Idioma

- **English** - [README.md](README.md)
- **日本語** - [README.ja.md](README.ja.md)
- **中文** - [README.zh.md](README.zh.md)
- **Español** - [README.es.md](README.es.md)
- **हिंदी** - [README.hi.md](README.hi.md)
- **Português** - [README.pt.md](README.pt.md)（atual）

> **Nota**: Todas as versões que não estão em inglês são traduzidas por máquina. Para obter as informações mais precisas, consulte a versão em inglês.

## Início Rápido

### Instalação

Baixe o binário mais recente para sua plataforma na [página de releases](https://github.com/superdaigo/gsecutil/releases), ou instale com Go:

```bash
go install github.com/superdaigo/gsecutil@latest
```

### Pré-requisitos

- [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install) instalado e autenticado
- Projeto do Google Cloud com a API do Secret Manager habilitada

### Autenticação

```bash
# Autentique com o gcloud
gcloud auth login

# Defina o projeto padrão
gcloud config set project YOUR_PROJECT_ID

# Ou defina a variável de ambiente
export GSECUTIL_PROJECT=YOUR_PROJECT_ID
```

## Uso Básico

Cada projeto normalmente tem seu próprio arquivo de configuração que armazena o ID do projeto, as convenções de nomenclatura de segredos e os atributos de metadados.

### 1. Criar um Arquivo de Configuração

Execute a configuração interativa para gerar um arquivo de configuração. Você será solicitado a informar o ID do projeto do Google Cloud, o prefixo do nome do segredo, os atributos padrão da lista e as credenciais de exemplo opcionais. O arquivo gerado é salvo como `gsecutil.conf` no diretório atual por padrão (use `--home` para salvar em `~/.config/gsecutil/gsecutil.conf`).

```bash
gsecutil config init
```

O arquivo de configuração é pesquisado nesta ordem:
1. Flag `--config` (se especificada)
2. Diretório atual: `gsecutil.conf`
3. Diretório home: `~/.config/gsecutil/gsecutil.conf`

### 2. Gerenciar Segredos

```bash
# Criar um segredo
gsecutil create database-password

# Obter a versão mais recente
gsecutil get database-password

# Copiar para a área de transferência
gsecutil get database-password --clipboard

# Listar todos os segredos
gsecutil list

# Atualizar um segredo
gsecutil update database-password

# Excluir um segredo
gsecutil delete database-password
```

### Exemplo de Configuração

```yaml
# ID do projeto (opcional se definido via ambiente ou gcloud)
project: "my-project-id"

# Prefixo de nome de segredo para organização da equipe
prefix: "team-shared-"

# Atributos padrão para exibir no comando list
list:
  attributes:
    - title
    - owner
    - environment

# Metadados de credenciais (os nomes são simples — o prefixo é adicionado automaticamente)
credentials:
  - name: "database-password"    # acessa "team-shared-database-password"
    title: "Production Database Password"
    environment: "production"
    owner: "backend-team"
```

> **O prefixo é transparente:** Quando um prefixo é configurado, sempre use nomes simples em comandos, configuração e arquivos CSV. O prefixo é adicionado e removido automaticamente.

Para opções detalhadas de configuração, consulte [docs/configuration.md](docs/configuration.md).

## Documentação

- **[Guia de Configuração](docs/configuration.md)** - Opções de configuração detalhadas e exemplos
- **[Referência de Comandos](docs/commands.md)** - Documentação completa de comandos
- **[Configuração de Log de Auditoria](docs/audit-logging.md)** - Habilite e use logs de auditoria
- **[Guia de Solução de Problemas](docs/troubleshooting.md)** - Problemas comuns e soluções
- **[Instruções de Build](BUILD.md)** - Compilar a partir do código-fonte
- **[Guia de Desenvolvimento](WARP.md)** - Desenvolvimento com WARP AI

## Licença

Este projeto está licenciado sob a Licença MIT - consulte o arquivo LICENSE para obter detalhes.

## Links Relacionados

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Documentação do Secret Manager](https://cloud.google.com/secret-manager/docs)
