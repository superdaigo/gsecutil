# gsecutil - Utilitário do Google Secret Manager

> **Nota sobre a tradução**: Este arquivo README foi traduzido automaticamente. Para obter as informações mais atualizadas e precisas, consulte a versão em inglês [README.md](README.md).
>
> **🆕 Nova funcionalidade**: v1.1.1 adiciona gerenciamento automático de versões para permanecer dentro do nível gratuito do Google Cloud (6 versões ativas). Consulte o README em inglês para mais detalhes.

🚀 **v1.1.0** - Um wrapper simplificado de linha de comando para Google Secret Manager com suporte a arquivos de configuração. `gsecutil` oferece comandos convenientes para operações comuns de segredos, facilitando que pequenas equipes gerenciem senhas e credenciais usando o Secret Manager do Google Cloud sem precisar de ferramentas dedicadas de gerenciamento de senhas.

**NOVO na v1.1.0**: Suporte a arquivos de configuração YAML, funcionalidade de prefixo e comandos aprimorados de lista e descrição com metadados personalizados da equipe.

## ✨ Recursos

### 🔐 **Gerenciamento Completo de Segredos**
- **Operações CRUD**: Criar, ler, atualizar, excluir segredos com comandos simplificados
- **Gerenciamento de versões**: Acessar qualquer versão, visualizar histórico de versões e metadados
- **Suporte multiplataforma** (Linux, macOS, Windows com suporte ARM64)
- **Integração com área de transferência** - copiar valores de segredos diretamente para a área de transferência
- **Entrada interativa e de arquivo** - prompts seguros ou carregamento de segredos baseado em arquivo

### 🛡️ **Gerenciamento Avançado de Acesso**
*(Introduzido na v1.0.0)*
- **Análise completa de políticas IAM** - ver quem tem acesso a segredos em qualquer nível
- **Verificação de permissões multinível** - análise de acesso a nível de segredo e projeto
- **Consciência de condições IAM** - suporte completo para políticas de acesso condicional com expressões CEL
- **Gerenciamento de principais** - conceder/revogar acesso para usuários, grupos e contas de serviço
- **Análise de todo o projeto** - identificar papéis de Editor/Proprietário que fornecem acesso ao Secret Manager

### 📊 **Auditoria e Conformidade**
- **Log de auditoria abrangente** - rastrear quem acessou segredos, quando e quais operações
- **Filtragem baseada em principais** - ver todos os segredos acessíveis por usuários/grupos específicos
- **Filtragem flexível** - por segredo, principal, tipo de operação, intervalo de tempo
- **Avaliação de condições** - entender quando o acesso condicional se aplica

### 🎯 **Pronto para Produção**
- **API consistente** - nomenclatura unificada de parâmetros em todos os comandos
- **Recursos empresariais** - condições IAM, análise a nível de projeto, auditoria de conformidade
- **Tratamento robusto de erros** - tratamento elegante de permissões ausentes e problemas de rede
- **Saída flexível** - formatos JSON, YAML, tabela com formatação rica

## Pré-requisitos

- [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install) instalado e autenticado
- Projeto Google Cloud com a API do Secret Manager habilitada
- Permissões IAM apropriadas para operações do Secret Manager

## Instalação

### Binários Pré-construídos

