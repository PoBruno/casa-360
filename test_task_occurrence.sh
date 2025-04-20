#!/bin/bash

BASE_URL="http://localhost:3001"

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Função para imprimir mensagens com timestamp
log() {
    echo -e "[$(date +"%Y-%m-%d %H:%M:%S")] $1"
}

# Cria um usuário
log "Criando usuário de teste..."
USER_RESPONSE=$(curl -s -X POST "$BASE_URL/users" -H "Content-Type: application/json" -d '{
    "name": "Usuário Teste"
}')
USER_ID=$(echo $USER_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$USER_ID" ]; then
    echo -e "${RED}✗ Erro: Falha ao criar usuário${NC}"
    exit 1
else
    echo -e "${GREEN}✓ Sucesso: Usuário criado com ID: $USER_ID${NC}"
fi

# Cria um grupo de pagadores
log "Criando grupo de pagadores..."
GROUP_RESPONSE=$(curl -s -X POST "$BASE_URL/payer-groups" -H "Content-Type: application/json" -d '{
    "name": "Grupo Teste"
}')
GROUP_ID=$(echo $GROUP_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$GROUP_ID" ]; then
    echo -e "${RED}✗ Erro: Falha ao criar grupo${NC}"
    exit 1
else
    echo -e "${GREEN}✓ Sucesso: Grupo criado com ID: $GROUP_ID${NC}"
fi

# Cria uma tarefa
log "Criando tarefa de teste..."
TASK_RESPONSE=$(curl -s -X POST "$BASE_URL/tasks" -H "Content-Type: application/json" -d '{
    "title": "Tarefa Teste",
    "description": "Descrição da tarefa teste",
    "start_date": "2023-01-01",
    "recurrence_cron": "0 0 1 * *",
    "subtasks": [{"title": "Subtarefa 1", "done": false}],
    "user_id": "'$USER_ID'",
    "payer_group_id": "'$GROUP_ID'"
}')
TASK_ID=$(echo $TASK_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$TASK_ID" ]; then
    echo -e "${RED}✗ Erro: Falha ao criar tarefa${NC}"
    exit 1
else
    echo -e "${GREEN}✓ Sucesso: Tarefa criada com ID: $TASK_ID${NC}"
fi

# Gera ocorrências
log "Gerando ocorrências de tarefas..."
curl -s -X POST "$BASE_URL/tasks/update-occurrences" > /dev/null

# Busca as ocorrências
log "Buscando ocorrências de tarefas..."
OCCURRENCE_RESPONSE=$(curl -s -X GET "$BASE_URL/task-occurrences")
TASK_OCCURRENCE_ID=$(echo $OCCURRENCE_RESPONSE | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)

if [ -z "$TASK_OCCURRENCE_ID" ]; then
    echo -e "${RED}✗ Erro: Falha ao encontrar ocorrências${NC}"
    exit 1
else
    echo -e "${GREEN}✓ Sucesso: Encontrada ocorrência com ID: $TASK_OCCURRENCE_ID${NC}"
fi

# Atualiza o status da ocorrência
log "Atualizando status da ocorrência..."
UPDATE_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X PUT "$BASE_URL/task-occurrences/$TASK_OCCURRENCE_ID" -H "Content-Type: application/json" -d '{
    "status": true
}')

HTTP_STATUS=$(echo "$UPDATE_RESPONSE" | grep "HTTP_STATUS:" | cut -d':' -f2)
RESPONSE_BODY=$(echo "$UPDATE_RESPONSE" | sed '/HTTP_STATUS:/d')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo -e "${GREEN}✓ Sucesso: Status da ocorrência atualizado${NC}"
    echo "Resposta: $RESPONSE_BODY"
else
    echo -e "${RED}✗ Erro ($HTTP_STATUS): Falha ao atualizar status da ocorrência${NC}"
    echo "Resposta: $RESPONSE_BODY"
fi

echo -e "Teste concluído!" 