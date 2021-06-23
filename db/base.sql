# create databases
CREATE DATABASE IF NOT EXISTS `coredb`;

# create root user and grant rights
GRANT ALL ON *.* TO 'root'@'%';


DROP TABLE IF EXISTS `gl_account`;

-- DDL for accounting sub system related tables
CREATE TABLE `gl_account` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `parent_id` BIGINT DEFAULT NULL,
  `gl_code` varchar(45) NOT NULL,
  `disabled` tinyint NOT NULL DEFAULT '0',
  `manual_journal_entries_allowed` tinyint NOT NULL DEFAULT '1',
  `account_usage` tinyint NOT NULL DEFAULT '2',
  `classification_enum` SMALLINT NOT NULL,
  `description` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `acc_gl_code` (`gl_code`),
  KEY `FK_ACC_0000000001` (`parent_id`),
  CONSTRAINT `FK_ACC_0000000001` FOREIGN KEY (`parent_id`) REFERENCES `acc_gl_account` (`id`)
) 


