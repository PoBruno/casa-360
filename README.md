# Casa360 API

API para gerenciamento de tarefas e finanças compartilhadas.

## Configuração

1. Crie um arquivo `.env` com as configurações:
```env
PORT=3001
DB_HOST=localhost
DB_PORT=5432
DB_USER=casa360
DB_PASSWORD=casa360
DB_NAME=casa360
```

2. Execute o PostgreSQL (recomendado usar Docker):
```bash
docker compose up -d
```

3. Execute a aplicação:
```bash
go run main.go
```

## Endpoints da API

### Usuários

- `POST /users` - Cria um novo usuário
  ```json
  {
    "name": "Nome do Usuário"
  }
  ```

- `GET /users` - Lista todos os usuários
- `GET /users/:id` - Busca um usuário pelo ID
- `PUT /users/:id` - Atualiza um usuário
- `DELETE /users/:id` - Remove um usuário

### Grupos de Pagadores

- `POST /payer-groups` - Cria um novo grupo
  ```json
  {
    "name": "Nome do Grupo"
  }
  ```

- `GET /payer-groups` - Lista todos os grupos
- `GET /payer-groups/:id` - Busca um grupo pelo ID
- `PUT /payer-groups/:id` - Atualiza um grupo
- `DELETE /payer-groups/:id` - Remove um grupo

#### Membros do Grupo

- `POST /payer-groups/:id/members` - Adiciona um membro ao grupo
  ```json
  {
    "user_id": "uuid",
    "percentage": 50.00
  }
  ```

- `GET /payer-groups/:id/members` - Lista membros do grupo
- `DELETE /payer-groups/:id/members/:member_id` - Remove um membro

### Centro de Custo

- `POST /finance-cc` - Cria um novo centro de custo
  ```json
  {
    "name": "Nome do CC",
    "parent_id": "uuid" // opcional
  }
  ```

- `GET /finance-cc` - Lista todos os centros de custo

### Moedas

- `POST /currencies` - Cria uma nova moeda
  ```json
  {
    "name": "Nome da Moeda",
    "symbol": "Símbolo",
    "value": 1.0000
  }
  ```

- `GET /currencies` - Lista todas as moedas

### Tarefas

- `POST /tasks` - Cria uma nova tarefa
  ```json
  {
    "title": "Título da Tarefa",
    "description": "Descrição",
    "start_date": "2024-01-01",
    "recurrence_cron": "0 0 * * 1",
    "subtasks": [],
    "user_id": "uuid",
    "payer_group_id": "uuid"
  }
  ```

- `GET /tasks` - Lista todas as tarefas
- `GET /tasks/:id` - Busca uma tarefa pelo ID
- `PUT /tasks/:id` - Atualiza uma tarefa
- `DELETE /tasks/:id` - Remove uma tarefa
- `POST /tasks/update-occurrences` - Atualiza ocorrências de todas as tarefas

#### Ocorrências de Tarefas

- `POST /tasks/:id/occurrences` - Gera ocorrências para uma tarefa
- `POST /task-occurrences` - Cria uma ocorrência manual
- `GET /task-occurrences` - Lista todas as ocorrências
- `PUT /task-occurrences/:id` - Atualiza uma ocorrência
- `DELETE /task-occurrences/:id` - Remove uma ocorrência

### Finanças

- `POST /finances` - Cria uma nova finança
  ```json
  {
    "title": "Título da Finança",
    "description": "Descrição",
    "type": false,
    "start_date": "2024-01-01",
    "end_date": null,
    "recurrence_days": 30,
    "amount": 100.00,
    "user_id": "uuid",
    "payer_group_id": "uuid",
    "finance_cc_id": "uuid",
    "currency_id": "uuid"
  }
  ```

- `GET /finances` - Lista todas as finanças
- `GET /finances/:id` - Busca uma finança pelo ID
- `PUT /finances/:id` - Atualiza uma finança
- `DELETE /finances/:id` - Remove uma finança
- `POST /finances/update-occurrences` - Atualiza ocorrências de todas as finanças

#### Ocorrências Financeiras

- `POST /finances/:id/occurrences` - Gera ocorrências para uma finança
- `POST /finance-occurrences` - Cria uma ocorrência manual
- `GET /finance-occurrences` - Lista todas as ocorrências
- `PUT /finance-occurrences/:id` - Atualiza uma ocorrência
- `DELETE /finance-occurrences/:id` - Remove uma ocorrência

### Dashboard e Carteiras

- `GET /occurrences/dashboard` - Retorna todas as ocorrências (tarefas e finanças)
- `GET /wallets/:user_id` - Retorna o último saldo da carteira do usuário
- `GET /transactions/:occurrence_id` - Lista transações de uma ocorrência

## Funcionalidades Automáticas

1. Quando uma ocorrência financeira é marcada como concluída (status = true):
   - Uma transação é registrada com o valor convertido pela taxa de câmbio
   - As carteiras dos usuários são atualizadas conforme seus percentuais no grupo
   - Receitas são somadas e despesas são subtraídas dos saldos

2. A view do dashboard unifica:
   - Ocorrências de tarefas e finanças
   - Informações detalhadas de todas as tabelas relacionadas
   - Valores convertidos para a moeda base
   - Responsáveis e grupos de pagadores 




