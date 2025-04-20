#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# URL base da API
BASE_URL="http://localhost:3001"

# Função para log
log() {
    echo -e "${YELLOW}[$(date '+%Y-%m-%d %H:%M:%S')] $1${NC}"
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

# Função para validar valor da carteira
validate_wallet() {
    local user_id=$1
    local expected_amount=$2
    local description=$3
    
    response=$(curl -s -w "%{http_code}" -X GET "$BASE_URL/wallets/$user_id")
    status_code=${response: -3}
    response_body=${response:0:${#response}-3}
    
    if [ $status_code -eq 200 ]; then
        actual_amount=$(echo $response_body | jq -r '.amount')
        if [ $(echo "$actual_amount == $expected_amount" | bc) -eq 1 ]; then
            echo -e "${GREEN}✓ Sucesso: Carteira $description - Valor correto: $actual_amount${NC}"
        else
            echo -e "${RED}✗ Erro: Carteira $description - Valor esperado: $expected_amount, Valor atual: $actual_amount${NC}"
        fi
    else
        echo -e "${RED}✗ Erro ($status_code): Falha ao buscar carteira $description${NC}"
    fi
}

# Variáveis para armazenar IDs
USER1_ID=""
USER2_ID=""
PAYER_GROUP_ID=""
PAYER_GROUP_MEMBER_ID=""
FINANCE_CC_ID=""
CURRENCY_ID=""
TASK_ID=""
TASK_OCCURRENCE_ID=""
FINANCE_ID=""
FINANCE_OCCURRENCE_ID=""

log "Iniciando testes da API Casa360"

# 1. Criar usuários
log "Testando criação de usuários"

response=$(curl -s -w "%{http_code}" -X POST $BASE_URL/users -H "Content-Type: application/json" -d '{
    "name": "João"
}')
status_code=${response: -3}
response_body=${response:0:${#response}-3}
test_response $status_code 201 "Criar usuário João"
USER1_ID=$(echo $response_body | jq -r '.id')

response=$(curl -s -w "%{http_code}" -X POST $BASE_URL/users -H "Content-Type: application/json" -d '{
    "name": "Maria"
}')
status_code=${response: -3}
response_body=${response:0:${#response}-3}
test_response $status_code 201 "Criar usuário Maria"
USER2_ID=$(echo $response_body | jq -r '.id')

# 2. Criar grupo de pagadores
log "Testando criação de grupo de pagadores"

response=$(curl -s -w "%{http_code}" -X POST $BASE_URL/payer-groups -H "Content-Type: application/json" -d '{
    "name": "Casa"
}')
status_code=${response: -3}
response_body=${response:0:${#response}-3}
test_response $status_code 201 "Criar grupo Casa"
PAYER_GROUP_ID=$(echo $response_body | jq -r '.id')

# 3. Adicionar membros ao grupo
log "Testando adição de membros ao grupo"

response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/payer-groups/$PAYER_GROUP_ID/members" -H "Content-Type: application/json" -d "{
    \"user_id\": \"$USER1_ID\",
    \"percentage\": 60.00
}")
status_code=${response: -3}
test_response $status_code 201 "Adicionar João ao grupo (60%)"

response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/payer-groups/$PAYER_GROUP_ID/members" -H "Content-Type: application/json" -d "{
    \"user_id\": \"$USER2_ID\",
    \"percentage\": 40.00
}")
status_code=${response: -3}
test_response $status_code 201 "Adicionar Maria ao grupo (40%)"

# 4. Criar centro de custo
log "Testando criação de centro de custo"

response=$(curl -s -w "%{http_code}" -X POST $BASE_URL/finance-cc -H "Content-Type: application/json" -d '{
    "name": "Moradia"
}')
status_code=${response: -3}
response_body=${response:0:${#response}-3}
test_response $status_code 201 "Criar CC Moradia"
FINANCE_CC_ID=$(echo $response_body | jq -r '.id')

# 5. Criar moeda
log "Testando criação de moeda"

response=$(curl -s -w "%{http_code}" -X POST $BASE_URL/currencies -H "Content-Type: application/json" -d '{
    "name": "Real",
    "symbol": "R$",
    "value": 1.0000
}')
status_code=${response: -3}
response_body=${response:0:${#response}-3}
test_response $status_code 201 "Criar moeda Real"
CURRENCY_ID=$(echo $response_body | jq -r '.id')

# 6. Criar tarefa
log "Testando criação de tarefa"

response=$(curl -s -w "%{http_code}" -X POST $BASE_URL/tasks -H "Content-Type: application/json" -d "{
    \"title\": \"Limpar Casa\",
    \"description\": \"Limpeza semanal\",
    \"start_date\": \"2025-02-01T00:00:00Z\",
    \"recurrence_cron\": \"0 8 * * 1\",
    \"subtasks\": \"[]\",
    \"user_id\": \"$USER1_ID\",
    \"payer_group_id\": \"$PAYER_GROUP_ID\"
}")
status_code=${response: -3}
response_body=${response:0:${#response}-3}
test_response $status_code 201 "Criar tarefa Limpar Casa"
TASK_ID=$(echo $response_body | jq -r '.id')

# 7. Gerar ocorrências da tarefa
log "Testando geração de ocorrências da tarefa"

response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/tasks/$TASK_ID/occurrences")
status_code=${response: -3}
test_response $status_code 200 "Gerar ocorrências da tarefa"

# 8. Criar finança
log "Testando criação de finança"

response=$(curl -s -w "%{http_code}" -X POST $BASE_URL/finances -H "Content-Type: application/json" -d "{
    \"title\": \"Aluguel\",
    \"description\": \"Aluguel mensal\",
    \"type\": false,
    \"start_date\": \"2025-02-01T00:00:00Z\",
    \"end_date\": null,
    \"recurrence_days\": 30,
    \"amount\": 1000.00,
    \"user_id\": \"$USER1_ID\",
    \"payer_group_id\": \"$PAYER_GROUP_ID\",
    \"finance_cc_id\": \"$FINANCE_CC_ID\",
    \"currency_id\": \"$CURRENCY_ID\"
}")
status_code=${response: -3}
response_body=${response:0:${#response}-3}
test_response $status_code 201 "Criar finança Aluguel"
FINANCE_ID=$(echo $response_body | jq -r '.id')

# 9. Gerar ocorrências da finança
log "Testando geração de ocorrências da finança"

response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/finances/$FINANCE_ID/occurrences")
status_code=${response: -3}
test_response $status_code 200 "Gerar ocorrências da finança"

# 10. Atualizar todas as ocorrências
log "Testando atualização de todas as ocorrências"

response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/tasks/update-occurrences")
status_code=${response: -3}
test_response $status_code 200 "Atualizar ocorrências de tarefas"

response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/finances/update-occurrences")
status_code=${response: -3}
test_response $status_code 200 "Atualizar ocorrências de finanças"

# 11. Buscar e atualizar uma ocorrência de tarefa
log "Testando atualização de ocorrência de tarefa"

response=$(curl -s -w "%{http_code}" -X GET "$BASE_URL/task-occurrences")
status_code=${response: -3}
response_body=${response:0:${#response}-3}
test_response $status_code 200 "Buscar ocorrências de tarefas"
TASK_OCCURRENCE_ID=$(echo $response_body | jq -r '.[0].id')

response=$(curl -s -w "%{http_code}" -X PUT "$BASE_URL/task-occurrences/$TASK_OCCURRENCE_ID" -H "Content-Type: application/json" -d '{
    "status": true
}')
status_code=${response: -3}
test_response $status_code 200 "Atualizar status da ocorrência de tarefa"

# 12. Buscar e atualizar uma ocorrência financeira
log "Testando atualização de ocorrência financeira"

response=$(curl -s -w "%{http_code}" -X GET "$BASE_URL/finance-occurrences")
status_code=${response: -3}
response_body=${response:0:${#response}-3}
test_response $status_code 200 "Buscar ocorrências financeiras"
FINANCE_OCCURRENCE_ID=$(echo $response_body | jq -r '.[0].id')

response=$(curl -s -w "%{http_code}" -X PUT "$BASE_URL/finance-occurrences/$FINANCE_OCCURRENCE_ID" -H "Content-Type: application/json" -d '{
    "amount": 2000.00,
    "status": true
}')
status_code=${response: -3}
test_response $status_code 200 "Atualizar valor e status da ocorrência financeira"

# 13. Validar carteiras
log "Testando valores das carteiras"

# João deve ter 60% de 2000 = 1200
validate_wallet "$USER1_ID" "1200.00" "João"

# Maria deve ter 40% de 2000 = 800
validate_wallet "$USER2_ID" "800.00" "Maria"

# Resumo dos testes
log "Testes concluídos!" 