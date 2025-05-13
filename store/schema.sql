-- Create ENUM type for staff role
-- CREATE TYPE IF NOT EXISTS role AS ENUM ('doctor', 'receptionist');
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
        CREATE TYPE role AS ENUM ('doctor', 'receptionist');
    END IF;
END $$;

-- Create ENUM type for patient's gender
-- CREATE TYPE IF NOT EXISTS gender AS ENUM ('male', 'female', 'other');
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'gender') THEN
        CREATE TYPE gender AS ENUM ('male', 'female', 'other');
    END IF;
END $$;

-- DROP TABLE IF EXISTS doctor;

-- Create table doctor
CREATE TABLE IF NOT EXISTS doctor (
    doctor_id UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    role   role NOT NULL DEFAULT 'doctor',
    specialization TEXT NULL,
    password_hash TEXT NOT NULL UNIQUE,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger to update updated_at column for doctor table
CREATE OR REPLACE FUNCTION doctor_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at := CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to updated_at column for doctor table
CREATE TRIGGER trigger_doctor_set_updated_at
BEFORE UPDATE ON doctor
FOR EACH ROW
EXECUTE FUNCTION doctor_set_updated_at();

-- DROP TABLE IF EXISTS staff;

-- Create table staff
CREATE TABLE IF NOT EXISTS staff (
    staff_id UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    role   role NOT NULL DEFAULT 'receptionist',
    password_hash TEXT NOT NULL UNIQUE,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger to update updated_at column for staff table
CREATE OR REPLACE FUNCTION staff_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at := CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to updated_at column for staff table
CREATE TRIGGER trigger_staff_set_updated_at
BEFORE UPDATE ON staff
FOR EACH ROW
EXECUTE FUNCTION staff_set_updated_at();

-- Create table patient
CREATE TABLE IF NOT EXISTS patient (
    patient_id UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    gender gender NOT NULL,
    age INT NOT NULL,
    contact VARCHAR(10) NOT NULL,
    symptoms TEXT NULL,
    assigned_to UUID NOT NULL,
    created_by UUID NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_assigned_to FOREIGN KEY (assigned_to) REFERENCES doctor(doctor_id) ON DELETE CASCADE,
    CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES staff(staff_id) ON DELETE CASCADE
);

-- Create trigger to update updated_at column for staff table
CREATE OR REPLACE FUNCTION set_patient_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at := CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to updated_at column for staff table
CREATE TRIGGER trigger_set_patient_updated_at
BEFORE UPDATE ON patient
FOR EACH ROW
EXECUTE FUNCTION set_patient_updated_at();


-- Clear existing data before inserting new data
TRUNCATE TABLE doctor CASCADE;
TRUNCATE TABLE staff CASCADE;
TRUNCATE TABLE patient CASCADE;

-- Insert data into the doctor table

INSERT INTO doctor (doctor_id, fullname, email, specialization, password_hash, updated_at, created_at) 
VALUES ('28844ae6-d482-441b-abd6-09db2a64c707', 'Harshit Raj', 'harshitraj@medi.go', 'general physician', '$2a$10$w44vZYNdhRQy/ccGFNxNSuLT0xZokcRk/so6EOI5UMyyEhXbEasLK', '2025-05-13 11:16:06.174262', '2025-05-13 11:16:06.174262'),
-- harshit@medigo
('84576c8c-9e88-494d-bdd9-e855247e11df', 'Raj Sinha', 'rajsinha@medi.go', 'geriatrics', '$2a$10$EuzlK9mLYrvPEJl1pAujWO86eOVi1dlQ92YPXrCtkqUrQglA2uPGK', '2025-05-13 11:16:06.174262', '2025-05-13 11:16:06.174262'),
-- rajsinha@medigo
('e58056e6-28e1-43de-afda-8c6e9363ddda', 'Lucy Mountain', 'mountain.lucy@medi.go', 'pediatrics', '$2a$10$I/SpcJ3SJddEccIujxo7GuuqDL7xYjPrlbRmkPCZrMCnWfBbSx2UC', '2025-05-13 11:16:06.174262', '2025-05-13 11:16:06.174262');
-- mountain.lucy@medigo


-- Insert data into the staff table

INSERT INTO staff (staff_id, fullname, email, password_hash, updated_at, created_at) 
VALUES ('f4a9c66b-8e38-419b-93c4-215d5cefb318', 'Kunal Kumar', 'kumarkunal@medi.go', '$2a$10$PK5bsZmREcQQzUjDMDerjedK5WnDqkEn65.qxQBKMkEa.gspzewOy', '2025-05-13 11:16:06.174262', '2025-05-13 11:16:06.174262'),
-- kumarkunal@medigo
('9746be12-07b7-42a3-b8ab-7d1f209b63d7', 'Priya Patel', 'priya@medi.go', '$2a$10$rKPPL4QzONHtY3sFxPS3.Oq5M/I.dDVZAClXeGptfLuTw59LxPvCu', '2025-05-13 11:16:06.174262', '2025-05-13 11:16:06.174262');
-- priya@medigo


-- Insert data into the patient table

INSERT INTO patient (patient_id, fullname, gender, age, contact, symptoms, assigned_to, created_by, updated_at, created_at) 
VALUES ('cc2c2a7d-2e21-4f59-b7b8-bd9e5e4cf04c', 'Ananya Desai', 'female', 8, '7894561238', 'Ananya Desai, an 8-year-old female, is experiencing a dry cough that has persisted for about a week and is not going away. According to her mother, there is no fever, but Ananya occasionally wheezes after physical activity and has been waking up at night due to the coughing. The mother is concerned about the ongoing symptoms and is seeking a pediatric evaluation', 'e58056e6-28e1-43de-afda-8c6e9363ddda', 'f4a9c66b-8e38-419b-93c4-215d5cefb318', '2025-05-13 11:16:06.174262','2025-05-13 11:16:06.174262'),

('404784eb-ba77-4f60-94ea-4a170be9fd7e', 'Aiden Scott', 'male', 27, '7412589635', 'Aiden Scott, a 27-year-old male, reports experiencing persistent headaches over the past two weeks. He describes the pain as a dull ache that starts in the temples and sometimes radiates to the back of the head. The headaches tend to worsen in the late afternoon, especially after prolonged screen time. He has also mentioned occasional blurred vision and difficulty concentrating. He is seeking a consultation with a general physician to determine the cause.', '28844ae6-d482-441b-abd6-09db2a64c707', 'f4a9c66b-8e38-419b-93c4-215d5cefb318', '2025-05-13 11:16:06.174262','2025-05-13 11:16:06.174262'),

('af2742da-95ef-4629-b189-2ec59ce24f90', 'Meera Nair', 'female', 32, '9632587417', 'Meera Nair, a 32-year-old female, reports experiencing ongoing fatigue and mild shortness of breath over the past three weeks. She notes that even routine tasks like climbing stairs leave her feeling unusually tired. She has also mentioned occasional lightheadedness and a general lack of energy throughout the day. She is requesting an appointment with a general physician to investigate the cause of these symptoms.', '28844ae6-d482-441b-abd6-09db2a64c707', 'f4a9c66b-8e38-419b-93c4-215d5cefb318', '2025-05-13 11:16:06.174262','2025-05-13 11:16:06.174262'),

('367a97e2-d7ab-4981-9164-947cd872028d', 'Devansh Kapoor', 'male', 31, '9988774455', 'Devansh Kapoor, a 31-year-old male, has been experiencing intermittent stomach discomfort and bloating for the past month. He reports that the symptoms often occur after meals, especially heavier ones, and are sometimes accompanied by mild nausea. He has also noticed occasional changes in bowel habits. Devansh is seeking a consultation with a general physician to evaluate the cause and get relief.', '28844ae6-d482-441b-abd6-09db2a64c707', 'f4a9c66b-8e38-419b-93c4-215d5cefb318', '2025-05-13 11:16:06.174262','2025-05-13 11:16:06.174262'),

('1501120c-5f2e-4c83-9de3-01be35edbb5f', 'Ava Wilson', 'female', 45, '3625147894', 'Ava Wilson, a 45-year-old female, reports experiencing frequent episodes of heartburn and acid reflux over the past several weeks. She notes that the discomfort usually worsens after eating spicy or fatty foods and is often more noticeable at night when lying down. She occasionally feels a burning sensation in her chest and a sour taste in her mouth. Ava is requesting an appointment with a general physician to discuss her symptoms and explore possible treatment options.', '28844ae6-d482-441b-abd6-09db2a64c707', '9746be12-07b7-42a3-b8ab-7d1f209b63d7', '2025-05-13 11:16:06.174262','2025-05-13 11:16:06.174262'),

('af188a46-236a-40e5-8186-edcdf6e34d9b', 'Benjamin Carter', 'male', 76, '9848751236', 'Benjamin Carter, a 76-year-old male, has been experiencing increasing joint pain and stiffness, particularly in his knees and lower back, over the past few months. He reports difficulty with mobility, especially in the mornings, and occasional swelling in his joints after prolonged sitting or walking. He also mentions feeling more fatigued than usual and experiencing trouble sleeping due to discomfort. Benjamin is seeking an evaluation with a geriatrics specialist to manage his symptoms and improve his quality of life.', '84576c8c-9e88-494d-bdd9-e855247e11df', '9746be12-07b7-42a3-b8ab-7d1f209b63d7', '2025-05-13 11:16:06.174262','2025-05-13 11:16:06.174262');
