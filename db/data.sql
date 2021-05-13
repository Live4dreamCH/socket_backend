use chat;

insert into users(u_id, psw, u_name, has_set_fmi)
values 
(1, 'mhq', 'mhq', 0),
(2, 'lch', 'lch', 0);

insert into convs(conv_id, is_group)
values (2, 0);

insert into conv_members(conv_id, mem_id)
values (2, 1),(2, 2);