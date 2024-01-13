DB_NAME := 'lighthouse'

# connect to the dev db
@dev:
    mariadb {{ DB_NAME }}

# drop the dev db
@drop:
    echo 'drop database if exists {{ DB_NAME }}' | mariadb

# connect to the live db
@live:
    ssh-to s -t mariadb {{ DB_NAME }}

# pull down the live db
@pull:
    importdb -f --host s --local {{ DB_NAME }} --remote {{ DB_NAME }}
