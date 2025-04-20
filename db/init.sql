-- Extensão para UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tabela de usuários
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL
);

-- Tabela de grupos de pagadores
CREATE TABLE payer_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL
);

-- Tabela de membros dos grupos de pagadores
CREATE TABLE payer_group_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payer_group_id UUID NOT NULL REFERENCES payer_groups(id),
    user_id UUID NOT NULL REFERENCES users(id),
    percentage DECIMAL(5,2) NOT NULL CHECK (percentage > 0 AND percentage <= 100),
    UNIQUE(payer_group_id, user_id)
);

-- Trigger para garantir que a soma dos percentuais seja 100%
CREATE OR REPLACE FUNCTION check_percentage_sum()
RETURNS TRIGGER AS $$
BEGIN
    IF (SELECT SUM(percentage) FROM payer_group_members WHERE payer_group_id = NEW.payer_group_id) > 100 THEN
        -- Escapamos '%%' para que o '%' seja interpretado como literal
        RAISE EXCEPTION 'A soma dos percentuais não pode exceder 100%%';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_percentage_sum_trigger
AFTER INSERT OR UPDATE ON payer_group_members
FOR EACH ROW
EXECUTE FUNCTION check_percentage_sum();

-- Tabela de centro de custo
CREATE TABLE finance_cc (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    parent_id UUID REFERENCES finance_cc(id)
);

-- Tabela de moedas
CREATE TABLE finance_currency (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    symbol TEXT NOT NULL,
    value DECIMAL(10,4) NOT NULL DEFAULT 1.0
);

-- Tabela de tarefas recorrentes
CREATE TABLE task_installments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,
    recurrence_cron TEXT NOT NULL,
    subtasks JSONB,
    user_id UUID REFERENCES users(id),
    payer_group_id UUID REFERENCES payer_groups(id)
);

-- Tabela de ocorrências de tarefas
CREATE TABLE task_occurrences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL REFERENCES task_installments(id),
    date DATE NOT NULL,
    status BOOLEAN DEFAULT false,
    user_id UUID REFERENCES users(id),
    payer_group_id UUID REFERENCES payer_groups(id),
    subtasks JSONB,
    UNIQUE(task_id, date)
);

-- Tabela de finanças recorrentes
CREATE TABLE finance_installments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT,
    type BOOLEAN NOT NULL, -- false = receita, true = despesa
    start_date DATE NOT NULL,
    end_date DATE,
    recurrence_days INTEGER NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    user_id UUID REFERENCES users(id),
    payer_group_id UUID REFERENCES payer_groups(id),
    finance_cc_id UUID REFERENCES finance_cc(id),
    currency_id UUID REFERENCES finance_currency(id)
);

-- Tabela de ocorrências financeiras
CREATE TABLE finance_occurrences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    finance_id UUID NOT NULL REFERENCES finance_installments(id),
    date DATE NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    status BOOLEAN DEFAULT false,
    UNIQUE(finance_id, date)
);

-- Nova tabela de carteiras financeiras
CREATE TABLE finance_wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, created_at)
);

-- Nova tabela de transações
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    finance_occurrence_id UUID NOT NULL REFERENCES finance_occurrences(id),
    amount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(finance_occurrence_id)
);

-- Função para processar transação e atualizar carteiras
CREATE OR REPLACE FUNCTION process_finance_occurrence()
RETURNS TRIGGER AS $$
DECLARE
    v_finance_installment finance_installments%ROWTYPE;
    v_currency finance_currency%ROWTYPE;
    v_transaction_amount DECIMAL(10,2);
    v_member RECORD;
    v_last_wallet finance_wallets%ROWTYPE;
    v_new_amount DECIMAL(10,2);
