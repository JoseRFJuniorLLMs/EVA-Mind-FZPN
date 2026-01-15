-- Migration v26: Computer Use Agent
-- Sistema de automação web com aprovação humana

-- Tabela de tarefas de automação
CREATE TABLE IF NOT EXISTS automation_tasks (
    id BIGSERIAL PRIMARY KEY,
    idoso_id BIGINT NOT NULL REFERENCES idosos(id) ON DELETE CASCADE,
    task_type VARCHAR(50) NOT NULL CHECK (task_type IN (
        'buy_medication',
        'schedule_appointment',
        'order_food',
        'request_ride',
        'other'
    )),
    service_name VARCHAR(100) NOT NULL, -- 'Drogasil', 'Doctoralia', 'iFood', 'Uber'
    task_params JSONB NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN (
        'pending',           -- Aguardando aprovação
        'approved',          -- Aprovado, aguardando execução
        'executing',         -- Em execução
        'completed',         -- Concluído com sucesso
        'failed',            -- Falhou
        'cancelled',         -- Cancelado
        'requires_approval'  -- Requer aprovação adicional
    )),
    approval_required BOOLEAN DEFAULT true,
    approved_by BIGINT REFERENCES usuarios(id),
    approved_at TIMESTAMP,
    execution_log JSONB,
    screenshots JSONB, -- Array de URLs de screenshots
    result JSONB,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    executed_at TIMESTAMP,
    completed_at TIMESTAMP
);

-- Tabela de aprovações pendentes
CREATE TABLE IF NOT EXISTS automation_approvals (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL REFERENCES automation_tasks(id) ON DELETE CASCADE,
    approver_id BIGINT NOT NULL REFERENCES usuarios(id),
    approval_status VARCHAR(20) CHECK (approval_status IN ('pending', 'approved', 'rejected')),
    approval_reason TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    responded_at TIMESTAMP
);

