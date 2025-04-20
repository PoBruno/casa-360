#!/bin/bash

# URL base da API
BASE_URL="http://localhost:3001"

echo "Buscando ocorrências de tarefas..."
RESPONSE=$(curl -s -X GET "$BASE_URL/task-occurrences")
TASK_OCCURRENCE_ID=$(echo $RESPONSE | jq -r '.[0].id')

if [ -z "$TASK_OCCURRENCE_ID" ]; then
    echo "Erro: Não foi possível encontrar ocorrências de tarefas"
    exit 1
fi

echo "ID da ocorrência: $TASK_OCCURRENCE_ID"

echo "Atualizando status da ocorrência..."
UPDATE_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X PUT "$BASE_URL/task-occurrences/$TASK_OCCURRENCE_ID" -H "Content-Type: application/json" -d '{
    "status": true
}')

HTTP_STATUS=$(echo "$UPDATE_RESPONSE" | grep "HTTP_STATUS:" | cut -d':' -f2)
RESPONSE_BODY=$(echo "$UPDATE_RESPONSE" | sed '/HTTP_STATUS:/d')

echo "Status HTTP: $HTTP_STATUS"
echo "Resposta: $RESPONSE_BODY"

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "Sucesso: Ocorrência atualizada com sucesso!"
else
    echo "Erro: Falha ao atualizar ocorrência"
fi 