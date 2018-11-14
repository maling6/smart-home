-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION pgcrypto;
SELECT gen_random_uuid();
create type flows_status as enum ('enabled', 'disabled');

CREATE TABLE flows (
  id                   bigserial                not null constraint flows_pkey primary key,
  name                 VarChar(255)             NOT NULL,
  description          Text                     NULL,
  status               flows_status             NOT NULL DEFAULT 'enabled',
  workflow_id          BIGINT CONSTRAINT flows_2_workflows_fk REFERENCES workflows (id) on update cascade on delete cascade,
  workflow_scenario_id BIGINT CONSTRAINT flows_2_workflow_scenarios_fk REFERENCES workflow_scenarios (id) on update cascade on delete cascade,
  created_at           timestamp with time zone,
  updated_at           timestamp with time zone null
);

create type flow_elements_status as enum ('enabled', 'disabled');
create type flow_elements_prototype_type as enum ('default', 'MessageHandler', 'MessageEmitter', 'Task', 'Gateway', 'Flow');

CREATE TABLE flow_elements (
  uuid           UUID PRIMARY KEY                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             DEFAULT gen_random_uuid(),
  name           VarChar(255)         NOT NULL,
  description    Text                 NULL,
  status         flow_elements_status not null DEFAULT 'enabled',
  graph_settings JSONB default '{}',
  flow_link      numeric              NULL,
  prototype_type flow_elements_prototype_type not null                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         default 'default',
  flow_id        BIGINT CONSTRAINT flow_elements_2_flows_fk REFERENCES flows (id) on update cascade on delete cascade     NULL,
  script_id      BIGINT CONSTRAINT flow_elements_2_scripts_fk REFERENCES scripts (id) on update cascade on delete cascade NULL,
  created_at     timestamp with time zone,
  updated_at     timestamp with time zone                                                                                 NULL
);

CREATE UNIQUE INDEX uuid_2_flow_elements_unq
  ON flow_elements (uuid);

CREATE TABLE connections (
  uuid           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name           VarChar(255)             NULL,
  element_from   UUID CONSTRAINT connections_from_2_flow_elements_fk REFERENCES flow_elements (uuid) on update cascade on delete cascade,
  element_to     UUID CONSTRAINT connections_to_2_flow_elements_fk REFERENCES flow_elements (uuid) on update cascade on delete cascade,
  flow_id        BIGINT CONSTRAINT connections_2_flows_fk REFERENCES flows (id) on update cascade on delete cascade,
  graph_settings JSONB            default '{}',
  point_from     smallint                 NOT NULL,
  point_to       smallint                 NOT NULL,
  direction      VarChar(255)             NOT NULL,
  created_at     timestamp with time zone,
  updated_at     timestamp with time zone null
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS connections CASCADE;
DROP TABLE IF EXISTS flow_elements CASCADE;
DROP TABLE IF EXISTS flows CASCADE;
drop extension pgcrypto;
drop type flows_status;
drop type flow_elements_status;
drop type flow_elements_prototype_type;