Baixe a versão mais recente para sua plataforma da [página de releases](https://github.com/superdaigo/gsecutil/releases):

| Plataforma | Arquitetura | Download |
|----------|--------------|----------|
| Linux | x64 | `gsecutil-linux-amd64-v{version}` |
| Linux | ARM64 | `gsecutil-linux-arm64-v{version}` |
| macOS | Intel | `gsecutil-darwin-amd64-v{version}` |
| macOS | Apple Silicon | `gsecutil-darwin-arm64-v{version}` |
| Windows | x64 | `gsecutil-windows-amd64-v{version}.exe` |

**Após o download:** Renomeie o binário para uso consistente:

```bash
# Exemplo Linux/macOS:
mv gsecutil-linux-amd64-v1.1.0 gsecutil
chmod +x gsecutil

# Exemplo Windows (PowerShell/Command Prompt):
ren gsecutil-windows-amd64-v1.1.0.exe gsecutil.exe
```

Isso permite usar `gsecutil` de forma consistente, independentemente da versão.

### Instalar com Go

```bash
go install github.com/superdaigo/gsecutil@latest
```

### Construir a partir do Código Fonte

Para instruções de construção abrangentes, consulte [BUILD.md](BUILD.md).

**Construção rápida:**
```bash
git clone https://github.com/superdaigo/gsecutil.git
cd gsecutil

# Construir para a plataforma atual
make build
# OU
./build.sh          # Linux/macOS
.\\build.ps1         # Windows

# Construir para todas as plataformas
make build-all
# OU
./build.sh all      # Linux/macOS
.\\build.ps1 all     # Windows
```

## Uso

### Opções Globais

- `-p, --project`: ID do projeto Google Cloud (também pode ser definido via variável de ambiente `GOOGLE_CLOUD_PROJECT`)

### Comandos

#### Get Secret (Obter Segredo)

Recupera um valor de segredo do Google Secret Manager. Por padrão, obtém a versão mais recente, mas você pode especificar qualquer versão:

```bash
# Obter a versão mais recente de um segredo
gsecutil get my-secret

# Obter versão específica (útil para rollbacks, debugging ou acesso a valores históricos)
gsecutil get my-secret --version 1
gsecutil get my-secret -v 3

# Obter segredo e copiar para área de transferência
gsecutil get my-secret --clipboard

# Obter versão específica com área de transferência
gsecutil get my-secret --version 2 --clipboard

# Obter segredo com metadados de versão (versão, tempo de criação, estado)
gsecutil get my-secret --show-metadata
gsecutil get my-secret -v 1 --show-metadata    # Versão antiga com metadados

# Obter segredo de projeto específico
gsecutil get my-secret --project my-gcp-project
```

**Suporte a Versões:**
- 🔄 **Versão mais recente**: Comportamento padrão quando `--version` não é especificado
- 📅 **Versões históricas**: Acesso a qualquer versão anterior por número (ex., `--version 1`, `--version 2`)
- 🔍 **Metadados de versão**: Use `--show-metadata` para ver detalhes da versão (tempo de criação, estado, ETag)
- 📋 **Suporte à área de transferência**: Funciona com qualquer versão usando `--clipboard`

## Configuração

### Variáveis de Ambiente

- `GOOGLE_CLOUD_PROJECT`: ID do projeto padrão (sobrescrito pela flag `--project`)

### Autenticação

`gsecutil` usa a mesma autenticação que o `gcloud`. Certifique-se de estar autenticado:

```bash
# Autenticar com gcloud
gcloud auth login

# Definir projeto padrão
gcloud config set project YOUR_PROJECT_ID

# Para contas de serviço (em CI/CD)
gcloud auth activate-service-account --key-file=service-account.json
```

### Autocompletar de Shell

`gsecutil` suporta autocompletamento de shell para bash, zsh, fish e PowerShell. Isso habilita o completamento por tab para comandos, flags e opções, tornando o CLI mais fácil de usar.

#### Instruções de Configuração

**Bash:**
```bash
# Temporário (apenas sessão atual)
source <(gsecutil completion bash)

# Instalação permanente (requer pacote bash-completion)
# Sistema completo (requer sudo)
sudo gsecutil completion bash > /etc/bash_completion.d/gsecutil

# Instalação local do usuário
gsecutil completion bash > ~/.local/share/bash-completion/completions/gsecutil

# Ou adicionar ao ~/.bashrc para carregamento automático
echo 'source <(gsecutil completion bash)' >> ~/.bashrc
```

**Zsh:**
```bash
# Temporário (apenas sessão atual)
source <(gsecutil completion zsh)

# Instalação permanente
gsecutil completion zsh > "${fpath[1]}/_gsecutil"

# Ou adicionar ao ~/.zshrc para carregamento automático
echo 'source <(gsecutil completion zsh)' >> ~/.zshrc
```

**Fish:**
```bash
# Temporário (apenas sessão atual)
gsecutil completion fish | source

# Instalação permanente
gsecutil completion fish > ~/.config/fish/completions/gsecutil.fish
```

**PowerShell:**
```powershell
# Adicionar ao perfil do PowerShell
gsecutil completion powershell | Out-String | Invoke-Expression

# Ou salvar no perfil para carregamento automático
gsecutil completion powershell >> $PROFILE
```

#### Recursos

Uma vez instalado, o autocompletamento de shell fornece:
- **Completamento de comandos**: Tab para completar subcomandos do `gsecutil` (`get`, `create`, `list`, etc.)
- **Completamento de flags**: Tab para completar flags como `--project`, `--version`, `--clipboard`
- **Sugestões inteligentes**: Completamentos conscientes do contexto baseados no comando atual
- **Texto de ajuda**: Descrições breves para comandos e flags (onde suportado)

#### Exemplo de Uso

```bash
# Digite e pressione Tab para ver comandos disponíveis
gsecutil <Tab>
# Mostra: access, auditlog, completion, create, delete, describe, get, help, list, update

# Digite comando parcial e pressione Tab para completar
gsecutil des<Tab>
# Completa para: gsecutil describe

# Completamento por tab também funciona para flags
gsecutil get my-secret --<Tab>
# Mostra: --clipboard, --project, --show-metadata, --version
```

**Nota**: Você pode precisar reiniciar seu shell ou fazer source do seu arquivo de configuração de shell para que o completamento tenha efeito.

## Segurança e Melhores Práticas

### Recursos de Segurança

- **Sem armazenamento persistente**: Valores de segredos nunca são registrados ou armazenados pelo `gsecutil`
- **Entrada segura**: Prompts interativos usam entrada de senha oculta
- **Área de transferência nativa do SO**: Operações da área de transferência usam APIs nativas seguras do SO
- **Delegação gcloud**: Todas as operações são delegadas para o CLI `gcloud` autenticado

### Melhores Práticas

- **Use `--force` com cuidado**: Sempre revise antes de usar `--force` em ambientes automatizados
- **Variáveis de ambiente**: Defina `GOOGLE_CLOUD_PROJECT` para evitar flags repetitivas `--project`
- **Controle de versão**: Use versões específicas de segredos em produção (`--version N`)
- **Auditoria regular**: Monitore acesso a segredos com `gsecutil auditlog secret-name`
- **Rotação de segredos**: Rotação regular de segredos usando `gsecutil update`

## Solução de Problemas

### Problemas Comuns

1. **"gcloud command not found"**
   - Certifique-se de que o Google Cloud SDK está instalado e `gcloud` está no seu PATH

2. **Erros de autenticação**
   - Execute `gcloud auth login` para autenticar
   - Verifique o acesso ao projeto: `gcloud config get-value project`

3. **Erros de permissão negada**
   - Certifique-se de que sua conta possui os papéis IAM necessários:
     - `roles/secretmanager.admin` (para todas as operações)
     - `roles/secretmanager.secretAccessor` (para operações de leitura)
     - `roles/secretmanager.secretVersionManager` (para operações de criação/atualização)

4. **Área de transferência não funciona**
   - Certifique-se de ter um ambiente gráfico (para Linux)
   - Em servidores sem interface gráfica, operações da área de transferência podem falhar graciosamente

### Modo de Debug

Adicione saída detalhada aos comandos gcloud definindo:

```bash
export CLOUDSDK_CORE_VERBOSITY=debug
```

## Documentação

- **[BUILD.md](BUILD.md)** - Instruções de construção abrangentes para todas as plataformas
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Diretrizes de contribuição e fluxo de trabalho de desenvolvimento
- **[WARP.md](WARP.md)** - Orientação de desenvolvimento para integração com terminal WARP AI
- **README.md** - Este arquivo, uso e visão geral

## Contribuição

Contribuições são bem-vindas! Consulte [CONTRIBUTING.md](CONTRIBUTING.md) para diretrizes detalhadas sobre como contribuir para este projeto, incluindo instruções de configuração para ambiente de desenvolvimento e hooks de pré-commit.

## Licença

Este projeto está licenciado sob a Licença MIT - consulte o arquivo LICENSE para detalhes.

## Projetos Relacionados

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager Documentation](https://cloud.google.com/secret-manager/docs)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
