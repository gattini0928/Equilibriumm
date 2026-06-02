# Equilibriumm - Sistema de Consultas Terapêuticas e Psiquiátricas

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge\&logo=go\&logoColor=white)
![HTMX](https://img.shields.io/badge/HTMX-3366CC?style=for-the-badge\&logo=htmx\&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-000000?style=for-the-badge\&logo=jsonwebtokens\&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge\&logo=postgresql\&logoColor=white)
![templ](https://img.shields.io/badge/templ-222222?style=for-the-badge\&logo=go\&logoColor=white)

## Sobre o Projeto

**Equilibriumm** é um sistema web desenvolvido em Go para gerenciamento de consultas entre pacientes, terapeutas e psiquiatras.

O projeto foi criado com foco em aprendizado, prática de arquitetura web, autenticação, renderização server-side, integração com HTMX, persistência de dados com PostgreSQL e organização de fluxos entre diferentes tipos de usuários.

> Este projeto é open source e possui finalidade educacional/portfólio. Não deve ser utilizado em ambiente real de saúde sem validações legais, segurança adequada, conformidade com LGPD e autorização dos órgãos competentes.

## Tecnologias Utilizadas

* Go
* templ
* HTMX
* PostgreSQL
* JWT
* HTML
* CSS
* Jitsi Meet API

## Funcionalidades

- Cadastro e autenticação de usuários.
- Login utilizando JWT.
- Perfis distintos para pacientes, terapeutas e psiquiatras.
- Agendamento de horários disponíveis.
- Vinculação de pacientes a terapeutas e psiquiatras.
- Gerenciamento de agendas profissionais.
- Reserva e cancelamento de consultas.
- Consultas online por videoconferência utilizando Jitsi Meet.
- Registro de anotações terapêuticas.
- Registro de prescrições psiquiátricas.
- Visualização de pacientes do terapeuta e psiquiatra.
- Visualização do histórico de consultas.
- Dashboard personalizado para cada tipo de usuário.

## Como Rodar o Projeto

```bash
git clone https://github.com/seu-usuario/seu-repositorio.git
cd seu-repositorio
go mod tidy
```

Depois, configure as variáveis de ambiente(env.example) e o banco de dados PostgreSQL conforme necessário.

```bash
make dev
```

## O que aprendi com este projeto

- Estruturação de aplicações web em Go.
- Renderização server-side utilizando templ.
- Integração com HTMX para atualizações parciais.
- Autenticação baseada em JWT.
- Modelagem de banco de dados relacional com PostgreSQL.
- Organização de arquitetura em camadas (routes, handlers, services e repositories).
- Integração de videoconferência utilizando Jitsi Meet API.

## Possíveis Melhorias

* Atualizar dinamicamente o perfil de terapeuta ou psiquiatra quando o paciente agenda um horário.
* Impedir inconsistências entre data da agenda e data real da consulta, por exemplo: agenda marcada para 10/06/2026, mas consulta finalizada/salva em 09/06/2026.
* Evitar que o paciente precise atualizar a página para receber corretamente o ID da consulta iniciada pelo terapeuta ou psiquiatra.
* Adicionar campo de avaliação/rating para terapeutas e psiquiatras.
* Exibir profissionais melhor avaliados na homepage, como uma seção de “Melhores Ranqueados”.
* Melhorar a interface visual do sistema no geral.
* Adicionar funcionalidade para alteração de informações pessoais.
* Permitir edição/correção de receitas ou registros preenchidos incorretamente durante a consulta.
* Melhorar validações de formulários.
* Melhorar responsividade para dispositivos móveis.
* Criar testes automatizados.
* Adicionar documentação mais detalhada das rotas e fluxos do sistema.
* Refinar regras de permissão entre pacientes, terapeutas e psiquiatras.
* Melhorar tratamento de erros e mensagens para o usuário.

## Status do Projeto

Projeto em desenvolvimento/finalização para fins de estudo e portfólio.

## Licença

Este projeto é open source. A licença pode ser definida posteriormente.
