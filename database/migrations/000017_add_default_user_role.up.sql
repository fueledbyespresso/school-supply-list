INSERT INTO role (role_id, role_name, role_desc)
VALUES (1, 'default', default);
INSERT INTO role_resource_bridge (rrb_id, can_add, can_view, can_edit, can_delete, resource_id, role_id)
VALUES (default, false, true, false, false, 1, 1);
INSERT INTO role_resource_bridge (rrb_id, can_add, can_view, can_edit, can_delete, resource_id, role_id)
VALUES (default, false, true, false, false, 2, 1);
INSERT INTO role_resource_bridge (rrb_id, can_add, can_view, can_edit, can_delete, resource_id, role_id)
VALUES (default, false, true, false, false, 3, 1);
INSERT INTO role_resource_bridge (rrb_id, can_add, can_view, can_edit, can_delete, resource_id, role_id)
VALUES (default, false, false, true, false, 4, 1);
INSERT INTO role_resource_bridge (rrb_id, can_add, can_view, can_edit, can_delete, resource_id, role_id)
VALUES (default, false, true, false, false, 5, 1);