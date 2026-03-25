INSERT INTO users (id, email, role, password_hash, created_at) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'admin@company.com', 'admin', '$2a$10$dummyhash_admin_placeholder', NOW()),
    ('550e8400-e29b-41d4-a716-446655440002', 'john.doe@company.com', 'user', '$2a$10$dummyhash_user_placeholder', NOW()),
    ('550e8400-e29b-41d4-a716-446655440003', 'jane.smith@company.com', 'user', '$2a$10$dummyhash_user2_placeholder', NOW())
ON CONFLICT (email) DO NOTHING;


INSERT INTO rooms (id, name, description, capacity, created_at) VALUES
    ('6ba7b810-9dad-11d1-80b4-00c04fd430c1', 'Room Alpha', 'Small room for 1-on-1 meetings, equipped with whiteboard', 4, NOW()),
    ('6ba7b811-9dad-11d1-80b4-00c04fd430c1', 'Room Beta', 'Medium conference room with video conferencing setup', 8, NOW()),
    ('6ba7b812-9dad-11d1-80b4-00c04fd430c1', 'Room Gamma', 'Large hall for team events and presentations', 20, NOW())
ON CONFLICT DO NOTHING;


INSERT INTO schedules (id, room_id, start_time, end_time, created_at) VALUES
    ('7c9e6679-7425-40de-944b-e07fc1f90ae1', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '09:00', '18:00', NOW()),
    ('7c9e667a-7425-40de-944b-e07fc1f90ae1', '6ba7b811-9dad-11d1-80b4-00c04fd430c1', '10:00', '19:00', NOW()),
    ('7c9e667b-7425-40de-944b-e07fc1f90ae1', '6ba7b812-9dad-11d1-80b4-00c04fd430c1', '08:00', '20:00', NOW())
ON CONFLICT (room_id) DO NOTHING;


INSERT INTO schedule_days (schedule_id, day_of_week) VALUES
    -- Room Alpha: Mon-Fri
    ('7c9e6679-7425-40de-944b-e07fc1f90ae1', 1),
    ('7c9e6679-7425-40de-944b-e07fc1f90ae1', 2),
    ('7c9e6679-7425-40de-944b-e07fc1f90ae1', 3),
    ('7c9e6679-7425-40de-944b-e07fc1f90ae1', 4),
    ('7c9e6679-7425-40de-944b-e07fc1f90ae1', 5),
    -- Room Beta: Mon-Fri
    ('7c9e667a-7425-40de-944b-e07fc1f90ae1', 1),
    ('7c9e667a-7425-40de-944b-e07fc1f90ae1', 2),
    ('7c9e667a-7425-40de-944b-e07fc1f90ae1', 3),
    ('7c9e667a-7425-40de-944b-e07fc1f90ae1', 4),
    ('7c9e667a-7425-40de-944b-e07fc1f90ae1', 5),
    -- Room Gamma: Mon-Sat
    ('7c9e667b-7425-40de-944b-e07fc1f90ae1', 1),
    ('7c9e667b-7425-40de-944b-e07fc1f90ae1', 2),
    ('7c9e667b-7425-40de-944b-e07fc1f90ae1', 3),
    ('7c9e667b-7425-40de-944b-e07fc1f90ae1', 4),
    ('7c9e667b-7425-40de-944b-e07fc1f90ae1', 5),
    ('7c9e667b-7425-40de-944b-e07fc1f90ae1', 6)
ON CONFLICT DO NOTHING;


INSERT INTO slots (id, room_id, start_time, end_time) VALUES
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b01', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 09:00:00+00', '2026-04-01 09:30:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b02', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 09:30:00+00', '2026-04-01 10:00:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b03', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 10:00:00+00', '2026-04-01 10:30:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b04', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 10:30:00+00', '2026-04-01 11:00:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b05', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 11:00:00+00', '2026-04-01 11:30:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b06', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 11:30:00+00', '2026-04-01 12:00:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b07', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 12:00:00+00', '2026-04-01 12:30:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b08', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 12:30:00+00', '2026-04-01 13:00:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b09', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 13:00:00+00', '2026-04-01 13:30:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b0a', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 13:30:00+00', '2026-04-01 14:00:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b0b', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 14:00:00+00', '2026-04-01 14:30:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b0c', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 14:30:00+00', '2026-04-01 15:00:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b0d', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 15:00:00+00', '2026-04-01 15:30:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b0e', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 15:30:00+00', '2026-04-01 16:00:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b0f', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 16:00:00+00', '2026-04-01 16:30:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b10', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 16:30:00+00', '2026-04-01 17:00:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b11', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 17:00:00+00', '2026-04-01 17:30:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b12', '6ba7b810-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 17:30:00+00', '2026-04-01 18:00:00+00')
ON CONFLICT (room_id, start_time, end_time) DO NOTHING;

INSERT INTO slots (id, room_id, start_time, end_time) VALUES
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4c01', '6ba7b811-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 10:00:00+00', '2026-04-01 10:30:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4c02', '6ba7b811-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 10:30:00+00', '2026-04-01 11:00:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4c03', '6ba7b811-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 14:00:00+00', '2026-04-01 14:30:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4c04', '6ba7b811-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 14:30:00+00', '2026-04-01 15:00:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4c05', '6ba7b811-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 18:30:00+00', '2026-04-01 19:00:00+00')
ON CONFLICT (room_id, start_time, end_time) DO NOTHING;

INSERT INTO slots (id, room_id, start_time, end_time) VALUES
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4d01', '6ba7b812-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 08:00:00+00', '2026-04-01 08:30:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4d02', '6ba7b812-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 12:00:00+00', '2026-04-01 12:30:00+00'),
    ('8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4d03', '6ba7b812-9dad-11d1-80b4-00c04fd430c1', '2026-04-01 19:30:00+00', '2026-04-01 20:00:00+00')
ON CONFLICT (room_id, start_time, end_time) DO NOTHING;


INSERT INTO bookings (id, slot_id, user_id, status, conference_link, created_at) VALUES
    ('9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c01', 
     '8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b03', 
     '550e8400-e29b-41d4-a716-446655440002', 
     'active', 
     'https://meet.company.com/abc123', 
     NOW()),
    
    ('9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c02', 
     '8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4d02', 
     '550e8400-e29b-41d4-a716-446655440003', 
     'active', 
     'https://meet.company.com/xyz789', 
     NOW()),
    
    ('9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c03', 
     '8f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4c03', 
     '550e8400-e29b-41d4-a716-446655440002', 
     'cancelled', 
     NULL, 
     NOW())
ON CONFLICT DO NOTHING;
