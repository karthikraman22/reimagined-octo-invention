# create databases
CREATE DATABASE IF NOT EXISTS `coredb`;

# create root user and grant rights
GRANT ALL ON *.* TO 'root'@'%';


DROP TABLE IF EXISTS gl_acct_jrnl_entry;
DROP TABLE IF EXISTS gl_acct;


CREATE table gl_acct (
  id UUID DEFAULT uuid_generate_v4 (),
  parent_id UUID DEFAULT NULL,
  name VARCHAR(48) DEFAULT NULL,
  gl_code VARCHAR(48) NOT NULL,
  disabled BOOL,
  description VARCHAR(256) DEFAULT NULL,
  allow_manual_entries BOOL,
  type VARCHAR(8) NOT NULL,
  usage VARCHAR(8) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE (gl_code),
  CONSTRAINT fk_ga_parent_id FOREIGN KEY (parent_id) REFERENCES gl_acct(id)
);

CREATE TABLE gl_acct_jrnl_entry (
  id UUID DEFAULT uuid_generate_v4 (),
  tranaction_id VARCHAR(50) NOT NULL,
  ext_ref VARCHAR(50) DEFAULT NULL,
  account_id UUID NOT NULL,
  entry_date DATE NOT NULL,
  amount DECIMAL(19,6) NOT NULL,
  description VARCHAR(256),
  reversal_id UUID DEFAULT NULL,
  manual_entry BOOL,
  entry_type INTEGER NOT NULL,
  tsignature BYTEA NOT NULL,  
  PRIMARY KEY (id),
  CONSTRAINT fk_gaje_account_id FOREIGN KEY (account_id) REFERENCES gl_acct(id),
  CONSTRAINT fk_gaje_reversal_id FOREIGN KEY (reversal_id) REFERENCES gl_acct_jrnl_entry(id)
);

CREATE INDEX idx_gaje_tranaction_id ON gl_acct_jrnl_entry(tranaction_id);
CREATE INDEX idx_gaje_account_id ON gl_acct_jrnl_entry(account_id);
CREATE INDEX idx_gaje_entry_date ON gl_acct_jrnl_entry(entry_date);
