#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# URL base da API
BASE_URL="http://localhost:3001"

# Função para log
log() {
    echo -e "${YELLOW}[$(date '+%Y-%m-%d %H:%M:%S')] $1${NC}"
}

# Função para separador de seção
section() {
    echo -e "${BLUE}=================================================${NC}"
    echo -e "${BLUE}== $1${NC}"
    echo -e "${BLUE}=================================================${NC}"
}

# Função para testar resposta
test_response() {
    if [ $1 -eq $2 ]; then
        echo -e "${GREEN}✓ Sucesso: $3${NC}"
        return 0
    else
        echo -e "${RED}✗ Erro ($1): $3${NC}"
        return 1
    fi
}

# Função para mostrar detalhes da resposta
show_response() {
    echo -e "${BLUE}Detalhes da resposta:${NC}"
    echo "$1" | jq
    echo ""
}

# Variáveis para armazenar IDs
USER1_ID=""
USER2_ID=""
PAYER_GROUP_ID=""
FINANCE_CC_ID=""
CURRENCY_ID=""
TASK_ID=""
TASK_OCCURRENCE_ID=""
FINANCE_ID=""
FINANCE_OCCURRENCE_ID=""

section "INICIANDO TESTES COMPLETOS DA API CASA360"

section "1. USUÁRIOS"
log "Criando usuário 1 (João)"
response=$(curl -s -X POST $BASE_URL/users -H "Content-Type: application/json" -d '{
    "name": "João"
}')
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/users -H "Content-Type: application/json" -d '{
    "name": "João"
}')
test_response $status_code 201 "Criar usuário João"
USER1_ID=$(echo $response | jq -r '.id')
show_response "$response"

log "Criando usuário 2 (Maria)"
response=$(curl -s -X POST $BASE_URL/users -H "Content-Type: application/json" -d '{
    "name": "Maria"
}')
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/users -H "Content-Type: application/json" -d '{
    "name": "Maria"
}')
test_response $status_code 201 "Criar usuário Maria"
USER2_ID=$(echo $response | jq -r '.id')
show_response "$response"

log "Listando todos os usuários"
response=$(curl -s -X GET $BASE_URL/users)
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X GET $BASE_URL/users)
test_response $status_code 200 "Listar usuários"
show_response "$response"

section "2. GRUPOS DE PAGADORES"
log "Criando grupo de pagadores (Casa)"
response=$(curl -s -X POST $BASE_URL/payer-groups -H "Content-Type: application/json" -d '{
    "name": "Casa"
}')
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/payer-groups -H "Content-Type: application/json" -d '{
    "name": "Casa"
}')
test_response $status_code 201 "Criar grupo Casa"
PAYER_GROUP_ID=$(echo $response | jq -r '.id')
show_response "$response"

log "Adicionando João ao grupo (60%)"
response=$(curl -s -X POST "$BASE_URL/payer-groups/$PAYER_GROUP_ID/members" -H "Content-Type: application/json" -d "{
    \"user_id\": \"$USER1_ID\",
    \"percentage\": 60.00
}")
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/payer-groups/$PAYER_GROUP_ID/members" -H "Content-Type: application/json" -d "{
    \"user_id\": \"$USER1_ID\",
    \"percentage\": 60.00
}")
test_response $status_code 201 "Adicionar João ao grupo (60%)"
show_response "$response"

log "Adicionando Maria ao grupo (40%)"
response=$(curl -s -X POST "$BASE_URL/payer-groups/$PAYER_GROUP_ID/members" -H "Content-Type: application/json" -d "{
    \"user_id\": \"$USER2_ID\",
    \"percentage\": 40.00
}")
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/payer-groups/$PAYER_GROUP_ID/members" -H "Content-Type: application/json" -d "{
    \"user_id\": \"$USER2_ID\",
    \"percentage\": 40.00
}")
test_response $status_code 201 "Adicionar Maria ao grupo (40%)"
show_response "$response"

log "Listando membros do grupo"
response=$(curl -s -X GET "$BASE_URL/payer-groups/$PAYER_GROUP_ID/members")
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/payer-groups/$PAYER_GROUP_ID/members")
test_response $status_code 200 "Listar membros do grupo"
show_response "$response"