```sql

-- INSERÇÕES DE DADOS DE EXEMPLO

-- 1) Usuários
INSERT INTO users (id, name) VALUES
  ('11111111-1111-1111-1111-111111111111','João'),
  ('22222222-2222-2222-2222-222222222222','Maria');

-- 2) Grupos de pagadores
INSERT INTO payer_groups (id, name) VALUES
  ('33333333-3333-3333-3333-333333333333','casa'),
  ('44444444-4444-4444-4444-444444444444','joao_solo'),
  ('55555555-5555-5555-5555-555555555555','maria_solo');

-- 3) Membros dos grupos (rateio)
INSERT INTO payer_group_members (payer_group_id, user_id, percentage) VALUES
  -- Grupo "casa": João 50%, Maria 50%
  ('33333333-3333-3333-3333-333333333333','11111111-1111-1111-1111-111111111111',50.00),
  ('33333333-3333-3333-3333-333333333333','22222222-2222-2222-2222-222222222222',50.00),
  -- Grupo solo de João: 100% João
  ('44444444-4444-4444-4444-444444444444','11111111-1111-1111-1111-111111111111',100.00),
  -- Grupo solo de Maria: 100% Maria
  ('55555555-5555-5555-5555-555555555555','22222222-2222-2222-2222-222222222222',100.00);

-- 4) Moedas
INSERT INTO finance_currency (id, name, symbol, value) VALUES
  ('66666666-6666-6666-6666-666666666666','Real','R$',1.0000),
  ('77777777-7777-7777-7777-777777777777','Dólar','$',5.8000);

-- 5) Centros de Custo (hierarquia pai/filho)
INSERT INTO finance_cc (id, name, parent_id) VALUES
  -- Pais
  ('88888888-8888-8888-8888-888888888888','Moradia',NULL),
  ('99999999-9999-9999-9999-999999999999','Serviços',NULL),
  -- Filhos
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa','Aluguel','88888888-8888-8888-8888-888888888888'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb','Água','99999999-9999-9999-9999-999999999999'),
  ('cccccccc-cccc-cccc-cccc-cccccccccccc','Luz','99999999-9999-9999-9999-999999999999'),
  ('dddddddd-dddd-dddd-dddd-dddddddddddd','Internet','99999999-9999-9999-9999-999999999999');

-- 6) Tarefas recorrentes (task_installments)
INSERT INTO task_installments (
    id, title, description, start_date, recurrence_cron,
    subtasks, user_id, payer_group_id
) VALUES
  -- "Dia do Lixo": toda terça e quinta, Maria responsável, rateio "casa"
  (
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee',
    'Dia do Lixo',
    'Recolher o lixo da casa',
    '2025-04-01',
    '0 0 * * 2,4',
    '[]',
    '22222222-2222-2222-2222-222222222222',  -- Maria
    '33333333-3333-3333-3333-333333333333'   -- casa
  ),
  -- "Caixinha dos Gator": a cada 2 dias, João responsável, rateio solo João
  (
    'ffffffff-ffff-ffff-ffff-ffffffffffff',
    'Caixinha dos Gator',
    'Contribuição quinzenal do Gator',
    '2025-04-01',
    '0 0 */2 * *',
    '[]',
    '11111111-1111-1111-1111-111111111111',  -- João
    '44444444-4444-4444-4444-444444444444'   -- joao_solo
  );

-- 7) Contas financeiras recorrentes (finance_installments)
INSERT INTO finance_installments (
    id, title, description, type, start_date, end_date,
    recurrence_days, amount, user_id, payer_group_id,
    finance_cc_id, currency_id
) VALUES
  -- Aluguel: despesa, inicia em 05/01/2025, 30 dias, Maria solo
  (
    '10101010-1010-1010-1010-101010101010',
    'Aluguel',
    'Pagamento do aluguel mensal',
    TRUE,              -- despesa
    '2025-01-05',
    NULL,
    30,
    1500.00,
    '22222222-2222-2222-2222-222222222222',  -- Maria
    '55555555-5555-5555-5555-555555555555',  -- maria_solo
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',  -- Aluguel CC
    '66666666-6666-6666-6666-666666666666'   -- BRL
  ),
  -- Água: despesa, inicia em 10/01/2025, 30 dias, João solo
  (
    '20202020-2020-2020-2020-202020202020',
    'Água',
    'Conta de água',
    TRUE,
    '2025-01-10',
    NULL,
    30,
    100.00,
    '11111111-1111-1111-1111-111111111111',  -- João
    '44444444-4444-4444-4444-444444444444',  -- joao_solo
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',  -- Água CC
    '66666666-6666-6666-6666-666666666666'   -- BRL
  ),
  -- Luz: despesa, inicia em 15/01/2025, 30 dias, rateio "casa"
  (
    '30303030-3030-3030-3030-303030303030',
    'Luz',
    'Conta de luz',
    TRUE,
    '2025-01-15',
    NULL,
    30,
    200.00,
    '22222222-2222-2222-2222-222222222222',  -- Maria
    '33333333-3333-3333-3333-333333333333',  -- casa
    'cccccccc-cccc-cccc-cccc-cccccccccccc',  -- Luz CC
    '66666666-6666-6666-6666-666666666666'   -- BRL
  ),
  -- Internet: despesa, inicia em 20/01/2025, 30 dias, rateio "casa"
  (
    '40404040-4040-4040-4040-404040404040',
    'Internet',
    'Assinatura de internet',
    TRUE,
    '2025-01-20',
    NULL,
    30,
    120.00,
    '11111111-1111-1111-1111-111111111111',  -- João
    '33333333-3333-3333-3333-333333333333',  -- casa
    'dddddddd-dddd-dddd-dddd-dddddddddddd',  -- Internet CC
    '77777777-7777-7777-7777-777777777777'   -- USD
  );

```