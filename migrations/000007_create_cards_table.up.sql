CREATE TABLE cards (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  category VARCHAR(255) NOT NULL,
  title VARCHAR(255) NOT NULL,
  body TEXT NOT NULL,
  type VARCHAR(255) NOT NULL,
  cost INTEGER,
  down_payment INTEGER,
  mortgage INTEGER,
  cash_flow INTEGER,
  trade_range_down INTEGER,
  trade_range_up INTEGER,
  inflation INTEGER,
  quantity INTEGER,
  active BOOLEAN NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
