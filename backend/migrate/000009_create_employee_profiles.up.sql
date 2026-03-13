CREATE TABLE employee_profiles (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    employee_code VARCHAR(100) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    
    corporate_email VARCHAR(255) UNIQUE,
    phone VARCHAR(20),
    telegram_handle VARCHAR(100),
    
    emergency_contact VARCHAR(255),
    
    position VARCHAR(255),
    department VARCHAR(255),
    
    birth_date DATE,
    avatar_url TEXT,
    hire_date DATE NOT NULL,
    dismissal_date DATE
);
