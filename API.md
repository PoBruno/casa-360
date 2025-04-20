# API Casa360 - Documentação

API para gerenciamento de tarefas e finanças compartilhadas, permitindo controle de recursos domésticos.

## Base URL

```
http://localhost:3001
```

## Autenticação

Atualmente, a API não requer autenticação.

## Formatos

- Todas as requisições e respostas utilizam o formato JSON.
- Datas devem ser enviadas no formato ISO 8601: `YYYY-MM-DDThh:mm:ssZ`.
- Os IDs são no formato UUID v4.

## Endpoints

### Usuários

#### Criar um usuário

```
POST /users
```

**Corpo da requisição:**
```json
{
  "name": "Nome do Usuário"
}
```

**Resposta (201 Created):**
```json
{
  "id": "uuid",
  "name": "Nome do Usuário"
}
```

#### Listar todos os usuários

```
GET /users
```

**Resposta (200 OK):**
```json
[
  {
    "id": "uuid",
    "name": "Nome do Usuário 1"
  },
  {
    "id": "uuid",
    "name": "Nome do Usuário 2"
  }
]
```

#### Buscar um usuário pelo ID

```
GET /users/:id
```

**Resposta (200 OK):**
```json
{
  "id": "uuid",
  "name": "Nome do Usuário"
}
```

#### Atualizar um usuário

```
PUT /users/:id
```

**Corpo da requisição:**
```json
{
  "name": "Novo Nome do Usuário"
}
```

**Resposta (200 OK):**
```json
{
  "id": "uuid",
  "name": "Novo Nome do Usuário"
}
```

#### Remover um usuário

```
DELETE /users/:id
```

**Resposta (204 No Content)**

### Grupos de Pagadores

#### Criar um grupo de pagadores

```
POST /payer-groups
```

**Corpo da requisição:**
```json
{
  "name": "Nome do Grupo"
}
```

**Resposta (201 Created):**
```json
{
  "id": "uuid",
  "name": "Nome do Grupo"
}
```

#### Listar todos os grupos de pagadores

```
GET /payer-groups
```

**Resposta (200 OK):**
```json
[
  {
    "id": "uuid",
    "name": "Nome do Grupo 1"
  },
  {
    "id": "uuid",
    "name": "Nome do Grupo 2"
  }
]
```

#### Buscar um grupo pelo ID

```
GET /payer-groups/:id
```

**Resposta (200 OK):**
```json
{
  "id": "uuid",
  "name": "Nome do Grupo"
}
```

#### Atualizar um grupo

```
PUT /payer-groups/:id
```

**Corpo da requisição:**
```json
{
  "name": "Novo Nome do Grupo"
}
```

**Resposta (200 OK):**
```json
{
  "id": "uuid",
  "name": "Novo Nome do Grupo"
}
```

#### Remover um grupo

```
DELETE /payer-groups/:id
```

**Resposta (204 No Content)**

#### Adicionar um membro ao grupo

```
POST /payer-groups/:id/members
```

**Corpo da requisição:**
```json
{
  "user_id": "uuid",
  "percentage": 50.00
}
```

**Resposta (201 Created):**
```json
{
  "id": "uuid",
  "payer_group_id": "uuid",
  "user_id": "uuid",
  "percentage": 50.00
}
```

**Observação:** A soma dos percentuais de todos os membros não pode exceder 100%.

#### Listar membros do grupo

```
GET /payer-groups/:id/members
```

**Resposta (200 OK):**
```json
[
  {
    "id": "uuid",
    "payer_group_id": "uuid",
    "user_id": "uuid",
    "percentage": 60.00,
    "user": {
      "id": "uuid",
      "name": "Nome do Usuário 1"
    }
  },
  {
    "id": "uuid",
    "payer_group_id": "uuid",
    "user_id": "uuid",
    "percentage": 40.00,
    "user": {
      "id": "uuid",
      "name": "Nome do Usuário 2"
    }
  }
]
```

#### Remover um membro do grupo

```
DELETE /payer-groups/:id/members/:member_id
```

**Resposta (204 No Content)**

### Centro de Custo

#### Criar um centro de custo

```
POST /finance-cc
```

**Corpo da requisição:**
```json
{
  "name": "Nome do Centro de Custo",
  "parent_id": "uuid" // opcional
}
```

