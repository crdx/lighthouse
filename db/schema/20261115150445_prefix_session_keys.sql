update sessions set k = concat('session:', k) where k not like 'session:%';