section "3. CENTRO DE CUSTO"
log "Criando centro de custo (Moradia)"
response=$(curl -s -X POST $BASE_URL/finance-cc -H "Content-Type: application/json" -d '{
    "name": "Moradia"
}')
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/finance-cc -H "Content-Type: application/json" -d '{
    "name": "Moradia"
}')
test_response $status_code 201 "Criar CC Moradia"
FINANCE_CC_ID=$(echo $response | jq -r '.id')
show_response "$response"

log "Listando centros de custo"
response=$(curl -s -X GET $BASE_URL/finance-cc)
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X GET $BASE_URL/finance-cc)
test_response $status_code 200 "Listar centros de custo"
show_response "$response"

section "4. MOEDAS"
log "Criando moeda (Real)"
response=$(curl -s -X POST $BASE_URL/currencies -H "Content-Type: application/json" -d '{
    "name": "Real",
    "symbol": "R$",
    "value": 1.0000
}')
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/currencies -H "Content-Type: application/json" -d '{
    "name": "Real",
    "symbol": "R$",
    "value": 1.0000
}')
test_response $status_code 201 "Criar moeda Real"
CURRENCY_ID=$(echo $response | jq -r '.id')
show_response "$response"

log "Listando moedas"
response=$(curl -s -X GET $BASE_URL/currencies)
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X GET $BASE_URL/currencies)
test_response $status_code 200 "Listar moedas"
show_response "$response"

section "5. TAREFAS"
log "Criando tarefa (Limpar Casa)"
response=$(curl -s -X POST $BASE_URL/tasks -H "Content-Type: application/json" -d "{
    \"title\": \"Limpar Casa\",
    \"description\": \"Limpeza semanal\",
    \"start_date\": \"$(date -d 'last month' '+%Y-%m-%d')T00:00:00Z\",
    \"recurrence_cron\": \"0 8 * * 1\",
    \"subtasks\": \"[{\\\"title\\\": \\\"Varrer chão\\\", \\\"done\\\": false}, {\\\"title\\\": \\\"Lavar louça\\\", \\\"done\\\": false}]\",
    \"user_id\": \"$USER1_ID\",
    \"payer_group_id\": \"$PAYER_GROUP_ID\"
}")
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/tasks -H "Content-Type: application/json" -d "{
    \"title\": \"Limpar Casa\",
    \"description\": \"Limpeza semanal\",
    \"start_date\": \"$(date -d 'last month' '+%Y-%m-%d')T00:00:00Z\",
    \"recurrence_cron\": \"0 8 * * 1\",
    \"subtasks\": \"[{\\\"title\\\": \\\"Varrer chão\\\", \\\"done\\\": false}, {\\\"title\\\": \\\"Lavar louça\\\", \\\"done\\\": false}]\",
    \"user_id\": \"$USER1_ID\",
    \"payer_group_id\": \"$PAYER_GROUP_ID\"
}")
test_response $status_code 201 "Criar tarefa Limpar Casa"
TASK_ID=$(echo $response | jq -r '.id')
show_response "$response"

log "Gerando ocorrências da tarefa"
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/tasks/update-occurrences")
test_response $status_code 200 "Gerar ocorrências da tarefa"

log "Listando ocorrências de tarefas"
response=$(curl -s -X GET "$BASE_URL/task-occurrences")
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/task-occurrences")
test_response $status_code 200 "Listar ocorrências de tarefas"
show_response "$response"

log "Selecionando primeira ocorrência da tarefa"
TASK_OCCURRENCE_ID=$(echo $response | jq -r '.[0].id')

log "Atualizando status da ocorrência da tarefa"
response=$(curl -s -X PUT "$BASE_URL/task-occurrences/$TASK_OCCURRENCE_ID" -H "Content-Type: application/json" -d '{
    "status": true
}')
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X PUT "$BASE_URL/task-occurrences/$TASK_OCCURRENCE_ID" -H "Content-Type: application/json" -d '{
    "status": true
}')
test_response $status_code 200 "Atualizar status da ocorrência da tarefa"
show_response "$response"