**Resposta (201 Created):**
```json
{
  "id": "uuid",
  "name": "Nome do Centro de Custo",
  "parent_id": "uuid" // null se não tiver parent
}
```

#### Listar todos os centros de custo

```
GET /finance-cc
```

**Resposta (200 OK):**
```json
[
  {
    "id": "uuid",
    "name": "Moradia",
    "parent_id": null
  },
  {
    "id": "uuid",
    "name": "Aluguel",
    "parent_id": "uuid" // ID do centro de custo "Moradia"
  }
]
```

### Moedas

#### Criar uma moeda

```
POST /currencies
```

**Corpo da requisição:**
```json
{
  "name": "Nome da Moeda",
  "symbol": "Símbolo",
  "value": 1.0000
}
```

**Resposta (201 Created):**
```json
{
  "id": "uuid",
  "name": "Nome da Moeda",
  "symbol": "Símbolo",
  "value": 1.0000
}
```

#### Listar todas as moedas

```
GET /currencies
```

**Resposta (200 OK):**
```json
[
  {
    "id": "uuid",
    "name": "Real",
    "symbol": "R$",
    "value": 1.0000
  },
  {
    "id": "uuid",
    "name": "Dólar",
    "symbol": "$",
    "value": 5.2000
  }
]
```

### Tarefas

#### Criar uma tarefa

```
POST /tasks
```

**Corpo da requisição:**
```json
{
  "title": "Título da Tarefa",
  "description": "Descrição da tarefa",
  "start_date": "2023-01-01T00:00:00Z",
  "recurrence_cron": "0 0 * * 1", // Toda segunda-feira
  "subtasks": [
    {
      "title": "Subtarefa 1",
      "done": false
    }
  ],
  "user_id": "uuid",
  "payer_group_id": "uuid"
}
```

**Resposta (201 Created):**
```json
{
  "id": "uuid",
  "title": "Título da Tarefa",
  "description": "Descrição da tarefa",
  "start_date": "2023-01-01T00:00:00Z",
  "recurrence_cron": "0 0 * * 1",
  "subtasks": [
    {
      "title": "Subtarefa 1",
      "done": false
    }
  ],
  "user_id": "uuid",
  "payer_group_id": "uuid"
}
```

#### Listar todas as tarefas

```
GET /tasks
```

**Resposta (200 OK):**
```json
[
  {
    "id": "uuid",
    "title": "Título da Tarefa 1",
    "description": "Descrição da tarefa 1",
    "start_date": "2023-01-01T00:00:00Z",
    "recurrence_cron": "0 0 * * 1",
    "subtasks": [],
    "user_id": "uuid",
    "payer_group_id": "uuid"
  },
  {
    "id": "uuid",
    "title": "Título da Tarefa 2",
    "description": "Descrição da tarefa 2",
    "start_date": "2023-01-01T00:00:00Z",
    "recurrence_cron": "0 0 * * 3",
    "subtasks": [],
    "user_id": "uuid",
    "payer_group_id": "uuid"
  }
]
```

#### Buscar uma tarefa pelo ID

```
GET /tasks/:id
```

**Resposta (200 OK):**
```json
{
  "id": "uuid",
  "title": "Título da Tarefa",
  "description": "Descrição da tarefa",
  "start_date": "2023-01-01T00:00:00Z",
  "recurrence_cron": "0 0 * * 1",
  "subtasks": [],
  "user_id": "uuid",
  "payer_group_id": "uuid"
}
```

#### Atualizar uma tarefa

```
PUT /tasks/:id
```

**Corpo da requisição:**
```json
{
  "title": "Novo Título da Tarefa",
  "description": "Nova descrição da tarefa",
  "start_date": "2023-01-01T00:00:00Z",
  "recurrence_cron": "0 0 * * 2", // Toda terça-feira
  "subtasks": [],
  "user_id": "uuid",
  "payer_group_id": "uuid"
}
```

**Resposta (200 OK):**
```json
{
  "id": "uuid",
  "title": "Novo Título da Tarefa",
  "description": "Nova descrição da tarefa",
  "start_date": "2023-01-01T00:00:00Z",
  "recurrence_cron": "0 0 * * 2",
  "subtasks": [],
  "user_id": "uuid",
  "payer_group_id": "uuid"
}
```

