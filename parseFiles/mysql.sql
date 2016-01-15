CREATE TABLE `sc_gi_card` (
  phone varchar(20) DEFAULT NULL,
  polno varchar(20) DEFAULT NULL,
  name varchar(20) DEFAULT NULL,
  card_no varchar(20) DEFAULT NULL,
  bank varchar(100) DEFAULT NULL,
  bank_code varchar(20) DEFAULT NULL,
  bank_branch varchar(100) DEFAULT NULL,
  bank_branch_code varchar(20) DEFAULT NULL,
  PRIMARY KEY (polno,name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `sc_gi_cfg_mail` (
  mail_addr varchar(30) DEFAULT NULL,
  name varchar(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