section "6. FINANÇAS"
log "Criando finança (Aluguel)"
response=$(curl -s -X POST $BASE_URL/finances -H "Content-Type: application/json" -d "{
    \"title\": \"Aluguel\",
    \"description\": \"Pagamento mensal\",
    \"type\": true,
    \"start_date\": \"$(date -d 'last month' '+%Y-%m-%d')T00:00:00Z\",
    \"end_date\": null,
    \"recurrence_days\": 30,
    \"amount\": 1500.00,
    \"user_id\": \"$USER1_ID\",
    \"payer_group_id\": \"$PAYER_GROUP_ID\",
    \"finance_cc_id\": \"$FINANCE_CC_ID\",
    \"currency_id\": \"$CURRENCY_ID\"
}")
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/finances -H "Content-Type: application/json" -d "{
    \"title\": \"Aluguel\",
    \"description\": \"Pagamento mensal\",
    \"type\": true,
    \"start_date\": \"$(date -d 'last month' '+%Y-%m-%d')T00:00:00Z\",
    \"end_date\": null,
    \"recurrence_days\": 30,
    \"amount\": 1500.00,
    \"user_id\": \"$USER1_ID\",
    \"payer_group_id\": \"$PAYER_GROUP_ID\",
    \"finance_cc_id\": \"$FINANCE_CC_ID\",
    \"currency_id\": \"$CURRENCY_ID\"
}")
test_response $status_code 201 "Criar finança Aluguel"
FINANCE_ID=$(echo $response | jq -r '.id')
show_response "$response"

log "Gerando ocorrências da finança"
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/finances/update-occurrences")
test_response $status_code 200 "Gerar ocorrências da finança"

log "Listando ocorrências financeiras"
response=$(curl -s -X GET "$BASE_URL/finance-occurrences")
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/finance-occurrences")
test_response $status_code 200 "Listar ocorrências financeiras"
show_response "$response"

log "Selecionando primeira ocorrência da finança"
FINANCE_OCCURRENCE_ID=$(echo $response | jq -r '.[0].id')

log "Atualizando status da ocorrência financeira para pago"
response=$(curl -s -X PUT "$BASE_URL/finance-occurrences/$FINANCE_OCCURRENCE_ID" -H "Content-Type: application/json" -d '{
    "amount": 1500.00,
    "status": true
}')
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X PUT "$BASE_URL/finance-occurrences/$FINANCE_OCCURRENCE_ID" -H "Content-Type: application/json" -d '{
    "amount": 1500.00,
    "status": true
}')
test_response $status_code 200 "Atualizar status da ocorrência financeira"
show_response "$response"

section "7. DASHBOARD"
log "Verificando dashboard de ocorrências"
response=$(curl -s -X GET "$BASE_URL/occurrences/dashboard")
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/occurrences/dashboard")
test_response $status_code 200 "Verificar dashboard de ocorrências"
show_response "$response"

section "8. CARTEIRAS"
log "Verificando carteira de João (esperado: -900.00, 60% de 1500.00 de despesa)"
response=$(curl -s -X GET "$BASE_URL/wallets/$USER1_ID")
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/wallets/$USER1_ID")
test_response $status_code 200 "Verificar carteira de João"
show_response "$response"

log "Verificando carteira de Maria (esperado: -600.00, 40% de 1500.00 de despesa)"
response=$(curl -s -X GET "$BASE_URL/wallets/$USER2_ID")
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/wallets/$USER2_ID")
test_response $status_code 200 "Verificar carteira de Maria"
show_response "$response"

section "9. TRANSAÇÕES"
log "Verificando transações da ocorrência financeira"
response=$(curl -s -X GET "$BASE_URL/transactions/$FINANCE_OCCURRENCE_ID")
status_code=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/transactions/$FINANCE_OCCURRENCE_ID")
test_response $status_code 200 "Verificar transações da ocorrência financeira"
show_response "$response"

section "TESTES CONCLUÍDOS"
log "Todos os testes foram executados. Verifique os resultados acima." 