#### Remover uma tarefa

```
DELETE /tasks/:id
```

**Resposta (204 No Content)**

#### Atualizar ocorrências de tarefas

```
POST /tasks/update-occurrences
```

Este endpoint gera automaticamente ocorrências para todas as tarefas, até a data atual, baseado na expressão CRON definida para cada tarefa.

**Resposta (200 OK):**
Evento SSE (Server-Sent Events) com atualizações em tempo real.

#### Gerar ocorrências para uma tarefa específica

```
POST /tasks/:id/occurrences
```

Este endpoint gera ocorrências para uma tarefa específica.

**Resposta (200 OK)**

### Ocorrências de Tarefas

#### Criar uma ocorrência de tarefa manualmente

```
POST /task-occurrences
```

**Corpo da requisição:**
```json
{
  "task_id": "uuid",
  "date": "2023-01-01T00:00:00Z",
  "status": false,
  "user_id": "uuid",
  "payer_group_id": "uuid",
  "subtasks": []
}
```

**Resposta (201 Created):**
```json
{
  "id": "uuid",
  "task_id": "uuid",
  "date": "2023-01-01T00:00:00Z",
  "status": false,
  "user_id": "uuid",
  "payer_group_id": "uuid",
  "subtasks": []
}
```

#### Listar todas as ocorrências de tarefas

```
GET /task-occurrences
```

**Resposta (200 OK):**
```json
[
  {
    "id": "uuid",
    "task_id": "uuid",
    "date": "2023-01-01T00:00:00Z",
    "status": false,
    "user_id": "uuid",
    "payer_group_id": "uuid",
    "subtasks": []
  },
  {
    "id": "uuid",
    "task_id": "uuid",
    "date": "2023-01-08T00:00:00Z",
    "status": true,
    "user_id": "uuid",
    "payer_group_id": "uuid",
    "subtasks": []
  }
]
```

#### Atualizar uma ocorrência de tarefa

```
PUT /task-occurrences/:id
```

**Corpo da requisição:**
```json
{
  "status": true, // atualizar status (completo)
  "user_id": "uuid", // opcional - atualizar responsável
  "payer_group_id": "uuid", // opcional - atualizar grupo pagador
  "subtasks": [] // opcional - atualizar subtarefas
}
```

**Resposta (200 OK):**
```json
{
  "id": "uuid",
  "task_id": "uuid",
  "date": "2023-01-01T00:00:00Z",
  "status": true,
  "user_id": "uuid",
  "payer_group_id": "uuid",
  "subtasks": []
}
```

#### Remover uma ocorrência de tarefa

```
DELETE /task-occurrences/:id
```

**Resposta (204 No Content)**

### Finanças

#### Criar uma finança

```
POST /finances
```

**Corpo da requisição:**
```json
{
  "title": "Título da Finança",
  "description": "Descrição da finança",
  "type": true, // true = despesa, false = receita
  "start_date": "2023-01-01T00:00:00Z",
  "end_date": "2023-12-31T00:00:00Z", // opcional
  "recurrence_days": 30,
  "amount": 1000.00,
  "user_id": "uuid",
  "payer_group_id": "uuid",
  "finance_cc_id": "uuid",
  "currency_id": "uuid"
}
```

**Resposta (201 Created):**
```json
{
  "id": "uuid",
  "title": "Título da Finança",
  "description": "Descrição da finança",
  "type": true,
  "start_date": "2023-01-01T00:00:00Z",
  "end_date": "2023-12-31T00:00:00Z",
  "recurrence_days": 30,
  "amount": 1000.00,
  "user_id": "uuid",
  "payer_group_id": "uuid",
  "finance_cc_id": "uuid",
  "currency_id": "uuid"
}
```

#### Listar todas as finanças

```
GET /finances
```

**Resposta (200 OK):**
```json
[
  {
    "id": "uuid",
    "title": "Aluguel",
    "description": "Pagamento mensal",
    "type": true,
    "start_date": "2023-01-01T00:00:00Z",
    "end_date": null,
    "recurrence_days": 30,
    "amount": 1500.00,
    "user_id": "uuid",
    "payer_group_id": "uuid",
    "finance_cc_id": "uuid",
    "currency_id": "uuid"
  },
  {
    "id": "uuid",
    "title": "Salário",
    "description": "Recebimento mensal",
    "type": false,
    "start_date": "2023-01-05T00:00:00Z",
    "end_date": null,
    "recurrence_days": 30,
    "amount": 3000.00,
    "user_id": "uuid",
    "payer_group_id": "uuid",
    "finance_cc_id": "uuid",
    "currency_id": "uuid"
  }
]
```

