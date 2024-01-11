DB_NAME := 'lighthouse'

# connect to the dev db
@connect-dev:
    mariadb {{ DB_NAME }}

# drop the dev db
@drop-dev:
    echo 'drop database if exists {{ DB_NAME }}' | mariadb

# connect to the live db
@connect-live:
    ssh-to s -t mariadb {{ DB_NAME }}

# pull down the live db
@pull-live:
    importdb -f --host s --local {{ DB_NAME }} --remote {{ DB_NAME }}