BEGIN
    -- Só processa quando o status muda para true
    IF NEW.status = true AND (TG_OP = 'INSERT' OR OLD.status = false) THEN
        -- Busca informações da finança
        SELECT * INTO v_finance_installment 
        FROM finance_installments 
        WHERE id = NEW.finance_id;

        -- Busca informações da moeda
        SELECT * INTO v_currency 
        FROM finance_currency 
        WHERE id = v_finance_installment.currency_id;

        -- Calcula o valor da transação considerando a taxa de câmbio
        v_transaction_amount := NEW.amount * v_currency.value;

        -- Registra a transação
        INSERT INTO transactions (finance_occurrence_id, amount)
        VALUES (NEW.id, v_transaction_amount);

        -- Para cada membro do grupo de pagadores, atualiza sua carteira
        FOR v_member IN 
            SELECT pgm.user_id, pgm.percentage 
            FROM payer_group_members pgm
            WHERE pgm.payer_group_id = v_finance_installment.payer_group_id
        LOOP
            -- Busca o último valor da carteira do usuário
            SELECT * INTO v_last_wallet 
            FROM finance_wallets 
            WHERE user_id = v_member.user_id 
            ORDER BY created_at DESC 
            LIMIT 1;

            -- Se não existir carteira, considera valor inicial 0
            IF v_last_wallet.amount IS NULL THEN
                v_new_amount := 0;
            ELSE
                v_new_amount := v_last_wallet.amount;
            END IF;

            -- Calcula o novo valor baseado no tipo (receita/despesa) e porcentagem
            IF v_finance_installment.type = false THEN -- Receita
                v_new_amount := v_new_amount + (v_transaction_amount * v_member.percentage / 100);
            ELSE -- Despesa
                v_new_amount := v_new_amount - (v_transaction_amount * v_member.percentage / 100);
            END IF;

            -- Insere novo registro na carteira
            INSERT INTO finance_wallets (user_id, amount)
            VALUES (v_member.user_id, v_new_amount);
        END LOOP;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para processar transações e atualizar carteiras
CREATE TRIGGER process_finance_occurrence_trigger
AFTER INSERT OR UPDATE ON finance_occurrences
FOR EACH ROW
EXECUTE FUNCTION process_finance_occurrence();

-- View para dashboard de ocorrências
CREATE OR REPLACE VIEW occurrences_dashboard AS
SELECT 
    'finance' as occurrence_type,
    fo.id,
    fo.date,
    fo.status,
    fi.title,
    fi.description,
    fi.type as finance_type,
    fo.amount,
    fc.symbol as currency_symbol,
    fc.value as currency_value,
    (fo.amount * fc.value) as amount_converted,
    fcc.name as cost_center,
    pg.name as payer_group,
    u.name as responsible_user
FROM 
    finance_occurrences fo
    INNER JOIN finance_installments fi ON fo.finance_id = fi.id
    LEFT JOIN finance_currency fc ON fi.currency_id = fc.id
    LEFT JOIN finance_cc fcc ON fi.finance_cc_id = fcc.id
    LEFT JOIN payer_groups pg ON fi.payer_group_id = pg.id
    LEFT JOIN users u ON fi.user_id = u.id
UNION ALL
SELECT 
    'task' as occurrence_type,
    to2.id,
    to2.date,
    to2.status,
    ti.title,
    ti.description,
    null as finance_type,
    null as amount,
    null as currency_symbol,
    null as currency_value,
    null as amount_converted,
    null as cost_center,
    pg.name as payer_group,
    u.name as responsible_user
FROM 
    task_occurrences to2
    INNER JOIN task_installments ti ON to2.task_id = ti.id
    LEFT JOIN payer_groups pg ON ti.payer_group_id = pg.id
    LEFT JOIN users u ON ti.user_id = u.id;

-- Índices para melhor performance
CREATE INDEX idx_task_occurrences_date ON task_occurrences(date);
CREATE INDEX idx_finance_occurrences_date ON finance_occurrences(date);
CREATE INDEX idx_task_installments_start_date ON task_installments(start_date);
CREATE INDEX idx_finance_installments_start_date ON finance_installments(start_date);
CREATE INDEX idx_finance_installments_end_date ON finance_installments(end_date);
CREATE INDEX idx_finance_wallets_user_created ON finance_wallets(user_id, created_at);
CREATE INDEX idx_transactions_created ON transactions(created_at);


