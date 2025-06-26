CREATE TABLE urls (
  id UUID PRIMARY KEY,
  original TEXT NOT NULL,
  short_code VARCHAR(10) UNIQUE NOT NULL,
  custom_code VARCHAR(20),
  domain TEXT,
  visibility VARCHAR(10) DEFAULT 'public',
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE visits (
  id UUID PRIMARY KEY,
  url_id UUID REFERENCES urls(id),
  timestamp TIMESTAMP,
  ip_address TEXT,
  country TEXT,
  browser TEXT,
  device TEXT
);
