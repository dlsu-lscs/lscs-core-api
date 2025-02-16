-- Table: positions
CREATE TABLE positions (
    position_id VARCHAR(10) PRIMARY KEY,  -- Position codes like MEM, JO, CT, etc.
    position_name VARCHAR(100) NOT NULL
);

-- Table: divisions
CREATE TABLE divisions (
    division_id VARCHAR(10) PRIMARY KEY,  -- Division codes like CORE, INT, EXT, OPS
    division_name VARCHAR(100) NOT NULL,
    division_head INT,
    CONSTRAINT fk_division_head FOREIGN KEY (division_head) REFERENCES members(id) ON DELETE SET NULL
);

-- Table: members
CREATE TABLE members (
    id INT UNIQUE PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    nickname VARCHAR(100),
    email VARCHAR(255) UNIQUE NOT NULL,
    telegram VARCHAR(100),
    position_id VARCHAR(10),
    committee_id VARCHAR(10),
    college VARCHAR(100),
    program VARCHAR(100),
    discord VARCHAR(100),
    CONSTRAINT fk_position FOREIGN KEY (position_id) REFERENCES positions(position_id) ON DELETE SET NULL,
    CONSTRAINT fk_committee FOREIGN KEY (committee_id) REFERENCES committees(committee_id) ON DELETE SET NULL
);

-- Table: committees
CREATE TABLE committees (
    committee_id VARCHAR(10) PRIMARY KEY, 
    committee_name VARCHAR(100) NOT NULL,
    committee_head INT,
    division_id VARCHAR(10),
    CONSTRAINT fk_committee_head FOREIGN KEY (committee_head) REFERENCES members(id) ON DELETE SET NULL,
    CONSTRAINT fk_committee_division FOREIGN KEY (division_id) REFERENCES divisions(division_id) ON DELETE SET NULL
);

-- Table: api_keys
CREATE TABLE api_keys (
    api_key_id INT AUTO_INCREMENT PRIMARY KEY,
    member_email VARCHAR(100) NOT NULL,
    api_key_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT fk_api_keys_member_email FOREIGN KEY (member_email) REFERENCES members(email) ON DELETE CASCADE
);