#### Buscar uma finança pelo ID

```
GET /finances/:id
```

**Resposta (200 OK):**
```json
{
  "id": "uuid",
  "title": "Aluguel",
  "description": "Pagamento mensal",
  "type": true,
  "start_date": "2023-01-01T00:00:00Z",
  "end_date": null,
  "recurrence_days": 30,
  "amount": 1500.00,
  "user_id": "uuid",
  "payer_group_id": "uuid",
  "finance_cc_id": "uuid",
  "currency_id": "uuid"
}
```

#### Atualizar uma finança

```
PUT /finances/:id
```

**Corpo da requisição:**
```json
{
  "title": "Novo Título da Finança",
  "description": "Nova descrição da finança",
  "type": true,
  "start_date": "2023-01-01T00:00:00Z",
  "end_date": "2023-12-31T00:00:00Z",
  "recurrence_days": 30,
  "amount": 1200.00,
  "user_id": "uuid",
  "payer_group_id": "uuid",
  "finance_cc_id": "uuid",
  "currency_id": "uuid"
}
```

**Resposta (200 OK):**
```json
{
  "id": "uuid",
  "title": "Novo Título da Finança",
  "description": "Nova descrição da finança",
  "type": true,
  "start_date": "2023-01-01T00:00:00Z",
  "end_date": "2023-12-31T00:00:00Z",
  "recurrence_days": 30,
  "amount": 1200.00,
  "user_id": "uuid",
  "payer_group_id": "uuid",
  "finance_cc_id": "uuid",
  "currency_id": "uuid"
}
```

#### Remover uma finança

```
DELETE /finances/:id
```

**Resposta (204 No Content)**

#### Atualizar ocorrências de finanças

```
POST /finances/update-occurrences
```

Este endpoint gera automaticamente ocorrências para todas as finanças, até a data atual, baseado no intervalo de recorrência definido para cada finança.

**Resposta (200 OK):**
Evento SSE (Server-Sent Events) com atualizações em tempo real.

#### Gerar ocorrências para uma finança específica

```
POST /finances/:id/occurrences
```

Este endpoint gera ocorrências para uma finança específica.

**Resposta (200 OK)**

### Ocorrências Financeiras

#### Criar uma ocorrência financeira manualmente

```
POST /finance-occurrences
```

**Corpo da requisição:**
```json
{
  "finance_id": "uuid",
  "date": "2023-01-01T00:00:00Z",
  "amount": 1000.00,
  "status": false
}
```

**Resposta (201 Created):**
```json
{
  "id": "uuid",
  "finance_id": "uuid",
  "date": "2023-01-01T00:00:00Z",
  "amount": 1000.00,
  "status": false
}
```

#### Listar todas as ocorrências financeiras

```
GET /finance-occurrences
```

**Resposta (200 OK):**
```json
[
  {
    "id": "uuid",
    "finance_id": "uuid",
    "date": "2023-01-01T00:00:00Z",
    "amount": 1500.00,
    "status": false
  },
  {
    "id": "uuid",
    "finance_id": "uuid",
    "date": "2023-01-31T00:00:00Z",
    "amount": 1500.00,
    "status": true
  }
]
```

#### Atualizar uma ocorrência financeira

```
PUT /finance-occurrences/:id
```

**Corpo da requisição:**
```json
{
  "amount": 1200.00, // opcional - atualizar valor
  "status": true // opcional - atualizar status (pago/recebido)
}
```

**Resposta (200 OK):**
```json
{
  "id": "uuid",
  "finance_id": "uuid",
  "date": "2023-01-01T00:00:00Z",
  "amount": 1200.00,
  "status": true
}
```

**Observação:** Quando o status muda para `true`, é criada uma transação e as carteiras dos usuários do grupo pagador são atualizadas de acordo com os percentuais definidos.

#### Remover uma ocorrência financeira

```
DELETE /finance-occurrences/:id
```

**Resposta (204 No Content)**

### Dashboard

