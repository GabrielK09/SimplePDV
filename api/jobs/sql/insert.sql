INSERT INTO users 
    (name, cpf, login, password, is_admin)

VALUES  
    ('Julcineia', '111.111.111-11', 'julci', $1, true),
    ('Gabriel', '222.222.222-22', 'gabi', $1, true)