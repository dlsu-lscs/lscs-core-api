-- Table: positions
CREATE TABLE positions (
    position_id VARCHAR(10) PRIMARY KEY,
    position_name VARCHAR(100) NOT NULL
);

-- Table: houses
CREATE TABLE houses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) UNIQUE,
    description TEXT
);

-- Table: docu_status
CREATE TABLE docu_status (
    id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);

-- Table: document_templates
CREATE TABLE document_templates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    doc_title VARCHAR(255) NOT NULL,
    doc_key VARCHAR(255) NOT NULL
);

-- Table: event_durations
CREATE TABLE event_durations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

-- Table: event_natures
CREATE TABLE event_natures (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

-- Table: event_types
CREATE TABLE event_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

-- Table: file_statuses
CREATE TABLE file_statuses (
    file_key VARCHAR(255) PRIMARY KEY,
    status ENUM('DRAFTING','DONE','REVISE','APPROVED') DEFAULT 'DRAFTING'
);

-- Table: fin_process_ref
CREATE TABLE fin_process_ref (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

-- Table: fin_status
CREATE TABLE fin_status (
    id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(200) NOT NULL
);

-- Table: progress_statuses
CREATE TABLE progress_statuses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    status VARCHAR(255) NOT NULL
);

-- Table: pub_req_status
CREATE TABLE pub_req_status (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- Table: terms
CREATE TABLE terms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    term INT NOT NULL,
    start_year INT NOT NULL,
    end_year INT NOT NULL
);

-- Table: members
CREATE TABLE members (
    id INT PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    nickname VARCHAR(100),
    email VARCHAR(255) UNIQUE NOT NULL,
    telegram VARCHAR(100),
    position_id VARCHAR(10),
    committee_id VARCHAR(10),
    college VARCHAR(255),
    program VARCHAR(255),
    discord VARCHAR(32),
    interests TEXT,
    contact_number VARCHAR(32),
    fb_link VARCHAR(255),
    house_id INT,
    FOREIGN KEY (position_id) REFERENCES positions(position_id) ON DELETE SET NULL,
    FOREIGN KEY (house_id) REFERENCES houses(id) ON DELETE SET NULL
);

-- Table: divisions
CREATE TABLE divisions (
    division_id VARCHAR(10) PRIMARY KEY,
    division_name VARCHAR(100) NOT NULL,
    division_head INT,
    FOREIGN KEY (division_head) REFERENCES members(id) ON DELETE SET NULL
);

-- Table: committees
CREATE TABLE committees (
    committee_id VARCHAR(10) PRIMARY KEY,
    committee_name VARCHAR(100) NOT NULL,
    committee_head INT,
    division_id VARCHAR(10),
    FOREIGN KEY (committee_head) REFERENCES members(id) ON DELETE SET NULL,
    FOREIGN KEY (division_id) REFERENCES divisions(division_id) ON DELETE SET NULL
);

ALTER TABLE members ADD CONSTRAINT fk_committee FOREIGN KEY (committee_id) REFERENCES committees(committee_id) ON DELETE SET NULL;

-- Table: api_keys
CREATE TABLE api_keys (
    api_key_id INT AUTO_INCREMENT PRIMARY KEY,
    member_email VARCHAR(100) NOT NULL,
    api_key_hash VARCHAR(255) NOT NULL,
    project VARCHAR(255) NOT NULL,
    allowed_origin VARCHAR(255) UNIQUE,
    is_dev BOOLEAN NOT NULL DEFAULT FALSE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    FOREIGN KEY (member_email) REFERENCES members(email) ON DELETE CASCADE
);

-- Table: events
CREATE TABLE events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    arn VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    committee_id VARCHAR(10) NOT NULL,
    type VARCHAR(100),
    nature_id INT NOT NULL,
    term_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    duration_id INT,
    brief_description TEXT,
    goals TEXT,
    objectives TEXT,
    strategies TEXT,
    measures TEXT,
    budget_allocation DECIMAL(10,2) DEFAULT 0.00,
    venue VARCHAR(255) DEFAULT 'Online',
    docu_head INT,
    fin_head INT,
    FOREIGN KEY (committee_id) REFERENCES committees(committee_id),
    FOREIGN KEY (nature_id) REFERENCES event_natures(id),
    FOREIGN KEY (term_id) REFERENCES terms(id),
    FOREIGN KEY (duration_id) REFERENCES event_durations(id),
    FOREIGN KEY (docu_head) REFERENCES members(id),
    FOREIGN KEY (fin_head) REFERENCES members(id)
);

-- Table: event_dates
CREATE TABLE event_dates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);

-- Table: event_docu_head
CREATE TABLE event_docu_head (
    event_id INT NOT NULL,
    member_id INT NOT NULL,
    PRIMARY KEY (event_id, member_id),
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE
);

-- Table: event_files
CREATE TABLE event_files (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT NOT NULL,
    file_key VARCHAR(255) NOT NULL,
    file_order INT NOT NULL,
    file_status INT NOT NULL DEFAULT 1,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);

-- Table: event_heads
CREATE TABLE event_heads (
    event_id INT NOT NULL,
    member_id INT NOT NULL,
    PRIMARY KEY (event_id, member_id),
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE
);

-- Table: event_participants
CREATE TABLE event_participants (
    event_id INT NOT NULL,
    student_id INT NOT NULL,
    student_name VARCHAR(255),
    student_email VARCHAR(255),
    notes VARCHAR(512),
    PRIMARY KEY (event_id, student_id),
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);

-- Table: event_trackers
CREATE TABLE event_trackers (
    event_id INT PRIMARY KEY,
    preacts_deadline DATETIME,
    preacts_status VARCHAR(100) NOT NULL DEFAULT 'INIT',
    postacts_deadline DATETIME,
    postacts_status VARCHAR(100) NOT NULL DEFAULT 'INIT',
    docu_drive_id VARCHAR(255),
    fin_drive_id VARCHAR(255),
    fin_preacts_deadline DATETIME,
    fin_preacts_status VARCHAR(100) DEFAULT 'INIT',
    fin_postacts_deadline DATETIME,
    fin_postacts_status VARCHAR(100) DEFAULT 'INIT',
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (fin_preacts_status) REFERENCES fin_status(id),
    FOREIGN KEY (fin_postacts_status) REFERENCES fin_status(id)
);

-- Table: fin_processes
CREATE TABLE fin_processes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tracker_id INT NOT NULL,
    process VARCHAR(50) NOT NULL,
    FOREIGN KEY (tracker_id) REFERENCES event_trackers(event_id),
    FOREIGN KEY (process) REFERENCES fin_process_ref(id)
);

-- Table: pub_requests
CREATE TABLE pub_requests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT,
    pub_head INT,
    pub_type TEXT,
    pub_drive_id VARCHAR(255) NOT NULL,
    pub_status VARCHAR(50),
    posting_date DATETIME NOT NULL,
    pub_details TEXT NOT NULL,
    pub_content TEXT NOT NULL,
    caption TEXT NOT NULL,
    opa_numbers VARCHAR(255) NOT NULL,
    for_posting BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL,
    requester_id INT,
    dimensions VARCHAR(255),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (requester_id) REFERENCES members(id)
);
