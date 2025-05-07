DROP DATABASE IF EXISTS postgres;
-- DROP DATABASE IF EXISTS kanbanify;

CREATE TYPE status AS ENUM ('backlog', 'inProgress', 'done');
CREATE TYPE variant AS ENUM ('bug', 'chore', 'task');

CREATE TABLE Issues (
  id SERIAL PRIMARY KEY,
  kanban_id VARCHAR(10) NOT NULL,
  title VARCHAR(255) NOT NULL,
  description VARCHAR(4000) NOT NULL,
  status status NOT NULL DEFAULT 'backlog',
  variant variant NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- INSERT INTO Issues (kanban_id, title, description, variant)
-- VALUES ('KAN-1', 'Test', 'Dette er en test creation from the SQL script', 'bug');