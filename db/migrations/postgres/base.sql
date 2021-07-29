# create databases
CREATE DATABASE IF NOT EXISTS "coredb";

# create root user and grant rights
GRANT ALL ON *.* TO 'root'@'%';


DROP TABLE IF EXISTS gl_journal;
DROP TABLE IF EXISTS gl_account;
DROP TABLE IF EXISTS gl_organization;


CREATE TABLE gl_organization (
  "id" UUID NOT NULL,
  "code" VARCHAR(18) NOT NULL,
  "name" VARCHAR(48) NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(), 
  "updated_at" TIMESTAMP DEFAULT NULL
);

ALTER TABLE gl_organization ADD CONSTRAINT pk_go_id PRIMARY KEY(id);
ALTER TABLE gl_organization ADD CONSTRAINT uk_go_code UNIQUE (code);


CREATE TABLE gl_account (
  "id" UUID NOT NULL,
  "number" VARCHAR(18) NOT NULL,
  "code" VARCHAR(18) NOT NULL,
  "parent_id" UUID DEFAULT NULL,
  "org_id" UUID NOT NULL,
  "disabled" BOOL,
  "description" VARCHAR(256) DEFAULT NULL,
  "allow_manual_entries" BOOL,
  "type" VARCHAR(8) NOT NULL,
  "usage" VARCHAR(8) NOT NULL,
  "tags" JSONB DEFAULT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(), 
  "updated_at" TIMESTAMP DEFAULT NULL
);

ALTER TABLE gl_account ADD CONSTRAINT pk_ga_id PRIMARY KEY(id);
ALTER TABLE gl_account ADD CONSTRAINT uk_ga_orgid_number UNIQUE ("org_id", "code");

ALTER TABLE gl_account ADD CONSTRAINT fk_ga_parent_id FOREIGN KEY (parent_id) REFERENCES gl_account(id);

ALTER TABLE gl_account ADD CONSTRAINT fk_ga_org_id FOREIGN KEY (org_id) REFERENCES gl_organization(id);



CREATE TABLE gl_journal (
  id UUID NOT NULL,
  tranaction_id VARCHAR(48) NOT NULL,
  ext_ref VARCHAR(48) DEFAULT NULL,
  account_number VARCHAR(18) NOT NULL,
  entry_date DATE NOT NULL,
  amount DECIMAL(19,6) NOT NULL,
  notes VARCHAR(256),
  reversal_id UUID DEFAULT NULL,
  payment_id UUID DEFAULT NULL,
  manual_entry BOOL,
  entry_type CHAR NOT NULL,
  tsignature BYTEA NOT NULL, 
  "tags" JSONB DEFAULT NULL, 
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE gl_journal ADD CONSTRAINT pk_gj_id PRIMARY KEY(id);

ALTER TABLE gl_journal ADD CONSTRAINT fk_gj_reversal_id FOREIGN KEY (reversal_id) REFERENCES gl_journal(id);

CREATE INDEX idx_gj_tranaction_id ON gl_journal(tranaction_id);
CREATE INDEX idx_gj_entry_date_acct_number ON gl_journal(entry_date, account_number);