#### Listar todas as ocorrências (tarefas e finanças)

```
GET /occurrences/dashboard
```

**Resposta (200 OK):**
```json
[
  {
    "occurrence_type": "finance",
    "id": "uuid",
    "date": "2023-01-01T00:00:00Z",
    "status": true,
    "title": "Aluguel",
    "description": "Pagamento mensal",
    "finance_type": true,
    "amount": 1500.00,
    "currency_symbol": "R$",
    "currency_value": 1.0000,
    "amount_converted": 1500.00,
    "cost_center": "Moradia",
    "payer_group": "Casa",
    "responsible_user": "João"
  },
  {
    "occurrence_type": "task",
    "id": "uuid",
    "date": "2023-01-01T00:00:00Z",
    "status": false,
    "title": "Limpar Casa",
    "description": "Limpeza semanal",
    "finance_type": null,
    "amount": null,
    "currency_symbol": null,
    "currency_value": null,
    "amount_converted": null,
    "cost_center": null,
    "payer_group": "Casa",
    "responsible_user": "Maria"
  }
]
```

### Carteiras

#### Obter a última carteira de um usuário

```
GET /wallets/:user_id
```

**Resposta (200 OK):**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "amount": 2500.00,
  "created_at": "2023-01-15T10:30:00Z"
}
```

### Transações

#### Listar transações de uma ocorrência financeira

```
GET /transactions/:occurrence_id
```

**Resposta (200 OK):**
```json
[
  {
    "id": "uuid",
    "finance_occurrence_id": "uuid",
    "amount": 1500.00,
    "created_at": "2023-01-01T12:00:00Z"
  }
]
```

## Exemplos de Uso

### Criar um Usuário

```bash
curl -X POST http://localhost:3001/users \
  -H "Content-Type: application/json" \
  -d '{"name": "João Silva"}'
```

### Criar uma Tarefa Recorrente

```bash
curl -X POST http://localhost:3001/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Limpeza Geral",
    "description": "Limpeza semanal da casa",
    "start_date": "2023-01-01T00:00:00Z",
    "recurrence_cron": "0 8 * * 1",
    "subtasks": [
      {"title": "Limpar sala", "done": false},
      {"title": "Limpar cozinha", "done": false}
    ],
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "payer_group_id": "123e4567-e89b-12d3-a456-426614174001"
  }'
```

### Atualizar Status de uma Ocorrência de Tarefa

```bash
curl -X PUT http://localhost:3001/task-occurrences/123e4567-e89b-12d3-a456-426614174002 \
  -H "Content-Type: application/json" \
  -d '{"status": true}'
```

### Criar uma Finança Recorrente

```bash
curl -X POST http://localhost:3001/finances \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Aluguel",
    "description": "Pagamento mensal de aluguel",
    "type": true,
    "start_date": "2023-01-01T00:00:00Z",
    "recurrence_days": 30,
    "amount": 1500.00,
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "payer_group_id": "123e4567-e89b-12d3-a456-426614174001",
    "finance_cc_id": "123e4567-e89b-12d3-a456-426614174003",
    "currency_id": "123e4567-e89b-12d3-a456-426614174004"
  }'
```

### Marcar uma Ocorrência Financeira como Paga

```bash
curl -X PUT http://localhost:3001/finance-occurrences/123e4567-e89b-12d3-a456-426614174005 \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1500.00,
    "status": true
  }'
```

## Funcionalidades Automáticas

1. **Geração de Ocorrências:** As ocorrências de tarefas e finanças são geradas automaticamente baseadas nas definições de recorrência de cada item.

2. **Transações Financeiras:** Quando uma ocorrência financeira é marcada como concluída (paga ou recebida), o sistema automaticamente:
   - Cria uma transação para o registro
   - Atualiza a carteira de cada usuário do grupo pagador com base no seu percentual
   - Para despesas, subtrai o valor da carteira de cada usuário
   - Para receitas, adiciona o valor à carteira de cada usuário

3. **Conversão de Moedas:** Todas as transações são armazenadas com o valor convertido de acordo com a taxa de câmbio da moeda.

## Códigos de Erro

- `400 Bad Request` - Requisição inválida ou dados malformados
- `404 Not Found` - Recurso não encontrado
- `500 Internal Server Error` - Erro interno do servidor 