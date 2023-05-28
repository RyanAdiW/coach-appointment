CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE coach_availability (
  id SERIAL PRIMARY KEY,
  user_id VARCHAR NULL,
  coach_name VARCHAR(255),
  timezone VARCHAR(255),
  day_of_week VARCHAR(255),
  available_at TIME WITH TIME ZONE,
  available_until TIME WITH TIME ZONE
);

INSERT INTO coach_availability (coach_name, timezone, day_of_week, available_at, available_until)
VALUES
  ('Christy Schumm', '(GMT-06:00) America/North_Dakota/New_Salem', 'Monday', '09:00-06:00', '17:30-06:00'),
  ('Christy Schumm', '(GMT-06:00) America/North_Dakota/New_Salem', 'Tuesday', '08:00-06:00', '16:00-06:00'),
  ('Christy Schumm', '(GMT-06:00) America/North_Dakota/New_Salem', 'Thursday', '09:00-06:00', '16:00-06:00'),
  ('Christy Schumm', '(GMT-06:00) America/North_Dakota/New_Salem', 'Friday', '07:00-06:00', '14:00-06:00'),
  ('Natalia Stanton Jr.', '(GMT-06:00) Central Time (US & Canada)', 'Tuesday', '08:00-06:00', '10:00-06:00'),
  ('Natalia Stanton Jr.', '(GMT-06:00) Central Time (US & Canada)', 'Wednesday', '11:00-06:00', '18:00-06:00'),
  ('Natalia Stanton Jr.', '(GMT-06:00) Central Time (US & Canada)', 'Saturday', '09:00-06:00', '15:00-06:00'),
  ('Natalia Stanton Jr.', '(GMT-06:00) Central Time (US & Canada)', 'Sunday', '08:00-06:00', '15:00-06:00'),
  ('Nola Murazik V', '(GMT-09:00) America/Yakutat', 'Monday', '08:00-09:00', '10:00-09:00'),
  ('Nola Murazik V', '(GMT-09:00) America/Yakutat', 'Tuesday', '11:00-09:00', '13:00-09:00'),
  ('Nola Murazik V', '(GMT-09:00) America/Yakutat', 'Wednesday', '08:00-09:00', '10:00-09:00'),
  ('Nola Murazik V', '(GMT-09:00) America/Yakutat', 'Saturday', '08:00-09:00', '11:00-09:00'),
  ('Nola Murazik V', '(GMT-09:00) America/Yakutat', 'Sunday', '07:00-09:00', '09:00-09:00'),
  ('Elyssa O''Kon', '(GMT-06:00) Central Time (US & Canada)', 'Monday', '09:00-06:00', '15:00-06:00'),
  ('Elyssa O''Kon', '(GMT-06:00) Central Time (US & Canada)', 'Tuesday', '06:00-06:00', '13:00-06:00'),
  ('Elyssa O''Kon', '(GMT-06:00) Central Time (US & Canada)', 'Wednesday', '06:00-06:00', '11:00-06:00'),
  ('Elyssa O''Kon', '(GMT-06:00) Central Time (US & Canada)', 'Friday', '08:00-06:00', '12:00-06:00'),
  ('Elyssa O''Kon', '(GMT-06:00) Central Time (US & Canada)', 'Saturday', '09:00-06:00', '16:00-06:00'),
  ('Elyssa O''Kon', '(GMT-06:00) Central Time (US & Canada)', 'Sunday', '08:00-06:00', '10:00-06:00'),
  ('Dr. Geovany Keebler', '(GMT-06:00) Central Time (US & Canada)', 'Thursday', '07:00-06:00', '14:00-06:00'),
  ('Dr. Geovany Keebler', '(GMT-06:00) Central Time (US & Canada)', 'Thursday', '15:00-06:00', '17:00-06:00');