-- Tabela de histórico de execução
CREATE TABLE IF NOT EXISTS automation_execution_log (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL REFERENCES automation_tasks(id) ON DELETE CASCADE,
    step_number INT NOT NULL,
    step_name VARCHAR(200) NOT NULL,
    step_status VARCHAR(50),
    screenshot_url TEXT,
    step_data JSONB,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Índices
CREATE INDEX idx_automation_tasks_idoso ON automation_tasks(idoso_id);
CREATE INDEX idx_automation_tasks_status ON automation_tasks(status);
CREATE INDEX idx_automation_tasks_type ON automation_tasks(task_type);
CREATE INDEX idx_automation_tasks_created ON automation_tasks(created_at DESC);
CREATE INDEX idx_automation_approvals_task ON automation_approvals(task_id);
CREATE INDEX idx_automation_approvals_approver ON automation_approvals(approver_id);
CREATE INDEX idx_automation_execution_task ON automation_execution_log(task_id);

-- View: Tarefas pendentes de aprovação
CREATE OR REPLACE VIEW v_pending_approvals AS
SELECT 
    t.id as task_id,
    t.idoso_id,
    i.nome as idoso_nome,
    t.task_type,
    t.service_name,
    t.task_params,
    t.created_at,
    EXTRACT(EPOCH FROM (NOW() - t.created_at))/60 as minutes_waiting
FROM automation_tasks t
JOIN idosos i ON i.id = t.idoso_id
WHERE t.status = 'pending'
  AND t.approval_required = true
ORDER BY t.created_at ASC;

-- View: Histórico de automações
CREATE OR REPLACE VIEW v_automation_history AS
SELECT 
    t.id,
    t.idoso_id,
    i.nome as idoso_nome,
    t.task_type,
    t.service_name,
    t.status,
    t.created_at,
    t.completed_at,
    EXTRACT(EPOCH FROM (t.completed_at - t.created_at)) as duration_seconds,
    t.result
FROM automation_tasks t
JOIN idosos i ON i.id = t.idoso_id
WHERE t.status IN ('completed', 'failed', 'cancelled')
ORDER BY t.created_at DESC;

-- View: Estatísticas de automação
CREATE OR REPLACE VIEW v_automation_stats AS
SELECT 
    task_type,
    service_name,
    COUNT(*) as total_tasks,
    COUNT(*) FILTER (WHERE status = 'completed') as completed,
    COUNT(*) FILTER (WHERE status = 'failed') as failed,
    COUNT(*) FILTER (WHERE status = 'cancelled') as cancelled,
    COUNT(*) FILTER (WHERE status = 'pending') as pending,
    AVG(EXTRACT(EPOCH FROM (completed_at - created_at))) FILTER (WHERE status = 'completed') as avg_duration_seconds,
    ROUND(100.0 * COUNT(*) FILTER (WHERE status = 'completed') / NULLIF(COUNT(*), 0), 2) as success_rate
FROM automation_tasks
GROUP BY task_type, service_name
ORDER BY total_tasks DESC;

-- Função: Criar tarefa de automação
CREATE OR REPLACE FUNCTION create_automation_task(
    p_idoso_id BIGINT,
    p_task_type VARCHAR(50),
    p_service_name VARCHAR(100),
    p_task_params JSONB,
    p_approval_required BOOLEAN DEFAULT true
) RETURNS BIGINT AS $$
DECLARE
    v_task_id BIGINT;
BEGIN
    INSERT INTO automation_tasks (
        idoso_id,
        task_type,
        service_name,
        task_params,
        approval_required,
        status
    ) VALUES (
        p_idoso_id,
        p_task_type,
        p_service_name,
        p_task_params,
        p_approval_required,
        CASE WHEN p_approval_required THEN 'pending' ELSE 'approved' END
    ) RETURNING id INTO v_task_id;
    
    RETURN v_task_id;
END;
$$ LANGUAGE plpgsql;

-- Função: Aprovar tarefa
CREATE OR REPLACE FUNCTION approve_automation_task(
    p_task_id BIGINT,
    p_approver_id BIGINT
) RETURNS VOID AS $$
BEGIN
    UPDATE automation_tasks
    SET status = 'approved',
        approved_by = p_approver_id,
        approved_at = NOW(),
        updated_at = NOW()
    WHERE id = p_task_id
      AND status = 'pending';
      
    INSERT INTO automation_approvals (task_id, approver_id, approval_status, responded_at)
    VALUES (p_task_id, p_approver_id, 'approved', NOW());
END;
$$ LANGUAGE plpgsql;

-- Função: Rejeitar tarefa
CREATE OR REPLACE FUNCTION reject_automation_task(
    p_task_id BIGINT,
    p_approver_id BIGINT,
    p_reason TEXT
) RETURNS VOID AS $$
BEGIN
    UPDATE automation_tasks
    SET status = 'cancelled',
        updated_at = NOW()
    WHERE id = p_task_id
      AND status = 'pending';
      
    INSERT INTO automation_approvals (task_id, approver_id, approval_status, approval_reason, responded_at)
    VALUES (p_task_id, p_approver_id, 'rejected', p_reason, NOW());
END;
$$ LANGUAGE plpgsql;

-- Função: Registrar passo de execução
CREATE OR REPLACE FUNCTION log_automation_step(
    p_task_id BIGINT,
    p_step_number INT,
    p_step_name VARCHAR(200),
    p_step_status VARCHAR(50),
    p_screenshot_url TEXT DEFAULT NULL,
    p_step_data JSONB DEFAULT NULL,
    p_error_message TEXT DEFAULT NULL
) RETURNS VOID AS $$
BEGIN
    INSERT INTO automation_execution_log (
        task_id,
        step_number,
        step_name,
        step_status,
        screenshot_url,
        step_data,
        error_message
    ) VALUES (
        p_task_id,
        p_step_number,
        p_step_name,
        p_step_status,
        p_screenshot_url,
        p_step_data,
        p_error_message
    );
END;
$$ LANGUAGE plpgsql;

-- Comentários
COMMENT ON TABLE automation_tasks IS 'Tarefas de automação web (Computer Use)';
COMMENT ON TABLE automation_approvals IS 'Aprovações de tarefas de automação';
COMMENT ON TABLE automation_execution_log IS 'Log detalhado de execução de automações';
COMMENT ON VIEW v_pending_approvals IS 'Tarefas aguardando aprovação humana';
COMMENT ON VIEW v_automation_history IS 'Histórico de automações executadas';
COMMENT ON VIEW v_automation_stats IS 'Estatísticas de performance de automações';
