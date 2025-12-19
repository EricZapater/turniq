DO $$
DECLARE
    v_customer_id UUID;
    v_sf_prod UUID;
    v_sf_assembly UUID;
    v_wc_lathe UUID;
    v_wc_drill UUID;
    v_wc_assembly UUID;
    v_user_id UUID;
    v_op_john UUID;
    v_op_jane UUID;
    v_op_pere UUID;
    v_shift_morn UUID;
    v_shift_after UUID;
    v_job1 UUID;
    v_job2 UUID;
    v_job3 UUID;
    v_job4 UUID;
    v_job5 UUID;
    v_job6 UUID;
    v_job7 UUID;
    v_job8 UUID;
BEGIN
    -- 1. Create Customer
    INSERT INTO customers (name, email, language, status, max_operators, max_workcenters, max_shop_floors, max_users, max_jobs)
    VALUES ('Test Customer', 'info@testcustomer.com', 'ca', 'active', 10, 10, 5, 5, 20)
    RETURNING id INTO v_customer_id;

    -- 2. Create User (Password: t3st_2026)
    INSERT INTO users (customer_id, email, password, username, is_admin, is_active)
    VALUES (v_customer_id, 'testuser@testcustomer.com', '$2a$10$Jad8aEphJiU3Ps.1aEqXb.KgkhkNjjOMhwWw.N8i18Pc1MIo1wP76', 'Test Admin', true, true)
    RETURNING id INTO v_user_id;

    -- 3. Create Shopfloors
    INSERT INTO shopfloors (customer_id, name) VALUES (v_customer_id, 'Producció Principal') RETURNING id INTO v_sf_prod;
    INSERT INTO shopfloors (customer_id, name) VALUES (v_customer_id, 'Muntatge') RETURNING id INTO v_sf_assembly;

    -- 4. Create Workcenters
    INSERT INTO workcenters (customer_id, shop_floor_id, name, is_active) VALUES (v_customer_id, v_sf_prod, 'Torn CNC 1', true) RETURNING id INTO v_wc_lathe;
    INSERT INTO workcenters (customer_id, shop_floor_id, name, is_active) VALUES (v_customer_id, v_sf_prod, 'Fresadora 1', true) RETURNING id INTO v_wc_drill;
    INSERT INTO workcenters (customer_id, shop_floor_id, name, is_active) VALUES (v_customer_id, v_sf_assembly, 'Banc de Muntatge A', true) RETURNING id INTO v_wc_assembly;

    -- 5. Create Shifts (Set to Today)
    INSERT INTO shifts (customer_id, shopfloor_id, name, color, start_time, end_time, is_active)
    VALUES (v_customer_id, v_sf_prod, 'Matí', '#FFD700', CURRENT_DATE + TIME '06:00:00', CURRENT_DATE + TIME '14:00:00', true)
    RETURNING id INTO v_shift_morn;

    INSERT INTO shifts (customer_id, shopfloor_id, name, color, start_time, end_time, is_active)
    VALUES (v_customer_id, v_sf_prod, 'Tarda', '#FFA500', CURRENT_DATE + TIME '14:00:00', CURRENT_DATE + TIME '22:00:00', true)
    RETURNING id INTO v_shift_after;

    -- 6. Create Operators
    INSERT INTO operators (customer_id, shop_floor_id, code, name, surname, is_active)
    VALUES (v_customer_id, v_sf_prod, 'OP001', 'Joan', 'Garcia', true)
    RETURNING id INTO v_op_john;

    INSERT INTO operators (customer_id, shop_floor_id, code, name, surname, is_active)
    VALUES (v_customer_id, v_sf_prod, 'OP002', 'Maria', 'Martínez', true)
    RETURNING id INTO v_op_jane;

    INSERT INTO operators (customer_id, shop_floor_id, code, name, surname, is_active)
    VALUES (v_customer_id, v_sf_assembly, 'OP003', 'Pere', 'Vila', true)
    RETURNING id INTO v_op_pere;

    -- 7. Create Jobs (Ampliats)
    INSERT INTO jobs (customer_id, shop_floor_id, workcenter_id, job_code, product_code, description, estimated_duration)
    VALUES (v_customer_id, v_sf_prod, v_wc_lathe, 'JOB-1001', 'P-X100', 'Mecanitzat Eix Principal', 120)
    RETURNING id INTO v_job1;

    INSERT INTO jobs (customer_id, shop_floor_id, workcenter_id, job_code, product_code, description, estimated_duration)
    VALUES (v_customer_id, v_sf_prod, v_wc_drill, 'JOB-1002', 'P-Y200', 'Forat Base', 60)
    RETURNING id INTO v_job2;

    INSERT INTO jobs (customer_id, shop_floor_id, workcenter_id, job_code, product_code, description, estimated_duration)
    VALUES (v_customer_id, v_sf_prod, v_wc_lathe, 'JOB-1003', 'P-Z300', 'Torn Bushings', 90)
    RETURNING id INTO v_job3;

    INSERT INTO jobs (customer_id, shop_floor_id, workcenter_id, job_code, product_code, description, estimated_duration)
    VALUES (v_customer_id, v_sf_prod, v_wc_drill, 'JOB-1004', 'P-A400', 'Fresatge Planxa', 180)
    RETURNING id INTO v_job4;

    INSERT INTO jobs (customer_id, shop_floor_id, workcenter_id, job_code, product_code, description, estimated_duration)
    VALUES (v_customer_id, v_sf_assembly, v_wc_assembly, 'JOB-2001', 'ASM-100', 'Muntatge Conjunt A', 240)
    RETURNING id INTO v_job5;

    INSERT INTO jobs (customer_id, shop_floor_id, workcenter_id, job_code, product_code, description, estimated_duration)
    VALUES (v_customer_id, v_sf_assembly, v_wc_assembly, 'JOB-2002', 'ASM-200', 'Muntatge Conjunt B', 150)
    RETURNING id INTO v_job6;

    INSERT INTO jobs (customer_id, shop_floor_id, workcenter_id, job_code, product_code, description, estimated_duration)
    VALUES (v_customer_id, v_sf_prod, v_wc_lathe, 'JOB-1005', 'P-B500', 'Eix Secundari', 100)
    RETURNING id INTO v_job7;

    INSERT INTO jobs (customer_id, shop_floor_id, workcenter_id, job_code, product_code, description, estimated_duration)
    VALUES (v_customer_id, v_sf_prod, v_wc_drill, 'JOB-1006', 'P-C600', 'Tapa Metàl·lica', 45)
    RETURNING id INTO v_job8;

    -- 8. Create Schedule Entries (Ampliats amb més diversitat)
    
    -- AVUI - Torn Matí
    INSERT INTO schedule_entries (customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id, date, start_time, end_time, is_completed)
    VALUES (v_customer_id, v_sf_prod, v_shift_morn, v_wc_lathe, v_job1, v_op_john, CURRENT_DATE, CURRENT_DATE + TIME '06:00:00', CURRENT_DATE + TIME '08:00:00', false);

    INSERT INTO schedule_entries (customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id, date, start_time, end_time, is_completed)
    VALUES (v_customer_id, v_sf_prod, v_shift_morn, v_wc_lathe, v_job3, v_op_john, CURRENT_DATE, CURRENT_DATE + TIME '08:30:00', CURRENT_DATE + TIME '10:00:00', false);

    -- AVUI - Fresadora Matí
    INSERT INTO schedule_entries (customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id, date, start_time, end_time, is_completed)
    VALUES (v_customer_id, v_sf_prod, v_shift_morn, v_wc_drill, v_job2, v_op_jane, CURRENT_DATE, CURRENT_DATE + TIME '06:00:00', CURRENT_DATE + TIME '07:00:00', false);

    INSERT INTO schedule_entries (customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id, date, start_time, end_time, is_completed)
    VALUES (v_customer_id, v_sf_prod, v_shift_morn, v_wc_drill, v_job8, v_op_jane, CURRENT_DATE, CURRENT_DATE + TIME '07:15:00', CURRENT_DATE + TIME '08:00:00', false);

    -- AVUI - Torn Tarda
    INSERT INTO schedule_entries (customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id, date, start_time, end_time, is_completed)
    VALUES (v_customer_id, v_sf_prod, v_shift_after, v_wc_lathe, v_job7, v_op_john, CURRENT_DATE, CURRENT_DATE + TIME '14:00:00', CURRENT_DATE + TIME '15:40:00', false);

    -- AVUI - Fresadora Tarda
    INSERT INTO schedule_entries (customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id, date, start_time, end_time, is_completed)
    VALUES (v_customer_id, v_sf_prod, v_shift_after, v_wc_drill, v_job4, v_op_jane, CURRENT_DATE, CURRENT_DATE + TIME '14:00:00', CURRENT_DATE + TIME '17:00:00', false);

    -- AVUI - Muntatge (sense horari específic)
    INSERT INTO schedule_entries (customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id, date, "order")
    VALUES (v_customer_id, v_sf_assembly, v_shift_morn, v_wc_assembly, v_job5, v_op_pere, CURRENT_DATE, 1);

    INSERT INTO schedule_entries (customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id, date, "order")
    VALUES (v_customer_id, v_sf_assembly, v_shift_morn, v_wc_assembly, v_job6, v_op_pere, CURRENT_DATE, 2);

    -- AHIR - Completades
    INSERT INTO schedule_entries (customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id, date, start_time, end_time, is_completed)
    VALUES (v_customer_id, v_sf_prod, v_shift_morn, v_wc_lathe, v_job1, v_op_john, CURRENT_DATE - INTERVAL '1 day', CURRENT_DATE - INTERVAL '1 day' + TIME '06:00:00', CURRENT_DATE - INTERVAL '1 day' + TIME '08:00:00', true);

    INSERT INTO schedule_entries (customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id, date, start_time, end_time, is_completed)
    VALUES (v_customer_id, v_sf_prod, v_shift_morn, v_wc_drill, v_job2, v_op_jane, CURRENT_DATE - INTERVAL '1 day', CURRENT_DATE - INTERVAL '1 day' + TIME '06:00:00', CURRENT_DATE - INTERVAL '1 day' + TIME '07:00:00', true);

    -- DEMÀ - Planificades
    INSERT INTO schedule_entries (customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id, date, start_time, end_time, is_completed)
    VALUES (v_customer_id, v_sf_prod, v_shift_morn, v_wc_lathe, v_job3, v_op_john, CURRENT_DATE + INTERVAL '1 day', CURRENT_DATE + INTERVAL '1 day' + TIME '06:00:00', CURRENT_DATE + INTERVAL '1 day' + TIME '07:30:00', false);

    INSERT INTO schedule_entries (customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id, date, start_time, end_time, is_completed)
    VALUES (v_customer_id, v_sf_prod, v_shift_morn, v_wc_drill, v_job4, v_op_jane, CURRENT_DATE + INTERVAL '1 day', CURRENT_DATE + INTERVAL '1 day' + TIME '06:00:00', CURRENT_DATE + INTERVAL '1 day' + TIME '09:00:00', false);

    -- 9. Create Payments (Billing Report)
    INSERT INTO payments (customer_id, amount, currency, payment_method, status, due_date, paid_at)
    VALUES (v_customer_id, 150.00, 'EUR', 'Visa', 'paid', CURRENT_DATE - INTERVAL '5 days', CURRENT_DATE - INTERVAL '4 days');

    INSERT INTO payments (customer_id, amount, currency, payment_method, status, due_date, paid_at)
    VALUES (v_customer_id, 200.50, 'EUR', 'Mastercard', 'pending', CURRENT_DATE + INTERVAL '10 days', NULL);

    -- 10. Create Time Entries (Ampliats)
    
    -- Joan (Torn) - Ahir
    INSERT INTO time_entries (operator_id, workcenter_id, check_in, check_out)
    VALUES (v_op_john, v_wc_lathe, CURRENT_DATE - INTERVAL '1 day' + TIME '06:00:00', CURRENT_DATE - INTERVAL '1 day' + TIME '14:00:00');

    -- Maria (Fresadora) - Ahir
    INSERT INTO time_entries (operator_id, workcenter_id, check_in, check_out)
    VALUES (v_op_jane, v_wc_drill, CURRENT_DATE - INTERVAL '1 day' + TIME '06:00:00', CURRENT_DATE - INTERVAL '1 day' + TIME '14:00:00');

    -- Pere (Muntatge) - Ahir
    INSERT INTO time_entries (operator_id, workcenter_id, check_in, check_out)
    VALUES (v_op_pere, v_wc_assembly, CURRENT_DATE - INTERVAL '1 day' + TIME '06:00:00', CURRENT_DATE - INTERVAL '1 day' + TIME '14:00:00');

    -- Joan (Torn) - Fa 2 dies
    INSERT INTO time_entries (operator_id, workcenter_id, check_in, check_out)
    VALUES (v_op_john, v_wc_lathe, CURRENT_DATE - INTERVAL '2 days' + TIME '06:00:00', CURRENT_DATE - INTERVAL '2 days' + TIME '14:00:00');

    -- Maria (Fresadora) - Fa 2 dies
    INSERT INTO time_entries (operator_id, workcenter_id, check_in, check_out)
    VALUES (v_op_jane, v_wc_drill, CURRENT_DATE - INTERVAL '2 days' + TIME '06:00:00', CURRENT_DATE - INTERVAL '2 days' + TIME '13:30:00');

    -- Joan (Torn) - Fa 3 dies - Torn tarda
    INSERT INTO time_entries (operator_id, workcenter_id, check_in, check_out)
    VALUES (v_op_john, v_wc_lathe, CURRENT_DATE - INTERVAL '3 days' + TIME '14:00:00', CURRENT_DATE - INTERVAL '3 days' + TIME '22:00:00');

    -- Maria (Fresadora) - Fa 3 dies
    INSERT INTO time_entries (operator_id, workcenter_id, check_in, check_out)
    VALUES (v_op_jane, v_wc_drill, CURRENT_DATE - INTERVAL '3 days' + TIME '06:00:00', CURRENT_DATE - INTERVAL '3 days' + TIME '14:00:00');

    -- Pere (Muntatge) - Fa 3 dies
    INSERT INTO time_entries (operator_id, workcenter_id, check_in, check_out)
    VALUES (v_op_pere, v_wc_assembly, CURRENT_DATE - INTERVAL '3 days' + TIME '06:00:00', CURRENT_DATE - INTERVAL '3 days' + TIME '14:00:00');

    -- Joan (Torn) - Fa 4 dies
    INSERT INTO time_entries (operator_id, workcenter_id, check_in, check_out)
    VALUES (v_op_john, v_wc_lathe, CURRENT_DATE - INTERVAL '4 days' + TIME '06:00:00', CURRENT_DATE - INTERVAL '4 days' + TIME '14:00:00');

    -- Maria (Fresadora) - Fa 4 dies - Mig dia (sortida anticipada)
    INSERT INTO time_entries (operator_id, workcenter_id, check_in, check_out)
    VALUES (v_op_jane, v_wc_drill, CURRENT_DATE - INTERVAL '4 days' + TIME '06:00:00', CURRENT_DATE - INTERVAL '4 days' + TIME '10:00:00');

    -- Joan (Torn) - Fa 5 dies - Torn tarda
    INSERT INTO time_entries (operator_id, workcenter_id, check_in, check_out)
    VALUES (v_op_john, v_wc_lathe, CURRENT_DATE - INTERVAL '5 days' + TIME '14:00:00', CURRENT_DATE - INTERVAL '5 days' + TIME '22:00:00');

    -- Pere (Muntatge) - Fa 5 dies
    INSERT INTO time_entries (operator_id, workcenter_id, check_in, check_out)
    VALUES (v_op_pere, v_wc_assembly, CURRENT_DATE - INTERVAL '5 days' + TIME '06:00:00', CURRENT_DATE - INTERVAL '5 days' + TIME '14:00:00');

END $$;