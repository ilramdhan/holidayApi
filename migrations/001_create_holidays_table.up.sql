CREATE TABLE IF NOT EXISTS holidays (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('national', 'collective_leave')),
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_holidays_date ON holidays(date);
CREATE INDEX idx_holidays_type ON holidays(type);
CREATE INDEX idx_holidays_is_active ON holidays(is_active);
CREATE INDEX idx_holidays_date_type ON holidays(date, type);

-- Create trigger to update updated_at timestamp
CREATE TRIGGER update_holidays_updated_at 
    AFTER UPDATE ON holidays
    FOR EACH ROW
BEGIN
    UPDATE holidays SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
