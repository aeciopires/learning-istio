# Plano executado 2 — mover o código do CLI para `agents/upgrade/app/`

> Continuação de [plan1.md](plan1.md): registra a reorganização de
> diretórios pedida depois da entrega inicial do CLI.

## Pedido

Mover os arquivos específicos do CLI para dentro de uma pasta `app`
dentro de `agents/upgrade`, ajustando todas as referências para manter
tudo funcionando.

## O que foi considerado "específico do CLI" (movido para `app/`)

`main.go`, `cmd/`, `internal/`, `go.mod`, `go.sum`, `VERSION` (precisa
ficar ao lado de `main.go` por causa do `//go:embed VERSION`),
`Dockerfile`, `.dockerignore` (contexto de build do Docker vira `app/`),
`testdata/`, `examples/`, `CHANGELOG.md` (changelog é do CLI
especificamente).

## O que ficou fora de `app/` (não é específico do CLI)

- `Makefile` - continua na raiz de `agents/upgrade/`, é o ponto de
  entrada único tanto para o CLI quanto para a automação de cluster.
- `README.md` - documenta as duas partes (CLI + automação), não só o CLI.
- `scripts/` - automação de cluster, por definição não é o CLI.
- `configs/versions.env` - é a ponte entre o CLI (`docs update --apply`
  escreve) e os scripts (`scripts/lib/common.sh` lê via `source`); mover
  para dentro de `app/` misturaria as duas responsabilidades de novo.
- `plans/` - documentação de planejamento do diretório como um todo.

## Decisões técnicas

- **Módulo Go**: o `module` em `go.mod` continua
  `github.com/aeciopires/learning-istio/agents/upgrade` (não virou
  `.../agents/upgrade/app`). O caminho do módulo Go não precisa refletir
  a posição no monorepo, e trocá-lo exigiria editar o import de todo
  arquivo `.go` do projeto sem nenhum ganho funcional (o módulo nunca é
  importado de fora dele mesmo).
- **`Makefile` passou a usar `go build -C app`, `go test -C app`,
  `go vet -C app`, `go run -C app .`, `go install -C app .`** (flag `-C`
  do próprio `go`, disponível desde Go 1.20) em vez de `cd app && ...`,
  para manter os targets como uma linha só. Uma exceção: `go tool cover`
  não aceita `-C` depois do subcomando, então virou `go -C app tool cover
  -func=coverage.out` (o `-C` é um flag do `go` em si, não do
  subcomando `tool`).
- **`golangci-lint`** não tem uma flag equivalente confiável para todos os
  casos, então o target `lint` usa `(cd $(APP_DIR) && golangci-lint run ./...)`
  em subshell.
- **Binário compilado continua em `agents/upgrade/bin/`** (não
  `agents/upgrade/app/bin/`), para não quebrar a instrução
  `./bin/istio-upgrade` já documentada no README - o `-o` do `go build -C
  app` usa um caminho relativo `../bin/istio-upgrade`, que resolve para
  fora de `app/`.
- **Docker**: o contexto de build passou a ser `app/` (onde `go.mod` e o
  `Dockerfile` agora vivem), via
  `docker build -f app/Dockerfile -t $(DOCKER_IMAGE) app`.
- **`.gitignore`** trocou os padrões `/bin/`, `/coverage.out`, `/.cache/`
  (ancorados na raiz de `agents/upgrade/`) para `/bin/`, `coverage.out`,
  `.cache/` (os dois últimos sem `/` inicial, para casar tanto com a raiz
  quanto com `app/`, já que `coverage.out` agora é gerado dentro de
  `app/`).

## Bug pré-existente encontrado durante o teste pós-mudança

Testando `make docs-update` depois da reorganização, apareceu um erro
(`no such file or directory` para `INSTALL_REQUIREMENTS.md`). Não era
causado pela movimentação de arquivos: `REPO_ROOT := $(abspath ..)` no
`Makefile` original já estava **errado** desde a criação (resolvia para
`agents/`, não para a raiz do repositório `learning-istio/`, dois níveis
acima de `agents/upgrade/`). Esse alvo nunca tinha sido testado
efetivamente antes (na sessão anterior, o `docs update --apply` real foi
rodado passando `--repo-root` manualmente na linha de comando, não via
`make docs-update-apply`). Corrigido para `REPO_ROOT := $(abspath ../..)`
e revalidado.

## Verificação

Depois da movimentação e dos ajustes no `Makefile`/`.gitignore`, foram
re-executados com sucesso: `make build` (binário funcional em
`agents/upgrade/bin/`), `make test`, `make lint` (0 issues), `make
coverage` (mesmos percentuais de antes), `make docker-build` +
`make docker-run ARGS="--version"`, `make versions-latest`,
`make analyze ARGS="..."` e `make docs-update` (com o `REPO_ROOT`
corrigido, voltou a reportar corretamente "already up to date" para todos
os arquivos já sincronizados).

`README.md` (diagrama mermaid de componentes, seção de instalação,
estrutura do projeto, seção de testes) e `CLAUDE.md` foram atualizados
para refletir os novos caminhos (`app/main.go`, `app/cmd/`,
`app/internal/...`, `app/VERSION`, `app/CHANGELOG.md`, `app/go.mod`,
`app/examples/`).
